package main

import (
	"regexp"
	"sort"

	"golang.org/x/tools/go/packages"
)

func (_ Global) UsesCommand(
	pkgs Pkgs,
	fset Fset,
) Commands {

	fn := func(
		args Args,
		items AllItems,
	) {
		if len(args) != 1 {
			pt("need an object pattern\n")
			return
		}

		pattern := regexp.MustCompile("(?i)" + args[0])
		var candidates []*Item
		for name, item := range items {
			if !pattern.MatchString(name) {
				continue
			}
			candidates = append(candidates, item)
		}

		if len(candidates) == 0 {
			pt("no such item\n")
			return
		}

		if len(candidates) > 1 {
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].FullName < candidates[j].FullName
			})
			pt("ambiguous pattern\n")
			for _, item := range candidates {
				pt("~ %s\n", item.FullName)
			}
			return
		}

		item := candidates[0]
		pt("%s\n", item.FullName)
		packages.Visit(pkgs, func(pkg *packages.Package) bool {
			for ident, obj := range pkg.TypesInfo.Uses {
				if obj != item.Object {
					continue
				}
				pt("%s\n", fset.Position(ident.Pos()))
			}
			return true
		}, nil)

	}

	return Commands{
		"uses": fn,
		"u":    fn,
	}
}
