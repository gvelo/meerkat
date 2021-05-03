package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	id := uuid.New()

	le := &LogEntry{
		Id:     id[:],
		Time:   time.Now(),
		Src:    "srcNode",
		OpType: OpType_CREATE_DATABASE,
		Desc: &LogEntry_DatabaseDesc{
			DatabaseDesc: &DatabaseDescriptor{
				Name:                      "db-name",
				NumOfPartitions:           4,
				PartitionAllocMap:         []*NodeList{{NodeId: [][]byte{id[:]}}},
				TableIngestionAliases:     map[string]string{"alias1": "col1"},
				TableQueryAliases:         map[string]string{"alias1": "col1"},
				AllowDynamicTableCreation: false,
			},
		},
	}

	s, err := json.Marshal(le)
	jsonpb.Marshaler{}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(s))

	jj := &LogEntry{}

	err = json.Unmarshal(s, jj)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(jj)

}

func TestName1(t *testing.T) {

	id := uuid.New()

	dd := &DatabaseDescriptor{
		Name:                      "db-name",
		NumOfPartitions:           4,
		PartitionAllocMap:         []*NodeList{{NodeId: [][]byte{id[:]}}},
		TableIngestionAliases:     map[string]string{"alias1": "col1"},
		TableQueryAliases:         map[string]string{"alias1": "col1"},
		AllowDynamicTableCreation: false,
	}

	s, err := json.Marshal(dd)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(s))

	jj := &DatabaseDescriptor{}

	err = json.Unmarshal(s, jj)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(jj)



}
