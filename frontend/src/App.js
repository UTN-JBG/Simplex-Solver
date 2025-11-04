import jsPDF from "jspdf";
import autoTable from "jspdf-autotable";
import { useEffect, useState } from "react";

 function App() {
  const [numVars, setNumVars] = useState(2);
  const [numConstr, setNumConstr] = useState(2);
  const [objective, setObjective] = useState(Array(2).fill(""));
  const [constraints, setConstraints] = useState(
    Array.from({ length: 2 }, () => Array(2).fill(""))
  );
  const [rhs, setRHS] = useState(Array(2).fill(""));
  const [type, setType] = useState("max");
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setObjective((prev) => {
      const next = [...prev];
      next.length = numVars;
      for (let i = 0; i < numVars; i++) if (next[i] === undefined) next[i] = "";
      return next;
    });
    setConstraints((prev) => {
      const next = prev.map((r) => {
        const nr = [...r];
        nr.length = numVars;
        for (let i = 0; i < numVars; i++) if (nr[i] === undefined) nr[i] = "";
        return nr;
      });
      next.length = numConstr;
      for (let i = 0; i < numConstr; i++) if (!next[i]) next[i] = Array(numVars).fill("");
      return next;
    });
    setRHS((prev) => {
      const next = [...prev];
      next.length = numConstr;
      for (let i = 0; i < numConstr; i++) if (next[i] === undefined) next[i] = "";
      return next;
    });
  }, [numVars, numConstr]);

  const styles = {
    page: { minHeight: "100vh", background: "#000", color: "#fff", fontFamily: "Segoe UI, Roboto, sans-serif", display: "flex", alignItems: "center", justifyContent: "center", padding: 20 },
    card: { width: "100%", maxWidth: 960, background: "#071426", borderRadius: 12, padding: 24, boxShadow: "0 6px 30px rgba(0,0,0,0.6)" },
    row: { display: "flex", gap: 12, alignItems: "center", marginBottom: 12 },
    label: { minWidth: 140, color: "#9fb4c9" },
    input: { padding: "8px 10px", borderRadius: 8, border: "1px solid rgba(255,255,255,0.06)", background: "#0b1720", color: "#e6eef3", width: "100%" },
    smallInput: { width: 80, padding: "6px 8px", borderRadius: 8, border: "1px solid rgba(255,255,255,0.06)", background: "#0b1720", color: "#e6eef3" },
    button: { padding: "10px 16px", background: "linear-gradient(90deg,#00d1ff,#0066ff)", border: "none", color: "#001217", fontWeight: 700, borderRadius: 10, cursor: "pointer" },
    matrix: { display: "grid", gap: 8, marginTop: 8 },
    matrixRow: { display: "flex", gap: 8 },
    previewJson: { marginTop: 12, padding: 12, background: "#00131a", borderRadius: 8, fontSize: 13, color: "#bfeafc" },
  };

  const setObjectiveAt = (i, v) => setObjective((s) => { const n = [...s]; n[i] = v; return n; });
  const setConstraintAt = (r, c, v) => setConstraints((s) => { const n = s.map((row) => [...row]); n[r][c] = v; return n; });
  const setRHSAt = (i, v) => setRHS((s) => { const n = [...s]; n[i] = v; return n; });

  const buildPayload = () => {
    const obj = objective.map((v) => Number(String(v).trim()));
    const A = constraints.map((row) => row.map((v) => Number(String(v).trim())));
    const b = rhs.map((v) => Number(String(v).trim()));
    return { objective: obj, constraints: A, rhs: b, type };
  };

  const handleSubmit = async () => {
    setLoading(true);
    try {
      const payload = buildPayload();
      console.log("Payload:", payload);
      const res = await fetch("http://localhost:8080/api/simplex", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.error || res.statusText || "Error en la solicitud");
      }
      const data = await res.json();
      setResult(data.result);
    } catch (err) {
      alert("Error: " + (err.message || err));
    } finally {
      setLoading(false);
    }
  };

  const handleDownloadPDF = () => {
    if (!result) return;
    const doc = new jsPDF({ unit: "pt", format: "a4" });
    let y = 40;
    doc.setFontSize(18);
    doc.text("Simplex Result", 40, y); y += 22;
    doc.setFontSize(12);
    doc.text(`Tipo: ${type}`, 40, y); y += 16;
    doc.text(`Valor óptimo: ${String(result.optimal)}`, 40, y); y += 20;

    const vars = Object.entries(result.variables || {}).map(([k, v]) => [k, String(v)]);
    if (vars.length) autoTable(doc,{ head: [["Variable", "Valor"]], body: vars, startY: y, theme: "grid", headStyles: { fillColor: [6, 86, 115] } });

    (result.tableaux_history || []).forEach((t, idx) => {
      const matrix = t.matrix || [];
      const start = doc.lastAutoTable ? doc.lastAutoTable.finalY + 12 : 200;
      doc.text(`Tabla ${idx + 1}`, 40, start);
      const head = Array.from({ length: matrix[0]?.length || 0 }, (_, i) => `C${i + 1}`);
      const body = matrix.map((r) => r.map((c) => String(c)));
      autoTable(doc,{ head: [head], body, startY: start + 6, theme: "grid", styles: { fontSize: 9 } });
    });

    doc.save("simplex_result.pdf");
  };

  return (
    <div style={styles.page}>
      <div style={styles.card}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 12 }}>
          <h2 style={{ margin: 0 }}>Simplex Solver</h2>
          <div style={{ color: "#7dd3fc", fontWeight: 600 }}>{type === "max" ? "Maximizar" : "Minimizar"}</div>
        </div>

        <div style={styles.row}>
          <div style={styles.label}>Variables:</div>
          <select style={styles.smallInput} value={String(numVars)} onChange={(e) => setNumVars(Number(e.target.value))}>
            {Array.from({ length: 10 }, (_, i) => i + 1).map((n) => <option key={n} value={n}>{n}</option>)}
          </select>

          <div style={{ minWidth: 120 }} />
          <div style={styles.label}>Restricciones:</div>
          <select style={styles.smallInput} value={String(numConstr)} onChange={(e) => setNumConstr(Number(e.target.value))}>
            {Array.from({ length: 10 }, (_, i) => i + 1).map((n) => <option key={n} value={n}>{n}</option>)}
          </select>

          <div style={{ marginLeft: "auto", display: "flex", gap: 8, alignItems: "center" }}>
            <select style={{ ...styles.smallInput, width: 120 }} value={type} onChange={(e) => setType(e.target.value)}>
              <option value="max">Max</option>
              <option value="min">Min</option>
            </select>
            <button style={styles.button} onClick={handleSubmit} disabled={loading}>{loading ? "Resolviendo..." : "Resolver"}</button>
          </div>
        </div>

        <div style={{ marginTop: 8 }}>
          <div style={{ marginBottom: 8, color: "#9fb4c9" }}>Objetivo (coeficientes)</div>
          <div style={{ display: "flex", gap: 8 }}>
            {objective.map((val, i) => (
              <input
                key={i}
                style={styles.smallInput}
                value={val}
                onChange={(e) => setObjectiveAt(i, e.target.value)}
                placeholder={`c${i + 1}`}
              />
            ))}
          </div>
        </div>

        <div style={{ marginTop: 14 }}>
          <div style={{ marginBottom: 8, color: "#9fb4c9" }}>Matriz de restricciones (fila × columna)</div>
          <div style={styles.matrix}>
            {constraints.map((row, rIdx) => (
              <div key={rIdx} style={styles.matrixRow}>
                <div style={{ width: 60, color: "#9fb4c9", display: "flex", alignItems: "center" }}>{`R${rIdx + 1}`}</div>
                {row.map((cell, cIdx) => (
                  <input
                    key={cIdx}
                    style={styles.smallInput}
                    value={cell}
                    onChange={(e) => setConstraintAt(rIdx, cIdx, e.target.value)}
                    placeholder={`a${rIdx + 1}${cIdx + 1}`}
                  />
                ))}
                <div style={{ width: 40, textAlign: "center", color: "#9fb4c9" }}>=</div>
                <input
                  style={{ ...styles.smallInput, width: 120 }}
                  value={rhs[rIdx]}
                  onChange={(e) => setRHSAt(rIdx, e.target.value)}
                  placeholder={`b${rIdx + 1}`}
                />
              </div>
            ))}
          </div>
        </div>
        {result && (
          <div style={{ marginTop: 16, padding: 12, background: "#00131a", borderRadius: 8 }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <strong>Resultado</strong>
              <button style={{ ...styles.button, padding: "8px 12px" }} onClick={handleDownloadPDF}>Descargar PDF</button>
            </div>
            <div style={{ marginTop: 8 }}>
              <div>Valor óptimo: <strong>{String(result.optimal)}</strong></div>
              <div style={{ marginTop: 8 }}>
                <strong>Variables</strong>
                <ul>
                  {Object.entries(result.variables || {}).map(([k, v]) => <li key={k}>{k}: {String(v)}</li>)}
                </ul>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;