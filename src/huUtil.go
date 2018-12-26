/*
@Time : 2018/12/10 9:37
@Author : RonanLuo
*/
package src

import (
	"github.com/phf/go-queue/queue"
)

type analyseData struct {
	huTable  map[int64][][]int
	cards    []int
	dupCards []int
	repCards *queue.Queue
	result   [][]int
	lzNum    int
	eye      int
}

func (data *analyseData) replace(value int) {
	data.repCards.PushBack(value)
	data.lzNum--
	data.cards[value]++
}

func (data *analyseData) unreplace(num int, sub bool) {
	for i := 0; i < num; i++ {
		v := data.repCards.PopBack()
		data.lzNum++
		if sub {
			val := v.(int)
			data.dupCards[val]--
			data.cards[val]--
		}
	}
}

func (data *analyseData) duplicate() *analyseData {
	t := &analyseData{
		cards:    make([]int, len(data.cards)),
		repCards: queue.New(),
		eye:      data.eye,
		lzNum:    data.lzNum,
	}
	copy(t.cards, data.cards)
	for i := 0; i < data.repCards.Len(); i++ {
		t.repCards.PushBack(data.repCards.Get(i))
	}
	return t
}

func (data *analyseData) isHu() [][]int {
	if data.lzNum%3 != 0 {
		return nil
	}
	return CheckHu(data.huTable, data.cards)
}

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
	//fmt.Println("CheckHu, result", cards, result)
	return result
}

func CheckHuWithLZ(huTable map[int64][][]int, cards []int, lzFlag map[int]bool) [][]int {
	//fmt.Println("CheckHuWithLZ, cards", cards, lzFlag)
	lzNum := 0
	dupCards := make([]int, len(cards))
	copy(dupCards, cards)
	for k := range lzFlag {
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
		leftNum := lzNum
		c1 := make([]int, len(dupCards))
		copy(c1, dupCards)
		c2 := make([]int, len(dupCards))
		copy(c2, dupCards)
		if dupCards[eye] >= 2 {
			c2[eye] -= 2
		} else {
			c1[eye] = 2
			c2[eye] = 0
			leftNum -= 2 - dupCards[eye]
		}
		aData := &analyseData{
			huTable:  huTable,
			lzNum:    leftNum,
			eye:      eye,
			repCards: queue.New(),
			cards:    c1,
			dupCards: c2,
		}
		rq := queue.New()
		iterateCards(rq, aData, 0)
		for i := 0; i < rq.Len(); i++ {
			t := rq.Get(i).(*analyseData)
			results = append(results, t.result...)
			//fmt.Println("CheckHuWithLZ, result", t.cards, t.repCards, t.result)
		}
	}
	return results
}

func iterateCards(results *queue.Queue, data *analyseData, pos int) {
	//fmt.Println("iterateCards", pos, data.cards, data.repCards, data.lzNum)
	if data.lzNum%3 == 0 {
		if ret := data.isHu(); ret != nil {
			dup := data.duplicate()
			dup.result = ret
			//fmt.Println("iterateCards, result", dup.cards, dup.result)
			results.PushBack(dup)
		}
		return
	}
	if pos >= len(data.cards) {
		return
	}
	if data.dupCards[pos] == 0 {
		iterateCards(results, data, pos+1)
		return
	}

	n0 := data.repCards.Len()

	if checkKeWithLz(data, pos) {
		n1 := data.repCards.Len()
		iterateCards(results, data, pos)
		n2 := data.repCards.Len()
		data.dupCards[pos] += 3
		data.unreplace(n2-n1, false)
		data.unreplace(n1-n0, true)
	}

	shunNum := 0
	for data.dupCards[pos] > 0 {
		if checkShunWithLz(data, pos) {
			shunNum++
		} else {
			break
		}
	}
	if shunNum > 0 {
		n1 := data.repCards.Len()
		if data.dupCards[pos] == 0 {
			iterateCards(results, data, pos+1)
		}
		n2 := data.repCards.Len()
		p := pos
		if m := pos % 9; m > 6 {
			p -= m - 6
		}
		for i := 0; i < shunNum; i++ {
			data.dupCards[p]++
			data.dupCards[p+1]++
			data.dupCards[p+2]++
		}
		data.unreplace(n2-n1, false)
		data.unreplace(n1-n0, true)
	}
}

func checkKeWithLz(data *analyseData, pos int) bool {
	if data.dupCards[pos]+data.lzNum < 3 {
		return false
	}
	if data.dupCards[pos] >= 3 {
		data.dupCards[pos] -= 3
	} else {
		need := 3 - data.dupCards[pos]
		if data.cards[pos]+need > 4 {
			return false
		}
		for i := 0; i < need; i++ {
			data.replace(pos)
		}
		data.dupCards[pos] = 0
	}
	return true
}

func checkShunWithLz(data *analyseData, pos int) bool {
	if pos/9 == 3 {
		return false
	}
	if pos > 6 {
		pos -= pos - 6
	}
	need := 0
	for i := pos; i <= pos+2; i++ {
		if data.dupCards[i] == 0 {
			if data.cards[i] >= 4 {
				return false
			} else {
				need++
			}
		}
	}
	if need > data.lzNum {
		return false
	}
	for i := pos; i <= pos+2; i++ {
		if data.dupCards[i] == 0 {
			data.replace(i)
		} else {
			data.dupCards[i]--
		}
	}
	return true
}
