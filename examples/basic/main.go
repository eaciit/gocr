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
	im := gocr.ImageToBinaryArray(image)
	gocr.ImageMatrixToImage(im, d+"/result/image.png")

	rs := gocr.CirucularScan(im)

	for i := 0; i < len(rs); i++ {
		gocr.ImageMatrixToImage(im.SliceArea(rs[i]), d+"/result/image_"+strconv.Itoa(i)+".png")
	}

}
