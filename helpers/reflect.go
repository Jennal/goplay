package helpers

import "reflect"

func GetValueFromPtr(ptr interface{}) interface{} {
	return reflect.ValueOf(ptr).Elem().Interface()
}
