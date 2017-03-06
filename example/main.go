package main

import (
	"os"

	"github.com/eaciit/gocr/train"
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

	// Load the sample data and save it in .gob file
	gocr.Train(samplePath+"sample1/", modelPath+"sample1/")

}
