// Do not edit. Bootstrap copy of /tmp/go-20170529-7428-4cl5ea/go/src/cmd/compile/internal/gc/mkbuiltin.go

//line /tmp/go-20170529-7428-4cl5ea/go/src/cmd/compile/internal/gc/mkbuiltin.go:1
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Generate builtin.go from builtin/runtime.go.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var stdout = flag.Bool("stdout", false, "write to stdout instead of builtin.go")

func main() {
	flag.Parse()

	var b bytes.Buffer
	fmt.Fprintln(&b, "// AUTO-GENERATED by mkbuiltin.go; DO NOT EDIT")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "package gc")

	mkbuiltin(&b, "runtime")

	out, err := format.Source(b.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if *stdout {
		_, err = os.Stdout.Write(out)
	} else {
		err = ioutil.WriteFile("builtin.go", out, 0666)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func mkbuiltin(w io.Writer, name string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join("builtin", name+".go"), nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	var interner typeInterner

	fmt.Fprintf(w, "var %sDecls = [...]struct { name string; tag int; typ int }{\n", name)
	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			if decl.Recv != nil {
				log.Fatal("methods unsupported")
			}
			if decl.Body != nil {
				log.Fatal("unexpected function body")
			}
			fmt.Fprintf(w, "{%q, funcTag, %d},\n", decl.Name.Name, interner.intern(decl.Type))
		case *ast.GenDecl:
			if decl.Tok != token.VAR {
				log.Fatal("unhandled declaration kind", decl.Tok)
			}
			for _, spec := range decl.Specs {
				spec := spec.(*ast.ValueSpec)
				if len(spec.Values) != 0 {
					log.Fatal("unexpected values")
				}
				typ := interner.intern(spec.Type)
				for _, name := range spec.Names {
					fmt.Fprintf(w, "{%q, varTag, %d},\n", name.Name, typ)
				}
			}
		default:
			log.Fatal("unhandled decl type", decl)
		}
	}
	fmt.Fprintln(w, "}")

	fmt.Fprintln(w)
	fmt.Fprintf(w, "func %sTypes() []*Type {\n", name)
	fmt.Fprintf(w, "var typs [%d]*Type\n", len(interner.typs))
	for i, typ := range interner.typs {
		fmt.Fprintf(w, "typs[%d] = %s\n", i, typ)
	}
	fmt.Fprintln(w, "return typs[:]")
	fmt.Fprintln(w, "}")
}

// typeInterner maps Go type expressions to compiler code that
// constructs the denoted type. It recognizes and reuses common
// subtype expressions.
type typeInterner struct {
	typs []string
	hash map[string]int
}

func (i *typeInterner) intern(t ast.Expr) int {
	x := i.mktype(t)
	v, ok := i.hash[x]
	if !ok {
		v = len(i.typs)
		if i.hash == nil {
			i.hash = make(map[string]int)
		}
		i.hash[x] = v
		i.typs = append(i.typs, x)
	}
	return v
}

func (i *typeInterner) subtype(t ast.Expr) string {
	return fmt.Sprintf("typs[%d]", i.intern(t))
}

func (i *typeInterner) mktype(t ast.Expr) string {
	switch t := t.(type) {
	case *ast.Ident:
		switch t.Name {
		case "byte":
			return "bytetype"
		case "rune":
			return "runetype"
		}
		return fmt.Sprintf("Types[T%s]", strings.ToUpper(t.Name))

	case *ast.ArrayType:
		if t.Len == nil {
			return fmt.Sprintf("typSlice(%s)", i.subtype(t.Elt))
		}
		return fmt.Sprintf("typArray(%s, %d)", i.subtype(t.Elt), intconst(t.Len))
	case *ast.ChanType:
		dir := "Cboth"
		switch t.Dir {
		case ast.SEND:
			dir = "Csend"
		case ast.RECV:
			dir = "Crecv"
		}
		return fmt.Sprintf("typChan(%s, %s)", i.subtype(t.Value), dir)
	case *ast.FuncType:
		return fmt.Sprintf("functype(nil, %s, %s)", i.fields(t.Params, false), i.fields(t.Results, false))
	case *ast.InterfaceType:
		if len(t.Methods.List) != 0 {
			log.Fatal("non-empty interfaces unsupported")
		}
		return "Types[TINTER]"
	case *ast.MapType:
		return fmt.Sprintf("typMap(%s, %s)", i.subtype(t.Key), i.subtype(t.Value))
	case *ast.StarExpr:
		return fmt.Sprintf("typPtr(%s)", i.subtype(t.X))
	case *ast.StructType:
		return fmt.Sprintf("tostruct(%s)", i.fields(t.Fields, true))

	default:
		log.Fatalf("unhandled type: %#v", t)
		panic("unreachable")
	}
}

func (i *typeInterner) fields(fl *ast.FieldList, keepNames bool) string {
	if fl == nil || len(fl.List) == 0 {
		return "nil"
	}
	var res []string
	for _, f := range fl.List {
		typ := i.subtype(f.Type)
		if len(f.Names) == 0 {
			res = append(res, fmt.Sprintf("anonfield(%s)", typ))
		} else {
			for _, name := range f.Names {
				if keepNames {
					res = append(res, fmt.Sprintf("namedfield(%q, %s)", name.Name, typ))
				} else {
					res = append(res, fmt.Sprintf("anonfield(%s)", typ))
				}
			}
		}
	}
	return fmt.Sprintf("[]*Node{%s}", strings.Join(res, ", "))
}

func intconst(e ast.Expr) int64 {
	switch e := e.(type) {
	case *ast.BasicLit:
		if e.Kind != token.INT {
			log.Fatalf("expected INT, got %v", e.Kind)
		}
		x, err := strconv.ParseInt(e.Value, 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		return x
	default:
		log.Fatalf("unhandled expr: %#v", e)
		panic("unreachable")
	}
}
