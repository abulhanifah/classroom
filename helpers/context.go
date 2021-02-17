package helpers

import (
	"github.com/abulhanifah/classroom/constants"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type Context echo.Context
type Ctx struct {
	C echo.Context
}

func SetContext(c echo.Context) Context {
	return Context(c)
}

func SetAuthContext(ctx Context, ti *Oauth2Token) {
	ctx.Set(constants.CtxAccessToken, ti.AccessToken)
	ctx.Set(constants.CtxClientID, ti.ClientID)
	ctx.Set(constants.CtxUserID, ti.UserID)

	lang := ctx.Request().Header.Get("Content-Language")
	if lang == "" {
		lang = GetUserLanguage(ctx, ti.UserID)
	}
	ctx.Set(constants.CtxLang, lang)
}

func GetUserLanguage(ctx Context, UserID string) string {
	user := User{}
	GetDB(ctx).Model(&User{}).Where(User{ID: UserID}).First(&user)
	return user.LanguageCode
}

func NewCtx(ctx echo.Context) *Ctx {
	return &Ctx{C: SetContext(ctx)}
}

func GetCtx(ctx Context) *Ctx {
	return &Ctx{C: ctx}
}

func (ctx *Ctx) DB() *gorm.DB {
	return ctx.C.Get(constants.CtxDB).(*gorm.DB)
}

func (ctx *Ctx) Tx() *gorm.DB {
	return ctx.C.Get(constants.CtxTx).(*gorm.DB)
}

func (ctx *Ctx) UserID() string {
	UserID := ""
	CtxUserID := ctx.C.Get(constants.CtxUserID)
	if CtxUserID != nil {
		UserID = CtxUserID.(string)
	}
	return UserID
}

func (ctx *Ctx) UserLang() string {
	lang := "id"
	ctxLang := ctx.C.Get(constants.CtxLang)
	if ctxLang != nil {
		lang = ctxLang.(string)
	}
	return lang
}

func (ctx *Ctx) ClientID() string {
	ClientID := ""
	CtxClientID := ctx.C.Get(constants.CtxClientID)
	if CtxClientID != nil {
		ClientID = CtxClientID.(string)
	}
	return ClientID
}

func (ctx *Ctx) AccessToken() string {
	AccessToken := ""
	CtxAccessToken := ctx.C.Get(constants.CtxAccessToken)
	if CtxAccessToken != nil {
		AccessToken = CtxAccessToken.(string)
	}
	return AccessToken
}

func (ctx *Ctx) ErrorMessage() map[string]interface{} {
	ErrorMessage := map[string]interface{}{}
	CtxErrorMessage := ctx.C.Get(constants.CtxErrorMessage)
	if CtxErrorMessage != nil {
		ErrorMessage = CtxErrorMessage.(map[string]interface{})
	}
	return ErrorMessage
}

func (ctx *Ctx) SetErrorMessage(ErrorMessage map[string]interface{}) {
	ctx.C.Set(constants.CtxErrorMessage, ErrorMessage)
}
