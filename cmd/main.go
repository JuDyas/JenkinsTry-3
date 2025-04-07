package main

import (
	handler "github.com/JuDyas/JenkinsTry-3/internal/handlers"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/wordcount", handler.CountWords)
	log.Fatal(e.Start(":8080"))
}
