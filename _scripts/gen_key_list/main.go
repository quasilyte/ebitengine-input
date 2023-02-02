//go:build generate

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
)

type keyInfo struct {
	VarName string
	KeyName string
}

func main() {
	outputFile := flag.String("o", "", "output file name")
	flag.Parse()
	if len(flag.Args()) != 1 {
		panic("expected example 1 positional argument: keys file path")
	}
	keysFile := flag.Args()[0]

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, keysFile, nil, 0)
	if err != nil {
		panic(err)
	}

	var keyVariables []keyInfo
	for _, decl := range astFile.Decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR {
			continue
		}
		for _, spec := range decl.Specs {
			spec := spec.(*ast.ValueSpec)
			if len(spec.Names) != 1 || len(spec.Values) != 1 {
				continue
			}
			id := spec.Names[0]
			if !strings.HasPrefix(id.Name, "Key") {
				continue
			}
			lit, ok := spec.Values[0].(*ast.CompositeLit)
			if !ok {
				continue
			}
			nameValue := getCompositeLitField(lit, "name")
			if nameValue == nil {
				panic(fmt.Sprintf("missing name field init value for %s", id))
			}
			nameString, ok := nameValue.(*ast.BasicLit)
			if !ok || nameString.Kind != token.STRING {
				panic(fmt.Sprintf("unexpected name field init for %s", id))
			}
			keyName, err := strconv.Unquote(nameString.Value)
			if err != nil {
				panic(err) // should never happen
			}
			if keyName == "" {
				panic(fmt.Sprintf("empty key name for %s", id))
			}
			keyVariables = append(keyVariables, keyInfo{
				VarName: id.Name,
				KeyName: keyName,
			})
		}
	}

	sort.SliceStable(keyVariables, func(i, j int) bool {
		return keyVariables[i].KeyName < keyVariables[j].KeyName
	})

	templateData := map[string]interface{}{
		"Keys": keyVariables,
	}

	var buf bytes.Buffer
	if err := outputTemplate.Execute(&buf, templateData); err != nil {
		panic(err)
	}

	prettyOutput, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if *outputFile == "" {
		fmt.Print(string(prettyOutput))
	} else {
		if err := os.WriteFile(*outputFile, prettyOutput, 0o0664); err != nil {
			panic(err)
		}
	}
}

func getCompositeLitField(lit *ast.CompositeLit, fieldName string) ast.Expr {
	for _, elem := range lit.Elts {
		kv, ok := elem.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		id, ok := kv.Key.(*ast.Ident)
		if ok && id.Name == fieldName {
			return kv.Value
		}
	}
	return nil
}

var outputTemplate = template.Must(template.New("output").Parse(`// Code generated by "_scripts/gen_key_list"; DO NOT EDIT.

package input

// allKeys contains all basic keys provided by a package.
// This slice is sorted by key names (see ActionKeyNames).
var allKeys = []Key{
{{- range $.Keys }}
    {{.VarName}},
{{- end }}
}
`))
