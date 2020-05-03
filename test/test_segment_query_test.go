package main

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/buffer"
	"meerkat/internal/executor"
	"meerkat/internal/schema"
	"meerkat/internal/storage"
	"meerkat/internal/util/testutil"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

// TODO(sebad): put more operators.
func TestNewBinaryUint32Operator(t *testing.T) {
	r := testCases("2006-01-02T15:04:09.00000", "2006-01-02T15:04:12.00000", "Error")
	assert.Len(t, r, 1)

	r = testCases("2006-01-02T15:04:05.00000", "2006-01-02T15:04:12.00000", "Error")
	assert.Len(t, r, 5)

}

func testCases(from, to, strToFind string) [][]string {

	f, _ := time.Parse("2006-01-02T15:04:05.00000", from)
	t, _ := time.Parse("2006-01-02T15:04:05.00000", to)

	fmt.Print("FROM ", f.Format("2006-01-02T15:04:05.00000"))
	fmt.Println(" TO ", t.Format("2006-01-02T15:04:05.00000"))

	indexInfo := createIndexInfo()
	s := createSegment(indexInfo)
	op := buildPhysicPlan(s, indexInfo, int(f.UnixNano()), int(t.UnixNano()), strToFind)

	op.Init()
	n := op.Next()
	if n == nil {
		println(" No result found")
	}

	res := make([][]string, 0)
	for r := 0; r < len(n); r++ {
		row := make([]string, 0)
		for i := 0; i < len(n[r]); i++ {
			fmt.Printf(" %s ", n[r][i])
			row = append(row, n[r][i])
		}
		res = append(res, row)
		fmt.Println("")
	}

	return res

}

func createSegment(indexInfo *schema.IndexInfo) *storage.Segment {

	buf := createBuffersFromFile(indexInfo, "/Users/sebad/desa/workspace_go/meerkat/test/data.csv")

	filePath := "/Users/sebad/meerkat/segments"

	uid, err := uuid.NewUUID()
	if err != nil {
		log.Error().Err(err)
	}
	sw := storage.NewSegmentWriter(filePath, uid, buf)

	err = sw.Write()
	if err != nil {
		log.Error().Err(err)
	}

	filePath = path.Join(filePath, uid.String())

	seg, err := storage.ReadSegment(filePath)

	if err != nil {
		log.Error().Err(err)
	}

	return seg
}

func buildPhysicPlan(s *storage.Segment, ii *schema.IndexInfo, from, to int, strToFind string) executor.StringOperator {

	sz := 200
	ctx := executor.NewContext(s, ii, sz)

	op1 := executor.NewTimeColumnScanOperator(ctx, executor.Between, from, to, "_ts", false)
	op2 := executor.NewStringColumnScanOperator(ctx, executor.Contains, strToFind, "message", false)
	op3 := executor.NewBinaryUint32Operator(ctx, executor.And, op1, op2)
	op4 := executor.NewMaterializeOperator(ctx, op3, nil)
	op5 := executor.NewTimeBucketOperator(ctx, op4, "1m")

	/* ag := []executor.Aggregation{
		{
			AggType: executor.Count,
			AggCol:  0,
		},
	}

	 keys := []int{0}

	 op6 := executor.NewHashAggregateOperator(ctx, op5, ag, keys)
	*/

	op6 := executor.NewColumnToRowOperator(ctx, op5)

	return op6
}

func createBuffersFromFile(indexInfo *schema.IndexInfo, file string) *buffer.Table {
	table := buffer.NewTable(indexInfo)

	csvFile, err := os.Open(file)
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't open the csv file ")
	}

	// Parse the file
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()

	fields := make(map[string]schema.Field, 0)
	fieldNameIdx := make(map[int]string, 0)

	for _, it := range indexInfo.Fields {
		fields[it.Name] = it
	}

	for i, r := range records {

		if i == 0 {
			for x, f := range r {
				fieldNameIdx[x] = f
			}
			continue
		}

		row := buffer.NewRow(len(indexInfo.Fields))
		for x, c := range r {

			f := fields[fieldNameIdx[x]]

			switch f.FieldType {
			case schema.FieldType_TIMESTAMP:
				if s, err := time.Parse("2006-01-02T15:04:05", c); err == nil {
					row.AddCol(f.Id, int(s.UnixNano()))
				} else {
					log.Error().Err(err)
				}
			case schema.FieldType_INT:
				if s, err := strconv.Atoi(c); err == nil {
					row.AddCol(f.Id, s)
				} else {
					log.Error().Err(err)
				}
			case schema.FieldType_FLOAT:
				if s, err := strconv.ParseFloat(r[i], 64); err == nil {
					row.AddCol(f.Id, s)
				} else {
					log.Error().Err(err)
				}
			case schema.FieldType_STRING:
				row.AddCol(f.Id, c)
			}

		}
		table.AppendRow(row)
	}

	return table
}

func createBuffers(indexInfo *schema.IndexInfo, testLen int, now time.Time) *buffer.Table {

	table := buffer.NewTable(indexInfo)

	for i := 0; i < testLen; i++ {
		r := buffer.NewRow(len(indexInfo.Fields))
		for _, f := range indexInfo.Fields {
			switch f.FieldType {
			case schema.FieldType_TIMESTAMP:
				t := now.UnixNano()
				d, _ := time.ParseDuration("1s")
				t += d.Milliseconds()
				r.AddCol(f.Id, t)
			case schema.FieldType_INT:
				if f.Nullable {
					if rand.Intn(3) == 2 {
						r.AddCol(f.Id, rand.Int())
					}
				} else {
					r.AddCol(f.Id, rand.Int())
				}
			case schema.FieldType_STRING:
				if f.Nullable {
					if rand.Intn(3) == 2 {
						r.AddCol(f.Id, testutil.RandomString(25))
					}
				} else {
					if rand.Intn(3) == 2 {
						r.AddCol(f.Id, "Error cagamos")
					} else {
						r.AddCol(f.Id, testutil.RandomString(25))
					}
				}
			}

		}

		table.AppendRow(r)

	}

	return table

}

func createIndexInfo() *schema.IndexInfo {

	ii := &schema.IndexInfo{
		Id:             "Log",
		Name:           "Log",
		Desc:           "Log",
		Created:        time.Time{},
		Updated:        time.Time{},
		PartitionAlloc: schema.PartitionAlloc{},
		Fields: []schema.Field{
			{
				Id:        "_ts",
				Name:      "_ts",
				Desc:      "",
				IndexId:   "Log",
				FieldType: schema.FieldType_TIMESTAMP,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
			{
				Id:        "message",
				Name:      "message",
				Desc:      "",
				IndexId:   "Log",
				FieldType: schema.FieldType_STRING,
				Nullable:  false,
				Created:   time.Time{},
				Updated:   time.Time{},
			},
		},
	}
	return ii

}
