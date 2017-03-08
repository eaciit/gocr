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
	// gocr.Prepare(d + "/")
	//
	// Load the sample data and save it in .gob file
	// err := gocr.Train(samplePath+"/Sample2/", modelPath+"/Sample2/")
	// if err != nil {
	// 	panic(err)
	// }

	// Load the model data and return it
	model, err := gocr.ReadModel(modelPath + "/Sample2/model.gob")
	if err != nil {
		panic(err)
	}

	fmt.Println("ModelImage count: ", len(model.ModelImages))

	// Read the image
	image, err := gocr.ReadImage(d + "/imagetext_1.png")
	if err != nil {
		panic(err)
	}

	// Convert to binaryArr
	iA := gocr.ImageToBinaryArray(image)
	gocr.ImageArrayToImage(iA, d+"/result/result.png")
	// Convert to gonum Matrix mat64.Dense
	data := gocr.ImageArrayToMatrix(iA)
	gocr.MatrixToImage(data, d+"/result/result2.png")

	// Scan and slice
	_, charss := gocr.LinearScan(data)

	// Save each line to image
	// for i, line := range lines {
	// 	gocr.MatrixToImage(line, d+"/result/line_"+strconv.Itoa(i)+".png")
	// }

	// Save each char to image
	// for i, chars := range charss {
	// 	for j, char := range chars {
	// 		gocr.MatrixToImage(char, d+"/result/char_"+strconv.Itoa(i)+strconv.Itoa(j)+".png")
	// 	}
	// }

	for _, chars := range charss {
		for _, char := range chars {
			fmt.Print(gocr.Predict(char, &model))
		}

		fmt.Println("")
	}

}
