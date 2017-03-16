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

func (c *Coordinate) IsInside(a *Area) bool {
	if a.topLeft.row == a.bottomRight.row && a.topLeft.col == a.bottomRight.col {
		return false
	}

	return c.row >= a.topLeft.row && c.col >= a.topLeft.col && c.row <= a.bottomRight.row && c.col <= a.bottomRight.col
}

type Area struct {
	topLeft     *Coordinate
	bottomRight *Coordinate
}

func NewArea(coor1, coor2 *Coordinate) *Area {

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

	return &Area{
		topLeft:     NewCoordinate(tlr, tlc),
		bottomRight: NewCoordinate(brr, brc),
	}
}

func (a *Area) Expand(c *Coordinate) {
	if c.row > a.bottomRight.row {
		a.bottomRight.row = c.row
	} else if c.row < a.topLeft.row {
		a.topLeft.row = c.row
	}

	if c.col > a.bottomRight.col {
		a.bottomRight.col = c.col
	} else if c.col < a.topLeft.col {
		a.topLeft.col = c.col
	}
}

func (a *Area) Include(r, c int) bool {
	return r >= a.topLeft.row && c >= a.topLeft.col && r <= a.bottomRight.row && c <= a.bottomRight.col
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

func (im ImageMatrix) SliceArea(a *Area) ImageMatrix {
	r, c := a.bottomRight.row-a.topLeft.row, a.bottomRight.col-a.topLeft.col
	slice := NewImageMatrix(r, c)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			slice[i][j] = im[a.topLeft.row+i][a.topLeft.col+j]
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
	a := NewArea(NewCoordinate(0, 0), NewCoordinate(r, c))

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
