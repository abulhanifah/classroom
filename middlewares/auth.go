package middlewares

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/abulhanifah/classroom/helpers"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if IsSkipAuth(c.Path()) {
			return next(c)
		} else {
			isTokenValid, isACLValid := CheckAuth(c)
			if !isTokenValid {
				return echo.NewHTTPError(401, "Token otentikasi tidak valid.")
			}
			if !isACLValid {
				return echo.NewHTTPError(404, "Pengguna tidak memiliki cukup izin untuk mengakses sumber daya.")
			}
			return next(c)
		}
	}
}

func CheckAuth(c echo.Context) (bool, bool) {
	ctx := helpers.SetContext(c)

	isTokenValid := false
	isACLValid := true
	tokenInfo, err := helpers.ValidateToken(ctx)
	if err == nil {
		helpers.SetAuthContext(ctx, tokenInfo)
		isTokenValid = true
		// isACLValid = helpers.CheckAcl(ctx, helpers.GetAclKeyFromRequest(ctx))
	}
	return isTokenValid, isACLValid
}

func IsSkipAuth(path string) bool {
	if strings.HasPrefix(path, "/callback") {
		return true
	}
	if strings.HasPrefix(path, "/auth") {
		return true
	}
	if strings.HasPrefix(path, "/oauth/access_token") {
		return true
	}
	return false
}
