package customhandlers

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mbatimel/RegressionAnalysis/internal/config"
	"github.com/mbatimel/RegressionAnalysis/internal/errors"
	"github.com/mbatimel/RegressionAnalysis/internal/models"
	regression "github.com/mbatimel/RegressionAnalysis/pkg/interfaces"
	"github.com/rs/zerolog/log"
)

const serviceName = "Regression"

func MlrRegression(ctx *fiber.Ctx, svc regression.Regression, observer string, vars []string, dataPoints []models.DataPoint) error {
	var (
		methodName = "MlrRegression"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":     "get",
			"path":       "/mlr",
			"observer":   observer,
			"vars":       vars,
			"dataPoints": dataPoints,
			"service":    serviceName,
			"took":       time.Since(begin).String(),
		}
		l := log.Info()
		if err != nil {
			if errors.Is(err, errors.ForbiddenError()) {
				l = log.Warn().Err(err)
			} else {
				l = log.Error().Err(err)
			}
		}
		l.Fields(fields).Msg("call")

		metrics.RequestLatency.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	defer func() {
		metrics.HttpCollector.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Add(1)
	}()

	formula, err := svc.MlrRegression(ctx.Context(), observer, vars, dataPoints)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, formula, nil)
	return err
}
func MlrRegressionCSV(ctx *fiber.Ctx, svc regression.Regression, observer string, vars []string, file multipart.File) error {
	var (
		methodName = "MlrRegressionCSV"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":   "get",
			"path":     "/mlrCSV",
			"observer": observer,
			"vars":     vars,
			"file":     file,
			"service":  serviceName,
			"took":     time.Since(begin).String(),
		}
		l := log.Info()
		if err != nil {
			if errors.Is(err, errors.ForbiddenError()) {
				l = log.Warn().Err(err)
			} else {
				l = log.Error().Err(err)
			}
		}
		l.Fields(fields).Msg("call")

		metrics.RequestLatency.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	defer func() {
		metrics.HttpCollector.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Add(1)
	}()

	formula, err := svc.MlrRegressionCSV(ctx.Context(), observer, vars, file)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, formula, nil)
	return err
}

func RidgeRegression(ctx *fiber.Ctx, svc regression.Regression, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) error {
	var (
		methodName = "RidgeRegression"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":      "get",
			"path":        "/ridge",
			"handlerName": methodName,
			"xdata":       XData,
			"ydata":       YData,
			"alpha":       alpha,
			"tol":         tol,
			"norma;ize":   normalize,
			"service":     serviceName,
			"took":        time.Since(begin).String(),
		}
		l := log.Info()
		if err != nil {
			if errors.Is(err, errors.ForbiddenError()) {
				l = log.Warn().Err(err)
			} else {
				l = log.Error().Err(err)
			}
		}
		l.Fields(fields).Msg("call")

		metrics.RequestLatency.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	defer func() {
		metrics.HttpCollector.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Add(1)
	}()

	res, err := svc.RidgeRegression(ctx.Context(), XData, YData, alpha, tol, normalize)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, res, nil)
	return err
}

func LassoRegression(ctx *fiber.Ctx, svc regression.Regression, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) error {
	var (
		methodName = "LassoRegression"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":      "get",
			"path":        "/lasso",
			"handlerName": methodName,
			"xdata":       XData,
			"ydata":       YData,
			"alpha":       alpha,
			"tol":         tol,
			"norma;ize":   normalize,
			"service":     serviceName,
			"took":        time.Since(begin).String(),
		}
		l := log.Info()
		if err != nil {
			if errors.Is(err, errors.ForbiddenError()) {
				l = log.Warn().Err(err)
			} else {
				l = log.Error().Err(err)
			}
		}
		l.Fields(fields).Msg("call")

		metrics.RequestLatency.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	defer func() {
		metrics.HttpCollector.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Add(1)
	}()

	res, err := svc.LassoRegression(ctx.Context(), XData, YData, alpha, tol, normalize)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, res, nil)
	return err
}

func ElasticNetRegression(ctx *fiber.Ctx, svc regression.Regression, params models.ElasticNetParams) error {
	var (
		methodName = "ElasticNetRegression"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":      "get",
			"path":        "/lasso",
			"handlerName": methodName,
			"params":      params,
			"service":     serviceName,
			"took":        time.Since(begin).String(),
		}
		l := log.Info()
		if err != nil {
			if errors.Is(err, errors.ForbiddenError()) {
				l = log.Warn().Err(err)
			} else {
				l = log.Error().Err(err)
			}
		}
		l.Fields(fields).Msg("call")

		metrics.RequestLatency.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	defer func() {
		metrics.HttpCollector.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Add(1)
	}()

	res, err := svc.ElasticNetRegression(ctx.Context(), params)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, res, nil)
	return err
}
