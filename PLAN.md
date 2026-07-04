# Editor And Community Plan

## Goal

Add an in-game editor where players can create nonogram picture levels by drawing directly, importing photos, color-correcting them, saving them into packs, and eventually sharing or downloading packs through a community browser backed by Supabase.

## Core Model

Separate artwork from the puzzle solution.

```go
type EditorCell struct {
	Color   color.RGBA
	Visible bool
	Filled  bool
}
```

- `Color` and `Visible` define the pixel art reveal.
- `Filled` defines the actual nonogram solution.
- This lets imported or hand-drawn art look good without forcing every visible pixel to be part of the puzzle unless desired.

## Phase 1: Local Editor

Add an `Editor` button on the main menu.

Editor features:

- New puzzle size presets: `8x8`, `10x10`, `15x15`, `20x20`.
- Pixel grid canvas.
- Pencil tool.
- Eraser tool.
- Fill bucket.
- Eyedropper.
- Palette.
- Undo and redo.
- Clear and reset.
- Grid toggle.
- Preview level.

Editor modes:

- `Art`: draw the reveal image.
- `Solution`: mark which cells are filled.
- `Preview`: test-play the puzzle.

## Phase 2: Solution Tools

Add helpers to turn artwork into a puzzle.

Features:

- Generate solution from visible pixels.
- Generate solution from brightness.
- Invert solution toggle.
- Manual solution editing.
- Warning if solution is empty or fully filled.
- Optional later: solvability check.

## Phase 3: Color Corrector

Add a color correction panel for drawn or imported art.

Controls:

- Brightness.
- Contrast.
- Saturation.
- Hue shift.
- Posterize or reduce color count.
- Alpha threshold.
- Replace selected color.
- Snap to palette.
- Preserve transparent cells toggle.

Apply options:

- Apply to whole image.
- Apply to selected color.
- Apply to visible cells only.

## Phase 4: Photo Import

For web builds, add browser file import.

Flow:

1. User selects image.
2. Crop or fit to square.
3. Choose grid size.
4. Downsample image into pixel cells.
5. Apply color correction.
6. Auto-generate solution.
7. Let user hand-edit art and solution.

Photo import controls:

- Crop position.
- Zoom.
- Brightness threshold.
- Alpha or background threshold.
- Invert.
- Color count.
- Preview before accepting.

## Phase 5: Local Packs

Add `My Packs`.

Pack features:

- Create pack.
- Rename pack.
- Add editor levels to pack.
- Rename levels.
- Reorder levels.
- Delete levels.
- Save locally in browser storage.
- Export pack as JSON.
- Import pack JSON.

Suggested pack format:

```json
{
  "id": "cute_food",
  "title": "Cute Food",
  "author": "alex",
  "version": 1,
  "levels": [
    {
      "id": "l1",
      "title": "Strawberry",
      "width": 10,
      "height": 10,
      "solution": [],
      "skeletonPixels": [],
      "revealPixels": []
    }
  ]
}
```

## Phase 6: Runtime Pack Loading

Refactor level loading so built-in levels, local packs, imported packs, and community packs all use the same path.

Use a common source interface:

```go
type LevelSource interface {
	ListPacks() []PackInfo
	LoadPuzzle(packID string, levelID string) (*Puzzle, error)
}
```

Sources:

- Built-in levels.
- Local saved packs.
- Imported JSON packs.
- Downloaded community packs.
- Supabase community packs.

## Phase 7: Community Browser

Add a `Community` button on the main menu.

Community views:

- Featured.
- New.
- Popular.
- Search.
- My Uploads.

Pack cards show:

- Cover thumbnail.
- Pack title.
- Author.
- Level count.
- Likes.
- Downloads.
- Play or download button.

Downloaded packs should cache locally.

## Phase 8: Supabase

Use Supabase for auth, pack metadata, uploaded pack JSON, stats, and reports.

Tables:

- `profiles`
- `packs`
- `pack_versions`
- `pack_stats`
- `reports`

Storage buckets:

- `pack-covers`
- Optional `source-images`.

Pack status:

- `draft`
- `published`
- `hidden`
- `removed`

## Phase 9: Upload Flow

From editor:

1. Validate pack.
2. Require sign-in.
3. Upload pack JSON.
4. Generate and upload cover thumbnail.
5. Publish pack.
6. Show it in Community.

Validation:

- Max pack size.
- Max level count.
- Valid puzzle dimensions.
- Non-empty solution.
- No fully-filled solution unless allowed.
- Safe title and description length.

## Phase 10: Moderation

Minimum moderation:

- Report button.
- Admin hide and remove.
- Upload rate limits.
- Max JSON and image size.
- Basic profanity filter.
- No external image URLs in public packs.

## Recommended Build Order

1. Local editor grid and drawing tools.
2. Art and solution separation.
3. Preview-play editor puzzle.
4. Local save, export, and import packs.
5. Photo import.
6. Color corrector.
7. Dynamic pack loading.
8. Community browser read-only.
9. Supabase upload and login.
10. Likes, reports, moderation, and polish.
