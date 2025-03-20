import React, { useState, useEffect } from "react";
import {
  ScatterChart,
  Scatter,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Line,
  LineChart,
  Legend,
} from "recharts";
import { schemeCategory10 } from "d3-scale-chromatic";

const Graphics = ({
  tableData = [],
  headers = [],
  datapoints = [],
  graphics = {},
}) => {
  // Состояние видимости графиков
  const [visibleGraphs, setVisibleGraphs] = useState({});
  useEffect(() => {
    const initialVisibility = {};

    // Добавляем видимость для графиков из headers
    headers.forEach((header) => {
      initialVisibility[header] = false; // По умолчанию все графики из headers скрыты
    });

    // Добавляем видимость для графиков из graphics
    Object.keys(graphics).forEach((key) => {
      initialVisibility[`График ${key}`] = false; // По умолчанию все графики из graphics скрыты
    });

    setVisibleGraphs(initialVisibility);
  }, [headers, graphics]);
  
  if (
    !tableData.length &&
    !datapoints.length &&
    !Object.keys(graphics).length
  ) {
    return <p>⚠️ Нет данных для построения графика.</p>;
  }

  // Данные из файлов или таблицы
  const datasets = datapoints.length
    ? headers.slice(0).map((header, index) => ({
        name: header,
        data: datapoints
          .map((point) => ({
            x: Number(point.vares[index]),
            y: Number(point.obs),
            label: `${header}: (${point.vares[index]}, ${point.obs})`,
          }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y)),
        color: schemeCategory10[index % schemeCategory10.length],
      }))
    : headers.slice(1).map((header, index) => ({
        name: header,
        data: tableData
          .map((row) => ({
            x: Number(row[index + 1]),
            y: Number(row[0]),
            label: `${header}: (${row[index + 1]}, ${row[0]})`,
          }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y)),
        color: schemeCategory10[index % schemeCategory10.length],
      }));

  // Графики (линии) из graphics
  const graphicsData = Object.entries(graphics).map(([key, values], index) => {
    const data = Object.entries(values)
      .map(([y, x]) => ({
        x: parseFloat(x),
        y: parseFloat(y),
        label: `График ${key}: (${x}, ${y})`,
      }))
      .filter((point) => !isNaN(point.x) && !isNaN(point.y))
      .sort((a, b) => a.x - b.x);

    return {
      name: `График ${key}`,
      data,
      color: "red", // Все графики красного цвета
    };
  });
  // Обработчик клика по легенде (скрытие/показ графиков)
  const handleLegendClick = (graphName) => {
    setVisibleGraphs((prev) => ({
      ...prev,
      [graphName]: !prev[graphName],
    }));
  };
  
  // Собираем все элементы легенды (сначала таблица, потом graphics)
  const legendItems = [
    ...datasets.map(({ name, color }) => ({ name, color })),
    ...graphicsData.map(({ name }) => ({ name, color: "red" })), // Графики красного цвета
  ];
  
  return (
    <div style={{ width: "100%", height: 500 }}>
      <h2>📊 График зависимостей Y от X</h2>

      {/* Кастомная легенда */}
      <div
        style={{
          display: "flex",
          flexWrap: "wrap",
          gap: "15px",
          marginBottom: "10px",
        }}
      >
        {legendItems.map(({ name, color }) => (
          <label
            key={name}
            style={{ cursor: "pointer", display: "flex", alignItems: "center" }}
          >
            <input
              type="checkbox"
              checked={visibleGraphs[name]}
              onChange={() => handleLegendClick(name)}
            />
            <span style={{ marginLeft: "5px", color }}>{name}</span>
          </label>
        ))}
      </div>

      <ResponsiveContainer width="100%" height={400}>
        <ScatterChart margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
          <CartesianGrid />
          <XAxis type="number" dataKey="x" name="X" />
          <YAxis type="number" dataKey="y" name="Y" />
          <Tooltip
            cursor={{ strokeDasharray: "3 3" }}
            formatter={(value, name, props) => props.payload.label}
          />

          {/* Отображаем точки из таблицы */}
          {datasets.map(({ name, data, color }, i) =>
            visibleGraphs[name] ? (
              <Scatter key={i} name={name} data={data} fill={color} />
            ) : null
          )}

          {/* Отображаем графики из graphics */}
          {graphicsData.map(({ name, data, color }, i) =>
            visibleGraphs[name] ? (
              <Line
                key={i}
                type="linear"
                dataKey="y"
                data={data}
                stroke={color}
                dot={{ fill: { color }, r: 4 }}
                name={name}
                connectNulls={true}
              />
            ) : null
          )}
        </ScatterChart>
        <LineChart
          width={500}
          height={300}
          data={graphicsData}
          margin={{
            top: 5,
            right: 30,
            left: 20,
            bottom: 5,
          }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="name" />
          <YAxis />
          <Tooltip />
          <Legend />
          {graphicsData.map(({ name, data, color }, i) =>
            !visibleGraphs[name] ? (
              <Line
                key={i}
                type="linear"
                dataKey="y"
                data={data}
                stroke={color}
                dot={{ fill: {color}, r: 2 }}
                name={name}
                connectNulls={true}
              />
            ) : null
          )}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default Graphics;
