package client

import (
	"encoding/json"
	pb "example/proto"
	"fmt"
	"net"

	"google.golang.org/protobuf/proto"
)

type connHandler struct {
	userAdress   string
	serverAdress string   // server => localhost:9202
	localAdress  string   // local => 192.168.0.202
	serverConn   net.Conn // server <=> client  |  server <=> user
	localConn    net.Conn
}

// type connectHandler struct {
// 	remoteAdress string   // server => localhost:9202
// 	localAdress  string   // local => 192.168.0.202
// 	conn         net.Conn // local conn => 192.168.0.202
// }

func (c *connHandler) sendConfig(configMap map[string]string) {
	byteConfig, err := json.Marshal(configMap)
	if err != nil {
		fmt.Printf("Can not encode configMap")
		return
	}

	message := pb.Message{
		MsgType:      1,
		UserAdress:   "",
		ServerAdress: "",
		LocalAdress:  "",
		Content:      byteConfig,
	}
	messageByte, err := proto.Marshal(&message)

	c.serverConn.Write(messageByte)
}
