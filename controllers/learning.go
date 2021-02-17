package controllers

import (
	"github.com/abulhanifah/classroom/helpers"
	"github.com/abulhanifah/classroom/models"
	"github.com/labstack/echo/v4"
)

func CheckInClassHandle(c echo.Context) error {
	o := echo.Map{}
	if err := c.Bind(&o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	var classId int
	var userId string
	if o["user_id"] != "" {
		userId = o["user_id"].(string)
	}
	if o["class_id"] != "" {
		classId = int(o["class_id"].(float64))
	}
	res := models.BookSeat(helpers.SetContext(c), classId, userId, "in")
	return helpers.Response(c, 201, res)
}

func CheckOutClassHandle(c echo.Context) error {
	o := echo.Map{}
	if err := c.Bind(&o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	var classId int
	var userId string
	if o["user_id"] != "" {
		userId = o["user_id"].(string)
	}
	if o["class_id"] != "" {
		classId = int(o["class_id"].(float64))
	}
	res := models.BookSeat(helpers.SetContext(c), classId, userId, "out")
	return helpers.Response(c, 201, res)
}
