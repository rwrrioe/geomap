package entities

//DB and geoJSON entities
//TODO add Object struct , separate objects and problems

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/ybru-tech/georm"
)

type DistrictDTO struct {
	Type       string `json:"type"`
	Properties struct {
		Osm_relation_id int    `json:"osm-relation-id,string"`
		Name            string `json:"name"`
		NameRu          string `json:"nameRu"`
	} `json:"properties"`
	Geometry struct {
		Type        string          `json:"type"`
		Coordinates [][][][]float64 `json:"coordinates"`
	} `json:"geometry"`
}

type District struct {
	DistrictID int                `gorm:"primaryKey;uniqueIndex:idx_distinctid"`
	NameRU     string             `gorm:"not null"`
	NameENG    string             `gorm:"not null"`
	Type       string             `gorm:"not null"`
	Geom       georm.MultiPolygon `gorm:"type:geometry(MultiPolygon,4326)"`
	Reputation float64
	Problems   []Problem `gorm:"foreignKey:DistrictID;references:DistrictID"`
}

type Problem struct {
	ProblemID   int `gorm:"primaryKey;autoIncrement;uniqueIndex:idx_problemid"`
	DistrictID  int
	District    District
	Geom        georm.Point `gorm:"type:geometry(Point,4326)"`
	Name        string      `gorm:"not null"`
	Description string      `gorm:"not null"`
	ImageURL    string      `gorm:"column:image_url"`
	Importance  float64     `gorm:"not null"`
	Status      string      `gorm:"not null"`
	TypeId      int         `gorm:"not null"`
}

type ProblemResponseDTO struct {
	ProblemID   int         `gorm:"column:problem_id"`
	DistrictID  int         `gorm:"column:district_id"`
	District    District    `gorm:"column:district"`
	Geom        georm.Point `gorm:"column:geom"`
	Name        string      `gorm:"name"`
	Description string      `gorm:"description"`
	Importance  float64     `gorm:"importance"`
	Status      string      `gorm:"status"`
}

type CreateProblemForm struct {
	ProblemName string `form:"problem_name" binding:"required"`
	ImageURL    string `form:"image_url"`
	Description string `form:"description"`
	TypeID      int    `form:"type_id" binding:"required"`
	Lat         float64
	Lon         float64
}

type ProblemType struct {
	TypeId   int    `gorm:"primaryKey;column:type_id"`
	TypeName string `gorm:"column:type"`
}

func MapToDistinct(dto DistrictDTO) *District {
	return &District{
		DistrictID: int(dto.Properties.Osm_relation_id),
		NameRU:     dto.Properties.NameRu,
		NameENG:    dto.Properties.Name,
		Type:       dto.Geometry.Type,
		Geom:       *multiPolygonFromDTO(dto),
	}
}

func multiPolygonFromDTO(dto DistrictDTO) *georm.MultiPolygon {
	mp := geom.NewMultiPolygon(geom.XY)

	for _, polygon := range dto.Geometry.Coordinates {
		var rings []*geom.LinearRing
		for _, ring := range polygon {
			flat := flattenRing(ring)
			flatRing := geom.NewLinearRingFlat(geom.XY, flat)
			rings = append(rings, flatRing)
		}
		poly := geom.NewPolygon(geom.XY)
		for _, r := range rings {
			poly.Push(r)
		}
		mp.Push(poly)
	}

	geormMP := georm.New(mp)
	return &geormMP
}

func flattenRing(ring [][]float64) []float64 {
	flat := []float64{}
	for _, p := range ring {
		flat = append(flat, p[0], p[1])
	}
	return flat
}

// USER ENTITIES
type User struct {
	ID   int
	Name string
	Role string // "guest", "admin", "user", etc...
}

// HEATMAP ENTITIES
type HeatMap struct {
	Max        int         `json:"max_points"`
	HeatPoints []HeatPoint `json:"heat_points"`
}

func NewHeatpMap() *HeatMap {
	return &HeatMap{}
}

type CachedHeatMap struct {
	HeatMapID int     `gorm:"column:heatmap_id;primaryKey;autoIncrement;-><-:create"`
	HeatMap   HeatMap `gorm:"column:heatmap_data;type:json"`
}

func (CachedHeatMap) TableName() string {
	return "cached_heatmaps"
}

type HeatPoint struct {
	Category int   `json:"category"`
	Point    Point `json:"point"`
}

type Point struct {
	DistrictId int     `json:"district_id"`
	Id         int     `json:"problem_id"`
	Lon        float64 `json:"lon"`
	Lat        float64 `json:"lat"`
	Importance float64 `json:"importance"`
}

// HEATMAP VALUER/SCANNER
func (h HeatMap) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *HeatMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, h)
}

type BreefAIResponse struct {
	DistrictID int    `json:"district_id"`
	Breef      string `json:"breef_answer"`
	Status     string `json:"status"`
}

type HeatMapResponse struct {
	Responses []BreefAIResponse `json:"responses"`
}

type CachedAnswer struct {
	AnswerID     int    `gorm:"primaryKey;autoIncrement;-><-:create"`
	ResponseText string `gorm:"column:response_text"`
	Status       string `gorm:"column:status"`
	RequestID    int    `gorm:"column:request_id"`
}

type ExtendedAIResponse struct {
	AnswerText string `json:"extended_answer"`
	Status     string `json:"status"`
}

var ProblemTypeMap = map[int]string{
	1: "ЖКХ",
	2: "Дороги и транспорт",
	3: "Гос.сервис",
	4: "Прочее",
}
