package agent

import (
	"math/rand"

	utils "gitlab.utc.fr/projet_ia04/Boid/utils"
	worldelements "gitlab.utc.fr/projet_ia04/Boid/worldelements"
)

type Boid struct {
	ImageWidth     int
	ImageHeight    int
	Position       utils.Vector2D
	Velocity       utils.Vector2D
	Acceleration   utils.Vector2D
	Species        int
	Dead           bool
	EscapePredator float64
	Marqued        bool
}

//GenerateBoids fonction qui permet de générer des agents Boids
func GenerateBoids() []*Boid {
	boids := make([]*Boid, utils.NumBoids)
	for i := range boids {
		w, h := utils.FishImage1.Size()

		// Pour éviter que les agents apparaisent au dessus ou en dessous des murs de bombes:
		// on les fait apparaitre  horizontalement en ligne au mileu  de l'écran:
		middle := utils.ScreenHeight / 2
		x, y := rand.Float64()*float64(utils.ScreenWidth-w), float64(middle)

		min, max := -utils.MaxForce, utils.MaxForce
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		s := rand.Intn(utils.NumSpecies)
		boids[i] = &Boid{
			ImageWidth:     w,
			ImageHeight:    h,
			Position:       utils.Vector2D{X: x, Y: y},
			Velocity:       utils.Vector2D{X: vx, Y: vy},
			Acceleration:   utils.Vector2D{X: 0, Y: 0},
			Species:        s,
			Dead:           false,
			EscapePredator: 80.0,
			Marqued:        false,
		}
	}
	return boids
}

//Update fonction qui permet de mettre à jour le Boid
func (boid *Boid) Update(level int, walls []*worldelements.Wall, boids []*Boid, predators []*Predator) {
	if !boid.CheckEdges() {
		if !boid.CheckWalls(walls) {
			boid.ApplyRules(boids, predators)
		}
	}
	boid.ApplyMovement()

	// Pour éviter que les poissons réussissent à s'échapper des murs de bombes
	// dans les niveaux supérieurs ou égal au niveau 4: (le dernier niveau pour le moment)
	// Dès que l'on detecte qu'ils ne sont pas où ils devraient être, on les fait réapparaitre au centre

	if level >= 4 && (boid.Position.Y <= 0 || boid.Position.Y >= float64(utils.ScreenHeight)) {
		boid.Position.Y = float64(utils.ScreenHeight) / 2
	}
}

func (boid *Boid) ApplyRules(restOfFlock []*Boid, predators []*Predator) {
	if !boid.Dead {
		alignSteering := utils.Vector2D{}
		alignTotal := 0
		cohesionSteering := utils.Vector2D{}
		cohesionTotal := 0
		separationSteering := utils.Vector2D{}
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
				goback := utils.Rotate(boid.Velocity, int(rand.Float64()*360))
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
					goback := utils.Rotate(boid.Velocity, int(rand.Float64()*360))
					goback.Multiply(2.5 * 3.0)
					boid.Acceleration.Add(goback)
				}
				if boid.Species == 3 || boid.Species == 4 {
					// Le boid alerte tous ses voisins qui partent dans une direction donnée si il n'est pas marqué
					for _, other := range restOfFlock {
						d := boid.Position.Distance(other.Position)
						if d < utils.SeparationPerception {
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
					if boid.Species == other.Species && d < utils.AlignPerception {
						alignTotal++
						alignSteering.Add(other.Velocity)
					}
					if boid.Species == other.Species && d < utils.CohesionPerception {
						cohesionTotal++
						cohesionSteering.Add(other.Position)
					}
					if d < utils.SeparationPerception {
						separationTotal++
						diff := boid.Position
						diff.Subtract(other.Position)
						diff.Divide(d)
						separationSteering.Add(diff)
						if other.Species != boid.Species {
							diff.Divide(d * utils.RepulsionFactorBtwnSpecies)
							separationSteering.Add(diff)
						}
					}
				}
			}
		}

		if separationTotal > 0 {
			separationSteering.Divide(float64(separationTotal))
			separationSteering.SetMagnitude(utils.MaxSpeed)
			separationSteering.Subtract(boid.Velocity)
			separationSteering.SetMagnitude(utils.MaxForce * 1.2)
		}
		if cohesionTotal > 0 {
			cohesionSteering.Divide(float64(cohesionTotal))
			cohesionSteering.Subtract(boid.Position)
			cohesionSteering.SetMagnitude(utils.MaxSpeed)
			cohesionSteering.Subtract(boid.Velocity)
			cohesionSteering.SetMagnitude(utils.MaxForce * 0.9)
		}
		if alignTotal > 0 {
			alignSteering.Divide(float64(alignTotal))
			alignSteering.SetMagnitude(utils.MaxSpeed)
			alignSteering.Subtract(boid.Velocity)
			alignSteering.Limit(utils.MaxForce)
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
		boid.Velocity.Limit(utils.MaxSpeed)
		boid.Acceleration.Multiply(0.0)
	}
}

func (boid *Boid) CheckEdges() bool {
	if boid.Position.X < 0 {
		boid.Position.X = utils.ScreenWidth
	} else if boid.Position.X > utils.ScreenWidth {
		boid.Position.X = 0
	}
	return false
}

//CheckWall fonction qui permet de vérifier la présence de murs
func (boid *Boid) CheckWalls(walls []*worldelements.Wall) bool {
	separationTotal := 0
	separationSteering := utils.Vector2D{}
	for _, wall := range walls {
		d := boid.Position.Distance(wall.Position)
		if d < utils.WallSeparationPerception {
			separationTotal++
			diff := boid.Position
			diff.Subtract(wall.Position)
			diff.Divide(d)
			separationSteering.Add(diff)
		}
	}
	if separationTotal > 0 {
		separationSteering.Divide(float64(separationTotal))
		separationSteering.SetMagnitude(utils.MaxSpeed)
		separationSteering.Subtract(boid.Velocity)
		separationSteering.SetMagnitude(utils.MaxForce * 1.2)
		boid.Acceleration.Add(separationSteering)
		return true
	}
	return false
}
