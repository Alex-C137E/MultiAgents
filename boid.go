package main

import(
	"fmt"
	"math"
	"math/rand"
	"time"
	gg "github.com/fogleman/gg"
	//"gg"
)

const ALPHA = 180

func main(){
		var monde []*Boid
		// ----- init monde.boids ------ //
		for i:=0 ; i<10; i++{
			boid := NewBoid(i,[2]float64{100 + 20*rand.Float64(),200 + 35*rand.Float64()})
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
									boid.Update(monde)
								}
								for _, boid:= range(monde){
									//boid.CheckPosition(1000.0, 1000.0)
									boid.Draw(dc)
									boid.CheckSpeed()
									boid.Update_pos(1000.0, 1000.0)
									//fmt.Println(boid.vitesse) // créer une structure contenant tous les boids ? -> utile pour les tests et les stats.
								}
								img := fmt.Sprintf("out_n%v.png",c)
								dc.Fill()
								dc.SavePNG(img)
								c+=1
								// travailler localement avec les boid -> créer des sous-monde.boidss.
	        }
				}
	    }()
		//
	  time.Sleep(10000 * time.Millisecond)
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
	cohesion map[*Boid]bool
	alignement map[*Boid]bool
	couleur int
}

func NewBoid(id int, position [2]float64) Boid{
  return Boid{
		id,
    position,
    [2]float64{1.0,1.0},
    [2]float64{2.0,3.0},
    20,
		5,
		10,
		20,
		10,
    [2]float64{15.0,15.0},
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		make(map[*Boid]bool),
		1,
  }
}

/*
func (b *Boid) CheckPosition(width float64, height float64){
	fmt.Println(b.position)
	if b.position[0] < 0.0{
		b.position[0] = width
	}
	if b.position[0] > width{
		b.position[0] = 0.0
	}
	if b.position[1] < 0.0{
		b.position[0] = height
	}
	if b.position[1] > height{
		b.position[1] = 0.0
	}
	fmt.Println(b.position)
}
*/
func (b *Boid) CheckSpeed(){
	// tod:  plutot checker la norme du vecteur
	if b.vitesse[0] > b.max_speed[0]{
		b.vitesse[0] = 5
	}
	if b.vitesse[1] > b.max_speed[1]{
		b.vitesse[1] = 5
	}
}
func (b * Boid) UpdateSpeed(boids map[*Boid]bool){
	// couche alignement
	var avg_speed [2]float64
	for boid,_ := range(boids){
		avg_speed[0] = avg_speed[0] + boid.vitesse[0]
		avg_speed[1] = avg_speed[1] + boid.vitesse[1]
	}
	if len(boids) != 0{
	b.vitesse[0] += avg_speed[0]/float64(len(boids))
	b.vitesse[1] += avg_speed[1]/float64(len(boids))
}
}

func (b *Boid) UpdateSeparation(boids map[*Boid]bool){
	delta_x := 0.0
	delta_y := 0.0
	for boid,_ := range(boids){
		delta_x = delta_x - boid.position[0]
		delta_y = delta_y - boid.position[1]
	}
	if len(boids) != 0 {
	b.vitesse[0] += delta_x/float64(len(boids))
	b.vitesse[1] += delta_y/float64(len(boids))
}
}

func (b *Boid) UpdateCohesion(boids map[*Boid]bool){
	delta_x := 0.0
	delta_y := 0.0
	for boid,_ := range(boids){
		delta_x = delta_x + boid.position[0]
		delta_y = delta_y + boid.position[1]
	}
	if len(boids) != 0{
		b.vitesse[0] += delta_x/float64(len(boids))
		b.vitesse[1] += delta_y/float64(len(boids))
	}
}

func (b *Boid) Update(boids []*Boid){
	b.GetCouches(boids)
	b.UpdateSeparation(b.separation)
	b.UpdateCohesion(b.cohesion)
	b.UpdateSpeed(b.alignement)
}

func (b *Boid) Update_pos(width float64, height float64){
	nouvelle_pos_x := b.position[0] + b.vitesse[0]
	nouvelle_pos_y := b.position[1] + b.vitesse[1]
	update := false
	if nouvelle_pos_x < 0{
		b.position[0] = width + (b.vitesse[0] + b.position[0]) //  b.vitesse[0] negatif mais plus grand en va que b.position(0)
		update = true
	}
	if nouvelle_pos_x > width{
		b.position[0] =  (b.vitesse[0] + (width - b.position[0])) //  b.vitesse[0] negatif mais plus grand en va que b.position(0)
		update =true
	}
	if nouvelle_pos_y < 0{
		b.position[0] = height + (b.vitesse[0] + b.position[0]) //  b.vitesse[0] negatif mais plus grand en va que b.position(0)
		update =true
	}
	if nouvelle_pos_y > height{
		b.position[0] =  (b.vitesse[0] + (height - b.position[0])) //  b.vitesse[0] negatif mais plus grand en va que b.position(0)
		update =true
	}
	if !(update){
		b.position[0] += b.vitesse[0] // pas de temps de 1 (d = vt)
		b.position[1] += b.vitesse[1]
	}
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
				b.cohesion[boid] = true
				//separation = append(separation, boid) // map plutot que tableau
			}else if (distance > b.rayonSeparation && distance <= b.rayonVitesse){
				b.alignement[boid] = true
				b.cohesion[boid] = true
				//alignement = append(alignement, boid)
			}else if (distance > b.rayonVitesse && distance <= b.rayonCohesion){
				b.cohesion[boid] = true
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
