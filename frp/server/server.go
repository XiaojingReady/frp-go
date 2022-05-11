package server

import (
	"encoding/json"
	pb "example/proto"
	"fmt"
	"net"

	"google.golang.org/protobuf/proto"
)

type Server struct {
	clientHandlers []*msgHandler
	userHandlers   []*msgHandler
	adressMap      map[string]*msgHandler
	portMap        map[string]string
}

func NewServer() *Server {
	return &Server{
		clientHandlers: make([]*msgHandler, 0),
		userHandlers:   make([]*msgHandler, 0),
		adressMap:      make(map[string]*msgHandler),
		portMap:        map[string]string{},
	}
}

func (server *Server) listenPort(port string, adress string, clientConn net.Conn) {
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
		userHandler := &msgHandler{
			userAdress:   userAdress,
			serverAdress: port,
			localAdress:  adress,
			clientConn:   clientConn,
			userConn:     userConn,
		}
		server.adressMap[userAdress] = userHandler
		server.userHandlers = append(server.userHandlers, userHandler)
		go server.handleUserMsg(userHandler)
	}
}

func (server *Server) handleUserMsg(userHandler *msgHandler) {
	// 将user的信息转发到client
	var buf [1024]byte
	for {
		n, err := userHandler.userConn.Read(buf[0:])
		if err != nil {
			fmt.Printf("Can not read msg: %v \n", err.Error())
			return
		}
		message := pb.Message{
			MsgType:      2,
			UserAdress:   userHandler.userAdress,
			ServerAdress: userHandler.serverAdress,
			LocalAdress:  userHandler.localAdress,
			Content:      buf[0:n],
		}
		messageByte, err := proto.Marshal(&message)
		userHandler.clientConn.Write(messageByte)
	}
}

func (server *Server) handleClientMsg(msg *msgHandler) {
	var buf [20480]byte
	for {
		n, err := msg.clientConn.Read(buf[0:])
		if err != nil {
			fmt.Printf("Can not read msg: %v \n", err.Error())
			msg.clientConn.Close()
			return
		}

		message := &pb.Message{}
		err = proto.Unmarshal(buf[0:n], message)
		if err != nil {
			panic(err)
		}

		fmt.Printf("message: %v\n", message)
		if message.MsgType == 1 {
			fmt.Println("get config message: ")
			// 接收配置表, 开放新的端口
			configMap := map[string]string{}
			err := json.Unmarshal(message.Content, &configMap)
			if err != nil {
				fmt.Printf("Can not message.Content: %v \n", err.Error())
				continue
			}
			for port, adress := range configMap {
				_, ok := server.portMap[port]
				if !ok {
					server.portMap[port] = adress
					go server.listenPort(port, adress, msg.clientConn)
				}
			}
		} else {
			// 普通转发信息
			userAdress := message.UserAdress
			fmt.Printf("get common message to %s \n", userAdress)
			handler, ok := server.adressMap[userAdress]
			if !ok {
				fmt.Printf("can not find handler: %s \n", userAdress)
				continue
			}
			fmt.Printf("msgMap-content = %v\n", string(message.Content))
			handler.userConn.Write(message.Content)
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
		clientHandler := &msgHandler{
			serverAdress: serverAdress,
			clientConn:   clientConn,
		}
		server.clientHandlers = append(server.clientHandlers, clientHandler)
		go server.handleClientMsg(clientHandler)
	}

}
