package gocr

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type Coordinate struct {
	row int
	col int
}

func NewCoordinate(r, c int) *Coordinate {
	return &Coordinate{
		row: r,
		col: c,
	}
}

func (c *Coordinate) N() *Coordinate {
	return &Coordinate{
		row: c.row - 1,
		col: c.col,
	}
}

func (c *Coordinate) E() *Coordinate {
	return &Coordinate{
		row: c.row,
		col: c.col + 1,
	}
}

func (c *Coordinate) S() *Coordinate {
	return &Coordinate{
		row: c.row + 1,
		col: c.col,
	}
}

func (c *Coordinate) W() *Coordinate {
	return &Coordinate{
		row: c.row,
		col: c.col - 1,
	}
}

func (c *Coordinate) NE() *Coordinate {
	return &Coordinate{
		row: c.row + 1,
		col: c.col + 1,
	}
}

func (c *Coordinate) SE() *Coordinate {
	return &Coordinate{
		row: c.row - 1,
		col: c.col + 1,
	}
}

func (c *Coordinate) SW() *Coordinate {
	return &Coordinate{
		row: c.row - 1,
		col: c.col - 1,
	}
}

func (c *Coordinate) NW() *Coordinate {
	return &Coordinate{
		row: c.row + 1,
		col: c.col - 1,
	}
}

func (c *Coordinate) IsInside(s *Square) bool {
	if s.Area() == 0 {
		return false
	}

	return c.row >= s.topLeft.row && c.col >= s.topLeft.col && c.row < s.bottomRight.row && c.col < s.bottomRight.col
}

func (c *Coordinate) DistanceTo(c2 *Coordinate) float64 {
	return math.Sqrt(math.Pow(float64(c2.row-c.row), 2) + math.Pow(float64(c2.col-c.col), 2))
}

func (c *Coordinate) HorizontalDistanceTo(c2 *Coordinate) float64 {
	dist := float64(c.col - c2.col)

	if dist < 0 {
		return -dist
	} else {
		return dist
	}
}

func (c *Coordinate) VerticalDistanceTo(c2 *Coordinate) float64 {
	dist := float64(c.row - c2.row)

	if dist < 0 {
		return -dist
	} else {
		return dist
	}
}

type Square struct {
	topLeft     *Coordinate
	bottomRight *Coordinate
}

func NewSquare(coor1, coor2 *Coordinate) *Square {

	tlr, tlc, brr, brc := 0, 0, 0, 0

	if coor1.row > coor2.row {
		tlr = coor2.row
		brr = coor1.row
	} else {
		tlr = coor1.row
		brr = coor2.row
	}

	if coor1.col > coor2.col {
		brc = coor1.col
		tlc = coor2.col
	} else {
		brc = coor2.col
		tlc = coor1.col
	}

	return &Square{
		topLeft:     NewCoordinate(tlr, tlc),
		bottomRight: NewCoordinate(brr, brc),
	}
}

func (s *Square) Width() int {
	return s.bottomRight.col - s.topLeft.col
}

func (s *Square) Height() int {
	return s.bottomRight.row - s.topLeft.row
}

func (s *Square) Area() int {
	return (s.bottomRight.row - s.topLeft.row) * (s.bottomRight.col - s.topLeft.col)
}

func (s *Square) Expand(c *Coordinate) {
	if c.row > s.bottomRight.row {
		s.bottomRight.row = c.row
	} else if c.row < s.topLeft.row {
		s.topLeft.row = c.row
	}

	if c.col > s.bottomRight.col {
		s.bottomRight.col = c.col
	} else if c.col < s.topLeft.col {
		s.topLeft.col = c.col
	}
}

func (s *Square) Include(r, c int) bool {
	return r >= s.topLeft.row && c >= s.topLeft.col && r <= s.bottomRight.row && c <= s.bottomRight.col
}

func (s *Square) DistancesTo(s2 *Square) []float64 {
	return []float64{
		s.topLeft.DistanceTo(s2.topLeft),
		s.topLeft.DistanceTo(s2.bottomRight),
		s.bottomRight.DistanceTo(s2.topLeft),
		s.bottomRight.DistanceTo(s2.bottomRight),
	}
}

func (s *Square) HorizontalDistancesTo(s2 *Square) []float64 {
	return []float64{
		s.topLeft.HorizontalDistanceTo(s2.topLeft),
		s.topLeft.HorizontalDistanceTo(s2.bottomRight),
		s.bottomRight.HorizontalDistanceTo(s2.topLeft),
		s.bottomRight.HorizontalDistanceTo(s2.bottomRight),
	}
}

func (s *Square) VerticalDistancesTo(s2 *Square) []float64 {
	return []float64{
		s.topLeft.VerticalDistanceTo(s2.topLeft),
		s.topLeft.VerticalDistanceTo(s2.bottomRight),
		s.bottomRight.VerticalDistanceTo(s2.topLeft),
		s.bottomRight.VerticalDistanceTo(s2.bottomRight),
	}
}

func (s *Square) NearestDistanceTo(s2 *Square) float64 {
	ds := s.DistancesTo(s2)
	m := ds[0]

	for _, e := range ds {
		if e < m {
			m = e
		}
	}

	return m
}

func (s *Square) AverageDistanceTo(s2 *Square) float64 {
	ds := s.DistancesTo(s2)
	sum := 0.0

	for _, e := range ds {
		sum += e
	}

	return sum / 4
}

func (s *Square) FarthestDistanceTo(s2 *Square) float64 {
	ds := s.DistancesTo(s2)
	m := ds[0]

	for _, e := range ds {
		if e > m {
			m = e
		}
	}

	return m
}

func (s *Square) NearestHorizontalDistanceTo(s2 *Square) float64 {
	ds := s.HorizontalDistancesTo(s2)
	m := ds[0]

	for _, e := range ds {
		if e < m {
			m = e
		}
	}

	return m
}

func (s *Square) AverageHorizontalDistanceTo(s2 *Square) float64 {
	ds := s.HorizontalDistancesTo(s2)
	sum := 0.0

	for _, e := range ds {
		sum += e
	}

	return sum / 4
}

func (s *Square) FarthestHorizontalDistanceTo(s2 *Square) float64 {
	ds := s.HorizontalDistancesTo(s2)
	m := ds[0]

	for _, e := range ds {
		if e > m {
			m = e
		}
	}

	return m
}

func (s *Square) NearestVerticalDistanceTo(s2 *Square) float64 {
	ds := s.VerticalDistancesTo(s2)
	m := ds[0]

	for _, e := range ds {
		if e < m {
			m = e
		}
	}

	return m
}

func (s *Square) AverageVerticalDistanceTo(s2 *Square) float64 {
	ds := s.VerticalDistancesTo(s2)
	sum := 0.0

	for _, e := range ds {
		sum += e
	}

	return sum / 4
}

func (s *Square) FarthestVerticalDistanceTo(s2 *Square) float64 {
	ds := s.VerticalDistancesTo(s2)
	m := ds[0]

	for _, e := range ds {
		if e > m {
			m = e
		}
	}

	return m
}

func (s *Square) Merge(a2 *Square) {
	s.Expand(a2.topLeft)
	s.Expand(a2.bottomRight)
}

type ImageVector []uint8

type ImageMatrix [][]uint8

func NewImageMatrix(r, c int) ImageMatrix {
	imageArray := make([][]uint8, r)
	for i := 0; i < r; i++ {
		imageArray[i] = make([]uint8, c)
	}

	return imageArray
}

func NewImageMatrixWithDefaultValue(r, c int, v uint8) ImageMatrix {
	imageArray := make([][]uint8, r)
	for i := 0; i < r; i++ {
		imageArray[i] = make([]uint8, c)
		for j := 0; j < c; j++ {
			imageArray[i][j] = v
		}
	}

	return imageArray
}

func (i ImageMatrix) Dims() (int, int) {
	if len(i) == 0 {
		return 0, 0
	}

	return len(i), len(i[0])
}

func (i ImageMatrix) At(r, c int) uint8 {
	return i[r][c]
}

func (i ImageMatrix) AtCoordinate(c *Coordinate) uint8 {
	return i[c.row][c.col]
}

func (i ImageMatrix) Set(r, c int, value uint8) {
	i[r][c] = value
}

func (i ImageMatrix) Slice(sr, er, sc, ec int) ImageMatrix {
	slice := NewImageMatrix(er-sr, ec-sc)

	for r := sr; r < er; r++ {
		for c := sc; c < ec; c++ {
			slice[r-sr][c-sc] = i[r][c]
		}
	}

	return slice
}

func (im ImageMatrix) SliceSquare(s *Square) ImageMatrix {
	r, c := s.bottomRight.row-s.topLeft.row, s.bottomRight.col-s.topLeft.col
	slice := NewImageMatrix(r, c)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			slice[i][j] = im[s.topLeft.row+i][s.topLeft.col+j]
		}
	}

	return slice
}

func (i ImageMatrix) NNInterpolation(tr, tc int) ImageMatrix {
	r, c := i.Dims()
	rRatio := float64(r) / float64(tr)
	cRatio := float64(c) / float64(tc)
	output := NewImageMatrix(tr, tc)

	for r := 0; r < tr; r++ {
		for c := 0; c < tc; c++ {
			pr := int(math.Floor(float64(r) * rRatio))
			pc := int(math.Floor(float64(c) * cRatio))

			output.Set(r, c, i.At(pr, pc))
		}
	}

	return output
}

// Erode using 4 Neighborhood
func (im ImageMatrix) Erode() {
	r, c := im.Dims()
	a := NewSquare(NewCoordinate(0, 0), NewCoordinate(r, c))

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			n := NewCoordinate(i, j)
			n1, n2, n3, n4 := n.N(), n.E(), n.W(), n.S()

			if n1.IsInside(a) && n2.IsInside(a) && n3.IsInside(a) && n4.IsInside(a) {
				sum := im.AtCoordinate(n) + im.AtCoordinate(n1) + im.AtCoordinate(n2) + im.AtCoordinate(n3) + im.AtCoordinate(n4)
				if sum > 0 {
					im.Set(n.row, n.col, 2)
				}
			}
		}
	}

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if im.At(i, j) == 2 {
				im.Set(i, j, 0)
			}
		}
	}

}

func (i ImageMatrix) Pad(top, bottom, left, right int, value uint8) ImageMatrix {
	sr, sc := i.Dims()
	nr := sr + top + bottom
	nc := sc + left + right
	output := NewImageMatrixWithDefaultValue(nr, nc, value)

	for r := 0; r < sr; r++ {
		for c := 0; c < sc; c++ {
			output[r+top][c+left] = i[r][c]
		}
	}

	return output
}

func (i ImageMatrix) Row(r int) ImageVector {
	return i[r]
}

func (i ImageMatrix) Col(c int) ImageVector {
	v := NewImageVector(len(i))

	for j := 0; j < len(i); j++ {
		v[j] = i[j][c]
	}

	return v
}

func (im ImageMatrix) Historgram() []int {
	hist := make([]int, 256)
	r, c := im.Dims()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			hist[im.At(i, j)] += 1
		}
	}

	return hist
}

type ImageMatrixs []ImageMatrix

func (is ImageMatrixs) Average() ImageMatrix {
	l := len(is)
	sr, sc := is[0].Dims()
	output := NewImageMatrix(sr, sc)

	for r := 0; r < sr; r++ {
		for c := 0; c < sc; c++ {
			sum := 0
			for i := 0; i < l; i++ {
				sum += int(is[i][r][c])
			}

			output[r][c] = uint8(sum / l)
		}
	}

	return output
}

func NewImageVector(l int) ImageVector {
	return make([]uint8, l)
}

func (v ImageVector) Sum() uint64 {
	var sum uint64 = 0

	for i := 0; i < len(v); i++ {
		sum += uint64(v[i])
	}

	return sum
}

// ========================= Convolution =========================

type fn func(int, bool) string

// Calculate size after forward
// ds: Data Size
// ks: Kernel / Filter size
// p: Padding
// s: Stride
func sizeAfter(ds, ks, p, s int) int {
	return ((ds - ks + 2*p) / s) + 1
}

func Sigmoid(x float64, deriv bool) float64 {
	if deriv {
		return x * (1 - x)
	}

	return 1 / (1 + math.Exp(-x))
}

func Relu(x float64, deriv bool) float64 {
	if deriv {
		if x > 0 {
			return 1
		} else {
			return 0
		}
	} else {
		if x > 0 {
			return x
		} else {
			return 0
		}
	}
}

func LeakyRelu(x float64, deriv bool) float64 {
	a := 0.001

	if deriv {
		if x > 0 {
			return 1
		} else {
			return a
		}
	} else {
		if x > 0 {
			return x
		} else {
			return a * x
		}
	}
}

func Reshape2DTo2D(x *mat64.Dense, r, c int) *mat64.Dense {
	_, sc := x.Dims()
	o := mat64.NewDense(r, c, nil)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			o.Set(i, j, x.At(i/c, (i*c+j)%sc))
		}
	}

	return o
}

func Im2col(x []*mat64.Dense, p, s, kd, kr, kc int) *mat64.Dense {
	dd := len(x)
	_, dc := x[0].Dims()
	dnr := kd * kr * kc
	pl := sizeAfter(dc, kc, p, s)
	dnc := dd * pl

	o := mat64.NewDense(dnr, dnc, nil)

	for i := 0; i < dnr; i++ {
		for j := 0; j < dnc; j++ {
			v := x[i/(dnr/dd)].At(i/kc%dd+(j/pl)*s, i%kc+(j%pl)*s)
			o.Set(i, j, v)
		}
	}

	return o
}

func Col2im(x *mat64.Dense, d, r, c int) []*mat64.Dense {
	_, sc := x.Dims()
	o := make([]*mat64.Dense, d)

	for i := 0; i < d; i++ {
		o[i] = mat64.NewDense(r, c, nil)
		for j := 0; j < r; j++ {
			for k := 0; k < c; k++ {
				o[i].Set(j, k, x.At(i, (j*c+k)%sc))
			}
		}
	}

	return o
}

func ReshapeKernel(x [][]*mat64.Dense) *mat64.Dense {
	n := len(x)
	d := len(x[0])
	r, c := x[0][0].Dims()
	nr := n
	nc := d * r * c

	o := mat64.NewDense(nr, nc, nil)

	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			v := x[i][j/(r*c)].At(j/c%r, j%c)
			o.Set(i, j, v)
		}
	}

	return o
}

func Convolve(x, k [][]*mat64.Dense, p, s int) [][]*mat64.Dense {
	kn := len(k)
	kd := len(k[0])
	kr, kc := k[0][0].Dims()
	xn := len(x)
	dr, dc := x[0][0].Dims()

	x_new := make([]*mat64.Dense, xn)
	k_new := ReshapeKernel(k)
	result := make([][]*mat64.Dense, xn)

	pr := sizeAfter(dr, kr, p, s)
	pc := sizeAfter(dc, kc, p, s)

	for i := 0; i < xn; i++ {
		x_new[i] = Im2col(x[i], p, s, kd, kr, kc)
	}

	for i := 0; i < xn; i++ {
		r := &mat64.Dense{}
		r.Mul(k_new, x_new[i])
		result[i] = Col2im(r, kn, pr, pc)
	}

	return result
}

func MaxPool(x *mat64.Dense, fr, fc int) (*mat64.Dense, *mat64.Dense) {
	r, c := x.Dims()
	frc := fr * fc
	nr := sizeAfter(r, fr, 0, fr)
	nc := sizeAfter(c, fc, 0, fc)
	o := mat64.NewDense(nr, nc, nil)
	sw := mat64.NewDense(nr, nc, nil)

	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			m := 0.0
			mi := 0

			for k := 0; k < frc; k++ {
				v := x.At(i*fr+k/fr, j*fc+k%fc)

				if m < v {
					m = v
					mi = k
				}
			}

			o.Set(i, j, m)
			sw.Set(i, j, float64(mi))
		}
	}

	return o, sw
}
