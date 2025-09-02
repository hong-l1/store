package demo

import (
	"awesomeProject1/demo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.Handle(http.MethodGet, "/objects/:filename", demo.GetObject)
	server.Handle(http.MethodPut, "/objects/:filename", demo.PutObject)
	server.Run(":8080")
}
