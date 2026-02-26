package analyzer

import (
	"go/ast"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func applyRules(pass *analysis.Pass, lit *ast.BasicLit, msg string) {
	if !beginsWithLowercase(msg) {
		diag := analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log messages should start with a lowercase letter",
		}

		if fix, ok := buildLowercaseFix(lit, msg); ok {
			diag.SuggestedFixes = []analysis.SuggestedFix{fix}
		}

		pass.Report(diag)
	}

	if containsNonASCII(msg) {
		pass.Reportf(lit.Pos(), "log messages must be written using English characters only")
	}

	if !allowSymbols && containsSymbols(msg) {
		pass.Reportf(lit.Pos(), "avoid punctuation, symbols, or emoji in log messages")
	}

	if looksSensitive(msg) {
		pass.Reportf(lit.Pos(), "possible sensitive information detected in log message")
	}
}

func buildLowercaseFix(lit *ast.BasicLit, msg string) (analysis.SuggestedFix, bool) {
	if msg == "" {
		return analysis.SuggestedFix{}, false
	}

	runes := []rune(msg)

	if len(runes) == 0 || !unicode.IsUpper(runes[0]) {
		return analysis.SuggestedFix{}, false
	}

	runes[0] = unicode.ToLower(runes[0])

	newLiteral := strconv.Quote(string(runes))

	return analysis.SuggestedFix{
		Message: "convert first letter to lowercase",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     lit.Pos(),
				End:     lit.End(),
				NewText: []byte(newLiteral),
			},
		},
	}, true
}

func looksSensitive(msg string) bool {
	lowered := strings.ToLower(msg)

	// substring patterns
	for _, k := range sensitiveSubstrings {
		if k != "" && strings.Contains(lowered, k) {
			return true
		}
	}

	// regex patterns 
	for _, re := range sensitiveRegexps {
		if re != nil && re.MatchString(msg) {
			return true
		}
	}

	return false
}
