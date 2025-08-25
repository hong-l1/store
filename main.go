package main

import (
	"awesomeProject1/demo"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	server := gin.Default()
	server.Handle(http.MethodGet, "/objects/:filename", GetObject)
	server.Handle(http.MethodPut, "/objects/:filename", PutObject)
	server.Run(":8080")
}
func GetObject(ctx *gin.Context) {
	name := ctx.Param("filename")
	log.Println(name)
	f, err := os.Open(filepath.Join(os.Getenv("STORAGE_ROOT"), name))
	if err != nil {
		ctx.JSON(http.StatusOK, demo.Result{Code: -1, Message: "系统错误"})
		return
	}
	io.Copy(ctx.Writer, f)
	ctx.JSON(http.StatusOK, demo.Result{Message: "Ok"})
}
func PutObject(ctx *gin.Context) {
	name := ctx.Param("filename")
	log.Println(name)
	f, err := os.Create(filepath.Join(os.Getenv("STORAGE_ROOT"), name))
	if err != nil {
		ctx.JSON(http.StatusOK, demo.Result{Code: -1, Message: "系统错误"})
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, ctx.Request.Body)
	ctx.JSON(http.StatusOK, demo.Result{Message: "上传成功"})
}
