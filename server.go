package main

import (
	"awesomeProject/configs"
	"awesomeProject/routes"
	"awesomeProject/util/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Starting server...")

	configs.LoadEnv()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.Sugar.Errorf("RECOVER: from error %s", err)
			return err
		},
	}))

	configs.ConnectDB()
	configs.ConnectRedis()

	routes.UserRoute(e)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
