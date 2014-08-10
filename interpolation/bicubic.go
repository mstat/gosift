package interpolation

import (
	// "fmt"
	"image"
	"image/color"
	"math"
)

type BicubicInterpolation struct {
}

func (interp *BicubicInterpolation) Smooth(srcImg, destImg image.Image, newWidth, newHeight, oldWidth, oldHeight int) {
	var u, v, srcX, srcY, k1, k2 float64
	var i, j int
	var rr, gg, bb float64

	scaleX, scaleY := calcFactors(newWidth, newHeight, oldWidth, oldHeight)
	maxX := oldWidth - 1
	maxY := oldHeight - 1

	for y := 0; y < newHeight; y++ {
		srcY = scaleX * float64(y)
		j = int(srcY)
		v = srcY - float64(j)

		for x := 0; x < newWidth; x++ {
			srcX = scaleY * float64(x)
			i = int(srcX)
			u = srcX - float64(i)

			rr, gg, bb = 0.0, 0.0, 0.0
			for m := -1; m < 3; m++ {
				k1 = cubic(v - float64(m))
				oj := j + m
				if oj < 0 {
					oj = 0
				} else if oj > maxY {
					oj = maxY
				}
				for n := -1; n < 3; n++ {
					k2 = k1 * cubic(float64(n)-u)
					oi := i + n
					if oi < 0 {
						oi = 0
					} else if oi > maxX {
						oi = maxX
					}
					r, g, b, _ := srcImg.At(oi, oj).RGBA()
					rr += k2 * float64(r%256)
					gg += k2 * float64(g%256)
					bb += k2 * float64(b%256)
				}
			}

			or, og, ob, oa := srcImg.At(i, j).RGBA()
			// fmt.Println(i, j, u, v, a, c, colors)
			// cr := calcCubicColor(colors, a, c)
			cr := &color.RGBA{restoreColor(or, rr), restoreColor(og, gg), restoreColor(ob, bb), uint8(oa % 256)}
			setColor(destImg, x, y, cr)
		}
	}
}

func restoreColor(source uint32, dest float64) uint8 {
	src := uint8(source % 256)
	dst := uint8(dest)
	if math.Abs(float64(src-dst)) > 5 {
		return src
	}
	return dst
}

func cubic(in float64) float64 {
	in = math.Abs(in)
	if in <= 1 {
		return in*in*(1.5*in-2.5) + 1.0
	}
	if in <= 2 {
		return in*(in*(2.5-0.5*in)-4.0) + 2.0
	}
	return 0
}
