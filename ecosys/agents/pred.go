package agents

import(
    gg "github.com/fogleman/gg";
  "math";
  "math/rand"
)
type Pred struct{
  Position [2]float64
  Vitesse [2]float64
  Max_speed [2]float64
  rayon_absorp float64
  Energie float64
}

func New_pred(position [2]float64, rayon float64) (p Pred){
  return Pred{
    position,
    [2]float64{1.5,1.5},
    [2]float64{10,10},
    rayon,
    1.0,
  }
}

func (p *Pred) Border(width float64, height float64){
  if (p.Position[0] < 0){
    p.Position[0] = width
  }
  if (p.Position[0] > width  ){
    p.Position[0] = 0
  }
  if (p.Position[1] < 0){
    p.Position[1] = height
  }
  if (p.Position[1] > height ){
    p.Position[1] = 0
  }
}

func (p *Pred) Update(){
  norme := math.Sqrt(math.Pow(p.Vitesse[0],2) + math.Pow(p.Vitesse[1],2))
	norme_max := math.Sqrt(math.Pow(p.Max_speed[0],2) + math.Pow(p.Max_speed[1],2))
  p.Vitesse[0] += 0.6*rand.Float64()
  p.Vitesse[1] += 0.6*rand.Float64()
  if norme >= norme_max{
		p.Vitesse[0] = 1.0 + 5*rand.Float64()
		p.Vitesse[1] = 1.0 + 5*rand.Float64()
	}
  p.Position[0] += p.Vitesse[0]
  p.Position[1] += p.Vitesse[1]
  p.Border(1000.0,1000.0)
}

func (p *Pred) GetEnergy(b *Boid){
  p.Energie += b.Energie
}

func (p *Pred) Draw(dc *gg.Context ){
	dc.SetRGB(100, 180, 0)
	dc.DrawCircle(p.Position[0],p.Position[1], 12)
	dc.Fill()
	}
