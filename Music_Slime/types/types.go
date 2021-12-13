package types

//Structure for the position of an agent
type Position struct {
	X float64 //Agent position x
	Y float64 //Agent position y
}

type TrailMap [][][]float64

type SpeciesMask [4]int

type Colour [4]float64

type Window struct {
	Width  float64
	Height float64
}

func GetWindowDefault() Window {
	return Window{500, 500}
}

const AgentsNum = 25000
const SpeciesNum = 3
const TrailWeight = 2.0
const DeltaTime = 0.02
const DiffuseRate = 3
const DecayRate = 0.5
