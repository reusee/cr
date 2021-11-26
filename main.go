package main

import (
	"fmt"
	"os"

	"github.com/reusee/dscope"
)

type Global struct{}

type Args []string

func main() {
	decls := dscope.Methods(new(Global))
	scope := dscope.New(decls...)
	scope.Call(func(
		cmds Commands,
		scope Scope,
	) {
		if len(os.Args) == 1 {
			for name := range cmds {
				pt("%s\n", name)
			}
			return
		}
		name := os.Args[1]
		fn, ok := cmds[name]
		if !ok {
			panic(fmt.Errorf("no such command: %s", name))
		}
		scope.Fork(func() Args {
			return os.Args[2:]
		}).Call(fn)
	})
}
