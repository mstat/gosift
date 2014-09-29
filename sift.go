/**
* Functions for detecting SIFT image features.
 */
package gosift

import (
	"fmt"
	"github.com/coraldane/resize"
	"image"
	"math"
)

const (
	/** default number of sampled intervals per octave. */
	SIFT_INTVLS = 3
	/** assumed gaussian blur for input image */
	SIFT_INIT_SIGMA = 0.5
	/** default sigma for initial guassian smoothing */
	SIFT_SIGMA = 1.6
	/** double image size before pyramid construction */
	SIFT_IMG_DBL = true
	/** default threshold on keypoint contrast |D(x)| */
	SIFT_CONTR_THR = 0.04
	/** default threshold on keypoint ratio of principle curvatures */
	SIFT_CURV_THR = 10
	/* width of border in which to ignore keypoints */
	SIFT_IMG_BORDER = 5
	/* maximum steps of keypoint interpolation before failure */
	SIFT_MAX_INTERP_STEPS = 5
)

func Sift(srcImg image.Image) {
	sift_features(srcImg, SIFT_INTVLS, SIFT_INIT_SIGMA, SIFT_IMG_DBL)
}

func sift_features(srcImg image.Image, intvls int, sigma float64, img_dbl bool) {
	/* build scale space pyramid; smallest dimension of top level is ~4 pixels */
	init_img := create_init_img(srcImg, img_dbl, sigma)

	initBounds := init_img.Bounds()
	octvs := int(math.Log(math.Min(float64(initBounds.Dx()), float64(initBounds.Dy())))/math.Log(2) - 2)
	gauss_pyr := build_gauss_pyr(init_img, octvs, intvls, sigma)
	dog_pyr := build_dog_pyr(gauss_pyr, octvs, intvls)
	fmt.Println(len(dog_pyr))
	scale_space_extrema(dog_pyr, SIFT_CONTR_THR, SIFT_CURV_THR)
}

/*
  Converts an image to 8-bit grayscale and Gaussian-smooths it.  The image is
  optionally doubled in size prior to smoothing.

  @param srcImg input image
  @param img_dbl if true, image is doubled in size prior to smoothing
  @param sigma total std of Gaussian smoothing
*/
func create_init_img(srcImg image.Image, img_dbl bool, sigma float64) image.Image {
	var sig_diff float64
	greyImg := GrayImage(srcImg)
	if img_dbl {
		greyImg = ResizeImageDouble(greyImg)
		sig_diff = (sigma*sigma - SIFT_INIT_SIGMA*SIFT_INIT_SIGMA*4)
	} else {
		sig_diff = (sigma*sigma - SIFT_INIT_SIGMA*SIFT_INIT_SIGMA)
	}
	greyImg = resize.GaussianSmooth(greyImg, sig_diff, 3)
	return greyImg
}

/*
  Builds Gaussian scale space pyramid from an image

  @param srcImg base image of the pyramid
  @param octvs number of octaves of scale space
  @param intvls number of intervals per octave
  @param sigma amount of Gaussian smoothing per octave

  @return Returns a Gaussian scale space pyramid as an octvs x (intvls + 3)
    array
*/
func build_gauss_pyr(srcImg image.Image, octvs, intvls int, sigma float64) [][]image.Image {
	k := math.Pow(2.0, 1.0/float64(intvls))
	sig := make([]float64, intvls+3)
	sig[0] = sigma
	sig[1] = sigma * math.Sqrt(k*k-1)
	for i := 2; i < intvls+3; i++ {
		sig[i] = sig[i-1] * k
	}

	gauss_pyr := make([][]image.Image, octvs)
	for row := 0; row < octvs; row++ {
		gauss_pyr[row] = make([]image.Image, intvls+3)
		for i := 0; i < intvls+3; i++ {
			if 0 == row && 0 == i {
				gauss_pyr[row][i] = srcImg
			} else if 0 == i {
				gauss_pyr[row][i] = downsample(gauss_pyr[row-1][intvls])
			} else {
				gauss_pyr[row][i] = resize.GaussianSmooth(gauss_pyr[row][i-1], sig[i], 3)
			}
		}
	}
	return gauss_pyr
}

/*
  Builds a difference of Gaussians scale space pyramid by subtracting adjacent
  intervals of a Gaussian pyramid

  @param gauss_pyr Gaussian scale-space pyramid
  @param octvs number of octaves of scale space
  @param intvls number of intervals per octave

  @return Returns a difference of Gaussians scale space pyramid as an
    octvs x (intvls + 2) array
*/
func build_dog_pyr(gauss_pyr [][]image.Image, octvs, intvls int) [][]*SiftImage {
	dog_pyr := make([][]*SiftImage, octvs)
	for o := 0; o < octvs; o++ {
		dog_pyr[o] = make([]*SiftImage, intvls+2)
		for i := 0; i < intvls+2; i++ {
			dog_pyr[o][i] = SubstractImage(gauss_pyr[o][i+1], gauss_pyr[o][i])
		}
	}
	return dog_pyr
}

/*
  Downsamples an image to a quarter of its size (half in each dimension)
  using nearest-neighbor interpolation

  @param img an image

  @return Returns an image whose dimensions are half those of img
*/
func downsample(srcImg image.Image) image.Image {
	bounds := srcImg.Bounds()
	return resize.Resize(bounds.Dx()/2, bounds.Dy()/2, srcImg, resize.Bilinear)
}

/*
  Detects features at extrema in DoG scale space.  Bad features are discarded
  based on contrast and ratio of principal curvatures.

  @param dog_pyr DoG scale space pyramid
  @param contr_thr low threshold on feature contrast
  @param curv_thr high threshold on feature ratio of principal curvatures
  @param storage memory storage in which to store detected features

  @return Returns an array of detected features whose scales, orientations,
    and descriptors are yet to be determined.
*/
func scale_space_extrema(dog_pyr [][]*SiftImage, contr_thr float64, curv_thr int) {
	octvs := len(dog_pyr)
	intvls := len(dog_pyr[0]) - 2
	prelim_contr_thr := 0.5 * contr_thr / float64(intvls) * 256
	for o := 0; o < octvs; o++ {
		imgBounds := dog_pyr[o][0].Bounds()
		for i := 1; i <= intvls; i++ {
			for r := SIFT_IMG_BORDER; r < imgBounds.Dy()-SIFT_IMG_BORDER; r++ {
				for c := SIFT_IMG_BORDER; c < imgBounds.Dx()-SIFT_IMG_BORDER; c++ {
					if math.Abs(pixval32f(dog_pyr[o][i], r, c)) > prelim_contr_thr {
						if is_extremum(dog_pyr, octvs, intvls, r, c) {

						}
					}
				}
			}
		}
	}
}

/*
  Determines whether a pixel is a scale-space extremum by comparing it to it's
  3x3x3 pixel neighborhood.

  @param dog_pyr DoG scale space pyramid
  @param octv pixel's scale space octave
  @param intvl pixel's within-octave interval
  @param r pixel's image row
  @param c pixel's image col

  @return Returns true if the specified pixel is an extremum (max or min) among
    it's 3x3x3 pixel neighborhood.
*/
func is_extremum(dog_pyr [][]*SiftImage, octv, intvl, r, c int) bool {
	val := pixval32f(dog_pyr[octv][intvl], r, c)
	/* check for maximum */
	if val > 0 {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				for k := -1; k <= 1; k++ {
					if val < pixval32f(dog_pyr[octv][intvl+i], r+j, c+k) {
						return false
					}
				}
			}
		}
	} else { /* check for minimum */
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				for k := -1; k <= 1; k++ {
					if val > pixval32f(dog_pyr[octv][intvl+i], r+j, c+k) {
						return false
					}
				}
			}
		}
	}
	return true
}

/*
  Computes the partial derivatives in x, y, and scale of a pixel in the DoG
  scale space pyramid.
  @param dog_pyr DoG scale space pyramid
  @param octv pixel's octave in dog_pyr
  @param intvl pixel's interval in octv
  @param r pixel's image row
  @param c pixel's image col

  @return Returns the slice of partial derivatives for pixel I
    { dI/dx, dI/dy, dI/ds }^T as a CvMat*
*/
func deriv_3D(dog_pyr [][]*SiftImage, octv, intvl int, r, c int) []float64 {
	di := make([]float64, 3)
	di[0] = pixval32f(dog_pyr[octv][intvl], r, c+1) - pixval32f(dog_pyr[octv][intvl], r, c-1)/2.0
	di[1] = pixval32f(dog_pyr[octv][intvl], r+1, c) - pixval32f(dog_pyr[octv][intvl], r-1, c)/2.0
	di[2] = pixval32f(dog_pyr[octv][intvl+1], r, c+1) - pixval32f(dog_pyr[octv][intvl-1], r, c)/2.0
	return di
}

/*
  Computes the 3D Hessian matrix for a pixel in the DoG scale space pyramid.

  @param dog_pyr DoG scale space pyramid
  @param octv pixel's octave in dog_pyr
  @param intvl pixel's interval in octv
  @param r pixel's image row
  @param c pixel's image col

  @return Returns the Hessian matrix (below) for pixel I as a CvMat*

  / Ixx  Ixy  Ixs \ <BR>
  | Ixy  Iyy  Iys | <BR>
  \ Ixs  Iys  Iss /
*/
func hessian_3D(dog_pyr [][]*SiftImage, octv, intvl int, r, c int) [][]float64 {
	hs := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		hs[i] = make([]float64, 3)
	}

	v := pixval32f(dog_pyr[octv][intvl], r, c)
	dxx := pixval32f(dog_pyr[octv][intvl], r, c+1) + pixval32f(dog_pyr[octv][intvl], r, c-1) - 2.0*v
	dyy := pixval32f(dog_pyr[octv][intvl], r+1, c) + pixval32f(dog_pyr[octv][intvl], r-1, c) - 2.0*v
	dss := pixval32f(dog_pyr[octv][intvl+1], r, c+1) + pixval32f(dog_pyr[octv][intvl-1], r, c) - 2.0*v
	dxy := (pixval32f(dog_pyr[octv][intvl], r+1, c+1) -
		pixval32f(dog_pyr[octv][intvl], r+1, c-1) -
		pixval32f(dog_pyr[octv][intvl], r-1, c+1) +
		pixval32f(dog_pyr[octv][intvl], r-1, c-1)) / 4.0
	dxs := (pixval32f(dog_pyr[octv][intvl+1], r, c+1) -
		pixval32f(dog_pyr[octv][intvl+1], r, c-1) -
		pixval32f(dog_pyr[octv][intvl-1], r, c+1) +
		pixval32f(dog_pyr[octv][intvl-1], r, c-1)) / 4.0
	dys := (pixval32f(dog_pyr[octv][intvl+1], r+1, c) -
		pixval32f(dog_pyr[octv][intvl+1], r-1, c) -
		pixval32f(dog_pyr[octv][intvl-1], r+1, c) +
		pixval32f(dog_pyr[octv][intvl-1], r-1, c)) / 4.0
	hs[0][0] = dxx
	hs[0][1] = dxy
	hs[0][2] = dxs
	hs[1][0] = dxy
	hs[1][1] = dyy
	hs[1][2] = dys
	hs[2][0] = dxs
	hs[2][1] = dys
	hs[2][2] = dss
	return hs
}
