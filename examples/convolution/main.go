package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/eaciit/gocr"
	"github.com/gonum/matrix/mat64"
)

var (
	modelPath = func() string {
		d, _ := os.Getwd()
		return d
	}() + "/../../model/"
)

func main() {
	d, _ := os.Getwd()

	image, _ := gocr.ReadImage(d + "/imagetext_3.png")
	s := gocr.NewCNNPredictorFromDir(modelPath + "tensor_4/")
	s.InputHeight, s.InputWidth = 64, 64

	strings := gocr.ScanToStrings(s, image)
	for _, s := range strings {
		fmt.Println(s)
	}
}

func scanAndPrint() {
	d, _ := os.Getwd()

	image, _ := gocr.ReadImage(d + "/imagetext_3.png")
	imageMatrix := gocr.ImageToBinaryArray(image)
	squaress, charss := gocr.CirucularScan(imageMatrix)

	inputSize := 64
	s := gocr.NewCNNPredictorFromDir(modelPath + "tensor_3/")

	for k, chars := range charss {
		datas := make([]gocr.ImageMatrix, len(chars))
		for i := 0; i < len(chars); i++ {
			datas[i] = gocr.PadAndResize(chars[i], inputSize, inputSize)
			gocr.ImageMatrixToImage(datas[i], d+"/result/char_"+strconv.Itoa(i)+".png", 255)
		}

		texts := s.Predicts(datas)
		for i, text := range texts {
			if i < len(squaress[k])-1 {
				if squaress[k][i].NearestHorizontalDistanceTo(squaress[k][i+1]) > float64(squaress[k][i].Width()) {
					print(text + " ")
					continue
				}
			}

			print(text)
		}

		fmt.Println("")
	}
}

func tryNum() {
	k1 := mat64.NewDense(2, 2, []float64{1, 2, 3, 4})
	k2 := mat64.NewDense(2, 2, []float64{5, 6, 7, 8})
	kernel := [][]*mat64.Dense{{k1, k2}, {k1, k2}}

	d1 := mat64.NewDense(4, 4, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	d2 := mat64.NewDense(4, 4, []float64{17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32})
	data := [][]*mat64.Dense{{d1, d2}}

	r := gocr.Convolve(data, kernel, 0, 2)
	fmt.Println("Result -")
	printMatrix(r[0][0])
	printMatrix(r[0][1])

	m, _ := gocr.MaxPool(d1, 2, 2)
	fmt.Println("Result - ")
	printMatrix(m)
}

func printMatrix(m *mat64.Dense) {
	r, c := m.Dims()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Print(m.At(i, j), " ")
		}
		fmt.Println("")
	}
}
