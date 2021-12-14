package variable

import (
	"github.com/hajimehoshi/ebiten"
)

var (
	BirdImage                  *ebiten.Image
	WallImage                  *ebiten.Image
	RepulsionFactorBtwnSpecies float64
	SeparationPerception       float64
	CohesionPerception         float64
)
