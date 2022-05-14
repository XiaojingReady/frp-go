package server

import (
	"fmt"
	"net"
)

type connHandler struct {
	userAdress   string   // 用户的地址
	serverAdress string   // 服务器的端口
	localAdress  string   // 本地机器的地址
	clientConn   net.Conn // 与客户端的数据通道
	userConn     net.Conn // 与用户的数据通道
	server       *Server
}

// 监听指定用户通道的数据
func (c *connHandler) listenUser() {
	defer func() {
		fmt.Printf("Close conn: User[%v] to Local[%v]\n", c.userAdress, c.localAdress)
		delete(c.server.adressMap, c.userAdress)
		c.userConn.Close()
	}()
	defer c.userConn.Close()
	var buf [MAXLEN]byte
	for {
		// 从通道中读取n字节的数据
		n, err := c.userConn.Read(buf[0:])
		if err != nil {
			fmt.Printf("listenUser: Can not read user conn, conn closed: %v \n", err.Error())
			return
		}
		// 封装消息
		writeByte, err := EncodeOneMsg(2, c.userAdress, c.serverAdress, c.localAdress, buf[0:n])
		if err != nil {
			fmt.Printf("listenUser: Encode failed, err:  = %v\n", err)
			return
		}
		// 发送消息
		// fmt.Printf("\nFrom user: %v\n\n", string(buf[0:n]))
		// fmt.Printf("Send Msg: From user[%v] to local[%v]: content: %v\n", c.userAdress, c.localAdress, string(buf[0:n]))
		code, err := c.clientConn.Write(writeByte)
		if err != nil {
			fmt.Printf("listenUser: Write failed, error code: %v, err:  = %v\n", code, err)
			return
		}
	}
}
