package gocr

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/anthonynsimon/bild/segment"
	"github.com/gonum/matrix/mat64"
)

const THRESHOLD_VALUE = 128

// Read image in path
func ReadImage(path string) (image.Image, error) {
	// Read the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Close the file later
	defer file.Close()

	// Decode the file to image (will decode any type of image .png, .jpg, .gif)
	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return src, nil
}

// Convert image to grayscale 2D array
func ImageToGraysclaeArray(src image.Image) [][]uint8 {
	// Convert the image to Grayscale
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

	// Initialize 2D array for the gray value of the image (row first)
	imageArr := make([][]uint8, gray.Bounds().Max.X)

	// Change the X,Y paradigma to Rows, Column
	for y := 0; y < gray.Bounds().Max.Y; y++ {
		// Intialize the column
		imageArr[y] = make([]uint8, gray.Bounds().Max.X)
		for x := 0; x < gray.Bounds().Max.X; x++ {
			imageArr[y][x] = gray.GrayAt(x, y).Y
		}
	}

	// Return the 2D array
	return imageArr
}

// Convert image to binary array
func ImageToBinaryArray(src image.Image) [][]uint8 {
	gray := segment.Threshold(src, THRESHOLD_VALUE)

	// Initialize 2D array for the gray value of the image (row first)
	imageArr := make([][]uint8, gray.Bounds().Max.Y)

	// Change the X,Y paradigma to Rows, Column
	for y := 0; y < gray.Bounds().Max.Y; y++ {
		// Intialize the column
		imageArr[y] = make([]uint8, gray.Bounds().Max.X)
		for x := 0; x < gray.Bounds().Max.X; x++ {
			imageArr[y][x] = gray.GrayAt(x, y).Y / 255
		}
	}

	return imageArr
}

// Convert image to binary matrix
func ImageToBinaryMatrix(src image.Image) *mat64.Dense {
	gray := segment.Threshold(src, THRESHOLD_VALUE)

	// Initialize 2D array for the gray value of the image (row first)
	matrix := mat64.NewDense(gray.Bounds().Max.Y, gray.Bounds().Max.X, nil)

	// Change the X,Y paradigma to Rows, Column
	for y := 0; y < gray.Bounds().Max.Y; y++ {

		for x := 0; x < gray.Bounds().Max.X; x++ {
			matrix.Set(y, x, float64(gray.GrayAt(x, y).Y/255))
		}
	}

	return matrix
}

// Binarize the given imageArr using
// Best algorithm based on this paper https://pdfs.semanticscholar.org/6347/5461213fdaa24e418c33454c72bdbbe8f8b4.pdf is Sauvola
// Sauvola Reference: http://www.mediateam.oulu.fi/publications/pdf/24.p
// TODO: Implement Sauvola algorithm
func SauvolaBinarization(imageArr [][]uint8) [][]uint8 {

	return nil
}

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

func ImageArrayToImage(imageArray [][]uint8, outPath string) error {
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
		return err
	}
	defer outfile.Close()

	png.Encode(outfile, gray)

	return nil
}

func MatrixToImage(matrix mat64.Matrix, outPath string) error {
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
		return err
	}
	defer outfile.Close()

	png.Encode(outfile, gray)

	return nil
}

// Resize the array to given height and width using Nearest Neighbor Interpolation
func NNInterpolation(data *mat64.Dense, outputHeight, outputWidth int) *mat64.Dense {

	r, c := data.Dims()

	xRatio := float64(c) / float64(outputWidth)
	yRatio := float64(r) / float64(outputHeight)

	output := mat64.NewDense(outputHeight, outputWidth, nil)

	for y := 0; y < outputHeight; y++ {
		for x := 0; x < outputWidth; x++ {
			py := int(math.Floor(float64(y) * yRatio))
			px := int(math.Floor(float64(x) * xRatio))

			output.Set(y, x, data.At(py, px))
		}
	}

	return output
}

// Find the distance of 2 give Dense using Euclidean Distance
func EuclideanDistance(m1, m2 *mat64.Dense) float64 {

	r1, c1 := m1.Dims()
	r2, c2 := m2.Dims()

	if r1 != r2 || c1 != c2 {
		panic("Dimension mismatch")
	}

	var sum float64 = 0.0

	for y := 0; y < r1; y++ {
		for x := 0; x < c1; x++ {
			sum += math.Pow(m1.At(y, x)-m2.At(y, x), 2)
		}
	}

	return math.Sqrt(sum)
}
