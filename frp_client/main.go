package main

import "client"

func main() {
	remoteAdress := "42.192.57.121:9200"
	configMap := map[string]string{
		":9202": "192.168.0.202:22",
		":9203": "192.168.0.203:22",
		":9204": "192.168.0.204:22",
	}
	client := &client.Client{}
	client.Run(remoteAdress, configMap)
}
