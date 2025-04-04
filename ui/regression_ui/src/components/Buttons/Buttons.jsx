import React from "react";

const Buttons = ({ MLR, uploadFile, handleFileChange, file, responses, headers }) => {
  return (
    <div>
      <h1>React cURL Buttons</h1>
      <button className="MLRButton" onClick={MLR}>Приступить к анализу</button>

      <h1>React File Upload</h1>
      <input type="file" onChange={handleFileChange} />
      <button className="UploadButton" onClick={uploadFile}>Приступить к анализу файла</button>
    </div>
  );
};

export default Buttons;
