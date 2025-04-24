import React, { useState } from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";

const Results = ({ title, data, datapoints = [], headers = [], recommended = false }) => {
  const [expanded, setExpanded] = useState(false);
  const { bestErr, graphics, ...restData } = data;

  return (
    <div className={`results-container ${recommended ? "recommended" : ""}`}>
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
          {graphics && (
            <div className="result-graph">
              <Graphics graphics={graphics} datapoints={datapoints} headers={headers} />
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default Results;
