package game

import (
	"image/color"
	"sync"

	ebiten "github.com/hajimehoshi/ebiten/v2"

	agent "gitlab.utc.fr/jedescam/projet_ia04/final/agent"
	settings "gitlab.utc.fr/jedescam/projet_ia04/final/settings"
	worldelements "gitlab.utc.fr/jedescam/projet_ia04/final/worldelements"
)

//Game type regroupe tous les états du monde
type Game struct {
	birds    []*agent.Bird
	settings settings.Settings
	window   settings.Window
	worldMap worldelements.Map
}

//NewGame fonction constructeur de Game, regroupe tous les états du monde
func NewGame() *Game {
	g := &Game{}

	//Initialize settings
	g.settings = settings.GetDefaultSettings()
	g.window = settings.GetWindowDefault()

	//Create birds
	g.birds = make([]*agent.Bird, g.settings.AgentsNum)
	for i := 0; i < len(g.birds); i++ {
		g.birds[i] = agent.NewBird(g.settings, i)
	}

	//Initialize birds
	var wg sync.WaitGroup
	wg.Add(g.settings.AgentsNum)
	for i := 0; i < g.settings.AgentsNum; i++ {
		go func(i int) {
			g.birds[i].Start(g.window)
			wg.Done()
		}(i)
	}
	wg.Wait()

	return g
}

//Update fonction appelée lorsque l'état du monde est mis à jour
func (g *Game) Update() error {
	//Update birds
	var wg sync.WaitGroup
	wg.Add(g.settings.AgentsNum)

	for i := 0; i < g.settings.AgentsNum; i++ {
		go func(i int) {
			g.birds[i].Update(g.window, g.birds)
			wg.Done()
		}(i)
	}

	wg.Wait()

	return nil
}

//Draw fonction permet d'afficher le monde en déléguant l'affichage aux éléments inclus dans le monde
func (g *Game) Draw(screen *ebiten.Image) {
	//Draw background
	screen.Fill(color.White)

	//Draw agents
	for i := 0; i < g.settings.AgentsNum; i++ {
		g.birds[i].Draw(screen, g.settings)
	}

	//Draw world elements
}

//Layout fonction permet de définir la taille de la fenêtre d'affichage
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.window.Width), int(g.window.Height)
}
