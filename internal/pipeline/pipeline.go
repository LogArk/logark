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

type Pipeline struct {
	Settings PipelineSettings `json:"settings"`
	Process  []Filter         `json:"process"`
	Outputs  []Output         `json:"Outputs"`
}

func Load(path string) (Pipeline, error) {
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

	yaml.Unmarshal(j, &p)

	return p, nil
}
