import { jsPDF } from "jspdf";
import Tableaux from "./Tableaux";

export default function SimplexResult({ result }) {
  const generatePDF = () => {
    const doc = new jsPDF();
    doc.text(`Valor óptimo: ${result.optimal}`, 10, 10);
    Object.entries(result.variables).forEach(([k,v], i) => {
      doc.text(`${k}: ${v}`, 10, 20 + i*10);
    });
    doc.save("simplex_result.pdf");
  };

  return (
    <div>
      <h2>Valor óptimo: {result.optimal}</h2>
      <h3>Variables:</h3>
      <ul>
        {Object.entries(result.variables).map(([k,v]) => <li key={k}>{k}: {v}</li>)}
      </ul>

      <h3>Tableaux:</h3>
      {result.tableaux_history.map((t, idx) => (
        <Tableaux key={idx} tableau={t} />
      ))}

      <button onClick={generatePDF}>Generar PDF</button>
    </div>
  );
}
