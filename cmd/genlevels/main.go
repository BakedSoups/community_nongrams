package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alex/nongrampictures/internal/pixelpuzzle"
	_ "golang.org/x/image/webp"
)

var (
	levelFilePattern   = regexp.MustCompile(`^(L[0-9]+)[-_](.+)_([0-9]+)(?:\.(?:png|webp))?$`)
	levelFolderPattern = regexp.MustCompile(`^(?:[Ll])?0*([0-9]+)[-_ ]+(.+)$`)
)

type levelSource struct {
	ID       string
	Title    string
	Source   string
	TileSize int
	Origin   string
}

func main() {
	levelsDir := flag.String("levels", "levels", "folder containing L1-name_16 spritesheet files")
	outRoot := flag.String("out", "assets/puzzles", "output puzzle root")
	embedRoot := flag.String("embed", "internal/assets/embedded/assets/puzzles", "embedded puzzle root to refresh")
	alphaThreshold := flag.Uint("alpha-threshold", 128, "alpha threshold for filled pixels, 0-255")
	useBackground := flag.Bool("background-empty", true, "when image is opaque, treat the top-left color as empty")
	flag.Parse()

	entries, err := os.ReadDir(*levelsDir)
	if err != nil {
		fatal(err)
	}

	sources, err := findLevelSources(*levelsDir, entries)
	if err != nil {
		fatal(err)
	}
	if len(sources) == 0 {
		fatal(fmt.Errorf("no level spritesheets found in %s", *levelsDir))
	}

	sources = dedupeLevelSources(sources)

	for _, level := range sources {
		out := filepath.Join(*outRoot, level.ID)

		puzzle, err := pixelpuzzle.GenerateSpriteSheet(pixelpuzzle.SpriteSheetOptions{
			ID:             level.ID,
			Title:          level.Title,
			Source:         level.Source,
			Out:            out,
			TileSize:       level.TileSize,
			AlphaThreshold: uint32(*alphaThreshold),
			UseBackground:  *useBackground,
		})
		if err != nil {
			fatal(err)
		}

		embedOut := filepath.Join(*embedRoot, level.ID)
		if err := os.MkdirAll(embedOut, 0o755); err != nil {
			fatal(err)
		}
		if err := pixelpuzzle.CopyFile(filepath.Join(out, "puzzle.json"), filepath.Join(embedOut, "puzzle.json")); err != nil {
			fatal(err)
		}

		fmt.Printf("generated %s from %s (%dx%d, json only)\n", filepath.Join(out, "puzzle.json"), level.Origin, puzzle.Width, puzzle.Height)
	}
}

func dedupeLevelSources(sources []levelSource) []levelSource {
	positions := map[string]int{}
	result := make([]levelSource, 0, len(sources))
	for _, source := range sources {
		if pos, ok := positions[source.ID]; ok {
			fmt.Fprintf(os.Stderr, "warning: duplicate level %s: using %s instead of %s\n", source.ID, source.Origin, result[pos].Origin)
			result[pos] = source
			continue
		}
		positions[source.ID] = len(result)
		result = append(result, source)
	}
	return result
}

func findLevelSources(levelsDir string, entries []os.DirEntry) ([]levelSource, error) {
	sources := make([]levelSource, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			source, ok, err := folderLevelSource(levelsDir, entry.Name())
			if err != nil {
				return nil, err
			}
			if ok {
				sources = append(sources, source)
			}
			continue
		}
		matches := levelFilePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}

		levelName := matches[1]
		artName := matches[2]
		tileSize, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, err
		}
		id := strings.ToLower(levelName)
		source := filepath.Join(levelsDir, entry.Name())
		sources = append(sources, levelSource{
			ID:       id,
			Title:    levelTitle(levelName, artName),
			Source:   source,
			TileSize: tileSize,
			Origin:   entry.Name(),
		})
	}
	return sources, nil
}

func folderLevelSource(levelsDir, name string) (levelSource, bool, error) {
	matches := levelFolderPattern.FindStringSubmatch(name)
	if matches == nil {
		return levelSource{}, false, nil
	}

	number, err := strconv.Atoi(matches[1])
	if err != nil {
		return levelSource{}, false, err
	}
	imagePath, err := levelFolderImage(filepath.Join(levelsDir, name))
	if err != nil {
		return levelSource{}, false, fmt.Errorf("%s: %w", name, err)
	}
	tileSize, err := inferredTileSize(imagePath)
	if err != nil {
		return levelSource{}, false, fmt.Errorf("%s: %w", name, err)
	}

	levelName := fmt.Sprintf("L%d", number)
	return levelSource{
		ID:       fmt.Sprintf("l%d", number),
		Title:    levelTitle(levelName, matches[2]),
		Source:   imagePath,
		TileSize: tileSize,
		Origin:   filepath.ToSlash(filepath.Join(name, filepath.Base(imagePath))),
	}, true, nil
}

func levelFolderImage(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var images []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !isLevelImage(entry.Name()) {
			continue
		}
		images = append(images, entry.Name())
	}
	if len(images) == 0 {
		return "", fmt.Errorf("no png or webp image found")
	}

	preferred := []string{"art.png", "art.webp", "sheet.png", "sheet.webp"}
	for _, want := range preferred {
		for _, imageName := range images {
			if strings.EqualFold(imageName, want) {
				return filepath.Join(dir, imageName), nil
			}
		}
	}
	if len(images) > 1 {
		return "", fmt.Errorf("found multiple images; name the source art.png, art.webp, sheet.png, or sheet.webp")
	}
	return filepath.Join(dir, images[0]), nil
}

func isLevelImage(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".png", ".webp":
		return true
	default:
		return false
	}
}

func inferredTileSize(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, err
	}
	if cfg.Width < cfg.Height*2 {
		return 0, fmt.Errorf("%s is %dx%d: expected a two-panel sheet", filepath.Base(path), cfg.Width, cfg.Height)
	}
	return cfg.Height, nil
}

func levelTitle(levelName, artName string) string {
	words := strings.FieldsFunc(artName, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}
	return levelName + " " + strings.Join(words, " ")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
