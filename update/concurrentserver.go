// 05_concurrentserver
package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// 重要！　服务器在处理客户端输入的命令时，环境是linux，换行符为一个字符/n,如果是windows下，换行符可能为/r/n，两个字符！！！！
// 重要！　无锁，待修改!
// 有bug　在等到客户端输入的时候，　客户端ctrl+c终止，到时服务器挂掉？

//nc 127.0.0.1 8000 连接到服务器
// who 查询在线client
//

// Name 默认设置　成string
type Client struct {
	C      chan string //用于发送数据的管道
	Name   string      //用户名
	Addr   string
	status bool // true 为在线， flase 为不在线
}

// 设计选择问题　在Client中设置登录状态　or 设置在线，所有登录过的用户集合
//保存在线用户cliAddr =====> Client
// 一定要设置成　struct的指针类型　否则无法修改值
var onlineMap map[string]*Client

// 在线消息发送，用于广播模式发送
var message = make(chan string)

// 离线消息存储
var offlineMsg map[string][]string = make(map[string][]string)

func WriteMsgToClient(cli *Client, conn net.Conn) {
	for msg := range cli.C {
		//给当前客户端发送信息
		conn.Write([]byte(msg + "\n"))
	}
}

func MakeMsg(cli *Client, msg string) (buf string) {
	buf = "(" + time_now() + ")" + "[" + cli.Addr + "]" + cli.Name + " : " + msg
	return
}

//处理用户连接
func HandleConn(conn net.Conn) {
	defer conn.Close()
	//获取客户端的网络地址
	cliAddr := conn.RemoteAddr().String()
	// 客户端名字默认为 网络端口地址
	cliName := make([]byte, 2048)
	//需要连接用户输入用户名
	conn.Write([]byte("请输入用户名"))

	// 阻塞等待用户输入用户名
	// 每一次登录，需要输入账号名
	n, _ := conn.Read(cliName)
	// 用户名格式检测
	// linux 回车键算一个字符　所以空字符长度为１
	for n == 1 {
		//isQuit <- true
		//对方断开　或者出问题
		conn.Write([]byte("请输入有效用户名，不能为空："))
		n, _ = conn.Read(cliName)
	}

	//linux环境下　输入最后一个字符是换行符/n,截取用户名！！！　注意windows环境下是 /r/n
	clientName := string(cliName[:n-1])

	//创建一个结构体 默认　用户名和网络地址一样
	cli := &Client{make(chan string), clientName, cliAddr, true}

	//新开一个　ｇｏｒｏｕｔｉｎｅ　专门给客户端发送信息
	// 注意，先开goroutine，才可以往cli.C中传数据,否则直接往cli.C传多个数据中会阻塞
	go WriteMsgToClient(cli, conn)

	//登录，如果用户　曾经登录过，不许要往onlineMap中添加数据
	//onlineMap[cliName]
	_, exist := onlineMap[clientName]

	if !exist {
		//把结构体添加到map
		onlineMap[clientName] = cli
	} else {
		// 用户曾经登录过聊天室
		delete(onlineMap, clientName)
		onlineMap[clientName] = cli
		cli.C <- MakeMsg(cli, "There are "+fmt.Sprint(len(offlineMsg[clientName])-1)+" messages that you haven't seen when you leave.")
		// 需要实现离线发送 把用户离线的这段时间收到的信息转发给用户
		for i := 1; i < len(offlineMsg[clientName]); i++ {
			cli.C <- offlineMsg[clientName][i]
		}

		delete(offlineMsg, clientName)
	}

	//广播某个人登录，在线
	message <- MakeMsg(cli, "login")
	//提示　我是谁
	cli.C <- MakeMsg(cli, "Login , my name is "+cli.Name)

	//对方是否主动退出
	isQuit := make(chan bool)
	hasData := make(chan bool) //对方是否有数据发送

	//新开一个　ｇｏｒｏｕｔｉｎｅ　接收用户发过来的数据
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				isQuit <- true
				//对方断开　或者出问题
				fmt.Println("conn.Read err = ", err)
				return
			}
			msg := string(buf[:n-1]) //通过　nc测试　多一个换行

			// who 查询在线人员　　ａll 查询所有在线离线人员
			if len(msg) == 3 && msg == "who" {
				i := 1
				//遍历map，给当前用户发送所有成员
				//conn.Write([]byte("user list: \n"))
				for _, tmp := range onlineMap {
					// 判断在线用户
					if tmp.status {
						msg = fmt.Sprint(i) + ") " + tmp.Addr + " : " + tmp.Name + " : " + strconv.FormatBool(tmp.status) + " \n"
						conn.Write([]byte(msg))
						i++
					}
				}
			} else if len(msg) == 3 && msg == "all" {
				i := 1
				//遍历map，给当前用户发送所有成员
				//conn.Write([]byte("user list: \n"))
				for _, tmp := range onlineMap {
					// 打印所有　onlineMap 在线，不在线都打印
					msg = fmt.Sprint(i) + ") " + tmp.Addr + " : " + tmp.Name + " : " + strconv.FormatBool(tmp.status) + " \n"
					conn.Write([]byte(msg))
					i++

				}
			} else if len(msg) == 4 && msg == "exit" {
				// 客户端输入　exit主动退出
				fmt.Println(cli.Addr+" : "+cli.Name, " exit")
				conn.Close()
			} else if len(msg) >= 8 && msg[:6] == "rename" { // 修改名字有bug
				//rename |mike
				//name := strings.Split(msg, "|")[1]
				//cli.Name = name
				//onlineMap[cliName] = cli
				//conn.Write([]byte("rename ok \n"))
			} else {
				//转发此内容
				message <- MakeMsg(cli, msg)
			}

			hasData <- true //代表有数据

		}
	}()
	for {
		//通过select 检测channel的流动
		select {
		case <-isQuit:
			//当前用户从map移除
			//delete(onlineMap, cliAddr)
			//不删除用户　设置登录状态
			onlineMap[cli.Name].status = false
			message <- MakeMsg(cli, " logout ") //广播谁下线了

			// 此时有用户下线了
			// 需要给下线的用户创建数据结构，用来存储其他用户离线发送的数据
			offlineMsg[cli.Name] = make([]string, 0)

			return
		case <-hasData:

		case <-time.After(60 * 60 * time.Second): //60s后超时了
			delete(onlineMap, cliAddr)
			message <- MakeMsg(cli, " time out leave out ") //广播谁下线了
			return
		}
	}
}

func Manager() {
	//给map　分派空间
	onlineMap = make(map[string]*Client)
	for {
		msg := <-message //没有消息前　这里会阻塞
		for _, cli := range onlineMap {
			// 如果 onlineMap中的用户状态为true则发送消息
			if cli.status {
				cli.C <- msg
			} else {
				// 如果　onlineMap中的用户状态为false，则实现离线发送，用户下次上线接受消息
				// TODO 待实现离线发送功能
				if msg == "\n" {
					fmt.Println("enter enter enter")
				}
				// 第一项　自己退出的消息也会添加进去
				offlineMsg[cli.Name] = append(offlineMsg[cli.Name], msg)
				fmt.Println(offlineMsg[cli.Name])
			}

		}
	}
}

func main() {
	fmt.Println(time_now())

	//监听
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("net.listen err = ", err)
		return
	}

	defer listener.Close()

	//新开一个　goroutine 转发消息　只要有消息来了　遍历map 给map每个成员都发送次消息
	go Manager()

	// 循环阻塞等待用户连接
	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			fmt.Println(" listener.accept err1 = ", err1)
			continue
		}
		go HandleConn(conn) //处理用户连接
	}

	fmt.Println("Hello World!")
}

func time_now() string {
	date_slice := make([]string, 0)
	time_slice := make([]string, 0)
	a, b, c := time.Now().Date()
	d, e, f := time.Now().Clock()
	date_slice = append(date_slice, fmt.Sprint(a), fmt.Sprint(b), fmt.Sprint(c))
	time_slice = append(time_slice, fmt.Sprint(d), fmt.Sprint(e), fmt.Sprint(f))
	date_parse := strings.Join(date_slice, "-")
	time_parse := strings.Join(time_slice, ":")

	return (date_parse + " " + time_parse)
}
