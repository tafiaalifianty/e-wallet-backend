package helper

import (
	"os"
	"strconv"
	"strings"
	"time"

	"assignment-golang-backend/internal/config"
	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"

	"github.com/golang-jwt/jwt/v4"
)

type IdTokenClaims struct {
	jwt.RegisteredClaims
	User *entity.TokenizedUser `json:"user"`
}

func GenerateJWT(user *entity.User) (string, error) {
	expMinute, _ := strconv.Atoi(config.GetEnv("TOKEN_EXP_MINUTE"))
	var idExp int64 = int64(expMinute * 60)
	unixTime := time.Now().Unix()
	tokenExp := unixTime + idExp

	tokenizedUser := &entity.TokenizedUser{
		ID:           int(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		WalletNumber: user.WalletNumber,
	}

	claims := &IdTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.GetEnv("TOKEN_ISSUER"),
			ExpiresAt: &jwt.NumericDate{Time: time.Unix(tokenExp, 0)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		User: tokenizedUser,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(
		[]byte(config.GetEnv("TOKEN_SECRET")),
	)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		token,
		&IdTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, &custom_error.InvalidToken{}
			}

			return []byte(os.Getenv("TOKEN_SECRET")), nil
		},
	)
}

func ParseAuthorizationHeader(authHeader string) (string, error) {
	authHeaderSplit := strings.Split(authHeader, "Bearer ")
	if len(authHeaderSplit) != 2 {
		return "", &custom_error.AuthHeaderNotAvailable{}
	}

	return strings.TrimSpace(authHeaderSplit[1]), nil
}
