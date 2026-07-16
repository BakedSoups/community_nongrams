package game

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/alex/nongrampictures/internal/assets"
	"github.com/alex/nongrampictures/internal/nonogram"
)

type editorMode uint8

const (
	editorModeArt editorMode = iota
	editorModeSolution
)

type editorTool uint8

const (
	editorToolPencil editorTool = iota
	editorToolEraser
	editorToolFill
	editorToolEyedropper
)

type editorCell struct {
	Color   color.RGBA `json:"color"`
	Visible bool       `json:"visible"`
	Filled  bool       `json:"filled"`
}

type editorState struct {
	Width      int          `json:"width"`
	Height     int          `json:"height"`
	Title      string       `json:"title"`
	Cells      []editorCell `json:"cells"`
	PaintColor color.RGBA   `json:"paintColor"`
	Mode       editorMode   `json:"mode"`
	Tool       editorTool   `json:"tool"`
}

type editorImportCell struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

type editorImportPayload struct {
	Width  int                `json:"width"`
	Height int                `json:"height"`
	Cells  []editorImportCell `json:"cells"`
}

type editorPack struct {
	ID        string             `json:"id"`
	Title     string             `json:"title"`
	Author    string             `json:"author"`
	Version   int                `json:"version"`
	UpdatedAt string             `json:"updatedAt"`
	Levels    []*nonogram.Puzzle `json:"levels"`
}

func newEditorState(size int) editorState {
	if size <= 0 {
		size = 10
	}
	cells := make([]editorCell, size*size)
	return editorState{
		Width:      size,
		Height:     size,
		Title:      fmt.Sprintf("Editor %dx%d", size, size),
		Cells:      cells,
		PaintColor: color.RGBA{72, 92, 121, 255},
		Mode:       editorModeArt,
		Tool:       editorToolPencil,
	}
}

func editorFromPackJSON(raw string) (editorState, error) {
	var pack editorPack
	if err := json.Unmarshal([]byte(raw), &pack); err != nil {
		return editorState{}, err
	}
	if len(pack.Levels) == 0 {
		return editorState{}, fmt.Errorf("pack has no levels")
	}
	p := pack.Levels[0]
	if p.Width <= 0 || p.Height <= 0 {
		return editorState{}, fmt.Errorf("pack level has invalid dimensions")
	}
	state := newEditorState(p.Width)
	state.Height = p.Height
	state.Title = p.Title
	state.Cells = make([]editorCell, p.Width*p.Height)
	if len(p.RevealRaw) == p.Height {
		for y, row := range p.RevealRaw {
			for x, value := range row {
				if x >= p.Width || value == "" || value == "transparent" {
					continue
				}
				c, ok := parseEditorHexColor(value)
				if !ok {
					continue
				}
				state.Cells[state.index(x, y)] = editorCell{Color: c, Visible: c.A > 0}
			}
		}
	}
	if len(p.SolutionRaw) == p.Height {
		for y, row := range p.SolutionRaw {
			for x, ch := range row {
				if x < p.Width && ch == '1' {
					state.Cells[state.index(x, y)].Filled = true
				}
			}
		}
	}
	return state, nil
}

func (e editorState) clone() editorState {
	next := e
	next.Cells = append([]editorCell(nil), e.Cells...)
	return next
}

func (e editorState) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < e.Width && y < e.Height
}

func (e editorState) index(x, y int) int {
	return y*e.Width + x
}

func (e *editorState) cell(x, y int) *editorCell {
	if !e.inBounds(x, y) {
		return nil
	}
	return &e.Cells[e.index(x, y)]
}

func (e *editorState) apply(x, y int) {
	cell := e.cell(x, y)
	if cell == nil {
		return
	}
	switch e.Tool {
	case editorToolPencil:
		cell.Color = e.PaintColor
		cell.Visible = true
	case editorToolEraser:
		cell.Visible = false
		cell.Filled = false
	case editorToolEyedropper:
		if cell.Visible {
			e.PaintColor = cell.Color
			e.Tool = editorToolPencil
		}
	case editorToolFill:
		e.fillArt(x, y)
	}
	e.autoSolutionFromVisible()
}

func (e *editorState) applyLine(x0, y0, x1, y1 int) {
	dx := absInt(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -absInt(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	for {
		e.apply(x0, y0)
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

func (e *editorState) fillArt(startX, startY int) {
	start := e.cell(startX, startY)
	if start == nil {
		return
	}
	targetColor := start.Color
	targetVisible := start.Visible
	if targetVisible && sameRGBA(targetColor, e.PaintColor) {
		return
	}
	queue := [][2]int{{startX, startY}}
	seen := make(map[[2]int]bool, e.Width*e.Height)
	for len(queue) > 0 {
		p := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if seen[p] || !e.inBounds(p[0], p[1]) {
			continue
		}
		seen[p] = true
		cell := e.cell(p[0], p[1])
		if cell.Visible != targetVisible || (targetVisible && !sameRGBA(cell.Color, targetColor)) {
			continue
		}
		cell.Color = e.PaintColor
		cell.Visible = true
		queue = append(queue, [2]int{p[0] + 1, p[1]}, [2]int{p[0] - 1, p[1]}, [2]int{p[0], p[1] + 1}, [2]int{p[0], p[1] - 1})
	}
}

func (e *editorState) autoSolutionFromVisible() {
	for i := range e.Cells {
		e.Cells[i].Filled = e.Cells[i].Visible
	}
}

func (e *editorState) autoSolutionFromBrightness() {
	for i := range e.Cells {
		cell := &e.Cells[i]
		if !cell.Visible {
			cell.Filled = false
			continue
		}
		luma := 0.2126*float64(cell.Color.R) + 0.7152*float64(cell.Color.G) + 0.0722*float64(cell.Color.B)
		cell.Filled = luma < 196
	}
}

func (e *editorState) invertSolution() {
	for i := range e.Cells {
		e.Cells[i].Filled = !e.Cells[i].Filled
	}
}

func (e *editorState) applyBrightness(delta int) {
	for i := range e.Cells {
		if !e.Cells[i].Visible {
			continue
		}
		e.Cells[i].Color.R = clampByte(int(e.Cells[i].Color.R) + delta)
		e.Cells[i].Color.G = clampByte(int(e.Cells[i].Color.G) + delta)
		e.Cells[i].Color.B = clampByte(int(e.Cells[i].Color.B) + delta)
	}
}

func (e *editorState) applySaturation(delta float64) {
	factor := 1 + delta
	for i := range e.Cells {
		if !e.Cells[i].Visible {
			continue
		}
		c := e.Cells[i].Color
		gray := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
		e.Cells[i].Color.R = clampByte(int(gray + (float64(c.R)-gray)*factor))
		e.Cells[i].Color.G = clampByte(int(gray + (float64(c.G)-gray)*factor))
		e.Cells[i].Color.B = clampByte(int(gray + (float64(c.B)-gray)*factor))
	}
}

func (e *editorState) posterize() {
	for i := range e.Cells {
		if !e.Cells[i].Visible {
			continue
		}
		e.Cells[i].Color.R = posterizeByte(e.Cells[i].Color.R)
		e.Cells[i].Color.G = posterizeByte(e.Cells[i].Color.G)
		e.Cells[i].Color.B = posterizeByte(e.Cells[i].Color.B)
	}
}

func (e *editorState) snapToPalette(palette []color.RGBA) {
	for i := range e.Cells {
		if !e.Cells[i].Visible {
			continue
		}
		e.Cells[i].Color = nearestPaletteColor(e.Cells[i].Color, palette)
	}
}

func (e editorState) pixelMatrix() [][]assets.PixelCell {
	matrix := make([][]assets.PixelCell, e.Height)
	for y := 0; y < e.Height; y++ {
		matrix[y] = make([]assets.PixelCell, e.Width)
		for x := 0; x < e.Width; x++ {
			cell := e.Cells[e.index(x, y)]
			matrix[y][x] = assets.PixelCell{Color: cell.Color, Visible: cell.Visible}
		}
	}
	return matrix
}

func (e editorState) solutionRows() []string {
	rows := make([]string, e.Height)
	for y := 0; y < e.Height; y++ {
		var b strings.Builder
		b.Grow(e.Width)
		for x := 0; x < e.Width; x++ {
			if e.Cells[e.index(x, y)].Filled {
				b.WriteByte('1')
			} else {
				b.WriteByte('0')
			}
		}
		rows[y] = b.String()
	}
	return rows
}

func (e editorState) rawPixels() [][]string {
	rows := make([][]string, e.Height)
	for y := 0; y < e.Height; y++ {
		rows[y] = make([]string, e.Width)
		for x := 0; x < e.Width; x++ {
			cell := e.Cells[e.index(x, y)]
			if !cell.Visible {
				continue
			}
			rows[y][x] = fmt.Sprintf("#%02X%02X%02X%02X", cell.Color.R, cell.Color.G, cell.Color.B, cell.Color.A)
		}
	}
	return rows
}

func (e editorState) puzzle() *nonogram.Puzzle {
	p := &nonogram.Puzzle{
		ID:          "editor_preview",
		Title:       e.Title,
		Width:       e.Width,
		Height:      e.Height,
		SolutionRaw: e.solutionRows(),
		SkeletonRaw: e.rawPixels(),
		RevealRaw:   e.rawPixels(),
	}
	_ = p.ParseSolution()
	return p
}

func (e editorState) packJSON() string {
	pack := editorPack{
		ID:        "local_editor_pack",
		Title:     "My Editor Pack",
		Author:    "local",
		Version:   1,
		UpdatedAt: time.Now().Format(time.RFC3339),
		Levels:    []*nonogram.Puzzle{e.puzzle()},
	}
	data, err := json.MarshalIndent(pack, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (e *editorState) importPayload(raw string) error {
	var payload editorImportPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return err
	}
	if payload.Width <= 0 || payload.Height <= 0 || len(payload.Cells) != payload.Width*payload.Height {
		return fmt.Errorf("invalid imported image payload")
	}
	e.Width = payload.Width
	e.Height = payload.Height
	e.Cells = make([]editorCell, len(payload.Cells))
	for i, cell := range payload.Cells {
		e.Cells[i] = editorCell{
			Color:   color.RGBA{R: cell.R, G: cell.G, B: cell.B, A: cell.A},
			Visible: cell.A > 20,
			Filled:  cell.A > 20,
		}
	}
	return nil
}

func sameRGBA(a, b color.RGBA) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B && a.A == b.A
}

func clampByte(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func posterizeByte(v uint8) uint8 {
	const step = 51
	return uint8(math.Round(float64(v)/step) * step)
}

func nearestPaletteColor(c color.RGBA, palette []color.RGBA) color.RGBA {
	if len(palette) == 0 {
		return c
	}
	best := palette[0]
	bestDist := colorDistance(c, best)
	for _, candidate := range palette[1:] {
		if d := colorDistance(c, candidate); d < bestDist {
			best = candidate
			bestDist = d
		}
	}
	best.A = c.A
	return best
}

func colorDistance(a, b color.RGBA) int {
	dr := int(a.R) - int(b.R)
	dg := int(a.G) - int(b.G)
	db := int(a.B) - int(b.B)
	return dr*dr + dg*dg + db*db
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func parseEditorHexColor(value string) (color.RGBA, bool) {
	if len(value) != 9 || value[0] != '#' {
		return color.RGBA{}, false
	}
	var rgba [4]uint8
	for i := 0; i < 4; i++ {
		n, ok := parseEditorHexByte(value[1+i*2 : 3+i*2])
		if !ok {
			return color.RGBA{}, false
		}
		rgba[i] = n
	}
	return color.RGBA{R: rgba[0], G: rgba[1], B: rgba[2], A: rgba[3]}, true
}

func parseEditorHexByte(value string) (uint8, bool) {
	if len(value) != 2 {
		return 0, false
	}
	var out uint8
	for i := 0; i < 2; i++ {
		ch := value[i]
		var n uint8
		switch {
		case ch >= '0' && ch <= '9':
			n = ch - '0'
		case ch >= 'a' && ch <= 'f':
			n = ch - 'a' + 10
		case ch >= 'A' && ch <= 'F':
			n = ch - 'A' + 10
		default:
			return 0, false
		}
		out = out*16 + n
	}
	return out, true
}

var editorPalette = []color.RGBA{
	{54, 52, 49, 255},
	{255, 252, 240, 255},
	{151, 83, 71, 255},
	{220, 145, 126, 255},
	{86, 115, 134, 255},
	{100, 132, 97, 255},
	{235, 194, 92, 255},
	{98, 76, 128, 255},
}
