package routers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gogateway/frontend"
)

func RegisterService(engine *gin.Engine) {
	engine.GET("/ws", func(c *gin.Context) {
		err := frontend.Instance().HandleRequest(c.Writer, c.Request)
		if err != nil {
			log.Error(err)
			return
		}
	})
}
