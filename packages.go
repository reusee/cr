package main

import (
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Pkgs = []*packages.Package

type Fset = *token.FileSet

func (_ Global) Pkgs(
	dir ModuleDir,
	loadConfig LoadConfig,
	ignorePatterns IgnorePatterns,
) (
	pkgs Pkgs,
	fset Fset,
) {

	loadConfig()

	var dirs []string
	ce(filepath.WalkDir(string(dir), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		for _, pattern := range *ignorePatterns {
			if pattern.MatchString(path) {
				return nil
			}
		}
		f, err := os.Open(path)
		ce(err)
		defer f.Close()
		names, err := f.Readdirnames(-1)
		ce(err)
		hasGoFile := false
		for _, name := range names {
			if strings.HasSuffix(name, ".go") {
				hasGoFile = true
				break
			}
		}
		if !hasGoFile {
			return nil
		}
		dirs = append(dirs, path)
		return nil
	}))

	fset = new(token.FileSet)
	var err error
	pkgs, err = packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedCompiledGoFiles |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedModule,
		Dir: string(dir),
		Env: append(os.Environ(), []string{
			"GOARCH=amd64",
		}...),
		Fset:  fset,
		Tests: true,
	}, dirs...)
	ce(err)
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(-1)
	}

	return
}

type IgnorePatterns *[]*regexp.Regexp

func (_ Global) IgnorePatterns() IgnorePatterns {
	var patterns []*regexp.Regexp
	return &patterns
}

func (_ Global) PackagesConfigFuncs(
	patterns IgnorePatterns,
) ConfigFuncs {
	return ConfigFuncs{
		"ignore": func(pattern string) {
			*patterns = append(*patterns, regexp.MustCompile(pattern))
		},
	}
}
