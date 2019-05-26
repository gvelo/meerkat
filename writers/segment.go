package writers

import (
	"eventdb/io"
	"eventdb/segment"
	"eventdb/segment/inmem"
	"eventdb/text"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
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
	colOffset := make([][]int, len(sw.segment.Columns))

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

	/*

			idx := sw.createAndLoadFieldIdx()

			err := sw.writePosting()
		err := sw.writeStore()
		if err != nil {
			return err
		}

		err = sw.writePosting()

			if err != nil {
				return err
			}

			// write all idx to disk

			for i, field := range sw.segment.IndexInfo.Fields {

				if field.Index {

					switch t := idx[i].(type) {
					case *inmem.BTrie:
						err := sw.writeBtrie(field, idx[i].(*inmem.BTrie))
						if err != nil {
							return err
						}
					case *inmem.SkipList:
						err := sw.writeSL(field, idx[i].(*inmem.SkipList))
						if err != nil {
							return err
						}
					default:
						// TODO:  in the case of ts this is nil, FIX IT
						if t != nil {
							sw.log.Panic().
								Str("index", sw.segment.IndexInfo.Name).
								Str("segmentID", sw.segment.ID).
								Msgf("Unknown index type: %T", t)

						}
					}
				}

			}

			err = sw.writeSegmentInfo()

			if err != nil {
				return err
			}

			return nil
	*/
}

func (sw *SegmentWriter) writeColumn(col inmem.Column) ([]int, error) {
	switch col.FieldInfo().Type {
	case segment.FieldTypeInt:
		return sw.writeColInt(col.(*inmem.ColumnInt))
	case segment.FieldTypeKeyword:
		return sw.writeColKeyword(col.(*inmem.ColumnStr))
	case segment.FieldTypeText:
		return sw.writeColText(col.(*inmem.ColumnStr))
	case segment.FieldTypeFloat:
		return sw.writeColFloat(col.(*inmem.ColumnFloat))
	default:
		sw.log.Panic().Msgf("invalid column type [%v]", col.FieldInfo().Type)
	}
	return nil, nil
}

func (sw *SegmentWriter) writeTSColumn(tsCol *inmem.ColumnTimeStamp) ([]int, error) {

	log := sw.log.With().Str("column", tsCol.FieldInfo().Name).Logger()

	if !tsCol.Sorted() {
		log.Debug().Msg("Sorting column")
		tsCol.Sort()
	}

	// TODO Replace by a compressed representation.

	f := filepath.Join(sw.path, tsCol.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)

	if err != nil {
		return nil, err
	}

	offset := make([]int, tsCol.Size())

	for i := 0; i < tsCol.Size(); i++ {
		offset[i] = bw.Offset
		err := bw.WriteVarInt(tsCol.Get(i))
		if err != nil {
			return nil, err
		}
	}

	return offset, nil

}

func (sw *SegmentWriter) writeColInt(col *inmem.ColumnInt) ([]int, error) {

	//log := sw.log.With().Str("column", col.FieldInfo().Name).Logger()

	posting := inmem.NewPostingStore()

	var u inmem.OnUpdate = func(n *inmem.SLNode, v interface{}) interface{} {
		if n.UserData == nil {
			n.UserData = posting.NewPostingList(uint32(v.(int)))
		} else {
			n.UserData.(*inmem.PostingList).Bitmap.Add(uint32(v.(int)))
		}
		return n.UserData
	}
	sl := inmem.NewSkipList(posting, u, inmem.IntComparator{})

	// TODO Replace by a compressed representation.

	f := filepath.Join(sw.path, col.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)

	if err != nil {
		return nil, err
	}

	offset := make([]int, col.Size())

	for i := 0; i < col.Size(); i++ {
		offset[i] = bw.Offset
		x := col.Get(i)
		err := bw.WriteVarInt(x)
		sl.Add(x, i) // save val -> offsets
		if err != nil {
			return nil, err
		}
	}

	if col.FieldInfo().Index {

		f = filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WritePosting(f, posting)
		if err != nil {
			return nil, err
		}

		f := filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WriteSkip(f, sl, 100)
		if err != nil {
			return nil, err
		}

	}

	return offset, nil

}

func (sw *SegmentWriter) writeColFloat(col *inmem.ColumnFloat) ([]int, error) {

	//log := sw.log.With().Str("column", col.FieldInfo().Name).Logger()

	posting := inmem.NewPostingStore()

	var u inmem.OnUpdate = func(n *inmem.SLNode, v interface{}) interface{} {
		if n.UserData == nil {
			n.UserData = posting.NewPostingList(uint32(v.(int)))
		} else {
			n.UserData.(*inmem.PostingList).Bitmap.Add(uint32(v.(int)))
		}
		return n.UserData
	}
	sl := inmem.NewSkipList(posting, u, inmem.Float64Comparator{})

	// TODO Replace by a compressed representation.

	f := filepath.Join(sw.path, col.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)

	if err != nil {
		return nil, err
	}

	offset := make([]int, col.Size())

	for i := 0; i < col.Size(); i++ {
		offset[i] = bw.Offset
		x := col.Get(i)
		bits := math.Float64bits(x)
		err := bw.WriteFixedUint64(bits)
		sl.Add(x, i) // save val -> offsets
		if err != nil {
			return nil, err
		}
	}

	if col.FieldInfo().Index {

		f = filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WritePosting(f, posting)
		if err != nil {
			return nil, err
		}

		f := filepath.Join(sw.path, col.FieldInfo().Name+idxExt)

		err = WriteSkip(f, sl, 100)
		if err != nil {
			return nil, err
		}

	}

	return offset, nil

}

func (sw *SegmentWriter) writeColKeyword(col *inmem.ColumnStr) ([]int, error) {

	//log := sw.log.With().Str("column", col.FieldInfo().Name).Logger()

	// TODO: optimization: if the column is not indexed we don't need a
	//       posting list.
	posting := inmem.NewPostingStore()

	trie := sw.buildTrie(col, posting)

	// now that we have the column cardinality ( from the btrie)  and the
	// column size we can chosee to write the column dictionary encoded or raw.

	// TODO Replace by a compressed representation.

	f := filepath.Join(sw.path, col.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)

	if err != nil {
		return nil, err
	}

	offset := make([]int, col.Size())

	for i := 0; i < col.Size(); i++ {
		offset[i] = bw.Offset
		err := bw.WriteString(col.Get(i))
		if err != nil {
			return nil, err
		}
	}

	if col.FieldInfo().Index {

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

	return offset, nil

}

func (sw *SegmentWriter) buildTrie(col *inmem.ColumnStr, posting *inmem.PostingStore) *inmem.BTrie {
	trie := inmem.NewBtrie(posting)
	for i := 0; i < col.Size(); i++ {
		trie.Add(col.Get(i), uint32(i))
	}
	return trie
}

func (sw *SegmentWriter) writeColText(col *inmem.ColumnStr) ([]int, error) {

	//log := sw.log.With().Str("column", col.FieldInfo().Name).Logger()

	posting := inmem.NewPostingStore()

	trie := inmem.NewBtrie(posting)

	// TODO Replace by a compressed representation.

	f := filepath.Join(sw.path, col.FieldInfo().Name+colExt)

	bw, err := io.NewBinaryWriter(f)

	if err != nil {
		return nil, err
	}

	tokenizer := text.NewTokenizer()

	offset := make([]int, col.Size())

	for i := 0; i < col.Size(); i++ {
		offset[i] = bw.Offset
		s := col.Get(i)
		err := bw.WriteString(s)
		if err != nil {
			return nil, err
		}
		if col.FieldInfo().Index {
			tokens := tokenizer.Tokenize(s)
			for _, token := range tokens {
				trie.Add(token, uint32(i))
			}
		}
	}

	if col.FieldInfo().Index {

		f = filepath.Join(sw.path, col.FieldInfo().Name+posExt)

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

	return offset, nil

}

func (sw *SegmentWriter) writeRowIndex(offsets [][]int) error {

	f := filepath.Join(sw.path, "Store"+idxExt)

	return WriteStoreIdx(f, offsets, 100)
}

/*

func byTs(a, b interface{}) bool {
	return a.(segment.Event)[tsField].(uint64) < b.(segment.Event)[tsField].(uint64)
}

func (sw *SegmentWriter) createAndLoadFieldIdx() []interface{} {

	idx := make([]interface{}, len(sw.segment.IndexInfo.Fields))

	for i, fInfo := range sw.segment.IndexInfo.Fields {
		switch fInfo.Type {
		case segment.FieldTypeTimestamp:
			// TODO add the proper index here , should I?
		case segment.FieldTypeInt:
			var u inmem.OnUpdate = func(n *inmem.SLNode, v interface{}) interface{} {
				if n.UserData == nil {
					n.UserData = sw.segment.PostingStore.NewPostingList(uint32(v.(int)))
				} else {
					n.UserData.(inmem.PostingList).Bitmap.Add(uint32(v.(int)))
				}
				return n
			}
			idx[i] = inmem.NewSkipList(sw.segment.PostingStore, u, inmem.Uint64Comparator{})
		case segment.FieldTypeFloat:
			var u inmem.OnUpdate = func(n *inmem.SLNode, v interface{}) interface{} {
				if n.UserData == nil {
					n.UserData = sw.segment.PostingStore.NewPostingList(uint32(v.(int)))
				} else {
					n.UserData.(inmem.PostingList).Bitmap.Add(uint32(v.(int)))
				}
				return n
			}
			idx[i] = inmem.NewSkipList(sw.segment.PostingStore, u, inmem.Float64Comparator{})
		case segment.FieldTypeKeyword:
			idx[i] = inmem.NewBtrie(sw.segment.PostingStore)
		case segment.FieldTypeText:
			idx[i] = inmem.NewBtrie(sw.segment.PostingStore)
		default:
			log.Panic().Int("Type", int(fInfo.Type)).Msg("Invalid Type")
		}
	}

	inmem.Sort(sw.segment.FieldStorage, byTs)

	for x, n := range sw.segment.FieldStorage {

		for i, info := range sw.segment.IndexInfo.Fields {

			if info.Index {

				switch info.Type {

				case segment.FieldTypeInt:
					idx := idx[i].(*inmem.SkipList)
					eventValue := n[info.Name].(uint64)
					idx.Add(eventValue, x)
				case segment.FieldTypeFloat:
					idx := idx[i].(*inmem.SkipList)
					eventValue := n[info.Name].(float64)
					idx.Add(eventValue, x)
				case segment.FieldTypeKeyword:
					idx := idx[i].(*inmem.BTrie)
					eventValue := n[info.Name].(string)
					idx.Add(eventValue, uint32(x))
				case segment.FieldTypeText:
					idx := idx[i].(*inmem.BTrie)
					eventValue := n.(segment.Event)[info.Name].(string)
					tokens := sw.segment.Tokenizer.Tokenize(eventValue)
					for _, token := range tokens {
						idx.Add(token, uint32(x))
					}
				case segment.FieldTypeTimestamp:
				//TODO Add to the proper index.
				default:
					log.Panic().Int("Type", int(info.Type)).Msg("Invalid Type")
				}

			}

		}

	}

	return idx
}

func (sw *SegmentWriter) writeStore() error {

	fileName := filepath.Join(sw.path, "store")

	_, err := WriteStore(fileName, sw.segment.FieldStorage, sw.segment.IndexInfo, 100)

	if err != nil {
		sw.log.Error().
			Err(err).
			Msg("error writing posting list")
		return errors.Wrapf(err, "error writing posting list segment=%v", sw.segment.ID)
	}

	return nil

}

func (sw *SegmentWriter) writePosting() error {

	fileName := filepath.Join(sw.path, "posting")

	err := WritePosting(fileName, sw.segment.PostingStore.Store)

	if err != nil {
		sw.log.Error().
			Err(err).
			Msg("error writing posting list")
		return errors.Wrapf(err, "error writing posting list segment=%v", sw.segment.ID)
	}

	return nil

}

func (sw *SegmentWriter) writeBtrie(field *segment.FieldInfo, btrie *inmem.BTrie) error {

	fileName := filepath.Join(sw.path, field.Name+idxExt)

	log := sw.log.
		With().
		Str("field", field.Name).
		Str("filename", fileName).
		Logger()

	log.Debug().Msg("writing btrie index")

	writer, err := newTrieWriter(fileName)

	if err != nil {
		sw.log.Error().
			Err(err).
			Msg("error creating btrie writer")
		return errors.Wrapf(err, "error creating btrie writer segment=%v field=%v", sw.segment.ID, field.Name)
	}

	defer writer.Close()

	err = writer.Write(btrie)

	if err != nil {
		sw.log.Error().
			Err(err).
			Msg("error writing btrie")
		return errors.Wrapf(err, "error writing btrie writer segment=%v field=%v", sw.segment.ID, field.Name)
	}

	// at this point, we don't defer writer closing so we can
	// report fsync and other os errors

	err = writer.Close()

	if err != nil {
		sw.log.Error().
			Err(err).
			Msg("error closing btrie writer")
		return errors.Wrapf(err, "error closing btrie writer segment=%v field=%v", sw.segment.ID, field.Name)

	}

	return nil
}

func (sw *SegmentWriter) writeSL(field *segment.FieldInfo, sl *inmem.SkipList) error {

	fileName := filepath.Join(sw.path, field.Name+idxExt)

	// TODO add a default for ixl
	err := WriteSkip(fileName, sl, 100)

	if err != nil {
		return err
	}

	return nil

}

*/

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
