package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Route struct {
	Method   string
	Path     string
	Handler  string
	Package  string
	Line     int
	FilePath string
}

type Parser interface {
	Parse(file string) ([]Route, error)
}

type GinParser struct{}

func NewGinParser() *GinParser {
	return &GinParser{}
}

func (p *GinParser) Parse(file string) ([]Route, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var routes []Route

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isGinMethod(call) {
			return true
		}

		if len(call.Args) < 2 {
			return true
		}

		pathLit, ok := call.Args[0].(*ast.BasicLit)
		if !ok {
			return true
		}

		path := stripQuotes(pathLit.Value)
		if path == "" {
			return true
		}

		method := extractMethod(call)
		if method == "" {
			return true
		}

		handler := extractHandler(call.Args[1])
		if handler == "" {
			return true
		}

		pos := fset.Position(call.Pos())

		routes = append(routes, Route{
			Method:   method,
			Path:     path,
			Handler:  handler,
			Package:  node.Name.Name,
			Line:     pos.Line,
			FilePath: file,
		})

		return true
	})

	return routes, nil
}

func isGinMethod(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name != "router" && ident.Name != "r" {
		return false
	}

	method := sel.Sel.Name
	return method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" ||
		method == "PATCH" || method == "HEAD" || method == "OPTIONS"
}

func extractMethod(call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	return sel.Sel.Name
}

func extractHandler(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Ident:
		return n.Name
	case *ast.SelectorExpr:
		if ident, ok := n.X.(*ast.Ident); ok {
			return ident.Name + "." + n.Sel.Name
		}
		return n.Sel.Name
	case *ast.FuncLit:
		return "anonymous"
	}
	return ""
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '`' && s[len(s)-1] == '`') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
