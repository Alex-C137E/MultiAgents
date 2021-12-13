package slime

import (
	types "gitlab.utc.fr/projet_ia04/musicslime/types"
)

type SpeciesSettings struct {
	MoveSpeed            float64 // augmentation => accélération
	TurnSpeed            float64 // augmentation => accélération
	SensorAngleDegrees   float64 // augmentation => instabilité via création de multiple chemin
	SensorOffsetDistance float64 // augmentation => périmmetre du senseur augmente => stabilité en gros groupe unique
	SensorSize           int     // augmentation => Stabilité
	Colour               types.Colour
	Displayed            bool
}

func NewSpeciesSettings(moveSpeed float64, turnSpeed float64, sensorAngleDegrees float64, sensorOffsetDistance float64, sensorSize int, colour types.Colour, displayed bool) *SpeciesSettings {
	return &SpeciesSettings{moveSpeed, turnSpeed, sensorAngleDegrees, sensorOffsetDistance, sensorSize, colour, displayed}
}

func GetSpeciesSettingsByID(id int, musicInfo string) *SpeciesSettings {
	var initColour types.Colour
	initColour[0] = 0
	initColour[1] = 0
	initColour[2] = 0
	initColour[3] = 1
	speciesSettings := NewSpeciesSettings(1, 1, 1, 1, 1, initColour, true)
	var colour types.Colour
	switch id {
	case 0:
		colour[0] = 0.53608304
		colour[1] = 0.062745094
		colour[2] = 1
		colour[3] = 0
		speciesSettings.Colour = colour
		speciesSettings.Displayed = false
		if musicInfo == "hard drop" || musicInfo == "very hard drop" {
			speciesSettings.Displayed = true
		}
	case 1:
		colour[0] = 0.10588236
		colour[1] = 0.6240185
		colour[2] = 0.83137256
		colour[3] = 1
		speciesSettings.Colour = colour
	case 2:
		colour[0] = 0.47160125
		colour[1] = 0.8301887
		colour[2] = 0.105731584
		colour[3] = 1
		speciesSettings.Colour = colour
	}
	switch musicInfo {
	case "small drop":
		speciesSettings.MoveSpeed = 2
		speciesSettings.TurnSpeed = 180
		speciesSettings.SensorAngleDegrees = 10
		speciesSettings.SensorOffsetDistance = 20
		speciesSettings.SensorSize = 1
	case "medium drop":
		speciesSettings.MoveSpeed = 10
		speciesSettings.TurnSpeed = 180
		speciesSettings.SensorAngleDegrees = 90
		speciesSettings.SensorOffsetDistance = 50
		speciesSettings.SensorSize = 1
	case "hard drop":
		speciesSettings.MoveSpeed = 20
		speciesSettings.TurnSpeed = 180
		speciesSettings.SensorAngleDegrees = 180
		speciesSettings.SensorOffsetDistance = 50
		speciesSettings.SensorSize = 1
	case "very hard drop":
		speciesSettings.MoveSpeed = 50
		speciesSettings.TurnSpeed = 180
		speciesSettings.SensorAngleDegrees = 180
		speciesSettings.SensorOffsetDistance = 50
		speciesSettings.SensorSize = 1
	}
	// return NewSpeciesSettings(10, 180, 10, 20, 1, colour) //explosion patern
	return speciesSettings
}
