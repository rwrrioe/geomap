import { useEffect } from "react";
import L from "leaflet";

export default function ProblemMarkers({ map, points, geoJsonFeatures }) {
  useEffect(() => {
    if (!map || !points?.length) return;

    const markers = [];

    points.forEach((p) => {
      const lat = Number(p.point.lat);
      const lon = Number(p.point.lon);
      if (isNaN(lat) || isNaN(lon)) return;

      // найти districtId по координатам или по district_id из точки
      const districtFeature = geoJsonFeatures.find(
        (f) => Number(f.properties["osm-relation-id"]) === Number(p.district_id)
      );
      const districtId = districtFeature?.properties["osm-relation-id"] || p.district_id;
      const problemId = p.point.problem_id;
      if (!districtId || !problemId) return;

      const marker = L.circleMarker([lat, lon], {
        radius: 8,
        color: "transparent",      // убрали обводку
        fillColor: "rgba(255,0,0,0.3)", // прозрачная точка
        fillOpacity: 0.7,
        weight: 0,
      }).addTo(map);

      marker.on("click", () => {
        window.open(`/heatmap/district/${districtId}/problems/${problemId}`, "_blank");
      });

      markers.push(marker);
    });

    return () => {
      markers.forEach((m) => map.removeLayer(m));
    };
  }, [map, points, geoJsonFeatures]);

  return null;
}
