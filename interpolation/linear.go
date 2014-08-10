package interpolation

import (
	"image"
	"image/color"
	"math"
)

type LinearInterpolation struct {
}

func (interp *LinearInterpolation) Smooth(srcImg, destImg image.Image, newWidth, newHeight, oldWidth, oldHeight int) {
	var u, v, srcX, srcY float64
	var i, j, i1, j1 int
	var colors [4]color.Color

	scaleX, scaleY := calcFactors(newWidth, newHeight, oldWidth, oldHeight)
	for y := 0; y < newHeight; y++ {
		srcY = scaleX * float64(y)
		j = int(srcY)
		j1 = int(math.Min(float64(j+1), float64(oldHeight-1)))
		v = srcY - float64(j)
		for x := 0; x < newWidth; x++ {
			srcX = scaleY * float64(x)
			i = int(srcX)
			i1 = int(math.Min(float64(i+1), float64(oldWidth-1)))
			u = srcX - float64(i)

			colors[0] = srcImg.At(i, j)
			colors[1] = srcImg.At(i, j1)
			colors[2] = srcImg.At(i1, j)
			colors[3] = srcImg.At(i1, j1)
			cr := mixColor(colors, i, j, u, v)
			setColor(destImg, x, y, cr)
		}
	}
}

func linearRGBA(colors [4]color.Color, i, j int, u, v float64, fn func(cr color.Color) uint8) uint8 {
	val := (1-u)*(1-v)*float64(fn(colors[0])) + (1-u)*v*float64(fn(colors[1])) +
		u*(1-v)*float64(fn(colors[2])) + u*v*float64(fn(colors[3]))
	return uint8(val)
}

func mixColor(colors [4]color.Color, i, j int, u, v float64) color.Color {
	r := linearRGBA(colors, i, j, u, v, getRed)
	g := linearRGBA(colors, i, j, u, v, getGreen)
	b := linearRGBA(colors, i, j, u, v, getBlue)
	a := linearRGBA(colors, i, j, u, v, getAlpha)
	return &color.RGBA{r, g, b, a}
}
