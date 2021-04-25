package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	plg "github.com/LogArk/logark/pkg/plugin"
)

type FilterAction struct {
	Name      string
	Plugin    plg.FilterPlugin
	OnSuccess *FilterPipeline
	OnFailure *FilterPipeline
}

type FilterPipeline struct {
	Filters []FilterAction
}

func buildPipeline(f []RawFilter) *FilterPipeline {
	var err error
	var rootPipeline FilterPipeline
	var plugins map[string]*plugin.Plugin = make(map[string]*plugin.Plugin)

	/* Get binary directory */
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range f {
		var fa FilterAction
		fa.Name = v.GetName()

		// Check if plugin is already loaded. If not, load it
		plug := plugins[fa.Name]
		if plug == nil {
			plug, err = plugin.Open(dir + "/plugins/" + fa.Name + ".so")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[fa.Name] = plug
		}

		newPlugin, err := plug.Lookup("New")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fa.Plugin = newPlugin.(func() plg.FilterPlugin)()
		fa.Plugin.Init(v.GetParams())
		fa.OnSuccess = buildPipeline(v.GetOnSuccess())
		fa.OnFailure = buildPipeline(v.GetOnFailure())

		rootPipeline.Filters = append(rootPipeline.Filters, fa)
	}
	return &rootPipeline
}
