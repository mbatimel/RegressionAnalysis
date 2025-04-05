import React from 'react';
import regressionDocs from './regressionDocs';

const Documentation = ({ selectedMethod }) => {
  const doc = regressionDocs[selectedMethod];

  if (!doc) return <p>Документация для метода не найдена.</p>;

  return (
    <div className="documentation-box">
      <h3>📖 {doc.title}</h3>
      <pre style={{ whiteSpace: 'pre-wrap', fontFamily: 'inherit' }}>
        {doc.content}
      </pre>
    </div>
  );
};

export default Documentation;
