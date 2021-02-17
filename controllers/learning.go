package controllers

import (
	"fmt"

	"github.com/abulhanifah/classroom/helpers"
	"github.com/abulhanifah/classroom/models"
	"github.com/labstack/echo/v4"
)

func CheckInClassHandle(c echo.Context) error {
	o := new(models.Classroom)
	if err := c.Bind(o); err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(400, err.Error())
	}
	fmt.Println(o)
	classId := 1
	userId := ""
	res := models.BookSeat(helpers.SetContext(c), classId, userId, "in")
	return helpers.Response(c, 201, res)
}

func CheckOutClassHandle(c echo.Context) error {
	o := echo.Map{}
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	classId := 1
	userId := ""
	res := models.BookSeat(helpers.SetContext(c), classId, userId, "out")
	return helpers.Response(c, 201, res)
}
