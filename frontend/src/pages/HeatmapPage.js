import React, { useEffect, useRef, useState } from "react";
import L from "leaflet";
import "leaflet.heat";
import { Container, Navbar, Button } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import districts from "../data/almaty.json";
import { getHeatmap, postBriefAnswers } from "../api";

const districtStyle = { color: "#1a73e8", weight: 0.7, fillOpacity: 0.05, opacity: 0.5 };
const districtHighlight = { color: "#ff7800", weight: 3, fillOpacity: 0.2, opacity: 1 };
const categoryColors = { 1: "rgba(0,0,255,0.5)", 2: "rgba(255,165,0,0.5)", 3: "rgba(0,200,0,0.5)", 4: "rgba(255,0,128,0.5)" };

const flattenMultiPolygon = (feature) =>
  feature.geometry.type === "MultiPolygon"
    ? feature.geometry.coordinates.map((coords) => ({
        type: "Feature",
        properties: feature.properties,
        geometry: { type: "Polygon", coordinates: coords },
      }))
    : [feature];

const flatFeatures = districts.features.flatMap(flattenMultiPolygon);

export default function HeatmapPage() {
  const mapRef = useRef(null);
  const geoJsonLayerRef = useRef(null);
  const heatLayersRef = useRef({});
  const [heatPoints, setHeatPoints] = useState([]);
  const [briefs, setBriefs] = useState([]);
  const [zoom, setZoom] = useState(12);
  const [layerToggles, setLayerToggles] = useState({ 1: true, 2: true, 3: true, 4: true });
  const [addingMode, setAddingMode] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const heatRes = await getHeatmap();
      setHeatPoints(heatRes.heat_points || []);
      const briefsRes = await postBriefAnswers();
      setBriefs(briefsRes.filter((b) => b.status === "ok"));
    } catch (err) {
      console.error(err);
    }
  };

  // Инициализация карты
  useEffect(() => {
    const mapDiv = document.getElementById("map");
    if (!mapDiv) return;
  
    const initMap = () => {
      if (mapRef.current) return;
  
      mapRef.current = L.map("map").setView([43.2389, 76.8897], 12);
  
      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: '&copy; <a href="http://osm.org">OpenStreetMap</a>',
      }).addTo(mapRef.current);
  
      mapRef.current.on("zoomend", () => setZoom(mapRef.current.getZoom()));
    };
  
    if (mapDiv.offsetWidth === 0 || mapDiv.offsetHeight === 0) {
      // ждем, пока div станет видимым
      const id = requestAnimationFrame(initMap);
      return () => cancelAnimationFrame(id);
    } else {
      initMap();
    }
  }, []);
  

  // Режим добавления: клик по карте
  useEffect(() => {
    if (!mapRef.current) return;

    const handleClick = (e) => {
      if (!addingMode) return;
      const { lat, lng } = e.latlng;

      let districtId = null;
      if (geoJsonLayerRef.current) {
        geoJsonLayerRef.current.eachLayer((layer) => {
          if (layer.getBounds().contains([lat, lng])) {
            districtId = layer.feature.properties["osm-relation-id"];
          }
        });
      }

      if (!districtId) {
        alert("Не удалось определить район. Увеличьте зум и попробуйте снова.");
        return;
      }

      navigate(`/geomap/heatmap/districts/${districtId}/problems/new?lat=${lat}&lon=${lng}`);
      setAddingMode(false);
    };

    mapRef.current.on("click", handleClick);
    return () => mapRef.current.off("click", handleClick);
  }, [addingMode, navigate]);

  // Районы
  useEffect(() => {
    if (!mapRef.current) return;

    if (zoom <= 13) {
      if (!geoJsonLayerRef.current) {
        geoJsonLayerRef.current = L.geoJSON(flatFeatures, {
          style: districtStyle,
          onEachFeature: (feature, layer) => {
            const districtName = feature.properties.nameRu || feature.properties.name || "Без названия";
            layer.bindTooltip(districtName, { direction: "auto" });
            layer.on("click", () => {
              geoJsonLayerRef.current.eachLayer((l) => geoJsonLayerRef.current.resetStyle(l));
              layer.setStyle(districtHighlight);
              if (layer.bringToFront) layer.bringToFront();

              const districtId = feature.properties["osm-relation-id"];
              const center = layer.getBounds().getCenter();
              const popupContent = document.createElement("div");
              popupContent.innerHTML = `
                <b>${districtName}</b>
                <div style="margin-top:8px; display:flex; gap:6px;">
                  <button class="btn btn-sm btn-primary" id="btn-analysis">Анализ</button>
                  <button class="btn btn-sm btn-primary" id="btn-problems">Проблемы</button>
                </div>
              `;
              const popup = L.popup({ closeOnClick: true, autoPan: true, className: "custom-popup" })
                .setLatLng(center)
                .setContent(popupContent)
                .openOn(mapRef.current);

              popupContent.querySelector("#btn-analysis").onclick = () => {
                window.open(`/geomap/heatmap/analysis/district/${districtId}`, "_blank");
              };

              popupContent.querySelector("#btn-problems").onclick = () => {
                window.open(`/geomap/heatmap/districts/${districtId}/problems`, "_blank");
              };
            });
          },
        }).addTo(mapRef.current);
      } else {
        mapRef.current.addLayer(geoJsonLayerRef.current);
      }
    } else if (geoJsonLayerRef.current) {
      mapRef.current.removeLayer(geoJsonLayerRef.current);
    }
  }, [zoom]);

  // Heatmap и точки
  // Heatmap и точки
useEffect(() => {
  if (!mapRef.current || !heatPoints.length || mapRef.current.getSize().x === 0) return;

  // удаляем старые слои
  Object.values(heatLayersRef.current).forEach((l) => mapRef.current.removeLayer(l));
  heatLayersRef.current = {};

  Object.keys(categoryColors).forEach((cat) => {
    if (!layerToggles[cat]) return;

    const points = heatPoints.filter((p) => Number(p.category) === Number(cat));
    if (!points.length) return;

    // heatmap слой
    const heatData = points.map((p) => [
      Number(p.point.lat),
      Number(p.point.lon),
      Math.max(0.01, Math.min(1, Number(p.point.importance) / 10)),
    ]);
    const heatLayer = L.heatLayer(heatData, {
      radius: 35,
      blur: 10,
      gradient: { 0.4: categoryColors[cat] },
      maxZoom: 17,
    });

    // geoJSON слой для кликов
    const geoJson = {
      type: "FeatureCollection",
      features: points.map((p) => ({
        type: "Feature",
        properties: { problem_id: p.point.problem_id, lat: p.point.lat, lon: p.point.lon },
        geometry: { type: "Point", coordinates: [Number(p.point.lon), Number(p.point.lat)] },
      })),
    };

    const pointLayer = L.geoJSON(geoJson, {
      pointToLayer: (_, latlng) =>
        L.circleMarker(latlng, {
          radius: 8,
          color: "transparent",
          fillColor: "transparent",
          fillOpacity: 0,
          weight: 0,
        }),
      onEachFeature: (feature, layer) => {
        layer.on("click", () => {
          if (mapRef.current.getZoom() <= 13) return;

          let districtOsmId = null;
          if (geoJsonLayerRef.current) {
            geoJsonLayerRef.current.eachLayer((l) => {
              if (l.getBounds().contains([feature.properties.lat, feature.properties.lon])) {
                districtOsmId = l.feature.properties["osm-relation-id"];
              }
            });
          }

          if (!districtOsmId) {
            console.warn("Не удалось определить район для точки", feature);
            return;
          }

          window.open(
            `/geomap/heatmap/districts/${districtOsmId}/problems/${feature.properties.problem_id}`,
            "_blank"
          );
        });
      },
    });

    // один слой на категорию
    const group = L.layerGroup([heatLayer, pointLayer]).addTo(mapRef.current);
    heatLayersRef.current[cat] = group;
  });
}, [heatPoints, layerToggles]);


  const toggleLayer = (cat) =>
    setLayerToggles((prev) => ({
      ...prev,
      [cat]: !prev[cat],
    }));

  const renderBriefs = (briefList) => {
    if (!briefList || briefList.length === 0)
      return <div style={{ fontSize: 12, color: "#999" }}>Нет данных</div>;
    return briefList.map((b, i) => {
      const districtName =
        districts.features.find((f) => Number(f.properties["osm-relation-id"]) === Number(b.district_id))
          ?.properties.nameRu || b.district_id;
      return (
        <div key={i} style={{ marginTop: 6, paddingTop: 6, borderTop: "1px solid #eee" }}>
          <div style={{ fontSize: 12, fontWeight: 600 }}>{districtName}</div>
          <div style={{ fontSize: 12 }}>{b.breef_answer}</div>
        </div>
      );
    });
  };

  return (
    <>
      <Navbar bg="dark" variant="dark" style={{ padding: "10px 20px", fontFamily: "Segoe UI, Roboto, system-ui" }}>
        <Container>
          <Navbar.Brand
            style={{ cursor: "pointer", fontWeight: 700, fontSize: 20 }}
            onClick={() => navigate("/geomap/heatmap/analysis/city/-1")}
          >
            Almaty Problems Geomap
          </Navbar.Brand>
        </Container>
      </Navbar>

      {/* Панель справа */}
      <div
        style={{
          position: "absolute",
          top: 60,
          right: 10,
          zIndex: 1000,
          width: 250,
          background: "rgba(255,255,255,0.95)",
          padding: 12,
          borderRadius: 8,
          boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
          maxHeight: "80vh",
          overflowY: "auto",
        }}
      >
        <div style={{ marginBottom: 8, fontWeight: 600 }}>Типы проблем</div>
        {[
          [1, "ЖКХ"],
          [2, "Дороги и транспорт"],
          [3, "Гос. сервис"],
          [4, "Прочее"],
        ].map(([cat, name]) => (
          <div key={cat} style={{ marginBottom: 6, display: "flex", alignItems: "center" }}>
            <input type="checkbox" checked={layerToggles[cat]} onChange={() => toggleLayer(cat)} />
            <span style={{ marginLeft: 8, flex: 1 }}>{name}</span>
            <Button
              size="sm"
              variant="link"
              onClick={() => window.open(`/geomap/heatmap/analysis/type/${cat}`, "_blank")}
            >
              Анализ
            </Button>
          </div>
        ))}

        <div style={{ marginTop: 12, fontWeight: 600 }}>Breefs</div>
        {renderBriefs(briefs)}

        {/* Кнопка добавить проблему */}
        <Button
          variant="success"
          style={{ marginTop: 12, width: "100%" }}
          onClick={() => setAddingMode(true)}
        >
          Добавить проблему
        </Button>

        {addingMode && (
          <div style={{ marginTop: 6, fontSize: 12, color: "#555" }}>
            Режим добавления: нажмите на карту в районе с проблемой
          </div>
        )}
      </div>

      <div id="map" style={{ height: "90vh", width: "100%" }}></div>
    </>
  );
}
