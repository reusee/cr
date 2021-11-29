package main

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Item struct {
	Pkg      *packages.Package
	FullName string
	Object   types.Object
}

type AllItems map[string]*Item

func (_ Global) AllItems(
	pkgs Pkgs,
	dir ModuleDir,
) AllItems {
	items := make(AllItems)

	modDir := string(dir)
	packages.Visit(pkgs, func(pkg *packages.Package) bool {
		if pkg.Module == nil {
			return false
		}
		if !strings.HasPrefix(pkg.Module.Dir, modDir) {
			return false
		}

		var walkScope func(scope types.Scope, path []string)
		walkScope = func(scope types.Scope, path []string) {
			names := scope.Names()
			for _, name := range names {
				obj := scope.Lookup(name)
				item := &Item{
					Pkg:      pkg,
					FullName: pkg.PkgPath + "." + strings.Join(append(path, name), "."),
					Object:   obj,
				}
				items[item.FullName] = item
			}
		}

		walkScope(*pkg.Types.Scope(), nil)

		return true
	}, nil)

	return items
}
