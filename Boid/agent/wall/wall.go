package wall

import (
	vector "gitlab.utc.fr/projet_ia04/Boid/utils/vector"
)

type Wall struct {
	ImageWidth  int
	ImageHeight int
	Position    vector.Vector2D
	TypeWall    int
}
