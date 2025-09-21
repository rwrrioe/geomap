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
      <Button as={Link} to="/heatmap" variant="secondary" className="mb-3">
        Back to map
      </Button>
      <Card style={{
  border: "1px solid #ddd",     // чёткая рамка
  borderRadius: "4px",          // почти квадратные углы (можно "0px" если вообще без скруглений)
  maxWidth: "800px",            // ограничиваем ширину для удобного чтения
  margin: "0 auto",             // центрируем
  boxShadow: "0 2px 6px rgba(0,0,0,0.05)" // лёгкая тень
}}>
  <Card.Body>
    <Card.Title style={{ fontWeight: 600, fontSize: "18px" }}>
      District analysis — {districtNameEn}
    </Card.Title>
    <Card.Text style={{
      whiteSpace: "pre-wrap",
      fontSize: "15px",
      lineHeight: "1.6" // удобное чтение текста
    }}>
      {analysis?.extended_answer || JSON.stringify(analysis, null, 2)}
    </Card.Text>
  </Card.Body>
</Card>

    </div>
  );
};

export default DistrictAnalysisPage;

