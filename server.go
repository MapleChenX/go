package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

// 功能注释
// 结构体的行为有两种：主要行为和普遍行为
// 主要行为：大功能行为 --- 放在结构体所在文件
// 普遍行为：服务大功能的行为，像私有方法一样 --- 一个专门的普遍行为文件

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("connect succ:", conn.RemoteAddr().String())

	user := NewUser(conn)
	this.OnlineMap[user.Name] = user
	this.BroadCast(user, "i fucking coming!")

	isAlive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "i fucking going!")
				return
			}
			if err != nil {
				fmt.Println("conn.Read err:", err)
				return
			}
			isAlive <- true
			msg := string(buf[:n-1])
			this.BroadCast(user, msg)
		}
	}()

	for {
		select {
		case <-isAlive:
		case <-time.After(time.Second * 5):
			this.T(user)
			return
		}
	}
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		for _, user := range this.OnlineMap {
			user.C <- msg
		}
	}
}

func (this *Server) Start() {
	// 启动服务器
	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// 关闭
	defer conn.Close()

	// 服务器运行中
	go this.ListenMessage()

	// 接受连接
	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Println("conn.Accept err:", err)
			continue
		}
		defer conn.Close()

		// 监听每个连接的消息
		go this.Handler(conn)
	}
}

func (this *Server) T(user *User) {
	user.sendMessage("i am going away now...")
	time.Sleep(time.Second * 1)
	_, ok := this.OnlineMap[user.Name]
	if ok {
		delete(this.OnlineMap, user.Name)
		close(user.C)
		err := user.conn.Close()
		if err != nil {
			log.Println("user.conn.Close err:", err)
			return
		}
	}
}
