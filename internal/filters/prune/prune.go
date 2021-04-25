package prune

import (
	"errors"

	"github.com/Jeffail/gabs"
	"github.com/LogArk/logark/pkg/plugin"
)

type PruneFilter struct {
	pruneType string
	keyMap    map[string]bool
}

func (p *PruneFilter) Init(params map[string]interface{}) error {
	p.keyMap = make(map[string]bool)
	// XXX kinda ugly way to get the first item, i know
	for key, value := range params {
		switch key {
		case "allow_list", "deny_list":
			p.pruneType = key
			keys := value.([]interface{})
			for _, k := range keys {
				p.keyMap[k.(string)] = true
			}
			return nil
		default:
			return errors.New("Prune: unexpected keyword: " + key)
		}
	}
	return nil
}

func (p PruneFilter) Exec(event *gabs.Container) (bool, error) {
	status := true

	switch p.pruneType {
	case "allow_list":
		p.allowList(event, p.keyMap)
	case "deny_list":
		p.denyList(event, p.keyMap)
	}

	return status, nil
}

func (p PruneFilter) allowList(event *gabs.Container, value map[string]bool) {
	keys, _ := event.S().ChildrenMap()
	for key := range keys {
		if !p.keyMap[key] {
			event.Delete(key)
		}
	}
}

func (p PruneFilter) denyList(event *gabs.Container, value interface{}) {
	keys, _ := event.S().ChildrenMap()
	for key := range keys {
		if p.keyMap[key] {
			event.Delete(key)
		}
	}
}

func New() plugin.FilterPlugin {
	var p PruneFilter

	return &p
}
