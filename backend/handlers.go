package main

import (
	"encoding/json"
	"net/http"
)

func getHabitations(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, district, lat, lng, population, risk_score, tier, slope_score, twi_score, landcover_score, rainfall_score FROM habitations ORDER BY risk_score DESC")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Habitation
	for rows.Next() {
		var h Habitation
		err := rows.Scan(&h.ID, &h.Name, &h.District, &h.Lat, &h.Lng, &h.Population, &h.RiskScore, &h.Tier, &h.SlopeScore, &h.TWIScore, &h.LandcoverScore, &h.RainfallScore)
		if err != nil {
			continue
		}
		results = append(results, h)
	}

	if results == nil {
		results = []Habitation{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func getRiskZones(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT zone_id, risk_score, tier, population, hazard_type, geometry FROM risk_zones")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Feature struct {
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
		Geometry   json.RawMessage        `json:"geometry"`
	}

	var features []Feature
	for rows.Next() {
		var zoneID, tier, hazardType, geoStr string
		var riskScore float64
		var pop int

		err := rows.Scan(&zoneID, &riskScore, &tier, &pop, &hazardType, &geoStr)
		if err != nil {
			continue
		}

		features = append(features, Feature{
			Type: "Feature",
			Properties: map[string]interface{}{
				"zone_id":     zoneID,
				"risk_score":  riskScore,
				"tier":        tier,
				"population":  pop,
				"hazard_type": hazardType,
			},
			Geometry: []byte(geoStr),
		})
	}

	if features == nil {
		features = []Feature{}
	}

	response := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": features,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getSites(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, lat, lng, capacity, suitability_score, land_cover_type, distance_to_road_km FROM sites ORDER BY suitability_score DESC")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Site
	for rows.Next() {
		var s Site
		err := rows.Scan(&s.ID, &s.Name, &s.Lat, &s.Lng, &s.Capacity, &s.SuitabilityScore, &s.LandCoverType, &s.DistanceToRoadKm)
		if err != nil {
			continue
		}
		results = append(results, s)
	}

	if results == nil {
		results = []Site{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
