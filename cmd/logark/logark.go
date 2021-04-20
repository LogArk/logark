package main

import (
	"github.com/Jeffail/gabs"
	"fmt"
	"github.com/LogArk/logark/internal/filters/mutate"
	"sigs.k8s.io/yaml"
)


type PipelineSettings struct {
	Workers uint `json:"workers"`
}

type Filter map[string]interface{}

func (f Filter) getName() string {
	return f["filter"].(string)
}

func (f Filter) getOnFailure() []Filter {
	of,_ := f["on_failure"].([]interface{})
	result := make([]Filter,0)
	for _,v := range of {
		newFilter,ok := v.(map[string]interface{})
		if ok {
			result = append(result, newFilter)
		}
	}
	return result
}

func (f Filter) getOnSuccess() []Filter {
	of,_ := f["on_success"].([]interface{})
	result := make([]Filter,0)
	for _,v := range of {
		newFilter,ok := v.(map[string]interface{})
		if ok {
			result = append(result, newFilter)
		}
	}
	return result
}

func (f Filter) getParams() map[string]interface{} {
	p,_ := f["params"].(map[string]interface{})
	return p
}

type Process []Filter

type Pipeline struct {
	Settings PipelineSettings `json:"settings"`
	Process Process `json:"process"`
}



func main() {

	fmt.Println("Loading pipeline...")

	y := []byte(`settings:
  workers: 5
inputs:
- type: beats
  port: 5048

process:
- filter: json
  params:
    source: "message"
    skip_on_invalid_json: true
  on_success:
  - filter: mutate
    add_field: ["source.env","DENVER"]
- filter: prune
  params:
    allowlist: ["data","type","v","sec","source","timestamp"]
- filter: test
  params:
    condition: "[v]"
  on_failure:
  - filter: drop
- filter: mutate
  params:
    add_field:
    - key: "source.env"
      value: "DENVER"
    - key: "source.version"
      value: "1.0.0"
- filter: mutate
  params:
    update_field:
    - key: "source.env"
      value: "toto"

      


outputs:
- output: stdout
  codec: json
- output: kafka
  codec: json
  topic_id: "aircraft_hub"
`)  

	j, _ := yaml.YAMLToJSON(y)
	fmt.Printf("%s\n",j)

	var p Pipeline

	yaml.Unmarshal(y, &p)

	fmt.Println(p)

	fmt.Println("Parsing log...")
	// Input
	e := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)

	jsonParsed, _ := gabs.ParseJSON(e)

	var executionStack []Filter

	// Initialize filter stack 
	for i:=len(p.Process)-1; i>=0; i-- {
		executionStack = append(executionStack, p.Process[i])
	}


	for len(executionStack) > 0 {
		n := len(executionStack) - 1

		// Pop last element
		V := executionStack[n]
		executionStack = executionStack[:n] 

		fmt.Println("----",V.getName(),"----")

		status := false
		switch V.getName() {
		case "mutate":
			status = mutate.ExecFilter(jsonParsed, V.getParams())
		default:
			fmt.Println("Cannot handle",V.getName())
		}
		if status {
			fmt.Println("Operation was succcesful, looking at on_success")
			success := V.getOnSuccess()
			for i:=len(success)-1; i>=0; i-- {
				executionStack = append(executionStack, success[i])
			}
		} else {
			fmt.Println("Operation was not succcesful, looking at on_failure")
			failure := V.getOnFailure()
			for i:=len(failure)-1; i>=0; i-- {
				executionStack = append(executionStack, failure[i])
			}
		}

	}
	// Output
	//out, _:= json.Marshal(f)
	fmt.Println(jsonParsed.String())

}