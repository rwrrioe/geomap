import React, { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import SpinnerOverlay from "../components/SpinnerOverlay";
import { getTypeAnalysis } from "../api";
import { Card, Button } from "react-bootstrap";

const TypeAnalysisPage = () => {
  const { typeId } = useParams();
  const [loading, setLoading] = useState(true);
  const [analysis, setAnalysis] = useState(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        setLoading(true);
        const res = await getTypeAnalysis(typeId);
        if (!cancelled) setAnalysis(res);
      } catch (e) {
        console.error(e);
        setAnalysis({ extended_answer: "No data or backend error.", status: "error" });
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [typeId]);

  if (loading) return <SpinnerOverlay text="Generating type analysis..." />;
  return (
    <div style={{ padding: 16 }}>
      <Button as={Link} to="/geomap/heatmap" variant="secondary" className="mb-3">Back to map</Button>
      <Card>
        <Card.Body>
          <Card.Title>Type analysis — {typeId}</Card.Title>
          <Card.Text style={{ whiteSpace: "pre-wrap" }}>
            {analysis && analysis.extended_answer ? analysis.extended_answer : JSON.stringify(analysis)}
          </Card.Text>
        </Card.Body>
      </Card>
    </div>
  );
};

export default TypeAnalysisPage;