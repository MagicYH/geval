package geval

// Gengin : todo
type Gengin interface {
	AddRule(ruleNode RuleNode, priority int) error
	AddData(name string, data interface{}) error
	AddFunc(name string, fun interface{}) error
	Eval() error
}

// BaseEngin : todo
type BaseEngine struct {
	dataCtx *DataContext
	funcCtx *FunContext
}
