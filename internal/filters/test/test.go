package test

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/Knetic/govaluate"
	"github.com/LogArk/logark/pkg/plugin"
)

type TestFilter struct {
	rawCondition string
	expression   *govaluate.EvaluableExpression
}

type EventParameter struct {
	Event *gabs.Container
}

func (ep EventParameter) Get(name string) (interface{}, error) {
	result := ep.Event.Path(name).Data()
	return result, nil
}

var functions = map[string]govaluate.ExpressionFunction{
	"strlen": func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return (float64)(length), nil
	},
	"exists": func(args ...interface{}) (interface{}, error) {
		return len(args) != 0, nil
	},
}

func (f *TestFilter) Init(params map[string]interface{}) error {
	var err error

	f.rawCondition = params["condition"].(string)

	/* Build expression */
	f.expression, err = govaluate.NewEvaluableExpressionWithFunctions(f.rawCondition, functions)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (f TestFilter) Exec(event *gabs.Container) (bool, error) {
	status := true

	myLog := EventParameter{Event: event}

	/* Evaluate result */
	result, err := f.expression.Eval(&myLog)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	/* If result is boolean, take result. Otherwise, evaluate if nil or not */
	switch v := result.(type) {
	case bool:
		status = v
	default:
		status = v != nil
	}

	return status, nil
}

func New() plugin.FilterPlugin {
	var p TestFilter

	return &p
}
