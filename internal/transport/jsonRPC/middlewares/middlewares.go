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

func (m *middleware) RidgeRegression(ctx context.Context, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) (formula map[string]interface{}, err error) {
	return m.regression.RidgeRegression(ctx, XData, YData, alpha, tol, normalize)
}

func (m *middleware) LassoRegression(ctx context.Context, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) (formula map[string]interface{}, err error) {
	return m.regression.LassoRegression(ctx, XData, YData, alpha, tol, normalize)
}

func (m *middleware) ElasticNetRegression(ctx context.Context, params models.ElasticNetParams) (formula map[string][]float64, err error) {
	return m.regression.ElasticNetRegression(ctx, params)
}

func Newmiddleware(regression external_service.Regression) external_service.Regression {
	return &middleware{
		regression: regression,
	}
}
