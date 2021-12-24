// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// This program is run via "go generate" (via a directive in sort.go)
// to generate zfunc_a.go & zfunc_b.go.

package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
)

var hackedFuncs = make(map[string]bool)

func main() {
	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, "sort.go", nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	af.Doc = nil
	af.Imports = nil
	af.Comments = nil

	var newDecl []ast.Decl
	for _, d := range af.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Recv != nil ||
			fd.Name.Name == "less" || fd.Name.IsExported() ||
			fd.Type.TypeParams == nil || len(fd.Type.TypeParams.List) != 1 {
			continue
		}
		field := fd.Type.TypeParams.List[0]
		if expr, ok := field.Type.(*ast.SelectorExpr); !ok ||
			expr.Sel.Name != "Ordered" ||
			len(field.Names) != 1 || field.Names[0].Name != "E" {
			continue
		}
		hackedFuncs[fd.Name.Name] = true
		fd.Type.TypeParams = nil
		newDecl = append(newDecl, fd)
	}
	af.Decls = newDecl
	ast.Walk(visitFunc(rewriteCalls), af)

	var out bytes.Buffer
	if err := format.Node(&out, fset, af); err != nil {
		log.Fatalf("format.Node: %v", err)
	}
	tpl := out.Bytes()

	funcPtn := regexp.MustCompile(`\nfunc `)
	lessPtn := regexp.MustCompile(`less\([^\),]+,[^\),]+\)`)

	src := funcPtn.ReplaceAll(tpl, []byte("\nfunc (lt lessFunc[E]) "))
	src = lessPtn.ReplaceAllFunc(src, func(origin []byte) []byte {
		out := make([]byte, len(origin)-2)
		out[0] = 'l'
		out[1] = 't'
		for i := 2; i < len(out); i++ {
			out[i] = origin[i+2]
		}
		return out
	})
	dumpOrDie("zfunc_a.go", src)

	src = funcPtn.ReplaceAll(tpl, []byte("\nfunc (lt refLessFunc[E]) "))
	src = lessPtn.ReplaceAllFunc(src, func(origin []byte) []byte {
		out := make([]byte, len(origin))
		out[0] = 'l'
		out[1] = 't'
		out[2] = '('
		out[3] = '&'
		pos := bytes.IndexByte(origin, ',')
		for i := 4; i < pos; i++ {
			out[i] = origin[i+1]
		}
		out[pos] = '&'
		for i := pos + 1; i < len(out); i++ {
			out[i] = origin[i]
		}
		return out
	})
	dumpOrDie("zfunc_b.go", src)
}

type visitFunc func(ast.Node) ast.Visitor

func (f visitFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

func rewriteCalls(n ast.Node) ast.Visitor {
	ce, ok := n.(*ast.CallExpr)
	if ok {
		ident, ok := ce.Fun.(*ast.Ident)
		if ok && hackedFuncs[ident.Name] {
			ident.Name = "lt." + ident.Name
		}
	}
	return visitFunc(rewriteCalls)
}

var header = `// Code generated from sort.go using genzfunc.go; DO NOT EDIT.

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

`

func dumpOrDie(filename string, src []byte) {
	src, err := format.Source(src)
	if err != nil {
		log.Fatalf("format.Source: %v on\n%s", err, src)
	}
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if _, err := out.WriteString(header); err != nil {
		log.Fatal(err)
	}
	if _, err := out.Write(src); err != nil {
		log.Fatal(err)
	}
}
