package handler

import (
	"encoding/json"
	"strings"
	"testing"

	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

var MockTokenizedUser *entity.TokenizedUser = &entity.TokenizedUser{
	ID:           1,
	Name:         "name",
	Email:        "email",
	WalletNumber: 1,
}

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	return r
}

func MakeRequestBody(dto interface{}) *strings.Reader {
	payload, _ := json.Marshal(dto)
	return strings.NewReader(string(payload))
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	return
}

func MiddlewareMockUser(ctx *gin.Context) {
	ctx.Set("user", MockTokenizedUser)
	ctx.Next()
}

func TestNew(t *testing.T) {
	service := &usecase.Services{}
	New(service)
}

func TestHandler_InitAPI(t *testing.T) {
	router := SetUpRouter()
	service := &usecase.Services{}
	handler := New(service)
	handler.InitAPI(router)
}
