package router

import (
	"auction-web/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

// NewGinRouter create a new gin router for each micro service
func NewGinRouter(isAuth bool) (router *gin.Engine) {
	router = gin.Default()

	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware)
	if !isAuth {
		router.Use(middlewares.VerifyToken)
	}

	return router
}
