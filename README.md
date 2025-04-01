# RegressionAnalysis



ДЛЯ MLR
curl -X POST http://localhost:9000/api/v1/mlr \
-H "Content-Type: application/json" \
-d '{
    "observer": "Murders per annum per 1,000,000 inhabitants",
    "vars": ["Inhabitants", "Percent with incomes below $5000", "Percent unemployed"],
    "dataPoints": [
        {"obs": 11.2, "vares": [587000, 16.5, 6.2]},
        {"obs": 13.4, "vares": [643000, 20.5, 6.4]},
        {"obs": 40.7, "vares": [635000, 26.3, 9.3]},
        {"obs": 5.3, "vares": [692000, 16.5, 5.3]},
        {"obs": 24.8, "vares": [1248000, 19.2, 7.3]}
    ]
}'


Ridge
curl -X POST http://localhost:9000/api/v1/ridge \
 -H "Content-Type: application/json" \
 -d '{
   "XData": [[0, 0], [1, 1], [2, 2]],
   "YData": [[0, 0], [1, 1], [2, 2]],
   "alpha": 1.0,
   "tol": 0.001,
   "normalize": true
 }'

ожидаемый ответ:
{
  "YPred": [
    [0.2, 0.2],
    [1.0, 1.0],
    [1.8, 1.8]
  ]
}



curl -X POST http://localhost:9000/api/v1/lasso \
-H "Content-Type: application/json" \
-d '{
  "XData": [[587000, 643000,635000,692000,1248000], [16.5, 20.5, 26.3, 16.5, 19.2], [6.2, 6.4, 9.3, 5.3, 7.3]],
  "YData": [[11.2, 13.4, 40.7, 5.3, 24.8], [11.2, 13.4, 40.7, 5.3, 24.8], [11.2, 13.4, 40.7, 5.3, 24.8]],
  "alpha": 0.5,
  "tol": 0.001,
  "normalize": false
}'




elasticnet
curl -X POST http://localhost:9000/api/v1/elasticnet \
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

MLR CSV
curl -X POST http://localhost:9000/api/v1/mlrCSV \
     -H "Content-Type: multipart/form-data" \
     -F "file=@/Users/macbook/Desktop/ДИПЛОМ/RegressionAnalysis/examples/autos3.csv"
     
curl -X POST http://localhost:9000/api/v1/mlrExcel \
     -H "Content-Type: multipart/form-data" \
     -F "file=@/Users/macbook/Desktop/ДИПЛОМ/RegressionAnalysis/examples/auto.xlsx"

/Users/macbook/Desktop/ДИПЛОМ/RegressionAnalysis/examples/autos2.csv