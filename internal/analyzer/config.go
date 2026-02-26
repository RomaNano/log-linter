package analyzer

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type FileConfig struct {
	AllowSymbols   *bool    `yaml:"allow_symbols"`
	Sensitive      []string `yaml:"sensitive"`
	SensitiveRegex []string `yaml:"sensitive_regex"`
}

var (
	cfgOnce sync.Once
	cfgErr  error
)

// ensureConfigLoaded applies config in the following order (low -> high priority):
// 1) defaults (from analyzer.go init)
// 2) golangci-lint plugin settings (ApplyPluginConfig)
// 3) LOGLINTER_CONFIG file (this function)  <-- highest priority
func ensureConfigLoaded() error {
	cfgOnce.Do(func() {
		path := strings.TrimSpace(os.Getenv("LOGLINTER_CONFIG"))
		if path == "" {
			return
		}

		b, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				// silently ignore missing config file
				return
			}
			cfgErr = fmt.Errorf("read LOGLINTER_CONFIG %q: %w", path, err)
			return
		}

		var fc FileConfig
		if err := yaml.Unmarshal(b, &fc); err != nil {
			cfgErr = fmt.Errorf("parse LOGLINTER_CONFIG %q: %w", path, err)
			return
		}

		if err := ApplyFileConfig(fc); err != nil {
			cfgErr = fmt.Errorf("apply LOGLINTER_CONFIG %q: %w", path, err)
			return
		}
	})

	return cfgErr
}

func ApplyFileConfig(fc FileConfig) error {
	if fc.AllowSymbols != nil {
		allowSymbols = *fc.AllowSymbols
	}

	if len(fc.Sensitive) > 0 {
		setSensitiveSubstrings(fc.Sensitive)
	}

	if len(fc.SensitiveRegex) > 0 {
		if err := setSensitiveRegex(fc.SensitiveRegex); err != nil {
			return err
		}
	}

	return nil
}

// ----- sensitive matchers (shared by plugin config + file config) -----

var (
	sensitiveSubstrings []string
	sensitiveRegexps    []*regexp.Regexp
)

func setSensitiveSubstrings(list []string) {
	out := make([]string, 0, len(list))
	for _, s := range list {
		s = strings.TrimSpace(strings.ToLower(s))
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	sensitiveSubstrings = out
	// keep legacy CSV in sync (optional)
	sensitiveCSV = strings.Join(out, ",")
}

func setSensitiveRegex(patterns []string) error {
	out := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return fmt.Errorf("invalid sensitive_regex %q: %w", p, err)
		}
		out = append(out, re)
	}
	if len(out) == 0 {
		return errors.New("sensitive_regex provided but no valid patterns found")
	}
	sensitiveRegexps = out
	return nil
}