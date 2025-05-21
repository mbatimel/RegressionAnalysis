package middlewares

import (
	"context"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	external_service "github.com/mbatimel/RegressionAnalysis/pkg/interfaces"
)

type middleware struct {
	regression external_service.Regression
}

func (m *middleware) MlrRegression(ctx context.Context, observer string, vars []string, dataPoints []models.DataPoint) (formula map[string]interface{}, err error) {
	return m.regression.MlrRegression(ctx, observer, vars, dataPoints)
}
func (m *middleware) MlrRegressionCSV(ctx context.Context, file []byte) (formula map[string]interface{}, err error) {
	return m.regression.MlrRegressionCSV(ctx, file)
}
func (m *middleware) MlrRegressionExcel(ctx context.Context, file []byte) (formula map[string]interface{}, err error) {
	return m.regression.MlrRegressionExcel(ctx, file)
}
func (m *middleware) OnlyMlrRegression(ctx context.Context, observer string, vars []string, dataPoints []models.DataPoint) (formula map[string]interface{}, err error) {
	return m.regression.OnlyMlrRegression(ctx, observer, vars, dataPoints)
}
func (m *middleware) OnlyMlrRegressionCSV(ctx context.Context, file []byte) (formula map[string]interface{}, err error) {
	return m.regression.OnlyMlrRegressionCSV(ctx, file)
}
func (m *middleware) OnlyMlrRegressionExcel(ctx context.Context, file []byte) (formula map[string]interface{}, err error) {
	return m.regression.OnlyMlrRegressionExcel(ctx, file)
}

func Newmiddleware(regression external_service.Regression) external_service.Regression {
	return &middleware{
		regression: regression,
	}
}
