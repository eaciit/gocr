package gocr

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/anthonynsimon/bild/segment"
)

// Read image in path
func ReadImage(path string) image.Image {
	// Read the file
	infile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	// Close the file later
	defer infile.Close()

	// Decode the file to image (will decode any type of image .png, .jpg, .gif)
	src, _, err := image.Decode(infile)
	if err != nil {
		panic(err)
	}

	return src
}

//Convert image to grayscale 2D array
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

	for x := 0; x < gray.Bounds().Max.X; x++ {
		// Intialize the column
		imageArr[x] = make([]uint8, gray.Bounds().Max.Y)
		for y := 0; y < gray.Bounds().Max.Y; y++ {
			imageArr[x][y] = gray.GrayAt(x, y).Y
		}
	}

	// Return the 2D array
	return imageArr
}

func ImageToBinaryArray(src image.Image) [][]uint8 {
	// FIXME: still finding the best Threshold
	gray := segment.Threshold(src, 128)

	// Initialize 2D array for the gray value of the image (row first)
	imageArr := make([][]uint8, gray.Bounds().Max.X)

	for x := 0; x < gray.Bounds().Max.X; x++ {
		// Intialize the column
		imageArr[x] = make([]uint8, gray.Bounds().Max.Y)
		for y := 0; y < gray.Bounds().Max.Y; y++ {
			imageArr[x][y] = gray.GrayAt(x, y).Y / 255
		}
	}

	return imageArr
}

// Binarize the given imageArr using
// Best algorithm based on this paper https://pdfs.semanticscholar.org/6347/5461213fdaa24e418c33454c72bdbbe8f8b4.pdf is Sauvola
// Sauvola Reference: http://www.mediateam.oulu.fi/publications/pdf/24.p
// TODO: Implement Sauvola algorithm
func SauvolaBinarization(imageArr [][]uint8) [][]uint8 {

	return nil
}
