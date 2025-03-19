import React, { useState } from 'react';
import './styles/App.css';
import './components/Buttons/Buttons.css';
import { Header } from './components/Header/Header';
import Documentation from './components/Documentation/Documentation';
import Tables from './components/Tables/Tables';
import Buttons from './components/Buttons/Buttons';
import Graphics from "./components/Graphics/Graphics";


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
  const [rows, setRows] = useState(1);
  const [cols, setCols] = useState(1);
  const XData = tableData.map((row) => row.slice(1).map(Number)); // Преобразуем X в числа
  const YData = tableData.map((row) => [Number(row[0])]); // Преобразуем Y в числа


  const handleMethodChange = (e) => {
    setSelectedMethod(e.target.value);
  };

  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  const uploadFile = async () => {
    if (!file) {
      alert("Выберите файл перед отправкой");
      return;
    }
    setTableData([]);
    setHeaders([]);
    const formData = new FormData();
    formData.append("file", file);
    
    let url = "/api/v1/mlrCSV";
    if (file.name.endsWith(".xlsx") || file.name.endsWith(".xls")) {
      url = "/api/v1/mlrExcel";
    }

    try {
      const response = await fetch(url, {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("Ошибка загрузки файла");
      }

      const data = await response.json();
      if (url.includes("mlrCSV")) {
        setMLRCSVData(data);
      } else {
        setMLRCSVData(data);
      }
      console.log(data);
    } catch (error) {
      console.error("Ошибка:", error);
    }
  };

  const MLR = async () => {
    if (!headers.length || !tableData.length) {
      alert("Введите данные в таблицу перед отправкой запроса");
      return;
    }
    setMLRCSVData(null);
    const requestData = {
      observer: "Y",
      vars: Array.from({ length: cols - 1 }, (_, i) => `X${i + 1}`),
      dataPoints: tableData.map((row) => ({
        obs: parseFloat(row[0]),
        vares: row.slice(1).map(Number),
      })),
    };

    try {
      const response = await fetch("/api/v1/mlr", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error("Request failed");
      }

      const data = await response.json();
      setMLRData(data);
      console.log(data);
    } catch (error) {
      console.error("Error:", error);
    }
  };

  const Ridge = async () => {
    if (!headers.length || !tableData.length) {
      alert("Введите данные в таблицу перед отправкой запроса");
      return;
    }
    setMLRCSVData(null);
    const requestData = {
      XData,
      YData,
      alpha: 1,
      tol: 0.001,
      normalize: false,
    };
  
    try {
      const response = await fetch("/api/v1/ridge", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestData),
      });
  
      if (!response.ok) throw new Error("Ridge Request failed");
  
      const data = await response.json();
      setRidgeData(data);
      console.log(data);
    } catch (error) {
      console.error("Error:", error);
    }
  };
  

  const Lasso = async () => {
    if (!headers.length || !tableData.length) {
      alert("Введите данные в таблицу перед отправкой запроса");
      return;
    }
    setMLRCSVData(null);
    const requestData = {
      XData,
      YData,
      alpha: 0.5,
      tol: 0.001,
      normalize: false,
    };
  
    try {
      const response = await fetch("/api/v1/lasso", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestData),
      });
  
      if (!response.ok) throw new Error("Lasso Request failed");
  
      const data = await response.json();
      setLassoResponse(data);
      console.log(data);
    } catch (error) {
      console.error("Error:", error);
    }
  };
  

  const ElasticNet = async () => {
    if (!headers.length || !tableData.length) {
      alert("Введите данные в таблицу перед отправкой запроса");
      return;
    }
    setMLRCSVData(null);
    const requestData = {
      XData,
      YData,
      params: {
        l1_ratio: 0.7,
        n_alphas: 20,
      },
    };
  
    try {
      const response = await fetch("/api/v1/elasticnet", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestData),
      });
  
      if (!response.ok) throw new Error("ElasticNet Request failed");
  
      const data = await response.json();
      setElasticnetResponse(data);
      console.log(data);
    } catch (error) {
      console.error("Error:", error);
    }
  };

  return (
    <div className="App">
      <Header isRed={true}>
        <span>Regression analytics</span>
      </Header>

      <Tables tableData={tableData} setTableData={setTableData} headers={headers} setHeaders={setHeaders} rows={rows} setRows={setRows} cols={cols} setCols={setCols} />
      <Graphics tableData={tableData} headers={responseMLRCSVData?.data?.names||headers} datapoints={responseMLRCSVData?.data?.datapoints || []}/>
      <Buttons
  MLR={MLR}
  Ridge={Ridge}
  Lasso={Lasso}
  ElasticNet={ElasticNet}
  uploadFile={uploadFile}
  handleFileChange={handleFileChange}
  file={file}
  responses={{ responseMLRData, responseRidgeData, lassoResponse, elasticnetResponse, responseMLRCSVData }}
  headers={headers} 
/>

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
      <Documentation selectedMethod={selectedMethod} />
    </div>
  );
}

export default App;
