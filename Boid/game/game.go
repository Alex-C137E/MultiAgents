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
	wall "gitlab.utc.fr/projet_ia04/Boid/agent/wall"
	flock "gitlab.utc.fr/projet_ia04/Boid/flock"
	constant "gitlab.utc.fr/projet_ia04/Boid/utils/constant"
	variable "gitlab.utc.fr/projet_ia04/Boid/utils/variable"
	vector "gitlab.utc.fr/projet_ia04/Boid/utils/vector"
)

type Vector2D = vector.Vector2D

type Game struct {
	Flock     flock.Flock
	Sync      chan string
	musicInfo string
	scoreFont font.Face

	currentLevel int
	levels       []*Level

	scores       []*Score
	currentScore Score

	polygon         []Vector2D
	polygonReleased string
	polygonSize     float64
	maxPolygonSize  float64

	initTime time.Time
	timeOut  int //entier corespondant au temps max de jeu (en minute)
}

func NewGame(c chan string, timeOut int) *Game {
	g := &Game{}
	g.Sync = c
	g.initTime = time.Now()
	g.timeOut = timeOut

	//Création de 5 niveaux:
	g.levels = make([]*Level, 5)

	// niveau 0: ce niveau est très simple:les poissons sont amenés à se regrouper et se stabiliser rapidement
	// rendant leur pêche facile.
	g.levels[0] = NewLevel(10000, 10, 300, 100, 300, 1, 4.0, 1000, 15)

	// niveau 1:  ce niveau est identique au niveau 0 à la différence que cette fois-ci les requins attaquent bien
	// plus: leur attribue density est diminué de tel sorte à ce qu’il attaque pour une quantité de poisson dans leur
	// champ d’attaque inférieur à celle du niveau 0. Le joueur doit donc être plus rapide pour ne pas se faire
	// manger tous ses poisons par les prédateurs.
	g.levels[1] = NewLevel(10000, 10, 300, 100, 16+10+300, 1, 4.0, 1000, 8)

	// niveau 2:  ce niveau est dans la continuité du niveau 1 mais ici: le facteur de cohésion diminue, et ceux de
	// répulsion intra et inter espèces augmentent ce qui diminue la stabilité dans le comportement des poissons les
	// rendant plus complexes à attraper: ils se regroupent moins, et se mélangent plus entre espèces. De plus, la
	// taille maximale du filet diminue.
	g.levels[2] = NewLevel(500, 100, 100, 75, 16+10+300, 2.0, 4.0, 700, 8)

	// niveau 3:  pour ce niveau, on reprend le niveau 2 et on augmente le niveau difficulté en réduisant le niveau
	// de stabilité de manière similaire à ce qui fut fait pour le niveau 2. En plus, on rend les prédateurs plus
	// agressifs en utilisant le même procédé utilisé dans le niveau 1.
	g.levels[3] = NewLevel(50, 100, 75, 75, 16+10+300, 2.0, 4.0, 700, 5)

	// niveau 4: pour ce niveau on reprend le niveau 3 et on  rajoute des murs/ bombes pour favoriser le chaos et
	// rendre plus difficile la tâche d’'attraper les poissons. De plus les poissons sont moins en cohésion et vont
	// plus vite et bien sûr, les prédateurs sont encore plus agressifs :).
	g.levels[4] = NewLevel(50, 100, 50, 75, 16+10+2*48, 2.0, 5.0, 700, 3)

	go func() {
		for {
			// lorsque l'agent  reçoit sur sa channel sync(bloquant): il reçoit une indication de la musique
			g.musicInfo = <-g.Sync
			// Il doit modifier un de ses paramêtres
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

	// Initialisation du jeu au niveau 0 (score = 0)
	g.setGame(0, 0)

	return g
}

func (g *Game) setGame(currentLevel int, initScore int) {
	g.currentLevel = currentLevel
	// Initialisation des variables vis à vis du niveau en cours:
	variable.RepulsionFactorBtwnSpecies = g.levels[g.currentLevel].RepulsionFactorBtwnSpecies
	variable.SeparationPerception = g.levels[g.currentLevel].SeparationPerception
	variable.CohesionPerception = g.levels[g.currentLevel].CohesionPerception
	variable.AlignPerception = g.levels[g.currentLevel].AlignPerception
	variable.NumWall = g.levels[g.currentLevel].numWall
	variable.MaxForce = g.levels[g.currentLevel].MaxForce
	variable.MaxSpeed = g.levels[g.currentLevel].MaxSpeed
	//Initialisation du score au niveau courant
	g.currentScore = *NewScore(g.currentLevel, initScore, 1)
	// Initialisation du filet(polygon)
	g.polygonReleased = "non"
	g.maxPolygonSize = g.levels[g.currentLevel].polygonSize

	// Initialisation des agents:
	rand.Seed(time.Hour.Milliseconds())
	g.Flock.Boids = make([]*boid.Boid, constant.NumBoids)
	g.Flock.Predators = make([]*boid.Predator, constant.NumPreda)
	g.Flock.Walls = make([]*wall.Wall, variable.NumWall)
	for i := range g.Flock.Boids {
		w, h := variable.FishImage1.Size()

		// Pour éviter que les agents apparaisent au dessus ou en dessous des murs de bombes:
		// on les fait apparaitre  horizontalement en ligne au mileu  de l'écran:
		middle := constant.ScreenHeight / 2
		x, y := rand.Float64()*float64(constant.ScreenWidth-w), float64(middle)

		min, max := -variable.MaxForce, variable.MaxForce
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		s := rand.Intn(constant.NumSpecies)
		g.Flock.Boids[i] = &boid.Boid{
			ImageWidth:     w,
			ImageHeight:    h,
			Position:       Vector2D{X: x, Y: y},
			Velocity:       Vector2D{X: vx, Y: vy},
			Acceleration:   Vector2D{X: 0, Y: 0},
			Species:        s,
			Dead:           false,
			EscapePredator: 80.0,
			Marqued:        false,
		}
	}
	for i := range g.Flock.Predators {
		w, h := variable.BirdImage.Size()

		// Pour éviter que les agents apparaisent au dessus ou en dessous des murs de bombes:
		// on les fait apparaitre  horizontalement en ligne au mileu  de l'écran:
		middle := constant.ScreenHeight / 2
		x, y := rand.Float64()*float64(constant.ScreenWidth-w), float64(middle+100)

		min, max := -variable.MaxForce, variable.MaxForce
		vx, vy := rand.Float64()*(max-min)+min, rand.Float64()*(max-min)+min
		g.Flock.Predators[i] = &boid.Predator{
			ImageWidth:   w,
			ImageHeight:  h,
			Position:     Vector2D{X: x, Y: y},
			Velocity:     Vector2D{X: vx, Y: vy},
			Acceleration: Vector2D{X: 0, Y: 0},
			Density:      g.levels[g.currentLevel].SharkDensity,
			Dist:         400,
			Angle:        10,
			V1:           Vector2D{X: 0, Y: 0},
			V2:           Vector2D{X: 0, Y: 0},
			R:            false,
		}
	}

	// Mise en place des murs/bombes en fonction du niveau
	wallIndex := 0
	//bombe oeil droit:
	if g.currentLevel >= 1 {
		wallIndex = g.eyeWallBomb(constant.ScreenWidth*0.5, constant.ScreenHeight*0.4, wallIndex)
		// bombe bouche
		wallIndex = g.mouthWallBomb(constant.ScreenWidth*0.4, constant.ScreenHeight*0.55, wallIndex)
		//bombe oeil gauche:
		wallIndex = g.eyeWallBomb(constant.ScreenWidth*0.4, constant.ScreenHeight*0.4, wallIndex)
	}
	if g.currentLevel > 3 {
		//Toit de Bombe:
		wallIndex = g.sideWallBomb(true, wallIndex)
		//Sol de Bombe:
		wallIndex = g.sideWallBomb(false, wallIndex)
	}

	if g.currentLevel < 4 {
		w, h := variable.SandImage.Size()
		for sx := 0; sx < 100; sx++ {
			g.Flock.Walls[wallIndex] = &wall.Wall{
				ImageWidth:  w,
				ImageHeight: h,
				Position:    Vector2D{X: (float64(h-5) * float64(sx)), Y: constant.ScreenHeight - 1},
				TypeWall:    1,
			}
			wallIndex++
		}
		for sx := 0; sx < 100; sx++ {
			g.Flock.Walls[wallIndex] = &wall.Wall{
				ImageWidth:  w,
				ImageHeight: h,
				Position:    Vector2D{X: (float64(h-5) * float64(sx)), Y: constant.ScreenHeight - (float64(h - 5))},
				TypeWall:    1,
			}
			wallIndex++
		}

		for sx := 0; sx < 100; sx++ {
			g.Flock.Walls[wallIndex] = &wall.Wall{
				ImageWidth:  w,
				ImageHeight: h,
				Position:    Vector2D{X: float64(h-5) * float64(sx), Y: -1},
				TypeWall:    2,
			}
			wallIndex++
		}
	}

	// Positionement des murs aléatoire: on garde au cas où...
	// for i := range g.Flock.Walls {
	// 	w, h := variable.WallImage.Size()
	// 	x, y := rand.Float64()*float64(constant.ScreenWidth-w+1000), rand.Float64()*float64(constant.ScreenWidth-h)
	// 	g.Flock.Walls[i] = &wall.Wall{
	// 		ImageWidth:  w,
	// 		ImageHeight: h,
	// 		Position:    Vector2D{X: x, Y: y},
	// 	}
	// }
}

func (g *Game) Update() error {
	// L'agent musique perturbe les agents boids afin de rendre le jeu plus complexe
	if g.musicInfo == "very hard drop" {
		variable.RepulsionFactorBtwnSpecies = 1000
		variable.SeparationPerception = 500
		variable.CohesionPerception = 10
	} else if g.musicInfo == "hard drop" {
		variable.RepulsionFactorBtwnSpecies = 800
		variable.SeparationPerception = 500
		variable.CohesionPerception = 100
	} else if g.musicInfo == "medium drop" {
		variable.RepulsionFactorBtwnSpecies = 500
		variable.SeparationPerception = 250
		variable.CohesionPerception = 200
	} else if g.musicInfo == "small drop" {
		variable.RepulsionFactorBtwnSpecies = 200
		variable.SeparationPerception = 100
		variable.CohesionPerception = 250
	} else if g.musicInfo == "1" { // raccourcis secret
		g.setGame(0, 0)
	} else if g.musicInfo == "2" {
		g.setGame(1, 0)
	} else if g.musicInfo == "3" {
		g.setGame(2, 0)
	} else if g.musicInfo == "4" {
		g.setGame(3, 0)
	} else if g.musicInfo == "5" {
		g.setGame(4, 0)
	} else { // g.musicInfo = "R"
		variable.RepulsionFactorBtwnSpecies = g.levels[g.currentLevel].RepulsionFactorBtwnSpecies
		variable.SeparationPerception = g.levels[g.currentLevel].SeparationPerception
		variable.CohesionPerception = g.levels[g.currentLevel].CohesionPerception
	}

	g.Flock.Logic(g.currentLevel)

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
		if g.polygonSize < g.maxPolygonSize {
			for i := 0; i < len(g.Flock.Boids); i++ {
				if !g.Flock.Boids[i].Dead && IsPointInPolygon(g.Flock.Boids[i].Position, g.polygon) {
					g.Flock.Boids[i].Dead = true
					g.currentScore.AddCollectedFish(g.Flock.Boids[i].Species)
				}
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
	w, h := variable.FishImage1.Size()
	for _, boid := range g.Flock.Boids {
		if !boid.Dead {
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Rotate(-1*math.Atan2(boid.Velocity.Y*-1, boid.Velocity.X) + math.Pi)
			op.GeoM.Translate(boid.Position.X, boid.Position.Y)

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
		if wall.TypeWall == 0 {
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			// op.GeoM.Rotate(math.Pi / 2)
			op.GeoM.Translate(wall.Position.X, wall.Position.Y)
			screen.DrawImage(variable.WallImage, &op)
		} else if wall.TypeWall == 1 {
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			// op.GeoM.Rotate(math.Pi / 2)
			op.GeoM.Translate(wall.Position.X, wall.Position.Y)
			screen.DrawImage(variable.SandImage, &op)
		}
	}

	//Draw GUI
	text.Draw(screen, "Level: "+strconv.Itoa(g.currentScore.Level), g.scoreFont, 32, 32, color.White)
	text.Draw(screen, "Score: "+strconv.Itoa(g.currentScore.Value), g.scoreFont, constant.ScreenWidth/2-50, 32, color.White)
	text.Draw(screen, "Lasso size: "+strconv.Itoa(int(g.polygonSize)), g.scoreFont, constant.ScreenWidth/2-100, constant.ScreenHeight-32, color.White)
	text.Draw(screen, "Max Lasso size: "+strconv.Itoa(int(g.maxPolygonSize)), g.scoreFont, 32, constant.ScreenHeight-32, color.White)

	if g.polygonSize > g.maxPolygonSize {
		text.Draw(screen, "Your lasso is too big!", g.scoreFont, constant.ScreenWidth/2-100, constant.ScreenHeight/2-32, color.RGBA{255, 12, 26, 255})
	}

	if g.nextLevel() {
		if g.currentLevel+1 == len(g.levels) {
			text.Draw(screen, "Well Done, you won with the score: "+strconv.Itoa(g.currentScore.Value), g.scoreFont, constant.ScreenWidth/2-50, constant.ScreenHeight/2+132, color.RGBA{53, 223, 26, 255})
		} else {
			g.scores = append(g.scores, &g.currentScore)
			g.setGame(g.currentLevel+1, g.currentScore.Value)
		}
	}

	if g.IsGameOver() {
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
	return true
	// pour simplifier supprimer code qui suit (sinon le filet fonctionne mal)
	// inside := false
	// i := 0
	// j := len(polygon) - 1
	// for i < len(polygon) {
	// 	if (polygon[i].Y > p.Y) != (polygon[j].Y > p.Y) && p.X < (polygon[j].X-polygon[i].X)*(p.Y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X {
	// 		inside = !inside
	// 	}
	// 	i++
	// 	j = i
	// }
	// return inside
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

func (g *Game) nextLevel() bool {
	for i := 0; i < len(g.Flock.Boids); i++ {
		if !g.Flock.Boids[i].Dead && g.currentScore.RequiredFishType == g.Flock.Boids[i].Species {
			return false
		}
	}
	return true
}

func (g *Game) IsGameOver() bool {
	now := time.Now()
	return now.Sub(g.initTime) > time.Duration(g.timeOut)*time.Minute
}

func (g *Game) mouthWallBomb(xPos float64, yPos float64, currentWallIndex int) int {
	w, h := variable.WallImage.Size()
	wallIndex := currentWallIndex
	changeInc := -1
	inc, subInc := 0, 0
	for b := 0; b < 10; b++ {
		x := xPos + float64(b*w)
		if inc == 4 && changeInc == -1 {
			changeInc = 0
		}
		if subInc == 3 && changeInc != 1 {
			changeInc = 1
			inc = 4
		}
		if changeInc == -1 {
			inc++
		} else if changeInc == 1 {
			inc--
		} else {
			subInc++
		}
		y := yPos + float64(inc*h)
		g.Flock.Walls[wallIndex] = &wall.Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    Vector2D{X: x, Y: y},
			TypeWall:    0,
		}
		wallIndex++
	}
	return wallIndex
}

func (g *Game) eyeWallBomb(xPos float64, yPos float64, currentWallIndex int) int {
	w, h := variable.WallImage.Size()
	wallIndex := currentWallIndex
	changeInc := false
	inc := 1
	for b := 0; b < 5; b++ {
		x := xPos + float64(b*w)
		if b%3 == 0 {
			changeInc = true
		}
		if inc == 1 {
			changeInc = false
		}
		if changeInc {
			inc--
		} else {
			inc++
		}
		y := yPos - float64(inc*h)
		g.Flock.Walls[wallIndex] = &wall.Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    Vector2D{X: x, Y: y},
		}
		wallIndex++
	}
	changeInc = false
	inc = 2
	for b := 0; b < 3; b++ {
		x := xPos + float64((b+1)*w)
		if b%2 == 0 {
			changeInc = true
		}
		if inc == 1 {
			changeInc = false
		}
		if changeInc {
			inc--
		} else {
			inc++
		}
		y := yPos + float64((inc-2)*h)
		g.Flock.Walls[wallIndex] = &wall.Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    Vector2D{X: x, Y: y},
			TypeWall:    0,
		}
		wallIndex++
	}
	return wallIndex
}

func (g *Game) sideWallBomb(top bool, currentWallIndex int) int {
	w, h := variable.WallImage.Size()
	wallIndex := currentWallIndex
	//pour le mur du bas:
	changeInc := true
	inc := 4
	//pour le mur du haut
	if top {
		changeInc = false
		inc = 1
	}
	for b := 0; b < 48; b++ {
		x := b*w + 10
		if b%4 == 0 {
			changeInc = true
		}
		if inc == 1 {
			changeInc = false
		}
		if changeInc {
			inc--
		} else {
			inc++
		}
		y := -inc*h + constant.ScreenHeight
		if top {
			y = inc * h
		}
		g.Flock.Walls[wallIndex] = &wall.Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    Vector2D{X: float64(x), Y: float64(y)},
			TypeWall:    0,
		}
		wallIndex++
	}
	return wallIndex
}
