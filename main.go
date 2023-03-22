package main

import (
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
	"text/template"
)

type Interface struct {
	Name     string
	Methods  MethodList
	Generics VariableList
}

type Method struct {
	Name    string
	Args    VariableList
	Results Results
}

type Results []string

func (r Results) String() string {
	if len(r) == 0 {
		return ""
	}
	if len(r) == 1 {
		return r[0]
	}

	return "(" + strings.Join(r, ", ") + ")"
}

func (r Results) Return() string {
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

func FieldListToVariableList(fl *ast.FieldList) VariableList {
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

func (vl VariableList) Types() (types []string) {
	types = make([]string, len(vl))
	i := 0
	for _, t := range vl {
		types[i] = t.Type
		i++
	}
	return
}

func (vl VariableList) Names() []string {
	names := make([]string, 0)
	for _, t := range vl {
		names = append(names, t.Name)
	}

	return names
}

func (vl VariableList) ToParams() string {
	params := make([]string, len(vl))

	for i, v := range vl {
		params[i] = v.Name + " " + v.Type
	}

	return strings.Join(params, ", ")
}

func (vl VariableList) ToGenericParams() string {
	return "[" + vl.ToParams() + "]"
}

func (vl VariableList) ToArgs() string {
	args := make([]string, len(vl))

	for i, v := range vl {
		args[i] = v.Name
	}

	return strings.Join(args, ", ")
}

func (vl VariableList) ToGenericArgs() string {
	return "[" + vl.ToArgs() + "]"
}

func (ts VariableList) StringWith(prefix, suffix string) string {
	if len(ts) == 0 {
		return ""
	}

	var types []string
	for _, t := range ts {
		types = append(types, t.Name+" "+t.Type)
	}

	return prefix + strings.Join(types, ", ") + suffix
}

func (ts VariableList) String() string {
	return ts.StringWith("", "")
}

func StringSliceFromFieldList(fl *ast.FieldList) []string {
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

func MethodListFromInterfaceType(t *ast.InterfaceType) MethodList {
	methods := []Method{}

	for _, m := range t.Methods.List {
		im := Method{}

		im.Name = m.Names[0].Name

		fn, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		im.Args = FieldListToVariableList(fn.Params)
		im.Results = StringSliceFromFieldList(fn.Results)

		methods = append(methods, im)
	}

	return methods
}

func InterfaceFromTypeSpec(t *ast.TypeSpec) (Interface, error) {
	im := Interface{}

	it, ok := t.Type.(*ast.InterfaceType)
	if !ok {
		return im, errors.New("not an *ast.InterfaceType")
	}

	im.Name = t.Name.Name
	im.Methods = MethodListFromInterfaceType(it)
	im.Generics = FieldListToVariableList(t.TypeParams)

	return im, nil
}

type File struct {
	pkg        string
	interfaces []Interface
}

func Inspect(src string, structs []string) (File, error) {
	file := File{}

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
			file.pkg = t.Name.Name
		case *ast.TypeSpec:
			if intf, err := InterfaceFromTypeSpec(t); err == nil {
				if _, ok := filter[intf.Name]; !ok && len(filter) > 0 {
					break
				}
				file.interfaces = append(file.interfaces, intf)
			}
		}
		return true
	})

	return file, nil
}

//go:embed mock.tmpl
var mock string

//func Render(w io.Writer, src string) error {
func (f File) Render(w io.Writer) error {
	tmpl := template.New("mock").
		Funcs(map[string]any{
			"lower": func(s string) string {
				b := []byte(s)
				b[0] = b[0] + 32
				return string(b)
			},
			"join": strings.Join,
			"args": func(params map[string]string) string {
				args := []string{}
				for k, v := range params {
					args = append(args, k+" "+v)
				}
				return strings.Join(args, ", ")
			},
			"renderTypeArgs": func(params []string) string {
				if len(params) == 0 {
					return ""
				}

				return "[" + strings.Join(params, ", ") + "]"
			},
		})

	tmpl, err := tmpl.Parse(mock)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, struct {
		PkgName    string
		Interfaces []Interface
	}{
		f.pkg,
		f.interfaces,
	})
	if err != nil {
		panic(err)
	}

	return nil
}

func Load(filename string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	file, err := os.Open(filename)
	if err != nil {
		return &buf, err
	}

	_, err = io.Copy(&buf, file)
	return &buf, err
}

var (
	filename     string
	structFilter string
)

func main() {

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

	src, err := Load(filename)
	if err != nil {
		panic(err)
	}

	file, err := Inspect(src.String(), strings.Split(structFilter, ","))
	if err != nil {
		panic(err)
	}

	file.Render(os.Stdout)
}
