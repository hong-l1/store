package main

import "log"

func main() {
	server := InitApiServer()
	log.Println("API Server starting on :8081")
	server.Run(":8081")
}
