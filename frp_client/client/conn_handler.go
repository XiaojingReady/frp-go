package client

import (
	"fmt"
	"net"
)

type connHandler struct {
	userAdress   string   // 用户的地址
	serverAdress string   // 服务器的端口
	localAdress  string   // 本地机器的地址
	serverConn   net.Conn // 与客户端的数据通道
	localConn    net.Conn // 与用户的数据通道
	server       *Client
}

// 监听指定本地机器通道的数据
func (c *connHandler) listenLocal() {
	defer func() {
		fmt.Printf("Close conn: User[%v] to Local[%v]\n", c.userAdress, c.localAdress)
		delete(c.server.adressMap, c.userAdress)
		c.localConn.Close()
	}()
	var buf [MAXLEN]byte
	for {
		// 从通道中读取n字节的数据
		n, err := c.localConn.Read(buf[0:])
		if err != nil {
			fmt.Printf("listenLocal: Can not read local conn, conn closed: %v \n", err.Error())
			return
		}
		// 封装消息
		writeByte, err := EncodeOneMsg(2, c.userAdress, c.serverAdress, c.localAdress, buf[0:n])
		if err != nil {
			fmt.Printf("listenLocal: Encode failed, err:  = %v\n", err)
			return
		}
		// 发送消息
		// fmt.Printf("\nFrom local: %v\n\n", string(buf[0:n]))
		// tmp := DecodeOneMsg(conn net.Conn)
		// fmt.Printf("Recive msg from local: From local[%v] to user[%v]: content: %v\n", c.localAdress, c.userAdress, string(buf[0:n]))
		code, err := c.serverConn.Write(writeByte)
		if err != nil {
			fmt.Printf("listenLocal: Write failed, error code: %v, err:  = %v\n", code, err)
			return
		}
	}
}
