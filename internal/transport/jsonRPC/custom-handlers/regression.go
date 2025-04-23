package customhandlers

import (
	"fmt"
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
			"method":     "post",
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
func MlrRegressionCSV(ctx *fiber.Ctx, svc regression.Regression, file []byte) error {
	var (
		methodName = "MlrRegressionCSV"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":  "post",
			"path":    "/mlrCSV",
			"file":    file,
			"service": serviceName,
			"took":    time.Since(begin).String(),
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
	formula, err := svc.MlrRegressionCSV(ctx.Context(), file)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, formula, nil)
	return err
}

func MlrRegressionExcel(ctx *fiber.Ctx, svc regression.Regression, file []byte) error {
	var (
		methodName = "MlrRegressionExcel"
		err        error
	)

	metrics := config.Metrics()
	defer func(begin time.Time) {
		fields := map[string]interface{}{
			"method":  "post",
			"path":    "/mlrExecl",
			"file":    file,
			"service": serviceName,
			"took":    time.Since(begin).String(),
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
	formula, err := svc.MlrRegressionExcel(ctx.Context(), file)
	if err != nil {
		sendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	sendResponse(ctx, log.Logger, formula, nil)
	return err
}
