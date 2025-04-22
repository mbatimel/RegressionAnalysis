import React, { useState } from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";

const Results = ({ title, data, datapoints = [], headers = [] }) => {
  const [expanded, setExpanded] = useState(false);
  const { bestErr, graphics, ...restData } = data;

  const toggleExpand = () => setExpanded((prev) => !prev);

  return (
    <div className="results-container">
      <h2>{title}</h2>

      <div className="results-summary" onClick={toggleExpand} style={{ cursor: "pointer", backgroundColor: "#f3f3f3", padding: "10px", borderRadius: "8px" }}>
        <strong>Ключевые метрики (нажмите, чтобы {expanded ? "скрыть" : "раскрыть"}):</strong>
        <pre style={{ margin: 0 }}>{JSON.stringify(bestErr, null, 2)}</pre>
      </div>

      {expanded && (
        <div className="results-content" style={{ marginTop: "10px" }}>
          {Object.entries(restData).map(([key, value]) => (
            <div key={key} className="result-item">
              <strong>{key}:</strong>{" "}
              {typeof value === "object" ? <pre>{JSON.stringify(value, null, 2)}</pre> : value}
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
