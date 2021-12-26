package predator

import (
	"math"
	"math/rand"

	boid "gitlab.utc.fr/projet_ia04/Boid/agent/boid"
	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
	constant "gitlab.utc.fr/projet_ia04/Boid/utils/constant"
	vector "gitlab.utc.fr/projet_ia04/Boid/utils/vector"
)

type Vector2D = vector.Vector2D

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

//Permet de faire une rotation de la matrice v d'un ang (en degré)
func Rotate(v Vector2D, ang int) Vector2D {
	aR := AngleToRadians(ang)
	oldX := v.X
	oldY := v.Y
	v.X = oldX*math.Cos(aR) - oldY*math.Sin(aR)
	v.Y = oldX*math.Sin(aR) + oldY*math.Cos(aR)
	return v
}

func Sign(p1 Vector2D, p2 Vector2D, p3 Vector2D) float64 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

//Check si un point est un triangle
func PointInTriangle(pt Vector2D, v1 Vector2D, v2 Vector2D, v3 Vector2D) bool {

	d1 := Sign(pt, v1, v2)
	d2 := Sign(pt, v2, v3)
	d3 := Sign(pt, v3, v1)

	has_neg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	has_pos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(has_neg && has_pos)
}

// Convertie un angle en radian
func AngleToRadians(angle int) float64 {
	return (math.Pi / 180) * float64(angle)
}

//Permet de créer le triangle correspondant au champ de vision
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
	l1x := x + float64(preda.Dist)*math.Cos(AngleToRadians(angU))
	l1y := y + float64(preda.Dist)*math.Sin(AngleToRadians(angU))

	l2x := x + float64(preda.Dist)*math.Cos(AngleToRadians(angL))
	l2y := y + float64(preda.Dist)*math.Sin(AngleToRadians(angL))

	p1 := Vector2D{X: l1x, Y: l1y}
	p2 := Vector2D{X: l2x, Y: l2y}

	mapP := make([]Vector2D, 2)
	mapP[0] = p1
	mapP[1] = p2
	return mapP
}

func (preda *Predator) ApplyRules(restOfFlock []*boid.Boid) {

	//preda.V1 = vPoint[0]
	//preda.V2 = vPoint[1]
	var dens int
	var newP Predator
	var new bool
	var densMax = 0
	var densMax2 = 0
	var proie1 *boid.Boid
	var proie2 *boid.Boid
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
			b := PointInTriangle(Boid.Position, preda.Position, vPoint[0], vPoint[1])
			b2 := false
			if new {
				b2 = PointInTriangle(Boid.Position, newP.Position, vPoint2[0], vPoint2[1])
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
			preda.Velocity = Rotate(preda.Velocity, rand.Intn(10))
		}
		if rand.Float64() > 0.99 {
			preda.Velocity = Rotate(preda.Velocity, -rand.Intn(10))
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
			if d <= 30 {
				preda.Velocity.Normalize()
				preda.Velocity.X = -preda.Velocity.X * (30 - d + 1)
				preda.Velocity.Y = -preda.Velocity.Y * (30 - d + 1)
				preda.R = true
			}
		}
	}
}

func (preda *Predator) ApplyMovement() {
	preda.Position.Add(preda.Velocity)
}
