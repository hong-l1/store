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
	svc1 service2.ObjectService
	l    logger2.Loggerv1
	svc2 service2.NodeService
}

func NewObjectHandler(svc1 service2.ObjectService, svc2 service2.NodeService, l logger2.Loggerv1) *DataServerHandler {
	return &DataServerHandler{svc1: svc1, l: l, svc2: svc2}
}
func (h *DataServerHandler) RegisterObjectRoute(server *gin.Engine) {
	v1 := server.Group("/objects")
	{
		v1.PUT("/:filename", h.PUT)
		v1.GET("/:filename", h.GET)
	}
	go h.svc2.ListenLocateMsg(0)
	go h.svc2.HeartbeatTicker()
}
func (h *DataServerHandler) PUT(ctx *gin.Context) {
	name := ctx.Param("filename")
	err := h.svc1.Upload(ctx.Request.Context(), name, ctx.Request.Body)
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
	reader, err := h.svc1.Download(ctx.Request.Context(), name)
	if err != nil {
		ctx.JSON(http.StatusOK, internal.Result{Code: 5, Message: "系统错误"})
		h.l.Error("文件获取失败",
			logger2.Error(err),
			logger2.String("文件名:", name),
		)
		return
	}
	n, err := io.Copy(ctx.Writer, reader)
	if err != nil {
		h.l.Error("文件写入响应失败",
			logger2.Error(err),
			logger2.String("文件名", name),
		)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	h.l.Info("写出的字节数", logger2.Int64("bytes", n))
}
