package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

const outDir = "assets/puzzles/test_001"
const uiDir = "assets/ui"

func main() {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(uiDir, 0o755); err != nil {
		log.Fatal(err)
	}
	write("skeleton.png", makeArt(true))
	write("full_art.png", makeArt(false))
	writeUI("home.png", makeHomeIcon())
	writeUI("gear.png", makeGearIcon())
	writeUI("pencil.png", makePencilIcon())
	writeUI("eraser.png", makeEraserIcon())
}

func write(name string, img image.Image) {
	path := filepath.Join(outDir, name)
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func writeUI(name string, img image.Image) {
	path := filepath.Join(uiDir, name)
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func makeArt(skeleton bool) image.Image {
	const size = 160
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	bg := color.RGBA{238, 229, 208, 255}
	if skeleton {
		bg = color.RGBA{226, 218, 202, 255}
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	pixel := 10
	palette := map[rune]color.RGBA{
		'1': {83, 103, 92, 255},
		'2': {93, 126, 104, 255},
		'3': {190, 115, 94, 255},
		'4': {68, 78, 84, 255},
	}
	if skeleton {
		palette = map[rune]color.RGBA{
			'1': {156, 150, 139, 115},
			'2': {156, 150, 139, 95},
			'3': {156, 150, 139, 85},
			'4': {156, 150, 139, 105},
		}
	}

	rows := []string{
		"0000110000000000",
		"0001221000000000",
		"0012222100000000",
		"0122222210000000",
		"0003333000000000",
		"0003333000000000",
		"0003333001100000",
		"0003333010010000",
		"0003333010010000",
		"0003333001100000",
		"0003333000000000",
		"0000330000000000",
		"0000330000000000",
		"0000330000000000",
		"0004444000000000",
		"0000000000000000",
	}

	for y, row := range rows {
		for x, ch := range row {
			c, ok := palette[ch]
			if !ok {
				continue
			}
			rect := image.Rect(x*pixel, y*pixel, x*pixel+pixel, y*pixel+pixel)
			draw.Draw(img, rect, &image.Uniform{c}, image.Point{}, draw.Src)
		}
	}
	return img
}

func iconCanvas() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, 64, 64))
}

func makeHomeIcon() image.Image {
	img := iconCanvas()
	ink := color.RGBA{45, 44, 40, 255}
	line(img, 16, 34, 32, 18, ink, 4)
	line(img, 32, 18, 48, 34, ink, 4)
	line(img, 21, 34, 21, 48, ink, 4)
	line(img, 43, 34, 43, 48, ink, 4)
	line(img, 21, 48, 43, 48, ink, 4)
	line(img, 28, 48, 28, 40, ink, 4)
	line(img, 36, 40, 36, 48, ink, 4)
	line(img, 28, 40, 36, 40, ink, 4)
	return img
}

func makeGearIcon() image.Image {
	img := iconCanvas()
	ink := color.RGBA{45, 44, 40, 255}

	fillRect(img, 28, 7, 8, 11, ink)
	fillRect(img, 28, 46, 8, 11, ink)
	fillRect(img, 7, 28, 11, 8, ink)
	fillRect(img, 46, 28, 11, 8, ink)
	fillRect(img, 14, 13, 8, 8, ink)
	fillRect(img, 42, 13, 8, 8, ink)
	fillRect(img, 14, 43, 8, 8, ink)
	fillRect(img, 42, 43, 8, 8, ink)

	circleStroke(img, 32, 32, 16, ink, 6)
	circleStroke(img, 32, 32, 5, ink, 4)
	return img
}

func makePencilIcon() image.Image {
	img := iconCanvas()
	ink := color.RGBA{45, 44, 40, 255}
	line(img, 18, 45, 43, 20, ink, 4)
	line(img, 24, 51, 49, 26, ink, 4)
	line(img, 18, 45, 24, 51, ink, 4)
	line(img, 43, 20, 49, 26, ink, 4)
	line(img, 49, 26, 54, 16, ink, 4)
	line(img, 43, 20, 54, 16, ink, 4)
	line(img, 14, 50, 23, 58, ink, 4)
	return img
}

func makeEraserIcon() image.Image {
	img := iconCanvas()
	ink := color.RGBA{45, 44, 40, 255}
	line(img, 19, 25, 47, 25, ink, 4)
	line(img, 47, 25, 40, 44, ink, 4)
	line(img, 40, 44, 13, 44, ink, 4)
	line(img, 13, 44, 19, 25, ink, 4)
	line(img, 38, 29, 33, 41, ink, 4)
	line(img, 19, 53, 45, 53, ink, 4)
	return img
}

func line(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA, width int) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		dot(img, x0, y0, width, c)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func circleStroke(img *image.RGBA, cx, cy, radius int, c color.RGBA, width int) {
	for a := 0; a < 360; a += 5 {
		rad := float64(a) * 0.01745329252
		dot(img, cx+int(mathCos(rad)*float64(radius)), cy+int(mathSin(rad)*float64(radius)), width, c)
	}
}

func dot(img *image.RGBA, cx, cy, width int, c color.RGBA) {
	r := width / 2
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if image.Pt(x, y).In(img.Bounds()) {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func fillRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			if image.Pt(xx, yy).In(img.Bounds()) {
				img.SetRGBA(xx, yy, c)
			}
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func mathSin(v float64) float64 {
	switch {
	case v < 1.5708:
		return v - v*v*v/6 + v*v*v*v*v/120
	case v < 3.1416:
		return mathSin(3.14159265359 - v)
	case v < 4.7124:
		return -mathSin(v - 3.14159265359)
	default:
		return -mathSin(6.28318530718 - v)
	}
}

func mathCos(v float64) float64 {
	return mathSin(v + 1.57079632679)
}
