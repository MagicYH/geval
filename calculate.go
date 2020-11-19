package geval

import (
	"errors"
	"fmt"
	"go/token"
	"reflect"
)

var numDict map[reflect.Kind]bool
var intDict map[reflect.Kind]bool
var typeFloat64 reflect.Type
var typeFloat32 reflect.Type
var typeInt reflect.Type
var typeString reflect.Type

func add(a, b interface{}) (interface{}, error) {
	if reflect.String == reflect.TypeOf(a).Kind() && reflect.String == reflect.TypeOf(b).Kind() {
		return fmt.Sprintf("%s%s", a, b), nil
	} else if IsNumber(reflect.TypeOf(a).Kind()) && IsNumber(reflect.TypeOf(b).Kind()) {
		return doNumMath(a, b, token.ADD.String())
	}
	return nil, fmt.Errorf("Math type not support with %T, %T", a, b)
}

func equ(a, b interface{}) (bool, error) {
	return reflect.DeepEqual(a, b), nil
}

func neq(a, b interface{}) (bool, error) {
	return !reflect.DeepEqual(a, b), nil
}

func getValueAndKind(input interface{}) (reflect.Value, reflect.Kind) {
	v := reflect.ValueOf(input)
	v = reflect.Indirect(v)
	return v, v.Kind()
}

// func addString(vA, vB reflect.Value) (reflect.Value, error) {
// 	kA := vA.Kind()
// 	kB := vB.Kind()
// 	if kA == reflect.String && kB != reflect.String {
// 		return nilValue, fmt.Errorf("ADD(+) can't be used between %s and %s", kA, kB)
// 	} else if kA != reflect.String && kB == reflect.String {
// 		return nilValue, fmt.Errorf("ADD(+) can't be used between %s and %s", kA, kB)
// 	}
// 	return reflect.ValueOf(fmt.Sprintf("%s%s", vA.String(), vB.String())), nil
// }

func doNumMath(a, b interface{}, op string) (interface{}, error) {
	afloat, err := interToFloat(a)
	if nil != err {
		return nil, err
	}
	bfloat, err := interToFloat(b)
	if nil != err {
		return nil, err
	}

	var ret interface{}
	switch op {
	case "+":
		ret = afloat + bfloat
	case "-":
		ret = afloat - bfloat
	case "*":
		ret = afloat * bfloat
	case "/":
		if 0 == bfloat {
			return nil, errors.New("Can not div with zero(0) value")
		}
		ret = afloat / bfloat
	case "<":
		ret = (afloat < bfloat)
	case ">":
		ret = (afloat > bfloat)
	case "<=":
		ret = (afloat <= bfloat)
	case ">=":
		ret = (afloat >= bfloat)
	case "++":
		ret = afloat + 1
	case "--":
		ret = afloat - 1
	default:
		return nil, fmt.Errorf("Operate(%s) not support", op)
	}
	return ret, nil
}

func interToFloat(v interface{}) (float64, error) {
	switch v.(type) {
	case int:
		return float64(v.(int)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	case *int:
		return float64(*v.(*int)), nil
	case *float32:
		return float64(*v.(*float32)), nil
	case *float64:
		return *v.(*float64), nil
	}
	return 0, fmt.Errorf("Can not conver %s value to float64", reflect.TypeOf(v).Elem())
}

func init() {
	initNumDict()
	initIntDict()

	typeFloat64 = reflect.TypeOf(float64(0.0))
	typeFloat32 = reflect.TypeOf(float32(0.0))
	typeString = reflect.TypeOf("")
	typeInt = reflect.TypeOf(int(0))
}

func initNumDict() {
	numDict = make(map[reflect.Kind]bool)
	numDict[reflect.Int] = true
	numDict[reflect.Int8] = true
	numDict[reflect.Int16] = true
	numDict[reflect.Int32] = true
	numDict[reflect.Int64] = true
	numDict[reflect.Uint] = true
	numDict[reflect.Uint8] = true
	numDict[reflect.Uint16] = true
	numDict[reflect.Uint32] = true
	numDict[reflect.Uint64] = true
	numDict[reflect.Float32] = true
	numDict[reflect.Float64] = true
}

func initIntDict() {
	intDict = make(map[reflect.Kind]bool)
	intDict[reflect.Int] = true
	intDict[reflect.Int8] = true
	intDict[reflect.Int16] = true
	intDict[reflect.Int32] = true
	intDict[reflect.Int64] = true
	intDict[reflect.Uint] = true
	intDict[reflect.Uint8] = true
	intDict[reflect.Uint16] = true
	intDict[reflect.Uint32] = true
	intDict[reflect.Uint64] = true
}

// IsNumber : Check if kind is number
func IsNumber(kind reflect.Kind) bool {
	_, ok := numDict[kind]
	return ok
}

// IsInt : Check kind is int
func IsInt(kind reflect.Kind) bool {
	_, ok := intDict[kind]
	return ok
}
