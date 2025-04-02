import React from "react";

const Buttons = ({ MLR, uploadFile, handleFileChange, file, responses, headers }) => {
  return (
    <div>
      <h1>React cURL Buttons</h1>
      <button className="MLRButton" onClick={MLR}>MLR</button>
      {responses.responseMLRData && (
        <div className="result-box">
          <h3>Формула:</h3>
          <p>Y = {headers[0]}</p>
          {headers.slice(1).map((name, i) => <p key={i}>X{i + 1} = {name}</p>)}
          <h3>Ответ от сервера:</h3>
          <pre>{JSON.stringify(responses.responseMLRData, null, 2)}</pre>
        </div>
      )}
      <h1>React File Upload</h1>
      <input type="file" onChange={handleFileChange} />
      <button className="UploadButton" onClick={uploadFile}>Отправить файл</button>
      {responses.responseMLRCSVData && (
        <div className="result-box">
          <h3>Ответ от сервера (MLR CSV):</h3>
          <pre>{JSON.stringify(responses.responseMLRCSVData, null, 2)}</pre>
        </div>
      )}
    </div>
  );
};

export default Buttons;
