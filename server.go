package main

import (
	"fmt"
	"net"
	"sync"
	"io"
)

type Server struct {
	// IP地址,端口
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex
 
	//消息广播
	Message chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}
// 监听Message广播消息channel的goroutine,一旦有消息就发送给全部在线用户
func (this *Server) ListenMessage(){
	for {
		msg := <- this.Message

		// 将msg发送给全部在线用户
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息方法
func (this *Server) BroadCast(user *User, msg string){
	// 发送的消息
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	// 将消息放入管道
	this.Message <- sendMsg
}

// 处理方法
func (this *Server) Handler(conn net.Conn) {
	// 初始化一个用户
	user := NewUser(conn)

	//用户上线, 将用户加入到onLineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线消息
	this.BroadCast(user, "用户上线 ^_^")

	// 接受客户端发送的消息
	go func(){
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)

		if n == 0{
			this.BroadCast(user, "下线 QAQ")
			return
		}

		if err != nil && err != io.EOF{
			fmt.Println("Conn Read err:", err)
			return
		}
		// 读取用户的消息 去除\n
		msg := string(buf[:n-1])

		//将得到的消息进行广播
		this.BroadCast(user, msg)
	}()


	// 当前handler阻塞
	select {}
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err :", err)
		return
	}

	// close listen socket
	defer listener.Close()


	// 启动监听Message的goroutine
	go this.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}

}