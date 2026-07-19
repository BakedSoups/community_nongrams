package main

import (
	"log"

	"github.com/BakedSoups/community_nongrams/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g, err := game.New("assets/puzzles/l1/puzzle.json")
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowTitle("Community Nongrams")
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
