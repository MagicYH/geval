package test

import (
	"reflect"
	"testing"

	"github.com/modern-go/reflect2"
)

func BenchmarkStdSetMap(b *testing.B) {
	var m map[string]string
	tMap := reflect.TypeOf(m)
	vM := reflect.MakeMapWithSize(tMap, 1)
	for i := 0; i < b.N; i++ {
		vM.SetMapIndex(reflect.ValueOf("hello"), reflect.ValueOf("world"))
	}
}

func BenchmarkSafeSetMap(b *testing.B) {
	var m map[string]string
	t2Map := reflect2.ConfigSafe.Type2(reflect.TypeOf(m)).(reflect2.MapType)

	k := "hello"
	v := "world"
	vM := t2Map.MakeMap(1)
	for i := 0; i < b.N; i++ {
		t2Map.SetIndex(vM, &k, &v)
	}
}

func BenchmarkUnsafeSetMap(b *testing.B) {
	var m map[string]string
	t2Map := reflect2.ConfigUnsafe.Type2(reflect.TypeOf(m)).(reflect2.MapType)

	k := "hello"
	v := "world"
	vM := t2Map.MakeMap(1)
	for i := 0; i < b.N; i++ {
		t2Map.SetIndex(vM, &k, &v)
	}
}

func BenchmarkStdGetMap(b *testing.B) {
	data := make(map[string]string)
	data["hello"] = "world"
	for i := 0; i < b.N; i++ {
		vData := reflect.ValueOf(data)
		vValue := vData.MapIndex(reflect.ValueOf("hello"))
		vValue.String()
	}
}

func BenchmarkUnsafeGetMap(b *testing.B) {
	data := make(map[string]string)
	data["hello"] = "world"
	t2Map := reflect2.ConfigUnsafe.Type2(reflect.TypeOf(data)).(reflect2.MapType)

	k := "hello"
	for i := 0; i < b.N; i++ {
		t2Map.GetIndex(&data, &k)
	}
}

func BenchmarkSafeGetMap(b *testing.B) {
	data := make(map[string]string)
	data["hello"] = "world"
	t2Map := reflect2.ConfigUnsafe.Type2(reflect.TypeOf(data)).(reflect2.MapType)

	k := "hello"
	for i := 0; i < b.N; i++ {
		t2Map.GetIndex(&data, &k)
	}
}
