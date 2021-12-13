package agents

import(
  gg "github.com/fogleman/gg"
)
type Source struct{
  Type int
  Position [2]float64
  Rayon float64
  Retention float64 // plus la retention est forte, plus le boid reste à cette source
  Energie_donnee float64 // unite donnee par unite de temps.
}

func New_source(Type int, position [2]float64) Source{
  return Source{
    Type,
    position,
    20,
    3.0,
    0.2,
  }
}

func (s *Source) Draw(dc *gg.Context ){
	dc.SetRGB(200, 400, 20)
	dc.DrawCircle(s.Position[0],s.Position[1], 20)
	dc.Fill()
	}
