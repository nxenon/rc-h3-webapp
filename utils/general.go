package utils

import "reflect"

func StructToMap(data interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	val := reflect.ValueOf(data)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		out[field.Name] = val.Field(i).Interface()
	}
	return out
}
