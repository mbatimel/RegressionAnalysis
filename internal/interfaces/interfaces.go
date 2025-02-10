package interfaces

// Package interfaces
// @tg version=0.0.1
// @tg backend=regression
// @tg title=`Regression API`
//
//go:generate tg transport --services . --out ../../internal/transport/jsonRPC/externalapi --outSwagger ../../swaggers/regression/swagger.yaml

import (
	"context"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/mat"
)

// Exeternal regression api ...
// @tg http-server metrics log
// @tg http-prefix=/api/v1
// @tg 200=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Resp200
// @tg 403=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err403
// @tg 405=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err405
// @tg 500=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err500
type Regression interface {
	// MlrRegression ...
	// @tg http-method=GET
	// @tg http-path=/mlr
	// @tg summary=`Ручка по рассчету MLR регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:MlrRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:MlrRegressionResp200
	MlrRegression(ctx context.Context, observer string, vars []string, dataPoints []models.DataPoint) (string, error)
	// RidgeRegression ...
	// @tg http-method=GET
	// @tg http-path=/ridge
	// @tg summary=`Ручка по рассчету ridge регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:RidgeRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:RidgeRegressionResp200
	RidgeRegression(ctx context.Context, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) (*mat.Dense, error)
	// LassoRegression ...
	// @tg http-method=GET
	// @tg http-path=/lasso
	// @tg summary=`Ручка по рассчету lasso регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:LassoRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:LassoRegressionResp200
	LassoRegression(ctx context.Context, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) (*mat.Dense, error)
	// ElasticNetRegression ...
	// @tg http-method=GET
	// @tg http-path=/elasticNet
	// @tg summary=`Ручка по рассчету Elastic Net регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:ElasticNetRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:ElasticNetRegressionResp200
	ElasticNetRegression(ctx context.Context, params models.ElasticNetParams) (map[string][]float64, error)
}
