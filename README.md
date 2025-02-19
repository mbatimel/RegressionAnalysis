# RegressionAnalysis



ДЛЯ MLR
curl -X GET http://localhost:9000/api/v1/mlr \
-H "Content-Type: application/json" \
-d '{
    "observer": "Murders per annum per 1,000,000 inhabitants",
    "vars": ["Inhabitants", "Percent with incomes below $5000", "Percent unemployed"],
    "dataPoints": [
        {"obs": 11.2, "vares": [587000, 16.5, 6.2]},
        {"obs": 13.4, "vares": [643000, 20.5, 6.4]},
        {"obs": 40.7, "vares": [635000, 26.3, 9.3]}
    ]
}'


Ridge
curl -X GET http://localhost:9000/api/v1/ridge \
-H "Content-Type: application/json" \
-d '{
  "XData": [[0, 0], [1, 1], [2, 2]],
  "YData": [[0, 0], [1, 1], [2, 2]],
  "alpha": 1.0,
  "tol": 0.001,
  "normalize": false
}'

ожидаемый ответ:
{
  "YPred": [
    [0.2, 0.2],
    [1.0, 1.0],
    [1.8, 1.8]
  ]
}



curl -X GET http://localhost:9000/api/v1/lasso \
-H "Content-Type: application/json" \
-d '{
  "XData": [[0, 0], [1, 1], [2, 2]],
  "YData": [[0, 0], [1, 1], [2, 2]],
  "alpha": 0.5,
  "tol": 0.001,
  "normalize": false
}'




elasticnet
curl -X GET http://localhost:9000/api/v1/elasticnet \
  -H "Content-Type: application/json" \
  -d '{
    "params":{
    "n_samples_train": 75,
    "n_samples_test": 150,
    "n_features": 500,
    "l1_ratio": 0.7,
    "n_alphas": 20
    }
  }'
