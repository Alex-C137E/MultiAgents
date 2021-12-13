package settings

import (
	"log"

	_ "image/png" //sinon ça fait une erreur quand on load l'image

	ebiten "github.com/hajimehoshi/ebiten/v2"
	ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//Settings type regroupe les paramêtres globaux de la simulation
type Settings struct {
	AgentsNum            int
	AgentLongueur        float64
	AgentRayonSeparation float64
	AgentRayonCohesion   float64
	AgentRayonVitesse    float64
	AgentMaxForce        float64
	BirdImage            *ebiten.Image
	BirdImage2           *ebiten.Image
}

//GetDefaultSettings fonction qui renvoie les paramêtres globaux par défaut de la simulation
func GetDefaultSettings() Settings {
	bImg, _, err := ebitenutil.NewImageFromFile("t2.png")
	bImg2, _, err := ebitenutil.NewImageFromFile("t3.png")

	if err != nil {
		log.Fatal(err)
	}

	return Settings{100, 5, 32, 30, 27, 30, bImg, bImg2}
}
