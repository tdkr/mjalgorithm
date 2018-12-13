/*
@Time : 2018/12/10 9:37
@Author : RonanLuo
*/
package hu

import (
	"fmt"
)

func CalcKey(matrix [][]int, pos []int) (key int64) {
	isContinue := false /*是否连续*/
	index := 0
	for row := 0; row < 4; row++ {
		for col := 0; col < 9; col++ {
			num := matrix[row][col]
			if num == 0 {
				continue
			}
			isContinue = row < 3 && col > 0 && matrix[row][col-1] > 0
			switch num {
			case 1: // 1:0, 10
				if isContinue {
					key = key << 1
				} else {
					key = (key << 2) + (1 << 2) - 2
				}
			case 2, 3, 4: // 2:110,1110, 3:11110, 111110, 4:1111110, 11111110
				bn := uint(2 * num)
				if isContinue {
					bn--
				}
				key = (key << bn) + (1 << bn) - 2
			}
			pos[index] = row*10 + col + 1
			index++
			//fmt.Println("calcKey", row, col, num, strconv.FormatInt(key, 2), isContinue)
		}
	}
	return
}

func CalcGroups(huTable map[int64][][]int, key int64, pos []int) (result [][]int) {
	if arr, ok := huTable[key]; ok {
		result = make([][]int, len(arr))
		for i, v := range arr {
			result[i] = make([]int, len(v))
			result[i][0] = pos[v[0]]
			result[i][1] = v[1]
			idx := 2
			num := v[1]
			for j := 0; j < num; j++ {
				result[i][j+idx] = pos[v[idx+j]]
			}
			result[i][1+num+1] = v[1+num+1]
			idx = 1 + num + 1 + 1
			num = v[1+num+1]
			for j := 0; j < num; j++ {
				result[i][j+idx] = pos[v[idx+j]]
			}
			//fmt.Printf("AnalyseHuInfo, key : %d, %v\n", result[i])
		}
	} else {
		fmt.Printf("AnalyseHuInfo, key %b not found\n", key)
	}
	return
}
