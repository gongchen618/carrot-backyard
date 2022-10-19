package main

import (
	"carrot-backyard/config"
	"carrot-backyard/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	e := echo.New()
	router.InitRouter(e.Group(config.C.App.Prefix))

	e.Use(middleware.CORS())

	log.Fatal(e.Start(config.C.App.Addr))
}
