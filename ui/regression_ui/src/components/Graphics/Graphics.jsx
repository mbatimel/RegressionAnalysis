import React, { useEffect, useState } from "react";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  LinearScale,
  Title,
  Tooltip,
  Legend,
  Filler,
  CategoryScale,
} from "chart.js";
import zoomPlugin from "chartjs-plugin-zoom";
import annotationPlugin from "chartjs-plugin-annotation";
ChartJS.register(LineElement, PointElement, LinearScale, Title, Tooltip, Legend, Filler, CategoryScale, zoomPlugin,annotationPlugin);

const Graphics = ({ tableData = [], headers = [], datapoints = [], graphics = {} }) => {
  const [visibleGraphs, setVisibleGraphs] = useState({});

  useEffect(() => {
    const initial = {};
    headers.forEach((header) => (initial[header] = true));
    Object.keys(graphics).forEach((key) => {
      initial[`Резулитат регрессии ${key}`] = true;
    });
    setVisibleGraphs(initial);
  }, [headers, graphics]);

  if (!tableData.length && !datapoints.length && !Object.keys(graphics).length) {
    return <p>⚠️ Нет данных для построения графика.</p>;
  }

  const datasets = datapoints.length
    ? headers.map((header, index) => ({
        label: header,
        data: datapoints
          .map((point) => ({ x: Number(point.vares[index]), y: Number(point.obs) }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y))
          .sort((a, b) => a.x - b.x),
        borderColor: `hsl(${index * 40}, 70%, 50%)`,
        tension: 0.3,
        hidden: !visibleGraphs[header],
      }))
    : headers.slice(1).map((header, index) => ({
        label: header,
        data: tableData
          .map((row) => ({ x: Number(row[index + 1]), y: Number(row[0]) }))
          .filter((point) => !isNaN(point.x) && !isNaN(point.y))
          .sort((a, b) => a.x - b.x),
        borderColor: `hsl(${index * 40}, 70%, 50%)`,
        tension: 0.3,
        hidden: !visibleGraphs[header],
      }));

  const graphicsData = Object.entries(graphics).map(([key, values], index) => ({
    label: `Резулитат регрессии ${key}`,
    data: Object.entries(values)
      .map(([y, x]) => ({ x: parseFloat(x), y: parseFloat(y) }))
      .filter((point) => !isNaN(point.x) && !isNaN(point.y))
      .sort((a, b) => a.x - b.x),
    borderColor: "red",
    borderWidth: 2,
    tension: 0,
    hidden: !visibleGraphs[`Резулитат регрессии ${key}`],
  }));

  const chartData = {
    datasets: [...datasets, ...graphicsData],
  };

  const options = {
    responsive: true,
    plugins: {
      legend: { display: true },
      tooltip: { mode: "nearest", intersect: false },
      zoom: {
        zoom: {
          wheel: { enabled: true },
          pinch: { enabled: true },
          mode: "xy",
        },
        pan: {
          enabled: true,
          mode: "xy",
        },
      },
      annotation: {
        annotations: {
          xZeroLine: {
            type: "line",
            xMin: 0,
            xMax: 0,
            borderColor: "blue",
            borderWidth: 2,
            label: {
              enabled: true,
              content: "x = 0",
              position: "start",
              backgroundColor: "rgba(0,0,255,0.1)",
              color: "blue",
            },
          },
          yZeroLine: {
            type: "line",
            yMin: 0,
            yMax: 0,
            borderColor: "green",
            borderWidth: 2,
            label: {
              enabled: true,
              content: "y = 0",
              position: "start",
              backgroundColor: "rgba(0,128,0,0.1)",
              color: "green",
            },
          },
        },
      },
    },
    scales: {
      x: {
        type: "linear",
        title: { display: true, text: "X" },
      },
      y: {
        type: "linear",
        title: { display: true, text: "Y" },
      },
    },
  };
  
  

  const handleLegendClick = (key) => {
    setVisibleGraphs((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  const legendItems = Object.keys(visibleGraphs);

  return (
    <div>
      <h2>📊 График зависимостей Y от X</h2>
      <div style={{ display: "flex", flexWrap: "wrap", gap: "15px", marginBottom: "10px" }}>
        {legendItems.map((key) => (
          <label key={key}>
            <input type="checkbox" checked={visibleGraphs[key]} onChange={() => handleLegendClick(key)} />
            <span style={{ marginLeft: 8 }}>{key}</span>
          </label>
        ))}
      </div>
      <Line data={chartData} options={options} />
    </div>
  );
};

export default Graphics;
