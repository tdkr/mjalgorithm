/*
@Time : 2018/12/10 9:37
@Author : RonanLuo
*/
package hu

import (
	"github.com/eapache/queue"
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
			pos[index] = row*9 + col
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
		//fmt.Printf("AnalyseHuInfo, key %b not found\n", key)
	}
	return
}

func CheckHu(huTable map[int64][][]int, cards []int) [][]int {
	pos := make([]int, 14)
	key := CalcKey(cards, pos)
	result := CalcGroups(huTable, key, pos)
	return result
}

func CheckHuWithLZ(huTable map[int64][][]int, cards []int, lzFlag map[int]bool) [][]int {
	lzNum := 0
	dupCards := make([]int, len(cards))
	copy(dupCards, cards)
	for k, _ := range lzFlag {
		lzNum += cards[k]
		dupCards[k] = 0
	}
	eyeArr := make([]int, 0)
	for val, num := range dupCards {
		if num > 0 {
			eyeArr = append(eyeArr, val)
		}
	}
	if lzNum >= 2 {
		for k := range lzFlag {
			eyeArr = append(eyeArr, k)
		}
	}
	results := make([][]int, 0)
	for _, eye := range eyeArr {
		backNum := lzNum
		repCards, ok := checkEyeWithLz(dupCards, eye, lzNum)
		if !ok {
			continue
		}
		lzNum -= len(repCards)
		ret := queue.New()
		rep := make([]int, 0, lzNum)
		iterateCards(ret, dupCards, lzNum, 0, rep)
		for i := 0; i < ret.Length(); i++ {
			tmpCards := make([]int, len(dupCards))
			copy(tmpCards, dupCards)
			item := ret.Get(i).([]int)
			for _, v := range item {
				tmpCards[v]++
			}
			tmpCards[eye] += 2
			t := CheckHu(huTable, tmpCards)
			if t != nil {
				results = append(results, t...)
			}
		}
		lzNum = backNum
		dupCards[eye] += 2
		for _, v := range repCards {
			dupCards[v]++
		}
	}
	return results
}

func iterateCards(results *queue.Queue, cards []int, lzNum int, pos int, repCards []int) {
	if pos >= len(cards) {
		return
	}
	// fmt.Println("iterateCards", pos, cards[pos], cards, repCards, lzNum)
	if lzNum == 0 {
		results.Add(repCards)
		return
	}
	if cards[pos] == 0 {
		iterateCards(results, cards, lzNum, pos+1, repCards)
		return
	}
	for i := 0; i < 2; i++ {
		n0 := len(repCards)
		blzNum := lzNum
		keNum, shunNum := 0, 0
		if i == 0 {
			if arr, ok := checkKeWithLz(cards, pos, lzNum); ok {
				repCards = append(repCards, arr...)
				lzNum -= len(arr)
				keNum++
			}
			for cards[pos] > 0 {
				if arr, ok := checkShunWithLz(cards, pos, lzNum); ok {
					repCards = append(repCards, arr...)
					lzNum -= len(arr)
					shunNum++
				} else {
					break
				}
			}
		} else {
			for cards[pos] > 0 {
				if arr, ok := checkShunWithLz(cards, pos, lzNum); ok {
					repCards = append(repCards, arr...)
					lzNum -= len(arr)
					shunNum++
				} else {
					break
				}
			}
			if arr, ok := checkKeWithLz(cards, pos, lzNum); ok {
				repCards = append(repCards, arr...)
				lzNum -= len(arr)
				keNum++
			}
		}
		if shunNum+keNum > 0 {
			n1 := len(repCards)
			iterateCards(results, cards, lzNum, pos+1, repCards)
			lzNum = blzNum
			n2 := len(repCards)
			for i := 0; i < n2-n1; i++ {
				repCards = repCards[:len(repCards)-1]
			}
			for i := 0; i < keNum; i++ {
				cards[pos] += 3
			}
			for i := 0; i < shunNum; i++ {
				cards[pos]++
				cards[pos+1]++
				cards[pos+2]++
			}
			// fmt.Println("restore1", repCards, n0, n1, n2, cards)
			for i := 0; i < n1-n0; i++ {
				tail := len(repCards) - 1
				cards[repCards[tail]]--
				repCards = repCards[:tail]
			}
			// fmt.Println("restore2", repCards, cards)
		}
	}
}

func checkEyeWithLz(cards []int, pos int, lzNum int) ([]int, bool) {
	if cards[pos]+lzNum < 2 {
		return nil, false
	}
	ret := make([]int, 0)
	if cards[pos] >= 2 {
		cards[pos] -= 2
	} else {
		for i := 0; i < 2-cards[pos]; i++ {
			ret = append(ret, pos)
		}
		cards[pos] = 0
	}
	return ret, true
}

func checkKeWithLz(cards []int, pos int, lzNum int) ([]int, bool) {
	if cards[pos]+lzNum < 3 {
		return nil, false
	}
	repCards := make([]int, 0)
	if cards[pos] >= 3 {
		cards[pos] -= 3
	} else {
		for i := 0; i < 3-cards[pos]; i++ {
			repCards = append(repCards, pos)
		}
		lzNum -= 3 - cards[pos]
		cards[pos] = 0
	}
	return repCards, true
}

func checkShunWithLz(cards []int, pos int, lzNum int) ([]int, bool) {
	if pos/9 == 3 || pos%9 > 6 {
		return nil, false
	}
	need := 0
	for i := pos; i <= pos+2; i++ {
		if cards[i] == 0 {
			need++
		}
	}
	if need > lzNum {
		return nil, false
	}
	repCards := make([]int, 0)
	for i := pos; i <= pos+2; i++ {
		if cards[i] == 0 {
			repCards = append(repCards, i)
		} else {
			cards[i]--
		}
	}
	return repCards, true
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
