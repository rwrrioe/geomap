import React, { useEffect, useRef, useState } from "react";
import L from "leaflet";
import "leaflet.heat";
import { Container, Navbar, Button } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import districts from "../data/almaty.json";
import { getHeatmap, postBriefAnswers } from "../api";

// стиль районов
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
  1: "rgba(0,0,255,5)",   // ЖКХ
  2: "rgba(255,165,0,6)", // Дороги и транспорт
  3: "rgba(0,200,0,6)",   // Гос. сервис
  4: "rgba(255,0,128,6)"  // Прочее
};

const HeatmapPage = () => {
  const mapRef = useRef(null);
  const [heatPoints, setHeatPoints] = useState([]);
  const [briefs, setBriefs] = useState([]);
  const [layerToggles, setLayerToggles] = useState({1:true,2:true,3:true,4:true});
  const [districtLayer, setDistrictLayer] = useState(null);
  const [heatLayers, setHeatLayers] = useState({});
  const navigate = useNavigate();

  // загрузка данных
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

    // создаем слой районов
    const style = { color: "#555", weight: 1, fillOpacity: 0.05 };
    const highlight = { color: "#ff7800", weight: 2, fillOpacity: 0.2 };

    const geoJson = L.geoJSON(districts, {
      style: districtStyle,
      onEachFeature: (feature, layer) => {
        const districtId = feature.properties["osm-relation-id"];
        const districtName = feature.properties.nameRu || districtId;

        // tooltip при наведении
        layer.bindTooltip(districtName, { direction: "auto" });

        // клик на район
        layer.on("click", () => {
          geoJson.eachLayer(l => geoJson.resetStyle(l));
          layer.setStyle(districtHighlight);

          // краткий бриф для всплывашки
          const brief = briefs.find(b => Number(b.district_id) === Number(districtId));
          const briefText = brief?.breef_answer || "Нет данных";

          const popupContent = L.DomUtil.create("div");
          const title = L.DomUtil.create("b", "", popupContent);
          title.innerText = districtName;
          const text = L.DomUtil.create("div", "", popupContent);
          text.innerHTML = briefText;

          const btn = L.DomUtil.create("button", "btn btn-sm btn-primary mt-2", popupContent);
          btn.innerText = "Получить детальный анализ";
          btn.onclick = () => window.open(`/heatmap/analysis/district/${districtId}`, "_blank");

          layer.bindPopup(popupContent).openPopup();
        });
      },
    }).addTo(mapRef.current);
    setDistrictLayer(geoJson);

    return () => {
      if (mapRef.current && districtLayer) mapRef.current.removeLayer(districtLayer);
    };
  }, [briefs]);

  // обновление heatmap слоев
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

      {/* меню справа */}
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
