package main

import (
	"testing"

	"github.com/reusee/sb"
)

func eq(t *testing.T, args ...any) {
	t.Helper()
	for i := 0; i < len(args); i += 2 {
		res := sb.MustCompare(
			sb.Marshal(args[i]),
			sb.Marshal(args[i+1]),
		)
		if res != 0 {
			t.Fatalf("not equal\n%#v\n%#v\n", args[i], args[i+1])
		}
	}
}
