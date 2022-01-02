package worldelements

import (
	utils "gitlab.utc.fr/projet_ia04/Boid/utils"
)

//Wall structure qui permet de définir un mur
type Wall struct {
	ImageWidth  int            //Longueur du mur
	ImageHeight int            //Hauteur du mur
	Position    utils.Vector2D //Position du mur
	TypeWall    int            //Type de mur
}

func GenerateMouthWallBomb(xPos float64, yPos float64) []*Wall {
	w, h := utils.WallImage.Size()
	var walls [10]*Wall
	changeInc := -1
	inc, subInc := 0, 0
	for b := 0; b < 10; b++ {
		x := xPos + float64(b*w)
		if inc == 4 && changeInc == -1 {
			changeInc = 0
		}
		if subInc == 3 && changeInc != 1 {
			changeInc = 1
			inc = 4
		}
		if changeInc == -1 {
			inc++
		} else if changeInc == 1 {
			inc--
		} else {
			subInc++
		}
		y := yPos + float64(inc*h)
		walls[b] = &Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    utils.Vector2D{X: x, Y: y},
			TypeWall:    0,
		}
	}
	return walls[:]
}

func GenerateEyeWallBomb(xPos float64, yPos float64) []*Wall {
	w, h := utils.WallImage.Size()
	var walls [8]*Wall
	changeInc := false
	inc := 1
	for b := 0; b < 5; b++ {
		x := xPos + float64(b*w)
		if b%3 == 0 {
			changeInc = true
		}
		if inc == 1 {
			changeInc = false
		}
		if changeInc {
			inc--
		} else {
			inc++
		}
		y := yPos - float64(inc*h)
		walls[b] = &Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    utils.Vector2D{X: x, Y: y},
		}
	}
	changeInc = false
	inc = 2
	for b := 0; b < 3; b++ {
		x := xPos + float64((b+1)*w)
		if b%2 == 0 {
			changeInc = true
		}
		if inc == 1 {
			changeInc = false
		}
		if changeInc {
			inc--
		} else {
			inc++
		}
		y := yPos + float64((inc-2)*h)
		walls[b+5] = &Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    utils.Vector2D{X: x, Y: y},
			TypeWall:    0,
		}
	}
	return walls[:]
}

func GenerateSideWallBomb(top bool) []*Wall {
	w, h := utils.WallImage.Size()
	var walls [48]*Wall
	//pour le mur du bas:
	changeInc := true
	inc := 4
	//pour le mur du haut
	if top {
		changeInc = false
		inc = 1
	}
	for b := 0; b < 48; b++ {
		x := b*w + 10
		if b%4 == 0 {
			changeInc = true
		}
		if inc == 1 {
			changeInc = false
		}
		if changeInc {
			inc--
		} else {
			inc++
		}
		y := -inc*h + utils.ScreenHeight
		if top {
			y = inc * h
		}
		walls[b] = &Wall{
			ImageWidth:  w,
			ImageHeight: h,
			Position:    utils.Vector2D{X: float64(x), Y: float64(y)},
			TypeWall:    0,
		}
	}
	return walls[:]
}
