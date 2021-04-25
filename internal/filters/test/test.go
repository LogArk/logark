package test

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/Knetic/govaluate"
)

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

func ExecFilter(event *gabs.Container, params map[string]interface{}) bool {
	status := true

	condition := params["condition"].(string)
	myLog := EventParameter{Event: event}

	/* Build expression */
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(condition, functions)
	if err != nil {
		fmt.Println(err)
		return false
	}

	/* Evaluate result */
	result, err := expression.Eval(&myLog)
	if err != nil {
		fmt.Println(err)
		return false
	}

	/* If result is boolean, take result. Otherwise, evaluate if nil or not */
	switch v := result.(type) {
	case bool:
		status = v
	default:
		status = v != nil
	}

	return status
}
