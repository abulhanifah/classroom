package helpers

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"gitlab.com/abulhanifah/classroom/configs"
)

type Oauth2Manager struct {
	Tx *gorm.DB

	// request auth code
	ClientID    string `json:"client_id" form:"client_id" query:"client_id"`
	RedirectUri string `query:"redirect_uri"`
	Scope       string `query:"scope"`
	State       string `query:"state"`

	// request access token
	// ClientID     string
	ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
	GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	Code         string `json:"code" form:"code" query:"code"`
	Username     string `json:"username" form:"username" query:"username"`
	Password     string `json:"password" form:"password" query:"password"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
}

type User struct {
	ID           string
	Email        string
	Phone        string
	Password     string
	LanguageCode string
}

type Oauth2Client struct {
	ID        string `gorm:"type:varchar(32)"`
	Name      string `gorm:"type:varchar(255)"`
	Secret    string `gorm:"type:varchar(255)"`
	UserID    string `gorm:"type:uuid;index:oauth2_client_user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Oauth2RedirectUrl struct {
	ID        string `gorm:"type:uuid"`
	Url       string `gorm:"type:varchar(255)"`
	ClientID  string `gorm:"type:varchar(32);index:oauth2_redirect_url_client_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Oauth2Token struct {
	ID              string    `json:"-" gorm:"type:uuid"`
	Code            string    `json:"-" gorm:"type:varchar(32);index:oauth2_token_code"`
	UserID          string    `json:"-" gorm:"type:uuid;index:oauth2_token_user_id"`
	ClientID        string    `json:"-" gorm:"type:varchar(32);index:oauth2_token_client_id"`
	AccessToken     string    `json:"access_token" gorm:"type:varchar(32);index:oauth2_token_access_token"`
	TokenType       string    `json:"token_type" gorm:"-"`
	ExpireIn        int       `json:"expire_in" gorm:"-"`
	RefreshToken    string    `json:"refresh_token" gorm:"type:varchar(32);index:oauth2_token_refresh_token"`
	AccessExpiredAt time.Time `json:"-"`
	ExpiredAt       time.Time `json:"-"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func ValidateToken(ctx Context) (*Oauth2Token, error) {
	AccessToken := ""
	auth := strings.Split(ctx.Request().Header.Get("Authorization"), " ")
	if len(auth) > 1 && auth[0] == "Bearer" {
		AccessToken = auth[1]
	}
	t := Oauth2Token{}
	err := GetDB(ctx).Model(&Oauth2Token{}).Where("access_token = ?", AccessToken).First(&t).Error
	if IsRecordNotFoundError(err) || t.AccessExpiredAt.Before(time.Now()) {
		return nil, constant.ErrAccessDenied
	}

	return &t, nil
}

// todo
func (o *Oauth2Manager) GenerateAuthToken(ctx Context) (string, error) {
	o.Tx = GetDB(ctx)

	t := Oauth2Token{}
	t.ClientID = o.ClientID
	t.Code = NewToken()
	return t.Code, nil
}

func (o *Oauth2Manager) GenerateAccessToken(ctx Context) (*Oauth2Token, error) {
	o.Tx = GetDB(ctx)
	if o.GrantType == constant.Oauth2AuthorizationCode {
		return o.Oauth2AuthorizationCode()
	} else if o.GrantType == constant.Oauth2PasswordCredentials {
		return o.Oauth2PasswordCredentials()
	} else if o.GrantType == constant.Oauth2ClientCredentials {
		return o.Oauth2ClientCredentials()
	} else if o.GrantType == constant.Oauth2Refreshing {
		return o.Oauth2Refreshing()
	} else if o.GrantType == constant.Oauth2Implicit {
		return o.Oauth2Implicit()
	} else {
		return &Oauth2Token{}, constant.ErrInvalidGrant
	}
}

// todo
func (o *Oauth2Manager) Oauth2AuthorizationCode() (*Oauth2Token, error) {
	return nil, constant.ErrUnsupportedGrantType
}

func (o *Oauth2Manager) Oauth2PasswordCredentials() (*Oauth2Token, error) {
	t, err := o.ClientLogin()
	if err != nil {
		return nil, err
	}
	t, err = o.UserLogin(t)
	if err != nil {
		return nil, err
	}
	return o.SaveToken(t)
}

func (o *Oauth2Manager) Oauth2ClientCredentials() (*Oauth2Token, error) {
	t, err := o.ClientLogin()
	if err != nil {
		return nil, err
	}
	return o.SaveToken(t)
}

func (o *Oauth2Manager) Oauth2Refreshing() (*Oauth2Token, error) {
	t, err := o.ClientLogin()
	if err != nil {
		return nil, err
	}
	t, err = o.CheckRefreshToken(t)
	if err != nil {
		return nil, err
	}
	return o.SaveToken(t)
}

// todo
func (o *Oauth2Manager) Oauth2Implicit() (*Oauth2Token, error) {
	return nil, constant.ErrUnsupportedGrantType
}

func (o *Oauth2Manager) SaveToken(t *Oauth2Token) (*Oauth2Token, error) {
	if t.ID == "" {
		t.ID = NewUUID()
	}
	o.Tx.Model(&Oauth2Token{}).Where(Oauth2Token{ID: t.ID}).Assign(Oauth2Token{
		ClientID:        t.ClientID,
		UserID:          t.UserID,
		TokenType:       o.GrantType,
		AccessToken:     NewToken(),
		RefreshToken:    NewToken(),
		ExpireIn:        int(configs.Get("OAUTH2_ACCESS_EXPIRE_IN").Duration().Seconds()),
		AccessExpiredAt: time.Now().Add(configs.Get("OAUTH2_ACCESS_EXPIRE_IN").Duration()),
		ExpiredAt:       time.Now().Add(configs.Get("OAUTH2_REFRESH_EXPIRE_IN").Duration()),
	}).FirstOrCreate(&t)
	return t, nil
}

func (o *Oauth2Manager) ClientLogin() (*Oauth2Token, error) {
	client := Oauth2Client{}
	err := o.Tx.Model(&Oauth2Client{}).Where("id = ?", o.ClientID).First(&client).Error
	if IsRecordNotFoundError(err) {
		return nil, constant.ErrInvalidClient
	}
	if o.ClientSecret != client.Secret {
		return nil, constant.ErrUnauthorizedClient
	}

	t := Oauth2Token{}
	t.ClientID = client.ID
	t.UserID = client.UserID
	return &t, nil
}

func (o *Oauth2Manager) UserLogin(t *Oauth2Token) (*Oauth2Token, error) {
	key := "email"
	v := Validation{}
	if !v.IsEmail(o.Username) {
		key = "phone"
	}

	user := User{}
	err := o.Tx.Model(&User{}).Where(key+" = ?", o.Username).First(&user).Error
	if IsRecordNotFoundError(err) || VerifyHash(user.Password, o.Password) != nil {
		return nil, constant.ErrAccessDenied
	}

	t.UserID = user.ID
	return t, nil
}

func (o *Oauth2Manager) CheckRefreshToken(t *Oauth2Token) (*Oauth2Token, error) {
	err := o.Tx.Model(&Oauth2Token{}).Where("refresh_token = ?", o.RefreshToken).First(&t).Error
	if o.RefreshToken == "" || IsRecordNotFoundError(err) || t.ExpiredAt.Before(time.Now()) {
		return nil, constant.ErrAccessDenied
	}
	return t, nil
}
