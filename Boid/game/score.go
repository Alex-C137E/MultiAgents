package game

type Score struct {
	Level            int
	Value            int
	RequiredFishType int
}

func NewScore(level int, currentScore int, requiredFishType int) *Score {
	return &Score{level, currentScore, requiredFishType}
}

func (s *Score) AddCollectedFish(fishType int) {
	if fishType == s.RequiredFishType {
		s.Value = s.Value + 100
	} else {
		s.Value = s.Value - 200
	}
}
