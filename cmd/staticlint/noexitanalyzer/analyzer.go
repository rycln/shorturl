/*
Package noexitanalyzer provides a static analysis analyzer that detects and reports
direct calls to os.Exit in the main function of the main package.
*/
package noexitanalyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "noexitanalyzer is an analyzer that prohibits direct calls to os.Exit in the main function of the main package."

// Analyzer is the main analyzer variable that should be imported and added to
// multichecker. It checks for direct calls to os.Exit in the main package.
//
// The analyzer requires the inspect.Analyzer to be run first and will:
//  1. Skip files in Go build cache
//  2. Only process files in the "main" package
//  3. Report any direct calls to os.Exit()
var Analyzer = &analysis.Analyzer{
	Name:     "noexit",
	Doc:      doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	file := pass.Files[0]
	filename := pass.Fset.File(file.Pos()).Name()

	if strings.Contains(filename, "/go-build/") {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.File:
			if node.Name.Name != "main" {
				return
			}

		case *ast.CallExpr:
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
