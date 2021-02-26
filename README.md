# Geval

**geval** is a engine that can let you run golang's code at runtime. It use package `reflect`, `go/ast`, `go/parser` to parse golang code realtime. Frequently-used grammar is support although benchmark is not good

## Quick Start
#### Download and install
    go get github.com/MagicYH/geval

#### Usage
```go
package main

import (
    "github.com/MagicYH/geval"
    "fmt"
)

func main(){
	a := 0
	rule := `
	a = 1
	`
	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("a", &a)

	node, err := geval.NewRuleNode(rule, nil)
	if nil != err {
		fmt.Println("New rule error: ", err)
		return
	}

	err = node.Eval(dataCtx)
	if nil != err {
		fmt.Println("Eval error: ", err)
		return
	}
	fmt.Println(a)
}
```

look [test](https://github.com/MagicYH/geval/tree/master/test) for more example

## Features

[x] **Value assignment**: Include map assignment and struct assignment

[x] **Function call**: Include struct function. Self define function inject is support

[x] **If block**: >, >=, <, <=, ==, !=, +, -, *, /

[x] **For block**: `break`, `continue` is support

[x] **Create slice, map**: Can create slice and map with base type (int, string, float). For example: `a := make(map[string]int)`, `a := []int{1, 2, 3}`

### Function inject
```go
package main

import (
    "github.com/MagicYH/geval"
    "fmt"
)

type Person struct {
	Name string
}

func (p *Person) Say(content string) string {
	out := fmt.Sprintf("%s: %s", p.Name, content)
	return out
}

func main(){
    d := make(map[string]float)
	p := Person{Name: "Lilei"}
	rule := `
		d["a"] = Pow(2, 3)
		p.Say("Hello")
	`

	dataCtx := geval.NewDataCtx()
	dataCtx.Bind("d", &d)
	dataCtx.Bind("p", &p)

	funCtx := geval.NewFunCtx()
	funCtx.Bind("Pow", math.Pow)

	node, err := geval.NewRuleNode(rule, funCtx)
	if nil != err {
		t.Error("New rule error: ", err)
		return
	}

	node.Eval(dataCtx)
	fmt.Println(d)
}
```