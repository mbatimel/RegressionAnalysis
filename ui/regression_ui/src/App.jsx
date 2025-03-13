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
  const [tableData, setTableData] = useState([]);
  const [headers, setHeaders] = useState([]);
  
  const handleRowsChange = (e) => setRows(e.target.value)||1;
  const handleColsChange = (e) => setCols(e.target.value)||1;
  const updateTableData = (rowIndex, colIndex, value) => {
    const newData = [...tableData];
    if (!newData[rowIndex]) newData[rowIndex] = [];
    newData[rowIndex][colIndex] = value;
    setTableData(newData);
  };
  const handleMethodChange = (e) => {
    setSelectedMethod(e.target.value);
  };
  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };
  const [rows, setRows] = useState(1);
  const [cols, setCols] = useState(1);
  const generateTable = () => {
  if (cols < 1 || rows < 1) return <p>Введите корректные размеры таблицы</p>;
  return(
    <table border="1">
      <thead>
        <tr>
          {["Y", ...Array(cols - 1).fill(0).map((_, i) => `X${i + 1}`)].map((label, index) => (
            <th key={index}>
              <input
                type="text"
                placeholder={label}
                value={headers[index] || ""}
                onChange={(e) => {
                  const newHeaders = [...headers];
                  newHeaders[index] = e.target.value;
                  setHeaders(newHeaders);
                }}
              />
            </th>
          ))}
        </tr>
      </thead>
      <tbody>
        {Array.from({ length: rows }).map((_, rowIndex) => (
          <tr key={rowIndex}>
            {Array.from({ length: cols }).map((_, colIndex) => (
              <td key={colIndex}>
                <input
                  type="text"
                  value={tableData[rowIndex]?.[colIndex] || ""}
                  onChange={(e) => updateTableData(rowIndex, colIndex, e.target.value)}
                />
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
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
    if (!headers.length || !tableData.length) {
      alert("Введите данные в таблицу перед отправкой запроса");
      return;
    }

    const requestData = {
      observer: "Y",
      vars: Array.from({ length: cols - 1 }, (_, i) => `X${i + 1}`),
      dataPoints: tableData.map(row => ({
        obs: parseFloat(row[0]),
        vares: row.slice(1).map(Number),
      })),
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
        {generateTable()}
    </div>
      <header className="App-header">
        <h1>React cURL Buttons</h1>
        <button className="MLRButton" onClick={MLR}>MLR</button>
      {responseMLRData && (
        <div className="result-box">
          <h3>Формула:</h3>
          <p>Y = {headers[0]}</p>
          {headers.slice(1).map((name, i) => <p key={i}>X{i + 1} = {name}</p>)}
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