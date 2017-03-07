package gocr

import (
	"encoding/csv"
	"encoding/gob"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/anthonynsimon/bild/segment"
)

type ModelImage struct {
	Label string
	Data  [][]uint8
}

type Model struct {
	Name        string
	ModelImages []ModelImage
}

// Read image in path
func ReadImage(path string) image.Image {
	// Read the file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	// Close the file later
	defer file.Close()

	// Decode the file to image (will decode any type of image .png, .jpg, .gif)
	src, _, err := image.Decode(file)
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

func ReadCSV(path string) [][]string {
	// Read the file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	// Close the file later
	defer file.Close()

	reader := csv.NewReader(file)
	var datas [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		datas = append(datas, record)
	}

	return datas
}

// Training method
// The train folder should include index.csv and images that inside the index.csv
func Train(sampleFolderPath string, modelPath string) {

	indexPath := sampleFolderPath + "/index.csv"
	indexData := ReadCSV(indexPath)

	// Initialize model
	model := Model{}

	// Read and binarize each image to array 0 and 1
	for _, elm := range indexData {
		image := ReadImage(sampleFolderPath + elm[0])
		binaryImageArray := ImageToBinaryArray(image)

		model.ModelImages = append(model.ModelImages, ModelImage{
			Label: elm[1],
			Data:  binaryImageArray,
		})
	}

	// Create the model file
	modelFile, err := os.Create(modelPath + "model.gob")
	if err != nil {
		panic(err)
	}

	// Close the file later
	defer modelFile.Close()

	// Create encoder
	encoder := gob.NewEncoder(modelFile)
	// Write the file
	if err := encoder.Encode(model); err != nil {
		panic(err)
	}

}
