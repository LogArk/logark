package main

import (
	"github.com/Jeffail/gabs"
	"github.com/LogArk/logark/pkg/plugin"
)

type Command struct {
	name  string
	key   string
	value interface{}
}

type MutatePlugin struct {
	commands []Command
}

func (p *MutatePlugin) Init(params map[string]interface{}) error {
	for key, value := range params {
		for _, v := range value.([]interface{}) {
			hash := v.(map[string]interface{})
			newCommand := Command{
				name:  key,
				key:   hash["key"].(string),
				value: hash["value"],
			}
			p.commands = append(p.commands, newCommand)
		}
	}
	return nil
}

func (p MutatePlugin) Exec(event *gabs.Container) (bool, error) {
	status := true

	for _, c := range p.commands {
		switch c.name {
		case "add_field":
			addField(event, c.key, c.value)
		case "update_field":
			updateField(event, c.key, c.value)
		default:
		}
	}
	return status, nil
}

func addField(event *gabs.Container, path string, fieldValue interface{}) {
	event.SetP(fieldValue, path)
}

func updateField(event *gabs.Container, path string, fieldValue interface{}) {
	data := event.Path(path).Data()
	if data != nil {
		event.SetP(fieldValue, path)
	}

}

func New() plugin.FilterPlugin {
	var p MutatePlugin
	return &p
}
