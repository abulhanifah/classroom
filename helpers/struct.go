package helpers

import (
	"reflect"

	"github.com/fatih/structs"
)

func GetStructName(s interface{}) string {
	if t := reflect.TypeOf(s); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func StructToMap(s interface{}) map[string]interface{} {
	st := structs.New(s)
	st.TagName = "json"
	m := st.Map()
	return m
}
