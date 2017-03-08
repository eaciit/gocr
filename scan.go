package gocr

import (
	"math"
	"os"

	"github.com/gonum/matrix/mat64"
	"github.com/ugorji/go/codec"
)

type Marker struct {
	Start int
	End   int
}

func (m Marker) thickness() int {
	return m.End - m.Start
}

// Read the model from a file and return Model
func ReadModel(path string) (Model, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return Model{}, err
	}

	// Close the file later
	defer file.Close()

	// Initialize decoder and empty Model
	decoder := codec.NewDecoder(file, new(codec.CborHandle))
	model := Model{}

	// Decode the file
	decoder.Decode(&model)

	return model, nil
}

// Find Marker of DarkArea from given Matrix
// Direction 0 means it will iterate every rows
// Direction 1 means it will iterate every columns
func MarkersOfMatrix(data *mat64.Dense, threshold float64, direction int) []Marker {
	r, c := data.Dims()
	markers := []Marker{}
	startMarker := -1
	isDarkArea := false

	n := r
	if direction == 1 {
		n = c
	}

	for i := 0; i < n; i++ {

		var colAvg float64 = 1

		if direction == 0 {
			colAvg = mat64.Sum(data.RowView(i)) / float64(c)
		} else {
			colAvg = mat64.Sum(data.ColView(i)) / float64(r)
		}

		if colAvg >= threshold {
			if isDarkArea && startMarker > 0 {
				markers = append(markers, Marker{
					Start: startMarker,
					End:   i,
				})
				startMarker = -1
				isDarkArea = false
			} else {
				startMarker = i
			}
		} else {
			isDarkArea = true
		}
	}

	return markers
}

// Scan the DarkArea and return it as []mat64.Dense for each line
// And [][]mat64.Dense for each characther
func LinearScan(data *mat64.Dense) ([]*mat64.Dense, [][]*mat64.Dense) {
	r, c := data.Dims()

	markers := MarkersOfMatrix(data, 0.9, 0)
	lines := []*mat64.Dense{}

	rowPadding := 5
	for _, marker := range markers {

		if marker.thickness() <= 5 {
			continue
		}

		start := marker.Start - rowPadding
		if start < 0 {
			start = 0
		}

		end := marker.End + rowPadding
		if end >= r {
			end = r
		}

		line := data.Slice(start, end, 0, c)
		lines = append(lines, mat64.DenseCopyOf(line))
	}

	charss := [][]*mat64.Dense{}

	for _, line := range lines {
		markers := MarkersOfMatrix(line, 0.97, 1)
		r, c := line.Dims()
		chars := []*mat64.Dense{}

		columnPadding := 2
		for _, marker := range markers {

			if marker.thickness() <= 5 {
				continue
			}

			start := marker.Start - columnPadding
			if start < 0 {
				start = 0
			}

			end := marker.End + columnPadding
			if end >= c {
				end = c
			}

			char := line.Slice(0, r, start, end)
			chars = append(chars, mat64.DenseCopyOf(char))
		}

		charss = append(charss, chars)
	}

	return lines, charss
}

// Predict the given image Dense to given model
func Predict(matrix *mat64.Dense, model *Model) string {

	// kNN with k = 1
	min := math.MaxFloat64
	predictedLabel := ""
	resizedMatrix := matrix
	mr, mc := model.ModelImages[0].Data.Dims()
	tr, tc := matrix.Dims()

	if mr != tr || mc != tc {
		resizedMatrix = NNInterpolation(matrix, mr, mc)
	}

	for _, modelImage := range model.ModelImages {
		distance := EuclideanDistance(resizedMatrix, modelImage.Data)

		if min > distance {
			min = distance
			predictedLabel = modelImage.Label
		}
	}

	return predictedLabel
}
