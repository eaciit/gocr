package main

import (
	"os"

	"github.com/eaciit/gocr"
)

var (
	samplePath = func() string {
		d, _ := os.Getwd()
		return d
	}() + "/../../train_data/"
)

var (
	modelPath = func() string {
		d, _ := os.Getwd()
		return d
	}() + "/../../model/"
)

func main() {
	// Downloading english font dataset and index.csv
	d, _ := os.Getwd()
	//
	// Load the sample data and save it in .gob file
	// err := gocr.TrainAverage(d+"/English/Fnt/", d+"/English/")
	// if err != nil {
	// 	panic(err)
	// }

	// Read the image
	image, err := gocr.ReadImage(d + "/imagetext_2.png")
	if err != nil {
		panic(err)
	}

	im := gocr.ImageToGraysclaeArray(image)
	gocr.ImageMatrixToImage(im, d+"/result/image.png", 255)

}
