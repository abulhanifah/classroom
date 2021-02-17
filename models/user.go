package models

import (
	"time"

	"github.com/abulhanifah/classroom/helpers"
)

type User struct {
	ID              string    `json:"id,omitempty" gorm:"type:char(36)"`
	Name            string    `json:"name,omitempty" gorm:"type:varchar(100)" validate:"required,max=100"`
	Gender          string    `json:"gender,omitempty" gorm:"type:varchar(6)" validate:"required,eq=male|eq=female"`
	Email           string    `json:"email,omitempty" gorm:"type:varchar(100)" validate:"required_without=phone,email,max=100"`
	Password        string    `json:"password,omitempty" validate:"required_with=email,max=100"`
	Phone           string    `json:"phone,omitempty" validate:"required_without=email,number"`
	IsPhoneVerified bool      `json:"is_phone_verified,omitempty"`
	IsEmailVerified bool      `json:"is_email_verified,omitempty"`
	RoleID          int       `json:"role.id,omitempty" gorm:"index:users_role_id"`
	RoleName        string    `json:"role.name,omitempty" query:"role.name,omitempty" gorm:"-"`
	LanguageCode    string    `json:"language_code,omitempty" gorm:"type:char(2)" validate:"max=2"`
	CountryCode     string    `json:"country_code,omitempty" gorm:"type:char(2)" validate:"max=2"`
	Province        string    `json:"province,omitempty" gorm:"type:varchar(50)" validate:"max=50"`
	City            string    `json:"city,omitempty" gorm:"type:varchar(75)" validate:"max=75"`
	Token           string    `json:"token,omitempty" gorm:"-"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

func (o *User) Schema() map[string]interface{} {
	return map[string]interface{}{
		"table": map[string]string{"name": "users", "as": "u"},
		"fields": map[string]map[string]string{
			"id":                {"name": "u.id", "as": "id"},
			"name":              {"name": "u.name", "as": "name"},
			"gender":            {"name": "u.gender", "as": "gender"},
			"email":             {"name": "u.email", "as": "email"},
			"phone":             {"name": "u.phone", "as": "phone"},
			"virtual_account":   {"name": "u.virtual_account", "as": "virtual_account"},
			"is_email_verified": {"name": "u.is_email_verified", "as": "is_email_verified"},
			"is_phone_verified": {"name": "u.is_phone_verified", "as": "is_phone_verified"},
			"country_code":      {"name": "u.country_code", "as": "country_code"},
			"province":          {"name": "u.province", "as": "province"},
			"city":              {"name": "u.city", "as": "city"},
			"language_code":     {"name": "u.language_code", "as": "language_code"},
			"role.id":           {"name": "r.id", "as": "role_id"},
			"role.name":         {"name": "r.name", "as": "role_name"},
			"created_at":        {"name": "u.created_at", "as": "created_at"},
			"updated_at":        {"name": "u.updated_at", "as": "updated_at"},
		},
		"relations": []map[string]string{
			{"name": "user_roles", "as": "r", "on": "r.id = u.role_id", "type": "hasOne"},
		},
	}
}

func (o *User) GetById(ctx helpers.Context, id string, params map[string][]string) map[string]interface{} {
	v := helpers.Validation{}
	key := "id"
	if !v.IsUUID(id) {
		key = "email"
		if !v.IsEmail(id) {
			key = "phone"
		}
	}
	return helpers.GetById(ctx, "users", key, id, params, o.Schema())
}

func (o *User) GetPaginated(ctx helpers.Context, params map[string][]string) map[string]interface{} {
	return helpers.GetPaginated(ctx, params, o.Schema())
}

func (o *User) Create(ctx helpers.Context) map[string]interface{} {
	isValid, msg := helpers.Validate(ctx, o)
	if !isValid {
		return msg
	}
	params, err := o.SetDefaultValue(ctx)
	if err != nil {
		return params
	}
	helpers.GetDB(ctx).Create(o)
	return helpers.GetById(ctx, "users", "id", o.ID, map[string][]string{}, o.Schema())
}

func (o *User) UpdateById(ctx helpers.Context) map[string]interface{} {
	helpers.GetDB(ctx).Model(User{}).Where("id = ?", o.ID).Updates(o)
	return helpers.GetById(ctx, "users", "id", o.ID, map[string][]string{}, o.Schema())
}

func (o *User) DeleteById(ctx helpers.Context) map[string]interface{} {
	id := o.ID
	helpers.GetDB(ctx).Model(User{}).Where("id = ?", o.ID).Delete(&User{})
	return helpers.DeletedMessage("users", "id", id)
}

func (o *User) SetDefaultValue(ctx helpers.Context) (map[string]interface{}, error) {
	var params map[string]interface{}
	var err error
	if o.ID == "" {
		o.ID = helpers.NewUUID()
	}
	if o.RoleID == 0 {
		o.RoleID = 2 // default is buyer
	}
	if o.CountryCode == "" {
		o.CountryCode = "ID"
	}
	params, err = o.IsEmailAvailable(ctx)
	if err != nil {
		return params, err
	}
	params, err = o.IsPhoneAvailable(ctx)
	if err != nil {
		return params, err
	}
	return params, nil
}

func (o *User) BeforeSave() error {
	hashedPassword, err := helpers.Hash(o.Password)
	if err != nil {
		return err
	}
	o.Password = string(hashedPassword)
	return nil
}

func (o *User) GetByKey(ctx helpers.Context, key, value string) (*User, error) {
	user := User{}
	err := helpers.GetDB(ctx).Model(&User{}).Where(key+" = ?", value).First(&user).Error
	if helpers.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &user, err
}

func (o *User) IsEmailAvailable(ctx helpers.Context) (map[string]interface{}, error) {
	i := User{}
	helpers.GetDB(ctx).Model(&User{}).Where("email = ?", o.Email).First(&i)
	if i.ID != "" {
		return helpers.ErrorEmailAlreadyBeenTaken(ctx)
	}
	return map[string]interface{}{}, nil
}

func (o *User) IsPhoneAvailable(ctx helpers.Context) (map[string]interface{}, error) {
	i := User{}
	helpers.GetDB(ctx).Model(&User{}).Where("phone = ?", o.Phone).First(&i)
	if i.ID != "" {
		return helpers.ErrorPhoneAlreadyBeenTaken(ctx)
	}
	return map[string]interface{}{}, nil
}
