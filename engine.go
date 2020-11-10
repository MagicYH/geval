package geval

// Gengine : todo
type Gengine interface {
	AddRule(ruleNode RuleNode, priority int) error
	AddData(name string, data interface{}) error
	AddFunc(name string, fun interface{}) error
	Eval() error
}

// BaseEngine : todo
type BaseEngine struct {
	dataCtx *DataContext
	funcCtx *FunContext
}
