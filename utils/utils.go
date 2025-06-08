package utils

import (
	"encoding/json"
	"os"
	"reflect"
)

func EntityHasDeleted(entity any) bool {
	hasDeleted := false
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	e := val.Type()

	for i := 0; i < e.NumField(); i++ {
		fieldName := e.Field(i).Name
		if fieldName == "Deleted" {
			hasDeleted = true
		}
	}
	return hasDeleted
}

func Map2Struct(ob map[string]any, entity any) error {
	jsonString, err := json.Marshal(ob)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, entity)
	return err
}

func IsEmpty(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func:
		return v.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	case reflect.String:
		return v.Len() == 0
	default:
		return false
	}
}

func Getenv(name, fallback string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return fallback
	}
	return value
}
