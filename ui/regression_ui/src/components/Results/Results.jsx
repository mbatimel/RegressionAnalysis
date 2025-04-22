import React, { useState } from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";

const Results = ({ title, data, datapoints = [], headers = [], recommended = false }) => {
  const [expanded, setExpanded] = useState(false);
  const { bestErr, graphics, ...restData } = data;

  return (
    <div className="results-container">
      <div
        className="result-summary"
        onClick={() => setExpanded(!expanded)}
        style={{ cursor: "pointer", backgroundColor: recommended ? "#e6f7ff" : "#f5f5f5", padding: "10px", borderRadius: "8px", border: "1px solid #ccc", marginBottom: "10px" }}
      >
        <h3 style={{ marginBottom: "5px" }}>{title}</h3>
        <div>
          <strong>MAE:</strong> {bestErr?.MAE?.toFixed(2)} | <strong>MSE:</strong> {bestErr?.MSE?.toFixed(2)} | <strong>R²:</strong> {bestErr?.R2?.toFixed(4)}
        </div>
        {recommended && (
          <div style={{ color: "#1890ff", fontWeight: "bold", marginTop: "5px" }}>
            ✅ Предлагаем к вашему рассмотрению этот метод регрессии
          </div>
        )}
        <div style={{ color: "#888", fontSize: "0.9em" }}>
          {expanded ? "Скрыть детали ⬆" : "Показать результаты ⬇"}
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
