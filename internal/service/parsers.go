package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

func flattenMatrix(m *mat.Dense) []float64 {
	r, c := m.Dims()
	values := make([]float64, 0, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			values = append(values, m.At(i, j))
		}
	}
	return values
}

func flattenBlasMatrix(m blas64.General) []float64 {
	flattened := make([]float64, 0, m.Rows*m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			flattened = append(flattened, m.Data[i*m.Stride+j])
		}
	}
	return flattened
}

func flattenNested(coefs []blas64.General) []float64 {
	var flat []float64

	for _, mat := range coefs {
		for r := 0; r < mat.Rows; r++ {
			for c := 0; c < mat.Cols; c++ {
				flat = append(flat, mat.Data[r*mat.Stride+c])
			}
		}
	}

	return flat
}

type float = float64

// PrepareTrainTestMatrices - parse data form models.DataPoint into mat.Dense
func PrepareTrainTestMatrices(dataPoints []models.DataPoint, trainRatio float64) (
	Xtrain, Ytrain, Xtest, Ytest *mat.Dense,
	testPoints []models.DataPoint,
	err error,
) {
	numOfSamples := len(dataPoints)
	if numOfSamples == 0 {
		return nil, nil, nil, nil, nil, fmt.Errorf("no data points provided")
	}
	numOfVars := len(dataPoints[0].Variables)

	// Перемешиваем и делим на train/test
	shuffled := make([]models.DataPoint, numOfSamples)
	copy(shuffled, dataPoints)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numOfSamples, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	splitIdx := int(trainRatio * float64(numOfSamples))
	trainPoints := shuffled[:splitIdx]
	testPoints = shuffled[splitIdx:]

	Xtrain = mat.NewDense(len(trainPoints), numOfVars, nil)
	Ytrain = mat.NewDense(len(trainPoints), 1, nil)
	for i, dp := range trainPoints {
		for j := 0; j < numOfVars; j++ {
			Xtrain.Set(i, j, dp.Variables[j])
		}
		Ytrain.Set(i, 0, dp.Observed)
	}

	Xtest = mat.NewDense(len(testPoints), numOfVars, nil)
	Ytest = mat.NewDense(len(testPoints), 1, nil)

	var wg sync.WaitGroup
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}
	rowChan := make(chan rowData, len(testPoints))

	for i := 0; i < len(testPoints); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dp := testPoints[i]
			rowChan <- rowData{
				index:    i,
				varRow:   append([]float64(nil), dp.Variables...),
				obsValue: dp.Observed,
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(rowChan)
	}()

	for row := range rowChan {
		for j := 0; j < numOfVars; j++ {
			Xtest.Set(row.index, j, row.varRow[j])
		}
		Ytest.Set(row.index, 0, row.obsValue)
	}

	return Xtrain, Ytrain, Xtest, Ytest, testPoints, nil
}
