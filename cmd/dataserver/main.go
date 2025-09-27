package main

import "log"

func main() {
	server := InitDateServer()
	log.Println("Data Server starting on :8080")
	server.Run(":8080")
}
