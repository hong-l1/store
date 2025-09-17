package api

import (
	"awesomeProject1/internal"
	logger2 "awesomeProject1/internal/pkg/zapx"
	service2 "awesomeProject1/internal/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type DataServerHandler struct {
	svc service2.ObjectService
	l   logger2.Loggerv1
}

func NewObjectHandler(svc service2.ObjectService, l logger2.Loggerv1) *DataServerHandler {
	return &DataServerHandler{svc: svc, l: l}
}
func (h *DataServerHandler) RegisterObjectRoute(server *gin.Engine) {
	v1 := server.Group("/objects")
	{
		v1.PUT("/:filename", h.PUT)
		v1.GET("/:filename", h.GET)
		//v1.DELETE("/objects/:name", objHandler.Delete)
	}
}
func (h *DataServerHandler) PUT(ctx *gin.Context) {
	name := ctx.Param("filename")
	err := h.svc.Upload(ctx.Request.Context(), name, ctx.Request.Body)
	if err != nil {
		h.l.Error("文件上传失败",
			logger2.Error(err),
			logger2.String("object", name),
		)
		ctx.JSON(http.StatusOK, internal.Result{Code: 5, Message: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, internal.Result{Message: "上传成功"})
}
func (h *DataServerHandler) GET(ctx *gin.Context) {
	name := ctx.Param("filename")
	reader, err := h.svc.Download(ctx.Request.Context(), name)
	if err != nil {
		ctx.JSON(http.StatusOK, internal.Result{Code: 5, Message: "系统错误"})
		h.l.Error("文件获取失败",
			logger2.Error(err),
			logger2.String("文件名:", name),
		)
		return
	}
	if _, err := io.Copy(ctx.Writer, reader); err != nil {
		h.l.Error("文件写入响应失败",
			logger2.Error(err),
			logger2.String("文件名", name),
		)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

//// 删除对象
//func (h *ObjectHandler) Delete(c *gin.Context) {
//	name := c.Param("name")
//
//	if err := h.svc.Delete(c.Request.Context(), name); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "delete success"})
//}
