import { useSearchParams, useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import { Button, Form, Container, Alert, Spinner } from "react-bootstrap";
import { useDropzone } from "react-dropzone";
import districtsData from "../data/almaty.json";

export default function AddProblemPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const lat = parseFloat(searchParams.get("lat"));
  const lon = parseFloat(searchParams.get("lon"));

  const [district, setDistrict] = useState(null);
  const [form, setForm] = useState({
    problem_name: "",
    description: "",
    type_id: 1,
  });

  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(null);

  // Определяем район
  useEffect(() => {
    if (isNaN(lat) || isNaN(lon)) return;

    const flatFeatures = districtsData.features.flatMap((feature) => {
      if (feature.geometry.type === "MultiPolygon") {
        return feature.geometry.coordinates.map((coords) => ({
          type: "Feature",
          properties: feature.properties,
          geometry: { type: "Polygon", coordinates: coords },
        }));
      }
      return [feature];
    });

    const found = flatFeatures.find((f) => {
      const bounds = f.geometry.coordinates.flat(2);
      const lats = bounds.filter((_, i) => i % 2 === 1);
      const lons = bounds.filter((_, i) => i % 2 === 0);
      const minLat = Math.min(...lats);
      const maxLat = Math.max(...lats);
      const minLon = Math.min(...lons);
      const maxLon = Math.max(...lons);
      return lat >= minLat && lat <= maxLat && lon >= minLon && lon <= maxLon;
    });

    if (found) setDistrict(found.properties);
  }, [lat, lon]);

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value });

  // drag & drop
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept: { "image/*": [] },
    onDrop: (acceptedFiles) => {
      setFile(acceptedFiles[0]);
    },
  });

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!district) {
      alert("Район не найден");
      return;
    }

    try {
      setLoading(true);

      const formData = new FormData();
      formData.append("problem_name", form.problem_name);
      formData.append("description", form.description);
      formData.append("type_id", form.type_id);
      formData.append("lat", lat);
      formData.append("lon", lon);
      if (file) {
        formData.append("file", file);
      }

      const res = await fetch(
        `${process.env.REACT_APP_API_URL}/heatmap/districts/${district["osm-relation-id"]}/problems?lat=${lat}&lon=${lon}`,
        {
          method: "POST",
          body: formData,
        }
      );

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || "Ошибка при добавлении");
      }

      const created = await res.json();

      setSuccess("Проблема успешно создана!");
      setTimeout(() => {
        navigate(
          `/heatmap/districts/${district["osm-relation-id"]}/problems/`
        );
      }, 1500);
    } catch (err) {
      console.error(err);
      alert("Не удалось добавить проблему: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container style={{ maxWidth: 600, marginTop: 40 }}>
      <h3>Добавить проблему</h3>

      {district && (
        <p>
          <strong>Район:</strong> {district.nameRu || district.name}
        </p>
      )}

      {success && (
        <Alert variant="success" className="mt-3">
          {success}
        </Alert>
      )}

      <Form onSubmit={handleSubmit}>
        <Form.Group className="mb-3">
          <Form.Label>Название</Form.Label>
          <Form.Control
            type="text"
            name="problem_name"
            value={form.problem_name}
            onChange={handleChange}
            required
          />
        </Form.Group>

        <Form.Group className="mb-3">
          <Form.Label>Описание</Form.Label>
          <Form.Control
            as="textarea"
            rows={3}
            name="description"
            value={form.description}
            onChange={handleChange}
          />
        </Form.Group>

        <div
          {...getRootProps()}
          className="mb-3 p-4 border border-primary rounded text-center"
          style={{
            background: isDragActive ? "#e3f2fd" : "#fafafa",
            cursor: "pointer",
          }}
        >
          <input {...getInputProps()} />
          {file ? (
            <p>Выбрано: {file.name}</p>
          ) : isDragActive ? (
            <p>Отпустите файл здесь...</p>
          ) : (
            <p>Перетащите фото сюда или кликните для выбора</p>
          )}
        </div>

        <Form.Group className="mb-3">
          <Form.Label>Тип проблемы</Form.Label>
          <Form.Select
            name="type_id"
            value={form.type_id}
            onChange={handleChange}
          >
            <option value={1}>ЖКХ</option>
            <option value={2}>Дороги и транспорт</option>
            <option value={3}>Гос. сервис</option>
            <option value={4}>Прочее</option>
          </Form.Select>
        </Form.Group>

        <Button type="submit" variant="primary" disabled={loading}>
          {loading ? <Spinner size="sm" animation="border" /> : "Добавить"}
        </Button>
        <Button
          variant="secondary"
          onClick={() => navigate(-1)}
          style={{ marginLeft: 10 }}
        >
          Отмена
        </Button>
      </Form>
    </Container>
  );
}
