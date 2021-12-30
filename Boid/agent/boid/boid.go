package boid

import (
	"math"
	"math/rand"

	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
	constant "gitlab.utc.fr/projet_ia04/Boid/utils/constant"
	variable "gitlab.utc.fr/projet_ia04/Boid/utils/variable"
	vector "gitlab.utc.fr/projet_ia04/Boid/utils/vector"
)

type Vector2D = vector.Vector2D

type Boid struct {
	ImageWidth     int
	ImageHeight    int
	Position       Vector2D
	Velocity       Vector2D
	Acceleration   Vector2D
	Species        int
	Dead           bool
	EscapePredator float64
	Marqued        bool
}

type Predator struct {
	ImageWidth   int
	ImageHeight  int
	Position     Vector2D
	Velocity     Vector2D
	Acceleration Vector2D
	Density      int
	Angle        int
	Dist         int
	V1           Vector2D
	V2           Vector2D
	R            bool
}

func (boid *Boid) ApplyRules(restOfFlock []*Boid, predators []*Predator) {
	if !boid.Dead {
		alignSteering := Vector2D{}
		alignTotal := 0
		cohesionSteering := Vector2D{}
		cohesionTotal := 0
		separationSteering := Vector2D{}
		separationTotal := 0

		istherepred := false
		// check predator presence
		if boid.Marqued {
			if boid.Species == 3 {
				goback := boid.Velocity
				goback.Multiply(2.5 * 3.0)
				boid.Acceleration.Add(goback)
			}
			if boid.Species == 4 {
				goback := vector.Rotate(boid.Velocity, int(rand.Float64()*360))
				goback.Multiply(2.5 * 3.0)
				boid.Acceleration.Add(goback)
			}
			boid.Marqued = false
			return
		}

		for _, pred := range predators {
			d := boid.Position.Distance(pred.Position)
			if d < boid.EscapePredator {
				istherepred = true
				// 180 retour espece 1
				if boid.Species == 1 {
					goback := boid.Velocity // on divise par 3 à la fin de la fonction
					goback.Multiply(2.5 * 3.0)
					boid.Acceleration.Add(goback)
				}
				if boid.Species == 2 {
					// eclatement de la population, orientation aleatoire
					goback := vector.Rotate(boid.Velocity, int(rand.Float64()*360))
					goback.Multiply(2.5 * 3.0)
					boid.Acceleration.Add(goback)
				}
				if boid.Species == 3 || boid.Species == 4 {
					// Le boid alerte tous ses voisins qui partent dans une direction donnée si il n'est pas marqué
					for _, other := range restOfFlock {
						d := boid.Position.Distance(other.Position)
						if d < variable.SeparationPerception {
							other.Marqued = true
						}
					}
				}
			}
		}

		if !(istherepred) {
			for _, other := range restOfFlock {
				d := boid.Position.Distance(other.Position)
				if boid != other {
					if boid.Species == other.Species && d < variable.AlignPerception {
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
		}

		if separationTotal > 0 {
			separationSteering.Divide(float64(separationTotal))
			separationSteering.SetMagnitude(variable.MaxSpeed)
			separationSteering.Subtract(boid.Velocity)
			separationSteering.SetMagnitude(variable.MaxForce * 1.2)
		}
		if cohesionTotal > 0 {
			cohesionSteering.Divide(float64(cohesionTotal))
			cohesionSteering.Subtract(boid.Position)
			cohesionSteering.SetMagnitude(variable.MaxSpeed)
			cohesionSteering.Subtract(boid.Velocity)
			cohesionSteering.SetMagnitude(variable.MaxForce * 0.9)
		}
		if alignTotal > 0 {
			alignSteering.Divide(float64(alignTotal))
			alignSteering.SetMagnitude(variable.MaxSpeed)
			alignSteering.Subtract(boid.Velocity)
			alignSteering.Limit(variable.MaxForce)
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
		boid.Velocity.Limit(variable.MaxSpeed)
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
		separationSteering.SetMagnitude(variable.MaxSpeed)
		separationSteering.Subtract(boid.Velocity)
		separationSteering.SetMagnitude(variable.MaxForce * 1.2)
		boid.Acceleration.Add(separationSteering)
		return true
	}
	return false
}

// ----------------------------------------------
// ---------------------------------- PREDATOR -------------------------- //
// ----------------------------------------------

func (preda *Predator) Vision() []Vector2D {
	x := preda.Position.X
	y := preda.Position.Y
	vx := preda.Velocity.X
	vy := preda.Velocity.Y

	//calculate angle between vect Velocity and x-axis
	angR := math.Atan2(vy, vx)
	angV := int(angR * 180 / math.Pi)
	// Calulate upper and lower angle
	angU := angV - preda.Angle
	angL := angV + preda.Angle

	// Calculate new point
	l1x := x + float64(preda.Dist)*math.Cos(vector.AngleToRadians(angU))
	l1y := y + float64(preda.Dist)*math.Sin(vector.AngleToRadians(angU))

	l2x := x + float64(preda.Dist)*math.Cos(vector.AngleToRadians(angL))
	l2y := y + float64(preda.Dist)*math.Sin(vector.AngleToRadians(angL))

	p1 := Vector2D{X: l1x, Y: l1y}
	p2 := Vector2D{X: l2x, Y: l2y}

	mapP := make([]Vector2D, 2)
	mapP[0] = p1
	mapP[1] = p2
	return mapP
}

func (preda *Predator) ApplyRules(restOfFlock []*Boid) {

	//preda.V1 = vPoint[0]
	//preda.V2 = vPoint[1]
	var dens int
	var newP Predator
	var new bool
	var densMax = 0
	var densMax2 = 0
	var proie1 *Boid
	var proie2 *Boid
	var vPoint2 []Vector2D

	//Tuer
	for _, Boid := range restOfFlock {
		if (preda.Position.Distance(Boid.Position)) < 15 {
			Boid.Dead = true
		}
	}

	//Vision + stratégie d'attaque
	vPoint := preda.Vision()
	new = false
	if (vPoint[0].X > constant.ScreenWidth || vPoint[0].X < 0) || (vPoint[0].Y > constant.ScreenHeight || vPoint[0].Y < 0) || (vPoint[1].X > constant.ScreenWidth || vPoint[1].X < 0) || (vPoint[1].Y > constant.ScreenHeight || vPoint[1].Y < 0) {
		new = true
		newP = *preda
		if vPoint[0].X > constant.ScreenWidth {
			newP.Position.X = preda.Position.X - constant.ScreenWidth
		} else if vPoint[0].X < 0 {
			newP.Position.X = preda.Position.X + constant.ScreenWidth
		}
		if vPoint[0].Y > constant.ScreenHeight {
			newP.Position.Y = preda.Position.Y - constant.ScreenHeight
		} else if vPoint[0].Y < 0 {
			newP.Position.Y = preda.Position.Y + constant.ScreenHeight
		}

		if vPoint[1].X > constant.ScreenWidth {
			newP.Position.X = preda.Position.X - constant.ScreenWidth
		} else if vPoint[1].X < 0 {
			newP.Position.X = preda.Position.X + constant.ScreenWidth
		}
		if vPoint[1].Y > constant.ScreenHeight {
			newP.Position.Y = preda.Position.Y - constant.ScreenHeight
		} else if vPoint[1].Y < 0 {
			newP.Position.Y = preda.Position.Y + constant.ScreenHeight
		}
	}
	preda.V1 = vPoint[0]
	preda.V2 = vPoint[1]
	if new {
		vPoint2 = newP.Vision()
		if (vPoint[1].X > constant.ScreenWidth || vPoint[1].X < 0) || (vPoint[1].Y > constant.ScreenHeight || vPoint[1].Y < 0) {
			preda.V2 = vPoint2[1]

		}
		if (vPoint[0].X > constant.ScreenWidth || vPoint[0].X < 0) || (vPoint[0].Y > constant.ScreenHeight || vPoint[0].Y < 0) {
			preda.V1 = vPoint2[0]
		}
	}

	for _, Boid := range restOfFlock {
		if !Boid.Dead {
			dens = 0
			b := vector.PointInTriangle(Boid.Position, preda.Position, vPoint[0], vPoint[1])
			b2 := false
			if new {
				b2 = vector.PointInTriangle(Boid.Position, newP.Position, vPoint2[0], vPoint2[1])
			}
			if b {
				//print(1)
				for _, other := range restOfFlock {
					if (Boid.Position.Distance(other.Position)) < 30 && !other.Dead {
						dens++
					}
				}
			}
			if dens > densMax {
				densMax = dens
				proie1 = Boid
			}
			dens = 0
			if b2 {
				//print(1)
				for _, other := range restOfFlock {
					if (Boid.Position.Distance(other.Position)) < 30 && !other.Dead {
						dens++
					}
				}
			}
			if dens > densMax2 {
				densMax2 = dens
				proie2 = Boid
			}
		}
	}
	if densMax >= densMax2 && densMax > preda.Density {
		Vit := vector.Vector2D{X: proie1.Position.X - preda.Position.X, Y: proie1.Position.Y - preda.Position.Y}
		Vit.Normalize()
		Vit.X = Vit.X * 10
		Vit.Y = Vit.Y * 10
		preda.Velocity = Vit

	} else if densMax2 > preda.Density {
		Vit := vector.Vector2D{X: proie2.Position.X - newP.Position.X, Y: proie2.Position.Y - newP.Position.Y}
		Vit.Normalize()
		Vit.X = Vit.X * 10
		Vit.Y = Vit.Y * 10
		preda.Velocity = Vit
	} else {
		vit := preda.Velocity
		if vit.X > 1 || -vit.X > 1 {
			vit.X = vit.X * 0.98
		}
		if vit.Y > 1 || -vit.Y > 1 {
			vit.Y = vit.Y * 0.98
		}
		preda.Velocity = vector.Vector2D{X: vit.X, Y: vit.Y}
		if rand.Float64() < 0.01 {
			preda.Velocity = vector.Rotate(preda.Velocity, rand.Intn(10))
		}
		if rand.Float64() > 0.99 {
			preda.Velocity = vector.Rotate(preda.Velocity, -rand.Intn(10))
		}

	}
}

func (preda *Predator) CheckEdges() bool {
	if preda.Position.X < 0 {
		preda.Position.X = constant.ScreenWidth
	} else if preda.Position.X > constant.ScreenWidth {
		preda.Position.X = 0
	}
	if preda.Position.Y < 0 {
		preda.Position.Y = constant.ScreenHeight
	} else if preda.Position.Y > constant.ScreenHeight {
		preda.Position.Y = 0
	}
	return false
}

func (preda *Predator) CheckWalls(walls []*wall.Wall) {
	if preda.R {
		preda.Velocity.Normalize()
		preda.Velocity.X = preda.Velocity.X * 2
		preda.Velocity.Y = preda.Velocity.Y * 2
		preda.R = false
	} else {
		for _, wall := range walls {
			d := preda.Position.Distance(wall.Position)
			if d <= 20 {
				preda.Velocity.Normalize()
				preda.Velocity.X = -preda.Velocity.X * (100 - d + 1)
				preda.Velocity.Y = -preda.Velocity.Y * (100 - d + 1)
				preda.R = true
				break
			}
		}
	}
}

func (preda *Predator) ApplyMovement() {
	preda.Position.Add(preda.Velocity)
}
