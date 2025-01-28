package interfacesq

import "github.com/mbatimel/RegressionAnalysis/internal/models"

type Regression interface {
	Mlr_regression(observer string, vars []string, dataPoints models.DataPoint)
}
