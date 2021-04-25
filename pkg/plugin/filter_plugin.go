package plugin

import "github.com/Jeffail/gabs"

type FilterPlugin interface {
	Init(map[string]interface{}) error
	Exec(*gabs.Container) (bool, error)
}
