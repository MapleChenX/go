package main

import (
	"fmt"
	"net"
)

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
	defer conn.Close()

	user := NewUser(conn)
	this.OnlineMap[user.Name] = user
	this.BroadCast(user, "i fucking coming!")

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
			msg := string(buf[:n-1])
			this.BroadCast(user, msg)
		}
	}()

	select {
	// 阻塞
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
	// 监听
	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// 关闭
	defer conn.Close()

	go this.ListenMessage()

	// 阻塞接受
	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Println("conn.Accept err:", err)
			continue
		}
		// 处理
		go this.Handler(conn)
	}
}
