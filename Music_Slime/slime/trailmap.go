package slime

import (
	"math"

	types "gitlab.utc.fr/projet_ia04/musicslime/types"
)

func Diffuse(trailMap types.TrailMap, turnUp int) types.TrailMap {
	//TODO diffuser le trail map au cours du temps
	width := types.GetWindowDefault().Width
	height := types.GetWindowDefault().Height

	diffusedTrailMap := trailMap

	for i := 0; i < len(trailMap); i++ {
		for j := 0; j < len(trailMap[i]); j++ {
			var sum [4]float64
			originalCol := trailMap[i][j]

			for offsetX := -1; offsetX <= 1; offsetX++ {
				for offsetY := -1; offsetY <= 1; offsetY++ {
					sampleX := math.Min(width-1.0, math.Max(0.0, float64(i+offsetX)))
					sampleY := math.Min(height-1.0, math.Max(0.0, float64(j+offsetY)))

					sum[0] = sum[0] + trailMap[int(sampleX)][int(sampleY)][0]
					sum[1] = sum[1] + trailMap[int(sampleX)][int(sampleY)][1]
					sum[2] = sum[2] + trailMap[int(sampleX)][int(sampleY)][2]
					sum[3] = sum[3] + trailMap[int(sampleX)][int(sampleY)][3]
				}
			}

			//moyenne de l'intensité des pxs autour du px a estomper
			var blurredCol [4]float64
			blurredCol[0] = sum[0] / 9
			blurredCol[1] = sum[1] / 9
			blurredCol[2] = sum[2] / 9
			blurredCol[3] = sum[3] / 9

			diffuseWeight := types.DiffuseRate * types.DeltaTime
			if diffuseWeight < 0.0 {
				diffuseWeight = 0.0
			}
			if diffuseWeight > 1 {
				diffuseWeight = 1.1
			}

			blurredCol[0] = originalCol[0]*(1-diffuseWeight) + blurredCol[0]*diffuseWeight
			blurredCol[1] = originalCol[1]*(1-diffuseWeight) + blurredCol[1]*diffuseWeight
			blurredCol[2] = originalCol[2]*(1-diffuseWeight) + blurredCol[2]*diffuseWeight
			blurredCol[3] = originalCol[3]*(1-diffuseWeight) + blurredCol[3]*diffuseWeight

			diffusedTrailMap[i][j][0] = math.Max(0.0, blurredCol[0]-types.DecayRate*types.DeltaTime)
			diffusedTrailMap[i][j][1] = math.Max(0.0, blurredCol[1]-types.DecayRate*types.DeltaTime)
			diffusedTrailMap[i][j][2] = math.Max(0.0, blurredCol[2]-types.DecayRate*types.DeltaTime)
			diffusedTrailMap[i][j][3] = math.Max(0.0, blurredCol[3]-types.DecayRate*types.DeltaTime)
		}
	}
	return diffusedTrailMap
}
