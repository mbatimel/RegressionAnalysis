package models
import (
	"github.com/mbatimel/RegressionAnalysis/internal/models"
)
type MLRRequest struct {
    Observer   string             `json:"observer"`
    Variables  []string           `json:"variables"`
    DataPoints []models.DataPoint 		`json:"data_points"`
}

// Входная структура для JSON-запроса
type RidgeRequest struct {
	XData     [][]float64 `json:"XData"`
	YData     [][]float64 `json:"YData"`
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
	XData     [][]float64 `json:"XData"`
	YData     [][]float64 `json:"YData"`
	Alpha     float64     `json:"alpha"`
	Tol       float64     `json:"tol"`
	Normalize bool        `json:"normalize"`
}

// Выходная структура для JSON-ответа
type LassoResponse struct {
	YPred [][]float64 `json:"YPred"`
}

type RequestBody struct {
	Data       [][]float64 `json:"data"`       // Матрица X
	Classes    []float64   `json:"classes"`    // Классы Y
	Alpha      float64     `json:"alpha"`      // Скорость обучения
	MeshStep   float64     `json:"mesh_step"`  // Шаг сетки
	Visualize  bool        `json:"visualize"`  // Флаг для визуализации
}

// ResponseBody структура для выходных данных
type ResponseBody struct {
	Accuracy float64 `json:"accuracy"` // Точность
}

