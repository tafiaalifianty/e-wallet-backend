package middlewares

import (
	"net/http"

	"assignment-golang-backend/internal/helper"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	tokenStr, err := helper.ParseAuthorizationHeader(authorizationHeader)
	if err != nil {
		helper.WriteErrorResponse(
			c,
			http.StatusUnauthorized,
			err.Error(),
			nil,
		)
		return
	}

	token, err := helper.ValidateToken(tokenStr)
	if err != nil || !token.Valid {
		helper.WriteErrorResponse(
			c,
			http.StatusUnauthorized,
			err.Error(),
			nil,
		)
		return
	}

	if claims, ok := token.Claims.(*helper.IdTokenClaims); ok {
		c.Set("user", claims.User)
		c.Next()
	} else {
		helper.WriteErrorResponse(
			c,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
			nil,
		)
		return
	}
}
