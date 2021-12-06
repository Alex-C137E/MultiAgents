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

type Particule struct {
	position [2]float64
	PHI float64
	BETA float64
	vitesse float64
	rayon float64
	Nb_voisins_gauche int
	Nb_voisins_droite int
	couleur int
}
//constructeur de Particule ?

func main(){
	// init
	rand.Seed(time.Now().UnixNano())

	var monde []*Particule
	for i := 0 ; i<2000; i++{
		x := -1.0 + rand.Float64()*2
		y := -1.0 + rand.Float64()*2
		new_part := NouvelleParticule([2]float64{500 + x*499,500 + y*499}, rand.Float64())
		monde = append(monde, &new_part)
	}
	// actualisation - timer
	c := 0
	ticker := time.NewTicker(5 * time.Millisecond)
  done := make(chan bool)
	go func() {
        for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
							fmt.Println("Tick at", t)
							dc := gg.NewContext(1000, 1000)
							for _, part := range(monde){
								part.Actualisation(1.0,monde, dc)
							}
							img := fmt.Sprintf("tout_n%v.png",c)
							dc.Fill()
							dc.SavePNG(img)
							c += 1
        }
			}
    }()

	time.Sleep(990000 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
	}

func NouvelleParticule(position [2]float64, Phi float64) Particule {
	return Particule{
		position,
		Phi,
		Phi,
		20.0,
		40.0,
		0,
		0,
		0,
	}
}

func (p Particule) Draw(dc *gg.Context ){
	switch p.couleur{
	case 1:
		dc.SetRGB255(255, 255, 0)
	case 2:
		dc.SetRGB255(0, 0, 255)
	case 3:
		dc.SetRGB255(255, 0, 255)
	case 4:
		dc.SetRGB255(88, 41, 0)
	case 0:
			dc.SetRGB255(0,255,0)
	}
	dc.DrawCircle(p.position[0],p.position[1], 5)
	dc.Fill()
	}

  //dc.SavePNG("out.png") // changer le format ici pour choper une image par particule ? -> il faudrait plutot choper l'ensemble (pas de save ici)

// mettre un mode ? particules, lime, boid.
func (p *Particule) Actualisation(delta_t float64, monde []*Particule, dc *gg.Context) {
	// nombre de voisins
	p.GetNombreVoisins(monde)
	//fmt.Println(p)
	Ntot := float64(p.Nb_voisins_droite + p.Nb_voisins_gauche)
	// update PHI
	//fmt.Println(Ntot)
	if math.Signbit(float64(p.Nb_voisins_droite - p.Nb_voisins_gauche)){
			p.PHI += ALPHA - p.BETA*Ntot
	}else{
			p.PHI += ALPHA + p.BETA*Ntot
	}
	// update position
	//last_p := p.position[0]
	angle := p.PHI*(math.Pi/2)/180.0
	p.position[0] += math.Cos(angle)*delta_t*p.vitesse // p.PHI en radianmath.Pi
	p.position[1] += math.Sin(angle)*delta_t*p.vitesse
	// update color
	if Ntot>35{
		p.couleur = 1 // jaune
	}else if (Ntot >16 && Ntot < 36){
		p.couleur = 2 // bleu
	}else if (Ntot == 15){
		p.couleur = 3 // magenta
	}else if ((Ntot >12 && Ntot < 15)){
		p.couleur = 4 // marron
	}else{
		p.couleur = 0 // vert
	}
	// en fonction du mode changer (mode 1 = PPS, mode 2 = boid
	p.Draw(dc)
}

func (p *Particule) GetNombreVoisins(monde []*Particule){
		p.Nb_voisins_droite = 0
		p.Nb_voisins_gauche = 0
		// on garde tout le monde, pas tres pertinent ...
		for _, part := range(monde){
			distance := math.Sqrt(math.Pow(p.position[0] - part.position[0],2) + math.Pow(p.position[1] - part.position[1],2))
			//fmt.Println("distance",distance)
			//fmt.Println("rayon",p.rayon)
			if distance < p.rayon {
				if p.position[0] - part.position[0] < 0{
					p.Nb_voisins_droite +=1
				}else{
					p.Nb_voisins_gauche += 1
				}
			}
			//fmt.Println("nb total",p.Nb_voisins_droite, p.Nb_voisins_gauche)
		}
}
