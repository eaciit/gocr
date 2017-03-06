package gocr

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// Read image in path to 2D array with size based on image width and height
func ReadImageWithPath(path string) [][]uint8 {
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

	// Initialize 2D array for the gray value of the image
	imageArr := make([][]uint8, gray.Bounds().Max.X, gray.Bounds().Max.Y)

	for x := 0; x < gray.Bounds().Max.X; x++ {
		for y := 0; y < gray.Bounds().Max.Y; y++ {
			imageArr[x][y] = gray.GrayAt(x, y).Y
		}
	}

	// Return the 2D array
	return imageArr
}
