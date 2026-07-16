package game

import "testing"

func TestEditorPaintingBuildsSolutionAutomatically(t *testing.T) {
	editor := newEditorState(4)
	editor.apply(1, 2)

	cell := editor.Cells[editor.index(1, 2)]
	if !cell.Visible || !cell.Filled {
		t.Fatalf("painted cell = %+v, want visible and filled", cell)
	}

	editor.Tool = editorToolEraser
	editor.apply(1, 2)
	cell = editor.Cells[editor.index(1, 2)]
	if cell.Visible || cell.Filled {
		t.Fatalf("erased cell = %+v, want hidden and unfilled", cell)
	}
}

func TestEditorLineConnectsFastPointerMovement(t *testing.T) {
	editor := newEditorState(6)
	editor.applyLine(0, 0, 5, 5)

	for i := 0; i < 6; i++ {
		cell := editor.Cells[editor.index(i, i)]
		if !cell.Visible || !cell.Filled {
			t.Fatalf("line cell (%d,%d) = %+v, want visible and filled", i, i, cell)
		}
	}
}
