package neural_network

type Matrix [][]float32
type Vector []float32

func NewVectorFromSlice(slice []float32) (v Vector) {
	return Vector(slice)
}

func NewVector(size int) (v Vector) {
	return make(Vector, size)
}

func NewMatrix(numRows, numColumns int) (m Matrix) {
	m = make(Matrix, numRows)
	for i := range m {
		m[i] = make([]float32, numColumns)
	}
	return m
}

func ComputeMatrixVectorMult(m Matrix, v Vector) (res Vector) {
	res = make([]float32, len(m))
	for i := range m {
		for j := range m[i] {
			res[i] += m[i][j] * v[j]
		}
	}

	return res
}

func ComputeVectorAdd(v1, v2 Vector) (res Vector) {
	res = make([]float32, len(v1))
	for i := range v1 {
		res[i] = v1[i] + v2[i]
	}
	return res
}

func VectorizeFunction(f func(float32) float32, v Vector) (res Vector) {
	res = make([]float32, len(v))
	for i := range v {
		res[i] = f(v[i])
	}
	return res
}
