package main

import (
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"

	agent "gitlab.utc.fr/projet_ia04/Boid/agent"
	game "gitlab.utc.fr/projet_ia04/Boid/game"
	utils "gitlab.utc.fr/projet_ia04/Boid/utils"
)

//init fonction permet de charger les images nécessaires au fonctionnement du jeu
func init() {
	fish, _, err := ebitenutil.NewImageFromFile("utils/images/chevron-up.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h := fish.Size()
	utils.BirdImage = ebiten.NewImage(w-20, h-10)
	op := &ebiten.DrawImageOptions{}
	utils.BirdImage.DrawImage(fish, op)

	bomb, _, err := ebitenutil.NewImageFromFile("utils/images/bomb.png")
	if err != nil {
		log.Fatal(err)
	}
	bombW, bombH := 24, 24
	utils.WallImage = ebiten.NewImage(bombW, bombH)
	utils.WallImage.DrawImage(bomb, op)

	sand, _, err := ebitenutil.NewImageFromFile("utils/images/sand2.png")
	if err != nil {
		log.Fatal(err)
	}
	sW, sH := sand.Size()
	utils.SandImage = ebiten.NewImage(sW, sH)
	utils.SandImage.DrawImage(sand, op)

	fish1, _, err := ebitenutil.NewImageFromFile("utils/images/poisson-2.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish1.Size()
	utils.FishImage1 = ebiten.NewImage(w, h)
	utils.FishImage1.DrawImage(fish1, op)

	fish2, _, err := ebitenutil.NewImageFromFile("utils/images/poisson-3.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish2.Size()
	utils.FishImage2 = ebiten.NewImage(w, h)
	utils.FishImage2.DrawImage(fish2, op)

	fish3, _, err := ebitenutil.NewImageFromFile("utils/images/poisson-5.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish1.Size()
	utils.FishImage3 = ebiten.NewImage(w, h)
	utils.FishImage3.DrawImage(fish3, op)

	preda, _, err := ebitenutil.NewImageFromFile("utils/images/poisson-4.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = fish1.Size()
	utils.PredImage = ebiten.NewImage(w+100, h+100)
	utils.PredImage.DrawImage(preda, op)

	back, _, err := ebitenutil.NewImageFromFile("utils/images/background.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h = back.Size()
	utils.BackgroundImage = ebiten.NewImage(w, h)
	utils.BackgroundImage.DrawImage(back, op)

}

func main() {
	c1 := make(chan string) // création de la chanel de sync pour la music:

	musicAgent := agent.NewMusicAgent("utils/music/jaws.mp3", "utils/music/jaws.wav", c1) //chargement de la musique des Dents de la mer
	musicAgent.Start()

	ebiten.SetWindowSize(utils.ScreenWidth, utils.ScreenHeight)
	ebiten.SetWindowTitle("Le Meilleur Jeu de Pêche SMA de la planète") //True story
	if err := ebiten.RunGame(game.NewGame(c1, 5)); err != nil {         //lancement du jeu
		log.Fatal(err)
	}
}
