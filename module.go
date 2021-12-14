package main

import (
	"os"
	"path/filepath"
)

type ModuleDir string

func (_ Global) ModuleDir() ModuleDir {
	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if dir == "" {
		dir = "."
	}
	dir, err := filepath.Abs(dir)
	ce(err)
	// find go.mod
	for {
		_, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil {
			return ModuleDir(dir)
		}
		nextDir := filepath.Dir(dir)
		if nextDir == dir {
			panic("go.mod not found")
		}
		dir = nextDir
	}
}
