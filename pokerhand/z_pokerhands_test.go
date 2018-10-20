package task01

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

// go test -v

func TestCompare(t *testing.T) {

	const num = 10000
	// 坑呀！！！！ 开头只能大写！！！！！！
	type PAIR struct {
		// board 5 张牌，　Alice Bob　手中各自２张牌
		Board  [5]string `json:"board"`
		Alice  [2]string `json:"alice"`
		Bob    [2]string `json:"bob"`
		Result int       `json:"result"`
	}
	type DATA struct {
		Matches [num]PAIR `json:"matches"`
	}

	// 从　52 张牌随机抽出　９张牌 , index 没重复项
	// math包中的rand!
	generalIndex := func() [9]int {
		var i int = 0
		var index [9]int
		for i < 9 {
			tmp := rand.Intn(52)
			//确定　tmp 不在 index数组中
			if i == 0 {
				//第一次生成的数组直接天剑进去
				index[i] = tmp
				i++
				continue
			} else {
				// 之后生成的随机数，确定tmp不在index数组中
				var flag bool = false
				for j := 0; j < i; j++ {
					if tmp == index[j] {
						flag = true
					}
				}
				if flag == true {
					continue
				} else {
					index[i] = tmp
					i++
				}
			}
		}

		return index
	}
	//fmt.Println(generalIndex())

	type mypair struct {
		alice    ZPokerHands
		bob      ZPokerHands
		myResult int
	}
	var total [num]mypair
	var totalIndex [num][9]int

	// 设置牌
	for i := 0; i < num; i++ {
		tmpIndex := generalIndex()
		totalIndex[i] = tmpIndex
		var aliceHand [2]USHORT
		var bobHand [2]USHORT
		var board [5]USHORT
		var alicehands [14]UCHAR
		var bobhands [14]UCHAR
		board[0] = USHORT(INITCARDS[tmpIndex[0]])
		board[1] = USHORT(INITCARDS[tmpIndex[1]])
		board[2] = USHORT(INITCARDS[tmpIndex[2]])
		board[3] = USHORT(INITCARDS[tmpIndex[3]])
		board[4] = USHORT(INITCARDS[tmpIndex[4]])
		aliceHand[0] = USHORT(INITCARDS[tmpIndex[5]])
		aliceHand[1] = USHORT(INITCARDS[tmpIndex[6]])
		bobHand[0] = USHORT(INITCARDS[tmpIndex[7]])
		bobHand[1] = USHORT(INITCARDS[tmpIndex[8]])

		//fmt.Printf("%x %x %x \n", board[0], aliceHand[0], bobHand[0])
		//fmt.Println(tmpIndex)

		// 一张牌生成两个数, 类似　3, 4, // 0x0304	　index = 28 ; 28/13=2 , 28%13 = 2
		//fmt.Println(tmpIndex[0]/13, " ", tmpIndex[0]%13)
		transform := func(i int) (color, value UCHAR) {
			color = UCHAR(i/13 + 1)
			value = UCHAR(i%13 + 2)
			return
		}
		//fmt.Println(transform(tmpIndex[0]))
		// alice 手牌
		alicehands[0], alicehands[1] = transform(tmpIndex[5])
		alicehands[2], alicehands[3] = transform(tmpIndex[6])
		// bob 手牌
		bobhands[0], bobhands[1] = transform(tmpIndex[7])
		bobhands[2], bobhands[3] = transform(tmpIndex[8])
		// board card
		color, value := transform(tmpIndex[0])
		alicehands[4], alicehands[5] = color, value
		bobhands[4], bobhands[5] = color, value
		color, value = transform(tmpIndex[1])
		alicehands[6], alicehands[7] = color, value
		bobhands[6], bobhands[7] = color, value
		color, value = transform(tmpIndex[2])
		alicehands[8], alicehands[9] = color, value
		bobhands[8], bobhands[9] = color, value
		color, value = transform(tmpIndex[3])
		alicehands[10], alicehands[11] = color, value
		bobhands[10], bobhands[11] = color, value
		color, value = transform(tmpIndex[4])
		alicehands[12], alicehands[13] = color, value
		bobhands[12], bobhands[13] = color, value

		total[i].alice.SetHands(aliceHand[0], aliceHand[1], board)
		total[i].alice.SetHandsUCHAR(alicehands)
		total[i].bob.SetHands(bobHand[0], bobHand[1], board)
		total[i].bob.SetHandsUCHAR(bobhands)

		//fmt.Printf("%x \n", total[i].alice.m_Cards[0].card)
		//fmt.Printf("%x \n", total[i].bob.m_Cards[0].card)
		//fmt.Println(totalIndex[i])
	}

	// 可以把测试开始没有结果的数据写到文件中
	var jsonData DATA

	for i := 0; i < num; i++ {
		singleValue := &total[i] //.alice.m_Cards[0]
		fomatString := func(number USHORT) string {
			return fmt.Sprintf("0x0%x", number)
		}
		jsonSingleValue := &jsonData.Matches[i]
		jsonSingleValue.Alice[0] = fomatString(singleValue.alice.m_Cards[0].card)
		jsonSingleValue.Alice[1] = fomatString(singleValue.alice.m_Cards[1].card)
		jsonSingleValue.Bob[0] = fomatString(singleValue.bob.m_Cards[0].card)
		jsonSingleValue.Bob[1] = fomatString(singleValue.bob.m_Cards[1].card)
		jsonSingleValue.Board[0] = fomatString(singleValue.alice.m_Cards[2].card)
		jsonSingleValue.Board[1] = fomatString(singleValue.alice.m_Cards[3].card)
		jsonSingleValue.Board[2] = fomatString(singleValue.alice.m_Cards[4].card)
		jsonSingleValue.Board[3] = fomatString(singleValue.alice.m_Cards[5].card)
		jsonSingleValue.Board[4] = fomatString(singleValue.alice.m_Cards[6].card)

		//fmt.Println(*jsonSingleValue)
	}
	bytes, _ := json.Marshal(jsonData)
	ioutil.WriteFile("./test_data_without_result.json", bytes, 0666)

	//测试开始　！！！！！
	start := time.Now()
	for i := 0; i < num; i++ {
		mypairPtr := &total[i]
		//fmt.Printf("board: %x %x %x %x %x ; alice: %x %x ; bob: %x %x ;\n",
		//	mypairPtr.alice.m_Cards[2].card, mypairPtr.alice.m_Cards[3].card, mypairPtr.alice.m_Cards[4].card, mypairPtr.alice.m_Cards[5].card, mypairPtr.alice.m_Cards[6].card,
		//	mypairPtr.alice.m_Cards[0].card, mypairPtr.alice.m_Cards[1].card, mypairPtr.bob.m_Cards[0].card, mypairPtr.bob.m_Cards[1].card)

		mypairPtr.alice.CheckHandsType()
		mypairPtr.bob.CheckHandsType()

		// 0：相等　１：t1大于t2 -1:t1小于t2
		mypairPtr.myResult = mypairPtr.alice.compHandsType(mypairPtr.alice.m_nTypeHands, mypairPtr.alice.m_pHands, mypairPtr.bob.m_nTypeHands, mypairPtr.bob.m_pHands)

		//fmt.Printf("alice : %d ; bob : %d \n", mypairPtr.alice.m_nTypeHands, mypairPtr.bob.m_nTypeHands)
		//fmt.Println(mypairPtr.myResult)
		//fmt.Println(totalIndex[i])
	}
	duration := time.Since(start)
	fmt.Printf("test %v compare, time: %v  \n", num, duration)
	// 测试结束！！！！

	// 把结果写到文件中
	for i := 0; i < num; i++ {
		jsonData.Matches[i].Result = total[i].myResult

		//fmt.Println(jsonData.Matches[i])
	}
	bytes, _ = json.Marshal(jsonData)
	ioutil.WriteFile("./test_data_with_result.json", bytes, 0666)

}
