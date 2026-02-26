package analyzer

import (
	"flag"
	"strings"
	"golang.org/x/tools/go/analysis"
)

var (
	allowSymbols bool
	sensitiveCSV string
)

var Analyzer = &analysis.Analyzer{
	Name: "loglinter",
	Doc:  "validates log message formatting and safety rules",
	Run:  run,
}

func init() {
	Analyzer.Flags = *flag.NewFlagSet("loglinter", flag.ContinueOnError)

	Analyzer.Flags.BoolVar(&allowSymbols, "allow-symbols", false, "allow punctuation/symbols in log messages")
	Analyzer.Flags.StringVar(&sensitiveCSV, "sensitive", "password,token,secret,api_key,apikey,credential,passwd", "comma-separated sensitive keywords")

	// from default CSV
	setSensitiveSubstrings(splitCSV(sensitiveCSV))
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Highest priority override: LOGLINTER_CONFIG
	if err := ensureConfigLoaded(); err != nil {
		return nil, err
	}

	// Keep substring matcher in sync if someone set sensitiveCSV via flags/plugin settings.
	// (If file config set Sensitive list explicitly, it already set matchers.)
	setSensitiveSubstrings(splitCSV(sensitiveCSV))

	for _, file := range pass.Files {
		inspectFile(pass, file)
	}
	return nil, nil
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := make([]string, 0, 8)
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}