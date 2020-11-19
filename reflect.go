package geval

import (
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type eface struct {
	rtype unsafe.Pointer
	data  unsafe.Pointer
}

func unpackEface(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}

func packEface(rtype unsafe.Pointer, ptr unsafe.Pointer) interface{} {
	var i interface{}
	e := (*eface)(unsafe.Pointer(&i))
	e.rtype = rtype
	e.data = ptr
	return i
}

func faceToPrt(obj interface{}) interface{} {
	cType := reflect.TypeOf(obj)
	pType := reflect.PtrTo(cType)
	typeFace := unpackEface(pType)
	e := (*eface)(unsafe.Pointer(&obj))
	e.rtype = typeFace.data
	return obj
}

func ptrElem(obj interface{}) interface{} {
	tObj := reflect.TypeOf(obj)
	kObj := tObj.Kind()
	if reflect.Ptr == kObj {
		return reflect2.Type2(tObj.Elem()).Indirect(obj)
	}
	return obj
}

func faceToReal(obj interface{}) interface{} {
	tObj := reflect.TypeOf(obj)
	if reflect.Interface != tObj.Kind() {
		return obj
	}

	typeFace := unpackEface(tObj)
	e := (*eface)(unsafe.Pointer(&obj))
	e.rtype = typeFace.data
	return obj
}
