package service

import (
	"fmt"
	"time"

	"math"

	"github.com/gofiber/fiber/v2/log"
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
	regr.L1Ratio = 10

	// Обучаем модель
	regr.Fit(variables, observed)

	// Делаем предсказание
	Ypred := mat.NewDense(numOfSamples, 1, nil)
	regr.Predict(variables, Ypred)

	fmt.Println("Predicted Y:\n", mat.Formatted(Ypred))

	// Возвращаем результаты
	res := map[string]interface{}{
		"Ypred":              fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"Coef":               fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Coef)),
		"XOffsetoef":         fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XOffset)),
		"XScale":             fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XScale)),
		"Intercept":          fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Intercept)),
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
	rss := &mat.VecDense{}
	rss.SubVec(Ypred.ColView(0), observed.ColView(0))
	rss.MulElemVec(rss, rss)
	res := map[string]interface{}{
		"Ypred":      fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"Coef":       fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Coef.T())),
		"XOffsetoef": fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XOffset.T())),
		"XScale":     fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.XScale.T())),
		"Intercept":  fmt.Sprintf("%.2f\n", mat.Formatted(regr.LinearRegression.Intercept.T())),
		"RSS":        mat.Sum(rss),
		"MaxIter":    regr.MaxIter,
		"Tol":        regr.Tol,
		"Alpha":      regr.Alpha,
		"L1Ratio":    regr.L1Ratio,
		"Selection":  regr.Selection,
		"WarmStart":  regr.WarmStart,
		"Positive":   regr.Positive,
		"CDResult":   regr.CDResult,
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
		"Ypred":      fmt.Sprintf("%.5f\n", mat.Formatted(Ypred)),
		"Coef":       fmt.Sprintf("%.7f\n", mat.Formatted(enet.LinearRegression.Coef)),
		"XOffsetoef": fmt.Sprintf("%.2f\n", mat.Formatted(enet.LinearRegression.XOffset)),
		"XScale":     fmt.Sprintf("%.2f\n", mat.Formatted(enet.LinearRegression.XScale)),
		"Intercept":  fmt.Sprintf("%.2f\n", mat.Formatted(enet.LinearRegression.Intercept)),
		"MaxIter":    enet.MaxIter,
		"Tol":        enet.Tol,
		"Alpha":      enet.Alpha,
		"L1Ratio":    enet.L1Ratio,
		"Selection":  enet.Selection,
		"WarmStart":  enet.WarmStart,
		"Positive":   enet.Positive,
		"CDResult":   enet.CDResult,
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
	regr.Alpha = 1e-5
	regr.MaxIter = 4

	beforeMinimize := func(problem optimize.Problem, initX []float64) {
		// check gradients
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
func svrChecking(dataPoints []models.DataPoint) (map[string]interface{}, error) {
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
		{kernel: "sigmoid", C: 1e3, gamma: .1},
		{kernel: "poly", gamma: 1, coef0: 1, C: 1e3, degree: 2},
		{kernel: "linear", C: 1e3},
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
		res = map[string]interface{}{
			"YPred " + opt.kernel: fmt.Sprintf("%.2f\n", mat.Formatted(Ypred[opt.kernel])),
			"Score " + opt.kernel: svr.Score(Xsc, Ysc),
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
	// Заполняем матрицы
	for i := 0; i < numOfSamples; i++ {
		for j := 0; j < numOfVars; j++ {
			variables.Set(i, j, dataPoints[i].Variables[j])
		}
		observed.Set(i, 0, dataPoints[i].Observed)
	}
	// Добавляем полиномиальные признаки
	poly := preprocessing.NewPolynomialFeatures(degree)
	poly.IncludeBias = false
	poly.Fit(variables, nil)
	Xp, _ := poly.Transform(variables, nil)

	_, nFeatures := Xp.Dims()
	_, nOutputs := observed.Dims()
	Ypred := mat.NewDense(nSamples, nOutputs, nil)

	best := make(map[string]string)
	bestLoss := math.Inf(1)
	bestTime := time.Second * 86400

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
			mlp.SetpackedParameters(i, 0)
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
		"Ypred": fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"best":  best,
		"acc":   metrics.AccuracyScore(observed, Ypred, true, nil),
	}

	return res, nil
}
