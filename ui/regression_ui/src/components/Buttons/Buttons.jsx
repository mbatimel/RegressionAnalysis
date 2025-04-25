import React from "react";
import "./Buttons.css";

const Buttons = ({ MLR, uploadFile, handleFileChange, file }) => {
  return (
    <div className="buttons-container">
      <div className="custom-file-upload">
        <label htmlFor="fileInput" className="custom-file-label">
          📁 Выбрать файл
        </label>
        <input
          id="fileInput"
          type="file"
          onChange={handleFileChange}
          className="file-input-hidden"
        />
        {file && <span className="file-name">{file.name}</span>}
      </div>

      <button className="MLRButton" onClick={MLR}>Анализировать таблицу</button>
      <button className="UploadButton" onClick={uploadFile}>Анализировать файл</button>
    </div>
  );
};

export default Buttons;
