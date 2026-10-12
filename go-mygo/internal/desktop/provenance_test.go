package desktop

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAdapterSnapshotDrift(t *testing.T) {
	root := filepath.Join("..", "..")
	raw, err := os.ReadFile(filepath.Join(root, "third_party", "mygo-go.source.json"))
	if err != nil {
		t.Fatal(err)
	}
	var pin struct {
		Repository, Revision string
		Files                map[string]string
	}
	if err := json.Unmarshal(raw, &pin); err != nil {
		t.Fatal(err)
	}
	if pin.Repository != "https://github.com/byteowlz/design-system" || pin.Revision != "c31fdc8ea0f685db15dcd5a9103f1d981db6c3de" {
		t.Fatal("unexpected canonical adapter pin")
	}
	for path, want := range pin.Files {
		data, err := os.ReadFile(filepath.Join(root, "third_party", "mygo-go", path))
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != want {
			t.Fatalf("canonical adapter snapshot drift: %s", path)
		}
	}
}
