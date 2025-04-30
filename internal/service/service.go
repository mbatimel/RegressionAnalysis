package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sync"

	"strconv"

	linearmodel "github.com/mbatimel/RegressionAnalysis/internal/linear_model"

	"github.com/xuri/excelize/v2"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	externalApi "github.com/mbatimel/RegressionAnalysis/pkg/interfaces"
	"github.com/rs/zerolog"
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

	rs.logger.Info().Msg("Starting checking on classifier ridge")
	ridgeCheck, err := ridgeChecking(dataPoints)
	if err != nil {
		ridgeCheck = nil
		rs.logger.Error().Msg("Ridge checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier lasso")
	lassoCheck, err := lassoChecking(dataPoints)
	if err != nil {
		lassoCheck = nil
		rs.logger.Error().Msg("lasso checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier elastic")
	elasticCheck, err := elasticChecking(dataPoints, 1000)
	if err != nil {
		elasticCheck = nil
		rs.logger.Error().Msg("elastic checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier logistic")
	logisticCheck, err := logisticChecking(dataPoints)
	if err != nil {
		logisticCheck = nil
		rs.logger.Error().Msg("Ridge checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier SVR")
	svrCheck, err := svrChecking(dataPoints)
	if err != nil {
		svrCheck = nil
		rs.logger.Error().Err(err).Msg("SRV checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier polynomial")
	polynomialCheck, err := polynomialChecking(dataPoints)
	if err != nil {
		polynomialCheck = nil
		rs.logger.Error().Msg("Polynomial checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier logChecking")
	logChecking, err := logChecking(dataPoints)
	if err != nil {
		logChecking = nil
		rs.logger.Error().Msg("logChecking checking is dead")
	}

	res := map[string]interface{}{
		"Ridge":       ridgeCheck,
		"Lasso":       lassoCheck,
		"Elastic":     elasticCheck,
		"Logistic":    logisticCheck,
		"SRV":         svrCheck,
		"Polynomial":  polynomialCheck,
		"data":        r,
		"names":       r.GetNames(),
		"coeff":       r.GetCoeffs(),
		"logChecking": logChecking,
		// "graphics":   makeGraphics(r), Убрал и не думаю что нам пока нужен MLR
		"datapoints": r.GetDataPoints(),
	}

	return res, nil
}
func (rs *regressionService) MlrRegressionCSV(ctx context.Context, file []byte) (map[string]interface{}, error) {
	r := new(linearmodel.Regression)
	reader := csv.NewReader(bytes.NewReader(file))
	reader.Comma = ';'
	dataPoints := make([]models.DataPoint, 0)
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

			dataPoint := models.DataPoint{Observed: observed, Variables: variables}
			dataPoints = append(dataPoints, dataPoint)
			dataChan <- dataPoint
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
	rs.logger.Info().Msg("Starting checking on classifier ridge")
	ridgeCheck, err := ridgeChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier lasso")
	lassoCheck, err := lassoChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("lasso checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier elastic")
	elasticCheck, err := elasticChecking(dataPoints, 1000)
	if err != nil {
		return nil, fmt.Errorf("elastic checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier logistic")
	logisticCheck, err := logisticChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier SVR")
	svrCheck, err := svrChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("SRV checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier polynomial")
	polynomialCheck, err := polynomialChecking(dataPoints)
	if err != nil {
		polynomialCheck = nil
		rs.logger.Error().Msg("Polynomial checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier logChecking")
	logChecking, err := logChecking(dataPoints)
	if err != nil {
		logChecking = nil
		rs.logger.Error().Msg("logChecking checking is dead")
	}
	fmt.Println(logChecking)
	res := map[string]interface{}{
		"Ridge":       ridgeCheck,
		"Lasso":       lassoCheck,
		"Elastic":     elasticCheck,
		"Logistic":    logisticCheck,
		"SRV":         svrCheck,
		"Polynomial":  polynomialCheck,
		"LogChecking": logChecking,
		"data":        r,
		"names":       r.GetNames(),
		"coeff":       r.GetCoeffs(),
		// "graphics":   makeGraphics(r), Убрал и не думаю что нам пока нужен MLR
		"datapoints": r.GetDataPoints(),
	}
	return res, nil
}
func (rs *regressionService) MlrRegressionExcel(ctx context.Context, file []byte) (map[string]interface{}, error) {
	r := new(linearmodel.Regression)
	dataPoints := make([]models.DataPoint, 0)
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

			dataPoint := models.DataPoint{Observed: observed, Variables: variables}
			dataPoints = append(dataPoints, dataPoint)
			dataChan <- dataPoint
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

	rs.logger.Info().Msg("Starting checking on classifier ridge")
	ridgeCheck, err := ridgeChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier lasso")
	lassoCheck, err := lassoChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("lasso checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier elastic")
	elasticCheck, err := elasticChecking(dataPoints, 10)
	if err != nil {
		return nil, fmt.Errorf("elastic checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier logistic")
	logisticCheck, err := logisticChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("Ridge checking is dead")
	}

	rs.logger.Info().Msg("Starting checking on classifier SVR")
	svrCheck, err := svrChecking(dataPoints)
	if err != nil {
		return nil, fmt.Errorf("SRV checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier polynomial")
	polynomialCheck, err := polynomialChecking(dataPoints)
	if err != nil {
		polynomialCheck = nil
		rs.logger.Error().Msg("Polynomial checking is dead")
	}
	rs.logger.Info().Msg("Starting checking on classifier logChecking")
	logChecking, err := logChecking(dataPoints)
	if err != nil {
		logChecking = nil
		rs.logger.Error().Err(err).Msg("logChecking checking is dead")
	}
	res := map[string]interface{}{
		"Ridge":       ridgeCheck,
		"Lasso":       lassoCheck,
		"Elastic":     elasticCheck,
		"Logistic":    logisticCheck,
		"SRV":         svrCheck,
		"Polynomial":  polynomialCheck,
		"data":        r,
		"names":       r.GetNames(),
		"coeff":       r.GetCoeffs(),
		"logChecking": logChecking,
		// "graphics":   makeGraphics(r), Убрал и не думаю что нам пока нужен MLR
		"datapoints": r.GetDataPoints(),
	}
	return res, nil
}

// Newservice - создание сервиса
func Newservice(logger zerolog.Logger) externalApi.Regression {
	return &regressionService{
		logger: logger,
	}
}
