package neural_network

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	MaxCPValue      = 20.00
	InputVectorSize = 768
)

type NeuralNetwork struct {
	layers  int // not including input layer
	weights []Matrix
	biases  []Vector
}

func NewNetwork(sizes []int) (nn NeuralNetwork) {
	nn.layers = len(sizes) - 1
	nn.biases = make([]Vector, len(sizes)-1)
	nn.weights = make([]Matrix, len(sizes)-1)

	for i := range nn.biases {
		nn.biases[i] = NewVector(sizes[1:][i])
	}

	for i := range nn.weights {
		nn.weights[i] = NewMatrix(sizes[1:][i], sizes[:len(sizes)-1][i])
	}

	return nn
}

func (nn *NeuralNetwork) Compute(input SparseVector) Vector {
	product := ComputeMatrixSparseVectorMult(nn.weights[0], input)
	sum := ComputeVectorAdd(product, nn.biases[0])
	activation := VectorizeFunction(ReLU, sum)

	for l := 1; l < nn.layers-1; l++ {
		product := ComputeMatrixVectorMult(nn.weights[l], activation)
		sum := ComputeVectorAdd(product, nn.biases[l])
		activation = VectorizeFunction(ReLU, sum)
	}

	product = ComputeMatrixVectorMult(nn.weights[nn.layers-1], activation)
	sum = ComputeVectorAdd(product, nn.biases[nn.layers-1])
	return VectorizeFunction(ModifiedTanh, sum)
}

func (nn *NeuralNetwork) LoadFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for {
		more := scanner.Scan()
		if !more {
			err := scanner.Err()
			if err != nil {
				panic(err)
			}
			break
		}

		layer, err := strconv.ParseInt(strings.TrimRight(scanner.Text(), "\n"), 10, 64)
		if err != nil {
			panic(err)
		}

		scanner.Scan()
		numWeights, err := strconv.ParseInt(strings.TrimRight(scanner.Text(), "\n"), 10, 64)
		if err != nil {
			panic(err)
		}

		scanner.Scan()
		numBiases, err := strconv.ParseInt(strings.TrimRight(scanner.Text(), "\n"), 10, 64)
		if err != nil {
			panic(err)
		}

		weights := []float32{}
		biases := []float32{}
		numColumns := len(nn.weights[layer][0])

		for i := int64(0); i < numWeights; i++ {
			scanner.Scan()
			line := scanner.Text()
			line = strings.TrimRight(line, "\n")
			weight, err := strconv.ParseFloat(line, 32)

			if err != nil {
				panic(err)
			}

			weights = append(weights, float32(weight))
		}

		for i := int64(0); i < numBiases; i++ {
			scanner.Scan()
			line := scanner.Text()
			line = strings.TrimRight(line, "\n")
			bias, err := strconv.ParseFloat(line, 32)

			if err != nil {
				panic(err)
			}

			biases = append(biases, float32(bias))
		}

		nn.biases[layer] = biases

		for row, i := 0, 0; i < len(weights); row, i = row+1, i+numColumns {
			rowOfWeights := weights[i : i+numColumns]
			nn.weights[layer][row] = rowOfWeights
		}
	}
}

func max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

func ReLU(x float32) float32 {
	return max(0, x)
}

func ModifiedTanh(x float32) float32 {
	return float32(math.Tanh(float64(x))) * MaxCPValue
}
