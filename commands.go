package main

import (
	"fmt"
	"reflect"

	"github.com/reusee/dscope"
)

type Commands map[string]any

var _ dscope.Reducer = Commands{}

func (_ Commands) Reduce(_ dscope.Scope, vs []reflect.Value) reflect.Value {
	ret := make(Commands)
	names := make(map[string]bool)
	for _, v := range vs {
		cmds := v.Interface().(Commands)
		for name, cmd := range cmds {
			if _, ok := names[name]; ok {
				panic(fmt.Errorf("duplicated command: %s", name))
			}
			ret[name] = cmd
		}
	}
	return reflect.ValueOf(ret)
}

func (_ Global) Commands() Commands {
	return nil
}
