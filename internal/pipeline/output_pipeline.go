package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	plg "github.com/LogArk/logark/pkg/plugin"
)

type OutputAction struct {
	Name   string
	Plugin plg.OutputPlugin
}

type OutputPipeline struct {
	Outputs []OutputAction
}

func buildOutputdPipeline(f []RawOutput) *OutputPipeline {
	var err error
	var rootPipeline OutputPipeline
	var plugins map[string]*plugin.Plugin = make(map[string]*plugin.Plugin)

	/* Get binary directory */
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range f {
		var oa OutputAction
		oa.Name = v.GetName()

		// Check if plugin is already loaded. If not, load it
		plug := plugins[oa.Name]
		if plug == nil {
			plug, err = plugin.Open(dir + "/plugins/outputs/" + oa.Name + ".so")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[oa.Name] = plug
		}

		newPlugin, err := plug.Lookup("New")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		oa.Plugin = newPlugin.(func() plg.OutputPlugin)()
		oa.Plugin.Init(v.GetParams())

		rootPipeline.Outputs = append(rootPipeline.Outputs, oa)
	}
	return &rootPipeline
}
