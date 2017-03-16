package gocr

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/ugorji/go/codec"
)

type Marker struct {
	Start int
	End   int
}

func (m Marker) thickness() int {
	return m.End - m.Start
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

// Find Marker of DarkSquare from given Matrix
// Direction 0 means it will iterate every rows
// Direction 1 means it will iterate every columns
func MarkersOfMatrix(data ImageMatrix, threshold float64, direction int) []Marker {
	r, c := data.Dims()
	markers := []Marker{}
	startMarker := -1
	isDarkSquare := false

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
			if isDarkSquare && startMarker > 0 {
				markers = append(markers, Marker{
					Start: startMarker,
					End:   i,
				})
				startMarker = -1
				isDarkSquare = false
			} else {
				startMarker = i
			}
		} else {
			isDarkSquare = true
		}
	}

	return markers
}

// Scan the DarkSquare and return it as []mat64.Dense for each line
// And [][]mat64.Dense for each characther
func LinearScan(data ImageMatrix) ([]ImageMatrix, [][]ImageMatrix) {
	r, c := data.Dims()

	markers := MarkersOfMatrix(data, 0.9, 0)
	lines := []ImageMatrix{}

	rowPadding := 10
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
		markers := MarkersOfMatrix(line, 0.95, 1)
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

func CirucularScan(image ImageMatrix) []*Square {
	r, c := image.Dims()
	start := NewCoordinate(0, 0)
	imageSquare := NewSquare(start, NewCoordinate(r, c))
	exploredSquares := []*Square{}
	resultsSquare := []*Square{}

	for i := 0; i < c; i++ {
		for j := 0; j < r; j++ {
			exist := false
			for _, ea := range exploredSquares {
				if ea.Include(j, i) {
					exist = true
					break
				}
			}

			if exist {
				continue
			}

			if image.At(j, i) == 0 {
				coor := NewCoordinate(j, i)
				result := NewSquare(coor, coor)
				vcs := []*Coordinate{}

				circleRun(image, coor, &vcs, imageSquare, result)

				resultsSquare = append(resultsSquare, result)
				exploredSquares = append(exploredSquares, result)
			}
		}
	}

	return resultsSquare
}

func circleRun(i ImageMatrix, c *Coordinate, vcs *[]*Coordinate, ia, rs *Square) {
	if c.IsInside(ia) {
		for _, cs := range *vcs {
			if c.row == cs.row && c.col == cs.col {
				return
			}
		}
		*vcs = append(*vcs, c)

		if i.At(c.row, c.col) == 0 {
			rs.Expand(c)

			circleRun(i, c.N(), vcs, ia, rs)
			circleRun(i, c.NE(), vcs, ia, rs)
			circleRun(i, c.E(), vcs, ia, rs)
			circleRun(i, c.SE(), vcs, ia, rs)
			circleRun(i, c.S(), vcs, ia, rs)
			circleRun(i, c.SW(), vcs, ia, rs)
			circleRun(i, c.W(), vcs, ia, rs)
			circleRun(i, c.NW(), vcs, ia, rs)
		}
	}
}

// ================================= Nearest Neighbor Scanner =================================
type NNScanner struct {
	model *Model
}

func NewNNScanner(model *Model) *NNScanner {
	return &NNScanner{
		model: model,
	}
}

func NewNNScannerFromFile(path string) *NNScanner {
	model, err := ReadModel(path)
	if err != nil {
		panic(err)
	}

	return &NNScanner{
		model: &model,
	}
}

// Predict using Nearest Neighbor method
func (s *NNScanner) Predict(matrix ImageMatrix) string {
	predictedLabel := ""
	mr, mc := s.model.ModelImages[0].Data.Dims()
	resizedMatrix := PadAndResize(matrix, mr, mc)

	min := math.MaxFloat64
	for _, modelImage := range s.model.ModelImages {
		distance := EuclideanDistance(resizedMatrix, modelImage.Data)

		fmt.Println(modelImage.Label, distance, " ")
		if min > distance {
			min = distance
			predictedLabel = modelImage.Label
		}
	}

	return predictedLabel
}

// ================================= Tensor CNN Scanner =================================

type CNNScanner struct {
	graph  *tf.Graph
	labels []string
}

func NewCNNScanner(graph *tf.Graph, labels []string) *CNNScanner {
	return &CNNScanner{
		graph:  graph,
		labels: labels,
	}
}

func NewCNNScannerFromDir(dir string) *CNNScanner {
	graph, labels := loadModel(dir)
	return NewCNNScanner(graph, labels)
}

func (s *CNNScanner) Predicts(images ImageMatrixs) []string {

	session, err := tf.NewSession(s.graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	kerasFlag, _ := tf.NewTensor(false)
	tensorImages, err := makeTensorFromImage(images)
	if err != nil {
		log.Fatal(err)
	}

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			s.graph.Operation("convolution2d_input_1").Output(0): tensorImages,
			s.graph.Operation("keras_learning_phase").Output(0):  kerasFlag,
		},
		[]tf.Output{
			s.graph.Operation("ArgMax_1").Output(0),
		},
		nil)

	if err != nil {
		log.Fatal(err)
	}

	predictions := output[0].Value().([]int64)
	for _, prediction := range predictions {
		fmt.Print(s.labels[prediction])
	}

	return nil
}

func makeTensorFromImage(images ImageMatrixs) (*tf.Tensor, error) {

	o := make([][][][]float32, len(images))

	for k := 0; k < len(images); k++ {
		o[k] = make([][][]float32, len(images[k]))
		for i := 0; i < len(images[k]); i++ {
			o[k][i] = make([][]float32, len(images[k][i]))
			for j := 0; j < len(images[k][i]); j++ {
				o[k][i][j] = make([]float32, 1)
				o[k][i][j][0] = float32(images[k][i][j])
			}
		}
	}

	tensor, err := tf.NewTensor(o)
	if err != nil {
		return nil, err
	}

	return tensor, nil
}

func loadModel(dir string) (*tf.Graph, []string) {
	var (
		modelFile  = filepath.Join(dir, "model.pb")
		labelsFile = filepath.Join(dir, "labels.txt")
	)

	model, err := ioutil.ReadFile(modelFile)
	if err != nil {
		log.Fatal(err)
	}

	graph := tf.NewGraph()
	if err = graph.Import(model, ""); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(labelsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}

	return graph, labels
}
