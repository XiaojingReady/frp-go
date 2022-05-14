package main

import (
	"example/client"
)


func main() {
	remoteAdress := "42.192.57.121:9200"
	configMap := map[string]string{
		":9202": "192.168.0.202:22",
		":9201": "192.168.0.201:22",
	}
	client := &client.Client{}
	client.Run(remoteAdress, configMap)

}
