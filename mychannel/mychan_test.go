package task03

import (
	"fmt"
	"testing"
	"time"
)

// 测试整个文件
// go test -v
// 潜在的一个bug已经测试出来，很有可能是go的bug，稍微抽取一下，可以去github提交这个issue

// 照搬 chan_test.go 中的测试源码 ， 好像在测试的时候，mychan中的异步发送有一个bug，先测出来再说
func TestChan(t *testing.T) {
	//defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	N := 200
	for chanCap := 0; chanCap < N; chanCap++ {
		{
			// Ensure that receive from empty chan blocks.
			c := make(chan int, chanCap)
			recv1 := false
			go func() {
				_ = <-c
				recv1 = true
			}()
			recv2 := false
			go func() {
				_, _ = <-c
				recv2 = true
			}()
			time.Sleep(time.Millisecond)
			if recv1 || recv2 {
				t.Fatalf("chan[%d]: receive from empty chan", chanCap)
			}
			// Ensure that non-blocking receive does not block.
			select {
			case _ = <-c:
				t.Fatalf("chan[%d]: receive from empty chan", chanCap)
			default:
			}
			select {
			case _, _ = <-c:
				t.Fatalf("chan[%d]: receive from empty chan", chanCap)
			default:
			}
			c <- 0
			c <- 0
		}
		{
			// Ensure that send to full chan blocks.
			c := make(chan int, chanCap)
			for i := 0; i < chanCap; i++ {
				c <- i
			}
			//sent := uint32(0)
			go func() {
				c <- 0
				//atomic.StoreUint32(&sent, 1)
			}()
			time.Sleep(time.Millisecond)
			//if atomic.LoadUint32(&sent) != 0 {
			//	t.Fatalf("chan[%d]: send to full chan", chanCap)
			//}
			// Ensure that non-blocking send does not block.
			select {
			case c <- 0:
				t.Fatalf("chan[%d]: send to full chan", chanCap)
			default:
			}
			<-c
		}
		{
			// Ensure that we receive 0 from closed chan.
			c := make(chan int, chanCap)
			for i := 0; i < chanCap; i++ {
				c <- i
			}
			close(c)
			for i := 0; i < chanCap; i++ {
				v := <-c
				if v != i {
					t.Fatalf("chan[%d]: received %v, expected %v", chanCap, v, i)
				}
			}
			if v := <-c; v != 0 {
				t.Fatalf("chan[%d]: received %v, expected %v", chanCap, v, 0)
			}
			if v, ok := <-c; v != 0 || ok {
				t.Fatalf("chan[%d]: received %v/%v, expected %v/%v", chanCap, v, ok, 0, false)
			}
		}
		{
			// Ensure that close unblocks receive.
			c := make(chan int, chanCap)
			done := make(chan bool)
			go func() {
				v, ok := <-c
				done <- v == 0 && ok == false
			}()
			time.Sleep(time.Millisecond)
			close(c)
			if !<-done {
				t.Fatalf("chan[%d]: received non zero from closed chan", chanCap)
			}
		}
	}
}

func TestMyChan(t *testing.T) {
	N := 200

	for chanCap := 0; chanCap <= N; chanCap++ {
		{
			// Ensure that receive from empty chan blocks.
			c := makemychan(int64(chanCap))
			recv1 := false
			go func() {
				_, _ = c.receive()
				recv1 = true
			}()
			if recv1 {
				panic("receive from empty chan")
			}
			// 缺少select
			c.send(0)
		}
		{
			// Ensure that send to full chan blocks.
			c := makemychan(int64(chanCap))
			for i := 0; i < chanCap; i++ {
				c.send(i)
			}
			go func() {
				c.send(0)
			}()
			time.Sleep(time.Microsecond)
			// 缺少select
			c.receive()
		}
		{
			// Ensure that we receive 0 from closed chan.
			c := makemychan(int64(chanCap))
			for i := 0; i < chanCap; i++ {
				c.send(i)
			}
			closemychan(c)
			for i := 0; i < chanCap; i++ {
				v, _ := c.receive()
				if v != i {
					t.Fatalf("chan[%d]: received %v, expected %v", chanCap, v, i)
				}
			}
			if v, _ := c.receive(); v != 0 {
				t.Fatalf("chan[%d]: received %v, expected %v", chanCap, v, 0)
			}
			if v, ok := c.receive(); v != 0 || ok {
				t.Fatalf("chan[%d]: received %v/%v, expected %v/%v", chanCap, v, ok, 0, false)

			}
		}
	}
}

// go test -v  -test.run TestMychan
func TestMychan(t *testing.T) {
	fmt.Println("test method : test mychan ")
}

// go test -v  -test.run TestMyStack
func TestMyStack(t *testing.T) {
	var myTestStack myStack
	myTestStack.push(0)
	myTestStack.push(2)

	a, err := myTestStack.pop()
	fmt.Println(a, " : ", err)
	//fmt.Println(myTestStack.pop())
	//fmt.Println(myTestStack.pop())
}

func TestMyQueue2(t *testing.T) {
	var myTestQueue2 myQueue2
	myTestQueue2.enQueue(0)
	myTestQueue2.enQueue(3)
	myTestQueue2.enQueue(6)

	a, has := myTestQueue2.deQueue()
	fmt.Println(has, " : ", a)
	a, has = myTestQueue2.deQueue()
	fmt.Println(has, " : ", a)
	a, has = myTestQueue2.deQueue()
	fmt.Println(has, " : ", a)
	a, has = myTestQueue2.deQueue()
	fmt.Println(has, " : ", a)
}

func TestSynchronousGOchan(t *testing.T) {
	// go 原装
	c := make(chan int)
	for i := 0; i < 10; i++ {
		go func(i int) { c <- i }(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
}

// go test -v  -test.run TestSynchronousMychan
func TestSynchronousMychan(t *testing.T) {
	// 自己实现
	testChan := makemychan(0)
	for i := 0; i < 10; i++ {
		go func(i int) { testChan.send(i) }(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(testChan.receive())
	}
}

func TestAsynchronousGOchan(t *testing.T) {
	c := make(chan int, 3)
	for i := 0; i < 10; i++ {
		go func(i int) { c <- i }(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}

}

func TestAsynchronousMychan(t *testing.T) {
	for i := 0; i < 200; i++ {
		testChan := makemychan(3)
		for i := 0; i < 10; i++ {
			go func(i int) {
				//time.Sleep(100 * time.Millisecond)
				testChan.send(i)
			}(i)
		}
		for i := 0; i < 10; i++ {
			//fmt.Println(testChan.receive()
			//time.Sleep(100 * time.Millisecond)
			// v: 0  v:1  之后好像测出来过一次bug
			_, _ = testChan.receive()
			//v,_ := testChan.receive()
			//fmt.Printf("v : %d \n", v)
		}
	}
}

//
func TestCloseGOchan(t *testing.T) {
	//
	exit := make(chan int)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("goroutine done.")
		close(exit)
	}()
	fmt.Println("main ...")
	v, err := <-exit
	fmt.Println(v, " ", err)

	/*
		// panic
		exit2 := make(chan int)
		go func() {
			fmt.Println("send 1 to exit2.")
			exit2 <- 1
		}()
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("close eixt2.")
			close(exit2)
		}()
		time.Sleep(4 * time.Second)
		v, err = <-exit2
		fmt.Println(v, " ", err)
	*/

	exit3 := make(chan int, 2)
	fmt.Println(cap(exit3), " ", len(exit3))

	for i := 0; i < 5; i++ {
		go func(i int) {
			exit3 <- i
			//fmt.Println(cap(exit3), " ", len(exit3))
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Println(cap(exit3), " ", len(exit3))
	//close(exit3)
	go func() {
		time.Sleep(3 * time.Second)
		close(exit3)
	}()

	time.Sleep(time.Second)

	for i := 0; i < 6; i++ {
		v, err = <-exit3
		fmt.Println("v : ", v, " err : ", err)
	}

}

func TestCloseMychan(t *testing.T) {
	exit := makemychan(0)
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("goroutine done.")
		closemychan(exit)
	}()
	fmt.Println("main ...")
	v, err := exit.receive()
	fmt.Println(v, " ", err)

	// panic
	/*
		exit2 := makemychan(0)
		go func() {
			fmt.Println("send 1 to exit2.")
			exit2.send(1)
		}()
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("close eixt2.")
			closemychan(exit2)
		}()
		time.Sleep(4 * time.Second)
		v, err = exit2.receive()
		fmt.Println(v, " ", err)
	*/

	exit3 := makemychan(2)

	for i := 0; i < 5; i++ {
		go func(i int) {
			exit3.send(i)
			//fmt.Println(cap(exit3), " ", len(exit3))
		}(i)
	}
	time.Sleep(time.Second)
	//closemychan(exit3)
	go func() {
		time.Sleep(3 * time.Second)
		closemychan(exit3)
	}()

	time.Sleep(time.Second)

	for i := 0; i < 6; i++ {
		v, err = exit3.receive()
		fmt.Println("v : ", v, " err : ", err)
	}
}

func TestSelectGOchan(t *testing.T) {
	var chan1 chan int = make(chan int)
	var chan2 chan int = make(chan int)
	var chs = []chan int{chan1, chan2}
	var numbers = []int{1, 2, 3, 4, 5}
	getNumber := func(i int) int {
		fmt.Println("numbers[", i, "]")
		return numbers[i]
	}
	getChan := func(i int) chan int {
		fmt.Println("chs[", i, "]")
		return chs[i]
	}

	select {
	case getChan(0) <- getNumber(2):
		fmt.Println("1th case is selected.")
	case getChan(1) <- getNumber(3):
		fmt.Println("2th case is selected.")
	default:
		fmt.Println("default!.")
	}

	var testSelectChan chan int = make(chan int, 1)
	select {
	case testSelectChan <- 0:
		fmt.Println("this line will be printed")
	}

	var closechan chan int = make(chan int, 2)
	close(closechan)
	select {
	case v, ok := <-closechan:
		fmt.Println("this must be run : ", v, " ok:", ok)
	default:
		fmt.Println("default : close buffered chan ,then receive data in the select")
	}

	var closechan2 chan int = make(chan int)
	close(closechan2)
	select {
	case v, ok := <-closechan:
		fmt.Println("this also must be run : ", v, " ok: ", ok)
	default:
		fmt.Println("default2 :")
	}
}

// go源码test  chan_test.go 中的测试部分代码
func TestSourceCode(t *testing.T) {
	// This test checks that select acts on the state of the channels at one
	// moment in the execution, not over a smeared time window.
	// In the test, one goroutine does:
	//	create c1, c2
	//	make c1 ready for receiving
	//	create second goroutine
	//	make c2 ready for receiving
	//	make c1 no longer ready for receiving (if possible)
	// The second goroutine does a non-blocking select receiving from c1 and c2.
	// From the time the second goroutine is created, at least one of c1 and c2
	// is always ready for receiving, so the select in the second goroutine must
	// always receive from one or the other. It must never execute the default case.
	n := 1000000
	done := make(chan bool, 1)
	for i := 0; i < n; i++ {
		c1 := make(chan int, 1)
		c2 := make(chan int, 1)
		c1 <- 1
		go func() {
			select {
			case <-c1:
			case <-c2:
			default:
				done <- false
				return
			}
			done <- true
		}()
		c2 <- 1
		select {
		case <-c1:
		default:
		}
		if !<-done {
			t.Fatal("no chan is ready")
		}
	}
}

// 下面测试程序没啥意思
// go test -v  -test.run TestChineseWhispers
func TestChineseWhispers(t *testing.T) {
	const n = 100000
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go func(left, right chan int) {
			left <- 1 + <-right
		}(left, right)
		left = right
	}
	go func(c chan int) { c <- 1 }(right)
	fmt.Println(<-leftmost)
}

func TestChineseWhispersUseMychan(t *testing.T) {
	const n = 100000
	leftmost := makemychan(0)
	right := leftmost
	left := leftmost
	for i := 0; i < n; i++ {
		right = makemychan(0)
		go func(left, right *mychan) {
			v, _ := right.receive()
			left.send(1 + v)
		}(left, right)
		left = right
	}
	go func(c *mychan) { c.send(1) }(right)
	fmt.Println(leftmost.receive())
}
