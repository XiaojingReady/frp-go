package client

import (
	pb "example/proto"
	"fmt"

	"google.golang.org/protobuf/proto"
	// "net"
)

type Client struct {
	serverHandler *connHandler
	localHandlers []*connHandler
	adressMap     map[string]*connHandler
}

func (client *Client) handlerServerMsg() {
	var buf [20480]byte
	for {
		n, err := client.serverHandler.serverConn.Read(buf[0:])
		if err != nil {
			fmt.Printf("Can not read msg: %v \n", err.Error())
			return
		}

		// 转发到local端口
		message := &pb.Message{}
		err = proto.Unmarshal(buf[0:n], message)
		if err != nil {
			panic(err)
		}
		// contentMap := Decode(buf[0:n])

		// userAdress := message.UserAdress
		handler, ok := client.adressMap[message.UserAdress]
		if !ok {
			fmt.Printf("Get a new ssh connect: %s \n", message.UserAdress)
			handler = &connHandler{
				userAdress:   message.UserAdress,
				serverAdress: message.ServerAdress,
				localAdress:  message.LocalAdress,
				serverConn:   client.serverHandler.serverConn,
				localConn:    Connect(message.LocalAdress),
			}
			client.adressMap[message.UserAdress] = handler
			go client.listenLocal(handler)
		}
		handler.localConn.Write(message.Content)
	}
}

func (client *Client) listenLocal(handler *connHandler) {
	// 将local的信息转发到server

	var buf [20480]byte
	for {
		n, err := handler.localConn.Read(buf[0:])
		if err != nil {
			fmt.Printf("Can not read msg: %v \n", err.Error())
			return
		}
		message := pb.Message{
			MsgType:      2,
			UserAdress:   handler.userAdress,
			ServerAdress: handler.serverAdress,
			LocalAdress:  handler.localAdress,
			Content:      buf[0:n],
		}
		messageByte, err := proto.Marshal(&message)
		fmt.Printf("userAdress: %s replay to server: %v \n", handler.userAdress, string(buf[0:n]))
		handler.serverConn.Write(messageByte)
	}
}

func (client *Client) Run(serverAdress string, configMap map[string]string) {
	fmt.Println("Client start")
	// 开启服务器的数据通道
	client.serverHandler = &connHandler{
		serverAdress: serverAdress, // server => localhost:9202
		serverConn:   Connect(serverAdress),
	}
	client.adressMap = map[string]*connHandler{}
	// 向服务器发送配置表
	client.serverHandler.sendConfig(configMap)
	client.handlerServerMsg()
}
