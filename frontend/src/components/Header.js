import React from "react";
import { useNavigate } from "react-router-dom";

const Header = () => {
  const navigate = useNavigate();
  return (
    <div style={{
      background: "#ffffff",
      borderBottom: "1px solid #e9ecef",
      padding: "10px 20px",
      display: "flex",
      alignItems: "center",
      justifyContent: "space-between"
    }}>
      <div
        onClick={() => navigate("/heatmap")}
        style={{
          fontWeight: 700,
          fontSize: 20,
          color: "#0d6efd",
          cursor: "pointer",
          fontFamily: "Segoe UI, Roboto, system-ui"
        }}
      >
        Almaty Problems Geomap
      </div>
    </div>
  );
};

export default Header;