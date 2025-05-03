package service

import (
	"fmt"
	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

func makeGraphicsForSVR(datapoints []models.DataPoint, Ypred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)
	yPred := denseToSlice(Ypred)

	for i := 0; i < len(datapoints[0].Variables); i++ {
		xyPlot := make(map[string]float64)
		for j := 0; j < len(datapoints); j++ {
			if i >= len(datapoints[j].Variables) {
				continue // избегаем выхода за границы
			}
			x := datapoints[j].Variables[i]
			xyPlot[fmt.Sprintf("%f", yPred[j])] = x
		}
		res[i] = xyPlot
	}

	return res
}
func makeGraphicsForPoly(datapoints []models.DataPoint, CoeffMatrics []blas64.General, Ypred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)

	_ = denseToSlice(Ypred)
	coeffMat := mat.NewDense(CoeffMatrics[0].Rows, CoeffMatrics[0].Cols, CoeffMatrics[0].Data)
	coeffs := denseToSlice(coeffMat)

	for i := 0; i < len(datapoints[0].Variables); i++ {
		xyPlot := make(map[string]float64)
		for j := 0; j < len(datapoints); j++ {
			if i >= len(datapoints[j].Variables) {
				continue // избегаем выхода за границы
			}
			x := datapoints[j].Variables[i]
			y := x * coeffs[i]
			xyPlot[fmt.Sprintf("%f", y)] = x
		}
		res[i] = xyPlot
	}
	return res
}

// func makeGraphics(r *linearmodel.Regression) map[int]map[string]float64 {
// 	res := make(map[int]map[string]float64)

// 	_ = r.GetCoeffs()
// 	datapoints := r.GetDataPoints()
// 	for i := 1; i < len(datapoints[0].Variables); i++ {
// 		xyPlot := make(map[string]float64)
// 		for j := 0; j < len(datapoints); j++ {
// 			x := datapoints[j].Variables[i-1]
// 			y := datapoints[j].Predicted
// 			xyPlot[fmt.Sprintf("%f", y)] = x
// 		}
// 		res[i] = xyPlot
// 	}

// 	return res
// }

func makeGraphicsFoBlas64(datapoints []models.DataPoint, matrixCoeff blas64.General, YPred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)

	_ = denseToSlice(YPred)
	coeffMat := mat.NewDense(matrixCoeff.Rows, matrixCoeff.Cols, matrixCoeff.Data)
	coeffs := denseToSlice(coeffMat)

	for i := 0; i < len(datapoints[0].Variables); i++ {
		xyPlot := make(map[string]float64)
		for j := 0; j < len(datapoints); j++ {
			if i >= len(datapoints[j].Variables) {
				continue // избегаем выхода за границы
			}
			x := datapoints[j].Variables[i]
			y := x * coeffs[i]
			xyPlot[fmt.Sprintf("%f", y)] = x
		}
		res[i] = xyPlot
	}
	return res
}
func makeGraphicsForOtherMethod(datapoints []models.DataPoint, matrixCoeff *mat.Dense, YPred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)
	_ = denseToSlice(YPred)
	coeffMat := mat.NewDense(matrixCoeff.RawMatrix().Rows, matrixCoeff.RawMatrix().Cols, matrixCoeff.RawMatrix().Data)
	coeffs := denseToSlice(coeffMat)
	for i := 0; i < len(datapoints[0].Variables); i++ {
		xyPlot := make(map[string]float64)
		for j := 0; j < len(datapoints); j++ {
			if i >= len(datapoints[j].Variables) {
				continue // избегаем выхода за границы
			}
			x := datapoints[j].Variables[i]
			y := x * coeffs[i]
			xyPlot[fmt.Sprintf("%f", y)] = x
		}
		res[i] = xyPlot
	}
	// for i := 0; i < len(datapoints[0].Variables); i++ {
	// 	var (
	// 		xMin, xMax float64
	// 		yMin, yMax float64
	// 		first      = true
	// 	)

	// 	for j := 0; j < len(datapoints); j++ {
	// 		if i >= len(datapoints[j].Variables) {
	// 			continue
	// 		}
	// 		x := datapoints[j].Variables[i]
	// 		y := x * coeffs[i]

	// 		if first {
	// 			xMin, xMax = x, x
	// 			yMin, yMax = y, y
	// 			first = false
	// 		} else {
	// 			if x < xMin {
	// 				xMin = x
	// 				yMin = y
	// 			}
	// 			if x > xMax {
	// 				xMax = x
	// 				yMax = y
	// 			}
	// 		}
	// 	}

	// 	// Сохраняем только две точки: начальную и конечную
	// 	xyPlot := make(map[string]float64)
	// 	xyPlot[fmt.Sprintf("%f", yMin)] = xMin
	// 	xyPlot[fmt.Sprintf("%f", yMax)] = xMax

	// 	res[i] = xyPlot
	// }

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
