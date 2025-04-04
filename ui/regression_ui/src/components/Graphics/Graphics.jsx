import React, { useState, useEffect } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { schemeCategory10 } from "d3-scale-chromatic";

const Graphics = ({ tableData = [], headers = [], datapoints = [], graphics = {} }) => {
  const [visibleGraphs, setVisibleGraphs] = useState({});

  useEffect(() => {
    const initialVisibility = {};
    headers.forEach((header) => {
      initialVisibility[header] = false;
    });
    Object.keys(graphics).forEach((key) => {
      initialVisibility[`Резулитат регрессии ${key}`] = false;
    });
    setVisibleGraphs(initialVisibility);
  }, [headers, graphics]);

  if (!tableData.length && !datapoints.length && !Object.keys(graphics).length) {
    return <p>⚠️ Нет данных для построения графика.</p>;
  }

  const datasets = datapoints.length
    ? headers.map((header, index) => ({
        name: header,
        data: datapoints
          .map((point) => ({ x: Number(point.vares[index]), y: Number(point.obs) }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y))
          .sort((a, b) => a.x - b.x),
        color: schemeCategory10[index % schemeCategory10.length],
      }))
    : headers.slice(1).map((header, index) => ({
        name: header,
        data: tableData
          .map((row) => ({ x: Number(row[index + 1]), y: Number(row[0]) }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y))
          .sort((a, b) => a.x - b.x),
        color: schemeCategory10[index % schemeCategory10.length],
      }));

  const graphicsData = Object.entries(graphics).map(([key, values]) => ({
    name: `График ${key}`,
    data: Object.entries(values)
      .map(([y, x]) => ({ x: parseFloat(x), y: parseFloat(y) }))
      .filter((point) => !isNaN(point.x) && !isNaN(point.y))
      .sort((a, b) => a.x - b.x),
    color: "red",
  }));

  const handleLegendClick = (graphName) => {
    setVisibleGraphs((prev) => ({
      ...prev,
      [graphName]: !prev[graphName],
    }));
  };

  const legendItems = [
    ...datasets.map(({ name, color }) => ({ name, color })),
    ...graphicsData.map(({ name }) => ({ name, color: "red" })),
  ];

  return (
    <div style={{ width: "100%", height: 500 }}>
      <h2>📊 График зависимостей Y от X</h2>
      <div style={{ display: "flex", flexWrap: "wrap", gap: "30px", marginBottom: "10px" }}>
        {legendItems.map(({ name, color }) => (
          <label key={name} style={{ cursor: "pointer", display: "flex", alignItems: "center" }}>
            <input type="checkbox" checked={visibleGraphs[name]} onChange={() => handleLegendClick(name)} />
            <span style={{ marginLeft: "10px", color }}>{name}</span>
          </label>
        ))}
      </div>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
          <CartesianGrid strokeDasharray="4 4" />
          <XAxis type="number" dataKey="x" name="X" />
          <YAxis type="number" dataKey="y" name="Y" />
          <Tooltip cursor={{ strokeDasharray: "3 3" }} />
          <Legend />
          {datasets.map(({ name, data, color }) =>
            visibleGraphs[name] ? <Line key={name} dataKey="y" data={data} stroke={color} name={name} /> : null
          )}
          {graphicsData.map(({ name, data, color }) =>
            visibleGraphs[name] ? <Line key={name} dataKey="y" data={data} stroke={color} name={name} /> : null
          )}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default Graphics;
