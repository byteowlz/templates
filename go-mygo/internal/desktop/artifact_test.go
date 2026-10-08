package desktop

import (
	"image"
	"image/png"
	"os"
	"testing"
)

func writePNG(t *testing.T, path string, picture image.Image) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, picture); err != nil {
		t.Fatal(err)
	}
}
