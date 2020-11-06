package geval

import (
	"fmt"
	"reflect"
)

type makeMapParam struct {
	tKey   reflect.Type
	tValue reflect.Type
}

// buildInMake : make
func buildInMake(param makeMapParam) interface{} {
	vMap := reflect.MakeMap(reflect.MapOf(param.tKey, param.tValue))
	return vMap.Interface()
}

// buildInLen : len
func buildInLen(param interface{}) int {
	vParam := reflect.ValueOf(param)
	switch vParam.Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.String, reflect.Chan:
		return vParam.Len()
	default:
		panic(fmt.Sprintf("Param type(%v) do not support len operate", vParam.Kind()))
	}
}
