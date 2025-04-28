import React, { useState } from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";

const Results = ({ title, data, datapoints = [], headers = [], recommended = false }) => {
  const [expanded, setExpanded] = useState(false);
  const { bestErr, graphics, ...restData } = data;

  // Функция для преобразования данных графиков в табличный формат
  const prepareGraphicsTables = (graphics) => {
    if (!graphics) return [];
    
    return Object.entries(graphics).map(([index, graphData]) => {
      // Преобразуем объект графика в массив точек {x, y}
      const tableData = Object.entries(graphData).map(([y, x]) => ({
        x: parseFloat(x),
        y: parseFloat(y)
      }));
      
      return {
        id: `graph-${index}`,
        tableData,
        graphData, // сохраняем оригинальные данные для графика
        title: `График ${parseInt(index) + 1}`
      };
    });
  };

  const graphicsTables = prepareGraphicsTables(graphics);

  return (
    <div className={`results-container ${expanded ? "expanded" : "collapsed"} ${recommended ? "recommended" : ""}`}>
      <div className="result-summary" onClick={() => setExpanded(!expanded)}>
        <div className="summary-header">
          <h3>{title}</h3>
          <div className="metrics">
            <span><strong>MAE:</strong> {bestErr?.MAE?.toFixed(2)}</span>
            <span><strong>MSE:</strong> {bestErr?.MSE?.toFixed(2)}</span>
            <span><strong>R²:</strong> {bestErr?.R2?.toFixed(4)}</span>
          </div>
          {recommended && (
            <div className="recommended-text">
              ✅ Предлагаем к вашему рассмотрению этот метод регрессии
            </div>
          )}
          <div className="toggle-details">
            {expanded ? "Скрыть детали ⬆" : "Показать результаты ⬇"}
          </div>
        </div>
      </div>

      {expanded && (
        <div className="results-content">
          {Object.entries(restData).map(([key, value]) => (
            <div key={key} className="result-item">
              <strong>{key}:</strong>{" "}
              {typeof value === "object" ? JSON.stringify(value, null, 2) : value}
            </div>
          ))}
          
          {graphicsTables.length > 0 && (
  <div className="graphics-section">
    <div className="graphics-grid">
      {graphicsTables.map((graph) => (
        <div key={graph.id} className="graph-container">
          <h4>{graph.title}</h4>

          {/* Таблица данных X и Y */}
          <div className="data-table-container">
            <table className="data-table">
              <thead>
                <tr>
                  <th>X (значение)</th>
                  <th>Y (предсказание)</th>
                </tr>
              </thead>
              <tbody>
                {graph.tableData.map((point, idx) => (
                  <tr key={idx}>
                    <td>{point.x.toFixed(4)}</td>
                    <td>{point.y.toFixed(4)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      ))}
    </div>
    {expanded && graphics && (
      <div className="result-graph">
        <Graphics graphics={graphics} datapoints={datapoints} headers={headers} />
      </div>
    )}
  </div>
)}

        </div>
      )}
    </div>
  );
};

export default Results;

