package gocr

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// Read image in given path
// Can open and decode some type of image (.png, .jpg, .gif)
func ReadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return src, nil
}

// Convert image to grayscale 2D array
func ImageToGraysclaeArray(src image.Image) ImageMatrix {
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			grayColor := gray.ColorModel().Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	imageArr := NewImageMatrix(gray.Bounds().Max.Y, gray.Bounds().Max.X)

	for y := 0; y < gray.Bounds().Max.Y; y++ {
		for x := 0; x < gray.Bounds().Max.X; x++ {
			imageArr[y][x] = gray.GrayAt(x, y).Y
		}
	}

	return imageArr
}

// Basic Thresholding
func Threshold(im ImageMatrix, thrs uint8) ImageMatrix {
	r, c := im.Dims()
	o := NewImageMatrix(r, c)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if im.At(i, j) >= thrs {
				o.Set(i, j, 1)
			} else {
				o.Set(i, j, 0)
			}
		}
	}

	return o
}

// Thresholding Using Otsu's Method
func OtsuThresh(im ImageMatrix) ImageMatrix {
	r, c := im.Dims()
	hist := im.Historgram()
	sumAll := 0

	for i := range hist {
		sumAll += i * hist[i]
	}

	sumBack, wBack, wFore, varMax, thrs := 0, 0, 0, 0.0, 0
	total := r * c

	for i := range hist {
		wBack += hist[i]

		if wBack == 0 {
			continue
		}

		wFore = total - wBack

		if wFore == 0 {
			break
		}

		sumBack += i * hist[i]
		mb := float64(sumBack) / float64(wBack)
		mf := float64(sumAll-sumBack) / float64(wFore)

		vb := float64(wBack*wFore) * math.Pow(mb-mf, 2)

		if vb > varMax {
			varMax = vb
			thrs = i
		}
	}

	return Threshold(im, uint8(thrs))
}

func AdaptiveThres(im ImageMatrix, bs int) ImageMatrix {
	r, c := im.Dims()
	o := NewImageMatrix(r, c)

	for i := 0; i < r/bs+1; i++ {
		for j := 0; j < c/bs+1; j++ {
			br, rc := (i+1)*bs, (j+1)*bs

			if br >= r {
				br = r - 1
			}

			if rc >= c {
				rc = c - 1
			}

			s := NewSquare(NewCoordinate(i*bs, j*bs), NewCoordinate(br, rc))
			o.SetSquare(s, OtsuThresh(im.SliceSquare(s)))
		}
	}

	return o
}

// Binarize the given imageArr using
// Best algorithm based on this paper https://pdfs.semanticscholar.org/6347/5461213fdaa24e418c33454c72bdbbe8f8b4.pdf is Sauvola
// Sauvola Reference: http://www.mediateam.oulu.fi/publications/pdf/24.p
// TODO: Implement Sauvola algorithm
func SauvolaBinarization(imageArr ImageMatrix) ImageMatrix {

	return nil
}

// Convert ImageMatrix to Image and save it to given path
func ImageMatrixToImage(imageArray ImageMatrix, outPath string, mul int) error {
	r := len(imageArray)
	c := len(imageArray[0])

	gray := image.NewGray(image.Rect(0, 0, c, r))
	for x := 0; x < c; x++ {
		for y := 0; y < r; y++ {
			grayColor := color.Gray{}
			grayColor.Y = imageArray[y][x] * uint8(mul)
			gray.Set(x, y, grayColor)
		}
	}

	outfile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	png.Encode(outfile, gray)

	return nil
}

// Adding pad to make a square matrix
// Then resize it to given row length and column length
func PadAndResize(matrix ImageMatrix, dr, dc int) ImageMatrix {
	resizedMatrix := matrix
	tr, tc := resizedMatrix.Dims()

	if tr > tc {
		left := (tr - tc) / 2
		right := tr - tc - left
		resizedMatrix = resizedMatrix.Pad(0, 0, left, right, 1)
	} else if tc > tr {
		top := (tc - tr) / 2
		bottom := tc - tr - top
		resizedMatrix = resizedMatrix.Pad(top, bottom, 0, 0, 1)
	}

	if dr != tr || dc != tc {
		resizedMatrix = resizedMatrix.NNInterpolation(dr, dc)
	}

	return resizedMatrix
}

// Find the distance of 2 give Dense using Euclidean Distance
func EuclideanDistance(m1, m2 ImageMatrix) float64 {

	r1, c1 := m1.Dims()
	r2, c2 := m2.Dims()

	if r1 != r2 || c1 != c2 {
		panic("Dimension mismatch")
	}

	var sum float64 = 0.0

	for y := 0; y < r1; y++ {
		for x := 0; x < c1; x++ {
			sum += math.Pow(float64(m1.At(y, x)-m2.At(y, x)), 2)
		}
	}

	return math.Sqrt(sum)
}
