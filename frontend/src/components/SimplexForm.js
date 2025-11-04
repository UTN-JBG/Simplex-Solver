import { useState } from "react";

export default function SimplexForm({ onResult }) {
  console.log({
          objective: objective.split(",").map(Number),
          constraints: constraints.split(";").map(row => row.split(",").map(Number)),
          rhs: rhs.split(",").map(Number),
          type: type,
      });

  const [objective, setObjective] = useState("");
  const [constraints, setConstraints] = useState("");
  const [rhs, setRHS] = useState("");
  const [type, setType] = useState("max"); // default: max

  const handleSubmit = async (e) => {
    e.preventDefault();

    const response = await fetch("http://localhost:8080/api/simplex", {

      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        objective: objective.split(",").map(Number),
        constraints: constraints.split(";").map(row => row.split(",").map(Number)),
        rhs: rhs.split(",").map(Number),
        type: type,
      }),
    });

    if (!response.ok) {
      alert("Error en la solicitud: " + response.statusText);
      return;
    }

    const data = await response.json();
    onResult(data.result);
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Simplex Solver</h2>

      <label>
        Objective (comma-separated):
        <input
          type="text"
          value={objective}
          onChange={(e) => setObjective(e.target.value)}
          placeholder="25,22"
        />
      </label>
      <br />

      <label>
        Constraints (rows separated by ';', values by ','):
        <input
          type="text"
          value={constraints}
          onChange={(e) => setConstraints(e.target.value)}
          placeholder="0.45,0.35;0.18,0.36;0.3,0.2"
        />
      </label>
      <br />

      <label>
        RHS (comma-separated):
        <input
          type="text"
          value={rhs}
          onChange={(e) => setRHS(e.target.value)}
          placeholder="1260000,900000,300000"
        />
      </label>
      <br />

      <label>
        Type:
        <select value={type} onChange={(e) => setType(e.target.value)}>
          <option value="max">Max</option>
          <option value="min">Min</option>
        </select>
      </label>
      <br />

      <button type="submit">Resolver Simplex</button>
    </form>
  );
}
