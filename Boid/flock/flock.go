package flock

import (
	agent "gitlab.utc.fr/projet_ia04/Boid/agent"
	worldelements "gitlab.utc.fr/projet_ia04/Boid/worldelements"
)

//Flock structure qui permet de garder en contexte tous les éléments du monde
type Flock struct {
	Boids     []*agent.Boid
	Walls     []*worldelements.Wall
	Predators []*agent.Predator
}

//Logic fonction permet de mettre à jours les agents
func (flock *Flock) Logic(level int) {
	for _, boid := range flock.Boids {
		go boid.Update(level, flock.Walls, flock.Boids, flock.Predators)
	}
	for _, preda := range flock.Predators {
		go preda.Update(flock.Walls, flock.Boids)
	}
}
