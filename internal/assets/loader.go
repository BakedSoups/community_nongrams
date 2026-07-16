package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alex/nongrampictures/internal/nonogram"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed embedded/assets/puzzles/*/* embedded/assets/ui/*
var embeddedFiles embed.FS

type PuzzleAssets struct {
	Puzzle         *nonogram.Puzzle
	Skeleton       *ebiten.Image
	Reveal         *ebiten.Image
	SkeletonPixels [][]PixelCell
	RevealPixels   [][]PixelCell
}

type PixelCell struct {
	Color   color.RGBA
	Visible bool
}

type UIIcons struct {
	Home       *ebiten.Image
	Gear       *ebiten.Image
	Pencil     *ebiten.Image
	Eraser     *ebiten.Image
	Fill       *ebiten.Image
	Eyedropper *ebiten.Image
	Eye        *ebiten.Image
}

func LoadPuzzleAssets(puzzlePath string) (*PuzzleAssets, error) {
	puzzle, err := nonogram.LoadPuzzle(puzzlePath)
	if err != nil {
		puzzle, err = nonogram.LoadPuzzleFS(embeddedFiles, embeddedPath(puzzlePath))
		if err != nil {
			return nil, err
		}
	}

	if len(puzzle.SkeletonRaw) > 0 && len(puzzle.RevealRaw) > 0 {
		skeletonPixels, err := pixelMatrixFromRaw(puzzle.SkeletonRaw, puzzle.Width, puzzle.Height)
		if err != nil {
			return nil, err
		}
		revealPixels, err := pixelMatrixFromRaw(puzzle.RevealRaw, puzzle.Width, puzzle.Height)
		if err != nil {
			return nil, err
		}
		return &PuzzleAssets{
			Puzzle:         puzzle,
			Skeleton:       imageFromPixelMatrix(skeletonPixels),
			Reveal:         imageFromPixelMatrix(revealPixels),
			SkeletonPixels: skeletonPixels,
			RevealPixels:   revealPixels,
		}, nil
	}

	base := filepath.Dir(puzzlePath)
	skeleton, err := loadImage(resolveAsset(base, puzzle.SkeletonArt))
	if err != nil {
		skeleton, err = loadEmbeddedImage(puzzle.SkeletonArt)
		if err != nil {
			return nil, err
		}
	}
	reveal, err := loadImage(resolveAsset(base, puzzle.RevealArt))
	if err != nil {
		reveal, err = loadEmbeddedImage(puzzle.RevealArt)
		if err != nil {
			return nil, err
		}
	}

	return &PuzzleAssets{
		Puzzle:         puzzle,
		Skeleton:       skeleton.Image,
		Reveal:         reveal.Image,
		SkeletonPixels: pixelMatrix(skeleton.Source, puzzle.Width, puzzle.Height),
		RevealPixels:   pixelMatrix(reveal.Source, puzzle.Width, puzzle.Height),
	}, nil
}

func ListPuzzlePaths() []string {
	paths := map[string]string{}
	addPuzzleDirs(paths, "assets/puzzles")
	addEmbeddedPuzzleDirs(paths)

	ids := make([]string, 0, len(paths))
	for id := range paths {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	result := make([]string, 0, len(ids))
	for _, id := range ids {
		result = append(result, paths[id])
	}
	return result
}

func addPuzzleDirs(paths map[string]string, root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		puzzlePath := filepath.ToSlash(filepath.Join(root, entry.Name(), "puzzle.json"))
		if _, err := os.Stat(puzzlePath); err == nil {
			paths[entry.Name()] = puzzlePath
		}
	}
}

func addEmbeddedPuzzleDirs(paths map[string]string) {
	const root = "embedded/assets/puzzles"
	entries, err := embeddedFiles.ReadDir(root)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		puzzlePath := "assets/puzzles/" + entry.Name() + "/puzzle.json"
		embeddedPuzzlePath := root + "/" + entry.Name() + "/puzzle.json"
		if _, err := fs.Stat(embeddedFiles, embeddedPuzzlePath); err == nil {
			if _, exists := paths[entry.Name()]; !exists {
				paths[entry.Name()] = puzzlePath
			}
		}
	}
}

func LoadUIIcons() (*UIIcons, error) {
	home, err := loadImageWithFallback("assets/ui/home.png")
	if err != nil {
		return nil, err
	}
	gear, err := loadImageWithFallback("assets/ui/gear.png")
	if err != nil {
		return nil, err
	}
	pencil, err := loadImageWithFallback("assets/ui/pencil.png")
	if err != nil {
		return nil, err
	}
	eraser, err := loadImageWithFallback("assets/ui/eraser.png")
	if err != nil {
		return nil, err
	}
	fill, err := loadImageWithFallback("assets/ui/fill.png")
	if err != nil {
		return nil, err
	}
	eyedropper, err := loadImageWithFallback("assets/ui/eyedropper.png")
	if err != nil {
		return nil, err
	}
	eye, err := loadImageWithFallback("assets/ui/eye.png")
	if err != nil {
		return nil, err
	}
	return &UIIcons{Home: home, Gear: gear, Pencil: pencil, Eraser: eraser, Fill: fill, Eyedropper: eyedropper, Eye: eye}, nil
}

func resolveAsset(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return filepath.Join(base, filepath.Base(path))
}

type loadedImage struct {
	Image  *ebiten.Image
	Source image.Image
}

func loadImage(path string) (loadedImage, error) {
	f, err := os.Open(path)
	if err != nil {
		return loadedImage{}, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return loadedImage{}, err
	}
	return loadedImage{Image: ebiten.NewImageFromImage(img), Source: img}, nil
}

func loadEmbeddedImage(path string) (loadedImage, error) {
	data, err := fs.ReadFile(embeddedFiles, embeddedPath(path))
	if err != nil {
		return loadedImage{}, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return loadedImage{}, err
	}
	return loadedImage{Image: ebiten.NewImageFromImage(img), Source: img}, nil
}

func loadImageWithFallback(path string) (*ebiten.Image, error) {
	img, err := loadImage(path)
	if err == nil {
		return img.Image, nil
	}
	embedded, err := loadEmbeddedImage(path)
	if err != nil {
		return nil, err
	}
	return embedded.Image, nil
}

func embeddedPath(path string) string {
	path = filepath.ToSlash(path)
	path = strings.TrimPrefix(path, "./")
	return filepath.ToSlash(filepath.Join("embedded", path))
}

func pixelMatrix(img image.Image, width, height int) [][]PixelCell {
	if width <= 0 || height <= 0 {
		return nil
	}
	bounds := img.Bounds()
	matrix := make([][]PixelCell, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([]PixelCell, width)
		srcY := bounds.Min.Y + int((float64(y)+0.5)*float64(bounds.Dy())/float64(height))
		if srcY >= bounds.Max.Y {
			srcY = bounds.Max.Y - 1
		}
		for x := 0; x < width; x++ {
			srcX := bounds.Min.X + int((float64(x)+0.5)*float64(bounds.Dx())/float64(width))
			if srcX >= bounds.Max.X {
				srcX = bounds.Max.X - 1
			}
			r, g, b, a := img.At(srcX, srcY).RGBA()
			matrix[y][x] = PixelCell{
				Color:   color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)},
				Visible: a > 0,
			}
		}
	}
	return matrix
}

func pixelMatrixFromRaw(raw [][]string, width, height int) ([][]PixelCell, error) {
	if len(raw) != height {
		return nil, fmt.Errorf("pixel art has %d rows, expected %d", len(raw), height)
	}
	matrix := make([][]PixelCell, height)
	for y, row := range raw {
		if len(row) != width {
			return nil, fmt.Errorf("pixel art row %d has width %d, expected %d", y, len(row), width)
		}
		matrix[y] = make([]PixelCell, width)
		for x, value := range row {
			c, visible, err := parsePixel(value)
			if err != nil {
				return nil, fmt.Errorf("pixel art row %d col %d: %w", y, x, err)
			}
			matrix[y][x] = PixelCell{Color: c, Visible: visible}
		}
	}
	return matrix, nil
}

func parsePixel(value string) (color.RGBA, bool, error) {
	if value == "" || value == "transparent" {
		return color.RGBA{}, false, nil
	}
	if len(value) != 9 || value[0] != '#' {
		return color.RGBA{}, false, fmt.Errorf("expected #RRGGBBAA, got %q", value)
	}
	var rgba [4]uint8
	for i := 0; i < 4; i++ {
		n, ok := parseHexByte(value[1+i*2 : 3+i*2])
		if !ok {
			return color.RGBA{}, false, fmt.Errorf("expected #RRGGBBAA, got %q", value)
		}
		rgba[i] = n
	}
	c := color.RGBA{R: rgba[0], G: rgba[1], B: rgba[2], A: rgba[3]}
	return c, c.A > 0, nil
}

func parseHexByte(value string) (uint8, bool) {
	var out uint8
	for _, ch := range value {
		var n uint8
		switch {
		case ch >= '0' && ch <= '9':
			n = uint8(ch - '0')
		case ch >= 'a' && ch <= 'f':
			n = uint8(ch-'a') + 10
		case ch >= 'A' && ch <= 'F':
			n = uint8(ch-'A') + 10
		default:
			return 0, false
		}
		out = out*16 + n
	}
	return out, true
}

func imageFromPixelMatrix(matrix [][]PixelCell) *ebiten.Image {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return ebiten.NewImage(1, 1)
	}
	img := image.NewRGBA(image.Rect(0, 0, len(matrix[0]), len(matrix)))
	for y, row := range matrix {
		for x, cell := range row {
			if cell.Visible {
				img.SetRGBA(x, y, cell.Color)
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}
