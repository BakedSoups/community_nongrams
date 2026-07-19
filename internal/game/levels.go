package game

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/BakedSoups/community_nongrams/internal/assets"
)

type levelInfo struct {
	ID        string
	Label     string
	Path      string
	Available bool
}

var (
	levelIDPattern = regexp.MustCompile(`^l([0-9]+)$`)
	gameLevels     = buildGameLevels()
)

func buildGameLevels() []levelInfo {
	available := map[int]string{}
	maxLevel := 0
	for _, puzzlePath := range assets.ListPuzzlePaths() {
		levelNumber, ok := levelNumberFromPath(puzzlePath)
		if !ok {
			continue
		}
		available[levelNumber] = puzzlePath
		if levelNumber > maxLevel {
			maxLevel = levelNumber
		}
	}
	if maxLevel == 0 {
		maxLevel = 1
	}

	levels := make([]levelInfo, 0, maxLevel)
	for number := 1; number <= maxLevel; number++ {
		id := fmt.Sprintf("l%d", number)
		puzzlePath, ok := available[number]
		levels = append(levels, levelInfo{
			ID:        id,
			Label:     strings.ToUpper(id),
			Path:      puzzlePath,
			Available: ok,
		})
	}
	return levels
}

func levelNumberFromPath(puzzlePath string) (int, bool) {
	id := path.Base(path.Dir(filepathSlash(puzzlePath)))
	matches := levelIDPattern.FindStringSubmatch(id)
	if matches == nil {
		return 0, false
	}
	number, err := strconv.Atoi(matches[1])
	return number, err == nil
}

func filepathSlash(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}
