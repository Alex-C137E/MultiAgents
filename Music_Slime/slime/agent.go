package slime

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	types "gitlab.utc.fr/projet_ia04/musicslime/types"
)

type Agent struct {
	PreviousPosition types.Position
	Position         types.Position    //Agent current position
	Angle            float64           //Agent current Angle
	SpeciesMask      types.SpeciesMask //Agent Mask
	CurrentSpecies   int               //Agent current species
	id               int               //Agent identification number
	// sync             chan string
	// musicInfo        string
	change bool
}

func NewAgent(position types.Position, angle float64, speciesMask types.SpeciesMask, currentSpecies int, id int /*c chan string*/) *Agent {
	return &Agent{position, position, angle, speciesMask, currentSpecies, id /*c, ""*/, false}
}

func (agent *Agent) sense(settings SpeciesSettings, sensorAngleOffset float64, trailMap *types.TrailMap) float64 {
	sensorAngle := agent.Angle + sensorAngleOffset

	var sensorDir [2]float64
	sensorDir[0] = math.Cos(sensorAngle)
	sensorDir[1] = math.Sin(sensorAngle)

	var sensorPos types.Position
	sensorPos.X = agent.Position.X + sensorDir[0]*settings.SensorOffsetDistance
	sensorPos.Y = agent.Position.Y + sensorDir[1]*settings.SensorOffsetDistance

	sensorCentreX := int(sensorPos.X)
	sensorCentreY := int(sensorPos.Y)

	var sum float64

	var senseWeight types.SpeciesMask
	// Permet de faire en sorte que les phéromone laissé par les autres espèces repousse une espèce différente
	senseWeight[0] = agent.SpeciesMask[0]*2 - 1
	senseWeight[1] = agent.SpeciesMask[1]*2 - 1
	senseWeight[2] = agent.SpeciesMask[2]*2 - 1
	senseWeight[3] = agent.SpeciesMask[3]*2 - 1

	for offsetX := -settings.SensorSize; offsetX <= settings.SensorSize; offsetX++ {
		for offsetY := -settings.SensorSize; offsetY <= settings.SensorSize; offsetY++ {
			width := types.GetWindowDefault().Width
			height := types.GetWindowDefault().Height

			sampleX := math.Min(width-float64(1), math.Max(float64(0), float64(sensorCentreX+offsetX)))
			sampleY := math.Min(height-float64(1), math.Max(float64(0), float64(sensorCentreY+offsetY)))
			if sampleX > width {
				sampleX = width
			}
			if sampleY > height {
				sampleY = height
			}
			trail := (*trailMap)[int(sampleX)][int(sampleY)]

			sum = sum + (float64(senseWeight[0]) * trail[0]) + (float64(senseWeight[1]) * trail[1]) + (float64(senseWeight[2]) * trail[2]) + (float64(senseWeight[3]) * trail[3])
		}
	}

	return sum
}

func (agent *Agent) Update(trailMap *types.TrailMap, musicInfo string) {
	settings := GetSpeciesSettingsByID(agent.CurrentSpecies, musicInfo)
	pos := agent.Position
	width := types.GetWindowDefault().Width
	height := types.GetWindowDefault().Height
	random := Hash(uint(pos.Y)*uint(width) + uint(pos.X) + Hash(uint(time.Now().Unix())*100000))

	// Steer based on sensory data
	sensorAngleRad := settings.SensorAngleDegrees * (3.1215 / 100)
	weightForward := agent.sense(*settings, 0, trailMap)
	weightLeft := agent.sense(*settings, sensorAngleRad, trailMap)
	weightRight := agent.sense(*settings, -sensorAngleRad, trailMap)

	randomSteerStrength := scaleToRange01(random)
	turnSpeed := settings.TurnSpeed * 2 * 3.1415

	// Enfonction de que "sentent" les agent devant, à droite et à gauche:
	// Continue in same direction
	if weightForward > weightLeft && weightForward > weightRight {
		agent.Angle += 0
	} else if weightForward < weightLeft && weightForward < weightRight {
		agent.Angle += (randomSteerStrength - 0.5) * 2 * turnSpeed * types.DeltaTime
	} else if weightRight > weightLeft { //Turn right
		agent.Angle -= randomSteerStrength * turnSpeed * types.DeltaTime
	} else if weightLeft > weightRight { //Turn left
		agent.Angle += randomSteerStrength * turnSpeed * types.DeltaTime
	}

	// Update position
	var direction types.Position
	direction.X = math.Cos(agent.Angle)
	direction.Y = math.Sin(agent.Angle)

	var newPosition types.Position
	newPosition.X = agent.Position.X + direction.X*types.DeltaTime*settings.MoveSpeed
	newPosition.Y = agent.Position.Y + direction.Y*types.DeltaTime*settings.MoveSpeed

	// Permet de faire en sorte que quand les agents rentrent dans le mur/bord de la map il chancge de direction
	// Clamp position to map boundaries, and pick new random move dir if hit boundary
	if newPosition.X < 0 || newPosition.X >= width || newPosition.Y < 0 || newPosition.Y >= height {
		random = Hash(random)
		randomAngle := scaleToRange01(random) * 2 * 3.1415

		newPosition.X = math.Min(width-float64(1), math.Max(float64(0), newPosition.X))
		newPosition.Y = math.Min(height-float64(1), math.Max(float64(0), newPosition.Y))

		agent.Angle = randomAngle
	} else {
		coord := newPosition
		oldTrail := (*trailMap)[int(coord.X)][int(coord.Y)]
		(*trailMap)[int(coord.X)][int(coord.Y)][0] = math.Min(1.0, oldTrail[0]+float64(agent.SpeciesMask[0])*float64(types.TrailWeight)*types.DeltaTime)
		(*trailMap)[int(coord.X)][int(coord.Y)][1] = math.Min(1.0, oldTrail[1]+float64(agent.SpeciesMask[1])*float64(types.TrailWeight)*types.DeltaTime)
		(*trailMap)[int(coord.X)][int(coord.Y)][2] = math.Min(1.0, oldTrail[2]+float64(agent.SpeciesMask[2])*float64(types.TrailWeight)*types.DeltaTime)
		(*trailMap)[int(coord.X)][int(coord.Y)][3] = math.Min(1.0, oldTrail[3]+float64(agent.SpeciesMask[3])*float64(types.TrailWeight)*types.DeltaTime)
	}

	agent.PreviousPosition = agent.Position
	agent.Position = newPosition

	if !agent.change && musicInfo == "very hard drop" {
		// Faire apparaitre les agents partout
		w := rand.Intn(int(width))
		h := rand.Intn(int(height))
		agent.Position.X = float64(w)
		agent.Position.Y = float64(h)
		agent.change = true
	} else if agent.change {
		agent.Position.X = width / 2
		agent.Position.Y = height / 2
		agent.change = false
	}
}

func (agent *Agent) Start(id int, c chan string) {
	// agent.sync = c
	agent.id = id

	width := types.GetWindowDefault().Width
	height := types.GetWindowDefault().Height

	//Setting default position
	agent.Position.X = width / 2
	agent.Position.Y = height / 2

	agent.PreviousPosition = agent.Position

	//Setting default angle
	agent.Angle = rand.Float64() * math.Pi * 2

	//Species of the agent
	agent.CurrentSpecies = rand.Intn(types.SpeciesNum)
	//fmt.Println("Choice", agent.CurrentSpecies, types.SpeciesNum)
	//agent.CurrentSpecies = 1
	switch agent.CurrentSpecies {
	case 0:
		agent.SpeciesMask = types.SpeciesMask{1, 0, 0, 0}
	case 1:
		agent.SpeciesMask = types.SpeciesMask{0, 1, 0, 0}
	case 2:
		agent.SpeciesMask = types.SpeciesMask{0, 0, 1, 0}
	case 3:
		agent.SpeciesMask = types.SpeciesMask{0, 0, 0, 1}
	}

	// Pas très fonctionel comme attendu
	// go func() {
	// 	for {
	// 		// lorsque l'agent  reçoit sur sa channel sync(bloquant): il reçoit une indication de la musique
	// 		agent.musicInfo = <-agent.sync
	// 		// Il doit modifier un de ses paramêtre
	// 	}
	// }()
}

func (agent *Agent) Draw(screen *ebiten.Image, musicInfo string) {
	settings := GetSpeciesSettingsByID(agent.CurrentSpecies, musicInfo)

	var color color.RGBA
	color.R = uint8(settings.Colour[0] * 255)
	color.G = uint8(settings.Colour[1] * 255)
	color.B = uint8(settings.Colour[2] * 255)
	color.A = uint8(settings.Colour[3] * 255)
	//ebitenutil.DrawLine(screen, agent.PreviousPosition.X, agent.PreviousPosition.Y, agent.Position.X, agent.Position.Y, color)
	if settings.Displayed {
		ebitenutil.DrawRect(screen, agent.Position.X, agent.Position.Y, 2, 2, color)
	}
}
