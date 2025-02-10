package config

import (
	metricsRegression "github.com/mbatimel/RegressionAnalysis/internal/metrics"
)

var metrics *metricsRegression.Metrics

func Metrics() *metricsRegression.Metrics {
	if metrics == nil {
		metrics = metricsRegression.CreateMetrics("mbatimel", "regression", nil)
	}

	return metrics
}
