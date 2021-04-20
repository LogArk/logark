package mutate

import (
	"github.com/Jeffail/gabs"
)

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