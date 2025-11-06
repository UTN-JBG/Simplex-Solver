import { jsPDF } from "jspdf";
import Tableaux from "./Tableaux";

export default function SimplexResult({ result }) {
  const optimalValue = result.optimal !== undefined && result.optimal !== null ? result.optimal.toFixed(2) : 'N/A';
  const variables = result.variables || {};
  const tableauxHistory = result.tableaux_history || [];  

  return (
    <div style={{ 
      marginTop: '20px', 
      padding: '20px', 
      border: '1px solid #00131a', 
      borderRadius: '8px',
      backgroundColor: '#00080e', 
      color: '#e6eef3' 
    }}>
      <h2 style={{ borderBottom: '2px solid #007bff', paddingBottom: '10px' , color: '#00d1ff'}}>Resultado</h2>
      
      <div style={{ marginBottom: '20px' }}>
        <p><strong>Tipo de problema:</strong> <span style={{ color: '#7dd3fc' }}>{result.type || 'N/A'}</span></p> 
        <h3 style={{ color: '#28a745' }}>Valor óptimo: <span style={{ color: '#fff' }}>{optimalValue}</span></h3>
        
        <h4>Variables:</h4>
        <ul style={{ listStyleType: 'disc', marginLeft: '20px' }}>
          {Object.entries(variables).map(([k, v]) => (
            <li key={k} style={{ marginBottom: '5px' }}>
              <strong>{k}:</strong> 
              <span style={{ color: '#fff' }}>{(typeof v === 'number') ? v.toFixed(2) : 'N/A'}</span>
            </li>
          ))}
        </ul>
      </div>

      <h3 style={{ marginTop: '30px', borderTop: '1px solid #0b1720', paddingTop: '20px', color: '#00d1ff' }}>
        Historial de Tablas (Pasos Intermedios)
      </h3>
      
      {tableauxHistory.length > 0 ? (
        tableauxHistory.map((t, idx) => (
          <Tableaux key={idx} tableau={t} index={idx} /> 
        ))
      ) : (
        <p style={{ color: '#9fb4c9' }}>No se generaron pasos intermedios (ej. solución inmediata).</p>
      )}

    </div>
  );
}
