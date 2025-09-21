import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import HeatmapPage from './pages/HeatmapPage';
import DistrictAnalysisPage from './pages/DistrictAnalysisPage';
import TypeAnalysisPage from './pages/TypeAnalysisPage';
import CityAnalysisPage from './pages/CityAnalysisPage';

function App(){
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/heatmap" replace />} />
      <Route path="/heatmap" element={<HeatmapPage />} />
      <Route path="/heatmap/analysis/district/:districtId" element={<DistrictAnalysisPage />} />
      <Route path="/heatmap/analysis/type/:typeId" element={<TypeAnalysisPage />} />
      <Route path="/heatmap/analysis/city/:cityId" element={<CityAnalysisPage />} />
    </Routes>
  );
}

export default App;