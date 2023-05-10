package handler

import (
	_ "auth/docs"
	"auth/internal/config"
	v1 "auth/internal/handler/v1"
	"auth/internal/service"
	"auth/pkg/auth"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{services: services, tokenManager: tokenManager}
}

func (h *Handler) InitRoutes(cfg config.Config) *gin.Engine {
	gin.SetMode(cfg.GIN.Mode)
	router := gin.New()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerv1 := v1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		handlerv1.Init(api)
	}
}
