package client

import (
	"fmt"
	"net"
)

func Connect(remoteAdress string) net.Conn {
	fmt.Printf("Connect to ip: %s \n", remoteAdress)
	remoteConn, err := net.Dial("tcp", remoteAdress)
	if err != nil {
		fmt.Printf("Failed: can not connect to %s \n", remoteAdress)
		return nil
	}
	return remoteConn
}

// func Encode(msgMap map[string][]byte) []byte {
// 	byteMsg, err := json.Marshal(msgMap)
// 	if err != nil {
// 		fmt.Printf("Can not encode msgMap\n")
// 		return []byte{}
// 	}
// 	return byteMsg
// }

// func Decode(byteMsg []byte) map[string][]byte {
// 	new_msg := map[string][]byte{}
// 	err := json.Unmarshal(byteMsg, &new_msg)
// 	if err != nil {
// 		fmt.Printf("Can not decode byteMsg")
// 		return map[string][]byte{}
// 	}
// 	return new_msg
// }
