package gocr

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type ImageMatrix [][]uint8

type ImageVector []uint8

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

func (i ImageMatrix) Pad(top, bottom, left, right int, value uint8) ImageMatrix {
	sr, sc := i.Dims()
	nr := sr + top + bottom
	nc := sc + left + right
	output := NewImageMatrixWithDefaultValue(nr, nc, 0)

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

func sizeAfterConvolve(ds, ks, p, s int) int {
	return ((ds - ks + 2*p) / s) + 1
}

func Im2col(x []*mat64.Dense, p, s, kd, kr, kc int) *mat64.Dense {
	dd := len(x)
	_, dc := x[0].Dims()
	dnr := kd * kr * kc
	pl := sizeAfterConvolve(dc, kc, p, s)
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

func Reshape2D(x *mat64.Dense, r, c int) *mat64.Dense {

	_, sc := x.Dims()
	o := mat64.NewDense(r, c, nil)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			o.Set(i, j, x.At(i/c, (i*c+j)%sc))
		}
	}

	return o
}
