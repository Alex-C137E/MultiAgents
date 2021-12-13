package agent

//Agent interface pour tous les agents de la simulation
type Agent interface {
	Start()  //Fonction permettant d'initialiser l'agent
	Update() //Fonction permettant de modifier le comportement de l'agent (position)
	Draw()   //Fonction permettant d'afficher l'agent
}
