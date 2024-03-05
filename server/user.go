package server

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	// 获取客户端连接的IP地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()

	return user
}

// 用户上线
func (this *User) Online() {
	//用户上线, 将用户加入到onLineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 广播当前用户上线消息
	this.server.BroadCast(this, "用户上线 ^_^")
}

// 用户下线
func (this *User) Offline() {
	//用户下线, 将用户加入到onLineMap中
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	// 广播当前用户下线消息
	this.server.BroadCast(this, "用户下线 Q_Q")

}

// 给当前用户客户端发送消息 指定用户
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + " 在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 修改用户名
		NewName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[NewName]
		if ok {
			this.SendMsg("当前用户名已使用")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[NewName] = this
			this.Name = NewName
			this.SendMsg("用户名已更改为:" + this.Name + "\n")
			this.server.mapLock.Unlock()
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 私聊
		// 消息格式 to|用户名|内容

		// 1.获取对方用户名, 消息
		content := strings.Split(msg, "|")
		remoteUser := content[1]
		message := content[2]
		// 2.根据用户名 得到对方User对象
		userObject, ok := this.server.OnlineMap[remoteUser]
		if ok {
			userObject.SendMsg(this.Name + "对您说:" + message + "\n")
		} else {
			this.SendMsg("用户不存在")
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听当前user channel的方法, 一旦有消息就直接发给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
