package main

import (
	"os"
	"regexp"
	"sort"

	"golang.org/x/tools/go/packages"
)

func (_ Global) UsesCommand(
	pkgs Pkgs,
	fset Fset,
) Commands {
	return Commands{

		"uses": func(
			args Args,
			items AllItems,
		) {
			if len(args) != 1 {
				pt("need an object pattern\n")
				os.Exit(-1)
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
				os.Exit(-1)
			}

			if len(candidates) > 1 {
				sort.Slice(candidates, func(i, j int) bool {
					return candidates[i].FullName < candidates[j].FullName
				})
				pt("ambiguous pattern\n")
				for _, item := range candidates {
					pt("%s\n", item.FullName)
				}
				os.Exit(-1)
			}

			item := candidates[0]
			packages.Visit(pkgs, func(pkg *packages.Package) bool {
				for ident, obj := range pkg.TypesInfo.Uses {
					if obj != item.Object {
						continue
					}
					pt("%s\n", fset.Position(ident.Pos()))
				}
				return true
			}, nil)

		},
	}
}
