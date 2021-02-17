package models

import (
	"github.com/abulhanifah/classroom/helpers"
)

type Setting struct {
	Key       string `json:"key,omitempty" gorm:"primary_key;type:varchar(100)"`
	Value     string `json:"value,omitempty"`
	CreatedAt string `json:"created_at,omitempty" gorm:"type:timestamp"`
	UpdatedAt string `json:"updated_at,omitempty" gorm:"type:timestamp"`
	DeletedAt string `json:"-"  gorm:"type:timestamp"`
}

func (o *Setting) Get(ctx helpers.Context) map[string]interface{} {
	fields := []Setting{}
	helpers.GetDB(ctx).Find(&fields)
	ret := map[string]interface{}{}
	for _, f := range fields {
		ret[f.Key] = f.Value
	}
	return ret
}

func (o *Setting) Update(ctx helpers.Context, data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		o := Setting{Key: k, Value: v.(string)}
		helpers.GetDB(ctx).Where(Setting{Key: k}).Assign(o).FirstOrCreate(&o)
	}
	return o.Get(ctx)
}
