package middlewares

import (
	"fmt"

	"github.com/abulhanifah/classroom/constants"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func TransactionHandler(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Ini", db)
			c.Set(constants.CtxDB, db)
			tx := db.Begin()
			c.Set(constants.CtxTx, tx)
			n := next(c)
			if c.Response().Status >= 200 && c.Response().Status < 400 {
				tx.Commit()
			} else {
				tx.Rollback()
			}
			return n
		}
	}
}
