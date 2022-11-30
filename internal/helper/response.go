package helper

import "github.com/gin-gonic/gin"

type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteSuccessResponse(
	ctx *gin.Context,
	code int,
	message string,
	data interface{},
) {
	ctx.JSON(code, JsonResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func WriteErrorResponse(
	ctx *gin.Context,
	code int,
	message string,
	data interface{},
) {
	ctx.AbortWithStatusJSON(code, JsonResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
