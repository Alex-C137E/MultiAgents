package game

import (
	"fmt"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"math"
	"strconv"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	ebitenutil "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	text "github.com/hajimehoshi/ebiten/v2/text"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	agent "gitlab.utc.fr/projet_ia04/Boid/agent"
	utils "gitlab.utc.fr/projet_ia04/Boid/utils"
	worldelements "gitlab.utc.fr/projet_ia04/Boid/worldelements"
)

//Graphics struct qui permet de prendre en charge toute la représentation graphique du jeu
type Graphics struct {
	scoreFont font.Face //Police utilisée pour l'interface

	polygon         []utils.Vector2D //Les points du lasso
	polygonReleased string           //L'état du lasso ("non" : pret à en créer un nouveau, "en cours" : en train d'etre fait par l'utilateur, "pret" : terminé, doit être remis à 0)
	polygonSize     float64          //Taille actuelle du lasso
	maxPolygonSize  float64          //Taille maximùale du lasso
}

//NewGraphics fonction constructeur du type Graphics
func NewGraphics() *Graphics {
	data, err := ioutil.ReadFile("Roboto-Black.ttf")
	if err != nil {
		fmt.Println(err)
	}

	ttf, err := truetype.Parse(data)
	if err != nil {
		fmt.Println(err)
	}

	op := truetype.Options{Size: 24, DPI: 72, Hinting: font.HintingFull}
	scoreFont := truetype.NewFace(ttf, &op)

	var emptyVectors []utils.Vector2D

	return &Graphics{scoreFont, emptyVectors, "non", 0.0, 0.0}
}

//ChangeMaxPolygonSize fonction qui permet de modifier la taille maximale du lasso
func (gra *Graphics) ChangeMaxPolygonSize(newMaxPolygonSize float64) {
	gra.polygonReleased = "non"
	gra.maxPolygonSize = newMaxPolygonSize
}

//UpdateGraphics Fonction permet de mettre à jour le lasso
func (gra *Graphics) UpdateGraphics(boids []*agent.Boid, score *Score) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) { //on est en train de cliquer
		mx, my := ebiten.CursorPosition()
		point := utils.Vector2D{X: float64(mx), Y: float64(my)}
		gra.polygon = append(gra.polygon, point) //on ajoute le point au lasso
		gra.polygonReleased = "en cours"
	} else if len(gra.polygon) > 0 { //on relache le click
		gra.polygonReleased = "pret" //le lasso est complet
	}

	if gra.polygonReleased != "no" {
		gra.polygonSize = utils.GetPolygonSize(gra.polygon)
	}

	if gra.polygonReleased == "pret" { //si le polygon est pret on peut attrapper des poissons
		if gra.polygonSize < gra.maxPolygonSize {
			for i := 0; i < len(boids); i++ {
				if !boids[i].Dead && utils.IsPointInPolygon(boids[i].Position, gra.polygon) {
					boids[i].Dead = true
					score.AddCollectedFish(boids[i].Species)
				}
			}
		}
		gra.polygonReleased = "non" //reset polygon
		gra.polygon = make([]Vector2D, 0)
	}
}

//DrawBackground fonction qui permet d'afficher le fond
func (gra *Graphics) DrawBackground(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	screen.DrawImage(utils.BackgroundImage, &op)
}

//DrawBoids fonction qui permet d'afficher les poissons
func (gra *Graphics) DrawBoids(screen *ebiten.Image, boids []*agent.Boid) {
	op := ebiten.DrawImageOptions{}
	w, h := utils.FishImage1.Size()
	for _, boid := range boids {
		if !boid.Dead {
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Rotate(-1*math.Atan2(boid.Velocity.Y*-1, boid.Velocity.X) + math.Pi)
			op.GeoM.Translate(boid.Position.X, boid.Position.Y)
			if boid.Species == 0 {
				screen.DrawImage(utils.FishImage1, &op)
			} else if boid.Species == 1 {
				screen.DrawImage(utils.FishImage2, &op)
			} else {
				screen.DrawImage(utils.FishImage3, &op)
			}
		}
	}
}

//DrawPredators fonction qui permet d'afficher les prédateurs
func (gra *Graphics) DrawPredators(screen *ebiten.Image, predators []*agent.Predator) {
	op := ebiten.DrawImageOptions{}
	w, h := utils.PredImage.Size()
	for _, preda := range predators {
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Rotate(-1*math.Atan2(preda.Velocity.Y*-1, preda.Velocity.X) + math.Pi)
		op.GeoM.Translate(preda.Position.X, preda.Position.Y)

		screen.DrawImage(utils.PredImage, &op)

		op.GeoM.Reset()
		//op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Translate(preda.V2.X, preda.V2.Y)
		screen.DrawImage(utils.BirdImage, &op)

		op.GeoM.Reset()
		//op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Translate(preda.V1.X, preda.V1.Y)
		screen.DrawImage(utils.BirdImage, &op)
	}
}

//DrawWalls fonction qui permet d'afficher les murs
func (gra *Graphics) DrawWalls(screen *ebiten.Image, walls []*worldelements.Wall) {
	op := ebiten.DrawImageOptions{}
	w, h := utils.WallImage.Size()
	for _, wall := range walls {
		if wall.TypeWall == 0 {
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Translate(wall.Position.X, wall.Position.Y)
			screen.DrawImage(utils.WallImage, &op)
		} else if wall.TypeWall == 1 {
			op.GeoM.Reset()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Translate(wall.Position.X, wall.Position.Y)
			screen.DrawImage(utils.SandImage, &op)
		}
	}
}

//DrawInterface fonction qui permet d'afficher les différents textes de l'interface.
func (g *Graphics) DrawInterface(screen *ebiten.Image, score Score, gameOver bool, end bool) {
	//Draw GUI
	text.Draw(screen, "Level: "+strconv.Itoa(score.Level+1), g.scoreFont, 32, 32, color.RGBA{100, 100, 100, 255})
	text.Draw(screen, "Score: "+strconv.Itoa(score.Value), g.scoreFont, utils.ScreenWidth/2-50, 32, color.RGBA{100, 100, 100, 255})
	text.Draw(screen, "Lasso size: "+strconv.Itoa(int(g.polygonSize)), g.scoreFont, utils.ScreenWidth/2-100, utils.ScreenHeight-32, color.RGBA{100, 100, 100, 255})
	text.Draw(screen, "Max Lasso size: "+strconv.Itoa(int(g.maxPolygonSize)), g.scoreFont, 32, utils.ScreenHeight-32, color.RGBA{100, 100, 100, 255})

	if g.polygonSize > g.maxPolygonSize {
		text.Draw(screen, "Your lasso is too big!", g.scoreFont, utils.ScreenWidth/2-100, utils.ScreenHeight/2-32, color.RGBA{255, 12, 26, 255})
	}

	if gameOver {
		text.Draw(screen, "Game over", g.scoreFont, utils.ScreenWidth/2-50, utils.ScreenHeight/2+132, color.RGBA{53, 223, 26, 255})
	}

	if end {
		text.Draw(screen, "you won: ", g.scoreFont, utils.ScreenWidth/2-50, utils.ScreenHeight/2+132, color.RGBA{53, 223, 26, 255})
	}

	//Draw polygon
	for i := 1; i < len(g.polygon); i++ {
		ebitenutil.DrawLine(screen, g.polygon[i-1].X, g.polygon[i-1].Y, g.polygon[i].X, g.polygon[i].Y, color.RGBA{120, 12, 200, 255})
	}

	if g.polygonReleased == "pret" || g.polygonReleased == "en cours" {
		ebitenutil.DrawLine(screen, g.polygon[0].X, g.polygon[0].Y, g.polygon[len(g.polygon)-1].X, g.polygon[len(g.polygon)-1].Y, color.RGBA{120, 12, 200, 255}) //link fist and last to complete shape
	}
}
