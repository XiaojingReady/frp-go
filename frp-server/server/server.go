package server

import (
	"encoding/json"
	"fmt"
	"net"
)

type Server struct {
	adressMap      map[string]*connHandler
	portMap        map[string]string
}

func NewServer() *Server {
	return &Server{
		adressMap:      make(map[string]*connHandler),
		portMap:        map[string]string{},
	}
}

func (server *Server) handleUser(port string, adress string, clientConn net.Conn) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Listen user failed! Error message: %v \n", err.Error())
		return
	}
	fmt.Printf("Listening, port-%v: adress-%v \n", port, adress)
	// 循环等待连接
	for {
		userConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Unable to accept a user connect, error: %s\n", err.Error())
			continue
		}

		userAdress := userConn.RemoteAddr().String()
		userHandler := &connHandler{
			userAdress:   userAdress,
			serverAdress: port,
			localAdress:  adress,
			clientConn:   clientConn,
			userConn:     userConn,
			server:       server,
		}
		server.adressMap[userAdress] = userHandler
		go userHandler.listenUser()
	}
}

func (server *Server) handleClient(clientConn net.Conn) {
	defer clientConn.Close()
	for {
		// 从client通道解析一条数据
		msg, err := DecodeOneMsg(clientConn)
		if err != nil {
			fmt.Printf("Can not decode one msg, error: %v\n", err.Error())
			return
		}
		// fmt.Printf("\nFrom client: %v\n\n", msg)
		// 处理消息
		if msg.MsgType == 1 {
			// 接收配置表, 开放新的端口
			
			configMap := map[string]string{}
			err := json.Unmarshal(msg.Content, &configMap)
			if err != nil {
				fmt.Printf("Parse config map failed: %v \n", err.Error())
				continue
			}
			// 开放端口等待user连接
			for port, adress := range configMap {
				_, ok := server.portMap[port]
				if !ok {
					server.portMap[port] = adress
					go server.handleUser(port, adress, clientConn)
				}
			}
		} else {
			// 普通转发信息
			// fmt.Printf("From client-common: from local[%v] to user[%v], content: %v\n", msg.LocalAdress, msg.UserAdress, string(msg.Content))
			handler, ok := server.adressMap[msg.UserAdress]
			if !ok {
				fmt.Printf("can not find handler: %s \n", msg.UserAdress)
				continue
			}
			// fmt.Printf("To user[%v], content: %v\n", msg.UserAdress, string(msg.Content))
			handler.userConn.Write(msg.Content)
		}
	}
}

func (server *Server) Run(serverAdress string) {
	// 开启本地服务端口, 用于接收客户端和用户端的连接
	serverServer, err := net.Listen("tcp", serverAdress)
	if err != nil {
		fmt.Printf("Listen client failed! Error message: %v \n", err.Error())
		return
	}
	fmt.Printf("Start listening, server Port: %s\n", serverAdress)

	// 循环等待连接
	for {
		clientConn, err := serverServer.Accept()
		if err != nil {
			fmt.Printf("Unable to accept a client connect, error: %s\n", err.Error())
			continue
		}
		go server.handleClient(clientConn)
	}
}
