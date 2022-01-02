package utils

import "math"

//IsPointInPolygon fonction qui détermine si un point est à l'intérieur d'un polygone
func IsPointInPolygon(p Vector2D, polygon []Vector2D) bool {
	minX := polygon[0].X
	maxX := polygon[0].X
	minY := polygon[0].Y
	maxY := polygon[0].Y
	for i := 1; i < len(polygon); i++ {
		minX = math.Min(polygon[i].X, minX)
		maxX = math.Max(polygon[i].X, maxX)
		minY = math.Min(polygon[i].Y, minY)
		maxY = math.Max(polygon[i].Y, maxY)
	}

	if p.X < minX || p.X > maxX || p.Y < minY || p.Y > maxY {
		return false
	}
	return true
}

//GetPolygonSize fonction qui permet d'obtenir la taille actuelle du lasso
func GetPolygonSize(polygon []Vector2D) float64 {
	if len(polygon) < 2 {
		return 0.0
	}
	sumDistance := 0.0
	for i := 1; i < len(polygon); i++ {
		p1 := polygon[i-1]
		p2 := polygon[i]
		distance := math.Sqrt((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y))
		sumDistance = sumDistance + distance
	}
	//add last distance
	p0 := polygon[0]
	pn := polygon[len(polygon)-1]
	return sumDistance + math.Sqrt((p0.X-pn.X)*(p0.X-pn.X)+(p0.Y-pn.Y)*(p0.Y-pn.Y))
}
