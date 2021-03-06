package routes

import (
	"awesomeProject/controllers"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/user", controllers.CreateUser)
	e.GET("/user/:userId", controllers.GetAUser)
	e.GET("/user/all", controllers.GetAllUsers)
}
