package task01

import "fmt"

// 改bug 调试逻辑　优化代码
// 瑕疵的的地方　：　sortCard()中排序顺序问题，同花排前面，公共牌和手牌排序　与　之后的isStraight()判断同花顺可能有冲突

type UCHAR byte
type USHORT uint //uint16

type SAME_CARDS_COUNT struct {
	pos, count UCHAR
}

type Split struct {
	value, color UCHAR
}
type POKERCARD struct {
	card  USHORT
	split Split
}

// 变量的大小写?
var INITCARDS = [...]int{
	0x0102, 0x0103, 0x0104, 0x0105, 0x0106, 0x0107, 0x0108, 0x0109, 0x010A, 0x010B, 0x010C, 0x010D, 0x010E, //方块
	0x0202, 0x0203, 0x0204, 0x0205, 0x0206, 0x0207, 0x0208, 0x0209, 0x020A, 0x020B, 0x020C, 0x020D, 0x020E, //梅花
	0x0302, 0x0303, 0x0304, 0x0305, 0x0306, 0x0307, 0x0308, 0x0309, 0x030A, 0x030B, 0x030C, 0x030D, 0x030E, //红桃
	0x0402, 0x0403, 0x0404, 0x0405, 0x0406, 0x0407, 0x0408, 0x0409, 0x040A, 0x040B, 0x040C, 0x040D, 0x040E, //黑桃
}

const (
	_                               = iota //占位符
	HAND_CARD_TYPE_PAIR_A                  //手牌为AA
	HAND_CARD_TYPE_FLSUH_AK                //同花AK
	HAND_CARD_TYPE_FLUSH_AQ_OR_AJ          //同花AQ或者AJ
	HAND_CARD_TYPE_UNFLUSH_AK              //非同花的AK
	HAND_CARD_TYPE_PAIR_J_TO_K             //对J到对K
	HAND_CARD_TYPE_UNFLUSH_AQ_OR_AJ        //非同花AQ或者AJ
	HAND_CARD_TYPE_PAIR_2_TO_10            //对2到对10
	HAND_CARD_TYPE_FLUSH_STRAIGHT          //手牌的同花顺子
	HAND_CARD_TYPE_FLUSH                   //手牌的同花
	HAND_CARD_TYPE_STRAIGHT                //手牌的顺子，即为两张点数连续的牌，注意: A,2为顺子，AK也为顺子
	HAND_CARD_TYPE_PIE                     //杂牌　11
)

const (
	_                   = iota //占位符　	０
	TYPE_PIE_CARD              //杂牌　		pie card
	TYPE_HIGH_CARD             //高牌　		high card
	TYPE_PAIR                  //一对　		pair
	TYPE_TWO_PAIRS             //两队　		two pairs
	TYPE_THREE_KIND            //三条　		three of a kind
	TYPE_STRAIGHT              //顺子　		straight
	TYPE_FLUSH                 //同花　		flush
	TYPE_FULL_HOUSE            //葫芦　		full house
	TYPE_FOUR_KIND             //四条　		four of a kind
	TYPE_STRAIGHT_FLUSH        //同花顺　	staight flush
	TYPE_ROYAL_FLUSH           //皇家同花顺　royal flush
)

const (
	TOTAL_OF_CARDS uint = 7
	TOTAL_OF_HANDS uint = 5
)

// PokerHands 一手牌 在德州扑克中　通常7张为一手牌
type ZPokerHands struct {

	//　原始7张牌
	m_Cards [TOTAL_OF_CARDS]POKERCARD
	// 优选出的最大５张牌 , c++这个指针估计后面在写go代码有问题 ???
	// 这里不用指针了
	//m_pHands *POKERCARD
	m_pHands [TOTAL_OF_HANDS]POKERCARD

	// 两张手牌
	m_HandCards [2]POKERCARD

	// 同样牌值计数结构
	m_SameCards [3]SAME_CARDS_COUNT
	// 花色统计数组, 记录７张牌中，４种花色各自数量
	m_ColorCount [4]UCHAR
	// 记录是否同花，０值不是同花，其他表示同花的花色　　 在函数coutFlush()里面设置的，５张牌同花
	m_nFlushColor UCHAR
	// 同样牌值计数器 ???
	m_nSameCount UCHAR
	// 对牌统计　７张牌有时会出现三个对子
	m_nPairCount UCHAR
	// 三条统计　７张牌有时会出现２个三条
	m_nThreeKindCount UCHAR
	// 最优的５张牌的类型
	m_nTypeHands UCHAR
}

// TODO constructor
// 构造函数 待写　！！！！
// ZPokerHands() { initialize(); }
// void initialize(void)

// public:
//////////////////////////////////////
// 对外部调用公开的函数定义

// 获取牌型 the poker hand type
func (z *ZPokerHands) GetHandsType() int {
	return int(z.m_nTypeHands)
}

// TODO ALL TEST

// 检测牌型并优选出最大的５张牌
func (z *ZPokerHands) CheckHandsType() { // (pHand *POKERCARD)
	// 断言是否SetHands设置过牌数据
	// assert(m_Cards[0].p_card != 0)
	if z.m_Cards[0].card == 0 {
		panic(" 断言是否SetHands设置过牌数据 ")
		return
	}

	// 德州扑克算法步骤
	// 这里的排序有问题

	z.coutFlush() // 2. 统计牌中各花色的数量　检测是否有５张同花色的牌
	// 源代码的注释是有问题的　这个　sortCard中会用到同花属性　m_nFlushColor
	z.sortCard() // 1. 将牌从大到小依次排列　共７张！！！！　这里的注释是错的，但是代码是正确的
	// sortCard() 修改源码
	z.countSameCards() // 3. 统计牌中复数牌的数量，例如对子、三条或四条

	// 4. 依次从最大牌型到最小牌型进行检测，并选出最优的５张
	if z.checkStraightFlush() { //　同花顺　royal flush & straight flush
		//fmt.Println("同花顺　royal flush & straight flush")
		return
	}
	if z.checkFourCard() { // 四条　four of a kind
		//fmt.Println("四条　four of a kind")
		return
	}
	if z.checkFullHouse() { // 葫芦  full house
		//fmt.Println("葫芦  full house")
		return
	}
	if z.checkFlush() { // 同花　flush
		//fmt.Println("同花　flush")
		return
	}
	if z.checkStraight() { // 顺子　straight
		//fmt.Println("顺子　straight")
		return
	}
	if z.checkThreeCard() { //　三条　three of a kind
		//fmt.Println("三条　three of a kind")
		return
	}
	if z.checkTwoPair() { // 两对　two pair
		//fmt.Println("两对　two pair")
		return
	}
	if z.checkPair() { // 一对 pair
		//fmt.Println("一对 pair")
		return
	}
	if z.checkHighCard() { //　高牌　high card
		//fmt.Println("高牌　high card")
		return
	}

	// 最差应为TYPE_PIE_CARD , 否则　Big Bug
	panic("最差应为TYPE_PIE_CARD , 否则　Big Bug ！！！")
}

// 这个是判断手上两张牌的类型　好像没啥用
func (z *ZPokerHands) GetHandCardTypeWithCard(card1, card2 USHORT) int {
	var nRet int = 1
	var sColor USHORT = 0xFF00
	var sValue USHORT = 0x00FF
	var sA USHORT = 0x000E //　值为Ａ
	var sJ USHORT = 0x000B // 值为Ｊ
	var s2 USHORT = 0x0002 // 值为２
	var sF USHORT = 0x0100 // 颜色为方块
	var sH USHORT = 0x0400 // 颜色为黑桃

	var bIsFlush bool = false
	var bIsStraight bool = false
	var bIsPair bool = false
	var bIsCardValid = false

	// 先排序　先按值排序　再按花色排序　保证card1一定大于card2
	//tmp := card1
	// 逻辑待搞清楚　！
	if (card1 & sValue) < (card2 & sValue) { // 先按值排序
		card1, card2 = card2, card1
	}
	if (card1 & sValue) == (card2 & sValue) { // 若值相等　按花色排序
		bIsPair = true // 对子
		if (card1 & sColor) < (card2 & sColor) {
			card1, card2 = card2, card1
		}
	}

	// 检测手牌的值是否有效
	if ((card1 & sValue) <= sH) && ((card1 & sValue) <= sA) {
		bIsCardValid = true
	}
	if (bIsCardValid && (card2&sColor) >= sF) && ((card2 & sValue) >= s2) {
		bIsCardValid = true
	}
	if !bIsCardValid {
		return nRet
	}

	// 逻辑又有点绕！
	// 做手牌统计
	if !bIsPair && (card1&sColor) == (card2&sColor) { // 手牌的同花
		bIsFlush = true
	}
	var nDiffer int = int((card1 & sValue) - (card2 & sValue))
	// 手牌的顺子 注意 此处 A与2 算为一个顺子 ?
	if !bIsPair && (nDiffer == 1 || nDiffer == 12) {
		bIsStraight = true
	}
	if !bIsFlush && !bIsStraight && !bIsPair {
		if (card1&sValue) == sA && ((card2&sValue) >= sJ && (card2&sValue) < sA) {
			return HAND_CARD_TYPE_UNFLUSH_AQ_OR_AJ
		}
		return HAND_CARD_TYPE_PIE
	}

	//同花
	if bIsFlush {
		if (card1 & sValue) == sA {
			if card2&sValue == (sA - 1) {
				return HAND_CARD_TYPE_FLSUH_AK
			}
			if (card2&sValue) >= sJ && (card2&sValue) < sA {
				return HAND_CARD_TYPE_FLUSH_AQ_OR_AJ
			}
		}
		if bIsStraight {
			return HAND_CARD_TYPE_FLUSH_STRAIGHT
		}
		return HAND_CARD_TYPE_FLUSH
	}

	if bIsStraight {
		if (card1&sValue) == sA && (card2&sValue) == (sA-1) {
			return HAND_CARD_TYPE_UNFLUSH_AK
		}
		return HAND_CARD_TYPE_STRAIGHT
	}

	if bIsPair {
		if (card1&sValue) >= s2 && (card1&sValue) < sJ {
			return HAND_CARD_TYPE_PAIR_2_TO_10 // 对２到对10
		}
		if (card1&sValue) >= sJ && (card1&sValue) < sA {
			return HAND_CARD_TYPE_PAIR_J_TO_K // 对Ｊ到对Ｋ
		}
		if (card1 & sValue) == sA {
			return HAND_CARD_TYPE_PAIR_A
		}
	}

	return nRet
}

// Initialize()方法并不太重要，　go语言中　var struct 时 default zero value
// 初始化内部所有变量　重复利用对象时可调用
func (z *ZPokerHands) Initialize() {
	z.m_nFlushColor = 0     // 初始化同花标志
	z.m_nSameCount = 0      // 复数牌数量
	z.m_nPairCount = 0      // 对数牌数量
	z.m_nThreeKindCount = 0 // 三条数量
	z.m_nTypeHands = 0      // 牌型，最小牌型应从１开始　TYPE_PIE_CARD
	// 这个还是这指针的问题
	z.m_pHands = [TOTAL_OF_HANDS]POKERCARD{}

	//　初始化　各种数组
	z.m_Cards = [TOTAL_OF_CARDS]POKERCARD{}
	z.m_ColorCount = [4]UCHAR{}
	z.m_SameCards = [3]SAME_CARDS_COUNT{}
}

// 设置牌值
//func (z ZPokerHands) setHands(c1,c2 USHORT, pBoardCard *USHORT) {
func (z *ZPokerHands) SetHands(c1, c2 USHORT, pBoardCard [TOTAL_OF_HANDS]USHORT) { // 不传指针　直接传数组
	z.m_Cards[0].card = c1
	z.m_Cards[1].card = c2
	z.m_Cards[2].card = pBoardCard[0]
	z.m_Cards[3].card = pBoardCard[1]
	z.m_Cards[4].card = pBoardCard[2]
	z.m_Cards[5].card = pBoardCard[3]
	z.m_Cards[6].card = pBoardCard[4]

	// 记录手牌
	z.m_HandCards[0].card = c1
	z.m_HandCards[1].card = c2
}

// go 没有函数重载呀　改名字了
func (z *ZPokerHands) SetHandsUCHAR(hands [TOTAL_OF_CARDS * 2]UCHAR) {
	// hands 数组越界的　error 检查
	for i := 0; i < int(TOTAL_OF_CARDS); i++ {
		z.m_Cards[i].split.color = hands[i*2]
		z.m_Cards[i].split.value = hands[i*2+1]

		// 断言牌值范围　２－１４
		// assert(m_Cards[i].p_value >= 2 && m_Cards[i].p_value <= 14);
		// 断言牌面花色　1 - 4 黑红梅方
		// assert(m_Cards[i].p_value >= 1 && m_Cards[i].p_color <= 4);

		// 投机取巧,嘿嘿
		if !(z.m_Cards[i].split.value >= 2 && z.m_Cards[i].split.value <= 14) {
			panic("牌值范围　２－１４")
		}
		if !(z.m_Cards[i].split.color >= 1 && z.m_Cards[i].split.color <= 4) {
			panic("牌面花色　1 - 4 黑红梅方")
		}
	}
}

// define COMP_VALUE (a,b) ((a) == (b) ? 0 : ((a)<(b) ? -1 : 1))
/*
func comp_value(a, b UCHAR) int {
	if a == b {
		return 0
	}

	if a < b {
		return -1
	} else {
		return 1
	}
}
*/

func (z *ZPokerHands) compHandsType(t1 UCHAR, p1 [TOTAL_OF_HANDS]POKERCARD, t2 UCHAR, p2 [TOTAL_OF_HANDS]POKERCARD) int {
	if t1 > t2 {
		return 1
	}
	if t1 < t2 {
		return -1
	}

	// 改写成内部函数
	comp_value := func(a, b UCHAR) int {
		if a == b {
			return 0
		}

		if a < b {
			return -1
		} else {
			return 1
		}
	}

	//比较第一张
	var nRet int = comp_value(p1[0].split.value, p2[0].split.value)
	// 第一张如果不相等　直接返回
	if nRet != 0 {
		return nRet
	}

	var nPos int = 0
	// 第一张相等的情况下　继续比较第２，３，４，５张牌
	switch t1 {
	case TYPE_HIGH_CARD: // 高牌　high card
	case TYPE_PIE_CARD: // 杂牌　pie card
	case TYPE_FLUSH: // 同花	flush
		nPos = 1
	case TYPE_PAIR: // 一对　pair
		nPos = 2 // 从第３张开始比较
	case TYPE_TWO_PAIRS: // 两对　two pairs
	case TYPE_THREE_KIND: // 三条　Three of a kind
		nPos = 3 // 从第４张牌开始比较
	case TYPE_FULL_HOUSE: // 葫芦　full house
	case TYPE_FOUR_KIND: // 四条　four of a kind
		nPos = 4 // 比较最后一张
	case TYPE_ROYAL_FLUSH: // 皇家同花顺　royal flush
	case TYPE_STRAIGHT_FLUSH: // 同花顺　straight flush
	case TYPE_STRAIGHT: // 顺子　straight
		return nRet
	default:
		//assert(false) // big Bug
		panic("big bug 假装让他们相等")
		return 0
	}
	// 循环比较第２张以后的牌
	for i := nPos; i < int(TOTAL_OF_HANDS) && nRet == 0; i++ {
		nRet = comp_value(p1[i].split.value, p2[i].split.value)
	}
	return nRet

}

// 检测牌型　并优选出最大的５张牌
// void checkHandsType(POKERCARD *pHand)

// int getHandCardTypeWithCard(USHORT card1, USHORT card2)

// 初始化内部所有变量，　重复利用对象时可调用

// 设置牌值

//private:
// 内部辅助函数定义开始

// 基本对的把　源代码中　m_Card 数组元素的类型是 POKERCARD　，　排序值转换了　card　属性　，　ｓplit.color 和　split.value都没换
// 将一组牌从大到小进行排序
func (z *ZPokerHands) sortCard() {
	var max USHORT
	// TOTAL_OF_CARDS uint unmitch
	for i := 0; i < int(TOTAL_OF_CARDS); i++ {
		max = USHORT(i)
		if z.m_Cards[i].card == 0 {
			panic("如果牌非法返回")
			return
		}
		j := i + 1
		for ; j < int(TOTAL_OF_CARDS); j++ {
			if z.m_Cards[j].card == 0 {
				panic("如果牌非法返回")
				continue
			}
			if z.m_Cards[j].split.value > z.m_Cards[max].split.value {
				max = USHORT(j)
				// 有瑕疵！！！！　上面已经将　max设置成j了，所以下面的条件判断　z.m_Cards[j].split.value == z.m_Cards[max].split.value
				// 下面　z.m_Cards[j].split.value == z.m_Cards[max].split.value　必然会相等的！
				// 个人修改成　else
			} else if z.m_nFlushColor != 0 &&
				z.m_Cards[j].split.value == z.m_Cards[max].split.value &&
				z.m_Cards[j].split.color == z.m_nFlushColor { // 有同花者 放前
				//fmt.Println(z.m_nFlushColor)
				max = USHORT(j)
			}
			// 有瑕疵　！！！
			// isStraight() 中　if z.m_Cards[i].split.value != z.m_Cards[i-1].split.value 　满足才设置　nCount 为０
			// 下面暂时去掉　会影响到　上面有同花的顺序
			/*
				if z.m_Cards[j].split.value == z.m_Cards[max].split.value &&
					z.m_Cards[j].card != z.m_HandCards[0].card &&
					z.m_Cards[j].card != z.m_HandCards[1].card { // 公共牌　排	前面
					//fmt.Println("z.m_Cards[j].split.value: ", z.m_Cards[j].split.value, " z.m_Cards[max].split.value: ", z.m_Cards[max].split.value)
					//fmt.Println("j : ", j, " max : ", max)
					//fmt.Println("公共牌　排前面")
					max = USHORT(j)
				}
			*/

		}
		if int(max) != i {

			z.m_Cards[i].card, z.m_Cards[max].card = z.m_Cards[max].card, z.m_Cards[i].card

			z.m_Cards[i].split.color, z.m_Cards[max].split.color = z.m_Cards[max].split.color, z.m_Cards[i].split.color
			z.m_Cards[i].split.value, z.m_Cards[max].split.value = z.m_Cards[max].split.value, z.m_Cards[i].split.value
		}
	}
}

// test_flush()
// 检查一组排中同样花色的牌是否大于5张
func (z *ZPokerHands) coutFlush() {
	//fmt.Println(" 检查一组排中同样花色的牌是否大于5张 ")
	for i := 0; i < int(TOTAL_OF_CARDS); i++ {
		if z.m_Cards[i].card == 0 {
			continue
		}
		z.m_ColorCount[z.m_Cards[i].split.color-1]++
		// UCHAR uint 有问题
		if uint(z.m_ColorCount[z.m_Cards[i].split.color-1]) >= TOTAL_OF_HANDS {
			z.m_nFlushColor = z.m_Cards[i].split.color
			break
		}
	}
}

// 统计一组7张牌中　对子　三条或者四条的数量
func (z *ZPokerHands) countSameCards() {

	i, nCount := 0, 0
	var start UCHAR = 0

	for i = 1; i < int(TOTAL_OF_CARDS); i++ {
		//fmt.Println("countSameCards")
		if z.m_Cards[i].card == 0 {
			panic(" countSameCards panic ")
			continue
		}
		if nCount == 0 {
			nCount = 1
			start = z.m_Cards[i-1].split.value
		}
		//fmt.Println("z.m_Cards[i].split.value: ", z.m_Cards[i].split.value)
		//fmt.Println(start)
		//fmt.Println(nCount)
		if start == z.m_Cards[i].split.value {
			nCount++
			//fmt.Println("nCount : ", nCount)
		} else {
			if nCount >= 2 {
				//fmt.Println("nCount : ", nCount)
				z.setSameCards(i, nCount)
				//nCount = 0	这个　nCount放错位置了，调了半天！！！
			}
			nCount = 0
		}
	}
	if nCount >= 2 {
		z.setSameCards(i, nCount)
	}
}

// 记录对子　三条　四条所在的位置及数量，　方便CountSameCards调用
func (z *ZPokerHands) setSameCards(nPos, nCount int) {
	z.m_SameCards[z.m_nSameCount].pos = UCHAR(nPos - nCount)
	z.m_SameCards[z.m_nSameCount].count = UCHAR(nCount)
	// assert(m_nSameCount < 3); //断言一组牌中不可能超过３组对子
	if z.m_nSameCount >= 3 {
		panic("m_nSameCount >= 3 ")
	}
	z.m_nSameCount++
	switch nCount {
	case 2:
		z.m_nPairCount++
	case 3:
		z.m_nThreeKindCount++
	}
}

// TODO 待test
// 检查是否顺子　bFlush = true 检查时是否考虑花色问题
// 返回　= -1 说明不是顺子
// 返回　>= 4 顺子最右边的牌在CARDS中的位置，　但Ace例外，　ACE是大牌，排在最左边
// int isStraight(bool bFlush = true)
// 没有默认参数！
func (z *ZPokerHands) isStraight(bFlush bool) int {
	i, nCount := 0, 0
	var pos UCHAR = 0
	for i < int(TOTAL_OF_CARDS-1) {
		i++
		//fmt.Println("test")
		if z.m_Cards[i-1].card == 0 {
			panic("isStraight() panic")
			return -1
		}
		if nCount == 0 {
			if bFlush && z.m_Cards[i-1].split.color != z.m_nFlushColor {
				// 如果顺子起始牌的花色是同花　则SKIP
				continue
			}

			nCount = 1
			pos = z.m_Cards[i-1].split.value
		}
		//
		if bFlush && z.m_Cards[i].split.color != z.m_nFlushColor { // 花色是否一样
			continue
		}

		// 如果顺子起始牌　与后而的Ｎ张牌构成顺子，　则COUNT++计数
		//fmt.Println("pos : ", pos)
		//fmt.Println("z.m_Cards[i].split.value：　", z.m_Cards[i].split.value)
		//fmt.Println("nCount: ", nCount)
		if pos == z.m_Cards[i].split.value+UCHAR(nCount) {
			// 检测　A,2,3,4,5
			// if(++nCount == 4 && pos == 5 && m_Cards[0].p_value = 14)
			nCount++

			if nCount == 4 && pos == 5 && z.m_Cards[0].split.value == 14 {
				if !bFlush || z.m_Cards[0].split.color == z.m_nFlushColor {
					return 100 + i
				}
			}
			// 够５张即是顺子
			if nCount >= int(TOTAL_OF_HANDS) {
				return i
			}
		} else {
			// 这里有问题把　，为什么　前后两个值相等　就不用设置为０了？
			// 顺序必须是同花在前！
			// 所以和　前面的 sortCard()的排序方式之间有很强关联
			if z.m_Cards[i].split.value != z.m_Cards[i-1].split.value {
				// if(pos == 14 && m_Cards[i-1].p_value == 5 && m_Cards[i].p_value ==5)
				nCount = 0
			}
			// 问题　！！！！
			// A X (不同色)5 5 4 3 2
			// 在第二个5的时候
			// 这里应该没问题，　同花的５应该排在前面　所以排序那里出错了
		}
	}
	return -1
}

///////////////////////////
// 提取优选牌算法函数定义开始
// 注意!  m_pHands[handpos].card = m_Cards[pos].card

// 提取同花顺子　或顺子　bFlush = true 同花选项　参数默认值！
// 从m_Cards 7张牌中筛选出最大的５张顺子, 如果是A,2,3,4,5 则Ace放在最右边
// 参数　bFlush = true 只筛选同花, bFlush = false 忽略花色
func (z *ZPokerHands) filterStraight(nPos int, bFlush bool) {
	//fmt.Println(" filterStraight() ")
	// handpos 从数组的最右端依次放入优牌
	var handpos int = int(TOTAL_OF_HANDS) - 1

	// nPos > 100 的情况为最小顺子　，　例如　A,2,3,4,5　，　否则为普通顺子
	if nPos > 100 {
		nPos -= 100
		nAcePos := 0
		if bFlush { //如果是同花，　则需循环获取到的同花色的那张Ace
			for i := 0; i < 2; i++ {
				if z.m_Cards[i].split.color == z.m_nFlushColor {
					nAcePos = i
					break
				}
			}
		}
		// 获取　A,2,3,4,5中的Ace牌, 此时Ace为小牌, 放入数组最右边
		z.m_pHands[handpos].card = z.m_Cards[nAcePos].card
		// 所有的　ｃａｒｄ　复制的时候，　都要添加代码的！！！！！！！
		z.m_pHands[handpos].split.color = z.m_Cards[nAcePos].split.color
		z.m_pHands[handpos].split.value = z.m_Cards[nAcePos].split.value

		handpos--
		// go 指针问题
	}
	//开始循环　从右到左提取顺子
	for handpos >= 0 {
		if bFlush { // bFlush = true 取同花顺子
			if z.m_Cards[nPos].split.color != z.m_nFlushColor {
				nPos--
				continue
			}
		} else { // bFlush = false 取普通顺子
			if handpos < int(TOTAL_OF_HANDS)-1 && z.m_pHands[handpos+1].split.value == z.m_Cards[nPos].split.value {
				nPos--
				continue
			}
			// 因为公共牌排前面，　顺子取前面的牌
			if handpos < int(TOTAL_OF_HANDS)-1 && nPos > 0 && z.m_Cards[nPos-1].split.value == z.m_Cards[nPos].split.value {
				nPos--
				continue
			}
		}
		// 提取顺子
		z.m_pHands[handpos].card = z.m_Cards[nPos].card
		z.m_pHands[handpos].split.color = z.m_Cards[nPos].split.color
		z.m_pHands[handpos].split.value = z.m_Cards[nPos].split.value

		handpos--
		nPos--
	}
	//如果顺子的第一张为A, 则为皇家同花顺
	if bFlush && z.m_pHands[0].split.value == 14 {
		z.m_nTypeHands = TYPE_ROYAL_FLUSH
	}
}

//提取四条　four of a kind
func (z *ZPokerHands) filterFourCard(nPos int) {
	//获取４条的４张牌
	z.m_pHands[0].card = z.m_Cards[nPos].card
	z.m_pHands[1].card = z.m_Cards[nPos+1].card
	z.m_pHands[2].card = z.m_Cards[nPos+2].card
	z.m_pHands[3].card = z.m_Cards[nPos+3].card

	z.m_pHands[0].split.color = z.m_Cards[nPos].split.color
	z.m_pHands[1].split.color = z.m_Cards[nPos+1].split.color
	z.m_pHands[2].split.color = z.m_Cards[nPos+2].split.color
	z.m_pHands[3].split.color = z.m_Cards[nPos+3].split.color

	z.m_pHands[0].split.value = z.m_Cards[nPos].split.value
	z.m_pHands[1].split.value = z.m_Cards[nPos+1].split.value
	z.m_pHands[2].split.value = z.m_Cards[nPos+2].split.value
	z.m_pHands[3].split.value = z.m_Cards[nPos+3].split.value

	// nPos = nPos = 0 ? 4 : 0 ;
	// 获取四条牌型的第5张　kickers 起脚牌
	if nPos == 0 {
		nPos = 4
	} else {
		nPos = 0
	}
	z.m_pHands[4].card = z.m_Cards[nPos].card

	z.m_pHands[4].split.color = z.m_pHands[nPos].split.color
	z.m_pHands[4].split.value = z.m_pHands[nPos].split.color
}

// 辅助函数 找出并取出第一组三条，从左至右放入优选数组
func (z *ZPokerHands) filterFirstThreeCard() {
	nPos := -1
	// 找出第一组三条所在的位置
	for i := 0; i < int(z.m_nSameCount); i++ {
		if z.m_SameCards[i].count == 3 {
			nPos = int(z.m_SameCards[i].pos)
			break
		}
	}

	// assert(nPos != -1)
	if nPos != -1 {
		//循环３次　取出三条
		for i := 0; i < 3; i++ {
			z.m_pHands[i].card = z.m_Cards[nPos+i].card

			z.m_pHands[i].split.color = z.m_Cards[nPos+i].split.color
			z.m_pHands[i].split.value = z.m_Cards[nPos+i].split.color
		}
	}
}

// 提取葫芦　full house
func (z *ZPokerHands) filterFullHouse() {
	var nPos int = 0
	//取full house牌型中的第一组三条
	z.filterFirstThreeCard()

	//m_nTreeKindCount > 1 说明是两个三条　否则是三条加两队或一对
	if z.m_nThreeKindCount > 1 {
		nPos = int(z.m_SameCards[1].pos)
		z.m_pHands[3].card = z.m_Cards[nPos].card
		z.m_pHands[4].card = z.m_Cards[nPos+1].card

		z.m_pHands[3].split.color = z.m_Cards[nPos].split.color
		z.m_pHands[3].split.color = z.m_Cards[nPos].split.color
		z.m_pHands[4].split.color = z.m_Cards[nPos].split.color
		z.m_pHands[4].split.color = z.m_Cards[nPos].split.color

	} else {
		//找出对子所在位置
		for i := 0; i < int(z.m_nSameCount); i++ {
			if z.m_SameCards[i].count == 2 {
				nPos = int(z.m_SameCards[i].pos)
				break
			}
		}
		// 取出最大对子的两张牌
		z.m_pHands[3].card = z.m_Cards[nPos].card
		z.m_pHands[4].card = z.m_Cards[nPos+1].card

		z.m_pHands[3].split.color = z.m_Cards[nPos].split.value
		z.m_pHands[3].split.value = z.m_Cards[nPos].split.value
		z.m_pHands[4].split.color = z.m_Cards[nPos].split.value
		z.m_pHands[4].split.value = z.m_Cards[nPos].split.value
	}
}

// 提取同花　flush, 从左至右　从大到小　选出５张同花
func (z *ZPokerHands) filterFlush() {
	nPos := 0

	for i := 0; i < int(TOTAL_OF_CARDS) && nPos < int(TOTAL_OF_HANDS); i++ {
		if z.m_Cards[i].card == 0 {
			continue
		}
		if z.m_Cards[i].split.color == z.m_nFlushColor {
			z.m_pHands[nPos].card = z.m_Cards[i].card

			z.m_pHands[nPos].split.color = z.m_Cards[i].split.color
			z.m_pHands[nPos].split.value = z.m_Cards[i].split.value
			nPos++
		}

	}
}

// 提取三条　three of a kind
func (z *ZPokerHands) filterThreeCard() {
	z.filterFirstThreeCard()

	var ch UCHAR = z.m_pHands[0].split.value
	for i, n := 0, 3; i < int(TOTAL_OF_CARDS) && n < int(TOTAL_OF_HANDS); i++ {
		if z.m_Cards[i].card == 0 {
			continue
		}

		if ch != z.m_Cards[i].split.value {
			z.m_pHands[n].card = z.m_Cards[i].card

			z.m_pHands[n].split.color = z.m_Cards[i].split.color
			z.m_pHands[n].split.value = z.m_Cards[i].split.value
			n++
		}
	}
}

// 提取两对　two pair
func (z *ZPokerHands) filterTwoPair() {
	//	assert(m_SameCards[0].count == 2);
	// 	assert(m_SameCards[1].count == 2);
	if z.m_SameCards[0].count != 2 {
		panic("filterPair , m_SameCards[0].count != 2")
	}
	if z.m_SameCards[1].count != 2 {
		panic("filterPair , m_SameCards[1].count != 2")
	}

	// 取出第一对
	nPos := z.m_SameCards[0].pos
	z.m_pHands[0].card = z.m_Cards[nPos].card
	z.m_pHands[1].card = z.m_Cards[nPos+1].card

	z.m_pHands[0].split.color = z.m_Cards[nPos].split.color
	z.m_pHands[0].split.value = z.m_Cards[nPos].split.value
	z.m_pHands[1].split.color = z.m_Cards[nPos+1].split.color
	z.m_pHands[1].split.value = z.m_Cards[nPos+1].split.value

	// 取出第二对
	nPos = z.m_SameCards[1].pos
	z.m_pHands[2].card = z.m_Cards[nPos].card
	z.m_pHands[3].card = z.m_Cards[nPos+1].card

	z.m_pHands[2].split.color = z.m_Cards[nPos].split.color
	z.m_pHands[2].split.value = z.m_Cards[nPos].split.value
	z.m_pHands[3].split.color = z.m_Cards[nPos+1].split.color
	z.m_pHands[3].split.value = z.m_Cards[nPos+1].split.value

	//取出最大牌　two pair 中的　kickers 起脚牌
	// 意思不太懂　待查 ？？
	for i := 0; i < int(TOTAL_OF_CARDS); i++ {
		if z.m_Cards[i].card != 0 && z.m_Cards[i].split.value != z.m_pHands[0].split.value && z.m_Cards[i].split.value != z.m_pHands[2].split.value {
			z.m_pHands[4].card = z.m_Cards[i].card

			z.m_pHands[4].split.color = z.m_Cards[i].split.color
			z.m_pHands[4].split.value = z.m_Cards[i].split.value
			break
		}
	}
}

// 提取对子　pair
func (z *ZPokerHands) filterPair() {
	if z.m_SameCards[0].count != 2 {
		panic("filterPair , m_SameCards[0].count != 2")
	}
	// 取出第一对
	var nPos int = int(z.m_SameCards[0].pos)
	z.m_pHands[0].card = z.m_Cards[nPos].card
	z.m_pHands[1].card = z.m_Cards[nPos+1].card

	z.m_pHands[0].split.color = z.m_Cards[nPos].split.color
	z.m_pHands[0].split.value = z.m_Cards[nPos].split.value
	z.m_pHands[1].split.color = z.m_Cards[nPos+1].split.color
	z.m_pHands[1].split.value = z.m_Cards[nPos+1].split.value

	//取出最大牌　pair中三张的kickers 起脚牌
	// 还是不懂意思　待查！！！
	for i, n := 0, 2; i < int(TOTAL_OF_CARDS) && n < int(TOTAL_OF_HANDS); i++ {
		if z.m_Cards[i].card != 0 && z.m_Cards[i].split.value != z.m_pHands[0].split.value {
			z.m_pHands[n].card = z.m_Cards[i].card

			z.m_pHands[n].split.color = z.m_Cards[i].split.color
			z.m_pHands[n].split.value = z.m_Cards[i].split.value
			n++
		}
	}
}

// 提取高牌 high card
func (z *ZPokerHands) filterHighCard() {
	//从大到小　　复制前５张牌
	for i := 0; i < int(TOTAL_OF_HANDS); i++ {
		if z.m_Cards[i].card == 0 {
			continue
		}
		z.m_pHands[i].card = z.m_Cards[i].card

		z.m_pHands[i].split.color = z.m_Cards[i].split.color
		z.m_pHands[i].split.value = z.m_Cards[i].split.value
	}
}

/////////////////////////
/// 检测优选牌算法函数定义开始

// 检查是否是同花顺　Royal Flush  (c++代码英文单词拼错了　Royal Fulsh)
func (z *ZPokerHands) checkStraightFlush() bool {
	if z.m_nFlushColor == 0 {
		return false
	}
	nPos := z.isStraight(true)
	//fmt.Println("checkStraightFlush() nPos : ", nPos)
	if nPos != -1 {
		//缺省为普通同花顺子
		z.m_nTypeHands = TYPE_STRAIGHT_FLUSH
		z.filterStraight(nPos, true)
		return true
	}
	return false
}

// 检测是否是四条　four of a kind
func (z *ZPokerHands) checkFourCard() bool {
	// 循环检测是否包括４条　防止类似　A,A,A,5,5,5,5 牌型
	for i := 0; i < int(z.m_nSameCount); i++ {
		if z.m_SameCards[i].count == 4 {
			// assert(m_nSameCount <= 2);
			// 4条的情况下, m_nSameCount 不应该大于２
			if z.m_nSameCount > 2 {
				panic("４条的情况下 m_nSameCount 不应该大于２")
			}
			z.m_nTypeHands = TYPE_FOUR_KIND
			// 调用filter程序取出４条　放入优选数组
			z.filterFourCard(int(z.m_SameCards[i].pos))
			return true
		}
	}
	return false
}

// 检测是否是葫芦　full house
func (z *ZPokerHands) checkFullHouse() bool {
	if z.m_nThreeKindCount < 1 {
		// 没有三条　不能组成　full house
		return false
	}
	if z.m_nThreeKindCount == 1 && z.m_nPairCount < 1 {
		return false
	}

	z.m_nTypeHands = TYPE_FULL_HOUSE
	z.filterFullHouse()
	return true
}

// 检测是否同花　flush
func (z *ZPokerHands) checkFlush() bool {
	if z.m_nFlushColor == 0 {
		return false
	}
	z.m_nTypeHands = TYPE_FLUSH
	z.filterFlush()

	return true
}

// 检测是否是顺子 straight
func (z *ZPokerHands) checkStraight() bool {
	// 非同花设置成 f
	var nPos int = z.isStraight(false)
	if nPos != -1 {
		//缺省为普通同花顺子
		z.m_nTypeHands = TYPE_STRAIGHT
		z.filterStraight(nPos, false)
		return true
	}
	return false
}

// 检测是否是三条　three of a kind
func (z *ZPokerHands) checkThreeCard() bool {
	if z.m_nThreeKindCount < 1 {
		return false
	}
	z.m_nTypeHands = TYPE_THREE_KIND
	z.filterThreeCard()

	return true
}

// 检查是否两对　two pair
func (z *ZPokerHands) checkTwoPair() bool {
	if z.m_nPairCount < 2 {
		return false
	}
	z.m_nTypeHands = TYPE_TWO_PAIRS
	z.filterTwoPair()

	return true
}

// 检查是否对子
func (z *ZPokerHands) checkPair() bool {
	if z.m_nPairCount < 1 {
		return false
	}
	z.m_nTypeHands = TYPE_PAIR
	z.filterPair()

	return true
}

// 检查高牌或杂牌　high card & pie card
func (z *ZPokerHands) checkHighCard() bool {
	z.m_nTypeHands = TYPE_HIGH_CARD
	z.filterHighCard()

	return true
}

// TODO  测试代码编写
func test_compHandsType() {
	fmt.Println()
	fmt.Println(" test_countSameCards(): ")
	fmt.Println("两套牌大小对比")
	var hands1 [14]UCHAR = [14]UCHAR{
		1, 5, // 0x0105	5
		4, 5, // 0x0405	5

		2, 5, // 0x0205	5
		3, 14, //0x030A	A
		3, 11, //0x030B	J
		4, 6, //0x0406	6
		4, 7} //0x0407	7

	var hands2 [14]UCHAR = [14]UCHAR{
		3, 4, // 0x0304	4
		1, 8, // 0x0108	8

		2, 5, // 0x0205	5
		3, 14, //0x030A	A
		3, 11, //0x030B	J
		4, 6, //0x0406	6
		4, 7} //0x0407	7

	var c11, c12 USHORT = 0x0105, 0x0405
	var c21, c22 USHORT = 0x0304, 0x0108
	var pBoardCard [5]USHORT = [5]USHORT{0x0205, 0x030A, 0x030B, 0x0406, 0x0407}

	// 可以直接 new() 构造　ZPokerHands 的指针
	var pokerhands1 ZPokerHands
	// pokerhands1 := new(ZPokerHands)
	pokerhands1.SetHandsUCHAR(hands1)
	pokerhands1.SetHands(c11, c12, pBoardCard)

	var pokerhands2 ZPokerHands
	// pokerhands2 := new(ZPokerHands)
	pokerhands2.SetHandsUCHAR(hands2)
	pokerhands2.SetHands(c21, c22, pBoardCard)

	pokerhands1.CheckHandsType()
	pokerhands2.CheckHandsType()

	// 0：相等　１：t1大于t2 -1:t1小于t2
	i := pokerhands1.compHandsType(pokerhands1.m_nTypeHands, pokerhands1.m_pHands, pokerhands2.m_nTypeHands, pokerhands2.m_pHands)
	fmt.Println("i: ", i)

}

func test_countSameCards() {
	fmt.Println()
	fmt.Println(" test_countSameCards(): ")
	var hands [14]UCHAR = [14]UCHAR{
		1, 8, // 0x0108		8
		4, 11, // 0x040B	J

		2, 8, // 0x0208		8
		2, 2, // 0x0202		2
		2, 9, // 0x0209		9
		3, 8, // 0x0308		8
		4, 14} // 0x040E    A
	// 上下两个顺序是要对上的　！！
	var c1, c2 USHORT = 0x0108, 0x040B
	var pBoardCard [5]USHORT = [5]USHORT{0x0208, 0x0202, 0x0209, 0x0308, 0x040E}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())

	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}

func test_sortCard() {
	fmt.Println()
	fmt.Println(" test_sortCard(): ")
	var hands [14]UCHAR = [14]UCHAR{
		2, 10, // 0x020A	10
		4, 11, // 0x040B	J

		2, 8, // 0x0208		8
		2, 2, // 0x0202		2
		2, 9, // 0x0209		9
		3, 8, // 0x0308		8
		4, 14} // 0x040E    A
	// 上下两个顺序是要对上的　！！
	var c1, c2 USHORT = 0x020A, 0x040B
	var pBoardCard [5]USHORT = [5]USHORT{0x0208, 0x0202, 0x0209, 0x0308, 0x040E}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())

	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}
func test_two_pair() {
	fmt.Println()
	fmt.Println(" test_two_pair(): ")

	var hands [14]UCHAR = [14]UCHAR{
		3, 4, // 0x0304	4
		1, 8, // 0x0108	8

		2, 4, // 0x0204	4
		3, 6, //0x030A	6
		3, 9, //0x0309	9
		4, 6, //0x0406	6
		4, 8} //0x0407	8

	var c1, c2 USHORT = 0x0304, 0x0108
	var pBoardCard [5]USHORT = [5]USHORT{0x0204, 0x030A, 0x0309, 0x0406, 0x0407}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())

	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}
func test_straight() {
	fmt.Println()
	fmt.Println(" test_straight(): ")

	var hands [14]UCHAR = [14]UCHAR{
		3, 4, // 0x0304	4
		1, 8, // 0x0108	8

		2, 5, // 0x0205	5
		3, 14, //0x030A	A
		3, 9, //0x0309	9
		4, 6, //0x0406	6
		4, 7} //0x0407	7

	var c1, c2 USHORT = 0x0304, 0x0108
	var pBoardCard [5]USHORT = [5]USHORT{0x0205, 0x030A, 0x0309, 0x0406, 0x0407}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())

	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}

func test_full_house() {
	fmt.Println()
	fmt.Println(" test_full_house(): ")
	var hands [14]UCHAR = [14]UCHAR{
		1, 8, // 0x0108		8
		4, 11, // 0x040B	J

		2, 8, // 0x0208		8
		2, 2, // 0x0202		2
		2, 11, // 0x020B	J
		3, 8, // 0x0308		8
		4, 14} // 0x040E    A

	var c1, c2 USHORT = 0x0108, 0x040B
	var pBoardCard [5]USHORT = [5]USHORT{0x0208, 0x0202, 0x020B, 0x0308, 0x040E}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	fmt.Println("m_nFlushColor: ", pokerhands.m_nFlushColor) // 打印哪一个颜色是同花
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())
	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}
func test_four_kind() {
	fmt.Println()
	fmt.Println(" test_four_kind(): ")
	var hands [14]UCHAR = [14]UCHAR{
		1, 8, // 0x0108		8
		4, 11, // 0x040B	J

		2, 8, // 0x0208		8
		2, 2, // 0x0202		2
		2, 11, // 0x020B	J
		3, 8, // 0x0308		8
		4, 8} // 0x0408    8

	var c1, c2 USHORT = 0x0108, 0x040B
	var pBoardCard [5]USHORT = [5]USHORT{0x0208, 0x0202, 0x020B, 0x0308, 0x0408}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	fmt.Println("m_nFlushColor: ", pokerhands.m_nFlushColor) // 打印哪一个颜色是同花
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())
	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}

func test_flush() {
	fmt.Println()
	fmt.Println(" test_flush(): ")
	var hands [14]UCHAR = [14]UCHAR{
		2, 14, // 0x020E	A
		4, 11, // 0x040B	J

		2, 8, // 0x0208		8
		2, 2, // 0x0202		2
		2, 9, // 0x0209		9
		2, 8, // 0x0208		8
		4, 5} // 0x0405		5
	var c1, c2 USHORT = 0x020E, 0x040B
	var pBoardCard [5]USHORT = [5]USHORT{0x0208, 0x0202, 0x0209, 0x0208, 0x0405}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	fmt.Println("m_nFlushColor: ", pokerhands.m_nFlushColor) // 打印哪一个颜色是同花
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())
	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}

func test_straight_flush() {
	fmt.Println()
	fmt.Println(" test_straight_flush(): ")
	var hands [14]UCHAR = [14]UCHAR{
		2, 14, // 0x020E	A
		4, 11, // 0x040B	J

		2, 3, // 0x0203		3
		2, 2, // 0x0202		2
		2, 5, // 0x0205		5
		2, 4, // 0x0204		4
		4, 5} // 0x0405		5
	var c1, c2 USHORT = 0x020E, 0x040B
	var pBoardCard [5]USHORT = [5]USHORT{0x0203, 0x0202, 0x0205, 0x0204, 0x0405}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	//fmt.Println("m_nFlushColor: ", pokerhands.m_nFlushColor) // 打印哪一个颜色是同花
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())

	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}

func test_royal_flush() {
	fmt.Println()
	fmt.Println(" test_royal_flush(): ")
	var hands [14]UCHAR = [14]UCHAR{
		2, 14, // 0x020E	A
		2, 11, // 0x020B	J

		1, 3, // 0x0103		3
		2, 12, // 0x020C	Q
		2, 10, // 0x020A	10
		2, 13, // 0x020D	K
		4, 5} // 0x0405		5
	var c1, c2 USHORT = 0x020E, 0x020B
	var pBoardCard [5]USHORT = [5]USHORT{0x0103, 0x020C, 0x020A, 0x020D, 0x0405}

	var pokerhands ZPokerHands
	pokerhands.SetHandsUCHAR(hands)
	pokerhands.SetHands(c1, c2, pBoardCard)

	pokerhands.CheckHandsType()
	//fmt.Println("m_nFlushColor: ", pokerhands.m_nFlushColor) // 打印哪一个颜色是同花
	fmt.Println("pokerhands.GetHandsType(): ", pokerhands.GetHandsType())

	for index, value := range pokerhands.m_Cards {
		fmt.Print("index: ", index, " : ", value.split.color, " , ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
	fmt.Println("pokerhands.m_pHands : ")
	for _, value := range pokerhands.m_pHands {
		fmt.Print("color: ", value.split.color, " , value: ", value.split.value, " ; ")
		fmt.Printf("%x\n", value.card)
	}
}

func test_GetHandCardTypeWithCard() {
	fmt.Println()
	fmt.Println(" test_GetHandCardTypeWithCard(): ")
	var pokerhands ZPokerHands
	// 手中两张牌的牌型　函数调用
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x030E)) // AA
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x020D)) // 同花AＫ
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x020C)) // 同花ＡＱ
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x020B)) // 同花　ＡＪ
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x030D)) // 非同花　ＡＫ
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020B, 0x030B)) // 对Ｊ
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x030C)) // 非同花ＡＱ
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x0202, 0x0302)) // 对２
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x0207, 0x0208)) // 同花顺子　７，８
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x0202)) // 同花顺子　A，2
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x0207, 0x0209)) // 同花 7 9
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x0302)) // 顺子 A,2 A,K
	fmt.Println(pokerhands.GetHandCardTypeWithCard(0x020E, 0x0307)) // 杂牌
}

func main() {

	test_GetHandCardTypeWithCard()

	test_straight_flush()

	test_royal_flush()

	test_flush()

	test_sortCard()

	test_countSameCards()

	test_straight()

	test_full_house()

	test_four_kind()

	test_two_pair()

	test_compHandsType()

}
