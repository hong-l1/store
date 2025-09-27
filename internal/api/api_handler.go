package api

import (
	"awesomeProject1/internal"
	logger2 "awesomeProject1/internal/pkg/zapx"
	service2 "awesomeProject1/internal/service"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiServerHandler struct {
	svc    service2.NodeService
	l      logger2.Loggerv1
	client *http.Client
}

func NewApiServerHandler(svc service2.NodeService, l logger2.Loggerv1) *ApiServerHandler {
	return &ApiServerHandler{
		svc:    svc,
		l:      l,
		client: &http.Client{Timeout: 3 * time.Second},
	}
}
func (h *ApiServerHandler) RegisterApiServerRoute(server *gin.Engine) {
	v1 := server.Group("/api")
	{
		v1.PUT("/:filename", h.PUT)
		v1.GET("/:filename", h.GET)
	}
	// 启动心跳监听和清理死节点
	go h.svc.ListenHeartbeats(0)
	go h.svc.CleanDeadNodes()
}
func (h *ApiServerHandler) PUT(ctx *gin.Context) {
	name := ctx.Param("filename")
	addr := h.svc.ChooseNode(ctx.Request.Context())
	if addr == "" {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "无可用节点"})
		return
	}
	url := fmt.Sprintf("http://%s/objects/%s", addr, name)
	req, _ := http.NewRequest("PUT", url, ctx.Request.Body)
	req.Header = ctx.Request.Header.Clone()
	resp, err := h.client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "上传失败"})
		return
	}
	defer resp.Body.Close()
	ctx.JSON(http.StatusOK, internal.Result{Message: "上传成功"})
}
func (h *ApiServerHandler) GET(ctx *gin.Context) {
	name := ctx.Param("filename")
	addr := h.svc.Locate(ctx.Request.Context(), name)
	if addr == "" {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "无可用节点"})
		return
	}
	url := fmt.Sprintf("http://%s/objects/%s", addr, name)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = ctx.Request.Header.Clone()
	resp, err := h.client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "下载失败"})
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(ctx.Writer, resp.Body)
	if err != nil {
		h.l.Error("下载转发失败",
			logger2.Error(err),
			logger2.String("object", name),
			logger2.String("node_addr", addr),
		)
	}
}
