package handler

import (
	"net/http"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/dto"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/helper"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	user := api.Group("/users")
	{
		user.GET("/info", h.GetUserInfo)
	}
}

func (h *Handler) GetUserInfo(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			nil,
		)
		return
	}

	res, err := h.services.User.FindByID(int(user.(*entity.TokenizedUser).ID))

	if _, ok := err.(*custom_error.NoDataFound); ok {
		helper.WriteErrorResponse(
			ctx,
			http.StatusNotFound,
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
		http.StatusOK,
		http.StatusText(http.StatusOK),
		dto.FormatUser(res),
	)
}
