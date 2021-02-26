package test

import (
	"fmt"
	"math"
	"testing"

	"github.com/MagicYH/geval"
)

type Position struct {
	Name string
}

type Profession struct {
	Name     string
	Salary   int
	Position Position
}

type Person struct {
	Name string
	Age  int
	Pro  Profession
}

func (p *Person) Say(content string) string {
	out := fmt.Sprintf("%s: %s", p.Name, content)
	return out
}

func TestAssignRule(t *testing.T) {
	dict := make(map[string]interface{})
	dict["int"] = 10
	dict["sub"] = make(map[string]interface{})
	dict["sub"].(map[string]interface{})["int"] = 2
	value := 0
	// stru := Person{}
	stru := Person{Pro: Profession{Name: "ss", Salary: 5}}
	stru.Pro.Name = "bbbb"

	rule := `
	tmp := "tmp"

	dict["int"] = 1
	dict["str"] = tmp
	dict["sub"]["int"] = 1
	dict["sub"]["str"] = tmp

	value = 10

	stru.Name = "Lilei"
	stru.Age = 20
	stru.Pro.Name = "student"
	stru.Pro.Salary = 10
	stru.Pro.Position.Name = "PM"
	`
	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("dict", &dict)
	dataCtx.Bind("value", &value)
	dataCtx.Bind("stru", &stru)

	node, err := geval.NewRuleNode(rule, nil)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		t.Error("Eval error: ", err)
		return
	}

	t.Log("dict: ", dict)
	t.Log("value: ", value)
	t.Log("stru", stru)

	if dict["int"] != 1 || dict["str"] != "tmp" || dict["sub"].(map[string]interface{})["int"] != 1 || dict["sub"].(map[string]interface{})["str"] != "tmp" {
		t.Error("Unexpect dict result")
		return
	}

	if value != 10 {
		t.Error("Unexpect value result")
		return
	}

	if stru.Name != "Lilei" || stru.Age != 20 || stru.Pro.Name != "student" || stru.Pro.Salary != 10 || stru.Pro.Position.Name != "PM" {
		t.Error("Unexpect stru result")
		return
	}
}

func TestMathAssign(t *testing.T) {
	a := 0
	b := 0
	c := 0
	d := make(map[string]int)
	rule := `
	a = 1 + 1.9
	b = 3 * 4
	c = 24 / 3
	d["a"] = 1 + 1
	d["b"] = 2 * (3 - 2) + 4 / 2
	d["c"] = 0
	d["c"]++
	d["c"]++
	d["c"]++
	`
	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("a", &a)
	dataCtx.Bind("b", &b)
	dataCtx.Bind("c", &c)
	dataCtx.Bind("d", &d)

	node, err := geval.NewRuleNode(rule, nil)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		t.Error("Eval error: ", err)
		return
	}
	t.Log("a: ", a)
	t.Log("b: ", b)
	t.Log("c: ", c)
	t.Log("d: ", d)
	if a != 2 || b != 12 || c != 8 {
		t.Error("result not right")
		return
	}

	if d["a"] != 2 || d["b"] != 4 || d["c"] != 3 {
		t.Error("result map not right")
		return
	}
}

func TestIf(t *testing.T) {
	a := 0
	b := 0
	c := 0
	rule := `
	one := 1

	// a = 0
	// b = 0
	// c = 0

	if one == 1 && (one == 1 || one == 2) {
		a = 1
	} else {
		a = 0
	}

	if one == 2 {
		b = 0
	} else {
		b = 1
	}

	if one == 3 {
		c = 0
	} else if one == 1 {
		c = 1
	} else {
		c = 0
	}
	`
	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("a", &a)
	dataCtx.Bind("b", &b)
	dataCtx.Bind("c", &c)

	node, err := geval.NewRuleNode(rule, nil)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	node.DumpAstTree()
	return
	err = node.Eval(dataCtx)
	if nil != err {
		t.Error("Eval error: ", err)
		return
	}
	t.Log("a: ", a)
	t.Log("b: ", b)
	t.Log("c: ", c)
	if a != 1 || b != 1 || c != 1 {
		t.Error("result not right")
		return
	}
}

func TestFun(t *testing.T) {
	d := make(map[string]interface{})
	p := Person{Name: "Lilei"}
	rule := `
		d["a"] = Pow(2, 3)
		d["b"] = "world"
		d["c"] = Sprintf("Hello %s, %f", d["b"], d["a"])
		d["d"] = 0
		d["e"], d["f"] = doubleAssign(10, 20)
		d["g"] = len("Hello world")

		if d["a"] >= Max(1, 8) {
			d["d"] = 1
		}

		p.Say("Hello")
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("d", &d)
	dataCtx.Bind("p", &p)

	funCtx := geval.NewFunCtx()
	funCtx.Bind("Sprintf", fmt.Sprintf)
	funCtx.Bind("Pow", math.Pow)
	funCtx.Bind("Max", math.Max)
	funCtx.Bind("doubleAssign", doubleAssign)

	node, err := geval.NewRuleNode(rule, funCtx)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		t.Error("Eval error: ", err)
		return
	}

	t.Log(d)
	if 8 != d["a"].(float64) {
		t.Error("Result a error")
		return
	}
	if "Hello world, 8.000000" != d["c"] {
		t.Error("Result c error")
		return
	}
	if 1 != d["d"].(int) || d["g"].(int) != 11 {
		t.Error("Result d error")
		return
	}
	if 10 != d["e"].(int) || 20 != d["f"].(int) {
		t.Error("Result e, f error")
		return
	}
}

func TestFor(t *testing.T) {
	d := make(map[string]int)
	rule := `
	d["a"] = 0
	d["b"] = 0
	d["c"] = 0
	for i := 0; i < 100; i++ {
		d["a"]++
	}

	for i := 0; i < 100; i++ {
		if d["b"] > 49 {
			break
		}
		d["b"]++
	}


	for i := 0; i < 100; i++ {
		if d["c"] > 59 {
			continue
		}
		d["c"]++
	}
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("d", &d)

	funCtx := geval.NewFunCtx()

	node, err := geval.NewRuleNode(rule, funCtx)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		t.Error("Eval error: ", err)
		return
	}

	t.Log(d)
	if d["a"] != 100 || d["b"] != 50 || d["c"] != 60 {
		t.Error("Result error")
		return
	}
}

func TestCompositeLit(t *testing.T) {
	var si []int
	var ss []string
	var sf []float32
	rule := `
	si = []int{1, 2, 3}
	ss = []string{"s1", "s2", "s3"}
	sf = []float32{1.1, 2.2, 3.3}
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("si", &si)
	dataCtx.Bind("ss", &ss)
	dataCtx.Bind("sf", &sf)

	funCtx := geval.NewFunCtx()

	node, err := geval.NewRuleNode(rule, funCtx)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		t.Error("Eval error: ", err)
		return
	}

	t.Log(si)
	t.Log(ss)
	t.Log(sf)
	if len(si) != 3 || si[0] != 1 || si[1] != 2 || si[2] != 3 {
		t.Error("si result error")
		return
	}
	if len(ss) != 3 || ss[0] != "s1" || ss[1] != "s2" || ss[2] != "s3" {
		t.Error("si result error")
		return
	}
	if len(sf) != 3 || sf[0] != 1.1 || sf[1] != 2.2 || sf[2] != 3.3 {
		t.Error("si result error")
		return
	}
}

func TestMap(t *testing.T) {
	var si map[string]int
	var ss map[string]string
	var sf map[string]float32
	var ii map[int]int
	rule := `
	si = make(map[string]int)
	si["a"] = 1

	ss = make(map[string]string)
	ss["a"] = "a"

	sf = make(map[string]float32)
	sf["a"] = 1.1

	ii = make(map[int]int)
	ii[10] = 10
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("si", &si)
	dataCtx.Bind("ss", &ss)
	dataCtx.Bind("sf", &sf)
	dataCtx.Bind("ii", &ii)

	funCtx := geval.NewFunCtx()

	node, err := geval.NewRuleNode(rule, funCtx)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		t.Error("Eval error: ", err)
		return
	}

	t.Log(si)
	t.Log(ss)
	t.Log(sf)
	t.Log(ii)
	if si["a"] != 1 {
		t.Error("si result error")
		return
	}

	if ss["a"] != "a" {
		t.Error("ss result error")
		return
	}

	if sf["a"] != 1.1 {
		t.Error("sf result error")
		return
	}

	if ii[10] != 10 {
		t.Error("ii result error")
		return
	}
}

func TestRealEval(t *testing.T) {
	data := make(map[string]interface{})
	data["aspect"] = 1.77
	data["cep"] = []int{0, 495}
	content := `
	data["bev"] = 0
	if data["aspect"] < 1 {
		data["bev"] = 1
	} else {
		cep := data["cep"]
		if ((cep[1] - cep[0]) / 495) <= 0.85 {
			data["bev"] = 1
		} else {
			data["bev"] = 0
		}
	}
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("data", &data)

	node, err := geval.NewRuleNode(content, nil)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}
	node.Eval(dataCtx)
	t.Log(data)

	if data["bev"] != 0 {
		t.Error("Result error")
	}
}

func doubleAssign(a, b int) (int, int) {
	return a, b
}
