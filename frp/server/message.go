package server

import "net"

type msgHandler struct {
	userAdress   string
	serverAdress string   // server => localhost:9202
	localAdress  string   // local => 192.168.0.202
	clientConn   net.Conn // server <=> client  |  server <=> user
	userConn     net.Conn
}

// func
