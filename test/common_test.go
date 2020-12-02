package test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/modern-go/reflect2"
)

func TestReflect2(t *testing.T) {
	data := make(map[string]interface{})
	subData := make(map[string]interface{})

	subData["content"] = "sub data content"

	data["hello"] = "world"
	data["sub_data"] = subData

	t.Log("Origin data: ", data)

	t2 := reflect2.TypeOf(data)
	mt2 := t2.(reflect2.MapType)

	key := "hello"
	retStr := mt2.GetIndex(&data, &key).(*interface{})
	// t.Log("hello type1: ", reflect.TypeOf(retStr), ", typ2: ", reflect2.TypeOf(retStr), ", value: ", retStr, ", realKind: ", reflect.TypeOf(retStr).Elem())
	t.Logf("hello type1: %v, type2: %T, type1*: %v, type2*: %v, value: %v, value*: %v", reflect.TypeOf(retStr), reflect2.TypeOf(retStr), reflect.TypeOf(*retStr), reflect2.TypeOf(*retStr), retStr, *retStr)

	t2pInter := reflect2.Type2(t2.Type1().Elem())
	ptr := mt2.UnsafeGetIndex(reflect2.PtrOf(&data), reflect2.PtrOf(&key))
	t.Logf("ptr: %v, strptr: %v", ptr, retStr)
	testRet := t2pInter.PackEFace(ptr).(*interface{})
	t.Logf("testRet: %v, type: %v", *testRet, t2pInter)

	// updateStr := (*retStr).(string)
	// updateStr = "Fuck"
	// t.Logf("After retStr update, updateStr: %v, *retStr: %v, data: %v", updateStr, *retStr, data)

	key = "sub_data"
	rretMap := mt2.GetIndex(&data, &key)
	retMap := rretMap.(*interface{})
	t.Log("sub_data kind: ", reflect.TypeOf(retMap), ", value: ", *retMap)

	t.Logf("retMap: %v, prtMap: %v", retMap, reflect2.PtrOf(retMap))

	retInter := *retMap
	efaceRetInter := (*eface)(unsafe.Pointer(&retInter))
	t.Logf("efaceRetInter: %v, data: %v", efaceRetInter.rtype, efaceRetInter.data)

	actualData := (*retMap).(map[string]interface{})

	efaceRetMap := (*eface)(unsafe.Pointer(&rretMap))
	t.Logf("rtype: %v, data: %v", efaceRetMap.rtype, efaceRetMap.data)

	t.Logf("retMap: %v, actualPtr: %v, ptr retMap: %v", retMap, reflect2.PtrOf(&actualData), reflect2.PtrOf(retMap))
	key = "content"
	subPtr := mt2.UnsafeGetIndex(reflect2.PtrOf(&actualData), reflect2.PtrOf(&key))
	subRet := t2pInter.PackEFace(subPtr).(*interface{})
	t.Log("sub_data content kind: ", reflect.TypeOf(subRet), ", value: ", *subRet)
}

func TestPtr(t *testing.T) {
	str := "Hello world"
	pstr := &str
	var face interface{}
	face = str
	pface := &face

	efaceFace := (*eface)(unsafe.Pointer(pface))
	t.Logf("pstr: %v, eface.rtype: %v, eface.data: %v, pface: %v", pstr, efaceFace.rtype, efaceFace.data, pface)

	pstr = (*string)(efaceFace.data)
	t.Log(*pstr)
}

func TestEface(t *testing.T) {
	str := "Hello world"
	pstr := &str
	var ptmp *string
	var face *eface

	face = unpackEface(str)
	t.Logf("Before change str, str: %v, *eface.data: %v, pstr: %v, eface.data: %v", str, *(*string)(face.data), pstr, face.data)
	ptmp = (*string)(face.data)
	*ptmp = "Fuck world"
	t.Logf("After change str, str: %v, *eface.data: %v, pstr: %v, eface.data: %v", str, *(*string)(face.data), pstr, face.data)

	face = unpackEface(&str)
	t.Logf("Before change &str, str: %v, *eface.data: %v, pstr: %v, eface.data: %v", str, *(*string)(face.data), pstr, face.data)
	ptmp = (*string)(face.data)
	*ptmp = "Fuck world"
	t.Logf("After change &str, str: %v, *eface.data: %v, pstr: %v, eface.data: %v", str, *(*string)(face.data), pstr, face.data)

	face = unpackEface(&pstr)
	t.Logf("Before change &pstr, str: %v, pstr: %v, *eface.data: %v", str, pstr, *(**string)(face.data))
	ptmp = *(**string)(face.data)
	*ptmp = "Enn hello world"
	t.Logf("After change &pstr, str: %v, pstr: %v, *eface.data: %v", str, pstr, *(**string)(face.data))

	var i, j interface{}
	i = &str
	face = unpackEface(i)
	t.Logf("interface i, type: %v, str: %v, *eface.data: %v, eface.data: %v", reflect.TypeOf(i), str, *(*string)(face.data), face.data)

	j = &i
	face = unpackEface(j)
	t.Logf("unpack 1 &i, type: %v, str: %v, *eface.data: %v, eface.data: %v", reflect.TypeOf(j), str, *(*interface{})(face.data), face.data)
	face = unpackEface(*(*interface{})(face.data))
	t.Logf("unpack 2 &i, str: %v, *eface.data: %v, eface.data: %v, eface.rtype: %v", str, *(*string)(face.data), face.data, (face.rtype))
	ptmp = (*string)(face.data)
	*ptmp = "Fuck again"
	t.Logf("Final str: %v", str)

	it := getPtr(unsafe.Pointer(&str), reflect.TypeOf(str))
	t.Logf("ptr %v: value: %v", it, *(it.(*string)))
	str = "kkk"
	t.Logf("ptr %v: value: %v", it, *(it.(*string)))

	i = str
	it = getFacePtr(i)
	t.Logf("ptr %v: value: %v", it, *(it.(*string)))
	str = "bbb"
	t.Logf("ptr %v: value: %v", it, *(it.(*string)))
}

func TestReflectSlice(t *testing.T) {
	s := []int{1, 2, 3}

	a := [3]int{1, 2, 3}

	t.Logf("sType: %v, aType: %v", reflect.TypeOf(s), reflect.TypeOf(a))
}

func TestEfaceSlice(t *testing.T) {
	s := []int{1, 2, 3}
	ps := &s
	var is interface{}
	is = s
	pis := &is

	sEface := unpackEface(s)
	pEface := unpackEface(ps)
	iEface := unpackEface(is)
	piEface := unpackEface(pis)

	t.Logf("type, s: %T, ps: %T, is: %T, pis: %T", s, ps, is, pis)
	t.Logf("ptr, s: %p, p: %p, is: %p, pis: %p", s, ps, is, pis)
	t.Logf("data, s: %v, p: %v, is: %v, pis: %v", sEface.data, pEface.data, iEface.data, piEface.data)

	index := 0
	type1 := reflect2.TypeOf(s).(reflect2.SliceType)
	v1 := type1.GetIndex(&s, index)
	t.Logf("v1: %T", v1)
}

func TestEfaceMap(t *testing.T) {
	m1 := make(map[string]interface{})
	m1["s"] = "Hello"

	m2 := make(map[string]string)
	m2["s"] = "Hello"

	s := "Hello"
	m3 := make(map[string]*string)
	m3["s"] = &s

	key := "s"

	type1 := reflect2.TypeOf(m1).(reflect2.MapType)
	v1 := type1.GetIndex(&m1, &key)

	type2 := reflect2.TypeOf(m2).(reflect2.MapType)
	v2 := type2.GetIndex(&m2, &key)

	type3 := reflect2.TypeOf(m3).(reflect2.MapType)
	v3 := type3.GetIndex(&m3, &key)

	v11 := ptrElem(v1)
	v22 := ptrElem(v2)
	v33 := ptrElem(v3)

	t.Logf("type, m1: %T, v1: %T, m2: %T, v2: %T, m3: %T, v3: %T, v11: %T, v22: %T, v33: %T", m1["s"], v1, m2["s"], v2, m3["s"], v3, v11, v22, v33)

	type2Map := reflect2.TypeOf(m1).(reflect2.MapType)
	newMap := type2Map.MakeMap(0)
	t.Logf("new map: %T, *newMap: %T", newMap, ptrElem(newMap))
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

func getFacePtr(obj interface{}) interface{} {
	var i interface{}
	typeFace := unpackEface(reflect.PtrTo(reflect.TypeOf(obj)))
	objFace := unpackEface(obj)

	e := (*eface)(unsafe.Pointer(&i))
	e.rtype = typeFace.data
	e.data = objFace.data
	return i
}

func getPtr(ptr unsafe.Pointer, tType reflect.Type) interface{} {
	var i interface{}
	typeFace := unpackEface(reflect.PtrTo(tType))

	e := (*eface)(unsafe.Pointer(&i))
	e.rtype = typeFace.data
	e.data = ptr
	return i
}

func ptrElem(obj interface{}) interface{} {
	tObj := reflect.TypeOf(obj)
	kObj := tObj.Kind()
	if reflect.Ptr == kObj {
		return reflect2.Type2(tObj.Elem()).Indirect(obj)
	}
	return obj
}

func faceToPrt(obj interface{}) interface{} {
	var i interface{}
	cType := reflect.TypeOf(obj)
	pType := reflect.PtrTo(cType)
	typeFace := unpackEface(pType)
	objFace := unpackEface(obj)

	e := (*eface)(unsafe.Pointer(&i))
	e.rtype = typeFace.data
	e.data = objFace.data
	return i
}

type eface struct {
	rtype unsafe.Pointer
	data  unsafe.Pointer
}
