package main

import "example/server"

func main() {
	clientAdress := ":9200"
	server := server.NewServer()
	server.Run(clientAdress)
}
