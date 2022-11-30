package main

import (
	"log"

	"assignment-golang-backend/database"
	"assignment-golang-backend/internal/config"
	"assignment-golang-backend/internal/handler"
	"assignment-golang-backend/internal/repository"
	"assignment-golang-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	database.Connect()

	r := gin.Default()

	rp := repository.New(database.Get())
	s := usecase.New(rp)
	h := handler.New(s)

	h.InitAPI(r)
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
