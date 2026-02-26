package analyzer_test

import (
	"testing"
	"os"

	"log-linter/internal/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMain(m *testing.M) {
	_ = os.Unsetenv("LOGLINTER_CONFIG")
	os.Exit(m.Run())
}

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "a", "c")
}

func TestSuggestedFixes(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.Analyzer, "b")
}
