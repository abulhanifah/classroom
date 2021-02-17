package controllers

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/abulhanifah/classroom/constants"
	"github.com/abulhanifah/classroom/helpers"
	"github.com/labstack/echo/v4"
)

func LoginHandle(c echo.Context) error {
	o := new(helpers.Oauth2Manager)
	o.Tx = helpers.GetDB(helpers.SetContext(c))
	o.RefreshToken = c.QueryParam("refresh_token")

	var tokenDecoded []byte
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	token := strings.Split(auth, " ")
	if len(token) > 1 {
		tokenDecoded, _ = base64.StdEncoding.DecodeString(token[1])
	} else {
		tokenDecoded, _ = base64.StdEncoding.DecodeString(c.QueryParam("token"))
	}
	user := strings.Split(string(tokenDecoded), ":")

	if len(user) > 1 {
		if user[0] == "refresh_token" {
			o.RefreshToken = user[1]
		} else {
			o.Username = user[0]
			o.Password = user[1]
		}
	} else {
		return echo.NewHTTPError(401, "Token otentikasi tidak valid.")
	}

	var err error
	t := new(helpers.Oauth2Token)
	if o.RefreshToken != "" {
		t, err = o.CheckRefreshToken(t)
	} else {
		t, err = o.UserLogin(t)
	}
	if err != nil {
		return echo.NewHTTPError(constants.StatusCodes[err], constants.Descriptions[err])
	}
	t, _ = o.SaveToken(t)
	return c.JSON(200, map[string]interface{}{"token": EncodeToken(t.UserID, t.AccessToken, t.RefreshToken)})
}
func EncodeToken(id, access_token, refresh_token string) string {
	u, _ := json.Marshal(map[string]interface{}{
		"id":            id,
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
	return base64.StdEncoding.EncodeToString([]byte(u))
}
