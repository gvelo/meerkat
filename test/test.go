package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {

	//createIndex()

	generateData()

}

func createIndex() {

	s := `{
			  "name": "test-index",
			  "Desc": "test",
			  "fields": [
				{
				  "name": "field1",
				  "desc": "field1",
				  "field_type": "INT"
				}
			  ],
			  "partitionAlloc": {
				"NumOfPartitions": 3,
				"AllocMap": {}
			  }
			}`

	Post("http://127.0.0.1:9090/index", s)

}

func Post(url string, body string) {
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(body))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status, resp.StatusCode, err)
	fmt.Println("body=", string(r))

}

func generateData() {

	s := ""

	for i := 0; i < 100; i++ {
		s = s + fmt.Sprintf("{\"field1\":%v}\n", i)
	}

	for i := 0; i < 5000; i++ {

		go Post("http://127.0.0.1:9090/testTable/ingest", s)

		time.Sleep(100 * time.Millisecond)

	}

}
