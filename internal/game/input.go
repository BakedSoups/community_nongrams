package game

import (
	"image"
	"time"

	"github.com/alex/nongrampictures/internal/nonogram"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) updateInput() {
	if g.mode == screenMainMenu {
		g.updateMainMenuInput()
		return
	}
	if g.mode == screenLevelSelect {
		g.updateLevelSelectInput()
		return
	}
	if g.mode == screenSettings {
		g.updateSettingsInput()
		return
	}
	if g.mode == screenEditor {
		g.updateEditorInput()
		return
	}
	if g.mode == screenCommunity {
		g.updateCommunityInput()
		return
	}
	if g.mode == screenReveal {
		g.updateRevealInput()
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.tool = nonogram.ToolFill
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.tool = nonogram.ToolMark
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = screenMainMenu
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.godModeFill()
		return
	}

	x, y, down, justPressed, justReleased := pointerState()
	if justReleased {
		g.pointerDown = false
		g.dragging = false
		g.lastCellX = -1
		g.lastCellY = -1
		g.strokeState = nonogram.CellEmpty
	}
	if !down {
		return
	}

	if justPressed {
		switch {
		case g.layout.fillTrigger.Contains(x, y):
			g.tool = nonogram.ToolFill
			return
		case g.layout.markTrigger.Contains(x, y):
			g.tool = nonogram.ToolMark
			return
		case g.layout.godModeButton.Contains(x, y):
			g.godModeFill()
			return
		case g.layout.menuButton.Contains(x, y):
			g.mode = screenMainMenu
			return
		case g.layout.settingsButton.Contains(x, y):
			g.mode = screenSettings
			return
		}
	}

	cellX, cellY, ok := g.layout.CellAt(x, y, g.board.Width, g.board.Height)
	if !ok {
		return
	}
	if justPressed {
		g.pushUndo()
		g.pointerDown = true
		g.strokeState = nonogram.TargetState(g.tool)
		if g.board.Cells[cellY][cellX] == g.strokeState {
			g.strokeState = nonogram.CellEmpty
		}
	}
	if !g.pointerDown && !g.dragging {
		return
	}
	if cellX == g.lastCellX && cellY == g.lastCellY {
		return
	}

	next, corrected := g.correctedStrokeState(cellX, cellY, g.strokeState)
	if g.board.SetCell(cellX, cellY, next) {
		if corrected {
			g.timePenalty += 10 * time.Second
			g.penaltyFlashUntil = time.Now().Add(900 * time.Millisecond)
			g.correctFlashUntil = time.Now().Add(850 * time.Millisecond)
			g.correctFlashX = cellX
			g.correctFlashY = cellY
			playWebSFX("correct")
		} else if next == nonogram.CellFilled {
			playWebSFX("pencil")
		} else if next == nonogram.CellMarked || next == nonogram.CellEmpty {
			playWebSFX("eraser")
		}
		if nonogram.IsSolved(g.board, g.puzzle.Solution) {
			g.completePuzzle()
		}
	}
	g.dragging = true
	g.lastCellX = cellX
	g.lastCellY = cellY
}

func (g *Game) correctedStrokeState(cellX, cellY int, attempted nonogram.CellState) (nonogram.CellState, bool) {
	if !g.autoCorrect || attempted == nonogram.CellEmpty {
		return attempted, false
	}
	if attempted == nonogram.CellFilled && !g.puzzle.Solution[cellY][cellX] {
		return nonogram.CellMarked, true
	}
	if attempted == nonogram.CellMarked && g.puzzle.Solution[cellY][cellX] {
		return nonogram.CellFilled, true
	}
	return attempted, false
}

func (g *Game) updateRevealInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.retry()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = screenLevelSelect
		return
	}

	x, y, _, justPressed, _ := pointerState()
	if !justPressed {
		return
	}
	if g.layout.retryButton.Contains(x, y) {
		g.retry()
		return
	}
	if g.layout.revealLevelsButton.Contains(x, y) {
		g.mode = screenLevelSelect
	}
}

func (g *Game) updateMainMenuInput() {
	x, y, _, justPressed, _ := pointerState()
	if !justPressed {
		return
	}
	switch {
	case mainLevelButton().Contains(x, y):
		g.mode = screenLevelSelect
	case mainEditorButton().Contains(x, y):
		g.mode = screenEditor
	case mainCommunityButton().Contains(x, y):
		g.mode = screenCommunity
	case mainSettingsButton().Contains(x, y):
		g.mode = screenSettings
	}
}

func (g *Game) updateCommunityInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = screenMainMenu
		return
	}
	x, y, _, justPressed, _ := pointerState()
	if justPressed && communityBackButton().Contains(x, y) {
		g.mode = screenMainMenu
	}
}

func (g *Game) updateEditorInput() {
	if raw := takeEditorImageImport(); raw != "" {
		g.pushEditorUndo()
		if err := g.editor.importPayload(raw); err != nil {
			g.showMenuNotice("import failed")
		} else {
			g.showMenuNotice("image imported")
		}
	}
	if raw := takeEditorPackImport(); raw != "" {
		g.pushEditorUndo()
		if editor, err := editorFromPackJSON(raw); err != nil {
			g.showMenuNotice("pack failed")
		} else {
			g.editor = editor
			g.showMenuNotice("pack imported")
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = screenMainMenu
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.undoEditor()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.editor.Mode = editorModeArt
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.editor.Mode = editorModeSolution
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.editor.Tool = editorToolPencil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.editor.Tool = editorToolEraser
	}

	x, y, down, justPressed, justReleased := pointerState()
	if justReleased {
		g.editorPointer = false
		g.editorLastX = -1
		g.editorLastY = -1
	}
	if justPressed {
		if g.handleEditorButton(x, y) {
			return
		}
		for i, c := range editorPalette {
			if editorPaletteRect(i).Contains(x, y) {
				g.editor.PaintColor = c
				return
			}
		}
		cellX, cellY, ok := editorCellAt(g.editor, x, y)
		if ok {
			g.pushEditorUndo()
			g.editor.apply(cellX, cellY)
			g.editorPointer = true
			g.editorLastX = cellX
			g.editorLastY = cellY
			return
		}
	}
	if !down || !g.editorPointer || g.editor.Tool == editorToolFill || g.editor.Tool == editorToolEyedropper {
		return
	}
	cellX, cellY, ok := editorCellAt(g.editor, x, y)
	if !ok || (cellX == g.editorLastX && cellY == g.editorLastY) {
		return
	}
	g.editor.apply(cellX, cellY)
	g.editorLastX = cellX
	g.editorLastY = cellY
}

func (g *Game) handleEditorButton(x, y int) bool {
	switch {
	case editorBackButton().Contains(x, y):
		g.mode = screenMainMenu
	case editorPreviewButton().Contains(x, y):
		g.loadEditorPuzzle()
	case editorSaveButton().Contains(x, y):
		g.saveEditor()
	case editorExportButton().Contains(x, y):
		g.exportEditor()
	case editorImportPackButton().Contains(x, y):
		if !requestEditorPackImport() {
			g.showMenuNotice("import unavailable")
		}
	case editorArtButton().Contains(x, y):
		g.editor.Mode = editorModeArt
	case editorSolutionButton().Contains(x, y):
		g.editor.Mode = editorModeSolution
	case editorPencilButton().Contains(x, y):
		g.editor.Tool = editorToolPencil
	case editorEraserButton().Contains(x, y):
		g.editor.Tool = editorToolEraser
	case editorFillButton().Contains(x, y):
		g.editor.Tool = editorToolFill
	case editorEyeButton().Contains(x, y):
		g.editor.Tool = editorToolEyedropper
	case editorAutoVisibleButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.autoSolutionFromVisible()
	case editorAutoBrightButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.autoSolutionFromBrightness()
	case editorInvertButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.invertSolution()
	case editorImportButton().Contains(x, y):
		if !requestEditorImageImport(g.editor.Width) {
			g.showMenuNotice("import unavailable")
		}
	case editorSize8Button().Contains(x, y):
		g.resetEditor(8)
	case editorSize10Button().Contains(x, y):
		g.resetEditor(10)
	case editorSize15Button().Contains(x, y):
		g.resetEditor(15)
	case editorSize20Button().Contains(x, y):
		g.resetEditor(20)
	case editorBrightDownButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.applyBrightness(-18)
	case editorBrightUpButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.applyBrightness(18)
	case editorSatDownButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.applySaturation(-0.18)
	case editorSatUpButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.applySaturation(0.18)
	case editorPosterizeButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.posterize()
	case editorSnapButton().Contains(x, y):
		g.pushEditorUndo()
		g.editor.snapToPalette(editorPalette)
	default:
		return false
	}
	return true
}

func (g *Game) updateLevelSelectInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = screenMainMenu
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.prevLevelPage()
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.nextLevelPage()
		return
	}

	x, y, _, justPressed, _ := pointerState()
	if !justPressed {
		return
	}
	if g.layout.levelBackButton.Contains(x, y) {
		g.mode = screenMainMenu
		return
	}
	if g.layout.levelPrevButton.Contains(x, y) {
		g.prevLevelPage()
		return
	}
	if g.layout.levelNextButton.Contains(x, y) {
		g.nextLevelPage()
		return
	}
	pageStart := g.levelPage * levelSelectPageSize
	for slot := 0; slot < levelSelectPageSize; slot++ {
		if levelTileRect(slot).Contains(x, y) {
			levelIndex := pageStart + slot
			if levelIndex < len(gameLevels) && gameLevels[levelIndex].Available {
				_ = g.loadLevel(levelIndex)
			} else {
				g.showMenuNotice("LW")
			}
			return
		}
	}
}

func (g *Game) updateSettingsInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.mode = screenMainMenu
		return
	}

	x, y, _, justPressed, _ := pointerState()
	if !justPressed {
		return
	}
	switch {
	case g.layout.soundButton.Contains(x, y):
		g.audioEnabled = !g.audioEnabled
		setWebMusicMuted(!g.audioEnabled)
	case g.layout.autoCorrectButton.Contains(x, y):
		g.autoCorrect = !g.autoCorrect
	case g.layout.settingsCloseButton.Contains(x, y):
		g.mode = screenMainMenu
	}
}

func pointerState() (int, int, bool, bool, bool) {
	x, y := ebiten.CursorPosition()
	down := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	justPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	justReleased := inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)

	touches := ebiten.AppendTouchIDs(nil)
	if len(touches) > 0 {
		tx, ty := ebiten.TouchPosition(touches[0])
		x, y = tx, ty
		down = true
		justPressed = inpututil.IsTouchJustReleased(touches[0]) == false && inpututil.TouchPressDuration(touches[0]) == 1
		justReleased = false
	}
	return x, y, down, justPressed, justReleased
}

type rect struct {
	x float64
	y float64
	w float64
	h float64
}

func (r rect) Contains(px, py int) bool {
	return float64(px) >= r.x && float64(px) <= r.x+r.w && float64(py) >= r.y && float64(py) <= r.y+r.h
}

func (r rect) ImageRect() image.Rectangle {
	return image.Rect(int(r.x), int(r.y), int(r.x+r.w), int(r.y+r.h))
}
