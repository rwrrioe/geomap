import axios from "axios";

const API_BASE = "http://localhost:8080"; // поменяй если бэкенд на другом хосте/порту

export const getHeatmap = async () => {
  const res = await axios.get(`${API_BASE}/heatmap`, { timeout: 150000 });
  return res.data;
};

export const postBriefAnswers = async () => {
  const res = await axios.post(`${API_BASE}/heatmap`, {}, { timeout: 200000 });
  // ожидаем массив ответов: [{ district_id: 3072217, breef_answer:"...", status:"ok" }, ...]
  return res.data.responses || res.data; // гибкость формата
};

export const getDistrictAnalysis = async (id) => {
  const res = await axios.get(`${API_BASE}/heatmap/analysis/district/${id}`, { timeout: 20000 });
  return res.data;
};

export const getTypeAnalysis = async (typeId) => {
  const res = await axios.get(`${API_BASE}/heatmap/analysis/type/${typeId}`, { timeout: 20000 });
  return res.data;
};

export const getCityAnalysis = async (cityId) => {
  const res = await axios.get(`${API_BASE}/heatmap/analysis/city/${cityId}`, { timeout: 20000 });
  return res.data;
};