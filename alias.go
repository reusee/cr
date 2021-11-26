package main

import (
	"errors"
	"fmt"

	"github.com/reusee/dscope"
	"github.com/reusee/e4"
)

var (
	ce = e4.Check
	pt = fmt.Printf
	is = errors.Is
)

type (
	any   = interface{}
	Scope = dscope.Scope
)
