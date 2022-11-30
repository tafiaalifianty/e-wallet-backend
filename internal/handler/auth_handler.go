package handler

import (
	"net/http"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/dto"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/helper"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

func (h *Handler) Login(c *gin.Context) {
	var input dto.LoginRequestBody
	err := c.ShouldBindJSON(&input)
	if err != nil {
		helper.WriteErrorResponse(
			c,
			http.StatusBadRequest,
			custom_error.InvalidRequestBody{}.Error(),
			nil,
		)
		return
	}

	res, err := h.services.Auth.Login(input.Email, input.Password)

	if err != nil {
		helper.WriteErrorResponse(
			c,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil,
		)
		return
	}

	helper.WriteSuccessResponse(
		c,
		http.StatusOK,
		http.StatusText(http.StatusOK),
		res,
	)
}

func (h *Handler) Register(ctx *gin.Context) {
	var input dto.RegisterRequestBody
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			custom_error.InvalidRequestBody{}.Error(),
			nil,
		)
		return
	}

	user := &entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	token, err := h.services.Auth.Register(user)

	if _, ok := err.(*custom_error.EmailAlreadyUsed); ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			err.Error(),
			nil,
		)

		return
	}

	if err != nil {
		helper.WriteErrorResponse(
			ctx,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			nil,
		)
		return
	}

	helper.WriteSuccessResponse(
		ctx,
		http.StatusCreated,
		http.StatusText(http.StatusCreated),
		token,
	)
}
