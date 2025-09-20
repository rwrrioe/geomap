import { useEffect } from "react";
import L from "leaflet";
import districts from "../data/almaty.json";


export default function DistrictLayer({ map, onClick }) {
  useEffect(() => {
    if (!map) return;

    const style = {
      color: "#555",
      weight: 1,
      opacity: 0.5,
      fillOpacity: 0.05, // почти прозрачные
    };

    const highlight = {
      color: "#000",
      weight: 2,
      opacity: 1,
      fillOpacity: 0.1,
    };

    const geoJson = L.geoJSON(districts, {
      style,
      onEachFeature: (feature, layer) => {
        layer.on("click", () => {
          layer.setStyle(highlight);
          onClick(feature.properties.id); // district_id
        });
      },
    });

    geoJson.addTo(map);

    return () => map.removeLayer(geoJson);
  }, [map, onClick]);

  return null;
}