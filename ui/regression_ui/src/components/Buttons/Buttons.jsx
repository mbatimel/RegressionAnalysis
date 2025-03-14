import React from "react";

const Buttons = ({ MLR, Ridge, Lasso, ElasticNet, uploadFile, handleFileChange, file, responses, headers }) => {
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

      <button className="RidgeButton" onClick={Ridge}>Ridge</button>
      {responses.responseRidgeData && (
        <div className="result-box">
          <h3>Ответ от сервера (Ridge):</h3>
          <pre>{JSON.stringify(responses.responseRidgeData, null, 2)}</pre>
        </div>
      )}

      <button className="LassoButton" onClick={Lasso}>Lasso</button>
      {responses.lassoResponse && (
        <div className="result-box">
          <h3>Ответ от сервера (Lasso):</h3>
          <pre>{JSON.stringify(responses.lassoResponse, null, 2)}</pre>
        </div>
      )}

      <button className="ElastNetButton" onClick={ElasticNet}>Elastic</button>
      {responses.elasticnetResponse && (
        <div className="result-box">
          <h3>Ответ от сервера (ElasticNet):</h3>
          <pre>{JSON.stringify(responses.elasticnetResponse, null, 2)}</pre>
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
