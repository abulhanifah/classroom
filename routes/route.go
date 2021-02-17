package routes

import (
	"github.com/abulhanifah/classroom/configs"
	"github.com/abulhanifah/classroom/controllers"
	"github.com/abulhanifah/classroom/middlewares"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(db *gorm.DB) *echo.Echo {
	e := echo.New()
	env := configs.GetConfig("APP_ENV").String()
	if env == "production" || env == "development" {
		e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize: 1 << 10, // 1 KB
		}))
	}
	e.Use(middlewares.TransactionHandler(db))
	e.Use(middlewares.Auth)

	e.GET("/api/login", controllers.LoginHandle)
	e.POST("/api/create_class", controllers.CreateClassHandle)
	e.POST("/api/check_in", controllers.CheckInClassHandle)
	e.POST("/api/check_out", controllers.CheckOutClassHandle)
	e.GET("/api/get_class_list", controllers.GetClassHandle)
	e.GET("/api/get_class_by_id/:id", controllers.GetClassHandle)

	return e
}
