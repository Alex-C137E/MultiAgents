package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

//////////////////////////////////////
type Vector2D struct {
	X float64
	Y float64
}

func (v *Vector2D) Add(v2 Vector2D) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector2D) Subtract(v2 Vector2D) {
	v.X -= v2.X
	v.Y -= v2.Y
}

func (v *Vector2D) Limit(max float64) {
	magSq := v.MagnitudeSquared()
	if magSq > max*max {
		v.Divide(math.Sqrt(magSq))
		v.Multiply(max)
	}
}

func (v *Vector2D) Normalize() {
	mag := math.Sqrt(v.X*v.X + v.Y*v.Y)
	v.X /= mag
	v.Y /= mag
}

func (v *Vector2D) SetMagnitude(z float64) {
	v.Normalize()
	v.X *= z
	v.Y *= z
}

func (v *Vector2D) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2D) Divide(z float64) {
	v.X /= z
	v.Y /= z
}

func (v *Vector2D) Multiply(z float64) {
	v.X *= z
	v.Y *= z
}

func (v Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v2.X-v.X, 2) + math.Pow(v2.Y-v.Y, 2))
}

//////////////////////////////////////

type Boid struct {
	ID              int
	position        Vector2D
	vitesse         Vector2D
	acceleration    Vector2D
	longueur        float64
	rayonSeparation float64
	rayonCohesion   float64
	rayonVitesse    float64
	max_force       float64
	max_speed       Vector2D
	separation      map[*Boid]bool
	Cohesion        map[*Boid]bool
	alignement      map[*Boid]bool
	evitement       int
	ty              int
}

type Flock struct {
	boids []*Boid
}

type Game struct {
	flock  Flock
	inited bool
}

var (
	birdImage  *ebiten.Image
	birdImage2 *ebiten.Image
)

const (
	screenWidth  = 1000
	screenHeight = 1000
	maxForce     = 30
)

func init() {
	bird, _, err := ebitenutil.NewImageFromFile("t2.png", ebiten.FilterDefault)
	bird2, _, err := ebitenutil.NewImageFromFile("t3.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	w, h := bird.Size()
	birdImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	birdImage2, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 1)
	birdImage.DrawImage(bird, op)
	birdImage2.DrawImage(bird2, op)
}

func (b *Boid) limit() {
	// regarder la norme plutot
	norme := math.Sqrt(math.Pow(b.vitesse.X, 2) + math.Pow(b.vitesse.Y, 2))
	norme_max := math.Sqrt(math.Pow(b.max_speed.X, 2) + math.Pow(b.max_speed.Y, 2))
	if norme >= norme_max {
		if b.vitesse.X > 0 {
			b.vitesse.X = b.max_speed.X / 2
		} else {
			b.vitesse.X = -(b.max_speed.X / 2)
		}
		if b.vitesse.Y > 0 {
			b.vitesse.Y = b.max_speed.Y / 2
		} else {
			b.vitesse.Y = -(b.max_speed.Y / 2)
		}
	}
}

func (b *Boid) ApplyForce(force Vector2D) {
	b.acceleration.X += force.X
	b.acceleration.Y += force.Y
}

//////////Calcul acceleration
func (b *Boid) separate() (t Vector2D) {
	steer := Vector2D{X: 0, Y: 0}
	for boid, _ := range b.separation {
		steer.X = steer.X - boid.position.X
		steer.Y = steer.Y - boid.position.Y
	}
	if len(b.separation) != 0 {
		steer.X = steer.X / float64(len(b.separation))
		steer.Y = steer.Y / float64(len(b.separation))
	}
	//steer.X -= b.vitesse.X
	//steer.Y -= b.vitesse.Y

	return steer
}

/*
func (b *Boid) evite() (t Vector2D) {
	steer := Vector2D{X: 0, Y: 0}
	for boid, _ := range b.evitement {
		steer.X = (steer.X - boid.vitesse.X) * 2
		steer.Y = (steer.Y - boid.vitesse.Y) * 2
	}
	if len(b.separation) != 0 {
		steer.X = steer.X / float64(len(b.separation))
		steer.Y = steer.Y / float64(len(b.separation))
	}
	//steer.X -= b.vitesse.X
	//steer.Y -= b.vitesse.Y

	return steer
}
*/

func (b *Boid) align() (t Vector2D) {
	alignement := Vector2D{X: 0, Y: 0}
	for boid, _ := range b.alignement {
		alignement.X = alignement.X + boid.vitesse.X
		alignement.Y = alignement.Y + boid.vitesse.Y
	}
	if len(b.alignement) != 0 {
		alignement.X = alignement.X / float64(len(b.alignement))
		alignement.Y = alignement.Y / float64(len(b.alignement))
	}
	return alignement
}

func (b *Boid) cohesion() (t Vector2D) {
	cohesion := Vector2D{X: 0, Y: 0}
	for boid, _ := range b.Cohesion {
		cohesion.X = cohesion.X + boid.position.X
		cohesion.Y = cohesion.Y + boid.position.Y
	}
	if len(b.Cohesion) != 0 {
		cohesion.X = cohesion.X / float64(len(b.Cohesion))
		cohesion.Y = cohesion.Y / float64(len(b.Cohesion))
	}
	return cohesion
}

func (b *Boid) random() (t Vector2D) {
	r := Vector2D{X: 0, Y: 0}
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

func (b *Boid) Update() {
	b.vitesse.X += b.acceleration.X
	b.vitesse.Y += b.acceleration.Y
	b.limit()
	b.position.X += b.vitesse.X
	b.position.Y += b.vitesse.Y
	b.acceleration.X = 0
	b.acceleration.Y = 0
	for k, _ := range b.separation {
		b.separation[k] = false
	}
	for k, _ := range b.alignement {
		b.alignement[k] = false
	}
	for k, _ := range b.Cohesion {
		b.Cohesion[k] = false
	}
}

func (flock *Flock) Logic() {
	var wg sync.WaitGroup
	wg.Add(len(flock.boids))
	for _, boid := range flock.boids {
		boid.GetCouches(flock.boids)
		boid.flock()
		boid.CheckEdges()
		boid.Update()
		wg.Done()
	}
}

func (boid *Boid) CheckEdges() {

	if boid.position.X < 0 {
		boid.position.X = 0
		boid.vitesse.X = -boid.vitesse.X
	} else if boid.position.X > screenWidth {
		boid.position.X = screenWidth
		boid.vitesse.X = -boid.vitesse.X
	}
	if boid.position.Y < 0 {
		boid.position.Y = 0
		boid.vitesse.Y = -boid.vitesse.Y
	} else if boid.position.Y > screenHeight {
		boid.position.Y = screenHeight
		boid.vitesse.Y = -boid.vitesse.Y
	}

	/*
		if boid.position.X < 0 {
			boid.position.X = screenWidth
		} else if boid.position.X > screenWidth {
			boid.position.X = 0
		}
		if boid.position.Y < 0 {
			boid.position.Y = screenHeight
		} else if boid.position.Y > screenHeight {
			boid.position.Y = 0
		}
	*/
}

func (b *Boid) GetCouches(boids []*Boid) {
	r := 0
	for _, boid := range boids {
		if boid.ty != b.ty && (math.Sqrt(math.Pow(b.position.X-boid.position.X, 2)+math.Pow(b.position.Y-boid.position.Y, 2))) <= 20 && b.ty == 1 && b.evitement != 2 {
			b.evitement = 1
			r = 1
		}
		if (math.Sqrt(math.Pow(b.position.X-boid.position.X, 2)+math.Pow(b.position.Y-boid.position.Y, 2))) < 20 && b.evitement == 2 && boid.ty != b.ty && b.ty == 1 {
			r = 1
		}
	}
	if r == 0 {
		b.evitement = 0
	}
	// Couche separation, cohesion, alignement
	for _, boid := range boids {
		if boid.ty == b.ty && b.evitement == 0 {
			if boid.ID == b.ID {
				continue
			}
			distance := math.Sqrt(math.Pow(b.position.X-boid.position.X, 2) + math.Pow(b.position.Y-boid.position.Y, 2))
			if distance <= b.rayonSeparation {
				b.separation[boid] = true
				b.alignement[boid] = true
				b.Cohesion[boid] = true
				//separation = append(separation, boid) // map plutot que tableau
			} else if distance > b.rayonSeparation && distance <= b.rayonVitesse {
				b.alignement[boid] = true
				b.Cohesion[boid] = true
				//alignement = append(alignement, boid)
			} else if distance > b.rayonVitesse && distance <= b.rayonCohesion {
				b.Cohesion[boid] = true
				//cohesion = append(cohesion, boid)
			}

		}
	}
}

func (b *Boid) flock() {
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

	//b.ApplyForce(r)

}

func NewBoid(id int, position Vector2D, vitesse Vector2D, acceleration Vector2D, MaxS Vector2D, ty int) Boid {
	return Boid{
		id,
		position,
		vitesse,
		acceleration,
		5,
		32,
		30,
		27,
		maxForce,
		MaxS,
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		0,
		ty,
	}
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()
	numBoids := 100

	rand.Seed(time.Hour.Milliseconds())
	g.flock.boids = make([]*Boid, numBoids)
	for i := range g.flock.boids {
		w, h := birdImage.Size()
		////Position de départ
		x, y := rand.Float64()*float64(screenWidth-w), rand.Float64()*float64(screenWidth-h)
		//min, max := -maxForce, maxForce
		////Vitesse de départ
		t := rand.Float64()
		ty := 1
		var vx, vy float64
		vx = 2
		vy = 2
		if t > 0.25 {
			vx, vy = -2, -2
		}
		if t > 0.5 {
			vx, vy = -2, 2
			ty = 2
		}
		if t > 0.75 {
			vx, vy = 2, -2
		}

		b := NewBoid(i, Vector2D{X: x, Y: y}, Vector2D{X: float64(vx), Y: float64(vy)}, Vector2D{X: 0, Y: 0}, Vector2D{X: 10, Y: 10}, ty)
		g.flock.boids[i] = &b
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !g.inited {
		g.init()
	}

	g.flock.Logic()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	op := ebiten.DrawImageOptions{}
	w, h := birdImage.Size()
	for _, boid := range g.flock.boids {
		op.GeoM.Reset()
		//op.GeoM.Scale(0.8, 0.8)
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Rotate(-1*math.Atan2(boid.vitesse.Y*-1, boid.vitesse.X) + math.Pi/2)
		op.GeoM.Translate(boid.position.X, boid.position.Y)
		if boid.ty == 1 {
			screen.DrawImage(birdImage, &op)
		} else {
			screen.DrawImage(birdImage2, &op)
		}

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
