package writers

import (
	"eventdb/io"
	"eventdb/segment"
	"eventdb/segment/inmem"
	"github.com/pkg/errors"
	"github.com/psilva261/timsort"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

const (
	idxExt = ".idx"
	binExt = ".bin"
)

type SegmentWriter struct {
	log     zerolog.Logger
	path    string
	segment *inmem.Segment
}

func NewSegmentWriter(path string, segment *inmem.Segment) *SegmentWriter {

	//TODO check for segment state

	sw := &SegmentWriter{
		segment: segment,
		path:    path,
	}

	sw.log = log.With().
		Str("component", "SegmentWriter").
		Str("segmentID", segment.ID).
		Str("indexName", segment.IndexInfo.Name).
		Logger()

	return sw

}

func (sw *SegmentWriter) Write() error {

	sw.log.Info().Msg("Starting to write segment")

	idx := sw.createAndLoadFieldIdx()

	err := sw.writePosting()

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

}

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

	timsort.Sort(sw.segment.FieldStorage, byTs)

	for x, n := range sw.segment.FieldStorage {

		for i, info := range sw.segment.IndexInfo.Fields {

			if info.Index {

				switch info.Type {

				case segment.FieldTypeInt:
					idx := idx[i].(*inmem.SkipList)
					eventValue := n.(segment.Event)[info.Name].(uint64)
					idx.Add(eventValue, x)
				case segment.FieldTypeFloat:
					idx := idx[i].(*inmem.SkipList)
					eventValue := n.(segment.Event)[info.Name].(float64)
					idx.Add(eventValue, x)
				case segment.FieldTypeKeyword:
					idx := idx[i].(*inmem.BTrie)
					eventValue := n.(segment.Event)[info.Name].(string)
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

	writer, err := NewTrieWriter(fileName)

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
	err = bw.WriteVarUInt32(sw.segment.EventID)

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
