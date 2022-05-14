package server

import (
	"bytes"
	"encoding/binary"
	pb "example/proto"
	"net"

	"google.golang.org/protobuf/proto"
)

const (
	MAXLEN int = 1024
)

//整形转换成字节
func IntToBytes(val int64) ([]byte, error) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuffer, binary.BigEndian, val); err != nil {
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}

//字节转换成整形
func BytesToInt(b []byte) (int64, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var iVal int64
	if err := binary.Read(bytesBuffer, binary.BigEndian, &iVal); err != nil {
		return -1, err
	}
	return iVal, nil
}

// 从指定通道中读取固定长度的字节
func ReadBytes(conn net.Conn, len int64) ([]byte, error) {
	buf := make([]byte, len)
	index := 0
	for {
		n, err := conn.Read(buf[index:])
		if err != nil {
			return buf, err
		}
		index += n
		if int64(n) >= len {
			break
		}
	}
	return buf, nil
}

// 从指定通道解析一条数据
func DecodeOneMsg(conn net.Conn) (*pb.Message, error) {
	// 读取长度
	byteLen, err := ReadBytes(conn, 8)
	if err != nil {
		// fmt.Printf("Read length failed: %v \n", err.Error())
		return nil, err
	}
	len_, _ := BytesToInt(byteLen)
	// fmt.Printf("Decode %v bytes", len_)
	// 读取指定长度的内容
	byteContent, err := ReadBytes(conn, len_)
	if err != nil {
		// fmt.Printf("Read content failed: %v \n", err.Error())
		return nil, err
	}
	// 解析消息
	msg := &pb.Message{}
	err = proto.Unmarshal(byteContent, msg)
	if err != nil {
		// fmt.Printf("Parse msg failed: %v \n", err.Error())
		return nil, err
	}
	// fmt.Printf("\n\n[Decode: len-%vbytes msg-%vbytes content-%vbytes]\n\n", len_, len(byteContent), len(msg.Content))
	return msg, nil
}

// 封装一条信息
func EncodeOneMsg(msgType int, userAdress string, serverAdress string, localAdress string, content []byte) ([]byte, error) {
	msg := pb.Message{
		MsgType:      2,
		UserAdress:   userAdress,
		ServerAdress: serverAdress,
		LocalAdress:  localAdress,
		Content:      content,
	}
	msgByte, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	// 获取消息长度
	msgLen := len(msgByte)
	lenByte, err := IntToBytes(int64(msgLen))
	if err != nil {
		return nil, err
	}
	// 拼接
	writeByte := append(lenByte[:], msgByte[:]...)
	// fmt.Printf("\n\n[Encode: content-%vbytes msg-%vbytes len-%vbytes total-%vbytes]\n\n", len(content), len(msgByte), len(lenByte), len(writeByte))
	return writeByte, nil
}