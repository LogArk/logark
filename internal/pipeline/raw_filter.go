package pipeline

type RawFilter map[string]interface{}

func (f RawFilter) GetName() string {
	return f["filter"].(string)
}

func (f RawFilter) GetOnFailure() []RawFilter {
	of, _ := f["on_failure"].([]interface{})
	result := make([]RawFilter, 0)
	for _, v := range of {
		newFilter, ok := v.(map[string]interface{})
		if ok {
			result = append(result, newFilter)
		}
	}
	return result
}

func (f RawFilter) GetOnSuccess() []RawFilter {
	of, _ := f["on_success"].([]interface{})
	result := make([]RawFilter, 0)
	for _, v := range of {
		newFilter, ok := v.(map[string]interface{})
		if ok {
			result = append(result, newFilter)
		}
	}
	return result
}

func (f RawFilter) GetParams() map[string]interface{} {
	p, _ := f["params"].(map[string]interface{})
	return p
}
