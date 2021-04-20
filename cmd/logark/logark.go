package main

import (
	"github.com/Jeffail/gabs"
	"fmt"
	"github.com/LogArk/logark/internal/filters/mutate"
)


type Node struct {


}


func main() {

	// Input
	e := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)

	jsonParsed, _ := gabs.ParseJSON(e)

	// Filter
	mutate.AddField(jsonParsed, "toto.new_field","value")
	mutate.AddField(jsonParsed, "Name","blah")
	mutate.UpdateField(jsonParsed, "Name", "new_blah")
	mutate.UpdateField(jsonParsed, "toto", "new_blah")
	// Output
	//out, _:= json.Marshal(f)
	fmt.Println(jsonParsed.String())

}