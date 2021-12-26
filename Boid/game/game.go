package game

import (
	"fmt"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	boid "gitlab.utc.fr/projet_ia04/Boid/agent/boid"
	predator "gitlab.utc.fr/projet_ia04/Boid/agent/predator"
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
	scoreFont font.Face

	scores       []Score
	currentScore Score

	polygon         []Vector2D
	polygonReleased string
	polygonSize     float64
	maxPolygonSize  float64
}

func NewGame(c chan string) *Game {
	g := &Game{}
	g.Sync = c

	g.polygonReleased = "non"
	g.maxPolygonSize = 700.0

	rand.Seed(time.Hour.Milliseconds())
	g.Flock.Boids = make([]*boid.Boid, constant.NumBoids)
	g.Flock.Walls = make([]*wall.Wall, constant.NumWalls)
	g.Flock.Predators = make([]*predator.Predator, constant.NumPreda)

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
			Dead:         false,
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
	for i := range g.Flock.Predators {
		w, h := variable.BirdImage.Size()
		x, y := rand.Float64()*float64(constant.ScreenWidth-w), rand.Float64()*float64(constant.ScreenWidth-h)
		min, max := -constant.MaxForce, constant.MaxForce
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		g.Flock.Predators[i] = &predator.Predator{
			ImageWidth:   w,
			ImageHeight:  h,
			Position:     Vector2D{X: x, Y: y},
			Velocity:     Vector2D{X: vx, Y: vy},
			Acceleration: Vector2D{X: 0, Y: 0},
			Density:      5,
			Dist:         400,
			Angle:        10,
			V1:           Vector2D{X: 0, Y: 0},
			V2:           Vector2D{X: 0, Y: 0},
			R:            false,
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

	data, err := ioutil.ReadFile("Roboto-Black.ttf")
	if err != nil {
		fmt.Println(err)
	}

	ttf, err := truetype.Parse(data)
	if err != nil {
		fmt.Println(err)
	}

	op := truetype.Options{Size: 24, DPI: 72, Hinting: font.HintingFull}
	g.scoreFont = truetype.NewFace(ttf, &op)

	//Init scoring
	g.currentScore = *NewScore(1, 1)

	return g
}

func (g *Game) Update() error {
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

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		point := Vector2D{X: float64(mx), Y: float64(my)}
		g.polygon = append(g.polygon, point)
		g.polygonReleased = "en cours"
	} else if len(g.polygon) > 0 {
		g.polygonReleased = "pret"
	}

	if g.polygonReleased != "no" {
		g.polygonSize = GetPolygonSize(g.polygon)
	}

	if g.polygonReleased == "pret" {
		for i := 0; i < len(g.Flock.Boids); i++ {
			if g.Flock.Boids[i].Dead == false && IsPointInPolygon(g.Flock.Boids[i].Position, g.polygon) {
				g.Flock.Boids[i].Dead = true
				g.currentScore.AddCollectedFish(g.Flock.Boids[i].Species)
			}
		}
		g.polygonReleased = "non" //reset polygon
		g.polygon = make([]Vector2D, 0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	screen.DrawImage(variable.BackgroundImage, &op)
	w, h := variable.BirdImage.Size()
	for _, boid := range g.Flock.Boids {
		if boid.Dead == false {
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Rotate(-1*math.Atan2(boid.Velocity.Y*-1, boid.Velocity.X) + math.Pi)
			op.GeoM.Translate(boid.Position.X, boid.Position.Y)
			// A la base les images sont des chevrons, mais il m'est impossible de changer proprement leur
			// couleur en fonction de leur espèce avec r,g,b et  op.ColorM.Scale(r, g, b, 1):
			// l'idée est donc de leur donner l'apparence de rectangles colorés que je peux remplir selon une couleur
			// donnée

			// r := 0.0
			// g := 0.0
			// b := 0.0
			if boid.Species == 0 {
				screen.DrawImage(variable.FishImage1, &op)
			} else if boid.Species == 1 {
				screen.DrawImage(variable.FishImage2, &op)
			} else {
				screen.DrawImage(variable.FishImage3, &op)
			}
		}
	}

	for _, preda := range g.Flock.Predators {
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Rotate(-1*math.Atan2(preda.Velocity.Y*-1, preda.Velocity.X) + math.Pi)
		op.GeoM.Translate(preda.Position.X, preda.Position.Y)

		screen.DrawImage(variable.PredImage, &op)

		op.GeoM.Reset()
		//op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Translate(preda.V2.X, preda.V2.Y)
		screen.DrawImage(variable.BirdImage, &op)

		op.GeoM.Reset()
		//op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Translate(preda.V1.X, preda.V1.Y)
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

	//Draw GUI
	text.Draw(screen, "Level: "+strconv.Itoa(g.currentScore.Level), g.scoreFont, 32, 32, color.White)
	text.Draw(screen, "Score: "+strconv.Itoa(g.currentScore.Value), g.scoreFont, constant.ScreenWidth/2-50, 32, color.White)
	text.Draw(screen, "Lasso size: "+strconv.Itoa(int(g.polygonSize)), g.scoreFont, constant.ScreenWidth/2-100, constant.ScreenHeight-32, color.White)
	text.Draw(screen, "Max Lasso size: "+strconv.Itoa(int(g.maxPolygonSize)), g.scoreFont, 32, constant.ScreenHeight-32, color.White)

	if g.polygonSize > g.maxPolygonSize {
		text.Draw(screen, "Your lasso is too big!", g.scoreFont, constant.ScreenWidth/2-100, constant.ScreenHeight/2-32, color.RGBA{255, 12, 26, 255})
	}

	if g.IsGameOver() == true {
		text.Draw(screen, "Game over", g.scoreFont, constant.ScreenWidth/2-50, constant.ScreenHeight/2+132, color.RGBA{53, 223, 26, 255})
	}

	//Draw polygon
	for i := 1; i < len(g.polygon); i++ {
		ebitenutil.DrawLine(screen, g.polygon[i-1].X, g.polygon[i-1].Y, g.polygon[i].X, g.polygon[i].Y, color.RGBA{120, 12, 200, 255})
	}

	if g.polygonReleased == "pret" || g.polygonReleased == "en cours" {
		ebitenutil.DrawLine(screen, g.polygon[0].X, g.polygon[0].Y, g.polygon[len(g.polygon)-1].X, g.polygon[len(g.polygon)-1].Y, color.RGBA{120, 12, 200, 255}) //link fist and last to complete shape
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constant.ScreenWidth, constant.ScreenHeight
}

func IsPointInPolygon(p Vector2D, polygon []Vector2D) bool {
	minX := polygon[0].X
	maxX := polygon[0].X
	minY := polygon[0].Y
	maxY := polygon[0].Y
	for i := 1; i < len(polygon); i++ {
		minX = math.Min(polygon[i].X, minX)
		maxX = math.Max(polygon[i].X, maxX)
		minY = math.Min(polygon[i].Y, minY)
		maxY = math.Max(polygon[i].Y, maxY)
	}

	if p.X < minX || p.X > maxX || p.Y < minY || p.Y > maxY {
		return false
	}

	inside := false

	i := 0
	j := len(polygon) - 1
	for i < len(polygon) {
		if (polygon[i].Y > p.Y) != (polygon[j].Y > p.Y) && p.X < (polygon[j].X-polygon[i].X)*(p.Y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X {
			inside = !inside
		}
		i++
		j = i
	}

	return inside
}

func GetPolygonSize(polygon []Vector2D) float64 {
	if len(polygon) < 2 {
		return 0.0
	}
	sumDistance := 0.0
	for i := 1; i < len(polygon); i++ {
		p1 := polygon[i-1]
		p2 := polygon[i]
		distance := math.Sqrt((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
		sumDistance = sumDistance + distance
	}
	//add last distance
	p0 := polygon[0]
	pn := polygon[len(polygon)-1]
	return sumDistance + math.Sqrt((p0.X-pn.X)*(p0.X-pn.X)+(p0.Y-pn.Y)*(p0.Y-pn.Y))
}

func (g *Game) IsGameOver() bool {
	for i := 0; i < len(g.Flock.Boids); i++ {
		if g.Flock.Boids[i].Dead == false && g.currentScore.RequiredFishType == g.Flock.Boids[i].Species {
			return false
		}
	}
	return true
}
