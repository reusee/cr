package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/reusee/dscope"
)

type Global struct{}

type Args []string

func main() {
	decls := dscope.Methods(new(Global))
	scope := dscope.NewMutable(decls...)
	scope.Call(func(
		cmds Commands,
		get dscope.GetScope,
	) {

		r := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			os.Stdout.Sync()
			line, err := r.ReadString('\n')
			if is(err, io.EOF) {
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
				pt("no such command\n")
				continue
			}
			var args Args
			for _, rs := range res[1:] {
				args = append(args, string(rs))
			}
			get().Fork(&args).Call(fn)

		}

	})
}
