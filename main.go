package main

import (
	"io"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/reusee/dscope"
)

type Global struct{}

type Args []string

func main() {
	decls := dscope.Methods(new(Global))
	scope := dscope.NewMutable(decls...)

	scope.Call(func(
		loadConfig LoadConfig,
	) {
		loadConfig()
	})

	scope.Call(func(
		cmds Commands,
	) {

		rl, err := readline.New("> ")
		ce(err)
		defer rl.Close()

		for {
			line, err := rl.Readline()
			if is(err, io.EOF) || is(err, readline.ErrInterrupt) {
				break
			}
			ce(err)

			var res [][]rune
			p := ParseTokens(&res, nil)
			for _, r := range []rune(line) {
				p, err = p(r)
				ce(err)
			}

			if len(res) == 0 {
				continue
			}

			name := string(res[0])
			fn, ok := cmds[name]
			if !ok {
				color.Red("no such command")
				continue
			}
			var args Args
			for _, rs := range res[1:] {
				args = append(args, string(rs))
			}
			scope.Fork(&args).Call(fn)

		}

	})
}
