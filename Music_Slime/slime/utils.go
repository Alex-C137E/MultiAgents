package slime

//Basic Hash function
func Hash(state uint) uint {
	state ^= 2747636419
	state *= 2654435769
	state ^= state >> 16
	state *= 2654435769
	state ^= state >> 16
	state *= 2654435769
	return state
}

func scaleToRange01(state uint) float64 {
	return float64(state) / 4294967295.0
}
