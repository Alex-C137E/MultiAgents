package game

type Level struct {
	RepulsionFactorBtwnSpecies float64
	SeparationPerception       float64
	CohesionPerception         float64
	AlignPerception            float64
	numWall                    int
	MaxForce                   float64
	MaxSpeed                   float64
	polygonSize                float64
}

func NewLevel(RepulsionFactorBtwnSpecies float64,
	SeparationPerception float64,
	CohesionPerception float64,
	AlignPerception float64,
	numWall int,
	MaxForce float64,
	MaxSpeed float64,
	polgonSize float64) *Level {
	return &Level{RepulsionFactorBtwnSpecies,
		SeparationPerception,
		CohesionPerception,
		AlignPerception,
		numWall,
		MaxForce,
		MaxSpeed,
		polgonSize}
}
