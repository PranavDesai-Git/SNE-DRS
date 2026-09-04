package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getHabitations(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, district, lat, lng, population, risk_score, tier, slope_score, twi_score, landcover_score, rainfall_score, pct_elderly, pct_children FROM habitations ORDER BY risk_score DESC")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Habitation
	for rows.Next() {
		var h Habitation
		err := rows.Scan(&h.ID, &h.Name, &h.District, &h.Lat, &h.Lng, &h.Population, &h.RiskScore, &h.Tier, &h.SlopeScore, &h.TWIScore, &h.LandcoverScore, &h.RainfallScore, &h.PctElderly, &h.PctChildren)
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
	log.Printf("[Backend] Returned %d habitations", len(results))
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
	log.Printf("[Backend] Returned %d risk zone features", len(features))
}

func getSites(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, lat, lng, capacity, suitability_score, land_cover_type, distance_to_road_km, current_rations, cots, medical_kits FROM sites ORDER BY suitability_score DESC")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Site
	for rows.Next() {
		var s Site
		err := rows.Scan(&s.ID, &s.Name, &s.Lat, &s.Lng, &s.Capacity, &s.SuitabilityScore, &s.LandCoverType, &s.DistanceToRoadKm, &s.CurrentRations, &s.Cots, &s.MedicalKits)
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
	log.Printf("[Backend] Returned %d sites", len(results))
}

func getReports(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, lat, lng, category, description, photo_url, reported_at, status, reviewed_by, reviewed_at FROM reports ORDER BY reported_at DESC")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []Report
	for rows.Next() {
		var rep Report
		err := rows.Scan(&rep.ID, &rep.Lat, &rep.Lng, &rep.Category, &rep.Description, &rep.PhotoURL, &rep.ReportedAt, &rep.Status, &rep.ReviewedBy, &rep.ReviewedAt)
		if err != nil {
			continue
		}
		results = append(results, rep)
	}

	if results == nil {
		results = []Report{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
	log.Printf("[Backend] Returned %d reports", len(results))
}

func getRoutes(w http.ResponseWriter, r *http.Request) {
	habitationID := r.URL.Query().Get("habitation_id")
	if habitationID == "" {
		http.Error(w, "Missing habitation_id", http.StatusBadRequest)
		return
	}

	// Mocking a route from the habitation to a safe site
	// In a real app, this would use osmnx/networkx or similar for risk-weighted routing
	// We'll just return a mock GeoJSON LineString
	var hLat, hLng float64
	err := db.QueryRow("SELECT lat, lng FROM habitations WHERE id = ?", habitationID).Scan(&hLat, &hLng)
	if err != nil {
		http.Error(w, "Habitation not found", http.StatusNotFound)
		return
	}

	var sLat, sLng float64
	err = db.QueryRow("SELECT lat, lng FROM sites ORDER BY suitability_score DESC LIMIT 1").Scan(&sLat, &sLng)
	if err != nil {
		http.Error(w, "No sites available", http.StatusInternalServerError)
		return
	}

	// Create a simple 3-point line representing a mock route
	midLat := (hLat + sLat) / 2.0
	midLng := (hLng + sLng) / 2.0

	routeGeoJSON := map[string]interface{}{
		"type": "FeatureCollection",
		"features": []map[string]interface{}{
			{
				"type": "Feature",
				"properties": map[string]interface{}{
					"route_id": "route_1",
					"risk_avoided": true,
				},
				"geometry": map[string]interface{}{
					"type": "LineString",
					"coordinates": [][]float64{
						{hLng, hLat},
						{midLng + 0.01, midLat + 0.01}, // Slight curve
						{sLng, sLat},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routeGeoJSON)
	log.Printf("[Backend] Returned mock route for habitation %s", habitationID)
}
