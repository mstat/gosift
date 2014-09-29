package gosift

import (
	"github.com/coraldane/resize"
	"image"
	"image/color"
	"image/draw"
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

/**
 * Substract two images. src1 - src2, the dimensions of src1, src2 must be same.
 * @param src1
 * @param src2
 * @return A new image with the given dimensions will be returned.
 */
func SubstractImage(src1, src2 image.Image) *SiftImage {
	bounds := src1.Bounds()
	source1 := ConvertImageToRGBA(src1)
	source2 := ConvertImageToRGBA(src2)
	destImg := newSiftImage(bounds)

	stride := source1.Stride
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			destImg.Pix[y*stride+x*4] = int(source1.Pix[y*stride+x*4] - source2.Pix[y*stride+x*4])
			destImg.Pix[y*stride+x*4+1] = int(source1.Pix[y*stride+x*4+1] - source2.Pix[y*stride+x*4+1])
			destImg.Pix[y*stride+x*4+2] = int(source1.Pix[y*stride+x*4+2] - source2.Pix[y*stride+x*4+2])
			destImg.Pix[y*stride+x*4+3] = int(source1.Pix[y*stride+x*4+3] - source2.Pix[y*stride+x*4+3])
		}
	}
	return destImg
}

func ConvertImageToRGBA(srcImg image.Image) *image.RGBA {
	bounds := srcImg.Bounds()
	destImg := image.NewRGBA(bounds)
	draw.Draw(destImg, bounds, srcImg, bounds.Min, draw.Src)
	return destImg
}

func ResizeImageHalf(srcImg image.Image) image.Image {
	oldBounds := srcImg.Bounds()
	width := oldBounds.Dx() / 2
	height := oldBounds.Dy() / 2
	return resize.Thumbnail(width, height, srcImg, resize.Bilinear)
}

func ResizeImageDouble(srcImg image.Image) image.Image {
	oldBounds := srcImg.Bounds()
	width := oldBounds.Dx() * 2
	height := oldBounds.Dy() * 2
	return resize.Resize(width, height, srcImg, resize.Bicubic)
}

func pixval32f(srcImg *SiftImage, x, y int) float64 {
	return float64(srcImg.Pix[y*srcImg.Stride+x*4])
}

// Keep value in [0,255] range.
func ClampUint8(in int) uint8 {
	if in < 0 {
		return 0
	}
	if in > 255 {
		return 255
	}
	return uint8(in)
}
