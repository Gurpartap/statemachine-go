package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"

	"github.com/Gurpartap/statemachine-go"
)

func main() {
	b, err := ioutil.ReadFile("examples/hcl/process.hcl")
	if err != nil {
		panic(err)
	}

	def := &statemachine.MachineDef{}
	err = hcl.Decode(def, string(b))
	if err != nil {
		panic(err)
	}

	defJSON, _ := json.MarshalIndent(def, "", "  ")
	fmt.Printf("%s\n", defJSON)
}
