package main

import (
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"

	music "gitlab.utc.fr/projet_ia04/Boid/agent/music"
	game "gitlab.utc.fr/projet_ia04/Boid/game"
	constant "gitlab.utc.fr/projet_ia04/Boid/utils/constant"
	variable "gitlab.utc.fr/projet_ia04/Boid/utils/variable"
)

func init() {
	fish, _, err := ebitenutil.NewImageFromFile("utils/fish/chevron-up.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h := fish.Size()
	variable.BirdImage = ebiten.NewImage(w-20, h-10)
	op := &ebiten.DrawImageOptions{}
	variable.BirdImage.DrawImage(fish, op)

	bomb, _, err := ebitenutil.NewImageFromFile("utils/fish/bomb.png")
	if err != nil {
		log.Fatal(err)
	}
	bombW, bombH := 24, 24
	variable.WallImage = ebiten.NewImage(bombW, bombH)
	variable.WallImage.DrawImage(bomb, op)

	fish1, _, err := ebitenutil.NewImageFromFile("utils/fish/poisson-2.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish1.Size()
	variable.FishImage1 = ebiten.NewImage(w, h)
	variable.FishImage1.DrawImage(fish1, op)

	fish2, _, err := ebitenutil.NewImageFromFile("utils/fish/poisson-3.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish2.Size()
	variable.FishImage2 = ebiten.NewImage(w, h)
	variable.FishImage2.DrawImage(fish2, op)

	fish3, _, err := ebitenutil.NewImageFromFile("utils/fish/poisson-5.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish1.Size()
	variable.FishImage3 = ebiten.NewImage(w, h)
	variable.FishImage3.DrawImage(fish3, op)

	preda, _, err := ebitenutil.NewImageFromFile("utils/fish/predT3.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish1.Size()
	variable.PredImage = ebiten.NewImage(w+100, h+100)
	variable.PredImage.DrawImage(preda, op)

	back, _, err := ebitenutil.NewImageFromFile("utils/fish/background.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = back.Size()
	variable.BackgroundImage = ebiten.NewImage(w, h)
	variable.BackgroundImage.DrawImage(back, op)

}

func main() {
	// création de la chanel de sync pour la music:
	c1 := make(chan string)

	musicAgent := music.NewMusicAgent("utils/music/jaws.mp3", "utils/music/jaws.wav", c1)
	musicAgent.Start()

	ebiten.SetWindowSize(constant.ScreenWidth, constant.ScreenHeight)
	ebiten.SetWindowTitle("Le Meilleur Jeu de Pêche SMA de la planète")
	if err := ebiten.RunGame(game.NewGame(c1, 5)); err != nil {
		log.Fatal(err)
	}
}
