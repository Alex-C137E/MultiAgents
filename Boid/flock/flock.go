package flock

import (
	boid "gitlab.utc.fr/projet_ia04/Boid/agent/boid"
	predator "gitlab.utc.fr/projet_ia04/Boid/agent/predator"
	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
)

type Flock struct {
	Boids     []*boid.Boid
	Walls     []*wall.Wall
	Predators []*predator.Predator
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
	for _, preda := range flock.Predators {
		if !preda.CheckEdges() {
			preda.ApplyRules(flock.Boids)
			preda.CheckWalls(flock.Walls)
		}
		preda.ApplyMovement()
	}
}
