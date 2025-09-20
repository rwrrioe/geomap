import { useEffect } from "react";
import L from "leaflet";
import "leaflet.heat";

const categoryColors = {
  1: "rgba(0,0,255,0.6)",   // ЖКХ
  2: "rgba(255,0,0,0.6)",   // Дороги
  3: "rgba(0,255,0,0.6)",   // Гос. сервис
  4: "rgba(255,165,0,0.6)"  // Прочее
};

export default function HeatLayer({ map, points }) {
  useEffect(() => {
    if (!map || !points?.length) return;

    const layers = [];

    Object.keys(categoryColors).forEach((cat) => {
      const catPoints = points
        .filter((p) => p.category === parseInt(cat))
        .map((p) => [
          p.point.lat,
          p.point.lon,
          p.point.importance / 10, // нормализация
        ]);

      if (catPoints.length > 0) {
        const heatLayer = L.heatLayer(catPoints, {
          radius: 25,
          blur: 10,
          maxZoom: 15,
          gradient: { 0.4: categoryColors[cat] },
        });
        heatLayer.addTo(map);
        layers.push(heatLayer);
      }
    });

    return () => layers.forEach((l) => map.removeLayer(l));
  }, [map, points]);

  return null;
}