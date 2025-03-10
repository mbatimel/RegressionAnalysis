import React, { useState } from 'react';
import './styles/App.css';
import { Header } from './components/Header/Header';

function App() {
  const [responseMLRData, setMLRData] = useState(null);
  const [responseRidgeData, setRidgeData] = useState(null);
  const [lassoResponse, setLassoResponse] = useState(null);
  const [elasticnetResponse, setElasticnetResponse] = useState(null);
  const [file, setFile] = useState(null);
  const [responseMLRCSVData, setMLRCSVData] = useState(null);
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };
  
  const uploadFile = async () => {
    if (!file) {
      alert("Выберите файл перед отправкой");
      return;
    }
    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch('/api/v1/mlrCSV', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Ошибка загрузки файла');
      }

      const data = await response.json();
      setMLRCSVData(data);
      console.log(data);
    } catch (error) {
      console.error('Ошибка:', error);
    }
  };

  const MLR = async () => {
    const requestData = {
      observer: "Murders per annum per 1,000,000 inhabitants",
      vars: [
        "Inhabitants",
        "Percent with incomes below $5000",
        "Percent unemployed"
      ],
      dataPoints: [
        { obs: 11.2, vares: [587000, 16.5, 6.2] },
        { obs: 13.4, vares: [643000, 20.5, 6.4] },
        { obs: 40.7, vares: [635000, 26.3, 9.3] },
        { obs: 5.3, vares: [692000, 16.5, 5.3] },
        { obs: 24.8, vares: [1248000, 19.2, 7.3] }
      ]
    };

    try {
      const response = await fetch('/api/v1/mlr', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error('Request failed');
      } 

      const data = await response.json();
      setMLRData(data);  
      console.log(data); 
    } catch (error) {
      console.error('Error:', error);
    }
  };
  const Ridge = async () => {
    const requestData = {
      XData: [
        [0, 0],
        [1, 1],
        [2, 2],
      ],
      YData: [
        [0, 0],
        [1, 1],
        [2, 2],
      ],
      alpha: 1,
      tol: 0.001,
      normalize: false,
    };

    try {
      const response = await fetch('/api/v1/ridge', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error('Request failed');
      }

      const data = await response.json();
      setRidgeData(data);  // Сохранение ответа в state
      console.log(data); // Вывод данных в консоль
    } catch (error) {
      console.error('Error:', error);
    }
  };
  const Lasso = async () => {
    const requestData = {
      XData: [
        [0, 0],
        [1, 1],
        [2, 2],
      ],
      YData: [
        [0, 0],
        [1, 1],
        [2, 2],
      ],
      alpha: 0.5,
      tol: 0.001,
      normalize: false,
    };

    try {
      const response = await fetch('/api/v1/lasso', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error('Lasso Request failed');
      }

      const data = await response.json();
      setLassoResponse(data);
      console.log(data);
    } catch (error) {
      console.error('Error:', error);
    }
  };

  // Функция для запроса ElasticNet
  const ElasticNet = async () => {
    const requestData = {
      params: {
        n_samples_train: 75,
        n_samples_test: 150,
        n_features: 500,
        l1_ratio: 0.7,
        n_alphas: 20,
      },
    };

    try {
      const response = await fetch('/api/v1/elasticnet', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error('ElasticNet Request failed');
      }

      const data = await response.json();
      setElasticnetResponse(data);
      console.log(data);
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <div className="App">
      <Header isRed={true}>
        <span>Regression analytics</span>
      </Header>
      <header className="App-header">
        <h1>React cURL Buttons</h1>
        <button onClick={MLR}>MLR </button>
            {responseMLRData && (
            <div>
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(responseMLRData, null, 2)}</pre>
            </div>
          )}
        <button onClick={Ridge}>Ridge</button>
        {responseRidgeData && (
            <div>
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(responseRidgeData, null, 2)}</pre>
            </div>
          )}
        <button onClick={Lasso}>Lasso</button>
        {lassoResponse && (
            <div>
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(lassoResponse, null, 2)}</pre>
            </div>
          )}
        <button onClick={ElasticNet}>Elastic</button>
        {elasticnetResponse && (
            <div>
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(elasticnetResponse, null, 2)}</pre>
            </div>
          )}

        <h1>React File Upload</h1>
        <input type="file" onChange={handleFileChange} />
        <button onClick={uploadFile}>Отправить файл</button>
        {responseMLRData && (
          <div>
            <h3>Ответ от сервера:</h3>
            <pre>{JSON.stringify(responseMLRCSVData, null, 2)}</pre>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;

