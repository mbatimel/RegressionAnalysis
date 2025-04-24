import React from "react";
import "./Buttons.css";

const Buttons = ({ MLR, uploadFile, handleFileChange, file }) => {
  return (
    <div className="buttons-container">
      <input type="file" onChange={handleFileChange} className="file-input" />
      <button className="MLRButton" onClick={MLR}>Анализировать таблицу</button>
      <button className="UploadButton" onClick={uploadFile}>Анализировать файл</button>
    </div>
  );
};

export default Buttons;
