package test

import (
	"testing"

	"github.com/MagicYH/geval"
)

func BenchmarkLoop(b *testing.B) {
	a := 0
	for i := 0; i < b.N; i++ { //use b.N for looping
		a = a + 1
	}
}

func BenchmarkMakSlice(b *testing.B) {
	var a []int
	for i := 0; i < b.N; i++ {
		a = []int{1, 2, 3}
	}
	if len(a) > 0 {

	}
}

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
	a := 0
	rule := `
	a++
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

	if a != b.N {
		b.Error("a result error")
		return
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
