package objects

import (
	logger "awesomeProject1/store/pkg/zapx"
	"github.com/gin-gonic/gin"
)

func put(ctx *gin.Context) {
	filename := ctx.Param("filename")
	c, err := storeObject(ctx.Request.Body, filename)
	if err != nil {
		logger.Error(err)
		return
	}
	ctx.JSON(c, err)
}
