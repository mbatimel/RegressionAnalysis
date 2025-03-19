package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sync"

	"strconv"

	"math"

	linearmodel "github.com/mbatimel/RegressionAnalysis/internal/linear_model"
	"github.com/xuri/excelize/v2"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	externalApi "github.com/mbatimel/RegressionAnalysis/pkg/interfaces"
	"github.com/rs/zerolog"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
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
	res := map[string]interface{}{
		"data":  r,
		"coeff": r.GetCoeffs(),
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
		"coeff":      r.GetCoeffs(),
		"datapoints": r.GetDataPoints(),
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
		"coeff":      r.GetCoeffs(),
		"datapoints": r.GetDataPoints(),
	}
	return res, nil
}
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

	regr := linearmodel.NewRidge()
	regr.Alpha = alpha
	regr.Tol = tol
	regr.Normalize = normalize
	regr.L1Ratio = 0

	regr.Fit(X, Y)

	Ypred := mat.NewDense(len(YData), len(YData[0]), nil)
	regr.Predict(X, Ypred)
	res := map[string]interface{}{
		"Ypred":              fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
		"LinearRegression":   regr.LinearRegression,
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
		"Ypred":            fmt.Sprintf("%.2f\n", mat.Formatted(Ypred)),
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

func Newservice(logger zerolog.Logger) externalApi.Regression {
	return &regressionService{
		logger: logger,
	}
}
