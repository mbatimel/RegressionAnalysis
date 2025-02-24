package models

import (
	"go/types"

	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/mat"
)

type MLRRequest struct {
	Observer   string             `json:"observer"`
	Variables  []string           `json:"variables"`
	DataPoints []models.DataPoint `json:"dataPoints"`
}

// Входная структура для JSON-запроса
type RidgeRequest struct {
	Xdata     [][]float64 `json:"xdata"`
	Ydata     [][]float64 `json:"ydata"`
	Alpha     float64     `json:"alpha"`
	Tol       float64     `json:"tol"`
	Normalize bool        `json:"normalize"`
}

// Выходная структура для JSON-ответа
type RidgeResponse struct {
	YPred [][]float64 `json:"YPred"`
}

// Входная структура для JSON-запроса
type LassoRequest struct {
	Xdata     [][]float64 `json:"xdata"`
	Ydata     [][]float64 `json:"ydata"`
	Alpha     float64     `json:"alpha"`
	Tol       float64     `json:"tol"`
	Normalize bool        `json:"normalize"`
}

// Выходная структура для JSON-ответа
type LassoResponse struct {
	YPred [][]float64 `json:"YPred"`
}

type RequestBody struct {
	Data      [][]float64 `json:"data"`      // Матрица X
	Classes   []float64   `json:"classes"`   // Классы Y
	Alpha     float64     `json:"alpha"`     // Скорость обучения
	MeshStep  float64     `json:"mesh_step"` // Шаг сетки
	Visualize bool        `json:"visualize"` // Флаг для визуализации
}

// ResponseBody структура для выходных данных
type ResponseBody struct {
	Accuracy float64 `json:"accuracy"` // Точность
}

type MlrRegressionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data string `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}
type RidgeRegressionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data *mat.Dense `json:"data"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	Error bool `json:"error"`
	// @tg example=``
	ErrorText        string    `json:"errorText"`
	AdditionalErrors types.Nil `json:"additionalErrors"`
}
type LassoRegressionResp200 struct {
	// @tg desc=`массив объектов оплат`
	Data *mat.Dense `json:"data"`
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

type RequestRegressionMlrRegression struct {
	Observer   string             `json:"observer,omitempty"`
	Vars       []string           `json:"vars,omitempty"`
	DataPoints []models.DataPoint `json:"dataPoints,omitempty"`
}
type ResponseRegressionMlrRegression struct {
	Result string `json:"result,omitempty"`
}
