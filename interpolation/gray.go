package interpolation

import (
	"image"
	"image/color"
)

/**
 * Convert an image into gray color.
 * @param srcImage
 * @return A new image with the given dimensions will be returned.
 */
func GrayImage(srcImg image.Image) image.Image {
	bounds := srcImg.Bounds()
	destImg := image.NewRGBA(bounds)
	var colorVal uint8
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			cr := srcImg.At(x, y)
			r, g, b, a := cr.RGBA()
			grey := (299*(r%256) + 587*(g%256) + 114*(b%256)) / 1000
			colorVal = uint8(grey % 256)
			destImg.Set(x, y, color.RGBA{colorVal, colorVal, colorVal, uint8(a % 256)})
		}
	}
	return destImg
}
