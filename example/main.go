package main

import (
	"fmt"
	"os"

	"github.com/eaciit/gocr"
)

var (
	samplePath = func() string {
		d, _ := os.Getwd()
		return d
	}() + "/../train_data/"
)

var (
	modelPath = func() string {
		d, _ := os.Getwd()
		return d
	}() + "/../model/"
)

func main() {
	// image := gocr.ReadImage(samplePath + "sample1/0.gif")
	// binaryArr := gocr.ImageToBinaryArray(image)

	// Downloading english font dataset and index.csv
	d, _ := os.Getwd()
	gocr.Prepare(d + "/")

	// Load the sample data and save it in .gob file
	gocr.Train(samplePath+"sample1/", modelPath+"sample1/")
	// Load the model data and return it
	model := gocr.ReadModel(modelPath + "sample1/model.gob")
	fmt.Println(len(model.ModelImages))
}
