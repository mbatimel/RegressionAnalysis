import React from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";

const Results = ({ title, data, datapoints = [], headers = []}) => {
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

      {data.graphics && (
        <div className="result-graph">
          <Graphics graphics={data.graphics} datapoints={datapoints} headers={headers}/>
        </div>
      )}
    </div>
  );
};

export default Results;