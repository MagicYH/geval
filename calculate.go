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

func add(vA, vB reflect.Value) (reflect.Value, error) {
	if reflect.String == vA.Kind() || reflect.String == vB.Kind() {
		return addString(vA, vB)
	}
	return doNumMath(vA, vB, token.ADD.String())
}

func equ(vA, vB reflect.Value) (reflect.Value, error) {
	return reflect.ValueOf(reflect.DeepEqual(vA.Interface(), vB.Interface())), nil
}

func neq(vA, vB reflect.Value) (reflect.Value, error) {
	return reflect.ValueOf(!reflect.DeepEqual(vA.Interface(), vB.Interface())), nil
}

func getValueAndKind(input interface{}) (reflect.Value, reflect.Kind) {
	v := reflect.ValueOf(input)
	v = reflect.Indirect(v)
	return v, v.Kind()
}

func addString(vA, vB reflect.Value) (reflect.Value, error) {
	kA := vA.Kind()
	kB := vB.Kind()
	if kA == reflect.String && kB != reflect.String {
		return nilValue, fmt.Errorf("ADD(+) can't be used between %s and %s", kA, kB)
	} else if kA != reflect.String && kB == reflect.String {
		return nilValue, fmt.Errorf("ADD(+) can't be used between %s and %s", kA, kB)
	}
	return reflect.ValueOf(fmt.Sprintf("%s%s", vA.String(), vB.String())), nil
}

func doNumMath(vA, vB reflect.Value, op string) (reflect.Value, error) {
	kA := vA.Kind()
	kB := vB.Kind()
	if reflect.Interface == kA || reflect.Ptr == kA {
		vA = vA.Elem()
		kA = vA.Kind()
	}
	if reflect.Interface == kB || reflect.Ptr == kB {
		vB = vB.Elem()
		kB = vB.Kind()
	}
	if !IsNumber(kA) || !IsNumber(kB) {
		return nilValue, fmt.Errorf("OP(%s) can't be used between %s and %s", op, kA, kB)
	}

	var ret interface{}
	afloat := vA.Convert(typeFloat64).Float()
	bfloat := vB.Convert(typeFloat64).Float()

	switch op {
	case "+":
		ret = afloat + bfloat
	case "-":
		ret = afloat - bfloat
	case "*":
		ret = afloat * bfloat
	case "/":
		if 0 == bfloat {
			return nilValue, errors.New("Can not div with zero(0) value")
		}
		ret = afloat / bfloat
	case "==":
		ret = (afloat == bfloat)
	case "<":
		ret = (afloat < bfloat)
	case ">":
		ret = (afloat > bfloat)
	case "!=":
		ret = (afloat != bfloat)
	case "<=":
		ret = (afloat <= bfloat)
	case ">=":
		ret = (afloat >= bfloat)
	case "++":
		ret = afloat + 1
	case "--":
		ret = afloat - 1
	default:
		return nilValue, fmt.Errorf("Operate(%s) not support", op)
	}
	return reflect.ValueOf(ret), nil
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
