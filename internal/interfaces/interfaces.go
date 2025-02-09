package interfacesq

import (
	"github.com/mbatimel/RegressionAnalysis/internal/models"
	"gonum.org/v1/gonum/mat"
)

type Regression interface {
 	MlrRegression(observer string, vars []string, dataPoints []models.DataPoint) (string, error)
	RidgeRegression(XData [][]float64,YData [][]float64,alpha float64,tol float64,normalize bool) (*mat.Dense, error) 
	LassoRegression(XData [][]float64,YData [][]float64,alpha float64,tol float64,normalize bool) (*mat.Dense, error) 
	ElasticNetRegression(params models.ElasticNetParams) (map[string][]float64, error)
}
