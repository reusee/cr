package main

import (
	"go/types"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

type TypeInfo struct {
	Package *packages.Package
}

type TypeInfoMap struct {
	*typeutil.Map
}

func (_ Global) TypeInfoMap(
	pkgs Pkgs,
) TypeInfoMap {
	m := TypeInfoMap{
		Map: new(typeutil.Map),
	}
	packages.Visit(pkgs, func(pkg *packages.Package) bool {
		for _, obj := range pkg.TypesInfo.Defs {
			typeName, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}
			m.Set(typeName.Type(), TypeInfo{
				Package: pkg,
			})
		}
		return true
	}, nil)
	return m
}
