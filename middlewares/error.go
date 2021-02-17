package middlewares

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"

	"github.com/abulhanifah/classroom/configs"
	"github.com/abulhanifah/classroom/helpers"
)

func ErrorHandler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    500,
			Message: "Internal Server Error",
		}
	}

	code := he.Code
	message := helpers.GetCtx(helpers.SetContext(c)).ErrorMessage()
	if len(message) == 0 {
		message = map[string]interface{}{
			"error": map[string]interface{}{
				"code":    code,
				"message": he.Message,
			},
		}
	}
	if code >= 500 && (configs.GetConfig("APP_ENV").String() == "production" || configs.GetConfig("APP_ENV").String() == "development") {
		trace := map[string]string{}
		for i := 0; i <= 15; i++ {
			fun, file, no, _ := runtime.Caller(i)
			if file != "" {
				trace[fmt.Sprintf("#%d", i)] = fmt.Sprintf("%s:%d on %s", file, no, runtime.FuncForPC(fun).Name())
			}
		}

		log := map[string]interface{}{}
		log["request"] = c.Request().Method + " " + c.Path()
		log["response"] = message
		log["trace"] = trace

		dataJson, _ := json.MarshalIndent(log, "", "  ")
		fmt.Println(dataJson)
		// b := map[string]string{"text": "```" + string(dataJson) + "```"}
		// helpers.RequestApi("POST", config.Get("SLACK_ERROR_URL").String(), b, map[string]string{})
	}
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			c.NoContent(he.Code)
		} else {
			c.JSON(code, message)
		}
	}
}
