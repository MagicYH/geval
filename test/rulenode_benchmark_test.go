package test

import (
	"testing"

	"github.com/MagicYH/geval"
)

// func BenchmarkLoop(b *testing.B) {
// 	a := 0
// 	for i := 0; i < b.N; i++ { //use b.N for looping
// 		a++
// 	}
// }

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
