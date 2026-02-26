package main

import (
	"fmt"

	"golang.org/x/tools/go/analysis"

	"log-linter/internal/analyzer"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	if err := analyzer.ApplyPluginConfig(conf); err != nil {
		return nil, fmt.Errorf("apply config: %w", err)
	}

	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}