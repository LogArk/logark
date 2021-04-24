package pipeline

type Filter map[string]interface{}

func (f Filter) GetName() string {
	return f["filter"].(string)
}

func (f Filter) GetOnFailure() []Filter {
	of, _ := f["on_failure"].([]interface{})
	result := make([]Filter, 0)
	for _, v := range of {
		newFilter, ok := v.(map[string]interface{})
		if ok {
			result = append(result, newFilter)
		}
	}
	return result
}

func (f Filter) GetOnSuccess() []Filter {
	of, _ := f["on_success"].([]interface{})
	result := make([]Filter, 0)
	for _, v := range of {
		newFilter, ok := v.(map[string]interface{})
		if ok {
			result = append(result, newFilter)
		}
	}
	return result
}

func (f Filter) GetParams() map[string]interface{} {
	p, _ := f["params"].(map[string]interface{})
	return p
}
