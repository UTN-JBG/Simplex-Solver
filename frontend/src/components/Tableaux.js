export default function Tableaux({ tableau, index}) {
  const headers = tableau.headers || []; 
  const matrix = tableau.matrix || [];
  return (
    <div style={{ marginBottom: "20px" }}>
      <h4>Tabla {index + 1}</h4>
      <table border="1">
        <thead>
          <tr>
            {headers.map((h, i) => <th key={i}>{h}</th>)}
          </tr>
        </thead>
        <tbody>
          {matrix.map((row, i) => (
            <tr key={i}>
              {row.map((val, j) => <td key={j}>{val.toFixed(2)}</td>)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
