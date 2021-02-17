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

func IsInArray(val string, array []string) bool {
	exist, _ := InArray(val, array)
	return exist
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
