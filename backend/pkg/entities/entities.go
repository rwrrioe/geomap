package entities

//DB and geoJSON entities
//TODO add Object struct , separate objects and problems

import (
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

	Problems []Problem `gorm:"foreignKey:DistrictID;references:DistrictID"`
}

type Problem struct {
	ProblemID   int `gorm:"primaryKey;uniqueIndex:idx_problemid"`
	DistrictID  int
	District    District
	Geom        georm.Point `gorm:"type:geometry(Point,4326)"`
	Name        string      `gorm:"not null"`
	Description string      `gorm:"not null"`
	Importance  int         `gorm:"not null"`
	Type        string      `gorm:"not null"`
	Status      string      `gorm:"not null"`
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


//HEATMAP ENTITIES
type HeatMap struct {
	Max        int
	HeatPoints []HeatPoint
}

type HeatPoint struct {
	Geom georm.Point `gorm:"type:geometry(Point,4326)"`
}


type BreefAnswer struct {
	Breef string `json:"breef_answer"`
	Status string  `json:"status"`
}

type ExtendedAnswer struct {
	Extended string `json:"extended_answer"`
	Status string  `json:"status"`
}  

type ProblemTypeMap [int]string {
	1 : "ЖКХ",
	2 : "Дороги и транспорт",
	3 : "Гос.сервис",
	4 : "Прочее",
}


func UnmapProblemType(id int) (string, error) {
	for _, t := range ProblemTypeMap {
		if t == id {
			return ProblemTypeMap[t]
		}
	} else {
		return _, fmt.Errorf("type key not found ")
	}
}