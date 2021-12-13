package agent

import (
	"math"
	"math/rand"

	ebiten "github.com/hajimehoshi/ebiten/v2"

	settings "gitlab.utc.fr/jedescam/projet_ia04/final/settings"
	vectors "gitlab.utc.fr/jedescam/projet_ia04/final/vectors"
)

//Bird type représente un agent qui se déplace de manière groupé dans le monde. Equivalent au comportement d'un oiseau
type Bird struct {
	ID              int
	position        vectors.Vector2D
	vitesse         vectors.Vector2D
	acceleration    vectors.Vector2D
	longueur        float64
	rayonSeparation float64
	rayonCohesion   float64
	rayonVitesse    float64
	maxForce        float64
	maxSpeed        vectors.Vector2D
	separation      map[*Bird]bool
	Cohesion        map[*Bird]bool
	alignement      map[*Bird]bool
	evitement       int
	ty              int
}

//NewBird constructeur qui permet de créer un nouvel agent
func NewBird(s settings.Settings, id int) *Bird {
	return &Bird{id, vectors.Vector2D{0, 0}, vectors.Vector2D{0, 0}, vectors.Vector2D{0, 0}, s.AgentLongueur, s.AgentRayonSeparation, s.AgentRayonCohesion, s.AgentRayonVitesse, s.AgentMaxForce, vectors.Vector2D{X: 10, Y: 10}, make(map[*Bird]bool), make(map[*Bird]bool), make(map[*Bird]bool), 0, 1}
}

//Start fonction qui initialise l'agent
func (b *Bird) Start(window settings.Window) {
	//Set agent default position randomly
	x, y := rand.Float64()*window.Width, rand.Float64()*window.Height
	b.position = vectors.Vector2D{X: x, Y: y}

	//Set agent default speed randomly
	t := rand.Float64()
	var vx, vy float64
	vx = 2
	vy = 2

	if t > 0.25 {
		vx, vy = -2, -2
	} else if t > 0.5 {
		vx, vy = -2, 2
		b.ty = 2
	} else if t > 0.75 {
		vx, vy = 2, -2
	}
	b.vitesse = vectors.Vector2D{X: float64(vx), Y: float64(vy)}
}

//ApplyForce applique une force sur l'agent
func (b *Bird) ApplyForce(force vectors.Vector2D) {
	b.acceleration.X += force.X
	b.acceleration.Y += force.Y
}

func (b *Bird) separate() (t vectors.Vector2D) {
	steer := vectors.Vector2D{X: 0, Y: 0}
	for bird := range b.separation {
		steer.X = steer.X - bird.position.X
		steer.Y = steer.Y - bird.position.Y
	}
	if len(b.separation) != 0 {
		steer.X = steer.X / float64(len(b.separation))
		steer.Y = steer.Y / float64(len(b.separation))
	}
	return steer
}

func (b *Bird) align() (t vectors.Vector2D) {
	alignement := vectors.Vector2D{X: 0, Y: 0}
	for bird := range b.alignement {
		alignement.X = alignement.X + bird.vitesse.X
		alignement.Y = alignement.Y + bird.vitesse.Y
	}
	if len(b.alignement) != 0 {
		alignement.X = alignement.X / float64(len(b.alignement))
		alignement.Y = alignement.Y / float64(len(b.alignement))
	}
	return alignement
}

func (b *Bird) cohesion() (t vectors.Vector2D) {
	cohesion := vectors.Vector2D{X: 0, Y: 0}
	for bird := range b.Cohesion {
		cohesion.X = cohesion.X + bird.position.X
		cohesion.Y = cohesion.Y + bird.position.Y
	}
	if len(b.Cohesion) != 0 {
		cohesion.X = cohesion.X / float64(len(b.Cohesion))
		cohesion.Y = cohesion.Y / float64(len(b.Cohesion))
	}
	return cohesion
}

func (b *Bird) random() (t vectors.Vector2D) {
	r := vectors.Vector2D{X: 0, Y: 0}
	u := rand.Float64()
	if u < 0.1 {
		r.X = rand.Float64() * 20
		r.Y = rand.Float64() * 20
	}
	if u < 0.5 {
		r.X = -r.X
		r.Y = -r.Y
	}
	return r
}

func (b *Bird) limit() {
	norme := math.Sqrt(math.Pow(b.vitesse.X, 2) + math.Pow(b.vitesse.Y, 2))
	normeMax := math.Sqrt(math.Pow(b.maxSpeed.X, 2) + math.Pow(b.maxSpeed.Y, 2))
	if norme >= normeMax {
		if b.vitesse.X > 0 {
			b.vitesse.X = b.maxSpeed.X / 2
		} else {
			b.vitesse.X = -(b.maxSpeed.X / 2)
		}
		if b.vitesse.Y > 0 {
			b.vitesse.Y = b.maxSpeed.Y / 2
		} else {
			b.vitesse.Y = -(b.maxSpeed.Y / 2)
		}
	}
}

//Update fonction qui permet de mettre à jour le comportement de l'agent
func (b *Bird) Update(window settings.Window, birds []*Bird) {
	b.getCouches(birds)
	b.flock()
	b.checkEdges(window)
	b.updatePosition()
}

//getCouches fonction qui met à jour les couches de l'agents
func (b *Bird) getCouches(birds []*Bird) {
	r := 0
	for _, bird := range birds {
		if bird.ty != b.ty && (math.Sqrt(math.Pow(b.position.X-bird.position.X, 2)+math.Pow(b.position.Y-bird.position.Y, 2))) <= 20 && b.ty == 1 && b.evitement != 2 {
			b.evitement = 1
			r = 1
		}
		if (math.Sqrt(math.Pow(b.position.X-bird.position.X, 2)+math.Pow(b.position.Y-bird.position.Y, 2))) < 20 && b.evitement == 2 && bird.ty != b.ty && b.ty == 1 {
			r = 1
		}
	}

	if r == 0 {
		b.evitement = 0
	}
	// Couche separation, cohesion, alignement
	for _, bird := range birds {
		if bird.ty == b.ty && b.evitement == 0 {
			if bird.ID == b.ID {
				continue
			}
			distance := math.Sqrt(math.Pow(b.position.X-bird.position.X, 2) + math.Pow(b.position.Y-bird.position.Y, 2))
			if distance <= b.rayonSeparation {
				b.separation[bird] = true
				b.alignement[bird] = true
				b.Cohesion[bird] = true
				//separation = append(separation, boid) // map plutot que tableau
			} else if distance > b.rayonSeparation && distance <= b.rayonVitesse {
				b.alignement[bird] = true
				b.Cohesion[bird] = true
				//alignement = append(alignement, boid)
			} else if distance > b.rayonVitesse && distance <= b.rayonCohesion {
				b.Cohesion[bird] = true
				//cohesion = append(cohesion, boid)
			}

		}
	}
}

//flock fonction
func (b *Bird) flock() {
	if b.evitement == 1 {
		b.vitesse.X = -b.vitesse.X
		b.vitesse.Y = -b.vitesse.Y
		b.evitement = 2
	} else if b.evitement == 0 {
		sep := b.separate()
		ali := b.align()
		coh := b.cohesion()
		r := b.random()
		b.ApplyForce(sep)
		b.ApplyForce(ali)
		b.ApplyForce(coh)
		b.ApplyForce(r)
	}
}

func (b *Bird) checkEdges(window settings.Window) {
	if b.position.X < 0 {
		b.position.X = 0
		b.vitesse.X = -b.vitesse.X
	} else if b.position.X > window.Width {
		b.position.X = window.Width
		b.vitesse.X = -b.vitesse.X
	}
	if b.position.Y < 0 {
		b.position.Y = 0
		b.vitesse.Y = -b.vitesse.Y
	} else if b.position.Y > window.Height {
		b.position.Y = window.Height
		b.vitesse.Y = -b.vitesse.Y
	}
}

//updatePosition fonction qui actualise la position de l'agent
func (b *Bird) updatePosition() {
	b.vitesse.X += b.acceleration.X
	b.vitesse.Y += b.acceleration.Y
	b.limit()
	b.position.X += b.vitesse.X
	b.position.Y += b.vitesse.Y
	b.acceleration.X = 0
	b.acceleration.Y = 0
	for k := range b.separation {
		b.separation[k] = false
	}
	for k := range b.alignement {
		b.alignement[k] = false
	}
	for k := range b.Cohesion {
		b.Cohesion[k] = false
	}
}

//Draw fonction permet d'afficher l'agent sur l'écran
func (b *Bird) Draw(screen *ebiten.Image, s settings.Settings) {
	op := ebiten.DrawImageOptions{}
	w, h := s.BirdImage.Size()
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(-1*math.Atan2(b.vitesse.Y*-1, b.vitesse.X) + math.Pi/2)
	op.GeoM.Translate(b.position.X, b.position.Y)

	switch b.ty {
	case 2:
		screen.DrawImage(s.BirdImage2, &op)
	default: //1 inclus
		screen.DrawImage(s.BirdImage, &op)
	}
}
