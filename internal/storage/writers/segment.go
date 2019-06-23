package writers

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"meerkat/internal/config"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"meerkat/internal/storage/text"
	"meerkat/internal/tools/encoding"
	"path/filepath"
)

const (
	idxExt = ".idx"
	binExt = ".bin"
	colExt = ".col"
	posExt = ".pos"
)

type SegmentWriter struct {
	log     zerolog.Logger
	path    string
	segment *inmem.Segment
}

func NewSegmentWriter(path string, segment *inmem.Segment) *SegmentWriter {

	log := log.With().
		Str("component", "SegmentWriter").
		Str("segmentID", segment.ID).
		Str("indexName", segment.IndexInfo.Name).
		Logger()

	if segment.State != inmem.InMem {
		log.Panic().Msgf("invalid segment state [%s]", segment.State)
	}

	sw := &SegmentWriter{
		segment: segment,
		path:    path,
		log:     log,
	}

	return sw

}

func (sw *SegmentWriter) Write() error {

	sw.log.Info().Msg("Starting to write segment")

	// Columns offsets.
	colOffset := make([][]*inmem.PageDescriptor, len(sw.segment.Columns))

	// TS column must be be processed first because it could
	// be sorted and and in this case it will determine the order
	// of the rest of segment columns.

	tsColumn := sw.segment.Columns[0].(*inmem.ColumnTimeStamp)

	var err error

	colOffset[0], err = sw.writeTSColumn(tsColumn)

	if err != nil {
		return err
	}

	for i, col := range sw.segment.Columns[1:] {
		col.SetSortMap(tsColumn.SortMap())
		colOffset[i+1], err = sw.writeColumn(col)
		if err != nil {
			return err
		}
	}

	// Write EventID --> offset idx
	err = sw.writeRowIndex(colOffset)

	if err != nil {
		return err
	}

	err = sw.writeSegmentInfo()

	if err != nil {
		return err
	}

	return nil
}

func (sw *SegmentWriter) writeColumn(col inmem.Column) ([]*inmem.PageDescriptor, error) {
	switch col.FieldInfo().Type {
	case segment.FieldTypeInt:
		return sw.writeColInt(col.(*inmem.ColumnInt))
	case segment.FieldTypeKeyword:
		return sw.writeColText(col.(*inmem.ColumnStr))
	case segment.FieldTypeText:
		return sw.writeColText(col.(*inmem.ColumnStr))
	case segment.FieldTypeFloat:
		return sw.writeColFloat(col.(*inmem.ColumnFloat))
	default:
		sw.log.Panic().Msgf("invalid column type [%v]", col.FieldInfo().Type)
	}
	return nil, nil
}

func (sw *SegmentWriter) writeTSColumn(tsCol *inmem.ColumnTimeStamp) (pd []*inmem.PageDescriptor, err error) {

	log := sw.log.With().Str("column", tsCol.FieldInfo().Name).Logger()

	posting := inmem.NewPostingStore()
	sl := inmem.NewSkipList(posting, inmem.IntComparator{})

	if !tsCol.Sorted() {
		log.Debug().Msg("Sorting column")
		tsCol.Sort()
	}

	// TODO Replace by a compressed representation.
	f := filepath.Join(sw.path, tsCol.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)
	defer bw.Close()

	if err != nil {
		return nil, err
	}

	page := new(inmem.Page)
	page.Total = 0
	page.StartID = 0

	// Page offset
	pd = make([]*inmem.PageDescriptor, 0)
	slice := make([]int, 0)

	for i := 0; i < tsCol.Size(); i++ {

		x := tsCol.Get(i)
		sl.Add(x, i) // save val -> ids
		slice = append(slice, x)
		page.Total++

		// TODO: por ahora fijo despues tenemos que ver....
		if i > 0 && i%config.PageSize == 0 {

			eh := encoding.NewEncoderHandler(tsCol.FieldInfo(), page)
			r := eh.DoEncode(slice)
			page.PayloadSize = len(r.([]byte))

			bw.WritePageHeader(page)
			bw.Write(r.([]byte))

			pd = append(pd, &inmem.PageDescriptor{StartID: i, Offset: bw.Offset})

			page := new(inmem.Page)
			page.Total = 0
			page.StartID = i

		}

	}

	if page.Total != 0 {
		eh := encoding.NewEncoderHandler(tsCol.FieldInfo(), page)
		r := eh.DoEncode(slice)
		page.PayloadSize = len(r.([]byte))

		bw.WritePageHeader(page)
		bw.Write(r.([]byte))

		pd = append(pd, &inmem.PageDescriptor{StartID: page.StartID, Offset: bw.Offset})
	}

	f = filepath.Join(sw.path, tsCol.FieldInfo().Name+idxExt)

	err = WritePosting(f, posting)
	if err != nil {
		return nil, err
	}

	f = filepath.Join(sw.path, tsCol.FieldInfo().Name+idxExt)

	err = WriteSkip(f, sl, config.SkipLevelSize)
	if err != nil {
		return nil, err
	}

	return

}

func (sw *SegmentWriter) writeColInt(col *inmem.ColumnInt) (pd []*inmem.PageDescriptor, err error) {

	//log := sw.log.With().Str("column", col.FieldInfo().Name).Logger()

	posting := inmem.NewPostingStore()
	sl := inmem.NewSkipList(posting, inmem.IntComparator{})

	f := filepath.Join(sw.path, col.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)
	defer bw.Close()

	if err != nil {
		return nil, err
	}

	page := new(inmem.Page)
	page.Total = 0
	page.StartID = 0

	// Page offset
	pd = make([]*inmem.PageDescriptor, 0)
	slice := make([]int, 0)

	for i := 0; i < col.Size(); i++ {

		x := col.Get(i)
		slice = append(slice, x)
		sl.Add(x, i) // save val -> ids
		page.Total++

		// TODO: por ahora fijo despues tenemos que ver....
		if i > 0 && i%config.PageSize == 0 {

			eh := encoding.NewEncoderHandler(col.FieldInfo(), page)
			r := eh.DoEncode(slice)
			page.PayloadSize = len(r.([]byte))

			bw.WritePageHeader(page)
			bw.Write(r.([]byte))

			pd = append(pd, &inmem.PageDescriptor{StartID: i, Offset: bw.Offset})

			page := new(inmem.Page)
			page.Total = 0
			page.StartID = i

		}

	}

	if page.Total != 0 {
		eh := encoding.NewEncoderHandler(col.FieldInfo(), page)
		r := eh.DoEncode(slice)
		page.PayloadSize = len(r.([]byte))

		bw.WritePageHeader(page)
		bw.Write(r.([]byte))

		pd = append(pd, &inmem.PageDescriptor{StartID: page.StartID, Offset: bw.Offset})
	}

	if col.FieldInfo().Index {

		f = filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WritePosting(f, posting)
		if err != nil {
			return nil, err
		}

		f := filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WriteSkip(f, sl, config.SkipLevelSize)
		if err != nil {
			return nil, err
		}

	}

	return
}

func (sw *SegmentWriter) writeColFloat(col *inmem.ColumnFloat) (pd []*inmem.PageDescriptor, err error) {

	//log := sw.log.With().Str("column", col.FieldInfo().Name).Logger()

	posting := inmem.NewPostingStore()
	sl := inmem.NewSkipList(posting, inmem.Float64Comparator{})

	// TODO Replace by a compressed representation.

	f := filepath.Join(sw.path, col.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)
	defer bw.Close()

	if err != nil {
		return nil, err
	}

	page := new(inmem.Page)
	page.Total = 0
	page.StartID = 0

	// Page offset
	pd = make([]*inmem.PageDescriptor, 0)
	slice := make([]uint64, 0)

	for i := 0; i < col.Size(); i++ {

		x := col.Get(i)
		bits := math.Float64bits(x)
		slice = append(slice, bits)

		sl.Add(x, i) // save val -> offsets
		page.Total++

		// TODO: por ahora fijo despues tenemos que ver....
		if i > 0 && i%config.PageSize == 0 {

			eh := encoding.NewEncoderHandler(col.FieldInfo(), page)
			r := eh.DoEncode(slice)
			page.PayloadSize = len(r.([]byte))

			bw.WritePageHeader(page)
			bw.Write(r.([]byte))

			pd = append(pd, &inmem.PageDescriptor{StartID: i, Offset: bw.Offset})

			page := new(inmem.Page)
			page.Total = 0
			page.StartID = i

		}

	}

	if page.Total != 0 {
		eh := encoding.NewEncoderHandler(col.FieldInfo(), page)
		r := eh.DoEncode(slice)
		page.PayloadSize = len(r.([]byte))

		bw.WritePageHeader(page)
		bw.Write(r.([]byte))

		pd = append(pd, &inmem.PageDescriptor{StartID: page.StartID, Offset: bw.Offset})
	}

	if col.FieldInfo().Index {

		f = filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WritePosting(f, posting)
		if err != nil {
			return nil, err
		}

		f := filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WriteSkip(f, sl, config.SkipLevelSize)
		if err != nil {
			return nil, err
		}

	}

	return

}

func (sw *SegmentWriter) writeColText(col *inmem.ColumnStr) (pd []*inmem.PageDescriptor, err error) {

	//log := sw.log.With().Str("column", col.FieldInfo().Name).Logger()

	// TODO: optimization: if the column is not indexed we don't need a
	//       posting list.
	posting := inmem.NewPostingStore()

	trie := sw.buildTrie(col, posting)

	f := filepath.Join(sw.path, col.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)
	defer bw.Close()

	if err != nil {
		return nil, err
	}

	page := new(inmem.Page)
	page.Total = 0
	page.StartID = 0

	// Page offset
	pd = make([]*inmem.PageDescriptor, 0)
	slice := make([]string, 0)

	for i := 0; i < col.Size(); i++ {

		x := col.Get(i)
		slice = append(slice, x)
		page.Total++

		// TODO: por ahora fijo despues tenemos que ver....
		if i > 0 && i%config.PageSize == 0 {

			eh := encoding.NewEncoderHandler(col.FieldInfo(), page)
			r := eh.DoEncode(slice)
			page.PayloadSize = len(r.([]byte))

			bw.WritePageHeader(page)
			bw.Write(r.([]byte))

			pd = append(pd, &inmem.PageDescriptor{StartID: i, Offset: bw.Offset})

			page := new(inmem.Page)
			page.Total = 0
			page.StartID = i

		}

	}

	if page.Total != 0 {
		eh := encoding.NewEncoderHandler(col.FieldInfo(), page)
		r := eh.DoEncode(slice)
		page.PayloadSize = len(r.([]byte))

		bw.WritePageHeader(page)
		bw.Write(r.([]byte))

		pd = append(pd, &inmem.PageDescriptor{StartID: page.StartID, Offset: bw.Offset})
	}

	if col.FieldInfo().Index {

		if col.FieldInfo().Type == segment.FieldTypeText {
			tokenizer := text.NewTokenizer()
			for i := 0; i < col.Size(); i++ {
				tokens := tokenizer.Tokenize(col.Get(i))
				for _, token := range tokens {
					trie.Add(token, uint32(i))
				}
			}
		}

		f = filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WritePosting(f, posting)
		if err != nil {
			return nil, err
		}

		f := filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WriteTrie(f, trie)
		if err != nil {
			return nil, err
		}

	}

	return

}

func (sw *SegmentWriter) buildTrie(col *inmem.ColumnStr, posting *inmem.PostingStore) *inmem.BTrie {
	trie := inmem.NewBtrie(posting)
	for i := 0; i < col.Size(); i++ {
		trie.Add(col.Get(i), uint32(i))
	}
	return trie
}

func (sw *SegmentWriter) writeRowIndex(offsets [][]*inmem.PageDescriptor) error {
	for i, col := range offsets {
		f := filepath.Join(sw.path, sw.segment.Columns[i].FieldInfo().Name+posExt)
		WriteStoreIdx(f, col, config.SkipLevelSize)
	}
	return nil
}

func (sw *SegmentWriter) writeSegmentInfo() error {

	file := filepath.Join(sw.path, "info")

	bw, err := io.NewBinaryWriter(file)

	if err != nil {
		return err
	}

	// Header
	err = bw.WriteHeader(io.SegmentInfo)

	if err != nil {
		return err
	}

	// Index name
	err = bw.WriteString(sw.segment.IndexInfo.Name)

	if err != nil {
		return err
	}

	// Field Count
	err = bw.WriteVarInt(len(sw.segment.IndexInfo.Fields))

	if err != nil {
		return err
	}

	// Field info
	for _, field := range sw.segment.IndexInfo.Fields {

		err = bw.WriteString(field.Name)
		if err != nil {
			return err
		}

		err = bw.WriteByte(byte(field.Type))
		if err != nil {
			return err
		}

		i := byte(0)
		if field.Index {
			i = 1
		}

		err = bw.WriteByte(i)
		if err != nil {
			return err
		}

	}

	// segment stats.

	// Event count
	err = bw.WriteVarUInt32(sw.segment.EventCount)

	if err != nil {
		return err
	}

	err = bw.Close()

	if err != nil {
		return err
	}

	return nil

	// TODO add field cardinality, max/min TS and SegmentID

}

func WriteSegment(path string, segment *inmem.Segment) error {

	sw := NewSegmentWriter(path, segment)
	return sw.Write()

}
