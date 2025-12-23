package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func expandPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	result := os.ExpandEnv(path)
	if strings.HasPrefix(result, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("expand %q: %w", path, err)
		}
		switch {
		case result == "~":
			result = home
		case strings.HasPrefix(result, "~/"):
			result = filepath.Join(home, result[2:])
		default:
			// Leave ~user untouched; uncommon on modern systems.
		}
	}

	return filepath.Clean(result), nil
}

func toEnvPrefix(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r - 32)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func defaultParallelism() int {
	if n := runtime.NumCPU(); n > 0 {
		return n
	}
	return 1
}

func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
