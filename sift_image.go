package gosift

import (
	"image"
	"image/color"
)

// SiftImage is an in memory RGBA image.  The red, green, blue, alpha samples are held in a
// single slice to increase resizing performance.
type SiftImage struct {
	// Pix holds the image's pixels, in red, green, blue, alpha order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []int
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *SiftImage) PixOffset(x, y int) int {
	return (y * p.Stride) + (x * 4)
}

func (p *SiftImage) Bounds() image.Rectangle {
	return p.Rect
}

func (p *SiftImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *SiftImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	return color.RGBA{
		ClampUint8(p.Pix[i+0]),
		ClampUint8(p.Pix[i+1]),
		ClampUint8(p.Pix[i+2]),
		ClampUint8(p.Pix[i+3]),
	}
}

func (p *SiftImage) Opaque() bool {
	return true
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *SiftImage) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &SiftImage{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &SiftImage{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// newSiftImage returns a new SiftImage with the given bounds.
func newSiftImage(r image.Rectangle) *SiftImage {
	w, h := r.Dx(), r.Dy()
	buf := make([]int, 4*w*h)
	return &SiftImage{Pix: buf, Stride: 4 * w, Rect: r}
}
