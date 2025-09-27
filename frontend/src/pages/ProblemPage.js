import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getProblem } from "../api";
import { Container, Spinner } from "react-bootstrap";

const STATUS_MAP = {
  created: { text: "Создано", bg: "#007bff" },
  processing: { text: "В процессе", bg: "#ffc107" },
  solved: { text: "Решено", bg: "#28a745" },
};

const TYPE_MAP = {
  1: "ЖКХ",
  2: "Дороги и транспорт",
  3: "Гос.сервис",
  4: "Прочее",
};

export default function ProblemPage() {
  const { districtId, problemId } = useParams();
  const navigate = useNavigate();
  const [problem, setProblem] = useState(null);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState(null);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      try {
        const res = await getProblem(districtId, problemId);
        if (!mounted) return;
        setProblem(res);
      } catch (e) {
        console.error(e);
        setErr(e);
      } finally {
        setLoading(false);
      }
    };
    load();
    return () => {
      mounted = false;
    };
  }, [districtId, problemId]);

  if (loading) {
    return (
      <Container style={{ paddingTop: 40, textAlign: "center" }}>
        <Spinner animation="border" />
        <div style={{ marginTop: 12 }}>Загрузка...</div>
      </Container>
    );
  }

  if (err) {
    return (
      <Container style={{ paddingTop: 40 }}>
        <div className="alert alert-danger">
          Ошибка загрузки: {err.message || String(err)}
        </div>
      </Container>
    );
  }

  if (!problem) {
    return (
      <Container style={{ paddingTop: 40 }}>
        <div className="alert alert-warning">Проблема не найдена.</div>
      </Container>
    );
  }

  const status = STATUS_MAP[problem.status] || {
    text: problem.status || "Неизвестно",
    bg: "#6c757d",
  };

  const typeText =
    TYPE_MAP[problem.problem_typeid] ||
    TYPE_MAP[problem.TypeID] ||
    problem.problem_typeid ||
    "-";

  const lat = problem.point?.lat || problem.geom?.lat || null;
  const lon = problem.point?.lon || problem.geom?.lon || null;

  return (
    <Container style={{ maxWidth: "900px", padding: "20px" }}>
      <button
        onClick={() => navigate(`/geomap/heatmap/districts/${districtId}/problems`)}
        style={{
          marginBottom: "20px",
          padding: "8px 16px",
          borderRadius: "6px",
          border: "1px solid #ccc",
          background: "#f9f9f9",
          cursor: "pointer",
        }}
      >
        ← Назад к списку
      </button>

      <div
        style={{
          background: "#fff",
          borderRadius: "10px",
          boxShadow: "0 2px 8px rgba(0,0,0,0.05)",
          overflow: "hidden",
        }}
      >
        {/* Фото */}
        <div style={{ width: "100%", height: "400px", background: "#f0f0f0" }}>
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
                color: "#aaa",
              }}
            >
              Нет фото
            </div>
          )}
        </div>

        {/* Текст */}
        <div style={{ padding: "16px 20px" }}>
          <h3 style={{ margin: "0 0 10px", fontWeight: 600 }}>
            {problem.problem_name || `Проблема #${problem.problem_id}`}
          </h3>

          <div style={{ margin: "6px 0 14px", display: "flex", gap: "10px" }}>
            <span
              style={{
                background: status.bg,
                color: "#fff",
                padding: "4px 10px",
                borderRadius: "12px",
                fontSize: "0.85rem",
                fontWeight: 500,
              }}
            >
              {status.text}
            </span>
            <span
              style={{
                border: "1px solid #ddd",
                padding: "4px 10px",
                borderRadius: "12px",
                fontSize: "0.85rem",
                color: "#555",
              }}
            >
              ID: {problem.problem_id}
            </span>
          </div>

          <p style={{ fontSize: "0.95rem", color: "#444", lineHeight: 1.5 }}>
            {problem.problem_desc || "Описание отсутствует"}
          </p>

          <div style={{ marginTop: 12, fontSize: "0.9rem", color: "#333" }}>
            <strong>Тип:</strong> {typeText}
          </div>

          {lat && lon && (
            <div style={{ marginTop: 12, fontSize: "0.9rem", color: "#333" }}>
              <strong>Координаты:</strong> {lat}, {lon}
              <br />
              <button
                style={{
                  marginTop: 6,
                  padding: "4px 10px",
                  fontSize: "0.85rem",
                  borderRadius: "6px",
                  border: "1px solid #007bff",
                  background: "#fff",
                  color: "#007bff",
                  cursor: "pointer",
                }}
                onClick={() =>
                  window.open(
                    `https://www.google.com/maps?q=${lat},${lon}`,
                    "_blank"
                  )
                }
              >
                Открыть в Google Maps
              </button>
            </div>
          )}

          <div
            style={{
              marginTop: 16,
              fontSize: "0.8rem",
              color: "#777",
              display: "flex",
              justifyContent: "space-between",
            }}
          >
            <span>Дата: {problem.created_at || problem.CreatedAt || "-"}</span>
            <span>Важность: {problem.importance ?? "-"}</span>
          </div>
        </div>
      </div>
    </Container>
  );
}
