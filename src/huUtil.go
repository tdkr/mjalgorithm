/*
@Time : 2018/12/10 9:37
@Author : RonanLuo
*/
package hu

import (
	"github.com/phf/go-queue/queue"
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
	//fmt.Println("CheckHuWithLZ, cards", cards, lzFlag)
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
		//fmt.Println("CheckHuWithLZ, eye", dupCards, eye, lzNum)
		tmpNum := lzNum
		repCards := queue.New()
		if !checkEyeWithLz(dupCards, eye, lzNum, repCards) {
			continue
		}
		lzNum -= repCards.Len()
		ret := queue.New()
		rep := queue.New()
		iterateCards(ret, dupCards, lzNum, 0, rep)
		for i := 0; i < ret.Len(); i++ {
			tmpCards := make([]int, len(dupCards))
			copy(tmpCards, dupCards)
			item := ret.Get(i).([]int)
			for _, v := range item {
				tmpCards[v]++
			}
			tmpCards[eye] += 2
			t := CheckHu(huTable, tmpCards)
			if t != nil {
				//fmt.Println("CheckHuWithLZ, result", tmpCards, t)
				results = append(results, t...)
			}
		}
		lzNum = tmpNum
		dupCards[eye] += 2
		for i := 0; i < repCards.Len(); i++ {
			v := repCards.Get(i).(int)
			dupCards[v]--
		}
	}
	return results
}

func iterateCards(results *queue.Queue, cards []int, lzNum int, pos int, repCards *queue.Queue) {
	if pos >= len(cards) {
		return
	}
	// fmt.Println("iterateCards", pos, cards[pos], cards, repCards, lzNum)
	if lzNum == 0 {
		t := make([]int, repCards.Len())
		for i := range t {
			t[i] = repCards.Get(i).(int)
		}
		results.PushBack(t)
		return
	}
	if cards[pos] == 0 {
		iterateCards(results, cards, lzNum, pos+1, repCards)
		return
	}

	n0 := repCards.Len()
	dupLz := lzNum

	if checkKeWithLz(cards, pos, lzNum, repCards) {
		n1 := repCards.Len()
		lzNum -= n1 - n0
		iterateCards(results, cards, lzNum, pos, repCards)
		n2 := repCards.Len()
		lzNum = dupLz
		for i := 0; i < n2-n1; i++ {
			repCards.PopBack()
		}
		cards[pos] += 3
		for i := 0; i < n1-n0; i++ {
			v := repCards.PopBack()
			cards[v.(int)]--
		}
	}

	shunNum := 0
	for cards[pos] > 0 {
		if checkShunWithLz(cards, pos, lzNum, repCards) {
			shunNum++
		} else {
			break
		}
	}
	if shunNum > 0 {
		n1 := repCards.Len()
		lzNum -= n1 - n0
		iterateCards(results, cards, lzNum, pos+1, repCards)
		n2 := repCards.Len()
		for i := 0; i < n2-n1; i++ {
			repCards.PopBack()
		}
		p := pos
		if m := pos % 9; m > 6 {
			p -= m - 6
		}
		for i := 0; i < shunNum; i++ {
			cards[p]++
			cards[p+1]++
			cards[p+2]++
		}
		for i := 0; i < n1-n0; i++ {
			v := repCards.PopBack()
			cards[v.(int)]--
		}
	}
}

func checkEyeWithLz(cards []int, pos int, lzNum int, repCards *queue.Queue) bool {
	if cards[pos]+lzNum < 2 {
		return false
	}
	if cards[pos] >= 2 {
		cards[pos] -= 2
	} else {
		for i := 0; i < 2-cards[pos]; i++ {
			repCards.PushBack(pos)
		}
		cards[pos] = 0
	}
	return true
}

func checkKeWithLz(cards []int, pos int, lzNum int, repCards *queue.Queue) bool {
	if cards[pos]+lzNum < 3 {
		return false
	}
	if cards[pos] >= 3 {
		cards[pos] -= 3
	} else {
		for i := 0; i < 3-cards[pos]; i++ {
			repCards.PushBack(pos)
		}
		lzNum -= 3 - cards[pos]
		cards[pos] = 0
	}
	return true
}

func checkShunWithLz(cards []int, pos int, lzNum int, repCards *queue.Queue) bool {
	if pos/9 == 3 {
		return false
	}
	if pos > 6 {
		pos -= pos - 6
	}
	need := 0
	for i := pos; i <= pos+2; i++ {
		if cards[i] == 0 {
			need++
		}
	}
	if need > lzNum {
		return false
	}
	for i := pos; i <= pos+2; i++ {
		if cards[i] == 0 {
			repCards.PushBack(i)
		} else {
			cards[i]--
		}
	}
	return true
}
