package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	//"log/slog"
	"golang.org/x/tools/go/analysis"
)

/* Если необходима проверка
func Demo() {
	slog.Info("Starting server")
}
*/

/* Если необходима проверка
func Demo() {
	slog.Info("Hello!!!")
}
*/

var supportedMethods = map[string]struct{}{
	"Debug":  {},
	"Info":   {},
	"Warn":   {},
	"Error":  {},
	"Fatal":  {},
	"DPanic": {},
	"Panic":  {},
}

func inspectFile(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isSupportedLoggerCall(pass, call) {
			return true
		}

		if len(call.Args) == 0 {
			return true
		}

		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok {
			return true
		}

		if lit.Kind != token.STRING {
			return true
		}

		msg, err := strconv.Unquote(lit.Value)
		if err != nil {
			return true
		}

		applyRules(pass, lit, msg)
		return true
	})
}

func isSupportedLoggerCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	method := sel.Sel.Name

	if _, ok := supportedMethods[method]; !ok {
		return false
	}

	// 1️⃣ Попытка через type-check (production режим)
	if recvType := pass.TypesInfo.TypeOf(sel.X); recvType != nil {
		if isSlog(recvType) || isZap(recvType) {
			return true
		}
	}

	// 2️⃣ Fallback для analysistest

	// slog.Info(...)
	if ident, ok := sel.X.(*ast.Ident); ok {
		if ident.Name == "slog" {
			return true
		}
	}

	// zap.L().Info(...) или zap.S().Info(...)
	if innerCall, ok := sel.X.(*ast.CallExpr); ok {
		if innerSel, ok := innerCall.Fun.(*ast.SelectorExpr); ok {
			if pkgIdent, ok := innerSel.X.(*ast.Ident); ok {
				if pkgIdent.Name == "zap" {
					return true
				}
			}
		}
	}

	return false
}

func unwrapNamed(t types.Type) *types.Named {
	switch tt := t.(type) {
	case *types.Named:
		return tt
	case *types.Pointer:
		if named, ok := tt.Elem().(*types.Named); ok {
			return named
		}
	}
	return nil
}

func isSlog(t types.Type) bool {
	named := unwrapNamed(t)
	if named == nil {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == "log/slog"
}

func isZap(t types.Type) bool {
	named := unwrapNamed(t)
	if named == nil {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == "go.uber.org/zap"
}