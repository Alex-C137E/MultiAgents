package game

type Score struct {
	Level            int
	Value            int
	RequiredFishType int
}

func NewScore(level int, requiredFishType int) *Score {
	return &Score{level, 0, requiredFishType}
}

func (s *Score) AddCollectedFish(fishType int) {
	if fishType == s.RequiredFishType {
		s.Value = s.Value + 100
	} else {
		s.Value = s.Value - 200
	}
}
