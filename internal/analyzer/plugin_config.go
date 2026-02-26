package analyzer

import (
	"fmt"
	"strings"
)

// Supported keys (case-insensitive):
// - allow-symbols: bool
// - sensitive: string (comma-separated)
// - sensitive-patterns: []any (list of strings)
// - sensitive-regex: []any (list of regex strings)
func ApplyPluginConfig(conf any) error {
	if conf == nil {
		return nil
	}

	m, ok := conf.(map[string]any)
	if !ok {
		return fmt.Errorf("unsupported config type %T (expected map)", conf)
	}

	for k, v := range m {
		switch strings.ToLower(k) {
		case "allow-symbols":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf("allow-symbols must be bool, got %T", v)
			}
			allowSymbols = b

		case "sensitive":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("sensitive must be string, got %T", v)
			}
			sensitiveCSV = s
			setSensitiveSubstrings(splitCSV(s))

		case "sensitive-patterns":
			list, err := anySliceToStrings(v)
			if err != nil {
				return fmt.Errorf("sensitive-patterns: %w", err)
			}
			setSensitiveSubstrings(list)

		case "sensitive-regex":
			list, err := anySliceToStrings(v)
			if err != nil {
				return fmt.Errorf("sensitive-regex: %w", err)
			}
			if err := setSensitiveRegex(list); err != nil {
				return err
			}
		}
	}

	return nil
}

func anySliceToStrings(v any) ([]string, error) {
	raw, ok := v.([]any)
	if !ok {
		// some decoders may return []interface{} (alias of []any) anyway; if not, fail
		return nil, fmt.Errorf("must be a YAML list, got %T", v)
	}

	out := make([]string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("list must contain strings, got %T", item)
		}
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out, nil
}