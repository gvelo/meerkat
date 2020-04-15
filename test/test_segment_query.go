package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"math/rand"
	"meerkat/internal/buffer"
	"meerkat/internal/executor"
	"meerkat/internal/schema"
	"meerkat/internal/storage"
	"meerkat/internal/util/testutil"
	"path"
	"time"
)

func main() {

	//createIndex()

	//generateData()

	execute()

}

func execute() {

	now := time.Now()

	indexInfo := createIndexInfo()
	s := createSegment(indexInfo, now)
	op := buildPhysicPlan(s, &indexInfo, now)

	op.Init()
	n := op.Next()
	if n == nil {
		println(" No result found")
	}

	for r := 0; r < len(n); r++ {
		for i := 0; i < len(n[r]); i++ {
			fmt.Printf(" %s ", n[r][i])
		}
		fmt.Println("")
	}

}

func createSegment(indexInfo schema.IndexInfo, now time.Time) *storage.Segment {

	buf := createBuffers(indexInfo, 250, now)

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

func buildPhysicPlan(s *storage.Segment, ii *schema.IndexInfo, now time.Time) executor.StringOperator {
	from := now
	to := time.Now()
	fmt.Print("FROM ", now.Format("2006-01-02T15:04:05.00000"))
	fmt.Print(" TO ", to.Format("2006-01-02T15:04:05.00000"), "\n")

	sz := 200
	ctx := executor.NewContext(s, ii)

	op1 := executor.NewTimeColumnScanOperator(ctx, executor.Between, int(from.UnixNano()), int(to.UnixNano()), "_ts", sz, false)
	op2 := executor.NewStringColumnScanOperator(ctx, executor.Contains, "Error", "message", sz, false)
	op3 := executor.NewBinaryUint32Operator(ctx, executor.And, op1, op2, sz)
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

func createBuffers(indexInfo schema.IndexInfo, testLen int, now time.Time) *buffer.Table {

	table := buffer.NewTable(indexInfo)

	for i := 0; i < testLen; i++ {
		r := buffer.NewRow(len(indexInfo.Fields))
		for _, f := range indexInfo.Fields {
			switch f.FieldType {
			case schema.FieldType_TIMESTAMP:
				t := int(now.UnixNano())
				t += rand.Intn(2000)
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

func createIndexInfo() schema.IndexInfo {

	return schema.IndexInfo{
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

}
