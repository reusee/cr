package main

import (
	"go/token"
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

		if len(args) == 0 {
			args = append(args, ".*")
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
		posSet := make(map[token.Position]bool)
		packages.Visit(pkgs, func(pkg *packages.Package) bool {
			for ident, obj := range pkg.TypesInfo.Uses {
				if obj != item.Object {
					continue
				}
				posSet[fset.Position(ident.Pos())] = true
			}
			return true
		}, nil)
		var poses []token.Position
		for pos := range posSet {
			poses = append(poses, pos)
		}
		sort.Slice(poses, func(i, j int) bool {
			posA := poses[i]
			posB := poses[j]
			if posA.Filename != posB.Filename {
				return posA.Filename < posB.Filename
			}
			return posA.Offset < posB.Offset
		})
		for _, pos := range poses {
			pt("%s\n", pos)
		}

	}

	return Commands{
		"uses": fn,
		"u":    fn,
	}
}
