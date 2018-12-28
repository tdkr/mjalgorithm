/*
@Time : 2018/12/10 9:37
@Author : RonanLuo
*/
package src

type HuGroup struct {
	EyeList  [][]int
	KeList   [][]int
	ShunList [][]int
}

type iterateData struct {
	huTable  map[int64][][]int
	norCards []int
	lzNum    int
	eye      int
	eyeList  [][]int
	keList   [][]int
	shunList [][]int
	results  []*HuGroup
}

func (data *iterateData) checkEye(pos int) bool {
	if data.norCards[pos]+data.lzNum < 2 {
		return false
	}
	eyeArr := make([]int, 2)
	for i := 0; i < 2; i++ {
		if data.norCards[pos] > 0 {
			data.norCards[pos]--
			eyeArr[i] = pos
		} else {
			data.lzNum--
			eyeArr[i] = -1
		}
	}
	data.eyeList = append(data.eyeList, eyeArr)
	return true
}

func (data *iterateData) checkKe(pos int) bool {
	if data.norCards[pos]+data.lzNum < 3 {
		return false
	}
	arr := make([]int, 3)
	for i := 0; i < 3; i++ {
		if data.norCards[pos] > 0 {
			data.norCards[pos]--
			arr[i] = pos
		} else {
			data.lzNum--
			arr[i] = -1
		}
	}
	data.keList = append(data.keList, arr)
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
		if data.norCards[i] == 0 {
			need++
		}
	}
	if need > data.lzNum {
		return false
	}
	arr := make([]int, 3)
	for i := 0; i <= 2; i++ {
		val := pos + i
		if data.norCards[val] > 0 {
			data.norCards[val]--
			arr[i] = val
		} else {
			data.lzNum--
			arr[i] = -1
		}
	}
	data.shunList = append(data.shunList, arr)
	return true
}

func (data *iterateData) revertEye() {
	arr := data.eyeList[len(data.eyeList)-1]
	for _, v := range arr {
		if v == -1 {
			data.lzNum++
		} else {
			data.norCards[v]++
		}
	}
	data.eyeList = data.eyeList[:len(data.eyeList)-1]
}

func (data *iterateData) revertKe() {
	arr := data.keList[len(data.keList)-1]
	for _, v := range arr {
		if v == -1 {
			data.lzNum++
		} else {
			data.norCards[v]++
		}
	}
	data.keList = data.keList[:len(data.keList)-1]
}

func (data *iterateData) revertShun() {
	arr := data.shunList[len(data.shunList)-1]
	for _, v := range arr {
		if v == -1 {
			data.lzNum++
		} else {
			data.norCards[v]++
		}
	}
	data.shunList = data.shunList[:len(data.shunList)-1]
}

func (data *iterateData) checkHu() []*HuGroup {
	if data.lzNum%3 != 0 {
		return nil
	}
	data.norCards[data.eye] += 2
	ret := CheckHuGroup(data.huTable, data.norCards)
	data.norCards[data.eye] -= 2
	if ret != nil {
		for _, v := range ret {
			v.EyeList = data.eyeList
			v.KeList = append(v.KeList, data.keList...)
			v.ShunList = append(v.ShunList, data.shunList...)
		}
		return ret
	}
	return nil
}

func CalcKey(cards []int) (key int64, pos []int) {
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

func CalcGroups(data []int, pos []int) *HuGroup {
	eyePos := 0
	eye := pos[data[eyePos]]
	kePos := 1
	keNum := data[kePos]
	shunPos := kePos + keNum + 1
	shunNum := data[shunPos]
	ret := &HuGroup{
		EyeList:  make([][]int, 1),
		KeList:   make([][]int, keNum),
		ShunList: make([][]int, shunNum),
	}
	ret.EyeList[0] = make([]int, 2)
	for i := 0; i < 2; i++ {
		ret.EyeList[0][i] = eye
	}
	for i := 0; i < keNum; i++ {
		ret.KeList[i] = make([]int, 3)
		val := pos[data[kePos+i+1]]
		for j := 0; j < 3; j++ {
			ret.KeList[i][j] = val
		}
	}
	for i := 0; i < shunNum; i++ {
		ret.ShunList[i] = make([]int, 3)
		val := pos[data[shunPos+i+1]]
		for j := 0; j < 3; j++ {
			ret.ShunList[i][j] = val + j
		}
	}
	// fmt.Println("CalcGroups", eye, keNum, shunNum, ret)
	return ret
}

func CheckHuGroup(huTable map[int64][][]int, cards []int) []*HuGroup {
	key, pos := CalcKey(cards)
	if data, ok := huTable[key]; ok {
		ret := make([]*HuGroup, len(data))
		for i, v := range data {
			ret[i] = CalcGroups(v, pos)
		}
		return ret
	}
	return nil
}

func CheckHuWithLZ(huTable map[int64][][]int, cards []int, lzList []int, lzFlag map[int]bool) []*HuGroup {
	//fmt.Println("CheckHuWithLZ, cards", cards, lzFlag)
	lzNum := len(lzList)
	eyeArr := make([]int, 0)
	for val, num := range cards {
		if num > 0 || lzFlag[val] {
			eyeArr = append(eyeArr, val)
		}
	}
	results := make([]*HuGroup, 0)
	for _, eye := range eyeArr {
		//fmt.Println("CheckHuWithLZ, eye", cards, eye, lzNum)
		aData := &iterateData{
			eye:      eye,
			huTable:  huTable,
			lzNum:    lzNum,
			norCards: cards,
			shunList: make([][]int, 0),
			keList:   make([][]int, 0),
			eyeList:  make([][]int, 0),
			results:  make([]*HuGroup, 0),
		}
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
	//fmt.Println("iterateCards", pos, data.norCards, data.lzNum)
	if data.lzNum%3 == 0 {
		if ret := data.checkHu(); ret != nil {
			data.results = append(data.results, ret...)
			return
		}
	}
	if pos >= len(data.norCards) {
		return
	}
	if data.norCards[pos] == 0 {
		iterateCards(data, pos+1)
		return
	}

	if data.checkKe(pos) {
		iterateCards(data, pos)
		data.revertKe()
	}

	shunNum := 0
	for data.norCards[pos] > 0 {
		if data.checkShun(pos) {
			shunNum++
		} else {
			break
		}
	}
	if shunNum > 0 {
		if data.norCards[pos] == 0 {
			iterateCards(data, pos+1)
		}
		for i := 0; i < shunNum; i++ {
			data.revertShun()
		}
	}
}
