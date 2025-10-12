package main

func main() {
	server := InitDateServer()
	server.Run(":8080")
}
