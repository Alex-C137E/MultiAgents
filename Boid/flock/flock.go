package flock

import (
	boid "gitlab.utc.fr/projet_ia04/Boid/agent/boid"
	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
)

type Flock struct {
	Boids []*boid.Boid
	Walls []*wall.Wall
}

func (flock *Flock) Logic() {
	for _, boid := range flock.Boids {
		if !boid.CheckEdges() {
			if !boid.CheckWalls(flock.Walls) {
				boid.ApplyRules(flock.Boids)
			}
		}
		boid.ApplyMovement()
	}
}
