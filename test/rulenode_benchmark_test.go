package test

import (
	"testing"

	"github.com/MagicYH/geval"
)

func BenchmarkEvalLoop(b *testing.B) {
	rule := `
	a := 0
	for i := 0; i < bn; i++ {
		a++
	}
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("bn", &b.N)

	funCtx := geval.NewFunCtx()
	node, err := geval.NewRuleNode(rule, funCtx)
	if nil != err {
		b.Error("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		b.Error("Eval error: ", err)
		return
	}
}

func BenchmarkMathEval(b *testing.B) {
	requestMade := 99.0
	requestSucceeded := 90.0
	rule := `
	(requestMade * requestSucceeded / 100) >= 90 
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("requestMade", &requestMade)
	dataCtx.Bind("requestSucceeded", &requestSucceeded)

	node, err := geval.NewRuleNode(rule, nil)
	if nil != err {
		b.Error("New rule error: ", err)
		return
	}
	for i := 0; i < b.N; i++ {
		node.Eval(dataCtx)
	}
}

func BenchmarkMakeSliceEval(b *testing.B) {
	a := 0
	rule := `
	a = []int{1, 2, 3}
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("a", &a)

	funCtx := geval.NewFunCtx()
	node, err := geval.NewRuleNode(rule, funCtx)
	if nil != err {
		b.Error("New rule error: ", err)
		return
	}
	for i := 0; i < b.N; i++ {
		node.Eval(dataCtx)
	}
}

func BenchmarkRealEval(b *testing.B) {
	data := make(map[string]interface{})
	data["aspect"] = 0.77
	data["cep"] = []int{0, 495}
	content := `
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
		b.Error("New rule error: ", err)
		return
	}
	for i := 0; i < b.N; i++ {
		node.Eval(dataCtx)
	}
	// b.Log(data)
	// if data["bev"] != 0 {
	// 	b.Error("Result error")
	// }
}
