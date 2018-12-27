package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var huTable map[int64][][]int

func BenchmarkCheckHu(b *testing.B) {
	cards := []int{
		0, 0, 0, 3, 4, 4, 3, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	for i := 0; i < b.N; i++ {
		CheckHuGroup(huTable, cards)
	}
}

func BenchmarkCheckHuWithLZ(b *testing.B) {
	cards := []int{
		//0, 0, 0, 3, 4, 4, 3, 0, 0,
		//0, 0, 0, 0, 0, 0, 0, 0, 0,
		//0, 0, 0, 0, 0, 0, 0, 0, 0,
		//0, 0, 0, 0, 0, 0, 0, 0, 0,

		3, 2, 1, 1, 1, 1, 1, 1, 3,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	lzFlag := map[int]bool{
		0: true,
		8: true,
	}
	lzList := []int{}
	for i, v := range cards {
		if lzFlag[i] {
			for j := 0; j < v; j++ {
				lzList = append(lzList, i)
			}
			cards[i] = 0
		}
	}
	for i := 0; i < b.N; i++ {
		CheckHuWithLZ(huTable, cards, lzList, lzFlag)
	}
}

func TestCheckHu(t *testing.T) {
	cards := []int{
		2, 3, 4, 3, 2, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	ret := CheckHuGroup(huTable, cards)
	fmt.Println("CheckHu", cards, ret)
}

func TestCheckHuWithLZ(t *testing.T) {
	cards := []int{
		//1, 1, 1, 0, 1, 1, 1, 0, 4,
		//0, 0, 0, 0, 0, 0, 0, 0, 0,
		//0, 0, 0, 0, 0, 0, 0, 0, 0,
		//0, 0, 0, 0, 0, 0, 1,

		3, 2, 1, 1, 1, 1, 1, 1, 3,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}
	lzFlag := map[int]bool{
		//33: true,

		0: true,
		1: true,
		8: true,
	}
	lzList := []int{}
	for i, v := range cards {
		if lzFlag[i] {
			for j := 0; j < v; j++ {
				lzList = append(lzList, i)
			}
			cards[i] = 0
		}
	}
	start := time.Now()
	var ret []*HuGroup
	for i := 0; i < 10000; i++ {
		ret = CheckHuWithLZ(huTable, cards, lzList, lzFlag)
	}
	fmt.Println("TestCheckHuWithLZ", time.Since(start))
	fmt.Println("TestCheckHuWithLZ", len(ret))
	for i, v := range ret {
		fmt.Println("TestCheckHuWithLZ, result", i, v.EyeList, v.KeList, v.ShunList)
	}
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	path := filepath.Join(wd, "output.json")
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	if e := json.Unmarshal(bytes, &huTable); e != nil {
		log.Fatal(e)
		return
	}
}
