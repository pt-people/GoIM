package main

import "net"

type User struct{
	Name string
	Addr string
	C	chan string
	conn net.Conn
}

//创建一个用户的API
func NewUser(conn net.Conn) *User{
	// 获取客户端连接的IP地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()

	return user
}

// 监听当前user channel的方法, 一旦有消息就直接发给客户端
func (this *User) ListenMessage(){
	for {
		msg := <- this.C
		
		this.conn.Write([]byte(msg + "\n"))
	}
}