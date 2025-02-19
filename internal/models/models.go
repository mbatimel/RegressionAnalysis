package models

type Matrix struct {
	Values [][]float64
	Cols   int64
	Rows   int64
}

type DataPoint struct {
	Observed  float64 `json:"obs,omitempty"`
	Variables []float64 `json:"vares,omitempty"`
	Predicted float64
	Error     float64
}

// ElasticNetParams - структура для входных данных API
type ElasticNetParams struct {
	NSamplesTrain int     `json:"n_samples_train" binding:"required"`
	NSamplesTest  int     `json:"n_samples_test" binding:"required"`
	NFeatures     int     `json:"n_features" binding:"required"`
	L1Ratio       float64 `json:"l1_ratio" binding:"required"`
	NAlphas       int     `json:"n_alphas" binding:"required"`
}

type LogisticRegressionParams struct {
	Alpha float64 `json:"alpha" binding:"required"`
}
