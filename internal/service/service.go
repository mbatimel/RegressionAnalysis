package service

import (	
	"github.com/mbatimel/RegressionAnalysis/internal/models"
)

type ServiceRegression struct {
}

func MlrRegression(observer string, vars []string, dataPoints models.DataPoint) {
	// r := new(Regression)
	// r.SetObserved(fmt.Sprintf("%s", observer))
	// for i, v := range vars {
	// 	r.SetVar(i, fmt.Sprintf("%s", v))
	// }
	// //FIXME: надо дописать и придумать как сделать сервисную часть

}

func RidgeRegression(){}

func LossoRegression(){}

func LogisticRegression(){}

func ElasticNetRegression(){}

func BayesRegression(){}

func Newservice() {

}
