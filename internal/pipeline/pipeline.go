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
	Outputs  []Output         `json:"Outputs"`
}

type Pipeline struct {
	raw RawPipeline

	Settings       *PipelineSettings
	FilterPipeline *FilterPipeline
	Outputs        []Output
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
	p.Outputs = p.raw.Outputs
	p.FilterPipeline = buildPipeline(p.raw.Filters)

	return p, nil
}
