package widgets

import (
	"os"
	"path/filepath"
	"testing"
)

func mustMkdir(t *testing.T, root string, rel string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(rel)), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
}
