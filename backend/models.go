package main

import "encoding/json"

type Habitation struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	District       string  `json:"district"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
	Population     int     `json:"population"`
	RiskScore      float64 `json:"risk_score"`
	Tier           string  `json:"tier"`
	SlopeScore     float64 `json:"slope_score"`
	TWIScore       float64 `json:"twi_score"`
	LandcoverScore float64 `json:"landcover_score"`
	RainfallScore  float64 `json:"rainfall_score"`
}

type RiskZone struct {
	ZoneID     string          `json:"zone_id"`
	RiskScore  float64         `json:"risk_score"`
	Tier       string          `json:"tier"`
	Population int             `json:"population"`
	HazardType string          `json:"hazard_type"`
	Geometry   json.RawMessage `json:"geometry"`
}

type Site struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	Capacity         int     `json:"capacity"`
	SuitabilityScore float64 `json:"suitability_score"`
	LandCoverType    string  `json:"land_cover_type"`
	DistanceToRoadKm float64 `json:"distance_to_road_km"`
}

type Report struct {
	ID          string  `json:"id"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	PhotoURL    string  `json:"photo_url"`
	ReportedAt  string  `json:"reported_at"`
	Status      string  `json:"status"`
	ReviewedBy  string  `json:"reviewed_by"`
	ReviewedAt  string  `json:"reviewed_at"`
}
