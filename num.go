package gocr

import "math"

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
