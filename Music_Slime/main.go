//go:build example
// +build example

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	game "gitlab.utc.fr/projet_ia04/musicslime/game"
	music "gitlab.utc.fr/projet_ia04/musicslime/music"
	types "gitlab.utc.fr/projet_ia04/musicslime/types"
	"log"
)

func main() {
	// création de la chanel de sync music-slime:
	c := make(chan string)
	musicAgent := music.NewMusicAgent("virtual-riot-the-darkest-night.mp3", "virtual-riot-the-darkest-night.wav", c)
	musicAgent.Start()
	ebiten.SetWindowSize(int(types.GetWindowDefault().Width), int(types.GetWindowDefault().Height))
	ebiten.SetWindowTitle("Music Slime Demo")
	if err := ebiten.RunGame(game.NewGame(c)); err != nil {
		log.Fatal(err)
	}
}
