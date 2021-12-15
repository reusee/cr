package main

import (
	"fmt"
	"go/types"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/emicklei/dot"
	"golang.org/x/tools/go/packages"
)

func (_ Global) TypeGraphCommand() Commands {

	fn := func(
		pkgs Pkgs,
		typeInfoMap TypeInfoMap,
		moduleDir ModuleDir,
		args Args,
	) {

		type FieldUse struct {
			Type string
			Use  string
		}
		fieldUses := make(map[FieldUse]bool)

		var pkgPattern *regexp.Regexp
		if len(args) > 0 {
			pkgPattern = regexp.MustCompile(args[0])
		} else {
			pkgPattern = regexp.MustCompile(`.*`)
		}

		modDir := string(moduleDir)
		set := func(t types.Type, use types.Type) {
			if _, ok := use.(*types.Named); !ok {
				return
			}
			v := typeInfoMap.At(use)
			if v == nil {
				return
			}
			useInfo := v.(TypeInfo)
			if useInfo.Package.Module == nil {
				return
			}
			if useInfo.Package.Module.Dir != modDir {
				return
			}
			info := typeInfoMap.At(t).(TypeInfo)
			if !pkgPattern.MatchString(info.Package.PkgPath) {
				return
			}
			if !pkgPattern.MatchString(useInfo.Package.PkgPath) {
				return
			}
			pt("%v %v\n", t, use)
			name := strings.TrimPrefix(t.String(), info.Package.Module.Path)
			useName := strings.TrimPrefix(use.String(), useInfo.Package.Module.Path)
			fieldUses[FieldUse{
				Type: name,
				Use:  useName,
			}] = true
		}

		var mark func(t types.Type)
		mark = func(typ types.Type) {
			switch underlying := typ.Underlying().(type) {

			case *types.Struct:
				pt("%v\n", typ)
				for i := 0; i < underlying.NumFields(); i++ {
					field := underlying.Field(i)
					fieldType := field.Type()
					set(typ, fieldType)
				}

			case *types.Slice:
				set(typ, underlying.Elem())

			case *types.Array:
				set(typ, underlying.Elem())

			case *types.Signature:
				//TODO

			case *types.Map:
				set(typ, underlying.Key())
				set(typ, underlying.Elem())

			case *types.Chan:
				set(typ, underlying.Elem())

			case *types.Basic:
			case *types.Interface:
			case *types.Pointer:

			default:
				panic(fmt.Errorf("unknown type: %T", underlying))
			}
		}

		packages.Visit(pkgs, func(pkg *packages.Package) bool {
			for _, name := range pkg.Types.Scope().Names() {
				obj := pkg.Types.Scope().Lookup(name)
				typeName, ok := obj.(*types.TypeName)
				if !ok {
					continue
				}
				mark(typeName.Type())
			}
			return true
		}, nil)

		graph := dot.NewGraph(dot.Directed)
		graph.NodeInitializer(func(n dot.Node) {
			n.Attr("shape", "rectangle")
		})

		nodes := make(map[string]dot.Node)
		getNode := func(name string) dot.Node {
			if node, ok := nodes[name]; ok {
				return node
			}
			node := graph.Node(name)
			nodes[name] = node
			return node
		}

		for use := range fieldUses {
			n1 := getNode(use.Type)
			n2 := getNode(use.Use)
			graph.Edge(n1, n2, "use")
		}

		cmd := exec.Command("dot", "-Tsvg", "-o", "type-graph.svg")
		cmd.Stdin = strings.NewReader(graph.String())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		ce(cmd.Run())

	}

	return Commands{
		"type-graph": fn,
		"tg":         fn,
	}
}
