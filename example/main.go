package main

import (
	"os"
	"strconv"

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

type Marker struct {
	Start int
	End   int
}

func (m Marker) thickness() int {
	return m.End - m.Start
}

func main() {
	// image := gocr.ReadImage(samplePath + "sample1/0.gif")
	// binaryArr := gocr.ImageToBinaryArray(image)

	// Downloading english font dataset and index.csv
	d, _ := os.Getwd()
	// gocr.Prepare(d + "/")
	//
	// // Load the sample data and save it in .gob file
	// gocr.Train(samplePath+"sample1/", modelPath+"sample1/")
	// // Load the model data and return it
	// model := gocr.ReadModel(modelPath + "sample1/model.gob")
	// fmt.Println(len(model.ModelImages))

	image := gocr.ReadImage(d + "/imagetext_3.png")
	iA := gocr.ImageToBinaryArray(image)
	gocr.ImageArrayToImage(iA, d+"/result/result.png")
	data := gocr.ImageArrayToMatrix(iA)
	gocr.MatrixToImage(data, d+"/result/result2.png")

	lines, charss := gocr.LinearScan(data)

	for i, line := range lines {
		gocr.MatrixToImage(line, d+"/result/line_"+strconv.Itoa(i)+".png")
	}

	for i, chars := range charss {
		for j, char := range chars {
			gocr.MatrixToImage(char, d+"/result/char_"+strconv.Itoa(i)+strconv.Itoa(j)+".png")
		}
	}

}
