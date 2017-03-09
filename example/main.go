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
	// Downloading english font dataset and index.csv
	d, _ := os.Getwd()
	//
	// Load the sample data and save it in .gob file
	// err := gocr.TrainAverage(d+"/English/Fnt/", d+"/English/")
	// if err != nil {
	// 	panic(err)
	// }

	// Read the image
	image, err := gocr.ReadImage(d + "/zero.png")
	if err != nil {
		panic(err)
	}
	data := gocr.ImageToBinaryArray(image)

	model, err := gocr.ReadModel(d + "/English/model.gob")
	model.ModelImages = append(model.ModelImages, gocr.ModelImage{
		Label: "0",
		Data:  data.NNInterpolation(128, 128),
	})

	s := gocr.NewScanner(&model)
	fmt.Println(s.Predict(data))
}
