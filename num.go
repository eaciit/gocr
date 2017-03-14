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
