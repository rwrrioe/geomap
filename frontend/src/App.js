import React from 'react';
import { Routes, Route, Navigate, BrowserRouter } from 'react-router-dom';
import HeatmapPage from './pages/HeatmapPage';
import DistrictAnalysisPage from './pages/DistrictAnalysisPage';
import TypeAnalysisPage from './pages/TypeAnalysisPage';
import CityAnalysisPage from './pages/CityAnalysisPage';
import ProblemPage from "./pages/ProblemPage";
import ProblemListPage from "./pages/ProblemListPage";
import AddProblemPage from "./pages/AddProblemPage"

function App(){
  return (
    <Routes>
       <Route index element={<Navigate to="heatmap" replace />} />
      <Route path="heatmap" element={<HeatmapPage />} />
      <Route path="heatmap/analysis/district/:districtId" element={<DistrictAnalysisPage />} />
      <Route path="heatmap/analysis/type/:typeId" element={<TypeAnalysisPage />} />
      <Route path="heatmap/analysis/city/:cityId" element={<CityAnalysisPage />} />
      <Route path="heatmap/districts/:districtId/problems/:problemId" element={<ProblemPage />} />
      <Route path="heatmap/districts/:districtId/problems/" element={<ProblemListPage />} />
      <Route path="heatmap/districts/:districtId/problems/new" element={<AddProblemPage />} />
    </Routes>
  );
}

export default App;