package mutate

import (
	"github.com/Jeffail/gabs"
)

func ExecFilter(event *gabs.Container, params map[string]interface{}) bool {
	status := true
	//Println("--- entering mutate.")
	//Println(params)
	for key, value := range params {
		switch key {
		case "add_field":
			for _, v := range value.([]interface{}) {
				hash := v.(map[string]interface{})
				AddField(event, hash["key"].(string), hash["value"])
			}

		case "update_field":
			for _, v := range value.([]interface{}) {
				hash := v.(map[string]interface{})
				UpdateField(event, hash["key"].(string), hash["value"])
			}
		default:
			//Println("Unknown param: ", key)
		}
	}

	//Println("--- exiting mutate.")
	return status
}

func getFieldFromName(event *interface{}, fieldName string) *interface{} {
	return nil
}

func AddField(event *gabs.Container, path string, fieldValue interface{}) {
	event.SetP(fieldValue, path)
}

func UpdateField(event *gabs.Container, path string, fieldValue interface{}) {
	data := event.Path(path).Data()
	if data != nil {
		event.SetP(fieldValue, path)
	}

}
