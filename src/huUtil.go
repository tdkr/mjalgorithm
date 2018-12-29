/*
@Time : 2018/12/10 9:37
@Author : RonanLuo
*/
package src

type iterateData struct {
	huTable   map[int64][][]int
	origCards []int
	iterCards []int
	lzNum     int
	eye       int
	keList    []int
	shunList  []int
	results   [][]int
}

func (data *iterateData) checkEye(pos int) bool {
	if data.iterCards[pos]+data.lzNum < 2 {
		return false
	}
	for i := 0; i < 2; i++ {
		if data.iterCards[pos] > 0 {
			data.iterCards[pos]--
		} else {
			data.lzNum--
		}
	}
	data.eye = pos
	return true
}

func (data *iterateData) checkKe(pos int) bool {
	if data.iterCards[pos]+data.lzNum < 3 {
		return false
	}
	for i := 0; i < 3; i++ {
		if data.iterCards[pos] > 0 {
			data.iterCards[pos]--
		} else {
			data.lzNum--
		}
	}
	data.keList = append(data.keList, pos)
	return true
}

func (data *iterateData) checkShun(pos int) bool {
	if pos/9 == 3 {
		return false
	}
	if pos > 6 {
		pos -= pos - 6
	}
	need := 0
	for i := pos; i <= pos+2; i++ {
		if data.iterCards[i] == 0 {
			need++
		}
	}
	if need > data.lzNum {
		return false
	}
	for i := 0; i <= 2; i++ {
		val := pos + i
		if data.iterCards[val] > 0 {
			data.iterCards[val]--
		} else {
			data.lzNum--
		}
	}
	data.shunList = append(data.shunList, pos)
	return true
}

func (data *iterateData) revertEye() {
	val := data.eye
	if data.iterCards[val]+2 <= data.origCards[val] {
		data.iterCards[val] += 2
	} else {
		data.lzNum += data.iterCards[val] + 2 - data.origCards[val]
		data.iterCards[val] = data.origCards[val]
	}
	data.eye = -1
}

func (data *iterateData) revertKe() {
	index := len(data.keList) - 1
	val := data.keList[index]
	if data.iterCards[val]+3 <= data.origCards[val] {
		data.iterCards[val] += 3
	} else {
		data.lzNum += data.iterCards[val] + 3 - data.origCards[val]
		data.iterCards[val] = data.origCards[val]
	}
	data.keList = append(data.keList[:index], data.keList[index+1:]...)
}

func (data *iterateData) revertShun() {
	index := len(data.shunList) - 1
	val := data.shunList[index]
	for i := 0; i < 3; i++ {
		v := val + i
		if data.iterCards[v] < data.origCards[v] {
			data.iterCards[v]++
		} else {
			data.lzNum++
		}
	}
	data.shunList = append(data.shunList[:index], data.shunList[index+1:]...)
}

func (data *iterateData) checkHu() [][]int {
	if data.lzNum%3 != 0 {
		return nil
	}
	data.iterCards[data.eye] += 2
	arr, ok := CheckHu(data.huTable, data.iterCards)
	data.iterCards[data.eye] -= 2
	if ok {
		ret := make([][]int, len(arr))
		for i, v := range arr {
			ke1, ke2 := v[1], len(data.keList)
			sh1, sh2 := v[1+ke1+1], len(data.shunList)
			t := make([]int, 3+ke1+ke2+sh1+sh2)
			t[0] = v[0]
			//刻子赋值
			t[1] = ke1 + ke2
			p1, p2 := 2, 2
			if ke1 > 0 {
				copy(t[p2:], v[p1:p1+ke1])
			}
			if ke2 > 0 {
				copy(t[p2+ke1:], data.keList)
			}
			//顺子赋值
			t[p2+ke1+ke2] = sh1 + sh2
			p1, p2 = p1+ke1+1, p2+ke1+ke2+1
			if sh1 > 0 {
				copy(t[p2:], v[p1:])
			}
			if sh2 > 0 {
				copy(t[p2+sh1:], data.shunList)
			}
			ret[i] = t
		}
		return ret
	}
	return nil
}

func calcKey(cards []int) (key int64, pos []int) {
	isContinue := false /*是否连续*/
	index := 0
	pos = make([]int, 14)
	for i := 0; i < 34; i++ {
		row := i / 9
		col := i % 9
		num := cards[i]
		if num == 0 {
			continue
		}
		isContinue = row < 3 && col > 0 && cards[i-1] > 0
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
		pos[index] = i
		index++
	}
	return
}

func calcGroups(data []int, pos []int) []int {
	ret := make([]int, len(data))
	eyePos := 0
	ret[eyePos] = pos[data[eyePos]]
	kePos := 1
	keNum := data[kePos]
	ret[kePos] = keNum
	shunPos := kePos + keNum + 1
	shunNum := data[shunPos]
	ret[shunPos] = shunNum
	for i := 0; i < keNum; i++ {
		ret[kePos+1+i] = pos[data[kePos+1+i]]
	}
	for i := 0; i < shunNum; i++ {
		ret[shunPos+1+i] = pos[data[shunPos+1+i]]
	}
	// fmt.Println("calcGroups", eye, keNum, shunNum, ret)
	return ret
}

func CheckHu(huTable map[int64][][]int, cards []int) ([][]int, bool) {
	key, pos := calcKey(cards)
	if arr, ok := huTable[key]; ok {
		retArr := make([][]int, len(arr))
		for i, v := range arr {
			ret := calcGroups(v, pos)
			retArr[i] = ret
		}
		return retArr, true
	}
	return nil, false
}

func CheckHuWithLZ(huTable map[int64][][]int, cards []int, lzList []int, lzFlag map[int]bool) [][]int {
	//fmt.Println("CheckHuWithLZ, cards", cards, lzFlag)
	lzNum := len(lzList)
	eyeArr := make([]int, 0)
	for val, num := range cards {
		if num > 0 || lzFlag[val] {
			eyeArr = append(eyeArr, val)
		}
	}
	results := make([][]int, 0)
	for _, eye := range eyeArr {
		//fmt.Println("CheckHuWithLZ, eye", cards, eye, lzNum)
		aData := &iterateData{
			huTable:   huTable,
			lzNum:     lzNum,
			origCards: cards,
			iterCards: make([]int, len(cards)),
			shunList:  make([]int, 0),
			keList:    make([]int, 0),
			results:   make([][]int, 0),
		}
		copy(aData.iterCards, aData.origCards)
		if aData.checkEye(eye) {
			iterateCards(aData, 0)
			results = append(results, aData.results...)
			aData.revertEye()
			//fmt.Println("CheckHuWithLZ, eye2", cards, aData.lzNum)
		}
	}
	return results
}

func iterateCards(data *iterateData, pos int) {
	//fmt.Println("iterateCards", pos, data.iterCards, data.lzNum)
	if data.lzNum%3 == 0 {
		if ret := data.checkHu(); ret != nil {
			data.results = append(data.results, ret...)
			return
		}
	}
	if pos >= len(data.iterCards) {
		return
	}
	if data.iterCards[pos] == 0 {
		iterateCards(data, pos+1)
		return
	}

	if data.checkKe(pos) {
		iterateCards(data, pos)
		data.revertKe()
	}

	shunNum := 0
	for data.iterCards[pos] > 0 {
		if data.checkShun(pos) {
			shunNum++
		} else {
			break
		}
	}
	if shunNum > 0 {
		if data.iterCards[pos] == 0 {
			iterateCards(data, pos+1)
		}
		for i := 0; i < shunNum; i++ {
			data.revertShun()
		}
	}
}
