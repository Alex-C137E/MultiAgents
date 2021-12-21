package boid

import (
	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
	constant "gitlab.utc.fr/projet_ia04/Boid/utils/constant"
	variable "gitlab.utc.fr/projet_ia04/Boid/utils/variable"
	vector "gitlab.utc.fr/projet_ia04/Boid/utils/vector"
)

type Vector2D = vector.Vector2D

type Boid struct {
	ImageWidth   int
	ImageHeight  int
	Position     Vector2D
	Velocity     Vector2D
	Acceleration Vector2D
	Species      int
	Dead         bool
}

func (boid *Boid) ApplyRules(restOfFlock []*Boid) {
	if !boid.Dead {
		alignSteering := Vector2D{}
		alignTotal := 0
		cohesionSteering := Vector2D{}
		cohesionTotal := 0
		separationSteering := Vector2D{}
		separationTotal := 0

		for _, other := range restOfFlock {
			d := boid.Position.Distance(other.Position)
			if boid != other {
				if boid.Species == other.Species && d < constant.AlignPerception {
					alignTotal++
					alignSteering.Add(other.Velocity)
				}
				if boid.Species == other.Species && d < variable.CohesionPerception {
					cohesionTotal++
					cohesionSteering.Add(other.Position)
				}
				if d < variable.SeparationPerception {
					separationTotal++
					diff := boid.Position
					diff.Subtract(other.Position)
					diff.Divide(d)
					separationSteering.Add(diff)
					if other.Species != boid.Species {
						diff.Divide(d * variable.RepulsionFactorBtwnSpecies)
						separationSteering.Add(diff)
					}
				}
			}
		}

		if separationTotal > 0 {
			separationSteering.Divide(float64(separationTotal))
			separationSteering.SetMagnitude(constant.MaxSpeed)
			separationSteering.Subtract(boid.Velocity)
			separationSteering.SetMagnitude(constant.MaxForce * 1.2)
		}
		if cohesionTotal > 0 {
			cohesionSteering.Divide(float64(cohesionTotal))
			cohesionSteering.Subtract(boid.Position)
			cohesionSteering.SetMagnitude(constant.MaxSpeed)
			cohesionSteering.Subtract(boid.Velocity)
			cohesionSteering.SetMagnitude(constant.MaxForce * 0.9)
		}
		if alignTotal > 0 {
			alignSteering.Divide(float64(alignTotal))
			alignSteering.SetMagnitude(constant.MaxSpeed)
			alignSteering.Subtract(boid.Velocity)
			alignSteering.Limit(constant.MaxForce)
		}

		boid.Acceleration.Add(alignSteering)
		boid.Acceleration.Add(cohesionSteering)
		boid.Acceleration.Add(separationSteering)
		boid.Acceleration.Divide(3)
	}
}

func (boid *Boid) ApplyMovement() {
	if !boid.Dead {
		boid.Position.Add(boid.Velocity)
		boid.Velocity.Add(boid.Acceleration)
		boid.Velocity.Limit(constant.MaxSpeed)
		boid.Acceleration.Multiply(0.0)
	}
}

// func (boid *Boid) CheckEdges() bool {
// 	separationTotal := 0
// 	separationSteering := Vector2D{}
// 	if boid.Position.X < 10 {
// 		separationTotal++
// 		Position := Vector2D{X: 0, Y: boid.Position.Y}
// 		d := boid.Position.Distance(Position)
// 		diff := boid.Position
// 		diff.Subtract(Position)
// 		diff.Divide(d)
// 		separationSteering.Add(diff)
// 	}
// 	if boid.Position.X > constant.ScreenWidth-10 {
// 		separationTotal++
// 		Position := Vector2D{X: constant.ScreenWidth, Y: boid.Position.Y}
// 		d := boid.Position.Distance(Position)
// 		diff := boid.Position
// 		diff.Subtract(Position)
// 		diff.Divide(d)
// 		separationSteering.Add(diff)
// 	}
// 	if boid.Position.Y < 10 {
// 		separationTotal++
// 		Position := Vector2D{X: boid.Position.X, Y: 0}
// 		d := boid.Position.Distance(Position)
// 		diff := boid.Position
// 		diff.Subtract(Position)
// 		diff.Divide(d)
// 		separationSteering.Add(diff)
// 	}
// 	if boid.Position.Y > constant.ScreenHeight-10 {
// 		separationTotal++
// 		Position := Vector2D{X: boid.Position.X, Y: constant.ScreenHeight}
// 		d := boid.Position.Distance(Position)
// 		diff := boid.Position
// 		diff.Subtract(Position)
// 		diff.Divide(d)
// 		separationSteering.Add(diff)
// 	}
// 	if separationTotal > 0 {
// 		separationSteering.Divide(float64(separationTotal))
// 		separationSteering.SetMagnitude(constant.MaxSpeed)
// 		separationSteering.Subtract(boid.Velocity)
// 		separationSteering.SetMagnitude(constant.MaxForce * 1.2)
// 		boid.Acceleration.Add(separationSteering)
// 		return true
// 	}
// 	return false
// }

// NO BORDER VERSION
func (boid *Boid) CheckEdges() bool {
	if boid.Position.X < 0 {
		boid.Position.X = constant.ScreenWidth
	} else if boid.Position.X > constant.ScreenWidth {
		boid.Position.X = 0
	}
	if boid.Position.Y < 0 {
		boid.Position.Y = constant.ScreenHeight
	} else if boid.Position.Y > constant.ScreenHeight {
		boid.Position.Y = 0
	}
	return false
}

func (boid *Boid) CheckWalls(walls []*wall.Wall) bool {
	separationTotal := 0
	separationSteering := Vector2D{}
	for _, wall := range walls {
		d := boid.Position.Distance(wall.Position)
		if d < constant.WallSeparationPerception {
			separationTotal++
			diff := boid.Position
			diff.Subtract(wall.Position)
			diff.Divide(d)
			separationSteering.Add(diff)
		}
	}
	if separationTotal > 0 {
		separationSteering.Divide(float64(separationTotal))
		separationSteering.SetMagnitude(constant.MaxSpeed)
		separationSteering.Subtract(boid.Velocity)
		separationSteering.SetMagnitude(constant.MaxForce * 1.2)
		boid.Acceleration.Add(separationSteering)
		return true
	}
	return false
}
