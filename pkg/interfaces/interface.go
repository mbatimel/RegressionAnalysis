// Package interfaces
// @tg version=0.0.1
// @tg backend=regression
// @tg title=`Regression API`
//
//go:generate tg transport --services . --out ../../internal/transport/jsonRPC/externalapi --outSwagger ../../swaggers/regression/swagger.yaml
package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
)

// regression
// @tg http-server metrics log
// @tg http-prefix=/api/v1
// @tg 200=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Resp200
// @tg 403=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err403
// @tg 405=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err405
// @tg 500=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err500
type Regression interface {
	// MlrRegression ...
	// @tg http-method=POST
	// @tg http-path=/mlr
	// @tg summary=`Ручка по рассчету MLR регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:MlrRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:MlrRegressionResp200
	MlrRegression(ctx context.Context, observer string, vars []string, dataPoints []models.DataPoint) (formula string, err error)
	// MlrRegressionCSV ...
	// @tg http-method=POST
	// @tg http-path=/mlrCSV
	// @tg summary=`Ручка по рассчету MLR регрессии из scv`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:MlrRegressionCSV
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:MlrRegressionResp200
	MlrRegressionCSV(ctx context.Context, observer string, vars []string, file multipart.File) (formula string, err error)
	//
	// RidgeRegression ...
	// @tg http-method=POST
	// @tg http-path=/ridge
	// @tg summary=`Ручка по рассчету ridge регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:RidgeRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:RidgeRegressionResp200
	RidgeRegression(ctx context.Context, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) (formula string, err error)
	//
	// LassoRegression ...
	// @tg http-method=POST
	// @tg http-path=/lasso
	// @tg summary=`Ручка по рассчету lasso регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:LassoRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:LassoRegressionResp200
	LassoRegression(ctx context.Context, XData [][]float64, YData [][]float64, alpha float64, tol float64, normalize bool) (formula string, err error)
	//
	// ElasticNetRegression ...
	// @tg http-method=POST
	// @tg http-path=/elasticNet
	// @tg summary=`Ручка по рассчету Elastic Net регрессии`
	// @tg http-response=github.com/mbatimel/RegressionAnalysis/internal/transport/jsonRPC/custom-handlers:ElasticNetRegression
	// @tg desc=`Ручка возвращает формулу и параметры`
	// @tg 400=github.com/mbatimel/RegressionAnalysis/swaggers/externalApi/models:Err400
	// @tg 200=github.com/mbatimel/RegressionAnalysis/pkg/models:ElasticNetRegressionResp200
	ElasticNetRegression(ctx context.Context, params models.ElasticNetParams) (formula map[string][]float64, err error)
	//TODO: MlrRegressionExcel
	// TODO: переделываем вывод  результата во всех функция
	// TODO: Логистик регрессия
}
