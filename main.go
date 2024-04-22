package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Interface struct {
	Name     string
	Methods  MethodList
	Generics GenericList
}

type Method struct {
	Name    string
	Params  ParamList
	Returns Returns
}

type Returns []string

func (r Returns) Declaration() string {
	if len(r) == 0 {
		return ""
	}
	if len(r) == 1 {
		return r[0]
	}

	return "(" + strings.Join(r, ", ") + ")"
}

func (r Returns) String() string {
	if len(r) == 0 {
		return ""
	}

	defaults := make([]string, len(r))
	for i, t := range r {
		var val string

		switch t {
		case "any":
			val = "nil"
		case "bool":
			val = "false"
		case "complex64", "complex128":
			val = "complex(0, 0)"
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint32", "uint64", "float32", "float64", "byte", "rune":
			val = "0"
		case "error":
			val = "errors.New(\"not implemented\")"
		case "string":
			val = "\"\""
		default:
			val = "*new(" + t + ")"
			if isPointer(t) {
				val = "new(" + t[1:] + ")"
			}
		}

		defaults[i] = val
	}

	return strings.Join(defaults, ", ")
}

type MethodList []Method

type Variable struct {
	Name string
	Type string
}

func isPointer(s string) bool {
	return len(s) > 1 && s[0] == '*'
}

type VariableList []Variable

func (vl VariableList) ToParams() string {
	params := make([]string, len(vl))

	for i, v := range vl {
		params[i] = v.Name + " " + v.Type
	}

	return strings.Join(params, ", ")
}

func (vl VariableList) ToArgs() string {
	args := make([]string, len(vl))

	for i, v := range vl {
		args[i] = v.Name
	}

	return strings.Join(args, ", ")
}

type ParamList struct {
	VariableList
}

func (pl ParamList) Declaration() string {
	if len(pl.VariableList) == 0 {
		return ""
	}

	return pl.ToParams()
}

func (pl ParamList) Args() string {
	args := make([]string, len(pl.VariableList))

	for i, v := range pl.VariableList {
		args[i] = v.Name
	}

	return strings.Join(args, ", ")
}

type GenericList struct {
	VariableList
}

func (gl GenericList) Declaration() string {
	if len(gl.VariableList) == 0 {
		return ""
	}

	return "[" + gl.ToParams() + "]"
}

func (gl GenericList) Args() string {
	if len(gl.VariableList) == 0 {
		return ""
	}

	args := make([]string, len(gl.VariableList))

	for i, v := range gl.VariableList {
		args[i] = v.Name
	}

	return "[" + strings.Join(args, ", ") + "]"
}

func NewVariableList(fl *ast.FieldList) VariableList {
	if fl == nil {
		return nil
	}

	var vars VariableList
	for _, f := range fl.List {
		t, ok := f.Type.(*ast.Ident)
		if !ok {
			continue
		}

		for _, n := range f.Names {
			vars = append(vars, Variable{Name: n.Name, Type: t.Name})
		}
	}

	return vars
}

func ExtractVariableNames(fl *ast.FieldList) []string {
	if fl == nil {
		return nil
	}

	sl := make([]string, len(fl.List))
	for i, f := range fl.List {
		switch t := f.Type.(type) {
		case *ast.Ident:
			sl[i] = t.Name
		case *ast.StarExpr:
			pt, ok := t.X.(*ast.Ident)
			if !ok {
				continue
			}
			sl[i] = "*" + pt.Name
		}
	}

	return sl
}

func NewMethodList(t *ast.InterfaceType) MethodList {
	methods := []Method{}

	for _, m := range t.Methods.List {
		im := Method{}

		im.Name = m.Names[0].Name

		fn, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		im.Params.VariableList = NewVariableList(fn.Params)
		im.Returns = ExtractVariableNames(fn.Results)

		methods = append(methods, im)
	}

	return methods
}

func NewInterface(t *ast.TypeSpec) (Interface, error) {
	im := Interface{}

	it, ok := t.Type.(*ast.InterfaceType)
	if !ok {
		return im, errors.New("not an *ast.InterfaceType")
	}

	im.Name = t.Name.String()
	im.Methods = NewMethodList(it)
	im.Generics.VariableList = NewVariableList(t.TypeParams)

	return im, nil
}

func Inspect(src []byte, structs []string) (Source, error) {
	file := Source{}

	filter := make(map[string]struct{})
	for _, intf := range structs {
		filter[intf] = struct{}{}
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.File:
			file.Pkg = t.Name.Name
		case *ast.TypeSpec:
			if intf, err := NewInterface(t); err == nil {
				if len(filter) > 0 {
					if _, ok := filter[intf.Name]; !ok {
						break
					}
				}

				file.Interfaces = append(file.Interfaces, intf)
			}
		}
		return true
	})

	return file, nil
}

type Source struct {
	Pkg        string
	Interfaces []Interface
}

//go:embed mock.tmpl
var stubTemplate string

func (src Source) Render(w io.Writer) error {
	if stubTemplate == "" {
		panic(errors.New("empty template"))
	}

	tmpl, err := template.New("mock").Parse(stubTemplate)
	if err != nil {
		panic(err)
	}

	importsCmd := exec.Command("goimports")
	stdin, err := importsCmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	//err = tmpl.Execute(w, struct{ Src Source }{src})
	err = tmpl.Execute(stdin, struct{ Src Source }{src})
	if err != nil {
		panic(err)
	}
	stdin.Close()

	out, err := importsCmd.Output()
	if err != nil {
		panic(err)
	}

	w.Write(out)

	return nil
}

func main() {
	var filename, structFilter string

	fs := flag.NewFlagSet("stub", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Printf(`gostub

Stub out your interfaces to use in tests or WIP code

Usage
stub -filename src.go -filter MyInterface

`)
		fs.PrintDefaults()
	}

	fs.StringVar(&filename, "filename", "", "the filename containing the interfaces you want to stub")
	fs.StringVar(&structFilter, "filter", "", "a CSV list of structs to filter by")

	fs.Parse(os.Args[1:])

	if filename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var filter []string
	if structFilter != "" {
		filter = strings.Split(structFilter, ",")
	}

	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	file, err := Inspect(src, filter)
	if err != nil {
		panic(err)
	}

	file.Render(os.Stdout)
}
