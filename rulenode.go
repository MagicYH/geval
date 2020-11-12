package geval

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"

	"github.com/modern-go/reflect2"
)

// RuleNode : base element of rule node
type RuleNode struct {
	fset    *token.FileSet
	astFile *ast.File
	dataCtx *DataContext
	funcCtx *FunContext
}

const TOKEN_BREAK = "TOKEN BREAK"
const TOKEN_CONTINUE = "TOKEN CONTINUE"

var nilValue reflect.Value

// NewRuleNode : Create a new rule node
func NewRuleNode(content string, funcCtx *FunContext) (*RuleNode, error) {
	src := "package main\nfunc main() {\n" + content + "\n}"

	var err error
	ruleNode := &RuleNode{}
	ruleNode.fset = token.NewFileSet()
	ruleNode.funcCtx = funcCtx
	ruleNode.astFile, err = parser.ParseFile(ruleNode.fset, "", src, parser.AllErrors)
	return ruleNode, err
}

// Eval : Run a node
func (ruleNode *RuleNode) Eval(dataCtx *DataContext) (err error) {
	ruleNode.dataCtx = dataCtx
	ast.Inspect(ruleNode.astFile.Decls[0].(*ast.FuncDecl), func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			// first func declare is main func
			_, err = ruleNode.evalBlockStmt(x.Body)
			return false

		default:
			return true
		}
	})

	return
}

func (ruleNode *RuleNode) evalBody(bodyNode *ast.BlockStmt) error {
	for _, stmt := range bodyNode.List {
		ruleNode.eval(stmt)
	}
	return nil
}

func (ruleNode *RuleNode) eval(node ast.Node) (reflect.Value, error) {
	switch n := node.(type) {
	case *ast.AssignStmt:
		err := ruleNode.evalAssignStmt(n)
		if nil != err {
			return nilValue, err
		}
		return nilValue, nil
		// case *ast.ExprStmt:
		// 	return ruleNode.evalExprStmt(n)
		// case *ast.BinaryExpr:
		// 	return ruleNode.evalBinaryExpr(n)
		// case *ast.ParenExpr:
		// 	return ruleNode.evalParenExpr(n)
		// case *ast.BasicLit:
		// 	return ruleNode.evalBasicLit(n)
		// case *ast.CallExpr:
		// 	return ruleNode.evalCallExpr(n)
		// case *ast.Ident:
		// 	return ruleNode.evalIdent(n)
		// case *ast.SelectorExpr:
		// 	return ruleNode.evalSelectorExpr(n)
		// case *ast.IfStmt:
		// 	return ruleNode.evalIfStmt(n)
		// case *ast.BlockStmt:
		// 	return ruleNode.evalBlockStmt(n)
		// case *ast.IndexExpr:
		// 	return ruleNode.evalIndexExpr(n)
		// case *ast.ForStmt:
		// 	return ruleNode.evalForStmt(n)
		// case *ast.IncDecStmt:
		// 	return ruleNode.evalIncDecStmt(n)
		// case *ast.BranchStmt:
		// 	return ruleNode.evalBranchStmt(n)
	}
	return nilValue, fmt.Errorf("Node type not support: %s", reflect.TypeOf(node).String())
}

func (ruleNode *RuleNode) evalAssignStmt(node *ast.AssignStmt) error {
	// Just support one element right side
	if 1 != len(node.Rhs) {
		return errors.New("Rhs's length should be one")
	}

	rValue, err := ruleNode.getData(node.Rhs[0])
	if nil != err {
		return err
	}

	// switch n := node.Rhs[0].(type) {
	// case *ast.CallExpr:
	// 	vFunc, err := ruleNode.getFunc(n)
	// 	if nil != err {
	// 		return err
	// 	}
	// 	if vFunc.Type().NumOut() != len(node.Lhs) {
	// 		nameFunc, _ := ruleNode.eval(n.Fun)
	// 		return fmt.Errorf("Result element number is not equal with function `%s` result count, expect: %d, real: %d", nameFunc.String(), len(node.Lhs), vFunc.Type().NumOut())
	// 	}

	// 	for i, setNode := range node.Lhs {
	// 		// get the value data, and convert to reflect.Value
	// 		err := ruleNode.setData(setNode, getInterfaceRealValue(rValue.Index(i).Interface().(reflect.Value)), node.Tok)
	// 		if nil != err {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// }

	return ruleNode.setData(node.Lhs[0], rValue, node.Tok)
}

func (ruleNode *RuleNode) evalParenExpr(node *ast.ParenExpr) (ret reflect.Value, err error) {
	return ruleNode.eval(node.X)
}

func (ruleNode *RuleNode) evalBinaryExpr(node *ast.BinaryExpr) (ret interface{}, err error) {
	var left, right reflect.Value
	// left, err = ruleNode.eval(node.X)
	left, err = ruleNode.getData(node.X)
	if nil != err {
		return
	}
	// If node.X is callExpr, than get the first result
	// left = getNodeFirstResult(node.X, left)

	// right, err = ruleNode.eval(node.Y)
	right, err = ruleNode.getData(node.Y)
	if nil != err {
		return
	}
	// If node.Y is callExpr, than get the first result
	// right = getNodeFirstResult(node.Y, right)

	switch node.Op.String() {
	case "+":
		return add(left, right)
	case "==":
		return equ(left, right)
	case "!=":
		return neq(left, right)
	case "-", "*", "/", "<", ">", "<=", ">=":
		return doNumMath(left, right, node.Op.String())
	}
	return nilValue, errors.New("Operate not define")
}

func (ruleNode *RuleNode) evalBasicLit(node *ast.BasicLit) (ret interface{}, err error) {
	switch node.Kind {
	case token.INT:
		return strconv.Atoi(node.Value)
	case token.FLOAT:
		return strconv.ParseFloat(node.Value, 64)
	case token.STRING:
		str := node.Value
		strLen := len(str)
		return str[1 : strLen-1], nil
	}
	return nil, fmt.Errorf("Basic token not support: %d", node.Kind)
}

func (ruleNode *RuleNode) evalExpr() (ret reflect.Value, err error) {
	return
}

func (ruleNode *RuleNode) evalExprStmt(node *ast.ExprStmt) (reflect.Value, error) {
	return ruleNode.eval(node.X)
}

func (ruleNode *RuleNode) evalIfStmt(node *ast.IfStmt) (ret reflect.Value, err error) {
	// run init
	if nil != node.Init {
		ruleNode.eval(node.Init)
	}
	// cond, err := ruleNode.eval(node.Cond)
	cond, err := ruleNode.getData(node.Cond)
	if nil != err {
		return nilValue, err
	}

	if cond.Bool() {
		return ruleNode.eval(node.Body)
	}
	if nil != node.Else {
		return ruleNode.eval(node.Else)
	}
	return nilValue, nil
}

func (ruleNode *RuleNode) evalBlockStmt(node *ast.BlockStmt) (ret reflect.Value, err error) {
	for _, stmt := range node.List {
		_, err = ruleNode.eval(stmt)
		if nil != err {
			return
		}
	}
	return
}

func (ruleNode *RuleNode) evalCallExpr(node *ast.CallExpr) (ret reflect.Value, err error) {
	var vFunc reflect.Value
	vFunc, err = ruleNode.getFunc(node)
	if nil != err {
		return
	}
	tFunc := vFunc.Type()
	numIn := tFunc.NumIn()

	var tLastIn reflect.Type
	if numIn > 0 {
		tLastIn = tFunc.In(numIn - 1)
		if tFunc.IsVariadic() {
			tLastIn = tLastIn.Elem()
		}
	}

	nodeFunc, isSel := node.Fun.(*ast.SelectorExpr)
	realInNum := len(node.Args)
	if isSel {
		realInNum++
	}
	if !tFunc.IsVariadic() && realInNum != numIn {
		return nilValue, fmt.Errorf("Call udf input number not right, expect: %d, real: %d", numIn, len(node.Args))
	}

	args := make([]reflect.Value, 0, realInNum)
	if isSel {
		selStru, err := ruleNode.getData(nodeFunc.X)
		if nil != err {
			return nilValue, err
		}
		args = append(args, selStru)
	}
	for _, n := range node.Args {
		param, err := ruleNode.getData(n)
		if nil != err {
			return nilValue, fmt.Errorf("Get fun params error: %v", err)
		}

		var expectType reflect.Type
		i := len(args)
		if i < numIn-1 {
			expectType = tFunc.In(i)
		} else {
			expectType = tLastIn
		}
		param, err = typeConvert(param, expectType)
		if nil != err {
			return nilValue, err
		}

		args = append(args, param)
	}

	return reflect.ValueOf(vFunc.Call(args)), nil
}

func (ruleNode *RuleNode) evalIdent(node *ast.Ident) (ret reflect.Value, err error) {
	// var err error
	// identity := Identity{}
	// identity.Name = node.Name
	// switch node.Name {
	// case "true":
	// 	identity.Value = true
	// 	identity.Kind = identKindInbuild
	// case "false":
	// 	identity.Value = false
	// 	identity.Kind = identKindInbuild
	// case "nil":
	// 	identity.Value = nil
	// 	identity.Kind = identKindInbuild
	// default:
	// 	identity.Kind = identKindVar
	// }
	// return identity, err
	return reflect.ValueOf(node.Name), nil
}

func (ruleNode *RuleNode) evalSelectorExpr(node *ast.SelectorExpr) (vFunc reflect.Value, err error) {
	vX, err := ruleNode.getData(node.X)
	if nil != err {
		return nilValue, err
	}
	tX := vX.Type()
	method, ok := tX.MethodByName(node.Sel.Name)
	if !ok {
		err = fmt.Errorf("Method %s not found", node.Sel.Name)
		return
	}

	vFunc = method.Func
	return
}

func (ruleNode *RuleNode) evalIndexExpr(node *ast.IndexExpr) (ret reflect.Value, err error) {
	return
}

func (ruleNode *RuleNode) evalForStmt(node *ast.ForStmt) (ret reflect.Value, err error) {
	// init loop
	if nil != node.Init {
		_, err = ruleNode.eval(node.Init)
		if nil != err {
			return
		}
	}

	if nil == node.Cond {
		err = fmt.Errorf("Nil for cond is not allow")
		return
	}

	var cond reflect.Value
	for {
		cond, err = ruleNode.eval(node.Cond)
		if nil != err {
			return nilValue, err
		}
		if !cond.Bool() {
			return
		}

		_, err = ruleNode.eval(node.Body)
		if nil != err {
			switch err.Error() {
			case TOKEN_BREAK:
				return nilValue, nil

			case TOKEN_CONTINUE:
				err = nil

			default:
				return
			}
		}

		if nil != node.Post {
			_, err = ruleNode.eval(node.Post)
			if nil != err {
				return
			}
		}
	}
}

func (ruleNode *RuleNode) evalIncDecStmt(node *ast.IncDecStmt) (ret reflect.Value, err error) {
	vX, err := ruleNode.getData(node.X)
	if nil != err {
		return nilValue, err
	}

	switch node.Tok.String() {
	case "++", "--":
		ret, err = doNumMath(vX, reflect.ValueOf(0.0), node.Tok.String())
		if nil != err {
			return
		}

		err = ruleNode.setData(node.X, ret, token.ASSIGN)
		return
	}
	return nilValue, errors.New("Operate not define")
}

func (ruleNode *RuleNode) evalBranchStmt(node *ast.BranchStmt) (ret reflect.Value, err error) {
	switch node.Tok {
	case token.BREAK:
		err = errors.New(TOKEN_BREAK)
	case token.CONTINUE:
		err = errors.New(TOKEN_CONTINUE)
	default:
		err = fmt.Errorf("Branch token not support: %v", node.Tok)
	}
	return
}

func (ruleNode *RuleNode) evalCompositeLit(node *ast.CompositeLit) (vSlice reflect.Value, err error) {
	switch n := node.Type.(type) {
	case *ast.ArrayType:
		ident, ok := n.Elt.(*ast.Ident)
		if !ok {
			err = fmt.Errorf("Elt not *ast.Ident type, realtype: %T", n.Elt)
		}

		typeName := ident.Name
		length := len(node.Elts)
		var elem reflect.Value
		var targetType reflect.Type
		switch typeName {
		case "int":
			targetType = typeInt
			vSlice = reflect.MakeSlice(reflect.SliceOf(typeInt), 0, length)
		case "string":
			targetType = typeString
			vSlice = reflect.MakeSlice(reflect.SliceOf(typeString), 0, length)
		case "float32":
			targetType = typeFloat32
			vSlice = reflect.MakeSlice(reflect.SliceOf(typeFloat32), 0, length)
		case "float64":
			targetType = typeFloat64
			vSlice = reflect.MakeSlice(reflect.SliceOf(typeFloat64), 0, length)
		default:
			err = fmt.Errorf("Slice element type not support: %v", typeName)
			return
		}
		for _, expr := range node.Elts {
			elem, err = ruleNode.getData(expr)
			if nil != err {
				return nilValue, err
			}
			elem, err = typeConvert(elem, targetType)
			if nil != err {
				return nilValue, err
			}
			vSlice = reflect.Append(vSlice, elem)
		}
	}
	return
}

func (ruleNode *RuleNode) evalMapType(node *ast.MapType) (param reflect.Value, err error) {
	keyIdent, ok := node.Key.(*ast.Ident)
	if !ok {
		err = fmt.Errorf("node.Key is not *ast.Ident, realtype: %T", node.Key)
		return
	}
	valueIdent, ok := node.Value.(*ast.Ident)
	if !ok {
		err = fmt.Errorf("node.Value is not *ast.Ident, realtype: %T", node.Value)
		return
	}

	tKey, err := getTypeWithName(keyIdent.Name)
	if nil != err {
		return nilValue, err
	}
	tValue, err := getTypeWithName(valueIdent.Name)
	if nil != err {
		return nilValue, err
	}

	makeMapParam := makeMapParam{tKey: tKey, tValue: tValue}
	return reflect.ValueOf(makeMapParam), nil
}

func (ruleNode *RuleNode) getData(node ast.Expr) (ret interface{}, err error) {
	switch n := node.(type) {
	case *ast.Ident:
		return ruleNode.identGet(n)

	case *ast.IndexExpr:
		var x, index reflect.Value
		x, err = ruleNode.getData(n.X)
		if nil != err {
			return x, err
		}
		index, err := ruleNode.getData(n.Index)
		if nil != err {
			return index, err
		}
		return getDataByIndex(x, index)

	case *ast.SelectorExpr:
		ret, err = ruleNode.getData(n.X)
		return getDataBySel(x, n.Sel.Name)
	case *ast.BasicLit:
		return ruleNode.evalBasicLit(n)
	case *ast.BinaryExpr:
		return ruleNode.evalBinaryExpr(n)
	case *ast.ParenExpr:
		return ruleNode.getData(n.X)
		// case *ast.CallExpr:
		// 	return ruleNode.evalCallExpr(n)
		// case *ast.CompositeLit:
		// 	return ruleNode.evalCompositeLit(n)
		// case *ast.MapType:
		// 	return ruleNode.evalMapType(n)
	}
	err = fmt.Errorf("Unexpect get node type: %T, value: %v", node, node)
	return
}

func (ruleNode *RuleNode) identGet(node *ast.Ident) (value interface{}, err error) {
	switch node.Name {
	case "true":
		value = true
	case "false":
		value = false
	case "nil":
		value = nil
	default:
		value, err = ruleNode.dataCtx.Get(node.Name)
	}
	return
}

func getDataByIndex(data interface{}, index interface{}) (ret interface{}, err error) {
	tData := reflect2.TypeOf(data)
	kData := tData.Type1().Kind()

	switch kData {
	case reflect.Map:
		tMap := tData.(reflect2.MapType)
		return tMap.GetIndex(&data, &index), nil

	case reflect.Slice:
		tSlice := tData.(reflect2.SliceType)
		return tSlice.GetIndex(&data, index.(int)), nil
	}
	err = fmt.Errorf("Unexpect data kind when get by index: %v", kData)
	return
}

func getDataBySel(data interface{}, field string) (ret interface{}, err error) {
	tData := reflect2.TypeOf(data)
	tStruct := tData.(reflect2.StructType)
	return tStruct.FieldByName(field).Get(data), nil
}

func (ruleNode *RuleNode) setData(node ast.Expr, value interface{}git , t token.Token) (err error) {
	var elem reflect.Value
	switch n := node.(type) {
	case *ast.Ident:
		err = ruleNode.identSet(n, value, t)
	case *ast.IndexExpr:
		elem, err = ruleNode.getData(n.X)
		if nil != err {
			return
		}
		var index reflect.Value
		index, err = ruleNode.getData(n.Index)
		if nil != err {
			return
		}
		err = setDataByIndex(elem, index, value)

	case *ast.SelectorExpr:
		elem, err = ruleNode.getData(n.X)
		if nil != err {
			return
		}

		err = setDataBySel(elem, n.Sel.Name, value)

	default:
		err = fmt.Errorf("Unexpect set node type: %T, value: %v", node, node)
	}
	return
}

func (ruleNode *RuleNode) identSet(node *ast.Ident, value reflect.Value, t token.Token) (err error) {
	switch node.Name {
	case "true":
		err = fmt.Errorf("Can not set to true")
	case "false":
		err = fmt.Errorf("Can not set to false")
	case "nil":
		err = fmt.Errorf("Can not set to nil")
	default:
		err = ruleNode.dataCtx.Set(node.Name, value)
		// ret, err := ruleNode.dataCtx.Get(node.Name)
		// switch t {
		// case token.DEFINE:
		// 	if nil != ret {
		// 		err = fmt.Errorf("Elem has exsits before, should not be define again, name: %s", node.Name)
		// 	} else {
		// 		err = ruleNode.dataCtx.Set(node.Name, convToRealType(value))
		// 	}

		// case token.ASSIGN:
		// 	if nil != err {
		// 		err = fmt.Errorf("Elem has not been define, define first, name: %s", node.Name)
		// 	} else {
		// 		err = ruleNode.dataCtx.Set(node.Name, convToRealType(value))
		// 	}

		// default:
		// 	err = fmt.Errorf("Assign token not found: %v", t)
		// }
	}
	return
}

func (ruleNode *RuleNode) getFunc(node *ast.CallExpr) (vFunc reflect.Value, err error) {

	switch node.Fun.(type) {
	case *ast.SelectorExpr:
		return ruleNode.eval(node.Fun)
	default:
		var vFuncName reflect.Value
		vFuncName, err = ruleNode.eval(node.Fun)
		if nil != err {
			return
		}
		funcName := vFuncName.String()
		udf, ok := ruleNode.funcCtx.data[funcName]
		if !ok {
			return nilValue, fmt.Errorf("Call udf fail, udf not found: %s", funcName)
		}

		vFunc = reflect.ValueOf(udf)
	}

	return vFunc, nil
}

func setDataByIndex(vData reflect.Value, vIndex reflect.Value, vValue reflect.Value) (err error) {
	kData := vData.Kind()
	if reflect.Ptr == kData || reflect.Interface == kData {
		vData = vData.Elem()
		kData = vData.Kind()
	}

	switch kData {
	case reflect.Map:
		_, err = setMapValue(vData, vIndex, vValue)

	case reflect.Slice:
		_, err = setSliceValue(vData, vIndex, vValue)

	default:
		err = fmt.Errorf("Unsupport set by index type: %v", kData)
	}
	return
}

func setDataBySel(vData reflect.Value, field string, vValue reflect.Value) (err error) {
	kData := vData.Kind()
	if reflect.Ptr == kData || reflect.Interface == kData {
		vData = vData.Elem()
		kData = vData.Kind()
	}
	if reflect.Struct != kData {
		return fmt.Errorf("Unexpect data kind when set by sel, type: %T, value: %v", vData.Type(), vData)
	}

	tData := vData.Type()
	_, ok := tData.FieldByName(field)
	if !ok {
		return fmt.Errorf("Field %s not found", field)
	}
	elem := vData.FieldByName(field)
	vValue, err = typeConvert(vValue, elem.Type())
	if nil != err {
		return
	}
	elem.Set(vValue)
	return
}

func getInterfaceRealValue(vValue reflect.Value) reflect.Value {
	if reflect.Interface != vValue.Kind() {
		return vValue
	}

	return reflect.ValueOf(convToRealType(vValue))
}

func getInterfaceRealType(vValue reflect.Value) reflect.Type {
	if reflect.Interface != vValue.Kind() {
		return vValue.Type()
	}

	value := convToRealType(vValue)
	return reflect.TypeOf(value)
}

func typeConvert(vValue reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	tValue := vValue.Type()
	if targetType != tValue && !tValue.ConvertibleTo(targetType) {
		return vValue, fmt.Errorf("Can not set value, variable type do not match, targetType: %s, sourceType: %s", targetType, tValue)
	}
	if targetType.Kind() != reflect.Interface && targetType != tValue {
		vValue = vValue.Convert(targetType)
	}
	return vValue, nil
}

func getNodeFirstResult(node ast.Expr, vValue reflect.Value) reflect.Value {
	switch node.(type) {
	case *ast.CallExpr:
		if vValue.Len() <= 0 {
			return nilValue
		}
		return vValue.Index(0).Interface().(reflect.Value)
	}
	return vValue
}

func getTypeWithName(name string) (t reflect.Type, err error) {
	switch name {
	case "int":
		t = typeInt
	case "string":
		t = typeString
	case "float32":
		t = typeFloat32
	case "float64":
		t = typeFloat64
	default:
		err = fmt.Errorf("Slice element type not support: %v", name)
		return
	}
	return
}

// DumpAstTree : Dump node's ast tree
func (ruleNode *RuleNode) DumpAstTree() {
	fset := token.NewFileSet()
	ast.Print(fset, ruleNode.astFile)
}

func init() {
	nilValue = reflect.ValueOf(nil)
}
