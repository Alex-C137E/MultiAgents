package game

//Level struct permet de définir un niveau
type Level struct {
	RepulsionFactorBtwnSpecies float64
	SeparationPerception       float64
	CohesionPerception         float64
	AlignPerception            float64
	numWall                    int
	MaxForce                   float64
	MaxSpeed                   float64
	polygonSize                float64
	SharkDensity               int
}

//NewLevel fonction constructeur d'un niveau
func NewLevel(RepulsionFactorBtwnSpecies float64,
	SeparationPerception float64,
	CohesionPerception float64,
	AlignPerception float64,
	numWall int,
	MaxForce float64,
	MaxSpeed float64,
	polgonSize float64,
	SharkDensity int) *Level {
	return &Level{RepulsionFactorBtwnSpecies,
		SeparationPerception,
		CohesionPerception,
		AlignPerception,
		numWall,
		MaxForce,
		MaxSpeed,
		polgonSize,
		SharkDensity}
}

//LoadLevels fonction qui permet de charger des niveaux préféfinis par défaut
func LoadLevels() []*Level {
	//Création de 5 niveaux:
	levels := make([]*Level, 5)

	// niveau 0: ce niveau est très simple:les poissons sont amenés à se regrouper et se stabiliser rapidement
	// rendant leur pêche facile.
	levels[0] = NewLevel(10000, 10, 300, 100, 300, 1, 4.0, 1000, 15)

	// niveau 1:  ce niveau est identique au niveau 0 à la différence que cette fois-ci les requins attaquent bien
	// plus: leur attribue density est diminué de tel sorte à ce qu’il attaque pour une quantité de poisson dans leur
	// champ d’attaque inférieur à celle du niveau 0. Le joueur doit donc être plus rapide pour ne pas se faire
	// manger tous ses poisons par les prédateurs.
	levels[1] = NewLevel(10000, 10, 300, 100, 16+10+300, 1, 4.0, 1000, 8)

	// niveau 2:  ce niveau est dans la continuité du niveau 1 mais ici: le facteur de cohésion diminue, et ceux de
	// répulsion intra et inter espèces augmentent ce qui diminue la stabilité dans le comportement des poissons les
	// rendant plus complexes à attraper: ils se regroupent moins, et se mélangent plus entre espèces. De plus, la
	// taille maximale du filet diminue.
	levels[2] = NewLevel(500, 100, 100, 75, 16+10+300, 2.0, 4.0, 700, 8)

	// niveau 3:  pour ce niveau, on reprend le niveau 2 et on augmente le niveau difficulté en réduisant le niveau
	// de stabilité de manière similaire à ce qui fut fait pour le niveau 2. En plus, on rend les prédateurs plus
	// agressifs en utilisant le même procédé utilisé dans le niveau 1.
	levels[3] = NewLevel(50, 100, 75, 75, 16+10+300, 2.0, 4.0, 700, 5)

	// niveau 4: pour ce niveau on reprend le niveau 3 et on  rajoute des murs/ bombes pour favoriser le chaos et
	// rendre plus difficile la tâche d’'attraper les poissons. De plus les poissons sont moins en cohésion et vont
	// plus vite et bien sûr, les prédateurs sont encore plus agressifs :).
	levels[4] = NewLevel(50, 100, 50, 75, 16+10+2*48, 2.0, 5.0, 700, 3)

	return levels
}
