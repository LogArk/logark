package pipeline

type Output map[string]interface{}

func (f Output) GetName() string {
	return f["output"].(string)
}

func (f Output) GetParams() map[string]interface{} {
	p, _ := f["params"].(map[string]interface{})
	return p
}
