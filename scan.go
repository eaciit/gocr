package gocr

import (
	"math"
	"os"

	"github.com/ugorji/go/codec"
)

type Marker struct {
	Start int
	End   int
}

func (m Marker) thickness() int {
	return m.End - m.Start
}

type Scanner struct {
	model *Model
}

func NewScanner(model *Model) *Scanner {
	return &Scanner{
		model: model,
	}
}

func NewScannerFromFile(path string) *Scanner {
	model, err := ReadModel(path)
	if err != nil {
		panic(err)
	}

	return &Scanner{
		model: &model,
	}
}

// Predict the given image Dense to given model using kNN with k = 1
func (s *Scanner) Predict(matrix ImageMatrix) string {

	min := math.MaxFloat64
	predictedLabel := ""
	resizedMatrix := matrix
	mr, mc := s.model.ModelImages[0].Data.Dims()
	tr, tc := resizedMatrix.Dims()

	if tr > tc {
		left := (tr - tc) / 2
		right := tr - tc - left
		resizedMatrix = resizedMatrix.Pad(0, 0, left, right, 1)
	} else if tc > tr {
		top := (tc - tr) / 2
		bottom := tc - tr - top
		resizedMatrix = resizedMatrix.Pad(top, bottom, 0, 0, 1)
	}

	if mr != tr || mc != tc {
		resizedMatrix = matrix.NNInterpolation(mr, mc)
	}

	for _, modelImage := range s.model.ModelImages {
		distance := EuclideanDistance(resizedMatrix, modelImage.Data)

		if min > distance {
			min = distance
			predictedLabel = modelImage.Label
		}
	}

	return predictedLabel
}

// Read the model from a file and return the Model
func ReadModel(path string) (Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return Model{}, err
	}
	defer file.Close()

	decoder := codec.NewDecoder(file, new(codec.CborHandle))
	model := Model{}
	decoder.Decode(&model)

	return model, nil
}

// Find Marker of DarkArea from given Matrix
// Direction 0 means it will iterate every rows
// Direction 1 means it will iterate every columns
func MarkersOfMatrix(data ImageMatrix, threshold float64, direction int) []Marker {
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
			colAvg = float64(data.Row(i).Sum()) / float64(c)
		} else {
			colAvg = float64(data.Col(i).Sum()) / float64(r)
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
func LinearScan(data ImageMatrix) ([]ImageMatrix, [][]ImageMatrix) {
	r, c := data.Dims()

	markers := MarkersOfMatrix(data, 0.9, 0)
	lines := []ImageMatrix{}

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
		lines = append(lines, line)
	}

	charss := [][]ImageMatrix{}

	for _, line := range lines {
		markers := MarkersOfMatrix(line, 0.97, 1)
		r, c := line.Dims()
		chars := []ImageMatrix{}

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
			chars = append(chars, char)
		}

		charss = append(charss, chars)
	}

	return lines, charss
}
