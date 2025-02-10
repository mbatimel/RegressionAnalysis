package service

import (
	"fmt"
	"math"

	linearmodel "github.com/mbatimel/RegressionAnalysis/internal/linear_model"
	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

func MlrRegression(observer string, vars []string, dataPoints []models.DataPoint) (string, error) {
	r := new(linearmodel.Regression)
	r.SetObserved(observer)
	for i, v := range vars {
		r.SetVar(i, v)
	}
	for _, dp := range dataPoints {
		linearDP := convertToLinearDataPoint(dp)
		r.Train(linearDP)
	}
	r.Run()

	// Вывод результатов (можно заменить на логирование или возврат результата)
	return fmt.Sprintf("Regression formula:\n%v\n", r.Formula), nil

}
func convertToLinearDataPoint(dp models.DataPoint) *models.DataPoint {
	return &models.DataPoint{
		Observed:  dp.Observed,
		Variables: dp.Variables,
	}
}
func RidgeRegression(
	XData [][]float64,
	YData [][]float64,
	alpha float64,
	tol float64,
	normalize bool,
) (*mat.Dense, error) {
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

	regr := linearmodel.NewRidge()
	regr.Alpha = alpha
	regr.Tol = tol
	regr.Normalize = normalize

	err := regr.Fit(X, Y)
	if err != nil {
		return nil, fmt.Errorf("ошибка обучения модели: %v", err)
	}

	Ypred := mat.NewDense(len(YData), len(YData[0]), nil)
	regr.Predict(X, Ypred)

	return Ypred, nil
}

// LassoRegression выполняет регрессию Lasso
func LassoRegression(
	XData [][]float64, // Входные данные (матрица признаков)
	YData [][]float64, // Целевые данные (матрица меток)
	alpha float64, // Гиперпараметр регуляризации
	tol float64, // Допустимая ошибка
	normalize bool, // Флаг нормализации данных
) (*mat.Dense, error) {
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

	// Обучаем модель
	err := regr.Fit(X, Y)
	if err != nil {
		return nil, fmt.Errorf("ошибка обучения модели: %v", err)
	}

	// Делаем предсказания
	Ypred := mat.NewDense(len(YData), len(YData[0]), nil)
	regr.Predict(X, Ypred)

	return Ypred, nil
}

// ElasticNetRegression - функция расчёта
func ElasticNetRegression(params models.ElasticNetParams) (map[string][]float64, error) {
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

func LogisticRegression() {}

func BayesRegression() {}

func Newservice() {

}
