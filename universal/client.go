package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(ip string, port int, name string) *Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 1234))
	if err != nil {
		fmt.Println("cao!")
	}

	return &Client{
		ServerIP:   ip,
		ServerPort: port,
		Name:       name,
		conn:       conn,
	}
}

func (this *Client) Run() {
	defer this.conn.Close()

	// 接收服务器回发的数据
	go this.recv()

	// 客户端发送数据
	this.send()
}

func (this *Client) recv() {
	buf := make([]byte, 4096)
	for {
		n, err := this.conn.Read(buf)
		if n == 0 {
			fmt.Println("服务器端关闭，客户端也退出")
			return
		}
		if err != nil {
			fmt.Println("conn.Read err:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}

func (this *Client) send() {
	var msg string
	for {
		fmt.Scanln(&msg)
		_, err := this.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("conn.Write err:", err)
		}
	}
}

func main() {
	client := NewClient("127.0.0.1", 1234, "zhangsan")
	client.Run()
	select {}
}
