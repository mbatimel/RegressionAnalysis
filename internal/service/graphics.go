package service

import (
	"fmt"
	"math"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

func makeGraphicsForLinear(datapoints []models.DataPoint, coeffs []float64) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)

	numFeatures := len(datapoints[0].Variables)
	numPoints := len(datapoints)

	// Получение коэффициентов
	intercept := coeffs[0]

	// Считаем средние значения признаков
	means := calculateMeans(datapoints)

	// Генерируем все комбинации степеней
	exponents := generateExponentCombinations(numFeatures, 1)

	for i := 0; i < numFeatures; i++ {
		xyPlot := make(map[string]float64)

		// Диапазон значений текущего признака
		xMin, xMax := datapoints[0].Variables[i], datapoints[0].Variables[i]
		for _, dp := range datapoints {
			x := dp.Variables[i]
			if x < xMin {
				xMin = x
			}
			if x > xMax {
				xMax = x
			}
		}
		step := (xMax - xMin) / float64(numPoints-1)

		// Строим точки для графика признака i
		for j := 0; j < numPoints; j++ {
			x := xMin + step*float64(j)

			// Вектор входа: xᵢ меняется, остальные — средние
			input := make([]float64, numFeatures)
			copy(input, means)
			input[i] = x

			// Строим полиномиальные признаки
			var polyFeatures []float64
			for _, powers := range exponents {
				term := 1.0
				for fIdx, power := range powers {
					term *= math.Pow(input[fIdx], float64(power))
				}
				polyFeatures = append(polyFeatures, term)
			}

			// Вычисляем y = intercept + Σ(coeffᵢ * polyFeatureᵢ)
			y := intercept
			for k := 0; k < len(polyFeatures); k++ {
				y += coeffs[k] * polyFeatures[k]
			}

			xyPlot[fmt.Sprintf("%.6f", y)] = x
		}

		res[i] = xyPlot
	}

	return res
}

func makeGraphicsForRidge(dp []models.DataPoint, m *mat.Dense, _ *mat.Dense) map[int]map[string]float64 {
	return makeGraphicsForLinear(dp, denseToSlice(m))
}
func makeGraphicsForLasso(dp []models.DataPoint, m *mat.Dense, _ *mat.Dense) map[int]map[string]float64 {
	return makeGraphicsForLinear(dp, denseToSlice(m))
}
func makeGraphicsForElastic(dp []models.DataPoint, m *mat.Dense, _ *mat.Dense) map[int]map[string]float64 {
	return makeGraphicsForLinear(dp, denseToSlice(m))
}

func makeGraphicsForSVR(datapoints []models.DataPoint, Ypred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)
	yPred := denseToSlice(Ypred)

	for i := 0; i < len(datapoints[0].Variables); i++ {
		xyPlot := make(map[string]float64)
		for j := 0; j < len(datapoints); j++ {
			if i >= len(datapoints[j].Variables) {
				continue
			}
			x := datapoints[j].Variables[i]
			xyPlot[fmt.Sprintf("%.6f", yPred[j])] = x
		}
		res[i] = xyPlot
	}
	return res
}

func makeGraphicsForPoly(datapoints []models.DataPoint, CoeffMatrics []blas64.General, _ *mat.Dense, degree int) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)

	numFeatures := len(datapoints[0].Variables)
	numPoints := 100 // Можно задать побольше для сглаживания графиков

	// Получение коэффициентов
	coeffs := denseToSlice(mat.NewDense(CoeffMatrics[0].Rows, CoeffMatrics[0].Cols, CoeffMatrics[0].Data))
	intercept := coeffs[len(coeffs)-1]

	// Считаем средние значения признаков
	means := calculateMeans(datapoints)

	// Генерируем все комбинации степеней
	exponents := generateExponentCombinations(numFeatures, degree)

	for i := 0; i < numFeatures; i++ {
		xyPlot := make(map[string]float64)

		// Диапазон значений текущего признака
		xMin, xMax := datapoints[0].Variables[i], datapoints[0].Variables[i]
		for _, dp := range datapoints {
			x := dp.Variables[i]
			if x < xMin {
				xMin = x
			}
			if x > xMax {
				xMax = x
			}
		}
		step := (xMax - xMin) / float64(numPoints-1)

		// Строим точки для графика признака i
		for j := 0; j < numPoints; j++ {
			x := xMin + step*float64(j)

			// Вектор входа: xᵢ меняется, остальные — средние
			input := make([]float64, numFeatures)
			copy(input, means)
			input[i] = x

			// Строим полиномиальные признаки
			var polyFeatures []float64
			for _, powers := range exponents {
				term := 1.0
				for fIdx, power := range powers {
					term *= math.Pow(input[fIdx], float64(power))
				}
				polyFeatures = append(polyFeatures, term)
			}

			// Вычисляем y = intercept + Σ(coeffᵢ * polyFeatureᵢ)
			y := intercept
			for k := 0; k < len(polyFeatures); k++ {
				y += coeffs[k] * polyFeatures[k]
			}

			xyPlot[fmt.Sprintf("%.6f", y)] = x
		}

		res[i] = xyPlot
	}

	return res
}

func generateExponentCombinations(numFeatures, degree int) [][]int {
	var result [][]int
	var recurse func(pos int, remDeg int, curr []int)
	recurse = func(pos int, remDeg int, curr []int) {
		if pos == numFeatures {
			if remDeg == 0 {
				comb := make([]int, numFeatures)
				copy(comb, curr)
				result = append(result, comb)
			}
			return
		}
		for i := 0; i <= remDeg; i++ {
			recurse(pos+1, remDeg-i, append(curr, i))
		}
	}
	for d := 1; d <= degree; d++ {
		recurse(0, d, []int{})
	}
	return result
}

func makeGraphicsForLog(testPoints []models.DataPoint, Ypred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)
	yPred := denseToSlice(Ypred)
	numFeatures := len(testPoints[0].Variables)

	for i := 0; i < numFeatures; i++ {
		xyPlot := make(map[string]float64)
		for j := 0; j < len(testPoints); j++ {
			if i >= len(testPoints[j].Variables) {
				continue
			}
			x := testPoints[j].Variables[i]
			y := math.Exp(yPred[j])
			xyPlot[fmt.Sprintf("%.6f", y)] = x
		}
		res[i] = xyPlot
	}
	return res
}

func makeGraphicsForLogistic(datapoints []models.DataPoint, matrixCoeff blas64.General, _ *mat.Dense) map[int]map[string]float64 {
	numFeatures := len(datapoints[0].Variables)
	coeffs := denseToSlice(mat.NewDense(matrixCoeff.Rows, matrixCoeff.Cols, matrixCoeff.Data))
	intercept := coeffs[0]
	means := calculateMeans(datapoints)
	numPoints := 100
	res := make(map[int]map[string]float64)

	for i := 0; i < numFeatures; i++ {
		xyPlot := make(map[string]float64)
		var xMin, xMax float64 = datapoints[0].Variables[i], datapoints[0].Variables[i]
		for _, dp := range datapoints {
			x := dp.Variables[i]
			if x < xMin {
				xMin = x
			}
			if x > xMax {
				xMax = x
			}
		}
		step := (xMax - xMin) / float64(numPoints-1)

		for j := 0; j < numPoints; j++ {
			x := xMin + step*float64(j)
			z := intercept
			for k := 0; k < numFeatures; k++ {
				if k == i {
					z += coeffs[k] * x
				} else {
					z += coeffs[k] * means[k]
				}
			}
			y := 1 / (1 + math.Exp(-z))
			xyPlot[fmt.Sprintf("%.6f", y)] = x
		}
		res[i] = xyPlot
	}
	return res
}

func denseToSlice(d *mat.Dense) []float64 {
	r, c := d.Dims()
	data := make([]float64, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = d.At(i, j)
		}
	}
	return data
}

func calculateMeans(datapoints []models.DataPoint) []float64 {
	numFeatures := len(datapoints[0].Variables)
	means := make([]float64, numFeatures)

	for _, dp := range datapoints {
		for i, v := range dp.Variables {
			means[i] += v
		}
	}
	for i := range means {
		means[i] /= float64(len(datapoints))
	}
	return means
}
