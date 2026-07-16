package game

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawEditor(screen *ebiten.Image) {
	screen.Fill(colPanel)
	drawScaledTextCentered(screen, "DRAW", rect{x: 150, y: 22, w: 240, h: 42}, 1.85, colInk)
	drawButton(screen, editorBackButton(), "back")
	drawButton(screen, editorUndoButton(), "undo")

	g.drawEditorGrid(screen)
	g.drawEditorToolbar(screen)
	g.drawEditorPalette(screen)
	g.drawEditorCanvasSizes(screen)
	g.drawEditorActions(screen)

	if time.Now().Before(g.menuNoticeUntil) {
		drawCenteredText(screen, g.menuNotice, rect{x: 120, y: 64, w: 300, h: 24}, colAccent)
	}
}

func (g *Game) drawEditorGrid(screen *ebiten.Image) {
	grid := editorGridRect(g.editor)
	drawRounded(screen, rect{x: grid.x - 8, y: grid.y - 8, w: grid.w + 16, h: grid.h + 16}, 8, colGridHeavy)
	vector.DrawFilledRect(screen, float32(grid.x), float32(grid.y), float32(grid.w), float32(grid.h), colWhite, false)
	cell := editorCellSize(g.editor)
	for y := 0; y < g.editor.Height; y++ {
		for x := 0; x < g.editor.Width; x++ {
			r := rect{x: grid.x + float64(x)*cell, y: grid.y + float64(y)*cell, w: cell, h: cell}
			if (x+y)%2 == 1 {
				vector.DrawFilledRect(screen, float32(r.x), float32(r.y), float32(r.w), float32(r.h), color.RGBA{244, 239, 224, 255}, false)
			}
			c := g.editor.Cells[g.editor.index(x, y)]
			if c.Visible {
				vector.DrawFilledRect(screen, float32(r.x), float32(r.y), float32(r.w), float32(r.h), c.Color, false)
			}
		}
	}
	for x := 0; x <= g.editor.Width; x++ {
		thick := float32(1)
		line := colGrid
		if x%5 == 0 {
			thick = 2
			line = colGridHeavy
		}
		xx := float32(grid.x + float64(x)*cell)
		vector.StrokeLine(screen, xx, float32(grid.y), xx, float32(grid.y+cell*float64(g.editor.Height)), thick, line, false)
	}
	for y := 0; y <= g.editor.Height; y++ {
		thick := float32(1)
		line := colGrid
		if y%5 == 0 {
			thick = 2
			line = colGridHeavy
		}
		yy := float32(grid.y + float64(y)*cell)
		vector.StrokeLine(screen, float32(grid.x), yy, float32(grid.x+cell*float64(g.editor.Width)), yy, thick, line, false)
	}
}

func (g *Game) drawEditorToolbar(screen *ebiten.Image) {
	drawButton(screen, editorPencilButton(), toolLabel("draw", g.editor.Tool == editorToolPencil))
	drawButton(screen, editorEraserButton(), toolLabel("erase", g.editor.Tool == editorToolEraser))
	drawButton(screen, editorFillButton(), toolLabel("fill", g.editor.Tool == editorToolFill))
	drawButton(screen, editorEyeButton(), toolLabel("pick", g.editor.Tool == editorToolEyedropper))
}

func (g *Game) drawEditorPalette(screen *ebiten.Image) {
	for i, c := range editorPalette {
		r := editorPaletteRect(i)
		drawRounded(screen, r, 5, c)
		if sameRGBA(c, g.editor.PaintColor) {
			drawRectOutline(screen, rect{x: r.x - 3, y: r.y - 3, w: r.w + 6, h: r.h + 6}, 3, colInk)
		}
	}
}

func (g *Game) drawEditorCanvasSizes(screen *ebiten.Image) {
	drawText(screen, "canvas", 44, 665, colMuted)
	drawButton(screen, editorSize8Button(), modeLabel("8", g.editor.Width == 8 && g.editor.Height == 8))
	drawButton(screen, editorSize10Button(), modeLabel("10", g.editor.Width == 10 && g.editor.Height == 10))
	drawButton(screen, editorSize15Button(), modeLabel("15", g.editor.Width == 15 && g.editor.Height == 15))
	drawButton(screen, editorSize20Button(), modeLabel("20", g.editor.Width == 20 && g.editor.Height == 20))
}

func (g *Game) drawEditorActions(screen *ebiten.Image) {
	drawButton(screen, editorImportButton(), "image")
	drawButton(screen, editorSaveButton(), "save")
	drawButton(screen, editorExportButton(), "export")
	drawButton(screen, editorPreviewButton(), "play")
}

func (g *Game) drawCommunity(screen *ebiten.Image) {
	drawMenuBackdrop(screen)
	drawScaledTextCentered(screen, "COMMUNITY", rect{x: 76, y: 46, w: 388, h: 52}, 2.1, colInk)
	panel := rect{x: 56, y: 226, w: 428, h: 336}
	drawRounded(screen, panel, 8, colWhite)
	drawRectOutline(screen, panel, 3, colGridHeavy)
	drawCenteredText(screen, "Pack browser scaffold", rect{x: panel.x, y: panel.y + 36, w: panel.w, h: 28}, colInk)
	drawCenteredText(screen, communityFetchStatus(), rect{x: panel.x + 24, y: panel.y + 116, w: panel.w - 48, h: 42}, colMuted)
	drawCenteredText(screen, "Supabase tables: profiles, packs, pack_versions, pack_stats, reports", rect{x: panel.x + 22, y: panel.y + 190, w: panel.w - 44, h: 42}, colMuted)
	drawButton(screen, communityBackButton(), "back")
}

func modeLabel(label string, active bool) string {
	if active {
		return "[" + label + "]"
	}
	return label
}

func toolLabel(label string, active bool) string {
	return modeLabel(label, active)
}

func editorGridArea() rect { return rect{x: 40, y: 98, w: 460, h: 430} }

func editorGridRect(e editorState) rect {
	area := editorGridArea()
	cell := editorCellSize(e)
	w := cell * float64(e.Width)
	h := cell * float64(e.Height)
	return rect{x: area.x + (area.w-w)/2, y: area.y + (area.h-h)/2, w: w, h: h}
}

func editorCellSize(e editorState) float64 {
	area := editorGridArea()
	return float64(int(minFloat(area.w/float64(e.Width), area.h/float64(e.Height))))
}

func editorCellAt(e editorState, px, py int) (int, int, bool) {
	grid := editorGridRect(e)
	if float64(px) < grid.x || float64(px) >= grid.x+grid.w || float64(py) < grid.y || float64(py) >= grid.y+grid.h {
		return 0, 0, false
	}
	cell := editorCellSize(e)
	x := int((float64(px) - grid.x) / cell)
	y := int((float64(py) - grid.y) / cell)
	return x, y, e.inBounds(x, y)
}

func editorBackButton() rect    { return rect{x: 40, y: 24, w: 82, h: 38} }
func editorUndoButton() rect    { return rect{x: 418, y: 24, w: 82, h: 38} }
func editorPencilButton() rect  { return rect{x: 40, y: 546, w: 105, h: 38} }
func editorEraserButton() rect  { return rect{x: 155, y: 546, w: 105, h: 38} }
func editorFillButton() rect    { return rect{x: 270, y: 546, w: 105, h: 38} }
func editorEyeButton() rect     { return rect{x: 385, y: 546, w: 105, h: 38} }
func editorSize8Button() rect   { return rect{x: 142, y: 651, w: 46, h: 32} }
func editorSize10Button() rect  { return rect{x: 198, y: 651, w: 50, h: 32} }
func editorSize15Button() rect  { return rect{x: 258, y: 651, w: 50, h: 32} }
func editorSize20Button() rect  { return rect{x: 318, y: 651, w: 50, h: 32} }
func editorImportButton() rect  { return rect{x: 40, y: 700, w: 105, h: 38} }
func editorSaveButton() rect    { return rect{x: 155, y: 700, w: 105, h: 38} }
func editorExportButton() rect  { return rect{x: 270, y: 700, w: 105, h: 38} }
func editorPreviewButton() rect { return rect{x: 385, y: 700, w: 105, h: 38} }
func communityBackButton() rect { return rect{x: 202, y: 650, w: 136, h: 42} }

func editorPaletteRect(index int) rect {
	return rect{x: 48 + float64(index)*56, y: 603, w: 36, h: 36}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
