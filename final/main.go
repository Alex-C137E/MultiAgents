package main

import (
	"fmt"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	game "gitlab.utc.fr/jedescam/projet_ia04/final/game"
	settings "gitlab.utc.fr/jedescam/projet_ia04/final/settings"
)

func main() {
	fmt.Println("Start of the program")

	ebiten.SetWindowSize(int(settings.GetWindowDefault().Width), int(settings.GetWindowDefault().Height))
	ebiten.SetWindowTitle("Final Project IA04")
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
