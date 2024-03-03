package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	serverIp   string
	serverPort int
	conn       net.Conn
	Name       string
	flag       int //当前客户端模式

}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端
	Client := &Client{
		serverIp:   serverIp,
		serverPort: serverPort,
		flag:       999,
	}

	// 链接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial, error", err)
	}
	Client.conn = conn
	// 返回对象
	return Client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "0.0.0.0", "设计服务器IP地址默认0.0.0.0")
	flag.IntVar(&serverPort, "port", 8888, "设计服务器端口默认8888")
}

func (Client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		Client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法的范围的数字<<<<")
		return false
	}
}
func (client *Client) UpdateName() bool {

	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&client.Name)
	command := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(command))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return false
	}
	return true

}

func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式...")
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式...")
			break
		case 3:
			//更新用户名模式
			client.UpdateName()
			break
		}
	}
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>服务器连接失败")
	}
	// 监听读取服务器发送的数据并输出
	go client.DealResponse()

	fmt.Println(">>>>服务器连接成功")

	// 启动客户端的业务
	client.Run()
}
