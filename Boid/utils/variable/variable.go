package variable

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

var (
	BirdImage                  *ebiten.Image
	FishImage1                 *ebiten.Image
	FishImage2                 *ebiten.Image
	FishImage3                 *ebiten.Image
	BackgroundImage            *ebiten.Image
	WallImage                  *ebiten.Image
	RepulsionFactorBtwnSpecies float64
	SeparationPerception       float64
	CohesionPerception         float64
)
