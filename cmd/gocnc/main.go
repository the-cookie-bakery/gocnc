package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"gocnc/engine"
)

func main() {
	game := engine.NewGame()

	if err := game.LoadSprites(); err != nil {
		log.Fatal("Failed to load sprites:", err)
	}

	game.SetupDemoLevel()

	if err := game.LoadAudio(); err != nil {
		log.Println("Audio not available:", err)
	}

	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("GoCNC")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
