package api

import (
	"awesomeProject1/internal"
	"awesomeProject1/internal/domain"
	logger2 "awesomeProject1/internal/pkg/zapx"
	service2 "awesomeProject1/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ApiServerHandler struct {
	nodesvc service2.NodeService
	l       logger2.Loggerv1
	client  *http.Client
	metasvc service2.MetadataService
}

func NewApiServerHandler(nodesvc service2.NodeService, l logger2.Loggerv1, metasvc service2.MetadataService) *ApiServerHandler {
	return &ApiServerHandler{
		nodesvc: nodesvc,
		l:       l,
		client:  &http.Client{Timeout: 3 * time.Second},
		metasvc: metasvc,
	}
}
func (h *ApiServerHandler) RegisterApiServerRoute(server *gin.Engine) {
	v1 := server.Group("/api")
	{
		v1.PUT("/:filename", h.PUT)
		v1.GET("/:filename", h.GET)
		v1.GET("/:filename/versions", h.Versions)
		v1.DELETE("/:filename", h.Delete)
	}
	go h.nodesvc.ListenHeartbeats(0)
	go h.nodesvc.CleanDeadNodes()
}
func (h *ApiServerHandler) PUT(ctx *gin.Context) {
	hash := h.GetHash(ctx)
	if hash == "" {
		ctx.JSON(http.StatusBadRequest, internal.Result{Message: "无hash值"})
		return
	}
	name := ctx.Param("filename")
	addr := h.nodesvc.ChooseNode(ctx.Request.Context())
	if addr == "" {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "无可用节点"})
		return
	}
	err := h.Dohttp(ctx, addr, name, http.MethodPut)
	_ = h.metasvc.PutNewMetadata(ctx, name, 0, ctx.Request.ContentLength)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "无可用节点"})
	}
	ctx.JSON(http.StatusOK, internal.Result{Message: "ok"})
	return
}
func (h *ApiServerHandler) GET(ctx *gin.Context) {
	name := ctx.Param("filename")
	versionStr := ctx.Query("version")
	var (
		err     error
		version int
	)
	if versionStr == "" {
		version = 0
	} else {
		version, err = strconv.Atoi(versionStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, internal.Result{Message: "version参数无效"})
			return
		}
	}
	meta, err := h.metasvc.GetMetadata(ctx, name, version)
	if err != nil {
		h.l.Warn("元数据查询失败", logger2.String("filename", name), logger2.Error(err))
		ctx.JSON(http.StatusBadRequest, internal.Result{Message: "无法找到元数据"})
		return
	}
	addr := h.nodesvc.Locate(ctx.Request.Context(), meta)
	h.l.Info("定位对象", logger2.String("filename", name), logger2.String("addr", addr))
	if addr == "" {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "无可用节点"})
		return
	}
	if err = h.Dohttp(ctx, addr, name, http.MethodGet); err != nil {
		h.l.Error("文件代理失败", logger2.String("addr", addr), logger2.Error(err))
		ctx.JSON(http.StatusBadGateway, internal.Result{Message: "节点请求失败"})
		return
	}
}
func (h *ApiServerHandler) GetVersionDetail(ctx *gin.Context) {
	name := ctx.Param("filename")
	versionstr := ctx.Param("version")
	version, err := strconv.Atoi(versionstr)
	if err != nil {
		h.l.Error("版本解析错误", logger2.String("filename", name),
			logger2.String("version", versionstr), logger2.Error(err))
		ctx.JSON(http.StatusBadRequest, internal.Result{Message: "版本号错误"})
		return
	}
	mate := h.metasvc.GetVersionDetail(ctx, name, version)
	if mate.Addr == "" {
		ctx.JSON(http.StatusNotFound, internal.Result{Message: "版本不存在"})
		return
	}
	if err := h.Dohttp(ctx, mate.Addr, mate.Hash, http.MethodGet); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "节点不可用"})
		return
	}
	return
}
func (h *ApiServerHandler) Delete(ctx *gin.Context) {
	name := ctx.Param("filename")
	meta, err := h.metasvc.SearchLatestVersion(ctx, name)
	if err != nil {
		h.l.Error("查询最新版本失败",
			logger2.String("filename", name),
			logger2.Error(err),
		)
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "元数据查询失败"})
		return
	}
	err = h.metasvc.PutNewMetadata(ctx, meta.Filename, meta.Version+1, 0)
	if err != nil {
		h.l.Error("写入删除标记失败",
			logger2.String("filename", name),
			logger2.Error(err),
		)
		ctx.JSON(http.StatusServiceUnavailable, internal.Result{Message: "删除标记写入失败"})
	}
	ctx.JSON(http.StatusOK, internal.Result{Message: "删除成功"})
	return
}
func (h *ApiServerHandler) Versions(ctx *gin.Context) {
	name := ctx.Param("filename")
	metas, err := h.metasvc.Versions(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, internal.Result{Message: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, metas)
	return
}
func (h *ApiServerHandler) Dohttp(ctx *gin.Context, addr, name, method string) error {
	url := fmt.Sprintf("http://%s/objects/%s", addr, name)
	req, _ := http.NewRequest(method, url, nil)
	req.Header = ctx.Request.Header.Clone()
	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	ctx.Status(resp.StatusCode)
	_, err = io.Copy(ctx.Writer, resp.Body)
	if err != nil {
		h.l.Error("下载转发失败",
			logger2.Error(err),
			logger2.String("object", name),
			logger2.String("node_addr", addr),
		)
		return err
	}
	return nil
}
func (h *ApiServerHandler) GetHash(ctx *gin.Context) string {
	digest := ctx.GetHeader("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}
