/*
@Time : 2018/11/6 11:21
@Author : RonanLuo
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"time"

	"github.com/tdkr/mjalgorithm/src"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func main() {
	PrintMemUsage()
	cards := []int{
		//1, 1, 1, 2, 3, 4, 6, 7, 8, 31, 31, 31, 33, 33,
		//1, 1, 1, 2, 3, 4, 5, 5,
		//1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7,
		//3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6,
		1, 4, 7,
		10, 13, 16,
		33, 33, 33, 33, 33, 33, 33, 33,
	}
	matrix := make([]int, 36)
	for _, v := range cards {
		matrix[v]++
	}

	file := "C:/Users/ronanluo/go/src/gogit.oa.com/266/mahjong/service/table/logic/utils/hu/output.json"
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("ReadFile failed", err)
		return
	}
	huTable := make(map[int64][][]int)
	if e := json.Unmarshal(bytes, &huTable); e != nil {
		log.Println("DecodeFile failed", e)
		return
	}
	bytes = nil

	var key int64
	var pos = make([]int, 14)
	start := time.Now()
	for i := 1; i < 10000000; i++ {
		key = hu.CalcKey(matrix, pos)
	}
	result := hu.CalcGroups(huTable, key, pos)
	fmt.Println("testFinished", time.Since(start))
	fmt.Println("result", matrix, result)
	PrintMemUsage()

	start = time.Now()
	lzFlag := make(map[int]bool)
	lzFlag[33] = true
	var ret [][]int
	for i := 0; i < 1; i++ {
		ret = hu.CheckHuWithLZ(huTable, matrix, lzFlag)
	}
	fmt.Println("testFinished", time.Since(start))
	fmt.Println("result", ret)
}
