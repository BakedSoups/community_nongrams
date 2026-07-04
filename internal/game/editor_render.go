package game

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawEditor(screen *ebiten.Image) {
	screen.Fill(colPanel)
	drawScaledTextCentered(screen, "EDITOR", rect{x: 36, y: 28, w: 220, h: 44}, 1.85, colInk)
	drawButton(screen, editorBackButton(), "back")
	drawButton(screen, editorPreviewButton(), "preview")
	drawButton(screen, editorSaveButton(), "save")
	drawButton(screen, editorExportButton(), "export")
	drawButton(screen, editorImportPackButton(), "import pack")

	g.drawEditorGrid(screen)
	g.drawEditorToolbar(screen)
	g.drawEditorPalette(screen)
	g.drawEditorColorPanel(screen)

	if time.Now().Before(g.menuNoticeUntil) {
		drawCenteredText(screen, g.menuNotice, rect{x: 0, y: 736, w: ScreenWidth, h: 28}, colAccent)
	}
}

func (g *Game) drawEditorGrid(screen *ebiten.Image) {
	grid := editorGridRect()
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
			if g.editor.Mode == editorModeSolution && c.Filled {
				drawRectOutline(screen, inset(r, 4), 3, colBlue)
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
	drawButton(screen, editorArtButton(), modeLabel("art", g.editor.Mode == editorModeArt))
	drawButton(screen, editorSolutionButton(), modeLabel("solution", g.editor.Mode == editorModeSolution))
	drawButton(screen, editorPencilButton(), toolLabel("pencil", g.editor.Tool == editorToolPencil))
	drawButton(screen, editorEraserButton(), toolLabel("eraser", g.editor.Tool == editorToolEraser))
	drawButton(screen, editorFillButton(), toolLabel("fill", g.editor.Tool == editorToolFill))
	drawButton(screen, editorEyeButton(), toolLabel("eye", g.editor.Tool == editorToolEyedropper))
	drawButton(screen, editorAutoVisibleButton(), "auto visible")
	drawButton(screen, editorAutoBrightButton(), "auto dark")
	drawButton(screen, editorInvertButton(), "invert")
	drawButton(screen, editorImportButton(), "import image")
	drawButton(screen, editorSize8Button(), "8")
	drawButton(screen, editorSize10Button(), "10")
	drawButton(screen, editorSize15Button(), "15")
	drawButton(screen, editorSize20Button(), "20")
}

func (g *Game) drawEditorPalette(screen *ebiten.Image) {
	for i, c := range editorPalette {
		r := editorPaletteRect(i)
		drawRounded(screen, r, 5, c)
		if sameRGBA(c, g.editor.PaintColor) {
			drawRectOutline(screen, rect{x: r.x - 3, y: r.y - 3, w: r.w + 6, h: r.h + 6}, 3, colInk)
		}
	}
	drawText(screen, "paint", 48, 610, colMuted)
	drawRounded(screen, rect{x: 48, y: 626, w: 64, h: 34}, 5, g.editor.PaintColor)
	drawRectOutline(screen, rect{x: 48, y: 626, w: 64, h: 34}, 2, colGridHeavy)
}

func (g *Game) drawEditorColorPanel(screen *ebiten.Image) {
	drawText(screen, fmt.Sprintf("%dx%d", g.editor.Width, g.editor.Height), 48, 686, colInk)
	drawText(screen, strings.ToUpper(editorModeName(g.editor.Mode)), 120, 686, colInk)
	drawButton(screen, editorBrightDownButton(), "bright -")
	drawButton(screen, editorBrightUpButton(), "bright +")
	drawButton(screen, editorSatDownButton(), "sat -")
	drawButton(screen, editorSatUpButton(), "sat +")
	drawButton(screen, editorPosterizeButton(), "posterize")
	drawButton(screen, editorSnapButton(), "snap")
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

func editorModeName(mode editorMode) string {
	if mode == editorModeSolution {
		return "solution"
	}
	return "art"
}

func editorGridRect() rect { return rect{x: 48, y: 128, w: 360, h: 360} }

func editorCellSize(e editorState) float64 {
	size := float64(e.Width)
	if e.Height > e.Width {
		size = float64(e.Height)
	}
	return editorGridRect().w / size
}

func editorCellAt(e editorState, px, py int) (int, int, bool) {
	grid := editorGridRect()
	cell := editorCellSize(e)
	x := int((float64(px) - grid.x) / cell)
	y := int((float64(py) - grid.y) / cell)
	return x, y, e.inBounds(x, y)
}

func editorBackButton() rect        { return rect{x: 424, y: 28, w: 82, h: 38} }
func editorPreviewButton() rect     { return rect{x: 424, y: 78, w: 82, h: 38} }
func editorSaveButton() rect        { return rect{x: 424, y: 128, w: 82, h: 38} }
func editorExportButton() rect      { return rect{x: 424, y: 178, w: 82, h: 38} }
func editorImportPackButton() rect  { return rect{x: 424, y: 228, w: 82, h: 38} }
func editorArtButton() rect         { return rect{x: 424, y: 286, w: 82, h: 34} }
func editorSolutionButton() rect    { return rect{x: 424, y: 328, w: 82, h: 34} }
func editorPencilButton() rect      { return rect{x: 424, y: 378, w: 82, h: 34} }
func editorEraserButton() rect      { return rect{x: 424, y: 420, w: 82, h: 34} }
func editorFillButton() rect        { return rect{x: 424, y: 462, w: 82, h: 34} }
func editorEyeButton() rect         { return rect{x: 424, y: 504, w: 82, h: 34} }
func editorAutoVisibleButton() rect { return rect{x: 48, y: 510, w: 108, h: 34} }
func editorAutoBrightButton() rect  { return rect{x: 166, y: 510, w: 92, h: 34} }
func editorInvertButton() rect      { return rect{x: 268, y: 510, w: 84, h: 34} }
func editorImportButton() rect      { return rect{x: 48, y: 704, w: 144, h: 34} }
func editorSize8Button() rect       { return rect{x: 424, y: 566, w: 36, h: 32} }
func editorSize10Button() rect      { return rect{x: 470, y: 566, w: 36, h: 32} }
func editorSize15Button() rect      { return rect{x: 424, y: 606, w: 36, h: 32} }
func editorSize20Button() rect      { return rect{x: 470, y: 606, w: 36, h: 32} }
func editorBrightDownButton() rect  { return rect{x: 198, y: 624, w: 86, h: 32} }
func editorBrightUpButton() rect    { return rect{x: 294, y: 624, w: 86, h: 32} }
func editorSatDownButton() rect     { return rect{x: 198, y: 664, w: 86, h: 32} }
func editorSatUpButton() rect       { return rect{x: 294, y: 664, w: 86, h: 32} }
func editorPosterizeButton() rect   { return rect{x: 392, y: 664, w: 114, h: 32} }
func editorSnapButton() rect        { return rect{x: 392, y: 624, w: 114, h: 32} }
func communityBackButton() rect     { return rect{x: 202, y: 650, w: 136, h: 42} }

func editorPaletteRect(index int) rect {
	return rect{x: 48 + float64(index%4)*44, y: 556 + float64(index/4)*44, w: 32, h: 32}
}
