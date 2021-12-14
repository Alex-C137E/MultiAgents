package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	music "gitlab.utc.fr/projet_ia04/Boid/agent/music"
	game "gitlab.utc.fr/projet_ia04/Boid/game"
	constant "gitlab.utc.fr/projet_ia04/Boid/utils/constant"
	variable "gitlab.utc.fr/projet_ia04/Boid/utils/variable"
	"log"
)

func init() {
	fish, _, err := ebitenutil.NewImageFromFile("utils/fish/chevron-up.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	w, h := fish.Size()
	variable.BirdImage, _ = ebiten.NewImage(w-20, h-10, 0)
	variable.WallImage, _ = ebiten.NewImage(w, h, 0)
	op := &ebiten.DrawImageOptions{}
	variable.BirdImage.DrawImage(fish, op)
	variable.WallImage.DrawImage(fish, op)
}

func main() {
	// création de la chanel de sync pour la music:
	c := make(chan string)
	musicAgent := music.NewMusicAgent("utils/music/virtual-riot-the-darkest-night.mp3", "utils/music/virtual-riot-the-darkest-night.wav", c)
	musicAgent.Start()

	ebiten.SetWindowSize(constant.ScreenWidth, constant.ScreenHeight)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(game.NewGame(c)); err != nil {
		log.Fatal(err)
	}
}
