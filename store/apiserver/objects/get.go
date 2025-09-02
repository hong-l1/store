package objects

import (
	"awesomeProject1/store"
	logger "awesomeProject1/store/pkg/zapx"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func GetObject(ctx *gin.Context) {
	filename := ctx.Param("filename")
	stream, err := getStream(filename)
	if err != nil {
		logger.Error(err)
		return
	}
	_, err = io.Copy(ctx.Writer, stream)
	if err != nil {
		logger.Error(err)
	}
	ctx.JSON(http.StatusOK, store.Result{
		Message: "ok"})
}
