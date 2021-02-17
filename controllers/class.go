package controllers

import (
	"github.com/abulhanifah/classroom/helpers"
	"github.com/abulhanifah/classroom/models"
	"github.com/labstack/echo/v4"
)

func GetClassHandle(c echo.Context) error {
	res := map[string]interface{}{}
	o := new(models.Classroom)
	ctx := helpers.SetContext(c)
	if c.Param("id") != "" {
		var id int = helpers.Convert(c.Param("id")).Int()
		res = o.GetClassById(ctx, id)
	} else {
		classes := []models.Classroom{}
		helpers.GetDB(ctx).Table("classrooms").Scan(&classes)
		res["count"] = len(classes)
		data := []map[string]interface{}{}
		if len(classes) > 0 {
			for _, d := range classes {
				data = append(data, map[string]interface{}{
					"class_id":   d.ID,
					"class_name": d.Name,
					"rows":       d.Rows,
					"columns":    d.Columns,
				})
			}
		}
		res["data"] = data
	}
	return c.JSON(200, res)
}

func CreateClassHandle(c echo.Context) error {
	o := new(models.Classroom)
	if err := c.Bind(o); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	res := o.CreateClass(helpers.SetContext(c))
	return helpers.Response(c, 201, res)
}
