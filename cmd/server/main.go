package main

import (
	"context"
	"log"

	_ "github.com/juanmaabanto/go-ms-beers/docs"

	"github.com/joho/godotenv"
	"github.com/juanmaabanto/go-ms-beers/common/managers"
	"github.com/juanmaabanto/go-ms-beers/common/middleware"
	"github.com/juanmaabanto/go-ms-beers/internal/ports"
	"github.com/juanmaabanto/go-ms-beers/internal/service"
	"github.com/juanmaabanto/go-ms-beers/internal/validations"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Beer API
// @version v1
// @description Specifying services for falabella.

// @contact.name Juan Manuel Abanto Mera
// @contact.url https://www.linkedin.com/in/juanmanuelabanto/
// @contact.email jmanuelabanto@gmail.com

// @license.name MIT License
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := echo.New()
	ctx := context.Background()

	application := service.NewApplication(ctx)

	Handler(ports.NewHttpServer(application), router)
	router.Logger.Fatal(router.Start(":3000"))
}

type ServerInterface interface {
	AddBeer(c echo.Context) error
	GetBeer(c echo.Context) error
	ListBeer(c echo.Context) error
	GetBoxPrice(c echo.Context) error
}

func Handler(si ServerInterface, router *echo.Echo) {
	if router == nil {
		router = echo.New()
	}

	router.Validator = validations.NewValidationUtil()

	loggerManager := managers.NewLoggerManager("https://logger.mydominio.pe/", "ms-beer")

	api := router.Group("/api/v1")

	api.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		LoggerErrorFunc: loggerManager.Error,
	}))

	//Swagger
	router.GET("/*", echoSwagger.WrapHandler)

	//beer
	api.GET("/beers", si.ListBeer)
	api.GET("/beers/:beerId", si.GetBeer)
	api.POST("/beers", si.AddBeer)
	api.GET("/beers/:beerId/boxprice", si.GetBoxPrice)
}
