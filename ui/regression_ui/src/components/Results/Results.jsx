import React, { useState } from "react";
import "./Results.css";
import Graphics from "./../Graphics/Graphics";
import 'katex/dist/katex.min.css';
import { BlockMath } from 'react-katex';

// Формулы с генерацией на основе параметров
const generateFormula = (title, params = {}) => {
  const { Coef = [], intercept = 0, degree = 2, lambda, lambda1, lambda2 } = params;

  switch (title) {
    case "Анализ через Linear регрессию":
      return `y = ${intercept} + ${Coef.map((c, i) => `${c} x_{${i + 1}}`).join(" + ")}`;

    case "Анализ через Polynomial регрессию": {
      let parts = [`${intercept}`];
      let idx = 0;
      for (let i = 0; i < Coef.length / degree; i++) {
        for (let d = 1; d <= degree; d++) {
          parts.push(`${Coef[idx++]} x_{${i + 1}}^{${d}}`);
        }
      }
      return `y = ${parts.join(" + ")}`;
    }

    case "Анализ через Ridge регрессию":
      return `y = X \\beta + ${lambda ?? '\\lambda'} \\|\\beta\\|^2`;

    case "Анализ через Lasso регрессию":
      return `y = X \\beta + ${lambda ?? '\\lambda'} \\sum |\\beta_i|`;

    case "Анализ через Elastic регрессию":
      return `y = X \\beta + ${lambda1 ?? '\\lambda_1'} \\|\\beta\\|^2 + ${lambda2 ?? '\\lambda_2'} \\sum |\\beta_i|`;

    case "Анализ через SRV регрессию":
      return "\\hat{y}(x) = \\sum_{i=1}^l (\\alpha_i - \\alpha_i^*) K(x_i, x) + b";

    case "Анализ через Logistic регрессию":
      return `P(y=1|x) = \\frac{1}{1 + e^{-(${intercept} + ${Coef.map((c, i) => `${c} x_{${i + 1}}`).join(" + ")})}}`;

    case "Анализ через LogChecking регрессию":
      return `y = ${intercept} + ${Coef[0]} \\ln(x)`;

    default:
      return null;
  }
};

const Results = ({ title, data, datapoints = [], headers = [], recommended = false }) => {
  const [expanded, setExpanded] = useState(false);
  if (!data) return null;
  const { bestErr, graphics, params } = data;

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
  const latexFormula = generateFormula(title, params);

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
