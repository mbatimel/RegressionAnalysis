import React from "react";
import "./Results.css";

const Results = ({ title, data }) => {
  return (
    <div className="results-container">
      <h2>{title}</h2>
      <div className="results-content">
        {Object.entries(data).map(([key, value]) => (
          <div key={key} className="result-item">
            <strong>{key}:</strong> {typeof value === "object" ? JSON.stringify(value, null, 2) : value}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Results;
