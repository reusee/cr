package main

import (
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/reusee/dscope"
	"golang.org/x/tools/go/packages"
)

type Pkgs = []*packages.Package

type Fset = *token.FileSet

func (_ Global) Pkgs(
	dir ModuleDir,
	ignorePatterns IgnorePatterns,
	buildTags BuildTags,
) (
	pkgs Pkgs,
	fset Fset,
) {

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
		hasGoModFile := false
		for _, name := range names {
			if strings.HasSuffix(name, ".go") {
				hasGoFile = true
			}
			if name == "go.mod" {
				hasGoModFile = true
			}
		}
		if !hasGoFile {
			return nil
		}
		if hasGoModFile {
			return nil
		}
		dirs = append(dirs, path)
		return nil
	}))

	var flags []string
	if len(buildTags) > 0 {
		flags = append(flags, "-tags="+strings.Join(buildTags, ","))
	}
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
		Fset:       fset,
		Tests:      true,
		BuildFlags: flags,
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

func (g Global) ReloadCommand() Commands {
	return Commands{
		"reload": func(
			mutate dscope.Mutate,
		) {
			mutate(g.Pkgs)
		},
	}
}

type BuildTags []string

func (_ Global) BuildTags() BuildTags {
	return nil
}

func (_ Global) BuildTagsFunc(
	mutate dscope.Mutate,
) ConfigFuncs {
	return ConfigFuncs{
		"tags": func(tags BuildTags) {
			mutate(&tags)
		},
	}
}
