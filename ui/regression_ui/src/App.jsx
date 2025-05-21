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
  const [analysisStep, setAnalysisStep] = useState("initial");

  React.useEffect(() => {
    if (analysisStep === "extended") {
      if (file) {
        uploadFile();
      } else {
        MLR();
      }
    }
  }, [analysisStep]);
  

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
  
    setHighlightBestMethod(false);
    setIsLoading(true);
    setTableData([]);
    setHeaders([]);
  
    const formData = new FormData();
    formData.append("file", file);
  
    let url = analysisStep === "initial"
      ? "/api/v1/onlymlrCSV"
      : "/api/v1/mlrCSV";
  
    if (file.name.endsWith(".xlsx") || file.name.endsWith(".xls")) {
      url = analysisStep === "initial"
        ? "/api/v1/onlymlrExcel"
        : "/api/v1/mlrExcel";
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
  
    setHighlightBestMethod(false);
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
  
    const url = analysisStep === "initial"
      ? "/api/v1/onlymlr"
      : "/api/v1/mlr";
  
    try {
      const response = await fetch(url, {
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
      .filter(([key, val]) => 
        !['datapoints', 'graphics'].includes(key) &&
        val?.bestErr
      )
      .map(([method, values]) => {
        const { R2 = -Infinity, MSE = Infinity, MAE = Infinity } = values.bestErr;
        return { method, R2, MSE, MAE };
      });
  
    if (!metrics.length) return null;
  
    const bestR2 = Math.max(...metrics.map(m => m.R2));
    const bestMSE = Math.min(...metrics.map(m => m.MSE));
    const bestMAE = Math.min(...metrics.map(m => m.MAE));
  
    const scored = metrics.map(m => {
      const r2Score = bestR2 ? m.R2 / bestR2 : 0;
      const mseScore = m.MSE ? bestMSE / m.MSE : 0;
      const maeScore = m.MAE ? bestMAE / m.MAE : 0;
      return {
        ...m,
        score: r2Score + mseScore + maeScore
      };
    });
  
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

{(responseMLRCSVData || responseMLRData) && analysisStep === "initial" && (
  
  <div className="continue-analysis">
  <h2>Можем продолжить анализ данных при помощи других методов</h2>
    <button 
    className="continue-button" 
    onClick={() => setAnalysisStep("extended")}
  >
    Продолжить анализ
  </button>
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
