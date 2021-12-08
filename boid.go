package main

import(
	"fmt"
	"math"
	"math/rand"
	"time"
	gg "github.com/fogleman/gg"
)

const ALPHA = 180

// introduire différents type de population

func main(){
		// regrouper monde comme tableau de Flock
		var monde []*Boid
		//var monde2 []*Boid
		// ----- init monde.boids ------ //
		for i:=0 ; i<10; i++{
			x := -1.0 + rand.Float64()*2
			y := -1.0 + rand.Float64()*2
			boid := NewBoid(i,[2]float64{500+20*x,500+20*y}, 1) // 20 c'est le rayon du flock ici ! 
			monde = append(monde, &boid)
		}
		// --- life ---- //
		ticker := time.NewTicker(10 * time.Millisecond)
	  done := make(chan bool)
		c := 0
		go func() {
	        for {
	            select {
	            case <-done:
	                return
	            case t := <-ticker.C:
								fmt.Println("Tick at", t)
								dc := gg.NewContext(1000, 1000)
								// ici on actualise tout
								for _,boid := range(monde){
									boid.GetCouches(monde)
                  boid.flock()
								}
								for _, boid := range(monde){
                  boid.Update()
									boid.Border(1000.0,1000.0)
                  boid.Draw(dc)
								}
								img := fmt.Sprintf("out_n%v.png",c)
								dc.Fill()
								dc.SavePNG(img)
								c+=1
								// travailler localement avec les boid -> créer des sous-monde.boidss.
	        }
				}
	    }()
	  time.Sleep(40000 * time.Millisecond)
	  ticker.Stop()
		done <- true
		fmt.Println("Ticker stopped")
}

type Boid struct {
	ID int
	position [2]float64
  vitesse [2]float64
  acceleration [2]float64
	longueur float64
	rayonSeparation float64
	rayonCohesion float64
	rayonVitesse float64
  max_force float64
  max_speed [2]float64
	separation map[*Boid]bool
	Cohesion map[*Boid]bool
	alignement map[*Boid]bool
	couleur int
}

func NewBoid(id int, position [2]float64, couleur int) Boid{
  return Boid{
		id,
    position,
    [2]float64{1.5,1.0},
    [2]float64{0.0,0.0},
    5,
		20,
		30,
		33,
		10,
    [2]float64{10.0,10.0},
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		couleur,
  }
}

// ne fonctionne pas. -> repenser l'espace (tore ?)
func (b *Boid) Border(width float64, height float64){
  if (b.position[0] < -b.longueur){
    b.position[0] = width - b.longueur
  }
  if (b.position[0] > width + b.longueur ){
    b.position[0] = b.longueur
  }
  if (b.position[1] < -b.longueur){
    b.position[1] = height - b.longueur
  }
  if (b.position[1] > height + b.longueur){
    b.position[1] = b.longueur
  }
}

func (b *Boid) limit(){
	// regarder la norme plutot
	norme := math.Sqrt(math.Pow(b.vitesse[0],2) + math.Pow(b.vitesse[1],2))
	norme_max := math.Sqrt(math.Pow(b.max_speed[0],2) + math.Pow(b.max_speed[1],2))
  if norme >= norme_max{
		b.vitesse[0] = 1.0 + 2*rand.Float64()
		b.vitesse[1] = 1.0 + 2*rand.Float64()
	}
}

func (b *Boid) Update(){
	b.vitesse[0] += b.acceleration[0]
  b.vitesse[1] += b.acceleration[1]
  b.limit()
  b.position[0] += b.vitesse[0]
  b.position[1] += b.vitesse[1]
  b.acceleration = [2]float64{0.0,0.0}
}

func (b *Boid) ApplyForce(force [2]float64){
  b.acceleration[0] += force[0]
  b.acceleration[1] += force[1]
}

func (b *Boid) flock(){
  sep := b.separate()
  ali := b.align()
  coh := b.cohesion()
  b.ApplyForce(sep)
  b.ApplyForce(ali)
  b.ApplyForce(coh)
}

func (b *Boid) separate() (t [2]float64){
  steer := [2]float64{0.0,0.0}
  for boid,_ := range(b.separation){
    steer[0] = steer[0] - boid.position[0]
    steer[1] = steer[1] - boid.position[1]
  }
  if len(b.separation) != 0 {
    steer[0] = steer[0]/float64(len(b.separation))
    steer[1] = steer[1]/float64(len(b.separation))
  }
	steer[0] -= b.vitesse[0]
	steer[1] -= b.vitesse[1]
  return steer
}

func (b *Boid) align() (t [2]float64){
  alignement := [2]float64{0.0, 0.0}
	for boid,_ := range(b.alignement){
		alignement[0] = alignement[0] + boid.vitesse[0]
		alignement[1] = alignement[1] + boid.vitesse[1]
	}
	if len(b.alignement) != 0{
  alignement[0] = alignement[0]/float64(len(b.alignement))
	alignement[1] = alignement[1]/float64(len(b.alignement))
}
  return alignement
}

func (b *Boid) cohesion() (t [2]float64){
  cohesion := [2]float64{0.0, 0.0}
	for boid,_ := range(b.Cohesion){
		cohesion[0] = cohesion[0] + boid.position[0]
		cohesion[1] = cohesion[1] + boid.position[1]
	}
	if len(b.Cohesion) != 0{
		cohesion[0] = cohesion[0]/float64(len(b.Cohesion))
		cohesion[1] = cohesion[1]/float64(len(b.Cohesion))
	}
  return cohesion
}


func (b *Boid) GetCouches(boids []*Boid){
	// Couche separation, cohesion, alignement
	for _,boid := range(boids){
		if boid.ID == b.ID{
			continue
		}
		distance := math.Sqrt(math.Pow(b.position[0] - boid.position[0],2) + math.Pow(b.position[1] - boid.position[1],2))
		if distance <= b.rayonSeparation{
				b.separation[boid] = true
				b.alignement[boid] = true
				b.Cohesion[boid] = true
				//separation = append(separation, boid) // map plutot que tableau
			}else if (distance > b.rayonSeparation && distance <= b.rayonVitesse){
				b.alignement[boid] = true
				b.Cohesion[boid] = true
				//alignement = append(alignement, boid)
			}else if (distance > b.rayonVitesse && distance <= b.rayonCohesion){
				b.Cohesion[boid] = true
				//cohesion = append(cohesion, boid)
			}
	}
}

func (b *Boid) Draw(dc *gg.Context ){
	switch b.couleur{
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
	dc.DrawCircle(b.position[0],b.position[1], b.longueur)
	dc.Fill()
	}
