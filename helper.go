package gocr

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/gonum/matrix/mat64"
)

func ImageArrayToMatrix(imageArray [][]uint8) *mat64.Dense {
	fA := make([]float64, len(imageArray)*len(imageArray[0]))
	w := len(imageArray)
	h := len(imageArray[0])

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			fA[x*h+y] = float64(imageArray[x][y])
		}
	}

	return mat64.NewDense(w, h, fA)
}

func ImageArrayToImage(imageArray [][]uint8, outPath string) {
	r := len(imageArray)
	c := len(imageArray[0])

	// Change row, column paradigman to x, y
	gray := image.NewGray(image.Rect(0, 0, c, r))
	for x := 0; x < c; x++ {
		for y := 0; y < r; y++ {
			grayColor := color.Gray{}
			grayColor.Y = imageArray[y][x] * 255
			gray.Set(x, y, grayColor)
		}
	}

	// Encode the grayscale image to the output file
	outfile, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	png.Encode(outfile, gray)
}

func MatrixToImage(matrix mat64.Matrix, outPath string) {
	h, w := matrix.Dims()

	// Change row, column paradigman to x, y
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			grayColor := color.Gray{}
			grayColor.Y = uint8(matrix.At(y, x)) * 255
			gray.Set(x, y, grayColor)
		}
	}

	// Encode the grayscale image to the output file
	outfile, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	png.Encode(outfile, gray)
}
