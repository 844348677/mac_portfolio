package task03

import (
	"fmt"
	"sync"
	//"time"
)

// my chan 个人另类实现 chan.go　功能
// 抽取　简化
/*
	1. chan 类型结构，　构造　chan 的方法
	２． receive send 方法，　收发信息
	３．　同步　异步方式　；　buffered,unbuffered
	(目前类型只支持 int ，　通用类型太底层，　估计比较繁琐)
	4. close()
	5. len() 方法
	(前五项以功能实现，进行了简单测试）
	6． select() 待实现 未完成！
	7. 类型匹配　估计把go源码中的抽取出来都很困难
	go mutex 有一点需要注意!!!!!!!!!!!
*/

/*
// runtime/type.go
type chantype struct {
	typ  _type
	elem *_type
	dir  uintptr
}
// 实现通用的　_type　struct 略难
*/

// 目前只支持 int 类型吧
type mstack interface {
	push(int) bool
	pop() (int, bool)
}

// 栈中数据的输入输出　同步问题
type myStack struct {
	length uint
	data   []int
}

// 目前栈的容量设置为　unlimited , 所有数据都能压入栈中，　返回true
func (ms *myStack) push(singleData int) bool {
	ms.data = append(ms.data, singleData)
	ms.length++
	return true
}
func (ms *myStack) pop() (int, bool) {
	if ms.length == 0 {
		return 0, false
	}
	tmp := ms.data[ms.length-1]
	ms.data = ms.data[0 : ms.length-1]
	ms.length--
	return tmp, true
}
func (ms myStack) isEmpty() bool {
	if ms.length == 0 {
		return true
	}
	return false
}

// Queue : Insertion and Deletion happends on different ends, FIFO
// 使用两个栈实现队列
type mqueue interface {
	enQueue(int) bool
	deQueue() (int, bool)
}

// 队列的实现　两种模式
// model1 : By making enQueue operation costly
// model2 : By making deQueue operation costly , Method 2 is definitely better than method 1.
type myQueue1 struct {
	stack1 *myStack
	stack2 *myStack
}

/*
enQueue(q, x)
  1) While stack1 is not empty, push everything from stack1 to stack2.
  2) Push x to stack1 (assuming size of stacks is unlimited).
  3) Push everything back to stack1.

dnQueue(q)
  1) If stack1 is empty then error
  2) Pop an item from stack1 and return it
*/

type myQueue2 struct {
	lock   sync.Mutex
	length int
	stack1 myStack
	stack2 myStack
}

/*
enQueue(q,  x)
  1) Push x to stack1 (assuming size of stacks is unlimited).

deQueue(q)
  1) If both stacks are empty then error.
  2) If stack2 is empty
       While stack1 is not empty, push everything from stack1 to stack2.
  3) Pop the element from stack2 and return it.
*/
func (mq *myQueue2) enQueue(singleData int) bool {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.stack1.push(singleData)
	mq.length++
	return true
}
func (mq *myQueue2) deQueue() (int, bool) {
	mq.lock.Lock()
	mq.lock.Unlock()
	if mq.stack1.isEmpty() && mq.stack2.isEmpty() {
		return 0, false
	} else if mq.stack2.isEmpty() {
		for !mq.stack1.isEmpty() {
			tmp, has := mq.stack1.pop()
			if !has {
				fmt.Errorf(" error in deQueue ! no more data")
			}
			mq.stack2.push(tmp)
		}
	}
	tmp, has := mq.stack2.pop()
	mq.length--
	return tmp, has
}

// 缓冲数据使用自建的Queue结构
type mychan struct {
	// 缓冲队列
	bufQueue myQueue2
	// 缓冲队列大小　可存储数据项数量
	dataQueueSize uint

	// 是否关闭
	closed uint32

	unbufferedData int
	publicMutex    sync.Mutex // 公共锁
	sendMutex      sync.Mutex // send 锁 unbuffered时使用
	receiveMutex   sync.Mutex // receive 锁 unbuffered时使用

	sendWG    sync.WaitGroup
	receiveWG sync.WaitGroup

	// unbuffered : hasData是否有数据,waitForData是否等待接收数据
	hasData     bool
	waitForData bool

	// buffered
	isFull         bool
	blockInReceive bool

	// 计数器　buffered send receive 中block的数量
	sendWaitCount    uint
	receiveWaitCount uint
}

func makemychan(size int64) *mychan {
	// 检查 size
	if size < 0 {
		panic("size < 0")
	}

	var c *mychan = new(mychan)

	c.dataQueueSize = uint(size)

	return c
}

// 收发方 G ，　是否? 维护发送和接收者等待队列 这个太底层了

// c <- x : c.send(x)
func (c *mychan) send(data int) {
	if c.dataQueueSize == 0 {
		mychanUnbufferedSend(c, data)
	} else {
		mychanBufferedSend(c, data)
	}

}

// <- c : c.receive()
func (c *mychan) receive() (int, bool) {
	var tmp int
	var err bool
	if c.dataQueueSize == 0 {
		tmp, err = mychanUnbufferedReceive(c)
	} else {
		tmp, err = mychanBufferedReceive(c)
	}

	return tmp, err
}

// 同步发送
func mychanUnbufferedSend(c *mychan, data int) bool {

	if c.dataQueueSize != 0 {
		panic("UnbufferedSend : dataQueueSize must be 0")
	}

	// send锁作用　：　如果当前goroutine　block　会释放公共锁，但不会释放send锁，　不会有另一个goroutine进来，　一直等待运行receive（）方法
	c.sendMutex.Lock()
	c.publicMutex.Lock()
	defer c.sendMutex.Unlock()

	if c.closed != 0 {
		c.publicMutex.Unlock()
		panic(" send to closed mychan ")
	}

	c.unbufferedData = data
	c.hasData = true
	// 如果此时　receive 接收数据 block　了, 通知可以取数据了
	if c.waitForData {
		//fmt.Println("aaa")
		c.receiveWG.Done()
	}

	// send中　block　等待　receive取数据
	c.sendWG.Add(1)
	c.publicMutex.Unlock()
	c.sendWG.Wait()

	// 如果　Wait此时取关闭channel,　在closemychan()中会报panic

	return true
}

func mychanUnbufferedReceive(c *mychan) (int, bool) {

	var tmp int
	// receive锁　: 类似　send锁的偶用
	c.receiveMutex.Lock()
	c.publicMutex.Lock()
	defer c.receiveMutex.Unlock()
	defer c.publicMutex.Unlock()

	if c.closed != 0 {
		return 0, false
	}

	// 此时没有数据　receive　会　block　在这里
	if !c.hasData {
		//fmt.Println("bbb")
		c.waitForData = true
		c.receiveWG.Add(1)
		// 公共锁　解锁
		c.publicMutex.Unlock()
		// 该　goroutine block 等待　发送数据
		c.receiveWG.Wait()
		// receive被唤醒，　只有在已经传数据的时候
		c.publicMutex.Lock()

		// 这里进行唤醒判断，　如果是由closemychan()方法中唤醒，　则直接return返回０值和false
		if c.closed != 0 {
			return 0, false
		}
	}
	// 接下来　的情况　send先调用block了，　或者receive先调用，send唤醒receive，send也block
	//两种情况 都是有数据, send　block
	tmp = c.unbufferedData
	c.hasData = false
	c.waitForData = false
	// 唤醒　send
	c.sendWG.Done()

	return tmp, true
}

// 异步发送
func mychanBufferedSend(c *mychan, data int) bool {
	c.publicMutex.Lock()
	//defer c.publicMutex.Unlock()

	if c.closed != 0 {
		panic(" mychanBufferedSend : send to closed mychan ")
	}

	// 如果队列满了　则block,　直到receive中　取出数据再唤醒
	if c.bufQueue.length == int(c.dataQueueSize) {

		//fmt.Println("c.bufQueue.length : ", c.bufQueue.length)
		c.isFull = true
		c.sendWG.Add(1)
		c.sendWaitCount++
		c.publicMutex.Unlock()
		// 等待　receive 取出数据　再唤醒
		// !!!!! 这里好像有个bug ， 有一次测试的时候出现过，之后但是测不出来了 ！！！！！！！！！！！！！！！！！
		// TODO
		// 尝试将队列变成 线程安全的？？
		// !!!!!!!!!!!!!!!!!!!!!!!!! 这里必须要睡一下 不睡就肯定报错，该问题go内部mutex有关
		//time.Sleep(time.Millisecond)
		// !!!!!!!!!!!!!!!!!!!!!!!!! 相关解释
		// Locking the mutex a second time in the same goroutine seems to be the trigger.
		// time.Sleep(time.Millisecond) ensures that the mutex enters "starvation" mode
		// 上面这个bug修改了，但是我这里这么用还是会有bug
		// 待仔细测试
		// sync/mutex.go
		// 两种模式， normal and starvation
		// 待阅读源码
		// 好像这个就应该算是go语言的一个bug！
		// 待写完select仔细阅读源码

		c.sendWG.Wait()
		c.publicMutex.Lock()

	}

	//fmt.Println(c.bufQueue.length)

	c.bufQueue.enQueue(data)

	// 如果　receive中block等待获取数据
	tmpindex := int(c.receiveWaitCount)

	for i := 0; i < tmpindex; i++ {
		c.receiveWG.Done()
		// 我们自信认为这里通知了，　之后receive一定会被唤醒，　block数量就会减少
		// 如果不在这里减掉的话，　send的goroutine再次抢到锁运行，会过多的运行Done()方法
		// 感觉这里和close() 一起用 有可能会有问题
		c.receiveWaitCount--
	}

	c.publicMutex.Unlock()
	return true
}
func mychanBufferedReceive(c *mychan) (int, bool) {

	var tmp int
	c.publicMutex.Lock()
	//defer c.publicMutex.Unlock()

	// 如果queue中的length为0， 则block住, 等待send中唤醒
	if c.bufQueue.length == 0 {
		// 这个放置的位置非常巧妙
		// 当channel已经关闭了，只有队列中的数量为0的时候才返回0，false
		if c.closed != 0 {
			c.publicMutex.Unlock()
			return 0, false
		}

		c.blockInReceive = true
		c.receiveWG.Add(1)
		c.receiveWaitCount++
		c.publicMutex.Unlock()
		c.receiveWG.Wait()
		c.publicMutex.Lock()

		// 如果此时是由　closemychan()唤醒的，则返回０值和false
		if c.closed != 0 {
			c.publicMutex.Unlock()
			return 0, false
		}
	}

	tmp, hasOne := c.bufQueue.deQueue()

	if !hasOne {
		panic(" queue don't have any data! ")
	}

	// 如果此时　send中满了　block
	//if c.isFull == true {
	//	c.isFull = false
	//	c.sendWG.Done()
	//}

	tmpindex := int(c.sendWaitCount)

	for i := 0; i < tmpindex; i++ {
		c.sendWG.Done()
		// 同上，　我们自信认为这里通知了，　之后send一定会被唤醒，　block数量就会减少
		c.sendWaitCount--
	}
	c.publicMutex.Unlock()
	return tmp, true
}

// close() 方法
func closemychan(c *mychan) {
	// chan关闭，无法将数据再次传给chan，　此时若有block存在send方法中，会panic
	// 有缓冲的chan中，关闭chan，若还有数据，可以读取。直到读完，再取数据会返回error
	// 无缓冲的chan，再次receive会返回error

	if c == nil {
		panic(" panic  closemychan c is nil")
	}
	c.publicMutex.Lock()

	// c.closed 为０是表示没有关闭
	if c.closed != 0 {
		c.publicMutex.Unlock()
		panic("close of closed mychan")
	}
	// 关闭channel
	c.closed = 1
	//fmt.Println("ccc: ", c.waitForData)
	// unbuffered　若receive处于block等待接收数据的时候，　解除阻塞 , 但是会ｒeturn  零值和false
	if c.waitForData {
		c.receiveWG.Done()
	}
	// ubuffered 若send处于block等待有人接收数据的时候，　报panic
	if c.hasData {
		panic(" panic : close whilst sending something to unbuffered channel which is waiting to be received ")
	}

	// buffered
	// 唤醒所有等待的　receive方法
	for i := 0; i < int(c.receiveWaitCount); i++ {
		c.receiveWG.Done()
		c.receiveWaitCount--
	}
	// 如果此时有send处于block　panic
	if c.sendWaitCount != 0 {
		panic("myclose : (buffered) panic: send on closed channel ")
	}
	c.publicMutex.Unlock()
}

func capMychan(c *mychan) int {
	return int(c.dataQueueSize)
}

func lenMychan(c *mychan) int {
	if c.dataQueueSize == 0 {
		return 0
	} else {
		c.publicMutex.Unlock()
		defer c.publicMutex.Unlock()
		return c.bufQueue.length
	}
}

// select 方法

// compiler implements
//
//	select {
//	case c <- v:
//		... foo
//	default:
//		... bar
//	}
//
// as
//
//	if selectnbsend(c, v) {
//		... foo
//	} else {
//		... bar
//	}
//

// compiler implements
//
//	select {
//	case v, ok = <-c:
//		... foo
//	default:
//		... bar
//	}
//
// as
//
//	if c != nil && selectnbrecv2(&v, &ok, c) {
//		... foo
//	} else {
//		... bar
//	}
//

// select() 方法去选择 send数组和receive数组哪个去运行，如果多余一个可执行，则随机选择一个。若没有可以执行的，则 。
//
// 写这个select()，确实有点头大，原因，前面设计不好，导致越来越复杂

// selected表示可以发送，发送的值是value
func selectSend(c *mychan, value int) (selected bool) {
	// 在这里发现，go判断是否能发送或者接受数据，集成在一个方法中了，设置true表示真正send数据，设置false表示能不能send数据
	// 前面设计不好，后面机会遇到麻烦！
	// 如果这里能往chan传值，则返回true，否则返回false

	c.send(value)
	//  因为send方法没有返回false的情况，所以运行完就是true

	return true
}

// selected表示可以是否可以接受数据，使用value接受数据 ， received表示是否mychan closed ？ select chan，若chan关闭 一定会拿到数据
func selectReceive(value int, received *bool, c *mychan) (selected bool) {
	// 如果channel closed，一定能接收到数据，
	// 如果是buffered chanel 其中queue中存在数据，一定能接收到数据

	v, ok := c.receive()
	value = v
	received = &ok

	return true
}

func selectMychan(sendSlice []*mychan, sendValue []int, receiveSlice []*mychan, receiveValue []*int, receiveds []*bool) bool {
	//select 中传入的参数介绍，这里是上层封装
	// select每一个case对应的一对chan和对应的值
	// case c <- v ; 参数chan c，value v，调用send方法
	// case v,ok <- c ; 参数value v，isclosed ok，chan c，调用receive方法

	// 还需要判断是否是default
	// 如果有default，send和receive阻塞，则运行default
	// 如果没有default，send和receive直到有一个能运行，如果都能运行，则随机运行一个

	// 这里资源竞争的问题很严重
	// 如果在外边设置判断能发送，再去发送的话，外边判断能发送，但调用发送的方法时锁的竞争导致阻塞。接收同理
	//

	return false
}
