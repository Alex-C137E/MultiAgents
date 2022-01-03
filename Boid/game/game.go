package game

import (
	_ "image/png"
	"math/rand"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"

	agent "gitlab.utc.fr/projet_ia04/Boid/agent"
	flock "gitlab.utc.fr/projet_ia04/Boid/flock"
	utils "gitlab.utc.fr/projet_ia04/Boid/utils"
	worldelements "gitlab.utc.fr/projet_ia04/Boid/worldelements"
)

type Vector2D = utils.Vector2D

type Game struct {
	Flock     flock.Flock
	Sync      chan string
	musicInfo string

	currentLevel int
	levels       []*Level

	scores       []*Score
	currentScore Score

	graphics *Graphics

	initTime time.Time
	timeOut  int //entier corespondant au temps max de jeu (en minute)
}

func NewGame(c chan string, timeOut int) *Game {
	g := &Game{}
	g.Sync = c
	g.initTime = time.Now()
	g.timeOut = timeOut

	g.levels = LoadLevels() //Chargement des niveaux par défaut

	go func() {
		for {
			// lorsque l'agent  reçoit sur sa channel sync(bloquant): il reçoit une indication de la musique
			g.musicInfo = <-g.Sync
			// Il doit modifier un de ses paramêtres
		}
	}()

	g.graphics = NewGraphics()

	// Initialisation du jeu au niveau 0 (score = 0)
	g.setGame(0, 0)

	return g
}

//setGame fonctio qui permet de changer le niveau
func (g *Game) setGame(currentLevel int, initScore int) {
	g.currentLevel = currentLevel
	// Initialisation des variables vis à vis du niveau en cours:
	utils.RepulsionFactorBtwnSpecies = g.levels[g.currentLevel].RepulsionFactorBtwnSpecies
	utils.SeparationPerception = g.levels[g.currentLevel].SeparationPerception
	utils.CohesionPerception = g.levels[g.currentLevel].CohesionPerception
	utils.AlignPerception = g.levels[g.currentLevel].AlignPerception
	utils.NumWall = g.levels[g.currentLevel].numWall
	utils.MaxForce = g.levels[g.currentLevel].MaxForce
	utils.MaxSpeed = g.levels[g.currentLevel].MaxSpeed
	//Initialisation du score au niveau courant
	g.currentScore = *NewScore(g.currentLevel, 0, 1)
	// Initialisation du filet(polygon)
	g.graphics.ChangeMaxPolygonSize(g.levels[g.currentLevel].polygonSize)

	// Initialisation des agents:
	rand.Seed(time.Hour.Milliseconds())

	g.Flock.Boids = agent.GenerateBoids()
	g.Flock.Predators = agent.GeneratePredators(g.levels[g.currentLevel].SharkDensity)

	// Mise en place des murs/bombes en fonction du niveau
	wallIndex := 0
	g.Flock.Walls = make([]*worldelements.Wall, utils.NumWall)
	//bombe oeil droit:
	if g.currentLevel >= 1 {
		walls := worldelements.GenerateEyeWallBomb(utils.ScreenWidth*0.5, utils.ScreenHeight*0.4)

		for i := 0; i < len(walls); i++ {
			g.Flock.Walls[wallIndex+i] = walls[i]
		}
		wallIndex = wallIndex + len(walls)

		// bombe bouche
		walls = worldelements.GenerateMouthWallBomb(utils.ScreenWidth*0.4, utils.ScreenHeight*0.55)

		for i := 0; i < len(walls); i++ {
			g.Flock.Walls[wallIndex+i] = walls[i]
		}
		wallIndex = wallIndex + len(walls)
		//bombe oeil gauche:
		walls = worldelements.GenerateEyeWallBomb(utils.ScreenWidth*0.4, utils.ScreenHeight*0.4)

		for i := 0; i < len(walls); i++ {
			g.Flock.Walls[wallIndex+i] = walls[i]
		}
		wallIndex = wallIndex + len(walls)
	}
	if g.currentLevel > 3 {
		//Toit de Bombe:
		walls := worldelements.GenerateSideWallBomb(true)

		for i := 0; i < len(walls); i++ {
			g.Flock.Walls[wallIndex+i] = walls[i]
		}
		wallIndex = wallIndex + len(walls)

		//Sol de Bombe:
		walls = worldelements.GenerateSideWallBomb(false)

		for i := 0; i < len(walls); i++ {
			g.Flock.Walls[wallIndex+i] = walls[i]
		}
		wallIndex = wallIndex + len(walls)
	}

	if g.currentLevel < 4 {
		w, h := utils.SandImage.Size()
		for sx := 0; sx < 100; sx++ {
			g.Flock.Walls[wallIndex] = &worldelements.Wall{
				ImageWidth:  w,
				ImageHeight: h,
				Position:    Vector2D{X: (float64(h-5) * float64(sx)), Y: utils.ScreenHeight - 1},
				TypeWall:    1,
			}
			wallIndex++
		}
		for sx := 0; sx < 100; sx++ {
			g.Flock.Walls[wallIndex] = &worldelements.Wall{
				ImageWidth:  w,
				ImageHeight: h,
				Position:    Vector2D{X: (float64(h-5) * float64(sx)), Y: utils.ScreenHeight - (float64(h - 5))},
				TypeWall:    1,
			}
			wallIndex++
		}

		for sx := 0; sx < 100; sx++ {
			g.Flock.Walls[wallIndex] = &worldelements.Wall{
				ImageWidth:  w,
				ImageHeight: h,
				Position:    Vector2D{X: float64(h-5) * float64(sx), Y: -1},
				TypeWall:    2,
			}
			wallIndex++
		}
	}
}

func (g *Game) Update() error {
	// L'agent musique perturbe les agents boids afin de rendre le jeu plus complexe
	if g.musicInfo == "very hard drop" {
		utils.RepulsionFactorBtwnSpecies = 1000
		utils.SeparationPerception = 500
		utils.CohesionPerception = 10
	} else if g.musicInfo == "hard drop" {
		utils.RepulsionFactorBtwnSpecies = 800
		utils.SeparationPerception = 500
		utils.CohesionPerception = 100
	} else if g.musicInfo == "medium drop" {
		utils.RepulsionFactorBtwnSpecies = 500
		utils.SeparationPerception = 250
		utils.CohesionPerception = 200
	} else if g.musicInfo == "small drop" {
		utils.RepulsionFactorBtwnSpecies = 200
		utils.SeparationPerception = 100
		utils.CohesionPerception = 250
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
		utils.RepulsionFactorBtwnSpecies = g.levels[g.currentLevel].RepulsionFactorBtwnSpecies
		utils.SeparationPerception = g.levels[g.currentLevel].SeparationPerception
		utils.CohesionPerception = g.levels[g.currentLevel].CohesionPerception
	}

	g.Flock.Logic(g.currentLevel)

	g.graphics.UpdateGraphics(g.Flock.Boids, &g.currentScore) //Update lasso

	if g.nextLevel() {
		if g.currentLevel+1 == len(g.levels) {
		} else {
			g.scores = append(g.scores, &g.currentScore)
			g.setGame(g.currentLevel+1, g.currentScore.Value)
		}
	}

	return nil
}

//Draw fonction affiche tous les élements du monde
func (g *Game) Draw(screen *ebiten.Image) {
	g.graphics.DrawBackground(screen)
	g.graphics.DrawBoids(screen, g.Flock.Boids)
	g.graphics.DrawPredators(screen, g.Flock.Predators)
	g.graphics.DrawWalls(screen, g.Flock.Walls)
	if g.currentLevel == len(g.levels)-1 && g.nextLevel() {
		s := 0
		for _, val := range g.scores {
			s += val.Value
		}
		g.currentScore.Value = s
		g.graphics.DrawInterface(screen, g.currentScore, false, true)
	} else {
		g.graphics.DrawInterface(screen, g.currentScore, g.IsGameOver(), false)
	}
}

//Layout fonction qui retourne la taille de la fenêtre
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return utils.ScreenWidth, utils.ScreenHeight
}

func (g *Game) nextLevel() bool {
	for i := 0; i < len(g.Flock.Boids); i++ {
		if !g.Flock.Boids[i].Dead && g.currentScore.RequiredFishType == g.Flock.Boids[i].Species {
			return false
		}
	}
	return true
}

//IsGameOver fonction qui détermine si le jeu est terminé ou non.
func (g *Game) IsGameOver() bool {
	now := time.Now()
	return now.Sub(g.initTime) > time.Duration(g.timeOut)*time.Minute
}
