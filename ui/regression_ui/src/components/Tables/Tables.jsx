import React from "react";

const Tables = ({
  tableData,
  setTableData,
  headers,
  setHeaders,
  rows,
  setRows,
  cols,
  setCols,
}) => {
  const handleRowsChange = (e) => setRows(Number(e.target.value) || 1);
  const handleColsChange = (e) => setCols(Number(e.target.value) || 1);

  const updateTableData = (rowIndex, colIndex, value) => {
    const newData = [...tableData];
    if (!newData[rowIndex]) newData[rowIndex] = [];
    newData[rowIndex][colIndex] = value;
    setTableData(newData);
  };

  return (
    <div>
      <div>
        <label>
          Количество наблюдений:
          <input type="number" value={rows} onChange={handleRowsChange} />
        </label>
      </div>
      <div>
        <label>
          Количество показателей:
          <input type="number" value={cols} onChange={handleColsChange} />
        </label>
      </div>
      <table border="1">
        <thead>
          <tr>
            {[
              "Y",
              ...Array(cols - 1)
                .fill(0)
                .map((_, i) => `X${i + 1}`),
            ].map((label, index) => (
              <th key={index}>
                <input
                  type="text"
                  placeholder={label}
                  value={headers[index] || ""}
                  onChange={(e) => {
                    const newHeaders = [...headers];
                    newHeaders[index] = e.target.value;
                    setHeaders(newHeaders);
                  }}
                />
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: rows }).map((_, rowIndex) => (
            <tr key={rowIndex}>
              {Array.from({ length: cols }).map((_, colIndex) => (
                <td key={colIndex}>
                  <input
                    type="text"
                    value={tableData[rowIndex]?.[colIndex] || ""}
                    onChange={(e) =>
                      updateTableData(rowIndex, colIndex, e.target.value)
                    }
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Tables;
