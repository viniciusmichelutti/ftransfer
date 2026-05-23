package fsutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

func SafeJoin(base, rel string) (string, error) {
	cleaned := filepath.Clean(rel)
	if filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("absolute path not allowed: %q", rel)
	}
	if strings.HasPrefix(cleaned, "..") || strings.Contains(cleaned, string(filepath.Separator)+"..") {
		return "", fmt.Errorf("path escapes base: %q", rel)
	}
	joined := filepath.Join(base, cleaned)
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absJoined, absBase+string(filepath.Separator)) && absJoined != absBase {
		return "", fmt.Errorf("path escapes base: %q", rel)
	}
	return joined, nil
}
