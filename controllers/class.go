package controllers

import (
	"github.com/abulhanifah/classroom/helpers"
	"github.com/abulhanifah/classroom/models"
	"github.com/labstack/echo/v4"
)

func GetClassListHandle(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{"token": "ss"})
}

func CreateClassHandle(c echo.Context) error {
	o := new(models.Classroom)
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	res := o.CreateClass(helpers.SetContext(c))
	return helpers.Response(c, 201, res)
}
