package interpolation

import (
	"image"
	"image/color"
)

// An InterpolationFunction provides the parameters that describe an
// interpolation kernel. It returns the number of samples to take
// and the kernel function to use for sampling.
type Interpolation interface {
	Smooth(srcImg, destImg image.Image, newWidth, newHeight, oldWidth, oldHeight int)
}

// Resize scales an image to new width and height using the interpolation function interp.
// A new image with the given dimensions will be returned.
// If one of the parameters width or height is set to 0, its size will be calculated so that
// the aspect ratio is that of the originating image.
func Resize(srcImg image.Image, newWidth, newHeight int, interp Interpolation) image.Image {
	oldBounds := srcImg.Bounds()
	oldWidth := oldBounds.Dx()
	oldHeight := oldBounds.Dy()

	destImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	interp.Smooth(srcImg, destImg, newWidth, newHeight, oldWidth, oldHeight)
	return destImg
}

// Calculates scaling factors using old and new image dimensions.
func calcFactors(newWidth, newHeight, oldWidth, oldHeight int) (scaleX, scaleY float64) {
	if newWidth == 0 {
		if newHeight == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = float64(oldHeight) / float64(newHeight)
			scaleX = scaleY
		}
	} else {
		scaleX = float64(oldWidth) / float64(newWidth)
		if newHeight == 0 {
			scaleY = scaleX
		} else {
			scaleY = float64(oldHeight) / float64(newHeight)
		}
	}
	return
}

func getRed(cr color.Color) uint8 {
	r, _, _, _ := cr.RGBA()
	return uint8(r)
}
func getGreen(cr color.Color) uint8 {
	_, g, _, _ := cr.RGBA()
	return uint8(g)
}
func getBlue(cr color.Color) uint8 {
	_, _, b, _ := cr.RGBA()
	return uint8(b)
}
func getAlpha(cr color.Color) uint8 {
	_, _, _, a := cr.RGBA()
	return uint8(a)
}

func setColor(destImg image.Image, x, y int, cr color.Color) {
	switch input := destImg.(type) {
	case *image.RGBA:
		input.Set(x, y, cr)
	case *image.NRGBA:
		input.Set(x, y, cr)
	case *image.RGBA64:
		input.Set(x, y, cr)
	default:

	}
}
