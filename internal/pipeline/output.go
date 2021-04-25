package pipeline

type RawOutput map[string]interface{}

func (f RawOutput) GetName() string {
	return f["output"].(string)
}

func (f RawOutput) GetParams() map[string]interface{} {
	p, _ := f["params"].(map[string]interface{})
	return p
}
