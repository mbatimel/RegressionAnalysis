package service

import (
	"fmt"
	"sync"

	"math"

	"github.com/mbatimel/RegressionAnalysis/internal/base"
	linearmodel "github.com/mbatimel/RegressionAnalysis/internal/linear_model"
	"github.com/mbatimel/RegressionAnalysis/internal/metrics"
	neuralnetwork "github.com/mbatimel/RegressionAnalysis/internal/neural_network"
	"github.com/mbatimel/RegressionAnalysis/internal/preprocessing"
	"github.com/mbatimel/RegressionAnalysis/internal/svm"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
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
func makeGraphics(r *linearmodel.Regression) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)

	coeff := r.GetCoeffs()
	datapoints := r.GetDataPoints()
	for i := 1; i < len(coeff); i++ {
		xyPlot := make(map[string]float64)
		for j := 0; j < len(datapoints); j++ {
			x := datapoints[j].Variables[i-1]
			y := datapoints[j].Predicted
			xyPlot[fmt.Sprintf("%f", y)] = x
		}
		res[i] = xyPlot
	}

	return res
}
func makeGraphicsFoBlas64(datapoints []models.DataPoint, matrixCoeff blas64.General, YPred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)

	// Конвертируем blas64.General в *mat.Dense для удобства работы
	coeffMat := mat.NewDense(matrixCoeff.Rows, matrixCoeff.Cols, matrixCoeff.Data)
	coeff := denseToSlice(coeffMat)
	yPred := denseToSlice(YPred)

	for i := 0; i < len(coeff); i++ {
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
func makeGraphicsForOtherMethod(datapoints []models.DataPoint, matrixCoeff *mat.Dense, YPred *mat.Dense) map[int]map[string]float64 {
	res := make(map[int]map[string]float64)

	coeff := denseToSlice(matrixCoeff)
	yPred := denseToSlice(YPred)

	for i := 0; i < len(coeff); i++ {
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
func ridgeChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Канал для сбора строк переменных и наблюдений
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}

	rowChan := make(chan rowData, numOfSamples)
	var wg sync.WaitGroup

	// Параллельно подготавливаем строки
	for i := 0; i < numOfSamples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dp := dataPoints[i]
			varRow := make([]float64, numOfVars)
			copy(varRow, dp.Variables)
			rowChan <- rowData{
				index:    i,
				varRow:   varRow,
				obsValue: dp.Observed,
			}
		}(i)
	}

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(rowChan)
	}()

	// Последовательно записываем в матрицы
	for row := range rowChan {
		for j := 0; j < numOfVars; j++ {
			variables.Set(row.index, j, row.varRow[j])
		}
		observed.Set(row.index, 0, row.obsValue)
	}

	// Создаем модель Ridge Regression
	regr := linearmodel.NewRidge()
	regr.Alpha = 1
	regr.Tol = 0.01
	regr.Normalize = true
	regr.L1Ratio = 10

	// Обучаем модель
	regr.Fit(variables, observed)

	// Делаем предсказание
	Ypred := mat.NewDense(numOfSamples, 1, nil)
	regr.Predict(variables, Ypred)

	// Возвращаем результаты
	res := map[string]interface{}{
		"ridge Ypred":              fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"ridge Coef":               fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Coef)),
		"ridge XOffsetoef":         fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XOffset)),
		"ridge XScale":             fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XScale)),
		"ridge Intercept":          fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Intercept)),
		"ridge ActivationFunction": regr.ActivationFunction,
		"graphics":                 makeGraphicsForOtherMethod(dataPoints, regr.Coef, Ypred),
	}

	return res, nil
}
func lassoChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Канал для сбора строк переменных и наблюдений
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}

	rowChan := make(chan rowData, numOfSamples)
	var wg sync.WaitGroup

	// Параллельно подготавливаем строки
	for i := 0; i < numOfSamples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dp := dataPoints[i]
			varRow := make([]float64, numOfVars)
			copy(varRow, dp.Variables)
			rowChan <- rowData{
				index:    i,
				varRow:   varRow,
				obsValue: dp.Observed,
			}
		}(i)
	}

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(rowChan)
	}()

	// Последовательно записываем в матрицы
	for row := range rowChan {
		for j := 0; j < numOfVars; j++ {
			variables.Set(row.index, j, row.varRow[j])
		}
		observed.Set(row.index, 0, row.obsValue)
	}

	// Создаем модель Lasso
	regr := linearmodel.NewLasso() // Предположим, что у вас есть Lasso модель
	regr.FitIntercept = true
	regr.Normalize = true
	regr.Alpha = 1e-5
	regr.L1Ratio = 1
	regr.MaxIter = 1e5
	regr.Tol = 1e-4
	// Обучаем модель
	regr.Fit(variables, observed)


	Ypred := mat.NewDense(numOfSamples, 1, nil)
	regr.Predict(variables, Ypred)
	rss := &mat.VecDense{}
	rss.SubVec(Ypred.ColView(0), observed.ColView(0))
	rss.MulElemVec(rss, rss)
	res := map[string]interface{}{
		"lasso Ypred":      fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"lasso Coef":       fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Coef.T())),
		"lasso XOffsetoef": fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XOffset.T())),
		"lasso XScale":     fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XScale.T())),
		"lasso Intercept":  fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Intercept.T())),
		"lasso RSS":        mat.Sum(rss),
		"lasso MaxIter":    regr.MaxIter,
		"lasso Tol":        regr.Tol,
		"lasso Alpha":      regr.Alpha,
		"lasso L1Ratio":    regr.L1Ratio,
		"lasso Selection":  regr.Selection,
		"lasso WarmStart":  regr.WarmStart,
		"lasso Positive":   regr.Positive,
		"lasso CDResult":   regr.CDResult,
		"graphics":         makeGraphicsForOtherMethod(dataPoints, regr.Coef, Ypred),
	}
	return res, nil
}
func elasticChecking(dataPoints []models.DataPoint, l1Ratio float64) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Канал для сбора строк переменных и наблюдений
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}

	rowChan := make(chan rowData, numOfSamples)
	var wg sync.WaitGroup

	// Параллельно подготавливаем строки
	for i := 0; i < numOfSamples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dp := dataPoints[i]
			varRow := make([]float64, numOfVars)
			copy(varRow, dp.Variables)
			rowChan <- rowData{
				index:    i,
				varRow:   varRow,
				obsValue: dp.Observed,
			}
		}(i)
	}

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(rowChan)
	}()

	// Последовательно записываем в матрицы
	for row := range rowChan {
		for j := 0; j < numOfVars; j++ {
			variables.Set(row.index, j, row.varRow[j])
		}
		observed.Set(row.index, 0, row.obsValue)
	}
	// Создаем модель ElasticNet
	enet := linearmodel.NewElasticNet()
	enet.Alpha = 1
	enet.Tol = 0.01
	enet.L1Ratio = l1Ratio

	// Обучаем модель
	enet.Fit(variables, observed)

	// Делаем предсказание
	Ypred := mat.NewDense(numOfSamples, 1, nil)
	enet.Predict(variables, Ypred)

	// Возвращаем результаты
	res := map[string]interface{}{
		"elastic Ypred":      fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"elastic Coef":       fmt.Sprintf("%.7f\n", mat.Formatted(enet.LinearRegression.Coef)),
		"elastic XOffsetoef": fmt.Sprintf("%.2f\n", mat.Formatted(enet.LinearRegression.XOffset)),
		"elastic XScale":     fmt.Sprintf("%.2f\n", mat.Formatted(enet.LinearRegression.XScale)),
		"elastic Intercept":  fmt.Sprintf("%.2f\n", mat.Formatted(enet.LinearRegression.Intercept)),
		"elastic MaxIter":    enet.MaxIter,
		"elastic Tol":        enet.Tol,
		"elastic Alpha":      enet.Alpha,
		"elastic L1Ratio":    enet.L1Ratio,
		"elastic Selection":  enet.Selection,
		"elastic WarmStart":  enet.WarmStart,
		"elastic Positive":   enet.Positive,
		"elastic CDResult":   enet.CDResult,
		"graphics":           makeGraphicsForOtherMethod(dataPoints, enet.Coef, Ypred),
	}

	return res, nil
}
func logisticChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Канал для сбора строк переменных и наблюдений
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}

	rowChan := make(chan rowData, numOfSamples)
	var wg sync.WaitGroup

	// Параллельно подготавливаем строки
	for i := 0; i < numOfSamples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dp := dataPoints[i]
			varRow := make([]float64, numOfVars)
			copy(varRow, dp.Variables)
			rowChan <- rowData{
				index:    i,
				varRow:   varRow,
				obsValue: dp.Observed,
			}
		}(i)
	}

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(rowChan)
	}()

	// Последовательно записываем в матрицы
	for row := range rowChan {
		for j := 0; j < numOfVars; j++ {
			variables.Set(row.index, j, row.varRow[j])
		}
		observed.Set(row.index, 0, row.obsValue)
	}

	// Создаем модель LogisticRegression
	regr := linearmodel.NewLogisticRegression()
	regr.Alpha = 1e-5
	regr.MaxIter = 4

	beforeMinimize := func(problem optimize.Problem, initX []float64) {
		fmt.Println("check minimize")
		// check gradients
		settings := &fd.Settings{Step: 1e-8}
		gradFromModel := make([]float64, len(initX))
		gradFromFD := make([]float64, len(initX))
		problem.Func(initX)
		problem.Grad(gradFromModel, initX)
		fd.Gradient(gradFromFD, problem.Func, initX, settings)
	}
	regr.SetbeforeMinimize(beforeMinimize)
	// we create an instance of our Classifier and fit the data.
	regr.Fit(variables, observed)
	accuracy := regr.Score(variables, observed)
	if accuracy >= 0.833 {
		fmt.Println("ok")
	} else {
		fmt.Printf("Accuracy:%.3f\n", accuracy)
	}
	Ypred := mat.NewDense(numOfSamples, 1, nil)
	regr.Predict(variables, Ypred)

	// Возвращаем результаты
	res := map[string]interface{}{
		"logistic Ypred":     fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"logistic Coef":      regr.Coef,
		"logistic Intercept": regr.Intercept,
		"logistic Tol":       regr.Tol,
		"logistic Alpha":     regr.Alpha,
		"graphics":           makeGraphicsFoBlas64(dataPoints, regr.Coef, Ypred),
	}

	return res, nil
}
func svrChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	var res map[string]interface{}
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Канал для сбора строк переменных и наблюдений
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}

	rowChan := make(chan rowData, numOfSamples)
	var wg sync.WaitGroup

	// Параллельно подготавливаем строки
	for i := 0; i < numOfSamples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dp := dataPoints[i]
			varRow := make([]float64, numOfVars)
			copy(varRow, dp.Variables)
			rowChan <- rowData{
				index:    i,
				varRow:   varRow,
				obsValue: dp.Observed,
			}
		}(i)
	}

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(rowChan)
	}()

	// Последовательно записываем в матрицы
	for row := range rowChan {
		for j := 0; j < numOfVars; j++ {
			variables.Set(row.index, j, row.varRow[j])
		}
		observed.Set(row.index, 0, row.obsValue)
	}

	randomState := base.NewLockedSource(7)
	xscaler := preprocessing.NewMinMaxScaler([]float64{-1, 1})
	yscaler := preprocessing.NewMinMaxScaler([]float64{-1, 1})
	Xsc, _ := xscaler.FitTransform(variables, nil)
	Ysc, _ := yscaler.FitTransform(observed, nil)
	Epsilon := 0.1 * yscaler.Scale.At(0, 0)
	Ypred := map[string]*mat.Dense{}
	// Определяем параметры для каждого типа ядра
	kernelOptions := []struct {
		kernel                  string
		C, gamma, coef0, degree float64
	}{
		{kernel: "rbf", C: 1e3, gamma: 0.1},
		{kernel: "sigmoid", C: 1e3, gamma: 0.1},
		{kernel: "poly", C: 1e3, gamma: 1, coef0: 1, degree: 2},
		{kernel: "linear", C: 1e3},
	}

	// Перебираем все варианты ядер
	for _, opt := range kernelOptions {
		fmt.Println(opt.kernel)
		Ypred[opt.kernel] = &mat.Dense{}

		// Настраиваем SVR
		svr := svm.NewSVR()
		svr.Kernel = opt.kernel
		svr.C = opt.C
		svr.Epsilon = Epsilon
		svr.Gamma = opt.gamma
		svr.Coef0 = opt.coef0
		svr.Degree = opt.degree
		svr.RandomState = randomState
		svr.Tol = math.Sqrt(Epsilon)
		svr.MaxIter = 5

		// Обучаем модель и делаем предсказания
		svr.Fit(Xsc, Ysc)
		svr.Predict(Xsc, Ypred[opt.kernel])

		// Обратное преобразование масштабирования
		Ypred[opt.kernel], _ = yscaler.InverseTransform(Ypred[opt.kernel], nil)

		res = map[string]interface{}{
			"svr YPred " + opt.kernel: fmt.Sprintf("%.2f\n", mat.Formatted(Ypred[opt.kernel])),
			"svr Score " + opt.kernel: svr.Score(Xsc, Ysc),
			"graphics":                makeGraphicsForSVR(dataPoints, Ypred[opt.kernel]),
		}

	}

	return res, nil
}
func polynomialChecking(dataPoints []models.DataPoint, degree int) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)
	nSamples, _ := variables.Dims()

	// Канал для сбора строк переменных и наблюдений
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}

	rowChan := make(chan rowData, numOfSamples)
	var wg sync.WaitGroup

	// Параллельно подготавливаем строки
	for i := 0; i < numOfSamples; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dp := dataPoints[i]
			varRow := make([]float64, numOfVars)
			copy(varRow, dp.Variables)
			rowChan <- rowData{
				index:    i,
				varRow:   varRow,
				obsValue: dp.Observed,
			}
		}(i)
	}

	// Закрытие канала после завершения всех горутин
	go func() {
		wg.Wait()
		close(rowChan)
	}()

	// Последовательно записываем в матрицы
	for row := range rowChan {
		for j := 0; j < numOfVars; j++ {
			variables.Set(row.index, j, row.varRow[j])
		}
		observed.Set(row.index, 0, row.obsValue)
	}

	mlp := neuralnetwork.NewMLPClassifier([]int{}, "logistic", "lbfgs", 1)
	mlp.WarmStart = false
	mlp.MaxIter = 400
	mlp.LearningRateInit = .11
	mlp.BatchSize = 118 //1,2,59,118

	variables, _ = preprocessing.NewStandardScaler().FitTransform(variables, nil)

	poly := preprocessing.NewPolynomialFeatures(degree)
	poly.IncludeBias = true

	poly.Fit(variables, observed)

	Xp, _ := poly.Transform(variables, observed)
	_, nOutputs := observed.Dims()
	Ypred := mat.NewDense(nSamples, nOutputs, nil)
	mlp.NLayers = len(dataPoints)
	mlp.Predict(Xp, Ypred)

	// Возвращаем результаты
	res := map[string]interface{}{
		"poly Ypred":    fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"poly accuracy": metrics.AccuracyScore(observed, Ypred, true, nil),
	}

	return res, nil
}
