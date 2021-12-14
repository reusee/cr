package main

import "testing"

func TestParser(t *testing.T) {
	type Spec struct {
		Input  string
		Tokens []string
	}

	specs := []Spec{
		{
			`  foo   bar  baz    `,
			[]string{
				"foo", "bar", "baz",
			},
		},
		{
			"`foo` 'bar' \"baz\"",
			[]string{
				"foo", "bar", "baz",
			},
		},
	}

	for _, spec := range specs {
		var res [][]rune
		p := ParseTokens(&res, nil)
		var err error
		for _, r := range spec.Input {
			p, err = p(r)
			ce(err)
		}
		if len(res) != len(spec.Tokens) {
			t.Fatal()
		}
		for i, item := range res {
			s := string(item)
			if s != spec.Tokens[i] {
				t.Fatal()
			}
		}
	}

}
