import React, { useState } from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";
import 'katex/dist/katex.min.css';
import { BlockMath } from 'react-katex';

const FORMULAS = {
  "Анализ через Linear регрессию": "y = \\beta_0 + \\beta_1 x",
  "Анализ через Polynomial регрессию": "y = \\beta_0 + \\beta_1 x + \\beta_2 x^2 + \\dots + \\beta_n x^n",
  "Анализ через Ridge регрессию": "y = X \\beta + \\lambda ||\\beta||^2",
  "Анализ через Lasso регрессию": "y = X \\beta + \\lambda \\sum |\\beta_i|",
  "Анализ через Elastic регрессию": "y = X \\beta + \\lambda_1 ||\\beta||^2 + \\lambda_2 \\sum |\\beta_i|",
  "Анализ через SRV регрессию": "\\hat{y} = \\frac{1}{T} \\sum_{t=1}^T h_t(x)",
  "Анализ через Logistic регрессию": "\\hat{y} = \\sum_{m=1}^{M} \\gamma_m h_m(x)"
};

const Results = ({ title, data, datapoints = [], headers = [], recommended = false }) => {
  const [expanded, setExpanded] = useState(false);
  if (!data) return null;
  const { bestErr, graphics, ...restData } = data;


  const prepareGraphicsTables = (graphics) => {
    if (!graphics) return [];
    return Object.entries(graphics).map(([index, graphData]) => {
      const tableData = Object.entries(graphData).map(([y, x]) => ({
        x: parseFloat(x),
        y: parseFloat(y)
      }));
      const graphType = data.resultType?.[index] || `График ${parseInt(index) + 1}`;
      return {
        id: `graph-${index}`,
        tableData,
        graphData,
        title: graphType
      };
    });
  };

  const graphicsTables = prepareGraphicsTables(graphics);
  const latexFormula = FORMULAS[title];

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
          {/* Формула */}
          {latexFormula && (
            <div className="formula-block">
              <h4>Формула метода:</h4>
              <BlockMath math={latexFormula} />
            </div>
          )}

          {/* Остальные данные */}
          {Object.entries(restData).map(([key, value]) => (
            <div key={key} className="result-item">
              <strong>{key}:</strong>{" "}
              {typeof value === "object" ? JSON.stringify(value, null, 2) : value}
            </div>
          ))}

          {/* Графики и таблицы */}
          {graphicsTables.length > 0 && (
            <div className="graphics-section">
              <div className="graphics-grid">
                {graphicsTables.map((graph) => (
                  <div key={graph.id} className="graph-container">
                    <h4>{graph.title}</h4>
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
