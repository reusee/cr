package main

import (
	"os"
	"path/filepath"
)

type ModuleDir string

func (_ Global) ModuleDir() ModuleDir {
	// find go.mod
	dir, err := filepath.Abs(".")
	ce(err)
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
