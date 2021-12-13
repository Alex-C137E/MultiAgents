package agents

import(
  gg "github.com/fogleman/gg"
  "math"
  "math/rand"
)

type Boid struct {
	ID int
  Type string // Proie ou Predateur, Raver ou Flic
	Position [2]float64
  Vitesse [2]float64
  Acceleration [2]float64
	Longueur float64
	RayonSeparation float64
	RayonCohesion float64
	RayonVitesse float64
  Max_force float64
  Max_speed [2]float64
	Separation map[*Boid]bool
	Cohesion map[*Boid]bool
	Alignement map[*Boid]bool
	Couleur int
	Energie float64 // entre 0 et 1
  Recharge bool
  Source_act *Source
}

type Monde struct{
  Agents [][][]*Boid
  Preds [][][]*Pred
  Sources []*Source
}

func New_boid(id int, position [2]float64, couleur int) Boid{
  return Boid{
		id,
    "proie",
    position,
    [2]float64{1.5,1.0},
    [2]float64{1.0,1.0},
    8,
		20,
		30,
		33,
		10,
    [2]float64{20.0,20.0},
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		couleur,
		rand.Float64(),
    false,
    nil,
  }
}


func (b *Boid) Is_alive() (live bool){
	if b.Energie < 0{
		return false
	}else{
		return true
	}
}

func (b *Boid) Border(width float64, height float64){
  if (b.Position[0] < -b.Longueur){
    b.Position[0] = width - b.Longueur
  }
  if (b.Position[0] > width + b.Longueur ){
    b.Position[0] = b.Longueur
  }
  if (b.Position[1] < -b.Longueur){
    b.Position[1] = height - b.Longueur
  }
  if (b.Position[1] > height + b.Longueur){
    b.Position[1] = b.Longueur
  }
}

func (b *Boid) Speed_limit(){
	// regarder la norme plutot
	norme := math.Sqrt(math.Pow(b.Vitesse[0],2) + math.Pow(b.Vitesse[1],2))
	norme_max := math.Sqrt(math.Pow(b.Max_speed[0],2) + math.Pow(b.Max_speed[1],2))
  if norme >= norme_max{
		b.Vitesse[0] = 1.0 + 5*rand.Float64()
		b.Vitesse[1] = 1.0 + 5*rand.Float64()
	}
}

func (b *Boid) Update(){
	b.Vitesse[0] += b.Acceleration[0]
  b.Vitesse[1] += b.Acceleration[1]
  b.Speed_limit()
  b.Position[0] += b.Vitesse[0]
  b.Position[1] += b.Vitesse[1]
	b.Energie -= 0.01
  b.Acceleration = [2]float64{0.0,0.0}
}

func (b *Boid) Apply_force(force [2]float64){
  b.Acceleration[0] += force[0]
  b.Acceleration[1] += force[1]
}

func (b *Boid) Flock(){
  /*
  A ce stade on a deja recuperer les voisins dans les différentes couches.
  Vérifier que les couches ne sont pas nulles ?
  */
  // Verifier si le Boid se trouve à proximite d'une source ?
  if b.Recharge && b.Energie < 0.8{ // actualiser le b.recharge
    b.Energie += b.Source_act.Energie_donnee
    b.Couleur = 3
    b.Vitesse[0] = 0
    b.Vitesse[1] = 0
    return // le boid ne peut pas bouger
  }
  b.Couleur = 1
  sep := b.Separate_force()
  ali := b.Align_force()
  coh := b.Cohesion_force()
  b.Apply_force(sep)
  b.Apply_force(ali)
  b.Apply_force(coh)
}

// Permet d'obtenir les autres agents se trouvant à différents rayons.
// TODO : regarder pour la couche de separation si les agents qui s'y trouvent sont tous du meme type -> simuler Bernoulli sinon, plus declencher signal
// de détresse si survit.
func (b *Boid) Get_couches(monde *Monde){
  // source
  for _,source := range(monde.Sources){
    distance := math.Sqrt(math.Pow(b.Position[0] - source.Position[0],2) + math.Pow(b.Position[1] - source.Position[1],2))
    if distance <= source.Rayon{
      b.Recharge = true
      b.Source_act = source
      break
    }
    b.Recharge = false
    b.Source_act = nil
  }
  // TODO : prendre en compte les bords dans le comptage des couches.
	//for _,boid := range(monde.Agents){ // etendre boids à tous les agents y compris les sources ?
  for i,_ := range(monde.Agents){
    for j,_ := range(monde.Agents[i]){
      boids := monde.Agents[i][j]
      if len(boids) == 0{
        continue
      }
      for _,boid := range(boids){
    		if boid.ID == b.ID{
    			continue
    		}
    		distance := math.Sqrt(math.Pow(b.Position[0] - boid.Position[0],2) + math.Pow(b.Position[1] - boid.Position[1],2))
    		if distance <= b.RayonSeparation{
    				b.Separation[boid] = true
    				b.Alignement[boid] = true
    				b.Cohesion[boid] = true
    				//separation = append(separation, boid) // map plutot que tableau
    			}else if (distance > b.RayonSeparation && distance <= b.RayonVitesse){
    				b.Alignement[boid] = true
    				b.Cohesion[boid] = true
    				//alignement = append(alignement, boid)
    			}else if (distance > b.RayonVitesse && distance <= b.RayonCohesion){
    				b.Cohesion[boid] = true
    				//cohesion = append(cohesion, boid)
    			}
        }
	}
}
}

// -------- forces qui s'appliquent à chaque agent (boid -> separation, cohesion, alignement) ------------------

func (b *Boid) Separate_force() (t [2]float64){
  steer := [2]float64{0.0,0.0}
  for boid,_ := range(b.Separation){
    steer[0] = steer[0] - boid.Position[0]
    steer[1] = steer[1] - boid.Position[1]
  }
  if len(b.Separation) != 0 {
    steer[0] = steer[0]/float64(len(b.Separation))
    steer[1] = steer[1]/float64(len(b.Separation))
  }
	steer[0] -= b.Vitesse[0]
	steer[1] -= b.Vitesse[1]
  return steer
}

func (b *Boid) Align_force() (t [2]float64){
  alignement := [2]float64{0.0, 0.0}
	for boid,_ := range(b.Alignement){
		alignement[0] = alignement[0] + boid.Vitesse[0]
		alignement[1] = alignement[1] + boid.Vitesse[1]
	}
	if len(b.Alignement) != 0{
  alignement[0] = alignement[0]/float64(len(b.Alignement))
	alignement[1] = alignement[1]/float64(len(b.Alignement))
}
  return alignement
}

func (b *Boid) Cohesion_force() (t [2]float64){
  cohesion := [2]float64{0.0, 0.0}
	for boid,_ := range(b.Cohesion){
		cohesion[0] = cohesion[0] + boid.Position[0]
		cohesion[1] = cohesion[1] + boid.Position[1]
	}
	if len(b.Cohesion) != 0{
		cohesion[0] = 1.2*cohesion[0]/float64(len(b.Cohesion))
		cohesion[1] = 1.2*cohesion[1]/float64(len(b.Cohesion))
	}
  return cohesion
}

// draw

func (b *Boid) Draw(dc *gg.Context ){
	switch b.Couleur{
	case 1:
		dc.SetRGB(200, 400, 0)
	case 2:
		dc.SetRGB(100, 400, 50)
	case 3:
		dc.SetRGB(100, 100, 300)
	case 4:
		dc.SetRGB(100, 400, 50)
	case 0:
	   dc.SetRGB(500,500,300)
	}
	dc.DrawCircle(b.Position[0],b.Position[1], b.Longueur)
	dc.Fill()
	}
