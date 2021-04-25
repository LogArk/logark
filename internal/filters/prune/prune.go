package prune

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func allowList(event *gabs.Container, value interface{}) {
	keyMap := make(map[string]bool)
	for _, v := range value.([]interface{}) {
		strKey := v.(string)
		keyMap[strKey] = true
	}
	keys, _ := event.S().ChildrenMap()
	for key := range keys {
		if !keyMap[key] {
			event.Delete(key)
		}
	}

}

func denyList(event *gabs.Container, value interface{}) {
	keyMap := make(map[string]bool)
	for _, v := range value.([]interface{}) {
		strKey := v.(string)
		keyMap[strKey] = true
	}
	keys, _ := event.S().ChildrenMap()
	for key := range keys {
		if keyMap[key] {
			event.Delete(key)
		}
	}

}

func ExecFilter(event *gabs.Container, params map[string]interface{}) bool {
	status := true
	for key, value := range params {
		switch key {
		case "allow_list":
			allowList(event, value)
		case "deny_list":
			denyList(event, value)
		default:
			fmt.Printf("Unknown param: %s\n", key)
		}
	}
	//Println("--- exiting mutate.")
	return status
}
