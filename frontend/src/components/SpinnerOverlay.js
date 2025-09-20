import React from "react";
import { Spinner } from "react-bootstrap";

const SpinnerOverlay = ({ text = "Loading..." }) => (
  <div style={{
    position: "absolute",
    top: 0, left: 0, right: 0, bottom: 0,
    display: "flex", alignItems: "center", justifyContent: "center",
    background: "rgba(255,255,255,0.6)", zIndex: 9999
  }}>
    <div className="text-center">
      <Spinner animation="border" role="status" className="mb-2" />
      <div>{text}</div>
    </div>
  </div>
);

export default SpinnerOverlay;