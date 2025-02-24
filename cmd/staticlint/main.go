package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"honnef.co/go/tools/analysis/facts/nilness"

	"honnef.co/go/tools/staticcheck"

	"github.com/kisielk/errcheck/errcheck"
	"github.com/mdempsky/maligned/passes/maligned"
)

func main() {
	analyzers := make([]*analysis.Analyzer, 0)
	for _, a := range staticcheck.Analyzers {
		analyzers = append(analyzers, a.Analyzer)
	}

	analyzers = append(analyzers,
		nilness.Analysis,
		shadow.Analyzer,
		printf.Analyzer,
		inspect.Analyzer,
		errcheck.Analyzer,
		maligned.Analyzer,
		osExitAnalyzer,
	)

	multichecker.Main(analyzers...)
}

var osExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit using",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Walk(&exitWalker{pass}, f)
	}

	return nil, nil
}

type exitWalker struct {
	pass *analysis.Pass
}

func (w *exitWalker) Visit(node ast.Node) ast.Visitor {
	if fn, ok := node.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "Exit" {
					w.pass.Reportf(n.Pos(), "нельзя использовать функцию os.Exit в main")
				}
			}
			return true
		})
	}
	return w
}
