package main

import (
	"fmt"

	"github.com/LogArk/logark/pkg/plugin"
)

type StdoutPlugin struct {
}

func (sp StdoutPlugin) Send(log []byte) error {
	fmt.Println(string(log))
	return nil
}

func (f *StdoutPlugin) Init(params map[string]interface{}) error {
	return nil
}

func New() plugin.OutputPlugin {
	var p StdoutPlugin

	return &p
}
