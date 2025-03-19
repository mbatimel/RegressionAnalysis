import React from "react";
import { ScatterChart, Scatter, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from "recharts";
import { schemeCategory10 } from "d3-scale-chromatic"; // Палитра цветов

const Graphics = ({ tableData = [], headers = [], datapoints = [] }) => {
  if (!tableData.length && !datapoints.length) {
    return <p>⚠️ Нет данных для построения графика.</p>;
  }

  // Если есть данные из файла, используем их
  const datasets = datapoints.length
    ? headers.slice(0).map((header, index) => ({
        name: header,
        data: datapoints
          .map((point) => ({
            x: Number(point.vares[index]), // X1, X2, X3, ...
            y: Number(point.obs), // Y
            label: `${header}: (${point.vares[index]}, ${point.obs})`, // Подпись точки
          }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y)), // Убираем NaN
        color: schemeCategory10[index % schemeCategory10.length], // Цвета из палитры
      }))
    : headers.slice(1).map((header, index) => ({
        name: header,
        data: tableData
          .map((row) => ({
            x: Number(row[index + 1]), // X1, X2, X3, ...
            y: Number(row[0]), // Y
            label: `${header}: (${row[index + 1]}, ${row[0]})`, // Подпись точки
          }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y)), // Убираем NaN
        color: schemeCategory10[index % schemeCategory10.length], // Цвета из палитры
      }));

  return (
    <div style={{ width: "100%", height: 450 }}>
      <h2>📊 График зависимостей Y от X</h2>
      <ResponsiveContainer width="100%" height={400}>
        <ScatterChart margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
          <CartesianGrid />
          <XAxis type="number" dataKey="x" name="X" />
          <YAxis type="number" dataKey="y" name="Y" />
          <Tooltip cursor={{ strokeDasharray: "3 3" }} formatter={(value, name, props) => props.payload.label} />
          <Legend />
          {datasets.map(({ name, data, color }, i) => (
            <Scatter key={i} name={name} data={data} fill={color} />
          ))}
        </ScatterChart>
      </ResponsiveContainer>
    </div>
  );
};

export default Graphics;