package apiserver

import (
	"awesomeProject1/store"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler(ctx *gin.Context) {
	m := ctx.Request.Method
	if m != http.MethodGet {
		ctx.JSON(http.StatusNotFound,
			store.Result{Message: "Not Found"})
		return
	}
	filename := ctx.Param("filename")
	info := Locate(filename)
	if len(info) == 0 {
		return
	}
	ctx.JSON(http.StatusOK, store.Result{
		Message: info,
	})
}
func Locate(name string, producer sarama.SyncProducer) string {
	producer.SendMessage(sarama.)
}

func Exist(name string) bool {
	return Locate(name) != ""
}
