package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/reusee/dscope"
	"github.com/reusee/starlarkutil"
	"go.starlark.net/starlark"
)

type LoadConfig func()

func (_ Global) LoadConfig(
	dir ModuleDir,
	scope Scope,
) LoadConfig {
	return func() {
		content, err := os.ReadFile(filepath.Join(string(dir), "cr.py"))
		if is(err, os.ErrNotExist) {
			return
		}
		var fns ConfigFuncs
		scope.Assign(&fns)
		pyFuncs := make(starlark.StringDict)
		for name, fn := range fns {
			pyFuncs[name] = starlarkutil.MakeFunc(name, fn)
		}
		_, err = starlark.ExecFile(
			new(starlark.Thread),
			"cr.py",
			content,
			pyFuncs,
		)
		ce(err)
	}
}

type ConfigFuncs map[string]any

func (_ Global) BuiltinConfigFuncs() ConfigFuncs {
	return ConfigFuncs{
		"pt": func(format string, args ...any) {
			fmt.Printf(format, args...)
		},
	}
}

var _ dscope.Reducer = ConfigFuncs{}

func (_ ConfigFuncs) Reduce(_ Scope, vs []reflect.Value) reflect.Value {
	return dscope.Reduce(vs)
}
