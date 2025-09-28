import React, { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import SpinnerOverlay from "../components/SpinnerOverlay";
import { getDistrictAnalysis } from "../api";
import { Card, Button } from "react-bootstrap";
import districts from "../data/almaty.json"; 

const DistrictAnalysisPage = () => {
  const { districtId } = useParams();
  const [loading, setLoading] = useState(true);
  const [analysis, setAnalysis] = useState(null);

  // ищем район по districtId
  const flatFeatures = districts.features.flatMap(feature => {
    if (feature.geometry.type === "MultiPolygon") {
      return feature.geometry.coordinates.map(coords => ({
        type: "Feature",
        properties: feature.properties,
        geometry: { type: "Polygon", coordinates: coords }
      }));
    }
    return [feature];
  });

  const district = flatFeatures.find(
    f => String(f.properties["osm-relation-id"]) === String(districtId)
  );

  const districtNameEn = district?.properties?.name || districtId;

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        setLoading(true);
        console.log("📡 Fetching analysis for districtId:", districtId);
        const res = await getDistrictAnalysis(districtId);
        if (!cancelled) setAnalysis(res);
      } catch (e) {
        console.error("❌ Error while fetching analysis:", e);
        setAnalysis({ extended_answer: "No data or backend error.", status: "error" });
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [districtId]);

  if (loading) {
    return <SpinnerOverlay text="Generating detailed district analysis..." />;
  }

  return (
    <div style={{ padding: 16 }}>
      <Card
        style={{
          border: "1px solid #ddd",
          borderRadius: "4px",
          maxWidth: "800px",
          margin: "0 auto",
          boxShadow: "0 2px 6px rgba(0,0,0,0.05)"
        }}
      >
        <Card.Body>
          <Card.Title style={{ fontWeight: 600, fontSize: "18px" }}>
            District analysis — {districtNameEn}
          </Card.Title>
          <Card.Text
            style={{
              whiteSpace: "pre-wrap",
              fontSize: "15px",
              lineHeight: "1.6"
            }}
          >
            {analysis?.extended_answer || JSON.stringify(analysis, null, 2)}
          </Card.Text>
        </Card.Body>
      </Card>

      {/* Кнопки снизу */}
      <div
        style={{
          marginTop: "24px",
          display: "flex",
          justifyContent: "space-between",
          maxWidth: "800px",
          marginLeft: "auto",
          marginRight: "auto"
        }}
      >
        <Button
          as={Link}
          to={`/heatmap/districts/${districtId}/problems`}
          variant="primary"
        >
          Список проблем
        </Button>
        <Button as={Link} to="/heatmap" variant="secondary">
          Назад к карте
        </Button>
      </div>
    </div>
  );
};

export default DistrictAnalysisPage;