import React from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";

const Results = ({ title, data, datapoints = [], headers = [] }) => {
    // Фильтруем поля, которые не должны отображаться пользователю
    const displayData = Object.entries(data).filter(
      ([key]) => !['graphics', 'datapoints'].includes(key)
    );
  
    return (
      <div className="results-container">
        <h2>{title}</h2>
        <div className="results-content">
          {displayData.map(([key, value]) => (
            <div key={key} className="result-item">
              <strong>{key}:</strong>{" "}
              {typeof value === "object" ? JSON.stringify(value, null, 2) : value}
            </div>
          ))}
        </div>
  
        {/* Эти поля не отображаются в основном контенте, но используются для графиков */}
        {data.graphics && (
          <div className="result-graph">
            <Graphics graphics={data.graphics} datapoints={datapoints} headers={headers} />
          </div>
        )}
      </div>
    );
  };
  
  export default Results;
