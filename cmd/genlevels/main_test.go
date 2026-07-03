package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestFolderLevelSourceInfersMetadataFromFolderAndImage(t *testing.T) {
	levelsDir := t.TempDir()
	levelDir := filepath.Join(levelsDir, "007-boat-house")
	if err := os.MkdirAll(levelDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writePNG(t, filepath.Join(levelDir, "art.png"), 24, 12)

	source, ok, err := folderLevelSource(levelsDir, "007-boat-house")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("folder was not recognized as a level")
	}
	if source.ID != "l7" {
		t.Fatalf("ID = %q, want l7", source.ID)
	}
	if source.Title != "L7 Boat House" {
		t.Fatalf("Title = %q, want L7 Boat House", source.Title)
	}
	if source.TileSize != 12 {
		t.Fatalf("TileSize = %d, want 12", source.TileSize)
	}
	if source.Origin != "007-boat-house/art.png" {
		t.Fatalf("Origin = %q, want 007-boat-house/art.png", source.Origin)
	}
}

func TestLevelFolderImageRejectsAmbiguousImages(t *testing.T) {
	levelDir := t.TempDir()
	writePNG(t, filepath.Join(levelDir, "one.png"), 20, 10)
	writePNG(t, filepath.Join(levelDir, "two.png"), 20, 10)

	if _, err := levelFolderImage(levelDir); err == nil {
		t.Fatal("expected ambiguous image error")
	}
}

func writePNG(t *testing.T, path string, width, height int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{uint8(x), uint8(y), 0, 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}
