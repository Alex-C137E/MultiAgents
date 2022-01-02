package game

//Score struct permet de définir le score actuel de l'utilisateur
type Score struct {
	Level            int
	Value            int
	RequiredFishType int
}

//NewScore fonction constructeur de score
func NewScore(level int, currentScore int, requiredFishType int) *Score {
	return &Score{level, currentScore, requiredFishType}
}

//AddCollectedFish fonction qui permet d'ajouter un poisson collecté au score actuel
func (s *Score) AddCollectedFish(fishType int) {
	if fishType == s.RequiredFishType {
		s.Value = s.Value + 100
	} else {
		s.Value = s.Value - 200
	}
}
