package main

import (
	"fmt"
	"github.com/coraldane/gosift/interpolation"
	"image/jpeg"
	"os"
)

func main() {
	file, err := os.Open("/Users/coral/Downloads/model/DSC_4142.jpg")
	defer file.Close()

	imgFile, err := os.Create("/Users/coral/Downloads/output.jpg")
	defer imgFile.Close()

	srcImg, err := jpeg.Decode(file)
	// destImg := interpolation.GrayImage(srcImg)
	bounds := srcImg.Bounds()
	width := bounds.Max.X * 2
	height := bounds.Max.Y * 2
	destImg := interpolation.Resize(srcImg, width, height, &interpolation.BicubicInterpolation{})

	err = jpeg.Encode(imgFile, destImg, &jpeg.Options{120})
	fmt.Println(err)
}
