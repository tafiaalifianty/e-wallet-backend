package handler

import (
	"net/http"

	"assignment-golang-backend/internal/helper"
	middlewares "assignment-golang-backend/internal/middleware"
	"assignment-golang-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *usecase.Services
}

func New(s *usecase.Services) *Handler {
	return &Handler{
		services: s,
	}
}

func (h *Handler) InitAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		h.initAuthRoutes(api)

		protected := api.Group("/")
		protected.Use(middlewares.AuthorizeJWT)

		h.initUserRoutes(protected)
		h.initTransactionRoutes(protected)
	}

	router.Static("/docs", "dist")

	router.NoRoute(func(ctx *gin.Context) {
		helper.WriteErrorResponse(
			ctx,
			http.StatusNotFound,
			"Page Not Found",
			nil,
		)
	})
}
