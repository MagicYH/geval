package geval

type Gengin interface {
	AddRule(ruleNode RuleNode, priority int) error
	AddData(name string, data interface{}) error
	AddFunc(name string, fun interface{}) error
	Eval() error
}

type BaseEngine struct {
	dataCtx *DataContext
	funcCtx *FunContext
}
