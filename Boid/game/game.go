package game

import (
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	boid "gitlab.utc.fr/projet_ia04/Boid/agent/boid"
	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
	flock "gitlab.utc.fr/projet_ia04/Boid/flock"
	constant "gitlab.utc.fr/projet_ia04/Boid/utils/constant"
	variable "gitlab.utc.fr/projet_ia04/Boid/utils/variable"
	vector "gitlab.utc.fr/projet_ia04/Boid/utils/vector"
)

type Vector2D = vector.Vector2D

type Game struct {
	Flock flock.Flock
	// Inited    bool
	Sync      chan string
	musicInfo string
}

func NewGame(c chan string) *Game {
	g := &Game{}
	g.Sync = c

	rand.Seed(time.Hour.Milliseconds())
	g.Flock.Boids = make([]*boid.Boid, constant.NumBoids)
	g.Flock.Walls = make([]*wall.Wall, constant.NumWalls)
	for i := range g.Flock.Boids {
		w, h := variable.BirdImage.Size()
		x, y := rand.Float64()*float64(constant.ScreenWidth-w), rand.Float64()*float64(constant.ScreenWidth-h)
		min, max := -constant.MaxForce, constant.MaxForce
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		s := rand.Intn(constant.NumSpecies)
		g.Flock.Boids[i] = &boid.Boid{
			ImageWidth:   w,
			ImageHeight:  h,
			Position:     Vector2D{X: x, Y: y},
			Velocity:     Vector2D{X: vx, Y: vy},
			Acceleration: Vector2D{X: 0, Y: 0},
			Species:      s,
		}
	}
	for i := range g.Flock.Walls {
		w, h := variable.WallImage.Size()
		x, y := rand.Float64()*float64(constant.ScreenWidth-w+1000), rand.Float64()*float64(constant.ScreenWidth-h)
		g.Flock.Walls[i] = &wall.Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    Vector2D{X: x, Y: y},
		}
	}
	// Variable Initialisations
	variable.RepulsionFactorBtwnSpecies = 100
	variable.SeparationPerception = 50
	variable.CohesionPerception = 300

	go func() {
		for {
			// lorsque l'agent  reçoit sur sa channel sync(bloquant): il reçoit une indication de la musique
			g.musicInfo = <-g.Sync
			// Il doit modifier un de ses paramêtre
		}
	}()

	return g
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.musicInfo == "very hard drop" {
		variable.RepulsionFactorBtwnSpecies = 500
		variable.SeparationPerception = 250
		variable.CohesionPerception = 100
	} else if g.musicInfo == "hard drop" {
		variable.RepulsionFactorBtwnSpecies = 400
		variable.SeparationPerception = 200
		variable.CohesionPerception = 150
	} else if g.musicInfo == "medium drop" {
		variable.RepulsionFactorBtwnSpecies = 300
		variable.SeparationPerception = 150
		variable.CohesionPerception = 200
	} else if g.musicInfo == "small drop" {
		variable.RepulsionFactorBtwnSpecies = 200
		variable.SeparationPerception = 100
		variable.CohesionPerception = 250
	} else {
		variable.RepulsionFactorBtwnSpecies = 100
		variable.SeparationPerception = 50
		variable.CohesionPerception = 300
	}

	g.Flock.Logic()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	op := ebiten.DrawImageOptions{}
	w, h := variable.BirdImage.Size()
	for _, boid := range g.Flock.Boids {
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Rotate(-1*math.Atan2(boid.Velocity.Y*-1, boid.Velocity.X) + math.Pi/2)
		op.GeoM.Translate(boid.Position.X, boid.Position.Y)
		// A la base les images sont des chevrons, mais il m'est impossible de changer proprement leur
		// couleur en fonction de leur espèce avec r,g,b et  op.ColorM.Scale(r, g, b, 1):
		// l'idée est donc de leur donner l'apparence de rectangles colorés que je peux remplir selon une couleur
		// donnée

		// r := 0.0
		// g := 0.0
		// b := 0.0
		if boid.Species == 0 {
			// l'espèce 0 est rouge
			variable.BirdImage.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
			// r = 255
		} else if boid.Species == 1 {
			// l'espèce 1 est bleu
			variable.BirdImage.Fill(color.NRGBA{0x00, 0x00, 0xff, 0xff})
			// g = 255
		} else {
			// l'espèce 2 est verte
			variable.BirdImage.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})
			// g = 255
		}
		// op.ColorM.Scale(r, g, b, 1)
		screen.DrawImage(variable.BirdImage, &op)
	}
	w, h = variable.WallImage.Size()
	for _, wall := range g.Flock.Walls {
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		// op.GeoM.Rotate(math.Pi / 2)
		op.GeoM.Translate(wall.Position.X, wall.Position.Y)
		variable.WallImage.Fill(color.Black)
		screen.DrawImage(variable.WallImage, &op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constant.ScreenWidth, constant.ScreenHeight
}
