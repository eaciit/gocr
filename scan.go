package gocr

import (
	"bufio"
	"fmt"
	"image"
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
func LinearScan(data ImageMatrix) [][]ImageMatrix {
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

	return charss
}

func CirucularScan(image ImageMatrix) ([][]*Square, [][]ImageMatrix) {
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

	charss := [][]ImageMatrix{}
	squaress := [][]*Square{}

	for _, result := range resultsSquare {
		if result.Width() < result.Height() {
			hat, i := findTopSquare(result, resultsSquare)
			if hat != nil {
				result.Merge(hat)
				resultsSquare = append(resultsSquare[:i], resultsSquare[i+1:]...)
			}
		}
	}

	for _, result := range resultsSquare {
		match := false
		for i, squares := range squaress {
			if result.AverageVerticalDistanceTo(squares[0]) < float64(squares[0].Height()) {
				charss[i] = append(charss[i], image.SliceSquare(result))
				squaress[i] = append(squaress[i], result)
				match = true
				break
			}
		}

		if !match {
			charss = append(charss, []ImageMatrix{image.SliceSquare(result)})
			squaress = append(squaress, []*Square{result})
		}
	}

	return squaress, charss
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

func findTopSquare(s *Square, squares []*Square) (*Square, int) {
	m := s.topLeft.col + (s.bottomRight.col-s.topLeft.col)/2
	ni := -1
	nd := float64((s.bottomRight.row - s.topLeft.row)) / 2

	for i := 0; i < len(squares); i++ {
		if squares[i].topLeft.col < m && m < squares[i].bottomRight.col {
			cd := s.NearestVerticalDistanceTo(squares[i])
			if nd > cd && cd > 0 {
				nd = cd
				ni = i
			}
		}
	}

	if ni < 0 {
		return nil, ni
	} else {
		return squares[ni], ni
	}
}

// ================================= Predictor =================================

type Predictor interface {
	inputWidth() int
	inputHeight() int
	Predicts(ImageMatrixs) []string
}

// ================================= Nearest Neighbor Predictor =================================
type NNPredictor struct {
	model *Model
}

func NewNNPredictor(model *Model) *NNPredictor {
	return &NNPredictor{
		model: model,
	}
}

func NewNNPredictorFromFile(path string) *NNPredictor {
	model, err := ReadModel(path)
	if err != nil {
		panic(err)
	}

	return &NNPredictor{
		model: &model,
	}
}

func (p *NNPredictor) inputHeight() int {
	r, _ := p.model.ModelImages[0].Data.Dims()
	return r
}

func (p *NNPredictor) inputWidth() int {
	_, c := p.model.ModelImages[0].Data.Dims()
	return c
}

func (p *NNPredictor) Predicts(images ImageMatrixs) []string {
	predictedLabels := make([]string, len(images))
	mr, mc := p.model.ModelImages[0].Data.Dims()

	for i, image := range images {
		resizedMatrix := PadAndResize(image, mr, mc)

		min := math.MaxFloat64
		for _, modelImage := range p.model.ModelImages {
			distance := EuclideanDistance(resizedMatrix, modelImage.Data)

			fmt.Println(modelImage.Label, distance, " ")
			if min > distance {
				min = distance
				predictedLabels[i] = modelImage.Label
			}
		}
	}

	return predictedLabels
}

// ================================= Tensor CNN Predictor =================================

type CNNPredictor struct {
	graph       *tf.Graph
	labels      []string
	InputHeight int
	InputWidth  int
}

func NewCNNPredictor(graph *tf.Graph, labels []string) *CNNPredictor {
	return &CNNPredictor{
		graph:  graph,
		labels: labels,
	}
}

func NewCNNPredictorFromDir(dir string) *CNNPredictor {
	graph, labels := loadModel(dir)
	return NewCNNPredictor(graph, labels)
}

func (p *CNNPredictor) inputHeight() int {
	return p.InputHeight
}

func (p *CNNPredictor) inputWidth() int {
	return p.InputWidth
}

func (p *CNNPredictor) Predicts(images ImageMatrixs) []string {

	session, err := tf.NewSession(p.graph, nil)
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
			p.graph.Operation("convolution2d_input_1").Output(0): tensorImages,
			p.graph.Operation("keras_learning_phase").Output(0):  kerasFlag,
		},
		[]tf.Output{
			p.graph.Operation("ArgMax_1").Output(0),
		},
		nil)

	if err != nil {
		log.Fatal(err)
	}

	predictions := output[0].Value().([]int64)
	result := make([]string, len(predictions))
	for i, prediction := range predictions {
		result[i] = p.labels[prediction]
	}

	return result
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

func ScanToStrings(p Predictor, image image.Image) []string {
	im := ImageToBinaryArray(image)
	squaress, charss := CirucularScan(im)
	results := []string{}

	for k, chars := range charss {
		datas := make([]ImageMatrix, len(chars))
		for i := 0; i < len(chars); i++ {
			datas[i] = PadAndResize(chars[i], p.inputHeight(), p.inputWidth())
		}

		texts := p.Predicts(datas)
		result := ""
		for i, text := range texts {
			if i < len(squaress[k])-1 {
				if squaress[k][i].NearestHorizontalDistanceTo(squaress[k][i+1]) > float64(squaress[k][i].Width()) {
					result += text + " "
					continue
				}
			}

			result += text
		}

		results = append(results, result)
	}

	return results
}
