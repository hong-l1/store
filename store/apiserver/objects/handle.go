package objects

import (
	"awesomeProject1/store/dataserver/handle"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(server *gin.Engine) {
	group := server.Group("/objects")
	group.GET("/:filename", handle.GetObjects)
	group.PUT("/:filename", handle.PutObjects)
}
