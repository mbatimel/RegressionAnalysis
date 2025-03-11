import React, { useState } from 'react';
import './styles/App.css';
import './styles/Buttons.css';
import { Header } from './components/Header/Header';
import Documentation from './components/Documentation/Documentation';

function App() {
  const [responseMLRData, setMLRData] = useState(null);
  const [responseRidgeData, setRidgeData] = useState(null);
  const [lassoResponse, setLassoResponse] = useState(null);
  const [elasticnetResponse, setElasticnetResponse] = useState(null);
  const [file, setFile] = useState(null);
  const [responseMLRCSVData, setMLRCSVData] = useState(null);
  const [selectedMethod, setSelectedMethod] = useState("MLR");
  
  const handleMethodChange = (e) => {
    setSelectedMethod(e.target.value);
  };
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };
  const [rows, setRows] = useState(0);
  const [cols, setCols] = useState(0);

  const handleRowsChange = (e) => setRows(e.target.value);
  const handleColsChange = (e) => setCols(e.target.value);

  const generateTable = () => {
    let table = [];
    for (let i = 0; i < rows; i++) {
      let row = [];
      for (let j = 0; j < cols; j++) {
        row.push(<td key={j}>Row {i + 1}, Col {j + 1}</td>);
      }
      table.push(<tr key={i}>{row}</tr>);
    }
    return table;
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
      <div>
      <div>
        <label>
          Количество строк:
          <input
            type="number"
            value={rows}
            onChange={handleRowsChange}
          />
        </label>
      </div>
      <div>
        <label>
          Количество столбцов:
          <input
            type="number"
            value={cols}
            onChange={handleColsChange}
          />
        </label>
      </div>
      <table border="1">
        <tbody>
          {generateTable()}
        </tbody>
      </table>
    </div>
      <header className="App-header">
        <h1>React cURL Buttons</h1>
        <button className="MLRButton" onClick={MLR}>MLR </button>
            {responseMLRData && (
            <div className="result-box">
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(responseMLRData, null, 2)}</pre>
            </div>
          )}
        <button className="RidgeButton" onClick={Ridge}>Ridge</button>
        {responseRidgeData && (
            <div className="result-box">
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(responseRidgeData, null, 2)}</pre>
            </div>
          )}
        <button className="LassoButton" onClick={Lasso}>Lasso</button>
        {lassoResponse && (
            <div className="result-box">
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(lassoResponse, null, 2)}</pre>
            </div>
          )}
        <button className="ElastNetButton" onClick={ElasticNet}>Elastic</button>
        {elasticnetResponse && (
            <div className="result-box">
              <h3>Ответ от сервера:</h3>
              <pre>{JSON.stringify(elasticnetResponse, null, 2)}</pre>
            </div>
          )}

        <h1>React File Upload</h1>
        <input type="file" onChange={handleFileChange} /> 
        <button className="UploadButton" onClick={uploadFile}>Отправить файл</button>
        {responseMLRData && (
          <div className="result-box">
            <h3>Ответ от сервера:</h3>
            <pre>{JSON.stringify(responseMLRCSVData, null, 2)}</pre>
          </div>
        )}
        <h2>📜 Документация по методам регрессии</h2>
        <div>
          <label>Выберите метод регрессии: </label>
          <select value={selectedMethod} onChange={handleMethodChange}>
            <option value="MLR">MLR (Multiple Linear Regression)</option>
            <option value="Ridge">Ridge Regression</option>
            <option value="Lasso">Lasso Regression</option>
            <option value="ElasticNet">Elastic Net Regression</option>
          </select>
        </div>

        {/* Подключаем компонент с документацией */}
        <Documentation selectedMethod={selectedMethod} />
      </header>
    </div>
  );
}

export default App;

