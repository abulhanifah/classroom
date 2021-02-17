package routes

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"gitlab.com/abulhanifah/classroom/configs"
)

func Init(db *gorm.DB) *echo.Echo {
	e := echo.New()
	env := configs.Get("APP_ENV").String()
	if env == "production" || env == "development" {
		e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize: 1 << 10, // 1 KB
		}))
	}

	return e
}
