package main

func main() {
	server := InitApiServer()
	server.Run(":8081")
}
