package models

import (
	"time"

	"github.com/abulhanifah/classroom/helpers"
)

type UserRole struct {
	ID        int       `json:"id,omitempty" form:"id,omitempty" query:"id,omitempty"`
	Name      string    `json:"name" gorm:"type:varchar(50)" validate:"required|string"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (o *UserRole) Schema() map[string]interface{} {
	return map[string]interface{}{
		"table": map[string]string{"name": "user_roles"},
		"fields": map[string]map[string]string{
			"id":         {"name": "id", "as": "id"},
			"name":       {"name": "name", "as": "name"},
			"created_at": {"name": "created_at", "as": "created_at"},
			"updated_at": {"name": "updated_at", "as": "updated_at"},
		},
	}
}

func (o *UserRole) GetById(ctx helpers.Context, id string, params map[string][]string) map[string]interface{} {
	return helpers.GetById(ctx, "user_roles", "id", id, params, o.Schema())
}

func (o *UserRole) GetPaginated(ctx helpers.Context, params map[string][]string) map[string]interface{} {
	return helpers.GetPaginated(ctx, params, o.Schema())
}

func (o *UserRole) Create(ctx helpers.Context) map[string]interface{} {
	isValid, msg := helpers.Validate(ctx, o)
	if !isValid {
		return msg
	}
	params, err := o.SetDefaultValue(ctx)
	if err != nil {
		return params
	}
	helpers.GetDB(ctx).Create(o)
	return helpers.GetById(ctx, "user_roles", "id", helpers.Convert(o.ID).String(), map[string][]string{}, o.Schema())
}

func (o *UserRole) UpdateById(ctx helpers.Context) map[string]interface{} {
	helpers.GetDB(ctx).Model(UserRole{}).Where("id = ?", o.ID).Updates(o)
	return helpers.GetById(ctx, "user_roles", "id", helpers.Convert(o.ID).String(), map[string][]string{}, o.Schema())
}

func (o *UserRole) DeleteById(ctx helpers.Context) map[string]interface{} {
	id := helpers.Convert(o.ID).String()
	helpers.GetDB(ctx).Model(UserRole{}).Where("id = ?", o.ID).Delete(&UserRole{})
	return helpers.DeletedMessage("user_roles", "id", id)
}

func (o *UserRole) SetDefaultValue(ctx helpers.Context) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	return params, nil
}
