package models

import (
	linearmodel "github.com/mbatimel/RegressionAnalysis/internal/linear_model"
	"go/types"
)

type MlrRegressionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data map[string]interface{} `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}
type MlrRegressionRespCSV200 struct {
	// @tg desc=`массив объектов оплат`
	Data map[string]interface{} `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}
type MlrRegressionRespExcel200 struct {
	// @tg desc=`массив объектов оплат`
	Data map[string]interface{} `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}
type RidgeRegressionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data map[string]interface{} `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}
type LassoRegressionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data map[string]interface{} `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}

type ElasticNetRegressionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data map[string][]float64 `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}

type ResponseRegressionMlrRegression struct {
	Result *linearmodel.Regression `json:"result,omitempty"`
}
