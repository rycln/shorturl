package noexitanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "noexitanalyzer is an analyzer that prohibits direct calls to os.Exit in the main function of the main package."

var Analyzer = &analysis.Analyzer{
	Name:     "noexit",
	Doc:      doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.File:
			// Проверяем, что мы в пакете main
			if node.Name.Name != "main" {
				return
			}

		case *ast.CallExpr:
			// Проверяем, что вызов является os.Exit
			sel, ok := node.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			ident, ok := sel.X.(*ast.Ident)
			if !ok {
				return
			}

			if ident.Name == "os" && sel.Sel.Name == "Exit" {
				pass.Reportf(node.Pos(), "direct call to os.Exit in main function of main package")
			}
		}
	})

	return nil, nil
}
