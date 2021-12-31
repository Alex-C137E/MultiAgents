package flock

import (
	boid "gitlab.utc.fr/projet_ia04/Boid/agent/boid"
	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
	"gitlab.utc.fr/projet_ia04/Boid/utils/constant"
)

type Flock struct {
	Boids     []*boid.Boid
	Walls     []*wall.Wall
	Predators []*boid.Predator
}

func (flock *Flock) Logic(level int) {
	for _, boid := range flock.Boids {
		if !boid.CheckEdges() {
			if !boid.CheckWalls(flock.Walls) {
				boid.ApplyRules(flock.Boids, flock.Predators)
			}
		}
		boid.ApplyMovement()

		// Pour éviter que les poissons réussissent à s'échapper des murs de bombes
		// dans les niveaux supérieurs ou égal au niveau 4: (le dernier niveau pour le moment)
		// Dès que l'on detecte qu'ils ne sont pas où ils devraient être, on les fait réapparaitre au centre

		if level >= 4 && (boid.Position.Y <= 0 || boid.Position.Y >= float64(constant.ScreenHeight)) {
			boid.Position.Y = float64(constant.ScreenHeight) / 2
		}

	}
	for _, preda := range flock.Predators {
		if !preda.CheckEdges() {
			preda.ApplyRules(flock.Boids)
			preda.CheckWalls(flock.Walls)
		}
		preda.ApplyMovement()
	}
}
