package geval

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var typeInterface reflect.Type
var typeMapStrInterface reflect.Type

// DataContext : DataContext is used to store temp data or bind data
type DataContext struct {
	data map[string]interface{}
}

// FunContext : FuncContext is used to store function that should inject into eval engine
type FunContext struct {
	data map[string]interface{}
}

// NewFunCtx : Get a new instance of FunContext, buildin function `make` and `len` is injected
func NewFunCtx() *FunContext {
	ctx := &FunContext{data: make(map[string]interface{})}

	// Bind build in function
	ctx.data["make"] = buildInMake
	ctx.data["len"] = buildInLen

	return ctx
}

// Bind : Inject self define function into eval engine
func (ctx *FunContext) Bind(name string, fun interface{}) error {
	if _, ok := ctx.data[name]; ok {
		return fmt.Errorf("Func '%s' have bind before", name)
	}
	ctx.data[name] = fun
	return nil
}

// NewDataCtx : Get a new instance of DataContext
func NewDataCtx() *DataContext {
	ctx := &DataContext{data: make(map[string]interface{})}
	return ctx
}

// Bind : Inject variable into datacontext, data must be ptr that it's value can be update in eval engine
func (ctx *DataContext) Bind(name string, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return errors.New("Must set ptr")
	}
	if _, ok := ctx.data[name]; ok {
		return fmt.Errorf("Variable '%s' have bind before", name)
	}
	ctx.data[name] = data
	return nil
}

// Get : Get one data from datacontext
func (ctx *DataContext) Get(name string) (value interface{}, err error) {
	//in dataContext
	fields := strings.Split(name, ".")
	data, ok := ctx.data[fields[0]]
	if !ok {
		return nil, errors.New("Fail to get variable")
	}

	if len(fields) > 1 {
		value, _ = getAttribute(data, fields[1:])
	} else {
		value = data
	}
	return
}

// Set set data, Now just support map[string]interface{} type
func (ctx *DataContext) Set(name string, value interface{}) error {
	fields := strings.Split(name, ".")
	data, ok := ctx.data[fields[0]]
	if !ok {
		_, err := setMapValue(reflect.ValueOf(ctx.data), fields[0], value)
		return err
	}

	if len(fields) > 1 {
		return setAttribute(data, fields[1:], value)
	}
	elem := reflect.ValueOf(data)
	// if elem.Kind() is Ptr, it will be en bind variable, otherwise it will be a temporary variable
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
		updateElem(elem, reflect.ValueOf(value))
	} else {
		// update temporary variable
		elem = reflect.ValueOf(ctx.data)
		setMapValue(elem, fields[0], value)
	}
	return nil
}

func getAttribute(obj interface{}, fieldNames []string) (interface{}, error) {
	value := reflect.ValueOf(obj)
	var attrVal reflect.Value
	fieldName := fieldNames[0]
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Map:
		attrVal = value.MapIndex(reflect.ValueOf(fieldName))
	case reflect.Slice:
		index, err := strconv.Atoi(fieldName)
		if nil != err {
			return nil, errors.New("Data is slice but field is not int")
		}
		attrVal = value.Index(index)
	case reflect.Interface:
		switch x := value.Interface().(type) {
		case map[string]interface{}, map[int]interface{}:
			attrVal = reflect.ValueOf(x).MapIndex(reflect.ValueOf(fieldName))
		case []interface{}:
			index, err := strconv.Atoi(fieldName)
			if nil != err {
				return nil, errors.New("Data is slice but field is not int")
			}
			attrVal = reflect.ValueOf(x).Index(index)
		}
	default:
		return nil, errors.New("Only support map and slice type")
	}
	if reflect.Map == value.Kind() {
		attrVal = value.MapIndex(reflect.ValueOf(fieldName))
	} else if reflect.Slice == value.Kind() {
		index, err := strconv.Atoi(fieldName)
		if nil != err {
			return nil, errors.New("Data is slice but field is not int")
		}
		attrVal = value.Index(index)
	} else {
		return nil, errors.New("Only support map and slice type")
	}

	if !attrVal.IsValid() {
		return nil, errors.New("not found")
	}

	ret := convToRealType(attrVal)
	if len(fieldNames) > 1 {
		return getAttribute(ret, fieldNames[1:])
	}
	return ret, nil
}

func setAttribute(obj interface{}, fieldNames []string, value interface{}) error {
	elem := reflect.ValueOf(obj)
	if reflect.Ptr != elem.Kind() {
		return errors.New("Target variable is not editable")
	}

	var err error
	fieldLen := len(fieldNames)
	elem = elem.Elem()
	for i, fieldName := range fieldNames {
		var v interface{}
		if fieldLen-1 == i {
			v = value
		}

		switch elem.Kind() {
		case reflect.Map:
			elem, err = setMapValue(elem, fieldName, v)
		case reflect.Slice:
			elem, err = setSliceValue(elem, fieldName, v)
		case reflect.Interface:
			switch x := elem.Interface().(type) {
			case map[string]interface{}, map[int]interface{}:
				elem, err = setMapValue(reflect.ValueOf(x), fieldName, v)
			case []interface{}, []string, []int, []float32, []float64, []bool:
				elem, err = setSliceValue(reflect.ValueOf(x), fieldName, v)
			default:
				err = errors.New("Interface change type not support")
			}
		default:
			err = errors.New("Data only support map and slice type")
		}
		if nil != err {
			return err
		}
	}

	return nil
}

func convToRealType(v reflect.Value) interface{} {
	switch x := v.Interface().(type) {
	default:
		return x
	}
}

func valueToInterface(v reflect.Value) interface{} {
	switch v.Type().Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(v.Uint())
	case reflect.Float32, reflect.Float64:
		return float32(v.Float())
	case reflect.Map:
		return v.Interface()
	case reflect.Array:
		return v.Interface()
	case reflect.Slice:
		return v.Interface()
	case reflect.Ptr:
		newPtr := reflect.New(v.Elem().Type())
		newPtr.Elem().Set(v.Elem())
		return newPtr.Interface()
	case reflect.Struct:
		if v.CanInterface() {
			return v.Interface()
		}
		return nil
	case reflect.Interface:
		return v.Interface()
	default:
		return nil
	}
}

func init() {
	typeInterface = reflect.TypeOf((*interface{})(nil)).Elem()
	typeMapStrInterface = reflect.MapOf(reflect.TypeOf(""), typeInterface)
}

func setMapValue(elem reflect.Value, fieldName string, value interface{}) (ret reflect.Value, err error) {
	tmpElem := elem.MapIndex(reflect.ValueOf(fieldName))
	if tmpElem.IsValid() {
		if nil == value {
			ret = tmpElem
			return
		}

		err = updateMapElem(elem, tmpElem.Type(), reflect.ValueOf(fieldName), reflect.ValueOf(value))
		return
	}

	if nil == value {
		// Create map[string]interface{}
		tmpElem = reflect.MakeMap(typeMapStrInterface)
		elem.SetMapIndex(reflect.ValueOf(fieldName), tmpElem)
		ret = tmpElem
	} else {
		err = updateMapElem(elem, elem.Type().Elem(), reflect.ValueOf(fieldName), reflect.ValueOf(value))
		ret = elem
	}

	return
}

func setSliceValue(elem reflect.Value, fieldName string, value interface{}) (ret reflect.Value, err error) {
	var index int
	index, err = strconv.Atoi(fieldName)
	if nil != err {
		return
	}

	// Out of range check
	if index >= elem.Len() {
		err = errors.New("Index out of range")
		return
	}

	tmpElem := elem.Index(index)
	if tmpElem.IsValid() {
		if nil == value {
			ret = tmpElem
			return
		}

		updateElem(tmpElem, reflect.ValueOf(value))
		return
	}

	err = errors.New("Do not increase slice")
	return
}

func canBeInt(value string) bool {
	_, err := strconv.Atoi(value)
	if nil != err {
		return false
	}
	return true
}

func updateElem(elem reflect.Value, value reflect.Value) error {
	if reflect.Interface == elem.Kind() {
		elem.Set(value)
	} else {
		typeElem := elem.Type()
		typeValue := value.Type()
		if !typeValue.ConvertibleTo(typeElem) {
			return fmt.Errorf("Assign type not math and can not be convert, targetType: %s, sourceType: %s", typeElem, typeValue)
		}
		elem.Set(value.Convert(typeElem))
	}
	return nil
}

func updateMapElem(elem reflect.Value, targetType reflect.Type, key reflect.Value, value reflect.Value) error {
	valueType := value.Type()
	if (targetType != valueType) && !valueType.ConvertibleTo(targetType) {
		return fmt.Errorf("Can not set value, variable type do not match, targetType: %s, sourceType: %s", targetType, valueType)
	}
	if targetType.Kind() != reflect.Interface && targetType != valueType {
		value = value.Convert(targetType)
	}

	elem.SetMapIndex(key, value)
	return nil
}
