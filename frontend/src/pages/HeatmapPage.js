import React, { useEffect, useRef, useState } from "react";
import L from "leaflet";
import "leaflet.heat";
import { Container, Navbar, Button } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import districts from "../data/almaty.json";
import { getHeatmap, postBriefAnswers } from "../api";

// Highlight и базовый стиль
const districtStyle = {
  color: "#1a73e8", // насыщенный синий
  weight: 0.7,
  fillOpacity: 0.05,
  opacity: 0.5,
};

const districtHighlight = {
  color: "#ff7800", // оранжевый при клике
  weight: 3,
  fillOpacity: 0.2,
  opacity: 1,
};


const categoryColors = {
  1: "rgba(0,0,255,0.5)",
  2: "rgba(255,165,0,0.5)",
  3: "rgba(0,200,0,0.5)",
  4: "rgba(255,0,128,0.5)"
};

// Функция разбивки MultiPolygon на отдельные Polygon
const flattenMultiPolygon = (feature) => {
  if (feature.geometry.type === "MultiPolygon") {
    return feature.geometry.coordinates.map(coords => ({
      type: "Feature",
      properties: feature.properties,
      geometry: { type: "Polygon", coordinates: coords }
    }));
  }
  return [feature];
};

const flatFeatures = districts.features.flatMap(flattenMultiPolygon);
const flatGeoJSON = { type: "FeatureCollection", features: flatFeatures };

const HeatmapPage = () => {
  const mapRef = useRef(null);
  const [heatPoints, setHeatPoints] = useState([]);
  const [briefs, setBriefs] = useState([]);
  const [layerToggles, setLayerToggles] = useState({1:true,2:true,3:true,4:true});
  const [heatLayers, setHeatLayers] = useState({});
  const navigate = useNavigate();

  useEffect(() => {
    fetchHeatmapAndBriefs();
  }, []);

  const fetchHeatmapAndBriefs = async () => {
    try {
      const heatRes = await getHeatmap();
      setHeatPoints(heatRes.heat_points || []);
      const briefsRes = await postBriefAnswers();
      setBriefs(briefsRes.filter(b => b.status === "ok"));
    } catch (err) {
      console.error("Ошибка загрузки данных:", err);
    }
  };

  useEffect(() => {
    if (!mapRef.current) {
      mapRef.current = L.map("map").setView([43.2389, 76.8897], 12);
      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: '&copy; <a href="http://osm.org">OpenStreetMap</a>',
      }).addTo(mapRef.current);
    }

    const geoJsonLayer = L.geoJSON(flatGeoJSON, {
      style: districtStyle,
      onEachFeature: (feature, layer) => {
        const districtName = feature.properties.nameRu || feature.properties.name || "Без названия";
      
        layer.bindTooltip(districtName, { direction: "auto" });
      
        layer.on("click", () => {
          geoJsonLayer.eachLayer(l => geoJsonLayer.resetStyle(l));
          layer.setStyle(districtHighlight);
          if (layer.bringToFront) layer.bringToFront();
        
          const districtName = feature.properties.nameRu || "Без названия";
          const brief = briefs.find(b => Number(b.district_id) === Number(feature.properties["osm-relation-id"]));
          const briefText = brief?.breef_answer || "Нет данных";
        
          const center = layer.getBounds().getCenter();
          const popupContent = document.createElement("div");
          popupContent.innerHTML = `
            <b>${districtName}</b>
            <div>${briefText}</div>
            <button class="btn btn-sm btn-primary mt-2">Получить детальный анализ</button>
          `;
        
          L.popup({
            closeOnClick: true,
            autoPan: true,
            className: "custom-popup"
          })
            .setLatLng(center)
            .setContent(popupContent)
            .openOn(mapRef.current);
        
            popupContent.querySelector("button").onclick = () => {
              const districtId = feature.properties["osm-relation-id"] || feature.properties.id;
              console.log("DEBUG districtId:", districtId, "properties:", feature.properties);
              if (!districtId) {
                console.warn("⚠️ У этого района нет districtId!", feature.properties);
                return;
              }
              window.open(`/heatmap/analysis/district/${districtId}`, "_blank");
            };
        });
      }
    }).addTo(mapRef.current);

    return () => {
      if (mapRef.current) mapRef.current.removeLayer(geoJsonLayer);
    };
  }, [briefs]);

  useEffect(() => {
    if (!mapRef.current || !heatPoints.length) return;

    const newHeatLayers = {};
    Object.keys(categoryColors).forEach(cat => {
      if (!layerToggles[cat]) return;

      const points = heatPoints
        .filter(p => Number(p.category) === Number(cat))
        .map(p => [
          Number(p.point.lat),
          Number(p.point.lon),
          Math.max(0.01, Math.min(1, Number(p.point.importance)/10))
        ]);

      if (points.length > 0) {
        const layer = L.heatLayer(points, {
          radius: 35,
          blur: 10,
          gradient: {0.4: categoryColors[cat]},
          maxZoom: 17
        }).addTo(mapRef.current);
        newHeatLayers[cat] = layer;
      }
    });

    Object.values(heatLayers).forEach(l => mapRef.current.removeLayer(l));
    setHeatLayers(newHeatLayers);
  }, [heatPoints, layerToggles]);

  const toggleLayer = cat => setLayerToggles(prev => ({...prev, [cat]: !prev[cat]}));

  return (
    <>
      <Navbar bg="dark" variant="dark" style={{padding:"10px 20px", fontFamily:"Segoe UI, Roboto, system-ui"}}>
        <Container>
          <Navbar.Brand
            style={{ cursor: "pointer", fontWeight:700, fontSize:20 }}
            onClick={()=>navigate("/heatmap/analysis/city/-1")}
          >
            Almaty Problems Geomap
          </Navbar.Brand>
        </Container>
      </Navbar>

      <div style={{
        position:"absolute",
        top:60,
        right:10,
        zIndex:1000,
        width:250,
        background:"rgba(255,255,255,0.95)",
        padding:12,
        borderRadius:8,
        boxShadow:"0 2px 8px rgba(0,0,0,0.1)",
        maxHeight:"80vh",
        overflowY:"auto"
      }}>
        <div style={{marginBottom:8,fontWeight:600}}>Типы проблем</div>
        {[[1,"ЖКХ"],[2,"Дороги и транспорт"],[3,"Гос. сервис"],[4,"Прочее"]].map(([cat,name])=>(
          <div key={cat} style={{marginBottom:6, display:"flex", alignItems:"center"}}>
            <input type="checkbox" checked={layerToggles[cat]} onChange={()=>toggleLayer(cat)} />
            <span style={{marginLeft:8, flex:1}}>{name}</span>
            <Button size="sm" variant="link" onClick={()=>window.open(`/heatmap/analysis/type/${cat}`,"_blank")}>Анализ</Button>
          </div>
        ))}

        <div style={{marginTop:12,fontWeight:600}}>Breefs</div>
        {briefs.length===0 && <div style={{fontSize:12,color:"#999"}}>Нет данных</div>}
        {briefs.map((b,i)=>{
          const districtName = districts.features.find(f=>Number(f.properties["osm-relation-id"])===Number(b.district_id))?.properties.nameRu || b.district_id;
          return (
            <div key={i} style={{marginTop:6,paddingTop:6,borderTop:"1px solid #eee"}}>
              <div style={{fontSize:12,fontWeight:600}}>{districtName}</div>
              <div style={{fontSize:12}}>{b.breef_answer}</div>
            </div>
          )
        })}
      </div>

      <div id="map" style={{height:"90vh", width:"100%"}}></div>
    </>
  );
};

export default HeatmapPage;
