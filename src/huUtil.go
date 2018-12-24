/*
@Time : 2018/12/10 9:37
@Author : RonanLuo
*/
package hu

import (
	"fmt"
)

func CalcKey(cards []int, pos []int) (key int64) {
	isContinue := false /*是否连续*/
	index := 0
	for row := 0; row < 4; row++ {
		for col := 0; col < 9; col++ {
			num := cards[row*9+col]
			if num == 0 {
				continue
			}
			isContinue = row < 3 && col > 0 && cards[row*9+col-1] > 0
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

func CheckHu(huTable map[int64][][]int, cards []int, lzFlag map[int]bool) [][]int {
	lzNum := 0
	dupCards := make([]int, len(cards))
	copy(dupCards, cards)
	for k, _ := range lzFlag {
		lzNum += cards[k]
		dupCards[k] = 0
	}
	if lzNum == 0 {
		pos := make([]int, 14)
		key := CalcKey(cards, pos)
		result := CalcGroups(huTable, key, pos)
		return result
	} else {
		for i := 0; i < 34; i++ {
			if cards[i] > 0 {
				if cards[i] < 2 {
					cards[i] = 0
				} else {
					cards[i] -= 2
				}
				checkHuWithEye(huTable, cards, lzNum, i)
			}
		}
	}
	return nil
}

func scanKeFromPos(cards []int, lzNum int, pos int) (bool, int) {
	if cards[pos] == 0 {
		return false, lzNum
	}
	if cards[pos]+lzNum < 3 {
		return false, lzNum
	}
	if cards[pos] < 3 {
		lzNum -= 3 - cards[pos]
		cards[pos] = 0
	} else {
		cards[pos] -= 3
	}
	return true, lzNum
}

func scanShunFromPos(cards []int, lzNum int, pos int) (bool, int) {
	if pos < 0 || pos > 6 {
		return false, lzNum
	}
	if cards[pos] > 0 && cards[pos+1] == 0 && cards[pos+2] == 0 {
		return false, lzNum
	}
	if cards[pos]+cards[pos+1]+cards[pos+2]+lzNum < 3 {
		return false, lzNum
	}
	for i := pos; i <= pos+2; i++ {
		if cards[i] == 0 {
			lzNum--
		} else {
			cards[i]--
		}
	}
	return true, lzNum
}

func iterateLaizi(cards []int, pos int) {

}

func checkHuWithEye(huTable map[int64][][]int, cards []int, lzNum int) {
	lzFlag := make(map[int]bool)
	lzList := make([]int, 0)
	for i:=0; i<=3;i++ {
		for j:=0; j<9; j++ {
			val := i*9 + j
			p1 := i*9
			p2 := p1 + 9
			if cards[val] > 0 {
				if i == 3 {
					lzFlag[val] = true
				} else {
					for k:=val-lzNum; k<=val+lzNum; k++ {
						if lzFlag[k] == false && k >= p1 && k <= p2 {
							lzFlag[k] = true
						}
					}
				}
			}
		}
	}
	
	leftNum := lzNum
	for i,v := range(lzList) {
		for j:=0; j < leftNum; j++ {
			cards[v] = j
		}
	}
}
