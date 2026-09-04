package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

func RunImport(db *sql.DB) {
	fmt.Println("Starting data import...")

	createTables(db)

	importHabitations(db, "../data/habitations.csv")
	importSites(db, "../data/sites.csv")
	importRiskZones(db, "../data/risk-zones.geojson")
	importReports(db, "../data/reports.csv")

	fmt.Println("Import completed successfully!")
}

func createTables(db *sql.DB) {
	habitationsTable := `
	CREATE TABLE IF NOT EXISTS habitations (
		id TEXT PRIMARY KEY,
		name TEXT,
		district TEXT,
		lat FLOAT,
		lng FLOAT,
		population INTEGER,
		risk_score FLOAT,
		tier TEXT,
		slope_score FLOAT,
		twi_score FLOAT,
		landcover_score FLOAT,
		rainfall_score FLOAT,
		pct_elderly FLOAT,
		pct_children FLOAT
	);`
	_, err := db.Exec(habitationsTable)
	if err != nil {
		log.Fatal("Failed to create habitations table:", err)
	}

	sitesTable := `
	CREATE TABLE IF NOT EXISTS sites (
		id TEXT PRIMARY KEY,
		name TEXT,
		lat FLOAT,
		lng FLOAT,
		capacity INTEGER,
		suitability_score FLOAT,
		land_cover_type TEXT,
		distance_to_road_km FLOAT,
		current_rations INTEGER,
		cots INTEGER,
		medical_kits INTEGER
	);`
	_, err = db.Exec(sitesTable)
	if err != nil {
		log.Fatal("Failed to create sites table:", err)
	}

	riskZonesTable := `
	CREATE TABLE IF NOT EXISTS risk_zones (
		zone_id TEXT PRIMARY KEY,
		risk_score FLOAT,
		tier TEXT,
		population INTEGER,
		hazard_type TEXT,
		geometry TEXT
	);`
	_, err = db.Exec(riskZonesTable)
	if err != nil {
		log.Fatal("Failed to create risk_zones table:", err)
	}
	reportsTable := `
	CREATE TABLE IF NOT EXISTS reports (
		id TEXT PRIMARY KEY,
		lat FLOAT,
		lng FLOAT,
		category TEXT,
		description TEXT,
		photo_url TEXT,
		reported_at TEXT,
		status TEXT,
		reviewed_by TEXT,
		reviewed_at TEXT
	);`
	_, err = db.Exec(reportsTable)
	if err != nil {
		log.Fatal("Failed to create reports table:", err)
	}
	fmt.Println("Tables checked/created.")
}

func importHabitations(db *sql.DB, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Skipping habitations (file not found): %s\n", filepath)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read CSV:", err)
	}

	stmt, err := db.Prepare("INSERT OR REPLACE INTO habitations (id, name, district, lat, lng, population, risk_score, tier, slope_score, twi_score, landcover_score, rainfall_score, pct_elderly, pct_children) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i, record := range records {
		if i == 0 {
			continue
		}
		lat, _ := strconv.ParseFloat(record[3], 64)
		lng, _ := strconv.ParseFloat(record[4], 64)
		pop, _ := strconv.Atoi(record[5])
		riskScore, _ := strconv.ParseFloat(record[6], 64)
		slope, _ := strconv.ParseFloat(record[8], 64)
		twi, _ := strconv.ParseFloat(record[9], 64)
		landcover, _ := strconv.ParseFloat(record[10], 64)
		rainfall, _ := strconv.ParseFloat(record[11], 64)
		elderly, _ := strconv.ParseFloat(record[12], 64)
		children, _ := strconv.ParseFloat(record[13], 64)

		_, err := stmt.Exec(record[0], record[1], record[2], lat, lng, pop, riskScore, record[7], slope, twi, landcover, rainfall, elderly, children)
		if err != nil {
			log.Printf("Failed to insert habitation %s: %v\n", record[0], err)
		}
	}
	fmt.Printf("Imported %d habitations.\n", len(records)-1)
}

func importSites(db *sql.DB, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Skipping sites (file not found): %s\n", filepath)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read CSV:", err)
	}

	stmt, err := db.Prepare("INSERT OR REPLACE INTO sites (id, name, lat, lng, capacity, suitability_score, land_cover_type, distance_to_road_km, current_rations, cots, medical_kits) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i, record := range records {
		if i == 0 {
			continue
		}
		lat, _ := strconv.ParseFloat(record[2], 64)
		lng, _ := strconv.ParseFloat(record[3], 64)
		cap, _ := strconv.Atoi(record[4])
		suit, _ := strconv.ParseFloat(record[5], 64)
		dist, _ := strconv.ParseFloat(record[7], 64)
		rations, _ := strconv.Atoi(record[8])
		cots, _ := strconv.Atoi(record[9])
		meds, _ := strconv.Atoi(record[10])

		_, err := stmt.Exec(record[0], record[1], lat, lng, cap, suit, record[6], dist, rations, cots, meds)
		if err != nil {
			log.Printf("Failed to insert site %s: %v\n", record[0], err)
		}
	}
	fmt.Printf("Imported %d sites.\n", len(records)-1)
}

func importReports(db *sql.DB, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Skipping reports (file not found): %s\n", filepath)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read CSV:", err)
	}

	stmt, err := db.Prepare("INSERT OR REPLACE INTO reports (id, lat, lng, category, description, photo_url, reported_at, status, reviewed_by, reviewed_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i, record := range records {
		if i == 0 {
			continue
		}
		lat, _ := strconv.ParseFloat(record[1], 64)
		lng, _ := strconv.ParseFloat(record[2], 64)

		_, err := stmt.Exec(record[0], lat, lng, record[3], record[4], record[5], record[6], record[7], record[8], record[9])
		if err != nil {
			log.Printf("Failed to insert report %s: %v\n", record[0], err)
		}
	}
	fmt.Printf("Imported %d reports.\n", len(records)-1)
}

func importRiskZones(db *sql.DB, filepath string) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Skipping risk zones (file not found): %s\n", filepath)
		return
	}

	type Feature struct {
		Properties struct {
			ZoneID     string  `json:"zone_id"`
			RiskScore  float64 `json:"risk_score"`
			Tier       string  `json:"tier"`
			Population int     `json:"population"`
			HazardType string  `json:"hazard_type"`
		} `json:"properties"`
		Geometry json.RawMessage `json:"geometry"`
	}

	type FeatureCollection struct {
		Features []Feature `json:"features"`
	}

	var fc FeatureCollection
	if err := json.Unmarshal(file, &fc); err != nil {
		log.Fatal("Failed to parse GeoJSON:", err)
	}

	stmt, err := db.Prepare("INSERT OR REPLACE INTO risk_zones (zone_id, risk_score, tier, population, hazard_type, geometry) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, f := range fc.Features {
		_, err := stmt.Exec(f.Properties.ZoneID, f.Properties.RiskScore, f.Properties.Tier, f.Properties.Population, f.Properties.HazardType, string(f.Geometry))
		if err != nil {
			log.Printf("Failed to insert risk zone %s: %v\n", f.Properties.ZoneID, err)
		}
	}
	fmt.Printf("Imported %d risk zones.\n", len(fc.Features))
}
