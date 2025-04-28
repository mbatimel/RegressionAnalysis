import React, { useState } from 'react';
import './styles/App.css';
import './components/Buttons/Buttons.css';
import { Header } from './components/Header/Header';
import Documentation from './components/Documentation/Documentation';
import Tables from './components/Tables/Tables';
import Buttons from './components/Buttons/Buttons';
// import Graphics from "./components/Graphics/Graphics";
import Results from "./components/Results/Results";



function App() {
  const [responseMLRData, setMLRData] = useState(null);
  const [file, setFile] = useState(null);
  const [responseMLRCSVData, setMLRCSVData] = useState(null);
  const [selectedMethod, setSelectedMethod] = useState("Nil");
  const [tableData, setTableData] = useState([]);
  const [headers, setHeaders] = useState([]);
  const [rows, setRows] = useState(1);
  const [cols, setCols] = useState(1);
  const [isLoading, setIsLoading] = useState(false);
  const [highlightBestMethod, setHighlightBestMethod] = useState(false);



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
  
    setIsLoading(true);
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
  
      if (!response.ok) throw new Error("Ошибка загрузки файла");
  
      const data = await response.json();
      setMLRCSVData(data);
      console.log(data);
    } catch (error) {
      console.error("Ошибка:", error);
    } finally {
      setIsLoading(false);
    }
  };
  

  const MLR = async () => {
    if (!headers.length || !tableData.length) {
      alert("Введите данные в таблицу перед отправкой запроса");
      return;
    }
  
    setIsLoading(true);
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
  
      if (!response.ok) throw new Error("Request failed");
  
      const data = await response.json();
      setMLRData(data);
      console.log(data);
    } catch (error) {
      console.error("Error:", error);
    } finally {
      setIsLoading(false);
    }
  };
  
  const getRecommendedMethod = (data) => {
    const metrics = Object.entries(data)
      .filter(([method]) => !['datapoints', 'graphics'].includes(method))
      .map(([method, values]) => ({
        method,
        R2: values.bestErr?.R2 ?? -Infinity,
        MSE: values.bestErr?.MSE ?? Infinity,
        MAE: values.bestErr?.MAE ?? Infinity,
      }));
  
    if (!metrics.length) return null;
  
    // Нормализуем метрики
    const bestR2 = Math.max(...metrics.map(m => m.R2));
    const bestMSE = Math.min(...metrics.map(m => m.MSE));
    const bestMAE = Math.min(...metrics.map(m => m.MAE));
  
    // Устанавливаем веса для расчёта близости к идеалу
    const scored = metrics.map(m => {
      const r2Score = m.R2 / bestR2;
      const mseScore = bestMSE / m.MSE;
      const maeScore = bestMAE / m.MAE;
      const totalScore = r2Score + mseScore + maeScore;
      return { ...m, score: totalScore };
    });
  
    // Возвращаем метод с максимальным "totalScore"
    return scored.sort((a, b) => b.score - a.score)[0].method;
  };
  const bestMethod =
  responseMLRCSVData?.data
    ? getRecommendedMethod(responseMLRCSVData.data)
    : responseMLRData?.data
    ? getRecommendedMethod(responseMLRData.data)
    : null;


  return (
    <div className="App">
      <Header isRed={true}>
        <span>Regression analytics</span>
      </Header>

      <Tables tableData={tableData} setTableData={setTableData} headers={headers} setHeaders={setHeaders} rows={rows} setRows={setRows} cols={cols} setCols={setCols} />
      {/* <Graphics tableData={tableData} headers={responseMLRCSVData?.data?.names||headers} datapoints={responseMLRCSVData?.data?.datapoints || []} graphics={responseMLRCSVData?.data?.graphics || responseMLRData?.data?.graphics }/> */}
      <Buttons
    MLR={MLR}
    uploadFile={uploadFile}
    handleFileChange={handleFileChange}
    file={file}
    responses={{ responseMLRData, responseMLRCSVData }}
    headers={headers} 
  />
  {isLoading && (
    <div className="loading-dots">
      <span></span><span></span><span></span>
      <p>Моделируем эксперименты и определяем метрики...</p>
    </div>
  )}

{(responseMLRCSVData || responseMLRData) && (
  <>
    <h2>Выберите метод регрессии:</h2>
    <button 
      className="highlight-best-button" 
      onClick={() => setHighlightBestMethod(true)}
    >
      Показать лучший метод
    </button>
    <div className="results-grid">
      { /* Здесь твои результаты */ }
    </div>
  </>
)}


{(responseMLRCSVData || responseMLRData) && (
  <div className="results-grid">
    {responseMLRCSVData &&
      Object.keys(responseMLRCSVData.data || {})
        .filter(method => !['datapoints', 'graphics', 'coeff', 'data', 'names'].includes(method))
        .map((method) => (
          <Results
            key={method}
            title={`Анализ через ${method} регрессию`}
            data={responseMLRCSVData.data[method]}
            datapoints={responseMLRCSVData?.data?.datapoints}
            headers={responseMLRCSVData?.data?.names}
            recommended={highlightBestMethod && bestMethod === method}
          />
        ))
    }

    {responseMLRData &&
      Object.keys(responseMLRData.data || {})
        .filter(method => !['datapoints', 'graphics', 'coeff', 'data', 'names'].includes(method))
        .map((method) => {
          const { graphics, datapoints: _, ...filteredData } = responseMLRData.data[method] || {};
          return (
            <Results
              key={method}
              title={`Анализ через ${method} регрессию`}
              data={{ ...filteredData, graphics }}
              datapoints={responseMLRData?.data?.datapoints}
              headers={responseMLRCSVData?.data?.names}
              recommended={highlightBestMethod && bestMethod === method}
            />
          );
        })
    }
  </div>
)}


      <h2>📜 Документация по методам регрессии</h2>
      <div>
        <select className="select-method" value={selectedMethod} onChange={handleMethodChange}>
          <option  value="nil">Выбор метода регрессии</option>
          <option  value="MLR">MLR (Multiple Linear Regression)</option>
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
