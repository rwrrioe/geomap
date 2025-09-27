import axios from "axios";



export const getHeatmap = async () => {
  const res = await axios.get(`http://${process.env.REACT_APP_API_URL}/heatmap`, { timeout: 150000 });
  return res.data;
};

export const postBriefAnswers = async () => {
  const res = await axios.post(`http://${process.env.REACT_APP_API_URL}/heatmap`, {}, { timeout: 200000 });
  // ожидаем массив ответов: [{ district_id: 3072217, breef_answer:"...", status:"ok" }, ...]
  return res.data.responses || res.data; // гибкость формата
};

export const getDistrictAnalysis = async (id) => {
  const res = await axios.get(`http://${process.env.REACT_APP_API_URL}/heatmap/analysis/district/${id}`, { timeout: 20000 });
  return res.data;
};

export const getTypeAnalysis = async (typeId) => {
  const res = await axios.get(`http://${process.env.REACT_APP_API_URL}/heatmap/analysis/type/${typeId}`, { timeout: 20000 });
  return res.data;
};

export const getCityAnalysis = async (cityId) => {
  const res = await axios.get(`http://${process.env.REACT_APP_API_URL}/heatmap/analysis/city/${cityId}`, { timeout: 20000 });
  return res.data;
};

export const getProblem = async (districtId, problemId) => {
  const res = await axios.get(`http://${process.env.REACT_APP_API_URL}/heatmap/districts/${districtId}/problems/${problemId}`, { timeout: 20000 });
  return res.data;
};
