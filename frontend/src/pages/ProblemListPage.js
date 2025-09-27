import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import ProblemListItem from "./ProblemListItem";
import districts from "../data/almaty.json";

export default function ProblemListPage() {
  const { districtId } = useParams();
  const navigate = useNavigate();
  const [problems, setProblems] = useState([]);

  // ищем название района
  const district = districts.features.find(
    f => String(f.properties["osm-relation-id"]) === String(districtId)
  );
  const districtName = district?.properties?.nameRu || `ID ${districtId}`;

  useEffect(() => {
    fetch(`http://${process.env.REACT_APP_API_URL}/heatmap/districts/${districtId}/problems`)
      .then((res) => res.json())
      .then((data) => setProblems(data))
      .catch((err) => console.error("Ошибка загрузки проблем:", err));
  }, [districtId]);

  return (
    <div style={{ maxWidth: "900px", margin: "0 auto", padding: "20px" }}>
      <button
        onClick={() => navigate("/geomap/heatmap")}
        style={{
          marginBottom: "20px",
          padding: "8px 16px",
          borderRadius: "6px",
          border: "1px solid #ccc",
          background: "#f9f9f9",
          cursor: "pointer",
        }}
      >
        ← Назад к карте
      </button>

      <h2 style={{ marginBottom: "16px", fontWeight: 600 }}>
        Проблемы: {districtName}
      </h2>

      <div
        style={{
          marginTop: "10px",
          background: "#fff",
          borderRadius: "10px",
          boxShadow: "0 2px 8px rgba(0,0,0,0.05)",
          overflow: "hidden"
        }}
      >
        {problems.length === 0 ? (
          <div style={{ padding: "20px", textAlign: "center", color: "#777" }}>
            Нет проблем в этом районе
          </div>
        ) : (
          problems.map((p) => (
            <ProblemListItem
              key={p.problem_id}
              problem={p}
              districtId={districtId}
            />
          ))
        )}
      </div>
    </div>
  );
}
