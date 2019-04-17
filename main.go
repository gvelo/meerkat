package main

import (
	"bufio"
	"eventdb/segment/inmem"
	"eventdb/writers"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blevesearch/segment"
)

var m map[string]bool
var trie *inmem.BTrie

func main() {

	trie = inmem.NewBtrie()

	file, err := os.Open("/home/gabrielvelo/Downloads/logs/feeder.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	m = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println("--------------------------------------------")
		// fmt.Println(scanner.Text())
		parseStr(scanner.Text())
		//fmt.Println("--------------------------------------------")
	}

	tk := make([]string, len(m))
	i := 0
	for k := range m {
		tk[i] = k
		i++
	}

	//sort.Strings(tk)

	l := 0
	for _, v := range tk {
		//fmt.Println(v)
		l = l + len(v)
	}

	triewriter, err := writers.NewTrieWriter("/tmp/feeder-trie.bin")
	if err != nil {
		log.Fatal(err)
	}

	_,err = triewriter.Write(trie)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("total de terminos ", len(m))
	fmt.Println("done.... ")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseStr(s string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(segment.SplitWords)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		s := scanner.Text()
		m[s] = true
		//addString(s)
		trie.Add(s,86587657)
	}
}
