//go:build !solution

package testifycheck

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "require",
	Doc:  "",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	names := map[string]string{
		"Nil":     "NoError",
		"Nilf":    "NoErrorf",
		"NotNil":  "Error",
		"NotNilf": "Errorf",
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			if _, ok := n.(*ast.ReturnStmt); ok {
				return false
			}

			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fn, _ := typeutil.Callee(pass.TypesInfo, ce).(*types.Func)
			if fn == nil {
				return true
			}
			if fn.Pkg().Name() != "require" && fn.Pkg().Name() != "assert" {
				return true
			}

			isErr := func(expr ast.Expr) bool {
				t := pass.TypesInfo.TypeOf(expr)
				if t == nil {
					return false
				}

				intf, ok := t.Underlying().(*types.Interface)
				if !ok {
					return false
				}

				return intf.NumMethods() == 1 && intf.Method(0).FullName() == "(error).Error"
			}

			argsLen := len(ce.Args)
			if argsLen < 1 {
				return true
			}
			if !isErr(ce.Args[0]) && !(argsLen > 1 && isErr(ce.Args[1])) {
				return true
			}

			if n, ok := names[fn.Name()]; ok {
				pass.Reportf(ce.Pos(), "use %s.%s instead of comparing error to nil", fn.Pkg().Name(), n)
				return false
			}
			return true
		})
	}

	return nil, nil
}
