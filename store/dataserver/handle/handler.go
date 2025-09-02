package handle

import (
	"awesomeProject1/store"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetObjects(c *gin.Context) {
	name := c.Param("filename")
	f, err := os.Open(filepath.Join(os.Getenv("STORAGE_ROOT"), name))
	if err != nil {
		c.JSON(http.StatusOK, store.Result{Code: -1, Message: "系统错误"})
		return
	}
	io.Copy(c.Writer, f)
	c.JSON(http.StatusOK, store.Result{Message: "Ok"})
}
func PutObjects(c *gin.Context) {
	name := c.Param("filename")
	f, err := os.Create(filepath.Join(os.Getenv("STORAGE_ROOT"), name))
	if err != nil {
		c.JSON(http.StatusOK, store.Result{Code: -1, Message: "系统错误"})
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, c.Request.Body)
	c.JSON(http.StatusOK, store.Result{Message: "上传成功"})
}
