package pipeline

import (
	"fmt"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

type PipelineSettings struct {
	Workers uint `json:"workers"`
}

type RawPipeline struct {
	Settings PipelineSettings `json:"settings"`
	Filters  []RawFilter      `json:"filters"`
	Outputs  []RawOutput      `json:"outputs"`
}

type Pipeline struct {
	raw RawPipeline

	Settings       *PipelineSettings
	FilterPipeline *FilterPipeline
	OutputPipeline *OutputPipeline
}

func Load(path string) (Pipeline, error) {
	var rp RawPipeline
	var p Pipeline

	yamlFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return p, err
	}
	defer yamlFile.Close()

	byteValue, _ := ioutil.ReadAll(yamlFile)
	j, _ := yaml.YAMLToJSON(byteValue)

	yaml.Unmarshal(j, &rp)

	p.raw = rp
	p.Settings = &p.raw.Settings
	p.FilterPipeline = buildFilterdPipeline(p.raw.Filters)
	p.OutputPipeline = buildOutputdPipeline(p.raw.Outputs)

	return p, nil
}
