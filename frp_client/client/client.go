package client

import (
	"encoding/json"
	"fmt"
	"net"
	// "net"
)

type Client struct {
	serverConn net.Conn
	adressMap  map[string]*connHandler
}

func (c *Client) sendConfig(configMap map[string]string) {
	// map转为[]byte
	byteConfig, err := json.Marshal(configMap)
	if err != nil {
		fmt.Printf("sendConfig: Can not encode configMap")
		return
	}
	// 封装消息
	writeByte, err := EncodeOneMsg(1, "", "", "", byteConfig)
	if err != nil {
		fmt.Printf("sendConfig: Encode failed, err:  = %v\n", err)
		return
	}
	// 发送消息
	code, err := c.serverConn.Write(writeByte)
	if err != nil {
		fmt.Printf("sendConfig: Write failed, error code: %v, err:  = %v\n", code, err)
		return
	}
}

// 监听指定用户通道的数据
func (c *Client) listenLocal(handler *connHandler) {
	var buf [MAXLEN]byte
	for {
		// 从通道中读取n字节的数据
		n, err := handler.localConn.Read(buf[0:])
		if err != nil {
			fmt.Printf("listenLocal: Read error: %v\n", err.Error())
			return
		}
		// 封装消息
		writeByte, err := EncodeOneMsg(2, handler.userAdress, handler.serverAdress, handler.localAdress, buf[0:n])
		if err != nil {
			fmt.Printf("listenLocal: Encode failed, err:  = %v\n", err)
			return
		}

		// 发送消息
		code, err := handler.serverConn.Write(writeByte)
		if err != nil {
			fmt.Printf("listenLocal: Write failed, error code: %v, err:  = %v\n", code, err)
			return
		}
	}
}

func (client *Client) listenServer() {
	defer client.serverConn.Close()
	for {
		// 从client通道解析一条数据
		msg, err := DecodeOneMsg(client.serverConn)
		if err != nil {
			fmt.Printf("Can not decode one msg, error: %v\n", err.Error())
			return
		}
		// log.Printf("\nFrom server: %v\n\n", msg)
		// 转发到local端口
		handler, ok := client.adressMap[msg.UserAdress]
		if !ok {
			fmt.Printf("Create new ssh to %s \n", msg.LocalAdress)
			localConn, err := net.Dial("tcp", msg.LocalAdress)
			// conn, err := Connect(msg.LocalAdress)
			if err != nil {
				fmt.Printf("listenServer: Connect failed, error: %v\n", err.Error())
				return
			}
			handler = &connHandler{
				userAdress:   msg.UserAdress,
				serverAdress: msg.ServerAdress,
				localAdress:  msg.LocalAdress,
				serverConn:   client.serverConn,
				localConn:    localConn,
				server:       client,
			}
			client.adressMap[msg.UserAdress] = handler

			go handler.listenLocal()
		}
		_, err = handler.localConn.Write(msg.Content)
		if err != nil {
			fmt.Printf("listenServer: Write failed, error: %v\n", err.Error())
			return
		}
	}
}

func (client *Client) Run(serverAdress string, configMap map[string]string) {
	fmt.Println("Client start...")
	// 连接服务器
	conn, err := Connect(serverAdress)
	if err != nil {
		fmt.Printf("Run: Can not connect %s: %v \n", serverAdress, err.Error())
		return
	}
	client.serverConn = conn
	client.adressMap = map[string]*connHandler{}
	// 向服务器发送配置表
	client.sendConfig(configMap)
	client.listenServer()
}
