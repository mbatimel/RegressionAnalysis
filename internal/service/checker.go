package service

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"math"

	"github.com/mbatimel/RegressionAnalysis/internal/base"
	linearmodel "github.com/mbatimel/RegressionAnalysis/internal/linear_model"
	"github.com/mbatimel/RegressionAnalysis/internal/metrics"
	neuralnetwork "github.com/mbatimel/RegressionAnalysis/internal/neural_network"
	"github.com/mbatimel/RegressionAnalysis/internal/preprocessing"
	"github.com/mbatimel/RegressionAnalysis/internal/svm"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

type float = float64

func ridgeChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)
	// Перемешиваем и делим на train/test
	shuffled := make([]models.DataPoint, numOfSamples)
	copy(shuffled, dataPoints)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numOfSamples, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	splitIdx := int(0.7 * float64(numOfSamples))
	trainPoints := shuffled[:splitIdx]
	testPoints := shuffled[splitIdx:]

	numTrain := len(trainPoints)
	numTest := len(testPoints)

	// Создаем матрицы обучения
	Xtrain := mat.NewDense(numTrain, numOfVars, nil)
	Ytrain := mat.NewDense(numTrain, 1, nil)
	for i, dp := range trainPoints {
		for j := 0; j < numOfVars; j++ {
			Xtrain.Set(i, j, dp.Variables[j])
		}
		Ytrain.Set(i, 0, dp.Observed)
	}

	// Создаем матрицы теста с параллельной загрузкой
	Xtest := mat.NewDense(numTest, numOfVars, nil)
	Ytest := mat.NewDense(numTest, 1, nil)
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}
	rowChan := make(chan rowData, numTest)
	var wg sync.WaitGroup
	for i := 0; i < numTest; i++ {
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

	// Обучение модели
	regr := linearmodel.NewRidge()
	regr.Alpha = 1
	regr.Tol = 0.01
	regr.Normalize = true
	regr.L1Ratio = 10
	regr.Fit(Xtrain, Ytrain)

	// Предсказание
	Ypred := mat.NewDense(numTest, 1, nil)
	regr.Predict(Xtest, Ypred)

	bestErr := make(map[string]float)
	r2score := metrics.R2Score(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	tmpScore, ok := bestErr["R2"]
	if !ok || r2score > tmpScore {
		bestErr["R2"] = r2score
	}
	mse := metrics.MeanSquaredError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	tmpScore, ok = bestErr["MSE"]
	if !ok || mse < tmpScore {
		bestErr["MSE"] = mse
	}
	mae := metrics.MeanAbsoluteError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	tmpScore, ok = bestErr["MAE"]
	if !ok || mae < tmpScore {
		bestErr["MAE"] = mae
	}
	if math.Sqrt(mse) > regr.Tol {
		fmt.Printf("Test %T normalize=%v r2score=%g mse=%g mae=%g\n", regr, true, r2score, mse, mae)
	}

	resultType := map[int]string{
		1: "linear",
		0: "linear",
		2: "linear",
	}

	res := map[string]interface{}{
		"ridge Ypred":              fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"ridge Coef":               fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Coef)),
		"ridge XOffsetoef":         fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XOffset)),
		"ridge XScale":             fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XScale)),
		"ridge Intercept":          fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Intercept)),
		"ridge ActivationFunction": regr.ActivationFunction,
		"graphics":                 makeGraphicsForOtherMethod(testPoints, regr.Coef, Ypred),
		"bestErr":                  bestErr,
		"resultType":               resultType,
	}

	return res, nil

}
func lassoChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)
	// Перемешивание и разбиение на train/test
	shuffled := make([]models.DataPoint, numOfSamples)
	copy(shuffled, dataPoints)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numOfSamples, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	splitIdx := int(0.7 * float64(numOfSamples))
	trainPoints := shuffled[:splitIdx]
	testPoints := shuffled[splitIdx:]

	numTrain := len(trainPoints)
	numTest := len(testPoints)

	Xtrain := mat.NewDense(numTrain, numOfVars, nil)
	Ytrain := mat.NewDense(numTrain, 1, nil)

	for i, dp := range trainPoints {
		for j, val := range dp.Variables {
			Xtrain.Set(i, j, val)
		}
		Ytrain.Set(i, 0, dp.Observed)
	}

	Xtest := mat.NewDense(numTest, numOfVars, nil)
	Ytest := mat.NewDense(numTest, 1, nil)

	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}
	rowChan := make(chan rowData, numTest)
	var wg sync.WaitGroup

	for i := 0; i < numTest; i++ {
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

	// Обучение модели
	regr := linearmodel.NewMultiTaskLasso()
	regr.FitIntercept = true
	regr.Normalize = true
	regr.Alpha = 1e-5
	regr.L1Ratio = 1
	regr.MaxIter = 1e5
	regr.Tol = 1e-4
	regr.Fit(Xtrain, Ytrain)

	Ypred := mat.NewDense(numTest, 1, nil)
	regr.Predict(Xtest, Ypred)

	rss := &mat.VecDense{}
	rss.SubVec(Ypred.ColView(0), Ytest.ColView(0))
	rss.MulElemVec(rss, rss)

	bestErr := map[string]float{
		"R2":  metrics.R2Score(Ytest, Ypred, nil, "").At(0, 0),
		"MSE": metrics.MeanSquaredError(Ytest, Ypred, nil, "").At(0, 0),
		"MAE": metrics.MeanAbsoluteError(Ytest, Ypred, nil, "").At(0, 0),
	}

	resultType := map[int]string{
		0: "linear",
		1: "linear",
		2: "linear",
	}

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
		"graphics":         makeGraphicsForOtherMethod(testPoints, regr.Coef, Ypred),
		"bestErr":          bestErr,
		"resultType":       resultType,
	}

	return res, nil

}
func elasticChecking(dataPoints []models.DataPoint, l1Ratio float64) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	numOfVars := len(dataPoints[0].Variables)
	// Перемешивание и разбиение на train/test
	shuffled := make([]models.DataPoint, numOfSamples)
	copy(shuffled, dataPoints)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numOfSamples, func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	splitIdx := int(0.7 * float64(numOfSamples))
	trainPoints := shuffled[:splitIdx]
	testPoints := shuffled[splitIdx:]

	numTrain := len(trainPoints)
	numTest := len(testPoints)

	// Подготовка Xtrain и Ytrain
	Xtrain := mat.NewDense(numTrain, numOfVars, nil)
	Ytrain := mat.NewDense(numTrain, 1, nil)
	for i, dp := range trainPoints {
		for j, val := range dp.Variables {
			Xtrain.Set(i, j, val)
		}
		Ytrain.Set(i, 0, dp.Observed)
	}

	// Подготовка Xtest и Ytest
	Xtest := mat.NewDense(numTest, numOfVars, nil)
	Ytest := mat.NewDense(numTest, 1, nil)

	// Используем горутины для параллельной подготовки тестовой части
	type rowData struct {
		index    int
		varRow   []float64
		obsValue float64
	}
	rowChan := make(chan rowData, numTest)
	var wg sync.WaitGroup

	for i := 0; i < numTest; i++ {
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

	// Создание и обучение модели
	enet := linearmodel.NewMultiTaskElasticNet()
	enet.Alpha = 1
	enet.Tol = 0.01
	enet.L1Ratio = l1Ratio

	enet.Fit(Xtrain, Ytrain)

	// Предсказание на тестовой выборке
	Ypred := mat.NewDense(numTest, 1, nil)
	enet.Predict(Xtest, Ypred)

	// Вычисление метрик
	bestErr := make(map[string]float)
	bestErr["R2"] = metrics.R2Score(Ytest, Ypred, nil, "").At(0, 0)
	bestErr["MSE"] = metrics.MeanSquaredError(Ytest, Ypred, nil, "").At(0, 0)
	bestErr["MAE"] = metrics.MeanAbsoluteError(Ytest, Ypred, nil, "").At(0, 0)

	resultType := map[int]string{
		0: "linear",
		1: "linear",
		2: "linear",
	}

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
		"graphics":           makeGraphicsForOtherMethod(testPoints, enet.Coef, Ypred),
		"bestErr":            bestErr,
		"resultType":         resultType,
	}

	return res, nil

}

func logisticChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	if numOfSamples == 0 {
		return nil, fmt.Errorf("no data points provided")
	}
	// Перемешиваем и делим на обучающую/тестовую выборки
	shuffled := make([]models.DataPoint, numOfSamples)
	copy(shuffled, dataPoints)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numOfSamples, func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	splitIdx := int(0.7 * float64(numOfSamples))
	trainPoints := shuffled[:splitIdx]
	testPoints := shuffled[splitIdx:]

	numOfVars := len(trainPoints[0].Variables)
	numTrain := len(trainPoints)
	numTest := len(testPoints)

	// Формируем матрицы X и Y для обучения
	Xtrain := mat.NewDense(numTrain, numOfVars, nil)
	Ytrain := mat.NewDense(numTrain, 1, nil)
	for i, dp := range trainPoints {
		for j, val := range dp.Variables {
			Xtrain.Set(i, j, val)
		}
		Ytrain.Set(i, 0, dp.Observed)
	}

	// Формируем матрицы X и Y для теста
	Xtest := mat.NewDense(numTest, numOfVars, nil)
	Ytest := mat.NewDense(numTest, 1, nil)
	for i, dp := range testPoints {
		for j, val := range dp.Variables {
			Xtest.Set(i, j, val)
		}
		Ytest.Set(i, 0, dp.Observed)
	}

	checkGradients := func(problem optimize.Problem, initX []float64) {
		settings := &fd.Settings{Step: 1e-8}
		gradFromModel := make([]float64, len(initX))
		gradFromFD := make([]float64, len(initX))
		problem.Func(initX)
		problem.Grad(gradFromModel, initX)
		fd.Gradient(gradFromFD, problem.Func, initX, settings)
	}

	// Обучаем модель
	regr := linearmodel.NewLogisticRegression()
	regr.Alpha = 1e-5
	regr.MaxIter = 4
	regr.BeforeMinimize = checkGradients
	regr.Fit(Xtrain, Ytrain)

	// Предсказываем на тестовой выборке
	Ypred := mat.NewDense(numTest, 1, nil)
	regr.Predict(Xtest, Ypred)

	// Вычисляем метрики
	bestErr := make(map[string]float)
	r2score := metrics.R2Score(Ytest, Ypred, nil, "").At(0, 0)
	bestErr["R2"] = r2score
	mse := metrics.MeanSquaredError(Ytest, Ypred, nil, "").At(0, 0)
	bestErr["MSE"] = mse
	mae := metrics.MeanAbsoluteError(Ytest, Ypred, nil, "").At(0, 0)
	bestErr["MAE"] = mae

	resultType := map[int]string{
		0: "linear",
		1: "linear",
		2: "linear",
	}

	res := map[string]interface{}{
		"logistic Ypred":     fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"logistic Coef":      regr.Coef,
		"logistic Intercept": regr.Intercept,
		"logistic Tol":       regr.Tol,
		"logistic Alpha":     regr.Alpha,
		"graphics":           makeGraphicsFoBlas64(testPoints, regr.Coef, Ypred),
		"bestErr":            bestErr,
		"resultType":         resultType,
	}

	return res, nil

}

func svrChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	// Перемешиваем копию dataPoints
	shuffled := make([]models.DataPoint, len(dataPoints))
	copy(shuffled, dataPoints)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	// Разделение 70/30
	splitIdx := int(float64(len(shuffled)) * 0.7)
	trainPoints := shuffled[:splitIdx]
	testPoints := shuffled[splitIdx:]

	var bestRes map[string]interface{}
	bestScore := struct {
		R2  float64
		MSE float64
		MAE float64
	}{
		R2:  math.Inf(-1), // Самое маленькое значение для начала
		MSE: math.Inf(1),  // Самое большое значение для начала
		MAE: math.Inf(1),
	}

	numTrain := len(trainPoints)
	numVars := len(trainPoints[0].Variables)

	// Матрицы для обучения
	Xtrain := mat.NewDense(numTrain, numVars, nil)
	Ytrain := mat.NewDense(numTrain, 1, nil)
	for i, dp := range trainPoints {
		for j, val := range dp.Variables {
			Xtrain.Set(i, j, val)
		}
		Ytrain.Set(i, 0, dp.Observed)
	}

	// Матрицы для теста
	numTest := len(testPoints)
	Xtest := mat.NewDense(numTest, numVars, nil)
	Ytest := mat.NewDense(numTest, 1, nil)
	for i, dp := range testPoints {
		for j, val := range dp.Variables {
			Xtest.Set(i, j, val)
		}
		Ytest.Set(i, 0, dp.Observed)
	}

	randomState := base.NewLockedSource(7)
	xscaler := preprocessing.NewMinMaxScaler([]float64{-1, 1})
	yscaler := preprocessing.NewMinMaxScaler([]float64{-1, 1})
	XtrainSc, _ := xscaler.FitTransform(Xtrain, nil)
	YtrainSc, _ := yscaler.FitTransform(Ytrain, nil)
	XtestSc, _ := xscaler.Transform(Xtest, nil)
	Epsilon := 0.1 * yscaler.Scale.At(0, 0)

	kernelOptions := []struct {
		kernel                  string
		C, gamma, coef0, degree float64
	}{
		{kernel: "linear", C: 1e3},
		{kernel: "rbf", C: 1e3, gamma: 0.1},
		{kernel: "poly", C: 1e3, gamma: 1, coef0: 200, degree: 2},
	}
	resultType := map[int]string{
		0: "cubic",
		1: "cubic",
		2: "cubic",
	}
	for _, opt := range kernelOptions {
		Ypred := &mat.Dense{}

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

		svr.Fit(XtrainSc, YtrainSc)
		svr.Predict(XtestSc, Ypred)

		Ypred, _ = yscaler.InverseTransform(Ypred, nil)

		// Метрики
		r2score := metrics.R2Score(Ytest, Ypred, nil, "").At(0, 0)
		mse := metrics.MeanSquaredError(Ytest, Ypred, nil, "").At(0, 0)
		mae := metrics.MeanAbsoluteError(Ytest, Ypred, nil, "").At(0, 0)

		isBetter := false
		if r2score > bestScore.R2 {
			isBetter = true
		} else if math.Abs(r2score-bestScore.R2) < 1e-3 {
			if mse < bestScore.MSE || mae < bestScore.MAE {
				isBetter = true
			}
		}

		if isBetter {
			bestScore.R2 = r2score
			bestScore.MSE = mse
			bestScore.MAE = mae

			bestRes = map[string]interface{}{
				"svr YPred " + opt.kernel: fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
				"graphics":                makeGraphicsForSVR(testPoints, Ypred),
				"bestErr": map[string]float{
					"R2":  r2score,
					"MSE": mse,
					"MAE": mae,
				},
				"kernel":     svr.Kernel,
				"resultType": resultType,
			}
		}
	}

	return bestRes, nil
}

func polynomialChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	if numOfSamples == 0 {
		return nil, fmt.Errorf("no data points provided")
	}

	// Перемешиваем и разделяем выборку
	rand.Seed(42) // фиксируем seed для воспроизводимости
	rand.Shuffle(len(dataPoints), func(i, j int) {
		dataPoints[i], dataPoints[j] = dataPoints[j], dataPoints[i]
	})

	trainSize := int(0.7 * float64(numOfSamples))
	trainData := dataPoints[:trainSize]
	testData := dataPoints[trainSize:]

	numOfVars := len(dataPoints[0].Variables)

	// Создание матриц обучения
	observedTrain := mat.NewDense(len(trainData), 1, nil)
	rawVarsTrain := mat.NewDense(len(trainData), numOfVars, nil)
	for i, dp := range trainData {
		for j, val := range dp.Variables {
			rawVarsTrain.Set(i, j, val)
		}
		observedTrain.Set(i, 0, dp.Observed)
	}

	// Создание матриц теста
	observedTest := mat.NewDense(len(testData), 1, nil)
	rawVarsTest := mat.NewDense(len(testData), numOfVars, nil)
	for i, dp := range testData {
		for j, val := range dp.Variables {
			rawVarsTest.Set(i, j, val)
		}
		observedTest.Set(i, 0, dp.Observed)
	}

	type polyResult struct {
		degree   int
		r2       float64
		yPred    *mat.Dense
		mlp      *neuralnetwork.MLPClassifier
		graphics interface{}
		bestErr  map[string]float64
	}

	var best polyResult
	best.r2 = -math.MaxFloat64 // initialize to lowest possible

	degreeToLabel := map[int]string{
		1: "linear",
		2: "quadratic",
		3: "cubic",
		4: "quartic",
		5: "quintic",
	}

	for degree := 1; degree <= 3; degree++ {
		variablesTrain := rawVarsTrain
		variablesTest := rawVarsTest
		checkGradients := func(problem optimize.Problem, initX []float64) {
			settings := &fd.Settings{Step: 1e-8}
			gradFromModel := make([]float64, len(initX))
			gradFromFD := make([]float64, len(initX))
			problem.Func(initX)
			problem.Grad(gradFromModel, initX)
			fd.Gradient(gradFromFD, problem.Func, initX, settings)
		}

		buf := []byte(`{"activation": "logistic", "alpha": 0.0001, "batch_size": "auto", "beta_1": 0.9, "beta_2": 0.999, "early_stopping": false, "epsilon": 1e-08, "hidden_layer_sizes": [], "learning_rate": "constant", "learning_rate_init": 0.001, "max_iter": 400, "momentum": 0.9, "n_iter_no_change": 10, "nesterovs_momentum": true, "power_t": 0.5, "random_state": 7, "shuffle": true, "solver": "lbfgs", "tol": 0.0001, "validation_fraction": 0.1, "verbose": false, "warm_start": false, "out_activation_": "tanh", "intercepts_": [[0.5082271055138958]], "coefs_": [[[-0.18963335144967644], [0.2744326667319166], [-0.0068960058868800505], [-0.1870170339590578], [0.33640123639043934], [0.14343164310877599], [-0.2840940844068544], [-0.06035740527894848], [-0.015548157556294752], [-0.09766841821748058], [-0.13516966516561582], [0.01180873002271984], [-0.37004002347719184], [-0.3146740174229507], [-0.010236340304847167], [0.034725564039145625], [0.07596312959511524], [0.07031424991074327], [0.03226286238715042], [-0.11777688776136522], [-0.0862585580460505], [0.046039278168215306], [-0.32297687193126345], [0.004283074654547827], [0.013040383833634088], [-0.047491825368820184], [-0.12259098577236986]]]}`)

		mlp := neuralnetwork.NewMLPClassifier([]int{}, "", "", 0)
		mlp.RandomState = base.NewLockedSource(2)
		err := mlp.Unmarshal(buf)
		if err != nil {
			return nil, fmt.Errorf("Error with unmarshal byte data")
		}
		mlp.MaxIter = 400
		mlp.LearningRateInit = 0.11
		mlp.BatchSize = 118
		mlp.BeforeMinimize = checkGradients

		poly := preprocessing.NewPolynomialFeatures(degree)
		poly.IncludeBias = false
		poly.Fit(variablesTrain, observedTrain)
		variablesTrain, _ = poly.FitTransform(variablesTrain, nil)
		variablesTest, _ = poly.FitTransform(variablesTest, nil)

		// Обучаем и предсказываем
		mlp.Fit(variablesTrain, observedTrain)
		Ypred := mat.NewDense(len(testData), 1, nil)
		mlp.Predict(variablesTest, Ypred)

		r2 := metrics.R2Score(observedTest, Ypred, nil, "").At(0, 0)
		if r2 > best.r2 {
			bestErr := map[string]float{
				"R2":  r2,
				"MSE": metrics.MeanSquaredError(observedTest, Ypred, nil, "").At(0, 0),
				"MAE": metrics.MeanAbsoluteError(observedTest, Ypred, nil, "").At(0, 0),
			}
			best = polyResult{
				degree:   degree,
				r2:       r2,
				yPred:    Ypred,
				mlp:      mlp,
				graphics: makeGraphicsForPoly(testData, mlp.Coefs, Ypred),
				bestErr:  bestErr,
			}
		}
	}

	resultType := map[int]string{
		0: degreeToLabel[best.degree],
		1: degreeToLabel[best.degree],
		2: degreeToLabel[best.degree],
	}

	detail := map[string]interface{}{
		"graphics":      best.graphics,
		"poly Ypred":    fmt.Sprintf("%.2f\n", mat.Formatted(best.yPred)),
		"Coeffs":        best.mlp.Coefs,
		"poly accuracy": metrics.AccuracyScore(observedTest, best.yPred, true, nil),
		"bestErr":       best.bestErr,
		"resultType":    resultType,
	}

	return detail, nil
}

func logChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
	numOfSamples := len(dataPoints)
	if numOfSamples == 0 {
		return nil, fmt.Errorf("no data points provided")
	}

	// Перемешиваем и разделяем выборку
	rand.Seed(42)
	rand.Shuffle(numOfSamples, func(i, j int) {
		dataPoints[i], dataPoints[j] = dataPoints[j], dataPoints[i]
	})
	trainSize := int(0.7 * float64(numOfSamples))
	trainData := dataPoints[:trainSize]
	testData := dataPoints[trainSize:]

	numOfVars := len(dataPoints[0].Variables)

	// Подготовка обучающих данных
	observedTrain := mat.NewDense(len(trainData), 1, nil)
	logVarsTrain := mat.NewDense(len(trainData), numOfVars, nil)
	for i, dp := range trainData {
		y := math.Log(dp.Observed)
		observedTrain.Set(i, 0, y)
		for j, x := range dp.Variables {
			if x <= 0 || math.IsInf(x, 0) || math.IsNaN(x) {
				logVarsTrain.Set(i, j, 0)
			} else {
				logVarsTrain.Set(i, j, math.Log(x))
			}
		}
	}

	// Подготовка тестовых данных
	observedTest := mat.NewDense(len(testData), 1, nil)
	logVarsTest := mat.NewDense(len(testData), numOfVars, nil)
	for i, dp := range testData {
		y := math.Log(dp.Observed)
		observedTest.Set(i, 0, y)
		for j, x := range dp.Variables {
			if x <= 0 || math.IsInf(x, 0) || math.IsNaN(x) {
				logVarsTest.Set(i, j, 0)
			} else {
				logVarsTest.Set(i, j, math.Log(x))
			}
		}
	}

	checkGradients := func(problem optimize.Problem, initX []float64) {
		settings := &fd.Settings{Step: 1e-8}
		gradFromModel := make([]float64, len(initX))
		gradFromFD := make([]float64, len(initX))
		problem.Func(initX)
		problem.Grad(gradFromModel, initX)
		fd.Gradient(gradFromFD, problem.Func, initX, settings)
	}

	buf := []byte(`{"activation": "tanh", "alpha": 0.0001, "batch_size": "auto", "beta_1": 0.9, "beta_2": 0.999, "early_stopping": false, "epsilon": 1e-08, "hidden_layer_sizes": [], "learning_rate": "constant", "learning_rate_init": 0.001, "max_iter": 400, "momentum": 0.9, "n_iter_no_change": 10, "nesterovs_momentum": true, "power_t": 0.5, "random_state": 7, "shuffle": true, "solver": "sgd", "tol": 0.0001, "validation_fraction": 0.1, "verbose": false, "warm_start": false, "out_activation_": "tanh", "intercepts_": [[0.5082271055138958]], "coefs_": [[[-0.18963335144967644], [0.2744326667319166], [-0.0068960058868800505], [-0.1870170339590578], [0.33640123639043934], [0.14343164310877599], [-0.2840940844068544], [-0.06035740527894848], [-0.015548157556294752], [-0.09766841821748058], [-0.13516966516561582], [0.01180873002271984], [-0.37004002347719184], [-0.3146740174229507], [-0.010236340304847167], [0.034725564039145625], [0.07596312959511524], [0.07031424991074327], [0.03226286238715042], [-0.11777688776136522], [-0.0862585580460505], [0.046039278168215306], [-0.32297687193126345], [0.004283074654547827], [0.013040383833634088], [-0.047491825368820184], [-0.12259098577236986]]]}`)
	mlp := neuralnetwork.NewMLPClassifier([]int{}, "", "", 0.)
	// mlp.RandomState = base.NewLockedSource(1)
	// mlp.Shuffle = true
	// mlp.LearningRateInit = .02
	// mlp.WeightDecay = .001
	// mlp.MaxIter = 10000
	// mlp.LossFuncName = "binary_log_loss"
	err := mlp.Unmarshal(buf)
	if err != nil {
		return nil, fmt.Errorf("Error with unmarshal byte data")
	}
	// mlp.RandomState = base.NewLockedSource(2)
	// mlp.WarmStart = false
	// mlp.LearningRateInit = 0.11
	// mlp.BatchSize = numOfSamples + 1
	mlp.BeforeMinimize = checkGradients

	// Обучение и предсказание
	mlp.Fit(logVarsTrain, observedTrain)

	YpredTest := mat.NewDense(len(testData), 1, nil)
	mlp.Predict(logVarsTest, YpredTest)

	// Метрики на тестовой выборке
	r2 := metrics.R2Score(observedTest, YpredTest, nil, "").At(0, 0)
	mse := metrics.MeanSquaredError(observedTest, YpredTest, nil, "").At(0, 0)
	mae := metrics.MeanAbsoluteError(observedTest, YpredTest, nil, "").At(0, 0)

	if math.IsInf(r2, 0) || math.IsNaN(r2) {
		r2 = 0
	}
	if math.IsInf(mse, 0) || math.IsNaN(mse) {
		mse = 0
	}
	if math.IsInf(mae, 0) || math.IsNaN(mae) {
		mae = 0
	}

	res := map[string]interface{}{
		"graphics":   makeGraphicsForPoly(testData, mlp.Coefs, YpredTest),
		"poly Ypred": fmt.Sprintf("%.2f\n", mat.Formatted(YpredTest)),
		"Coeffs":     mlp.Coefs,
		"bestErr": map[string]float{
			"R2":  r2,
			"MSE": mse,
			"MAE": mae,
		},
		"resultType": map[int]string{
			0: "log",
			1: "log",
			2: "log",
		},
	}

	return res, nil
}
