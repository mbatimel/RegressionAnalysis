package models

type Matrix struct {
	Values [][]float64
	Cols   int64
	Rows   int64
}

type DataPoint struct {
	Y float64   `json:"y"`
	X []float64 `json:"x"`
}
