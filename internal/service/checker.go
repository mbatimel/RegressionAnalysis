package service

import (
	"fmt"

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

func ridgeChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint) (map[string]interface{}, error) {
	regr := linearmodel.NewRidge()
	regr.Alpha = 1
	regr.Tol = 0.01
	regr.Normalize = true
	regr.L1Ratio = 10
	regr.Fit(Xtrain, Ytrain)

	Ypred := mat.NewDense(Xtest.RawMatrix().Rows, 1, nil)
	regr.Predict(Xtest, Ypred)

	bestErr := map[string]float64{}
	r2score := metrics.R2Score(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	bestErr["R2"] = r2score
	mse := metrics.MeanSquaredError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	bestErr["MSE"] = mse
	mae := metrics.MeanAbsoluteError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	bestErr["MAE"] = mae

	if math.Sqrt(mse) > regr.Tol {
		fmt.Printf("Test %T normalize=%v r2score=%g mse=%g mae=%g\n", regr, true, r2score, mse, mae)
	}

	resultType := map[int]string{0: "Линейный", 1: "Линейный", 2: "Линейный"}

	res := map[string]interface{}{
		"ridge Ypred": fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"Coef":        fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Coef)),
		"graphics":    makeGraphicsForRidge(testPoints, regr.Coef, Ypred),
		"bestErr":     bestErr,
		"resultType":  resultType,
		"params": map[string]interface{}{
			"coefficients": flattenMatrix(regr.Coef),
			"intercept":    regr.Intercept,
			"lambda":       regr.Alpha,
		},
	}

	return res, nil
}
func MLRChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint) (map[string]interface{}, error) {
	mlr := linearmodel.NewLinearRegression()
	
	
	mlr.Normalize = true

	mlr.Fit(Xtrain, Ytrain)

	Ypred := mat.NewDense(Xtest.RawMatrix().Rows, 1, nil)
	mlr.Predict(Xtest, Ypred)

	bestErr := map[string]float64{}
	r2score := metrics.R2Score(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	bestErr["R2"] = r2score
	mse := metrics.MeanSquaredError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	bestErr["MSE"] = mse
	mae := metrics.MeanAbsoluteError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
	bestErr["MAE"] = mae

	resultType := map[int]string{0: "Линейный", 1: "Линейный", 2: "Линейный"}

	res := map[string]interface{}{
		"ridge Ypred": fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"Coef":        fmt.Sprintf("%.2f\n", mat.Formatted(mlr.Coef)),
		"graphics":    makeGraphicsForRidge(testPoints, mlr.Coef, Ypred),
		"bestErr":     bestErr,
		"resultType":  resultType,
		"params": map[string]interface{}{
			"coefficients": flattenMatrix(mlr.Coef),
			"intercept":    mlr.Intercept,
		},
	}

	return res, nil
}

func lassoChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint) (map[string]interface{}, error) {
	regr := linearmodel.NewMultiTaskLasso()
	regr.FitIntercept = true
	regr.Normalize = true
	regr.Alpha = 1e-5
	regr.L1Ratio = 1
	regr.MaxIter = 1e5
	regr.Tol = 1e-4
	regr.Fit(Xtrain, Ytrain)

	Ypred := mat.NewDense(Xtest.RawMatrix().Rows, 1, nil)
	regr.Predict(Xtest, Ypred)

	rss := &mat.VecDense{}
	rss.SubVec(Ypred.ColView(0), Ytest.ColView(0))
	rss.MulElemVec(rss, rss)

	bestErr := map[string]float64{
		"R2":  metrics.R2Score(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
		"MSE": metrics.MeanSquaredError(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
		"MAE": metrics.MeanAbsoluteError(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
	}

	resultType := map[int]string{0: "Линейный", 1: "Линейный", 2: "Линейный"}

	res := map[string]interface{}{
		"lasso Ypred": fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"Coef":        fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Coef.T())),
		"graphics":    makeGraphicsForLasso(testPoints, regr.Coef, Ypred),
		"bestErr":     bestErr,
		"resultType":  resultType,
		"params": map[string]interface{}{
			"coefficients": flattenMatrix(regr.Coef),
			"intercept":    regr.Intercept,
			"lambda":       regr.Alpha,
		},
	}

	return res, nil
}

func elasticChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint, l1Ratio float64) (map[string]interface{}, error) {

	// Создание и обучение модели
	enet := linearmodel.NewMultiTaskElasticNet()
	enet.Alpha = 1
	enet.Tol = 0.01
	enet.L1Ratio = l1Ratio

	enet.Fit(Xtrain, Ytrain)

	// Предсказание на тестовой выборке
	numTest := Xtest.RawMatrix().Rows
	Ypred := mat.NewDense(numTest, 1, nil)
	enet.Predict(Xtest, Ypred)

	// Вычисление метрик
	bestErr := map[string]float{
		"R2":  metrics.R2Score(Ytest, Ypred, nil, "").At(0, 0),
		"MSE": metrics.MeanSquaredError(Ytest, Ypred, nil, "").At(0, 0),
		"MAE": metrics.MeanAbsoluteError(Ytest, Ypred, nil, "").At(0, 0),
	}

	resultType := map[int]string{
		0: "Линейный",
		1: "Линейный",
		2: "Линейный",
	}

	res := map[string]interface{}{
		"elastic Ypred": fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"Coef":          fmt.Sprintf("%.7f\n", mat.Formatted(enet.LinearRegression.Coef)),
		"graphics":      makeGraphicsForElastic(testPoints, enet.Coef, Ypred),
		"bestErr":       bestErr,
		"resultType":    resultType,
		"params": map[string]interface{}{
			"coefficients": flattenMatrix(enet.Coef),
			"intercept":    enet.Intercept,
			"lambda1":      enet.Alpha * (1 - enet.L1Ratio),
			"lambda2":      enet.Alpha * enet.L1Ratio,
		},
	}

	return res, nil
}

func logisticChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint) (map[string]interface{}, error) {
	// checkGradients := func(problem optimize.Problem, initX []float64) {
	// 	settings := &fd.Settings{Step: 1e-8}
	// 	gradFromModel := make([]float64, len(initX))
	// 	gradFromFD := make([]float64, len(initX))
	// 	problem.Func(initX)
	// 	problem.Grad(gradFromModel, initX)
	// 	fd.Gradient(gradFromFD, problem.Func, initX, settings)
	// }

	// Обучение модели
	regr := linearmodel.NewLogisticRegression()
	regr.Alpha = 1e-5
	regr.MaxIter = 1000
	regr.Alpha = 1
	regr.Tol = 0.01

	// regr.BeforeMinimize = checkGradients

	regr.Fit(Xtrain, Ytrain)

	// Предсказание
	numTest := Xtest.RawMatrix().Rows
	Ypred := mat.NewDense(numTest, 1, nil)
	regr.Predict(Xtest, Ypred)

	// Метрики
	bestErr := map[string]float{
		"R2":  metrics.R2Score(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
		"MSE": metrics.MeanSquaredError(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
		"MAE": metrics.MeanAbsoluteError(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
	}

	resultType := map[int]string{
		0: "Линейный",
		1: "Линейный",
		2: "Линейный",
	}

	res := map[string]interface{}{
		"logistic Ypred": fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"Coef":           regr.Coef,
		"graphics":       makeGraphicsForLogistic(testPoints, regr.Coef, Ypred),
		"bestErr":        bestErr,
		"resultType":     resultType,
		"params": map[string]interface{}{
			"coefficients": flattenBlasMatrix(regr.Coef),
			"intercept":    regr.Intercept,
		},
	}

	return res, nil
}

func svrChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint) (map[string]interface{}, error) {
	const maxTrainSize = 50

	if Xtrain.RawMatrix().Rows > maxTrainSize {
		Xtrain, Ytrain = getRandomSubset(Xtrain, Ytrain, maxTrainSize)
	}

	var bestRes map[string]interface{}
	bestScore := struct {
		R2  float64
		MSE float64
		MAE float64
	}{
		R2:  math.Inf(-1),
		MSE: math.Inf(1),
		MAE: math.Inf(1),
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
		0: "Квадратичный",
		1: "Квадратичный",
		2: "Квадратичный",
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
		svr.MaxIter = 400

		svr.Fit(XtrainSc, YtrainSc)
		svr.Predict(XtestSc, Ypred)

		Ypred, _ = yscaler.InverseTransform(Ypred, nil)

		r2score := metrics.R2Score(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
		mse := metrics.MeanSquaredError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
		mae := metrics.MeanAbsoluteError(Ytest, Ypred, nil, "variance_weighted").At(0, 0)

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
			if opt.kernel == "linear" {
				resultType = map[int]string{
					0: "Линейный",
					1: "Линейный",
					2: "Линейный",
				}
			}
			if opt.kernel == "rbf" {
				resultType = map[int]string{
					0: "Кубический",
					1: "Кубический",
					2: "Кубический",
				}
			}
			if opt.kernel == "poly" {
				resultType = map[int]string{
					0: "Квадратичный",
					1: "Квадратичный",
					2: "Квадратичный",
				}
			}
			bestRes = map[string]interface{}{
				"svr YPred " + opt.kernel: fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
				"graphics":                makeGraphicsForSVR(testPoints, Ypred),
				"bestErr": map[string]float64{
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

func polynomialChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint) (map[string]interface{}, error) {
	rawVarsTrain := mat.DenseCopyOf(Xtrain)
	rawVarsTest := mat.DenseCopyOf(Xtest)

	type polyResult struct {
		degree   int
		r2       float64
		yPred    *mat.Dense
		mlp      *neuralnetwork.MLPClassifier
		graphics interface{}
		bestErr  map[string]float64
	}

	var best polyResult
	best.r2 = -math.MaxFloat64

	degreeToLabel := map[int]string{
		1: "Линейный",
		2: "Квадратичный",
		3: "Кубический",
		4: "Четвертичная дробь",
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
			return nil, fmt.Errorf("unmarshal error")
		}
		mlp.MaxIter = 400
		mlp.LearningRateInit = 0.11
		mlp.BatchSize = 118
		mlp.BeforeMinimize = checkGradients

		poly := preprocessing.NewPolynomialFeatures(degree)
		poly.IncludeBias = false
		poly.Fit(variablesTrain, Ytrain)
		variablesTrain, _ = poly.FitTransform(variablesTrain, nil)
		variablesTest, _ = poly.FitTransform(variablesTest, nil)

		mlp.Fit(variablesTrain, Ytrain)
		Ypred := mat.NewDense(len(testPoints), 1, nil)
		mlp.Predict(variablesTest, Ypred)

		r2 := metrics.R2Score(Ytest, Ypred, nil, "variance_weighted").At(0, 0)
		if r2 > best.r2 {
			bestErr := map[string]float64{
				"R2":  r2,
				"MSE": metrics.MeanSquaredError(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
				"MAE": metrics.MeanAbsoluteError(Ytest, Ypred, nil, "variance_weighted").At(0, 0),
			}
			best = polyResult{
				degree:   degree,
				r2:       r2,
				yPred:    Ypred,
				mlp:      mlp,
				graphics: makeGraphicsForPoly(testPoints, mlp.Coefs, Ypred, degree),
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
		"coefficients":  flattenNested(best.mlp.Coefs),
		"OutActivation": best.mlp.OutActivation,
		"poly accuracy": metrics.AccuracyScore(Ytest, best.yPred, true, nil),
		"bestErr":       best.bestErr,
		"resultType":    resultType,
		"params": map[string]interface{}{
			"coefficients": flattenNested(best.mlp.Coefs),
			"degree":       degreeToLabel[best.degree],
			"intercept":    best.mlp.Intercepts,
		},
	}

	return detail, nil
}

func logChecking(Xtrain, Ytrain, Xtest, Ytest *mat.Dense, testPoints []models.DataPoint) (map[string]interface{}, error) {

	// Перемешиваем и строим матрицы для логарифмического преобразования
	numOfVars := Xtrain.RawMatrix().Cols
	logVarsTrain := mat.NewDense(Xtrain.RawMatrix().Rows, numOfVars, nil)
	observedTrain := mat.NewDense(Ytrain.RawMatrix().Rows, 1, nil)
	for i := 0; i < Xtrain.RawMatrix().Rows; i++ {
		observedTrain.Set(i, 0, math.Log(Ytrain.At(i, 0)))
		for j := 0; j < numOfVars; j++ {
			x := Xtrain.At(i, j)
			if x <= 0 || math.IsInf(x, 0) || math.IsNaN(x) {
				logVarsTrain.Set(i, j, 0)
			} else {
				logVarsTrain.Set(i, j, math.Log(x))
			}
		}
	}

	// Для тестовой выборки
	logVarsTest := mat.NewDense(Xtest.RawMatrix().Rows, numOfVars, nil)
	observedTest := mat.NewDense(Ytest.RawMatrix().Rows, 1, nil)
	for i := 0; i < Xtest.RawMatrix().Rows; i++ {
		observedTest.Set(i, 0, math.Log(Ytest.At(i, 0)))
		for j := 0; j < numOfVars; j++ {
			x := Xtest.At(i, j)
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
		return nil, fmt.Errorf("Error with unmarshal byte data: %v", err)
	}
	mlp.MaxIter = 400
	mlp.LearningRateInit = 0.11
	mlp.BatchSize = 118
	mlp.BeforeMinimize = checkGradients
	mlp.BeforeMinimize = checkGradients

	// Обучение и предсказание
	mlp.Fit(logVarsTrain, observedTrain)

	YpredTest := mat.NewDense(len(testPoints), 1, nil)
	mlp.Predict(logVarsTest, YpredTest)

	// Метрики на тестовой выборке
	r2 := metrics.R2Score(observedTest, YpredTest, nil, "variance_weighted").At(0, 0)
	mse := metrics.MeanSquaredError(observedTest, YpredTest, nil, "variance_weighted").At(0, 0)
	mae := metrics.MeanAbsoluteError(observedTest, YpredTest, nil, "variance_weighted").At(0, 0)

	// Обработка бесконечных и NaN значений
	if math.IsInf(r2, 0) || math.IsNaN(r2) {
		r2 = 0
	}
	if math.IsInf(mse, 0) || math.IsNaN(mse) {
		mse = 0
	}
	if math.IsInf(mae, 0) || math.IsNaN(mae) {
		mae = 0
	}
	// Результаты
	res := map[string]interface{}{
		"graphics":   makeGraphicsForLog(testPoints, YpredTest),
		"poly Ypred": fmt.Sprintf("%.2f\n", mat.Formatted(YpredTest)),
		"params": map[string]interface{}{
			"coefficients": flattenNested(mlp.Coefs),
			"intercept":    mlp.Intercepts,
		},
		"OutActivation": mlp.OutActivation,
		"bestErr": map[string]float{
			"R2":  r2,
			"MSE": mse,
			"MAE": mae,
		},
		"resultType": map[int]string{
			0: "Логарифмический",
			1: "Логарифмический",
			2: "Логарифмический",
		},
	}

	return res, nil
}
