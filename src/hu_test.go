package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
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
		CheckHu(huTable, cards)
	}
}

func BenchmarkCheckHuWithLZ(b *testing.B) {
	cards := []int{
		0, 0, 0, 3, 4, 4, 3, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	lzFlag := map[int]bool{
		4: true,
		5: true,
	}
	for i := 0; i < b.N; i++ {
		CheckHuWithLZ(huTable, cards, lzFlag)
	}
}

func TestCheckHu(t *testing.T) {
	cards := []int{
		0, 0, 0, 3, 4, 4, 3, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	ret := CheckHu(huTable, cards)
	fmt.Println("CheckHu", cards, ret)
}

func TestCheckHuWithLZ(t *testing.T) {
	cards := []int{
		0, 0, 0, 3, 4, 4, 3, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	lzFlag := map[int]bool{
		4: true,
		5: true,
	}
	ret := CheckHuWithLZ(huTable, cards, lzFlag)
	fmt.Println("CheckHuWithLZ", cards, lzFlag, ret)
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
