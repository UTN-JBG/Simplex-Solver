export default function Tableaux({ tableau }) {
  return (
    <div style={{ marginBottom: "20px" }}>
      <h4>Tabla</h4>
      <table border="1">
        <thead>
          <tr>
            {tableau.headers.map(h => <th key={h}>{h}</th>)}
          </tr>
        </thead>
        <tbody>
          {tableau.matrix.map((row, i) => (
            <tr key={i}>
              {row.map((val, j) => <td key={j}>{val.toFixed(2)}</td>)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
