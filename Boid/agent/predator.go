package agent

import (
	"math"
	"math/rand"

	utils "gitlab.utc.fr/projet_ia04/Boid/utils"
	worldelements "gitlab.utc.fr/projet_ia04/Boid/worldelements"
)

type Predator struct {
	ImageWidth   int
	ImageHeight  int
	Position     utils.Vector2D
	Velocity     utils.Vector2D
	Acceleration utils.Vector2D
	Density      int
	Angle        int
	Dist         int
	V1           utils.Vector2D
	V2           utils.Vector2D
	R            bool
}

func GeneratePredators(sharkDensity int) []*Predator {
	predators := make([]*Predator, utils.NumPreda)

	for i := range predators {
		w, h := utils.BirdImage.Size()

		// Pour éviter que les agents apparaisent au dessus ou en dessous des murs de bombes:
		// on les fait apparaitre  horizontalement en ligne au mileu  de l'écran:
		middle := utils.ScreenHeight / 2
		x, y := rand.Float64()*float64(utils.ScreenWidth-w), float64(middle+100)

		min, max := -utils.MaxForce, utils.MaxForce
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		predators[i] = &Predator{
			ImageWidth:   w,
			ImageHeight:  h,
			Position:     utils.Vector2D{X: x, Y: y},
			Velocity:     utils.Vector2D{X: vx, Y: vy},
			Acceleration: utils.Vector2D{X: 0, Y: 0},
			Density:      sharkDensity,
			Dist:         400,
			Angle:        10,
			V1:           utils.Vector2D{X: 0, Y: 0},
			V2:           utils.Vector2D{X: 0, Y: 0},
			R:            false,
		}
	}
	return predators
}

//Update fonction qui permet de mettre à jour le prédateur
func (preda *Predator) Update(walls []*worldelements.Wall, boids []*Boid) {
	if !preda.CheckEdges() {
		preda.ApplyRules(boids)
		preda.CheckWalls(walls)
	}
	preda.ApplyMovement()
}

func (preda *Predator) Vision() []utils.Vector2D {
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
	l1x := x + float64(preda.Dist)*math.Cos(utils.AngleToRadians(angU))
	l1y := y + float64(preda.Dist)*math.Sin(utils.AngleToRadians(angU))

	l2x := x + float64(preda.Dist)*math.Cos(utils.AngleToRadians(angL))
	l2y := y + float64(preda.Dist)*math.Sin(utils.AngleToRadians(angL))

	p1 := utils.Vector2D{X: l1x, Y: l1y}
	p2 := utils.Vector2D{X: l2x, Y: l2y}

	mapP := make([]utils.Vector2D, 2)
	mapP[0] = p1
	mapP[1] = p2
	return mapP
}

func (preda *Predator) ApplyRules(restOfFlock []*Boid) {
	var dens int
	var newP Predator
	var new bool
	var densMax = 0
	var densMax2 = 0
	var proie1 *Boid
	var proie2 *Boid
	var vPoint2 []utils.Vector2D

	//Tuer les boid se situant à une distance inférieure à 10
	for _, Boid := range restOfFlock {
		if (preda.Position.Distance(Boid.Position)) < 10 {
			Boid.Dead = true
		}
	}

	//Calcule de la postion des points correspondant au champ de vision
	vPoint := preda.Vision()
	new = false
	//Calcule nouveau point si 1 des premiers se situe à l'exterieur de la carte
	if (vPoint[0].X > utils.ScreenWidth || vPoint[0].X < 0) || (vPoint[0].Y > utils.ScreenHeight || vPoint[0].Y < 0) || (vPoint[1].X > utils.ScreenWidth || vPoint[1].X < 0) || (vPoint[1].Y > utils.ScreenHeight || vPoint[1].Y < 0) {
		new = true
		//on crée également un objet predateur "fantome" (situer à X - ScreenWidth et Y + ScreenHeight) pour gérer la detection des boids de l'autre coté de la carte
		newP = *preda
		if vPoint[0].X > utils.ScreenWidth {
			newP.Position.X = preda.Position.X - utils.ScreenWidth
		} else if vPoint[0].X < 0 {
			newP.Position.X = preda.Position.X + utils.ScreenWidth
		}

		if vPoint[1].X > utils.ScreenWidth {
			newP.Position.X = preda.Position.X - utils.ScreenWidth
		} else if vPoint[1].X < 0 {
			newP.Position.X = preda.Position.X + utils.ScreenWidth
		}
	}
	preda.V1 = vPoint[0]
	preda.V2 = vPoint[1]
	if new {
		vPoint2 = newP.Vision()
		if vPoint[1].X > utils.ScreenWidth || vPoint[1].X < 0 {
			preda.V2 = vPoint2[1]

		}
		if vPoint[0].X > utils.ScreenWidth || vPoint[0].X < 0 {
			preda.V1 = vPoint2[0]
		}
	}

	//Pour chaque boid
	for _, Boid := range restOfFlock {
		if !Boid.Dead {
			dens = 0
			b := utils.PointInTriangle(Boid.Position, preda.Position, vPoint[0], vPoint[1])
			b2 := false
			//Si un point de vision dépasse les limites de la carte
			if new {
				b2 = utils.PointInTriangle(Boid.Position, newP.Position, vPoint2[0], vPoint2[1])
			}
			//Si il est dans le champ de vision du prédateur
			if b {
				//On calcule le nombre de voisin se situant à moins de 30 pixels de distance
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
			//Si un point de vision dépasse les limites de la carte
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
	//Si une "zone" à une densité supérieur à preda.Density
	if densMax >= densMax2 && densMax > preda.Density {
		//On modifie le vecteur vitesse du prédateur pour le faire diriger dans la zone
		Vit := utils.Vector2D{X: proie1.Position.X - preda.Position.X, Y: proie1.Position.Y - preda.Position.Y}
		Vit.Normalize()
		Vit.X = Vit.X * 10
		Vit.Y = Vit.Y * 10
		preda.Velocity = Vit

	} else if densMax2 > preda.Density {
		Vit := utils.Vector2D{X: proie2.Position.X - newP.Position.X, Y: proie2.Position.Y - newP.Position.Y}
		Vit.Normalize()
		Vit.X = Vit.X * 10
		Vit.Y = Vit.Y * 10
		preda.Velocity = Vit
	} else {
		//Sinon on le fait ralentire jusqu'à une vitesse de 1
		vit := preda.Velocity
		if vit.X > 1 || -vit.X > 1 {
			vit.X = vit.X * 0.98
		}
		if vit.Y > 1 || -vit.Y > 1 {
			vit.Y = vit.Y * 0.98
		}
		//Si sa vitesse est éguale à 1 et que aucune zone ne dépasse la Densité minimum, on fait tourner aléatoirement le vecteur vitesse du prédateur
		preda.Velocity = utils.Vector2D{X: vit.X, Y: vit.Y}
		if rand.Float64() < 0.01 {
			preda.Velocity = utils.Rotate(preda.Velocity, rand.Intn(10))
		}
		if rand.Float64() > 0.99 {
			preda.Velocity = utils.Rotate(preda.Velocity, -rand.Intn(10))
		}

	}
}

func (preda *Predator) CheckEdges() bool {
	if preda.Position.X < 0 {
		preda.Position.X = utils.ScreenWidth
	} else if preda.Position.X > utils.ScreenWidth {
		preda.Position.X = 0
	}

	return false
}

func (preda *Predator) CheckWalls(walls []*worldelements.Wall) {
	if preda.R {
		preda.Velocity.Normalize()
		preda.Velocity.X = preda.Velocity.X * 2
		preda.Velocity.Y = preda.Velocity.Y * 2
		preda.R = false
	} else {
		for _, wall := range walls {
			d := preda.Position.Distance(wall.Position)
			//Si le prédateur est à moins de 20px d'un mur
			if d <= 20 {
				//On lui fait faire demi tour (et sotir de la zone d'action du mur)
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
