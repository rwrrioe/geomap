import React from "react";
import { useNavigate } from "react-router-dom";

const statusColors = {
  created: { bg: "#007bff", text: "Создано" },
  processing: { bg: "#ffc107", text: "В процессе" },
  solved: { bg: "#28a745", text: "Решено" },
};

export default function ProblemListItem({ problem, districtId }) {
  const navigate = useNavigate();
  const status =
    statusColors[problem.status] || { bg: "#6c757d", text: "Неизвестно" };

  return (
    <div
      onClick={() =>
        navigate(`/geomap/heatmap/districts/${districtId}/problems/${problem.problem_id}`)
      }
      style={{
        display: "flex",
        alignItems: "flex-start",
        padding: "14px 16px",
        borderBottom: "1px solid #eee",
        cursor: "pointer",
        transition: "background 0.2s",
      }}
      onMouseEnter={(e) => (e.currentTarget.style.background = "#f9f9f9")}
      onMouseLeave={(e) => (e.currentTarget.style.background = "#fff")}
    >
      {/* Фото */}
      <div
        style={{
          width: "110px",
          height: "110px",
          flexShrink: 0,
          background: "#f0f0f0",
          borderRadius: "8px",
          overflow: "hidden",
          marginRight: "16px",
        }}
      >
        {problem.image_url ? (
          <img
            src={problem.image_url}
            alt={problem.problem_name}
            style={{ width: "100%", height: "100%", objectFit: "cover" }}
          />
        ) : (
          <div
            style={{
              width: "100%",
              height: "100%",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              fontSize: "0.85rem",
              color: "#aaa",
            }}
          >
            Нет фото
          </div>
        )}
      </div>

      {/* Текст */}
      <div style={{ flex: 1 }}>
        <div style={{ fontWeight: "600", fontSize: "1.05rem", marginBottom: "4px" }}>
          {problem.problem_name}
        </div>

        <div
          style={{
            margin: "6px 0",
            display: "flex",
            alignItems: "center",
            gap: "10px",
          }}
        >
          <span
            style={{
              background: status.bg,
              color: "#fff",
              padding: "2px 8px",
              borderRadius: "12px",
              fontSize: "0.8rem",
              fontWeight: 500,
            }}
          >
            {status.text}
          </span>
          <span style={{ fontSize: "0.85rem", color: "#666" }}>
            ID: {problem.problem_id}
          </span>
        </div>

        <div style={{ fontSize: "0.9rem", color: "#444" }}>
          {problem.problem_desc?.length > 80
            ? problem.problem_desc.slice(0, 80) + "..."
            : problem.problem_desc || "Описание отсутствует"}
        </div>
      </div>
    </div>
  );
}
