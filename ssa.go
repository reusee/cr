package main

import (
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type SSAProgram = *ssa.Program

type SSAPkgs = []*ssa.Package

func (_ Global) SSA(
	ps Pkgs,
) (
	program SSAProgram,
	pkgs SSAPkgs,
) {
	program, pkgs = ssautil.AllPackages(ps, ssa.PrintPackages)
	program.Build()
	return
}

func (_ Global) SSACommand() Commands {
	return Commands{
		"ssa": func(
			ssaProgram SSAProgram,
		) {
		},
	}
}
