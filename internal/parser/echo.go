package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type EchoParser struct{}

func NewEchoParser() *EchoParser {
	return &EchoParser{}
}

func (p *EchoParser) Parse(file string) ([]Route, error) {
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

		if !isEchoMethod(call) {
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

func isEchoMethod(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	method := strings.ToUpper(sel.Sel.Name)
	return method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" ||
		method == "PATCH" || method == "HEAD" || method == "OPTIONS"
}
