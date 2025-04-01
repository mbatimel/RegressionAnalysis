package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sync"
	"time"

	"strconv"

	"math"

	"github.com/gofiber/fiber/v2/log"
	"github.com/mbatimel/RegressionAnalysis/internal/base"
	linearmodel "github.com/mbatimel/RegressionAnalysis/internal/linear_model"
	"github.com/mbatimel/RegressionAnalysis/internal/metrics"
	neuralnetwork "github.com/mbatimel/RegressionAnalysis/internal/neural_network"
	"github.com/mbatimel/RegressionAnalysis/internal/preprocessing"
	"github.com/mbatimel/RegressionAnalysis/internal/svm"
	"github.com/xuri/excelize/v2"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	externalApi "github.com/mbatimel/RegressionAnalysis/pkg/interfaces"
	"github.com/rs/zerolog"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

type regressionService struct {
	logger zerolog.Logger
}

func (rs *regressionService) MlrRegression(ctx context.Context, observer string, vars []string, dataPoints []models.DataPoint) (map[string]interface{}, error) {
	r := new(linearmodel.Regression)
	r.SetObserved(observer)
	for i, v := range vars {
		r.SetVar(i, v)
	}
	for _, dp := range dataPoints {
		r.Train(linearmodel.DataPoint(dp.Observed, dp.Variables))
	}
	if err := r.Run(); err != nil {
		return nil, fmt.Errorf("failed to train model: %w", err)
	}
	fmt.Println(r)

	rs.logger.Info().Msg("Starting checking on classifier ridge")
	ridgeCheck, err := ridgeChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier lasso")
	lassoCheck, err := lassoChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier elastic")
	elasticCheck, err := elasticChecking(dataPoints, 1000)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier logistic")
	logisticCheck, err := logisticChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier SVR")
	svrCheck, err := SVRChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier polynomial")
	polynomialCheck, err := polynomialChecking(dataPoints, 3)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}
	res := map[string]interface{}{
		"ridge":      ridgeCheck,
		"lasso":      lassoCheck,
		"elastic":    elasticCheck,
		"logistic":   logisticCheck,
		"SRV":        svrCheck,
		"polynomial": polynomialCheck,
		"data":       r,
		"names":      r.GetNames(),
		"coeff":      r.GetCoeffs(),
		"graphics":   makeGraphics(r),
	}
	return res, nil

}

func (rs *regressionService) MlrRegressionCSV(ctx context.Context, file []byte) (map[string]interface{}, error) {
	r := new(linearmodel.Regression)
	reader := csv.NewReader(bytes.NewReader(file))
	reader.Comma = ';'

	// Читаем заголовки
	header, err := reader.Read()
	if err != nil {
		rs.logger.Println("Ошибка чтения заголовков")
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	rs.logger.Println("Заголовки CSV прочитаны")
	r.SetObserved(header[0])
	for i := 1; i < len(header); i++ {
		r.SetVar(i-1, header[i])
	}

	// Канал для передачи данных
	dataChan := make(chan models.DataPoint, 100)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	// Запускаем обработку строк в отдельных горутинах
	const workerCount = 10
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range dataChan {
				r.Train(linearmodel.DataPoint(record.Observed, record.Variables))
			}
		}()
	}

	// Читаем строки и отправляем их в канал
	go func() {
		defer close(dataChan)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- fmt.Errorf("error reading CSV row: %w", err)
				return
			}

			observed, err := strconv.ParseFloat(record[0], 64)
			if err != nil {
				errChan <- fmt.Errorf("invalid observed value: %w", err)
				return
			}

			variables := make([]float64, len(record)-1)
			for j := 1; j < len(record); j++ {
				variables[j-1], err = strconv.ParseFloat(record[j], 64)
				if err != nil {
					errChan <- fmt.Errorf("invalid variable value: %w", err)
					return
				}
			}

			dataChan <- models.DataPoint{Observed: observed, Variables: variables}
		}
	}()

	// Ждем завершения всех горутин
	wg.Wait()
	close(errChan)

	// Проверяем ошибки
	if err, ok := <-errChan; ok {
		return nil, err
	}

	// Запускаем расчет модели
	if err := r.Run(); err != nil {
		return nil, fmt.Errorf("failed to train model: %w", err)
	}
	res := map[string]interface{}{
		"data":       r,
		"names":      r.GetNames(),
		"coeff":      r.GetCoeffs(),
		"datapoints": r.GetDataPoints(),
		"graphics":   makeGraphics(r),
	}
	return res, nil
}

func (rs *regressionService) MlrRegressionExcel(ctx context.Context, file []byte) (map[string]interface{}, error) {
	r := new(linearmodel.Regression)
	reader := bytes.NewReader(file)
	xlFile, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if err := xlFile.Close(); err != nil {
			rs.logger.Error().Err(fmt.Errorf("failed to close excel: %w", err)).Msg("call")
		}
	}()

	rows, err := xlFile.GetRows("Лист1")
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from Excel: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("not enough data in Excel file")
	}

	// Читаем заголовки
	header := rows[0]
	if len(header) < 2 {
		return nil, fmt.Errorf("invalid header format: at least one independent variable required")
	}

	r.SetObserved(header[0])
	for i := 1; i < len(header); i++ {
		r.SetVar(i-1, header[i])
	}
	rs.logger.Println("Заголовки Excel прочитаны")
	// Канал для передачи данных
	dataChan := make(chan models.DataPoint, 100)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	const workerCount = 10
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range dataChan {
				r.Train(linearmodel.DataPoint(record.Observed, record.Variables))
			}
		}()
	}
	go func() {
		defer close(dataChan)
		for _, row := range rows[1:] {
			if len(row) != len(header) {
				errChan <- fmt.Errorf("data row length does not match header length")
				break
			}
			observed, err := strconv.ParseFloat(row[0], 64)
			if err != nil {
				errChan <- fmt.Errorf("invalid observed value: %w", err)
				break
			}

			variables := make([]float64, len(row)-1)
			for j := 1; j < len(row); j++ {
				variables[j-1], err = strconv.ParseFloat(row[j], 64)
				if err != nil {
					errChan <- fmt.Errorf("invalid variable value: %w", err)
					break
				}
			}

			dataChan <- models.DataPoint{Observed: observed, Variables: variables}
		}
	}()
	// Ждем завершения всех горутин
	wg.Wait()
	close(errChan)

	// Проверяем ошибки
	if err, ok := <-errChan; ok {
		return nil, err
	}

	// Запускаем расчет модели
	if err := r.Run(); err != nil {
		return nil, fmt.Errorf("failed to train model: %w", err)
	}

	res := map[string]interface{}{
		"data":       r,
		"names":      r.GetNames(),
		"coeff":      r.GetCoeffs(),
		"datapoints": r.GetDataPoints(),
		"graphics":   makeGraphics(r),
	}
	return res, nil
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
func ridgeChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Заполняем матрицы
	for i := 0; i < numOfSamples; i++ {
		for j := 0; j < numOfVars; j++ {
			variables.Set(i, j, dataPoints[i].Variables[j])
		}
		observed.Set(i, 0, dataPoints[i].Observed)
	}

	// Создаем модель Ridge Regression
	regr := linearmodel.NewRidge()
	regr.Alpha = 1
	regr.Tol = 0.01
	regr.Normalize = true
	regr.L1Ratio = 0

	// Обучаем модель
	regr.Fit(variables, observed)

	// Делаем предсказание
	Ypred := mat.NewDense(numOfSamples, 1, nil)
	regr.Predict(variables, Ypred)

	fmt.Println("Predicted Y:\n", mat.Formatted(Ypred))

	// Возвращаем результаты
	res := map[string]interface{}{
		"Ypred":              fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"LinearRegression":   regr.LinearRegression,
		"Coef":               regr.LinearRegression.Coef,
		"Solver":             regr.Solver,
		"Tol":                regr.Tol,
		"Alpha":              regr.Alpha,
		"L1Ratio":            regr.L1Ratio,
		"ActivationFunction": regr.ActivationFunction,
	}

	return res, nil
}
func lassoChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Заполняем матрицы
	for i := 0; i < numOfSamples; i++ {
		for j := 0; j < numOfVars; j++ {
			variables.Set(i, j, dataPoints[i].Variables[j])
		}
		observed.Set(i, 0, dataPoints[i].Observed)
	}

	// Создаем модель Lasso
	regr := linearmodel.NewLasso() // Предположим, что у вас есть Lasso модель
	regr.Alpha = 1
	regr.Tol = 0.01
	regr.Normalize = true
	// Обучаем модель
	regr.Fit(variables, observed)

	// Делаем предсказание
	Ypred := mat.NewDense(numOfSamples, 1, nil)
	regr.Predict(variables, Ypred)

	fmt.Println("Predicted Y:\n", mat.Formatted(Ypred))
	res := map[string]interface{}{
		"Ypred":            fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"LinearRegression": regr.LinearRegression,
		"MaxIter":          regr.MaxIter,
		"Tol":              regr.Tol,
		"Alpha":            regr.Alpha,
		"L1Ratio":          regr.L1Ratio,
		"Selection":        regr.Selection,
		"WarmStart":        regr.WarmStart,
		"Positive":         regr.Positive,
		"CDResult":         regr.CDResult,
	}
	return res, nil
}
func elasticChecking(dataPoints []models.DataPoint, l1Ratio float64) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Заполняем матрицы
	for i := 0; i < numOfSamples; i++ {
		for j := 0; j < numOfVars; j++ {
			variables.Set(i, j, dataPoints[i].Variables[j])
		}
		observed.Set(i, 0, dataPoints[i].Observed)
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

	fmt.Println("Predicted Y:\n", mat.Formatted(Ypred))

	// Возвращаем результаты
	res := map[string]interface{}{
		"Ypred":            fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"LinearRegression": enet.LinearRegression,
		"Coef":             enet.LinearRegression.Coef,
		"Tol":              enet.Tol,
		"Alpha":            enet.Alpha,
		"L1Ratio":          enet.L1Ratio,
	}

	return res, nil
}
func logisticChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Заполняем матрицы
	for i := 0; i < numOfSamples; i++ {
		for j := 0; j < numOfVars; j++ {
			variables.Set(i, j, dataPoints[i].Variables[j])
		}
		observed.Set(i, 0, dataPoints[i].Observed)
	}

	fmt.Println("X (variables):\n", mat.Formatted(variables))
	fmt.Println("Y (observed):\n", mat.Formatted(observed))

	// Создаем модель LogisticRegression
	regr := linearmodel.NewLogisticRegression()
	regr.Tol = 0.01
	regr.Alpha = 1

	// Обучаем модель
	regr.Fit(variables, observed)

	// Делаем предсказание
	Ypred := mat.NewDense(numOfSamples, 1, nil)
	regr.Predict(variables, Ypred)

	fmt.Println("Predicted Y:\n", mat.Formatted(Ypred))

	// Возвращаем результаты
	res := map[string]interface{}{
		"Ypred":     fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"Coef":      regr.Coef,
		"Intercept": regr.Intercept,
		"Tol":       regr.Tol,
		"Alpha":     regr.Alpha,
	}

	return res, nil
}
func SVRChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	var res map[string]interface{}
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)

	// Заполняем матрицы
	for i := 0; i < numOfSamples; i++ {
		for j := 0; j < numOfVars; j++ {
			variables.Set(i, j, dataPoints[i].Variables[j])
		}
		observed.Set(i, 0, dataPoints[i].Observed)
	}

	fmt.Println("X (variables):\n", mat.Formatted(variables))
	fmt.Println("Y (observed):\n", mat.Formatted(observed))
	randomState := base.NewLockedSource(7)
	xscaler := preprocessing.NewMinMaxScaler([]float64{-1, 1})
	yscaler := preprocessing.NewMinMaxScaler([]float64{-1, 1})
	Xsc, _ := xscaler.FitTransform(variables, nil)
	Ysc, _ := yscaler.FitTransform(observed, nil)
	Epsilon := 0.1 * yscaler.Scale.At(0, 0)
	Ypred := map[string]*mat.Dense{}
	for _, opt := range []struct {
		kernel                  string
		C, gamma, coef0, degree float64
	}{
		{kernel: "rbf", C: 1e3, gamma: .1},
		{kernel: "linear", C: 1e3},
		{kernel: "poly", gamma: 1, coef0: 1, C: 1e3, degree: 2},
	} {
		Ypred[opt.kernel] = &mat.Dense{}
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
		svr.Fit(Xsc, Ysc)
		svr.Predict(Xsc, Ypred[opt.kernel])
		Ypred[opt.kernel], _ = yscaler.InverseTransform(Ypred[opt.kernel], nil)
		fmt.Println(base.MatStr(variables, observed, Ypred[opt.kernel]))
		res["Ypred"] = fmt.Sprintf("%.2f\n", mat.Formatted(Ypred[opt.kernel]))

	}

	return res, nil
}
func polynomialChecking(dataPoints []models.DataPoint, degree int) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)

	// Создаем матрицы X (variables) и Y (observed)
	observed := mat.NewDense(numOfSamples, 1, nil)          // Y - вектор (numOfSamples × 1)
	variables := mat.NewDense(numOfSamples, numOfVars, nil) // X - матрица (numOfSamples × numOfVars)
	nSamples, _ := observed.Dims()
	// Заполняем матрицы
	for i := 0; i < numOfSamples; i++ {
		for j := 0; j < numOfVars; j++ {
			variables.Set(i, j, dataPoints[i].Variables[j])
		}
		observed.Set(i, 0, dataPoints[i].Observed)
	}

	fmt.Println("X (variables):\n", mat.Formatted(variables))
	fmt.Println("Y (observed):\n", mat.Formatted(observed))

	// Добавляем полиномиальные признаки
	poly := preprocessing.NewPolynomialFeatures(degree)
	poly.IncludeBias = false
	poly.Fit(variables, nil)
	Xp, _ := poly.Transform(variables, nil)

	_, nFeatures := Xp.Dims()
	_, nOutputs := observed.Dims()
	Ypred := mat.NewDense(nSamples, nOutputs, nil)

	Alpha := 1.
	mlp := neuralnetwork.NewMLPClassifier([]int{}, "logistic", "adam", Alpha)
	mlp.BatchSize = nSamples

	// we allocate Coef here because we use it for loss and grad tests before Fit
	mlp.Initializer(observed.RawMatrix().Cols, []int{nFeatures, nOutputs}, true, false)
	mlp.WarmStart = true
	mlp.Shuffle = false
	var J float64
	loss := func() float64 {
		mlp.MaxIter = 1
		mlp.Fit(Xp, observed)
		return mlp.Loss
	}
	chkLoss := func(context string, expectedLoss float64) {
		if math.Abs(J-expectedLoss) > 1e-3 {
			log.Errorf("%s J=%g expected:%g", context, J, expectedLoss)
		}
	}
	chkGrad := func(context string, expectedGradient []float64) {
		actualGradient := mlp.GetpackedGrads()[:len(expectedGradient)]

		fmt.Printf("%s grad=%v expected %v\n", context, actualGradient, expectedGradient)
		for j := 0; j < len(expectedGradient); j++ {
			if !scalar.EqualWithinAbs(expectedGradient[j], actualGradient[j], 1e-4) {
				log.Errorf("%s grad=%v expected %v", context, actualGradient, expectedGradient)
				return
			}
		}
	}
	mlp.Alpha = 10.

	J = loss()
	chkLoss("At test theta", 3.164)
	chkGrad("at test theta", []float64{0.3460, 0.1614, 0.1948, 0.2269, 0.0922})
	// try different solvers

	best := make(map[string]string)
	bestLoss := math.Inf(1)
	bestTime := time.Second * 86400

	// // test Fit with various base.Optimizer
	var Optimizers = []string{
		"sgd",
		// "adagrad",
		// "rmsprop",
		// "adadelta",
		"adam",
		"lbfgs",
	}

	checkGradients := func(problem optimize.Problem, initX []float64) {
		settings := &fd.Settings{Step: 1e-8}
		gradFromModel := make([]float64, len(initX))
		gradFromFD := make([]float64, len(initX))
		problem.Func(initX)
		problem.Grad(gradFromModel, initX)
		fd.Gradient(gradFromFD, problem.Func, initX, settings)
		for i := range initX {
			if math.Abs(gradFromFD[i]-gradFromModel[i]) > 1e-4 {
				panic(fmt.Errorf("bad gradient, expected:\n%.3f\ngot:\n%.3f", gradFromFD, gradFromModel))
			}
		}
	}

	for _, optimizer := range Optimizers {
		testSetup := optimizer
		mlp := neuralnetwork.NewMLPClassifier([]int{}, "logistic", optimizer, 1)
		mlp.RandomState = base.NewLockedSource(1)
		mlp.Initializer(observed.RawMatrix().Cols, []int{nFeatures, nOutputs}, true, false)
		for i := range mlp.GetpackedParameters() {
			mlp.GetpackedParameters()[i] = 0
		}
		mlp.WarmStart = true
		mlp.MaxIter = 400
		mlp.LearningRateInit = .11
		mlp.BatchSize = 118 //1,2,59,118
		mlp.SetbeforeMinimize(checkGradients)

		start := time.Now()
		mlp.Fit(Xp, observed)
		elapsed := time.Since(start)
		J := mlp.Loss

		if J < bestLoss {
			bestLoss = J
			best["best for loss"] = testSetup + fmt.Sprintf("(%g)", J)
		}
		if elapsed < bestTime {
			bestTime = elapsed
			best["best for time"] = testSetup + fmt.Sprintf("(%s)", elapsed)
		}
		mlp.Predict(Xp, Ypred)
		accuracy := metrics.AccuracyScore(observed, Ypred, true, nil)
		// accuracy should be over 0.83
		expectedAccuracy := 0.8305
		if accuracy < expectedAccuracy {
			log.Errorf("%s accuracy=%.3g expected:%.3g", optimizer, accuracy, expectedAccuracy)
		}
	}

	// Возвращаем результаты
	res := map[string]interface{}{
		"Ypred":  fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"best":   best,
		"acc":    metrics.AccuracyScore(observed, Ypred, true, nil),
		"allRes": mlp,
	}

	return res, nil
}

// выполнение Ridge регрессии
func (rs *regressionService) RidgeRegression(
	ctx context.Context,
	XData [][]float64,
	YData [][]float64,
	alpha float64,
	tol float64,
	normalize bool,
) (map[string]interface{}, error) {
	X := mat.NewDense(len(XData), len(XData[0]), nil)
	Y := mat.NewDense(len(YData), len(YData[0]), nil)

	for i := range XData {
		for j := range XData[i] {
			X.Set(i, j, XData[i][j])
		}
	}
	for i := range YData {
		for j := range YData[i] {
			Y.Set(i, j, YData[i][j])
		}
	}
	fmt.Println("%.2f\n", mat.Formatted(X))
	fmt.Println("%.2f\n", mat.Formatted(Y))
	regr := linearmodel.NewRidge()
	regr.Alpha = alpha
	regr.Tol = tol
	regr.Normalize = normalize
	regr.L1Ratio = 0

	regr.Fit(X, Y)

	Ypred := mat.NewDense(len(YData), len(YData[0]), nil)
	fmt.Println("%.2f\n", mat.Formatted(Ypred))
	regr.Predict(X, Ypred)
	res := map[string]interface{}{
		"Ypred":              fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"LinearRegression":   regr.LinearRegression,
		"Coef":               regr.LinearRegression.Coef,
		"Solver":             regr.Solver,
		"Tol":                regr.Tol,
		"Alpha":              regr.Alpha,
		"L1Ratio":            regr.L1Ratio,
		"ActivationFunction": regr.ActivationFunction,
	}
	return res, nil
}

// LassoRegression выполняет регрессию Lasso
func (rs *regressionService) LassoRegression(
	ctx context.Context,
	XData [][]float64, // Входные данные (матрица признаков)
	YData [][]float64, // Целевые данные (матрица меток)
	alpha float64, // Гиперпараметр регуляризации
	tol float64, // Допустимая ошибка
	normalize bool, // Флаг нормализации данных
) (map[string]interface{}, error) {
	// Преобразование данных в матрицы Gonum
	X := mat.NewDense(len(XData), len(XData[0]), nil)
	Y := mat.NewDense(len(YData), len(YData[0]), nil)

	for i := range XData {
		for j := range XData[i] {
			X.Set(i, j, XData[i][j])
		}
	}
	for i := range YData {
		for j := range YData[i] {
			Y.Set(i, j, YData[i][j])
		}
	}

	// Создаем модель Lasso
	regr := linearmodel.NewLasso() // Предположим, что у вас есть Lasso модель
	regr.Alpha = alpha
	regr.Tol = tol
	regr.Normalize = normalize
	regr.Fit(X, Y)

	// Делаем предсказания
	Ypred := mat.NewDense(len(YData), len(YData[0]), nil)
	regr.Predict(X, Ypred)
	res := map[string]interface{}{
		"Ypred":            fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"LinearRegression": regr.LinearRegression,
		"MaxIter":          regr.MaxIter,
		"Tol":              regr.Tol,
		"Alpha":            regr.Alpha,
		"L1Ratio":          regr.L1Ratio,
		"Selection":        regr.Selection,
		"WarmStart":        regr.WarmStart,
		"Positive":         regr.Positive,
		"CDResult":         regr.CDResult,
	}
	return res, nil
}

// ElasticNetRegression - функция расчёта
func (rs *regressionService) ElasticNetRegression(ctx context.Context, params models.ElasticNetParams) (map[string][]float64, error) {
	// Генерация данных
	rand.Seed(0)
	coef := mat.NewDense(params.NFeatures, 1, nil)
	for feat := 0; feat < 50; feat++ { // Только 10% признаков
		coef.Set(feat, 0, rand.NormFloat64())
	}

	X := mat.NewDense(params.NSamplesTrain+params.NSamplesTest, params.NFeatures, nil)
	x := X.RawMatrix().Data
	for i := range x {
		x[i] = rand.NormFloat64()
	}

	Y := &mat.Dense{}
	Y.Mul(X, coef)

	// Разделение на обучающую и тестовую выборки
	rowslice := func(X mat.RawMatrixer, start, end int) *mat.Dense {
		rm := X.RawMatrix()
		return mat.NewDense(end-start, rm.Cols, rm.Data[start*rm.Stride:end*rm.Stride])
	}
	Xtrain := rowslice(X, 0, params.NSamplesTrain)
	Xtest := rowslice(X, params.NSamplesTrain, params.NSamplesTrain+params.NSamplesTest)
	Ytrain := rowslice(Y, 0, params.NSamplesTrain)
	Ytest := rowslice(Y, params.NSamplesTrain, params.NSamplesTrain+params.NSamplesTest)

	// Генерация alpha и расчёт ошибок
	logalphas := make([]float64, params.NAlphas)
	for i := range logalphas {
		logalphas[i] = -5 + 8*float64(i)/float64(params.NAlphas)
	}
	trainErrors := make([]float64, params.NAlphas)
	testErrors := make([]float64, params.NAlphas)

	for ialpha, logalpha := range logalphas {
		enet := linearmodel.NewElasticNet()
		enet.L1Ratio = params.L1Ratio
		enet.Alpha = math.Pow(10, logalpha)
		enet.Fit(Xtrain, Ytrain)
		trainErrors[ialpha] = enet.Score(Xtrain, Ytrain)
		testErrors[ialpha] = enet.Score(Xtest, Ytest)
	}

	// Формируем результат
	result := map[string][]float64{
		"logalphas":   logalphas,
		"trainErrors": trainErrors,
		"testErrors":  testErrors,
	}

	return result, nil
}

// Newservice - создание сервиса
func Newservice(logger zerolog.Logger) externalApi.Regression {
	return &regressionService{
		logger: logger,
	}
}
