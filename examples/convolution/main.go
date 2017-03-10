package main

import (
	"fmt"

	"github.com/eaciit/gocr"
	"github.com/gonum/matrix/mat64"
)

func main() {

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
