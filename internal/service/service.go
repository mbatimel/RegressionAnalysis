package service

import (
	"fmt"

	"github.com/mbatimel/RegressionAnalysis/internal/math"
	"github.com/mbatimel/RegressionAnalysis/internal/models"
)

type ServiceRegression struct {
}

func Mlr_regression(observer string, vars []string, dataPoints models.DataPoint) {
	r := new(math.Regression)
	r.SetObserved(fmt.Sprintf("%s", observer))
	for i, v := range vars {
		r.SetVar(i, fmt.Sprintf("%s", v))
	}
	//FIXME: надо дописать и придумать как сделать сервисную часть

}

func Newservice() {

}
