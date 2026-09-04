# SIH26191 — Hazard Red Zones & Relocation Platform
### Master Build Guide — 3-District NE Pilot (Meghalaya / Arunachal Pradesh / Assam)

---

## 1. The problem, precisely

India's disaster-hazard data infrastructure is genuinely mature — NDEM (ISRO/NRSC, operational since 2013) already does multi-hazard zonation and real-time monitoring at national scale. What doesn't exist, evidenced concretely by Joshimath (a nationally-watched crisis where, a year in, zero families had been relocated, site selection was ad hoc, and one candidate site provoked resident protest for lack of a transparent case) is the **decision layer**: turning hazard data into a ranked, actionable, capacity-aware relocation plan before a disaster forces the issue.

That decision layer is what this platform is.

---

## 2. Product structure — one platform, three escalating modes

Not three separate apps. One shared data backbone (hazard scores, population, candidate relocation sites), three views that activate progressively. Each mode is a **complete, valid submission on its own** — this is also literally the build order.

**Three real-world scenarios this plays out in:**
- **Routine planning (Mode 1), any normal week** — an official reviews the ranked habitation list, checks why the top village is flagged, sees which safe site has capacity for it. Proactive planning, months before a crisis.
- **Pre-monsoon escalation (Mode 2), as a storm approaches** — rainfall nears the cloudburst threshold, Red Zones expand dynamically, evacuation routes and capacity-matched shelter assignments get computed before the disaster hits.
- **Active event (Mode 3), during an unfolding disaster** — a citizen report of a blockage gets confirmed and the evacuation routes reroute automatically around it.

### Mode 1 — Long-term (baseline, must complete)
Standing risk assessment and relocation planning, used year-round for proactive government planning.
- Interactive map, color-coded by risk tier (Immediate / Short-term / Medium-term — plain language, not raw scores)
- Click any habitation → hazard breakdown, population, contributing factors
- Ranked, sortable, exportable priority list
- Candidate relocation sites with suitability + estimated capacity
- Data provenance panel (sources, last updated) — this is what makes the transparent MCDA scoring credible to a non-technical government user

### Mode 2 — Pre-disaster (build if Mode 1 finishes with time to spare)
Activates when risk crosses a threshold (rainfall slider, or an IMD-style warning).
- Same map, now showing **risk-weighted evacuation routes** — shortest path that avoids other hazard zones, not just shortest distance
- **Capacity-matched displacement assignment** — which habitation's population goes to which site, respecting each site's real capacity limit
- Near-real-time flood layer from Sentinel-1 SAR (sees through monsoon cloud cover, genuinely operational-grade technique)

### Mode 3 — During-disaster (stretch, only if well ahead of schedule)
Activates only during an active event.
- Citizen-submitted reports overlaid live on the map (PetaBencana-style: structured, two-step-verified, not raw unverified text)
- Confirmed reports of a blocked road trigger automatic rerouting in Mode 2's routing layer
- Built as a fully functional pipeline, demoed with seeded sample reports (you won't have real users during a hackathon demo — that's fine and expected, be upfront about it)

### Action Center (template-based, no LLM, no GPU dependency)
Click a habitation, get a readable summary generated from a plain string template — no model, no inference, nothing that can hallucinate a wrong number in front of judges:
```python
f"{name} is flagged {tier} priority, hazard score {score}/1.0. "
f"Nearest viable relocation site: {site} (capacity {capacity}, {distance}km away)."
```
This is deliberately simple. The platform's actual AI/ML substance is Mode 1's Random Forest–informed MCDA scoring — that's what satisfies the "AI-driven" bar, not this panel. This template just presents that output readably, with zero implementation risk and zero tension with the transparent-scoring pitch angle used elsewhere in this build.

---

## 3. Tech stack

| Component | Tool | Mode |
|---|---|---|
| GIS processing | QGIS (SRTM-Downloader, Bhuvan Web Services, Point Sampling Tool plugins) | 1 |
| Elevation source | CartoDEM (Cartosat-1, via Bhuvan Open Data Archive, 30m, ~8-10m vertical accuracy) | 1 |
| Terrain features | QGIS Slope + SAGA toolbox (TWI, curvature) — from CartoDEM for the 2 landslide districts, SRTM for Assam | 1 |
| ML weight derivation | scikit-learn `RandomForestClassifier` on GSI's historical landslide points | 1 |
| Risk scoring | Transparent MCDA weighted sum (weights informed by the RF's feature importances) | 1 |
| Boundaries | `yashveeeeeeer/india-geodata` (merges LGD, Survey of India, Bhuvan, DataMeet) OR DataMeet (`datta07/INDIAN-SHAPEFILES`) — not GADM, given Arunachal Pradesh border sensitivity | 1 |
| Roads/infra | Geofabrik NE Zone `.osm.pbf`, opened directly in QGIS | 1, 2 |
| Population | WorldPop 100m raster (regional base) + Census 2011 village-level figures via `india-geodata` OR GHS-POP 100m grid | 1 |
| Flood hazard (static) | Bhuvan Flood Hazard Zonation Atlas / India-WRIS (registration required — register early if used) | 1 |
| Flood hazard (near-real-time) | EOS-04 via Bhoonidhi (primary), Sentinel-1 SAR via Google Earth Engine (fallback, free, server-side, no local GPU) | 2 |
| Routing | `osmnx` + `networkx`, risk-weighted edge costs | 2 |
| Capacity assignment | Greedy priority-first allocation (Python/pandas) | 2 |
| Citizen reports | Svelte/React form → Go/FastAPI `POST /reports`, rule-based two-report confirmation | 3 |
| Dashboard | Svelte + Leaflet OR React (Vite) + TailwindCSS + React-Leaflet, served by a Go/FastAPI backend | 1, 2, 3 |
| Data storage | SQLite, populated by a local import script — not a live cross-machine connection | 1, 2, 3 |

**Team split (2-person):** one person owns the data/ML/GIS pipeline (QGIS, the RF model, MCDA scoring, routing, capacity matching) and exports results to static GeoJSON/CSV files. The other owns the backend + frontend (this pilot supports either Go + Svelte/Leaflet or FastAPI + React), serving that data via REST endpoints. **Data flow, precisely:** the GIS/ML teammate exports fresh GeoJSON/CSV files whenever they recompute something and hands them off (shared repo, Slack, synced folder — anything that isn't a live connection); the backend teammate runs a small **import script** that loads those files into a local **SQLite** database, and the API queries SQLite rather than re-parsing files per request. This gets you real `WHERE`/`ORDER BY`/`JOIN` ergonomics for things like sorting `/habitations` by risk score, without either machine needing to be network-reachable during dev or demo. Explicitly avoid: a live database server (MySQL/Postgres) with the compute machine pushing writes to the backend machine over the network mid-demo — it adds a cross-laptop dependency that can fail in front of judges for no real benefit over the file → import script → SQLite flow. The backend should serve precomputed data, not compute live on request, so demo-day performance doesn't depend on the GIS pipeline running in the request path. Agree on the exact API contract before either side starts building:
- `GET /habitations` → ranked list (id, name, tier, population, risk_score)
- `GET /risk-zones` → GeoJSON, properties: `risk_score`, `tier`, `population`
- `GET /sites` → candidate relocation sites + capacity
- `GET /routes?habitation_id=X` → Mode 2 routing output
- `POST /reports` / `GET /reports` → Mode 3 citizen reports
- `GET /summary?habitation_id=X` → triggers the template-based Action Center summary

No GPU or CUDA required anywhere. Google Earth Engine's compute happens on Google's servers for free under education use, not on your laptop.

**Local vs. remote, precisely:** everything runs fully offline on your machine once source data is downloaded (QGIS processing, the RF model, MCDA scoring, routing, capacity assignment, and the dashboard itself) — the one genuine exception is the SAR flood analysis, which runs on Google Earth Engine's servers by design. For demo day, **pre-run the GEE export ahead of time and store the result locally** rather than querying it live in front of judges — removes venue-WiFi risk from your only remote-dependent step, and the dashboard itself needs no internet to run either way.

---

## 3.5 API data contract (finalized)

This is the exact file format and SQLite schema for each endpoint, agreed between the GIS/ML teammate (exports these files) and the backend teammate (imports them into SQLite, serves them). The GIS/ML teammate should export to these formats directly — column names and structure are locked, not placeholders.

### `GET /habitations`
**Purpose:** a ranked list of villages/habitations, most at-risk first.

Exported as CSV:
```csv
id,name,district,lat,lng,population,risk_score,tier,slope_score,twi_score,landcover_score,rainfall_score
HAB001,Sorbhog,Barpeta,26.4521,90.9871,3200,8.7,Immediate,0.82,0.75,0.60,0.90
HAB002,Bahari,Barpeta,26.3890,91.0120,1850,6.2,Short-term,0.55,0.60,0.45,0.70
```

SQLite schema:
```sql
CREATE TABLE habitations (
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
    rainfall_score FLOAT
);
```
`tier` is plain-language text (Immediate/Short-term/Medium-term), not a numeric bucket — the frontend shouldn't have to do threshold logic. The `*_score` columns are the contributing-factor breakdown shown when a habitation is clicked; confirm the real feature-name set with the ML pipeline once the RF model's features are final.

### `GET /risk-zones`
**Purpose:** the map visualization layer — hazard zone polygons, not points, colored by tier.

Exported as GeoJSON:
```json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {
        "zone_id": "RZ001",
        "risk_score": 8.7,
        "tier": "Immediate",
        "population": 3200,
        "hazard_type": "flood"
      },
      "geometry": {
        "type": "Polygon",
        "coordinates": [[[90.9850, 26.4500], [90.9900, 26.4500], [90.9900, 26.4550], [90.9850, 26.4550], [90.9850, 26.4500]]]
      }
    }
  ]
}
```

SQLite schema:
```sql
CREATE TABLE risk_zones (
    zone_id TEXT PRIMARY KEY,
    risk_score FLOAT,
    tier TEXT,
    population INTEGER,
    hazard_type TEXT,
    geometry TEXT
);
```
SQLite has no native polygon type — `geometry` stores the raw GeoJSON geometry object as a TEXT string; the API reassembles the FeatureCollection wrapper around it on read. No spatial indexing needed since this is display-only, not queried spatially.

### `GET /sites`
**Purpose:** candidate relocation sites — where at-risk populations could actually go.

Exported as CSV:
```csv
id,name,lat,lng,capacity,suitability_score,land_cover_type,distance_to_road_km
SITE001,Bahjani Chapori,26.5210,91.0340,1200,7.8,grassland,1.2
SITE002,Kalgachia Uplands,26.4780,90.9560,850,6.5,cropland,2.4
```

SQLite schema:
```sql
CREATE TABLE sites (
    id TEXT PRIMARY KEY,
    name TEXT,
    lat FLOAT,
    lng FLOAT,
    capacity INTEGER,
    suitability_score FLOAT,
    land_cover_type TEXT,
    distance_to_road_km FLOAT
);
```
Confirm with the ML pipeline whether `capacity` is a person-count or household-count — this matters once capacity-matched assignment logic compares it against habitation `population`, and a silent unit mismatch there wouldn't throw an error, it'd just produce wrong assignments.

### `GET /routes?habitation_id=X` (Mode 2)
**Purpose:** shortest path from an at-risk habitation to its assigned site, weighted to avoid other hazard zones — not just shortest distance.

Status: **not finalized.** Route computation itself (graph pathfinding over roads) is cheap enough to run live per-request in the server — this doesn't violate the "no live compute" rule, since it's pathfinding over an already-built graph, not model inference. What genuinely can't be live is the hazard layer feeding into which edges are weighted as dangerous, if that layer depends on live SAR reprocessing. Resolve with the ML pipeline owner: is the road graph (with hazard-weighted edges) handed off once as a static file, or does it need updating mid-event? For the hackathon demo, this may not need to be real-time; the production version (Round 2) should be.

### `POST /reports` / `GET /reports` (Mode 3)
**Purpose:** ground-truth citizen reports during an active disaster — feeding live updates into the routing layer faster than any precomputed hazard data can react (e.g., a bridge that just collapsed).

**Confirmation model — rule-based two-step confirmation or human review.** 
In an actual disaster, waiting for a second independent report within 500m/30min before a hazard shows up anywhere is a bigger risk than trusting a single ground report. Alternative to two-report corroboration:
- A report goes live on the map and into the routing graph **immediately on submission**, marked `unconfirmed` — visually distinct (e.g. dashed/yellow marker vs solid/red), so it's actionable right away but flagged as not yet verified.
- An admin review panel lets a human mark it `confirmed` or `rejected` afterward — this is a credibility/cleanup step, not a gate. Confirmation never blocks the report from affecting routing in the moment.

SQLite schema:
```sql
CREATE TABLE reports (
    id TEXT PRIMARY KEY,
    lat FLOAT,
    lng FLOAT,
    category TEXT,
    description TEXT,
    photo_url TEXT,
    reported_at DATETIME,
    status TEXT,              -- "unconfirmed" | "confirmed" | "rejected"
    reviewed_by TEXT,
    reviewed_at DATETIME
);
```

---

## 4. Build sequence

### Before Day 1 — get organized first
Set up a clean local folder structure before opening QGIS: `/raw_data/<district>/`, `/processed/`, `/qgis_project/`, `/exports/`. Set your QGIS project CRS to EPSG:4326 before loading a single layer. Then work **one layer at a time** — load, clip to your district boundary, visually verify it actually overlaps correctly, only then move to the next layer. Debugging one layer at a time is dramatically easier than discovering a CRS mismatch after five layers are already loaded. Do district boundaries first — they become your clipping mask for everything else — then DEM, then the remaining layers. Fully complete and verify one district before starting the next, rather than processing all 3 in parallel on the first pass.

### Phase 0 — Pre-hackathon prep (do this before the clock starts, if your round allows it)
Downloading and clipping data during the actual 3 days is pure time risk. If your specific SIH round permits pre-work (check the rules — the internal/idea round typically does): pull CartoDEM (East Khasi Hills, West Kameng) and SRTM (Assam), GSI points, Bhuvan layers, WorldCover, population data (GHS-POP or WorldPop), administrative boundaries, and the Geofabrik NE extract for your 3 districts now, clip and reproject them in QGIS, and export clean GeoJSON/GeoParquet files ahead of time. Everything below then starts from clean local data instead of racing to acquire it.

**CartoDEM specifically needs extra lead time here** — unlike SRTM's one-click QGIS plugin, CartoDEM requires creating a Bhuvan account, browsing/ordering tiles per district (not a single bounding-box pull), and manually mosaicking multiple tiles in QGIS if your district spans more than one tile. Do this download during Phase 0, not Day 1 — Bhuvan's portal is also known to be slow/intermittent, so budget real slack time.

**Village/habitation-level coverage should also be checked in Phase 0, before you commit to it as the join layer.** If using DataMeet, verify village-level polygons actually exist and align sensibly. If using `india-geodata`, verify its Census 2011 joins. If coverage is thin, the honest fallback is manual digitization from Census 2011 village shapefiles or Survey of India topo sheets — not a drop-in dataset swap.

### Mode 1 — Long-term (Days 1–2, must complete)

**Day 1 — data acquisition, all in QGIS:**
1. Install QGIS + the SRTM-Downloader, Bhuvan Web Services, and Point Sampling Tool plugins. (CartoDEM isn't served through that plugin; it's pulled manually from Bhuvan's Open Data Archive, ideally already done in Phase 0.)
2. Pull CartoDEM (East Khasi Hills, West Kameng) and SRTM DEM (Assam), GSI landslide points, Bhuvan flood hazard, ESA WorldCover, population, and boundaries for your 3 districts.
3. Open the Geofabrik NE Zone `.osm.pbf` directly in QGIS for roads/hospitals.
4. Reproject everything to EPSG:4326; verify layers visually overlap correctly.

**Day 2 — the model and the score:**
1. Derive slope (QGIS Slope tool), TWI and curvature (SAGA toolbox) from your DEM.
2. Build a labeled point set: GSI landslide points (positive) + randomly sampled stable-area points (negative). Use Point Sampling Tool to attach slope/TWI/curvature/land-cover/elevation to every point — no code.
3. Train with a spatially-blocked split (hold out one whole district as test, not a random split — random splits leak spatial autocorrelation and inflate accuracy):
   ```python
   from sklearn.ensemble import RandomForestClassifier
   from sklearn.metrics import roc_auc_score

   train = points_df[points_df["district"] != "West Kameng"]
   test = points_df[points_df["district"] == "West Kameng"]
   features = ["slope", "twi", "curvature", "landcover_class", "elevation", "dist_to_drainage"]

   rf = RandomForestClassifier(n_estimators=300, random_state=42).fit(train[features], train["is_landslide"])
   print("AUC:", roc_auc_score(test["is_landslide"], rf.predict_proba(test[features])[:,1]))
   ```
4. Extract feature importances, rescale to sum to 1, use as MCDA weights. If your positive-sample count is under ~50, blend these with domain-reasoned weights rather than trusting them outright:
   ```python
   importances = rf.feature_importances_
   weights = dict(zip(features, importances / importances.sum()))
   ```
5. Compute the final Risk Score as a transparent weighted sum in QGIS's Raster Calculator using these weights — this stays fully auditable, which is the whole point.
6. Join population (QGIS "Join Attributes by Location"), compute `Priority = Risk Score × Population Density`, export as ranked CSV + GeoJSON risk zones.
7. Export the final risk zones, ranked habitations, and sites as clean GeoJSON/CSV matching the agreed API contract (Section 3.5) — your teammate's dashboard consumes these directly, built in parallel with your Day 2 work above.

*Checkpoint: a complete, working, defensible submission exists here even if nothing else gets built.*

### Mode 2 — Pre-disaster (attempt after Mode 1 is solid)

1. **Near-real-time flood layer — EOS-04 (RISAT-1A) via Bhoonidhi, preferred.** ISRO's own SAR satellite, explicitly used in published studies for flood mapping in Assam specifically — a strong, genuinely on-point "Indian data" pitch line. **Register at bhoonidhi.nrsc.gov.in on Day 1**, before you need it. Browse & Order the EOS-04 Medium Resolution ScanSAR product for your Assam district, then apply a simple backscatter threshold to extract flood extent.
   **Fallback — Sentinel-1 SAR via Google Earth Engine**, if Bhoonidhi registration doesn't clear in time. Free, runs server-side (no local GPU needed), browser-based:
   ```javascript
   var s1 = ee.ImageCollection('COPERNICUS/S1_GRD')
     .filterBounds(assamDistrictGeometry)
     .filterDate('2026-06-01', '2026-09-01')
     .filter(ee.Filter.eq('instrumentMode', 'IW'));
   var flooded = s1.select('VV').mean().lt(-17);
   Export.image.toDrive({image: flooded, region: assamDistrictGeometry, scale: 10});
   ```
2. **Risk-weighted routing:** build a road graph from your already-downloaded Geofabrik file, not a live API call — keeps this step fully local and demo-safe if venue WiFi is unreliable. Penalize edges that pass through high-risk zones (not just physical distance):
   ```python
   import osmnx as ox, networkx as nx

   # Build from the local Geofabrik extract, not a live Overpass query
   G = ox.graph_from_xml("northeast_zone.osm", network_type="drive")
   for u, v, data in G.edges(data=True):
       data["weight"] = data["length"] * (1 + risk_score_along_edge(u, v))
   route = nx.shortest_path(G, source=habitation_node, target=site_node, weight="weight")
   ```
   (`osmnx` reads `.osm` XML, not `.pbf` directly — convert once with `osmium cat northeast-latest.osm.pbf -o northeast_zone.osm`, a single command-line step, not a recurring dependency.)
3. **Capacity-matched assignment** — sort habitations by priority descending, greedily assign each to its nearest site with remaining capacity, decrementing as you go:
   ```python
   sites = sites_df.set_index("site_id")["capacity"].to_dict()
   assignments = []
   for _, hab in habitations_df.sort_values("priority", ascending=False).iterrows():
       nearest = sorted(sites, key=lambda s: distance(hab, s))
       for site in nearest:
           if sites[site] >= hab["population"]:
               assignments.append((hab["name"], site))
               sites[site] -= hab["population"]
               break
   ```
4. Surface routes and assignments on the dashboard, triggered by the rainfall slider crossing threshold.

### Mode 3 — During-disaster (only if well ahead of schedule)

1. Frontend form (map-click or GPS location, category dropdown, optional photo, auto-timestamp) posting to `POST /reports`.
2. Human review OR rule-based two-step confirmation (simplified from PetaBencana's chatbot flow): a report goes live on the map immediately as `unconfirmed` so it's actionable, then flips to "confirmed" by a human reviewer or if a second report lands nearby.
3. **Any** unconfirmed-or-confirmed road-blockage report removes that edge from the Mode 2 routing graph immediately, triggering automatic reroute — confirmation status affects the marker's visual trust level, not whether it acts. Hazard reports temporarily boost that area's live risk score the same way.
4. For the demo: seed 5–8 realistic sample reports. State plainly that this is a functional pipeline shown with sample data, not a claim of a live user base — that honesty is a strength, not a weakness, in front of judges.

---

## 5. Methodology notes worth stating explicitly in your pitch

- **On existing infrastructure:** you're not duplicating NDEM's hazard-monitoring layer — you're building the relocation-decision layer on top of it. Naming NDEM explicitly shows awareness of the ecosystem rather than reinventing what already works.
- **On "real-time":** honestly means an automated refresh pipeline on the best available public data (SAR re-pulls, rainfall thresholds), not literal live government IoT — no team gets that access in days, and judges know it.
- **On ML vs. MCDA:** the final score is a transparent, auditable formula. ML's role is deriving the weights from real historical data, not replacing the formula with a black box — this matters specifically because the output drives real relocation decisions.
- **On the DEM choice:** CartoDEM (Cartosat-1) was chosen over SRTM specifically because SRTM's C-band radar has known accuracy problems in exactly this terrain type — canopy penetration issues in dense hill forest and layover/shadow artifacts on steep slopes, both common across Meghalaya and Arunachal Pradesh. CartoDEM's stereo-optical derivation avoids those radar-specific artifacts and carries better published vertical accuracy (~8-10m vs. SRTM's ~16m), which matters directly here since slope, TWI, and curvature — the terrain features feeding the RF model and MCDA score — are derivatives of the DEM. An error in the DEM propagates into every hazard score downstream, so this is one of the few dataset swaps with genuine accuracy stakes, not just a sourcing preference.

---

## 6. If you fall behind schedule

- Mode 1 alone is complete and defensible — never sacrifice finishing it for a partial Mode 2.
- Within Mode 2, the SAR flood layer is the single highest-impact addition if you can only do one piece.
- The dashboard is never optional at any mode — an unfinished model with no live way to show it reads as incomplete regardless of backend sophistication.
- **If CartoDEM acquisition (account setup, tile ordering, mosaicking) is eating into Day 1 and Phase 0 pre-work wasn't possible:** fall back to SRTM via QGIS's SRTM-Downloader plugin rather than let the DEM block Mode 1. State the accuracy trade-off honestly if asked rather than silently downgrading — a working Mode 1 on SRTM beats a stalled Mode 1 waiting on Bhuvan.

---

## 7. Dataset list (final)

| Dataset | Source | Format | Used for |
|---|---|---|---|
| Elevation (DEM) — East Khasi Hills & West Kameng | CartoDEM (Cartosat-1) via Bhoonidhi | GeoTIFF, ~32m | Accuracy evidence favors CartoDEM specifically in hilly terrain — used here for the landslide districts |
| Elevation (DEM) — Assam district | SRTM 30m via QGIS's SRTM-Downloader plugin | GeoTIFF | Flat terrain, negligible accuracy gap vs. CartoDEM, simpler to acquire — no reason to add friction here |
| Landslide inventory | GSI NLSM points via Bhuvan Web Services (WMS), or GSI Bhukosh portal export | Vector points | Hazard layer AND ML training labels (Mode 1) — dual purpose |
| Flood hazard (static) | Bhuvan Flood Hazard Zonation Atlas (WMS, no registration); fall back to India-WRIS if Assam isn't covered in the WMS layer | Vector/raster | Assam district hazard layer (Mode 1) |
| Land cover | ESA WorldCover 10m — esa-worldcover.org directly, not via Bhuvan | GeoTIFF | Exposure scoring, relocation-site exclusion (forest/water) |
| Population (regional base layer) | WorldPop 100m grid OR GHS-POP 100m grid | GeoTIFF | Population weighting for the broader regional view/priority ranking |
| Population (top-ranked habitations) | Census 2011 village-level population figures, joined via boundary datasets | CSV + Shapefile | Real named-village precision for the specific habitations shown in the demo/pitch |
| Admin boundaries | `yashveeeeeeer/india-geodata` (merges LGD, Survey of India, Bhuvan, and DataMeet) OR DataMeet (`datta07/INDIAN-SHAPEFILES`) — not GADM | Shapefile | Joining population/risk to habitation-level units |
| Roads/hospitals | PMGSY rural roads + habitation points (primary) OR Geofabrik North-Eastern Zone `.osm.pbf` (supplementary/alternative) | PMGSY: Shapefile/Parquet; OSM: PBF | Site accessibility (Mode 1), routing graph (Mode 2), habitation points |

**Derived locally, not downloaded:** slope, TWI, curvature (from DEM via QGIS/SAGA), distance to drainage/roads (from DEM + Geofabrik layer).
**Live remote source (Mode 2 only):** EOS-04 (ISRO) via Bhoonidhi, preferred. Sentinel-1 via Google Earth Engine as fallback — pre-run and cache locally before demo day rather than querying live.
**Optional:** NASA POWER API for historical rainfall context.

**Why PMGSY over OSM as primary for roads/habitations (if using PMGSY):** OSM is crowdsourced, so completeness tracks contributor density — remote NE hill districts are exactly the kind of area where volunteer mapping can be sparse. PMGSY is field-surveyed by the Ministry of Rural Development specifically for rural connectivity — completeness doesn't depend on whether a volunteer happened to visit.
**Handling OSM's known rural-mapping gap (if using OSM):** OpenStreetMap's own community documentation acknowledges Meghalaya and Arunachal Pradesh are among the least-mapped regions in the country. This matters beyond convenience: your site-accessibility scoring (Mode 1) and risk-weighted routing (Mode 2) both inherit this gap. Name this explicitly in the data-provenance panel and in your pitch.

**Resolution — the honest limiting factor:** your finest layer doesn't set overall precision, your coarsest one does. Population grids (GHS-POP or WorldPop) are the bottleneck for habitation-level ranking specifically — small NE hill hamlets can be smaller than a single grid cell. Primary mitigation: village/habitation-level boundary polygons joined directly to Census population figures, give you real named-village precision for your top-ranked habitations.

**Handling outdated Census data — and why this is a strength, not a weakness, if framed correctly:** Census 2011 is 15 years stale as your population baseline. Don't cite it raw:
1. Use GHS-POP or WorldPop's growth-modeled estimates, which redistribute/project population using satellite-detected built-up area change, not just the frozen 2011 count.
2. Cite the National Commission on Population's official "Population Projections for India and States, 2011–2036" for district-level scaling.
3. For your top-ranked habitations specifically, sanity-check against recent built-up-area growth visible in ESA WorldCover's time series.

India's actual 2021 census was delayed by COVID and is only happening now — house-listing began April 2026, full enumeration results aren't expected until after March 2027. State this explicitly to judges; it converts a weakness into a credibility point.

---

## 8. Final pitch structure & script

### Pitch Structure
**1. Hook (30–45s):** Open with Joshimath — a year after India's most-watched land subsidence crisis, not a single family had been relocated, because nothing connected hazard data to an actionable plan.
**2. The gap (30s):** NDEM already gives India excellent hazard-monitoring infrastructure. This platform builds the decision layer that turns monitoring into action.
**3. What was built (1–2 min):** One platform, three escalating modes — Long-term baseline, Pre-disaster, During-disaster — piloted across 3 districts.
**4. Technical credibility (1 min):** ML-informed but transparent (auditable formula, not a black box), Indian data where it counts (EOS-04, GSI, CartoDEM, Bhuvan), fully offline-capable, and handles India's real population-data currency gap honestly.
**5. Live demo (2–3 min):** Show ranked breakdown, drag rainfall slider, submit/confirm citizen report and watch rerouting.
**6. Scaling roadmap (30s):** Architecture extends directly to full Northeast coverage and coastal states.
**7. Close:** "We didn't rebuild India's hazard data... We built the part that was still missing."

### Full Script (~6–7 min spoken)
*(Working name used below: **Suraksha Setu**)*

**The Hook.** In January 2023, cracks started tearing through the walls of Joshimath, Uttarakhand. Within weeks it became a national crisis — buildings tilting, hundreds of families evacuated overnight. The state government promised relocation within 20 days. A year later, journalists went back. Not a single family had actually been moved. Officials were still saying they needed new surveys to figure out where people could even go. And when one relocation site was finally proposed, residents protested — because there was no transparent case for why that site, and not somewhere else. This wasn't a data problem. ISRO and GSI already had the hazard information the whole time. It was a decision problem.

**The Gap.** India already has excellent hazard-monitoring infrastructure. NDEM, run by ISRO, has been mapping floods, landslides, and other hazards nationally since 2013. We're not here to reinvent that. What we built is the layer that's missing on top of it — the system that turns hazard data into an actual, ranked, ready-to-act-on relocation plan. Before the next Joshimath, not two years after.

**What We Built.** One platform, three modes that scale with urgency. Baseline mode continuously assesses risk and ranks specific habitations — not districts, actual named villages — by urgency, and identifies exactly which relocation sites nearby can take them, with real computed capacity, not guesses. When conditions escalate — rainfall approaching cloudburst intensity — it switches into pre-disaster mode: evacuation routes that specifically avoid other hazard zones, populations matched to shelters respecting real capacity limits. During an unfolding event, it takes in verified citizen reports — a bridge down, a road flooding — the same crowdsourcing approach running successfully in Jakarta for years — and reroutes around what's actually happening, live.

**Technical Credibility.** None of this is a black box — the risk score is a fully transparent formula a district officer can open and understand. We used ML responsibly: trained on GSI's own historical landslide records to learn which factors mattered, using that to set the formula's weights instead of guessing. AI-informed, not AI-obscured. Built on India's own data wherever it counts — EOS-04, GSI, Bhuvan — not borrowed global proxies. Runs fully offline on a single laptop, which matters specifically because these districts lose connectivity exactly when disaster strikes. And we didn't hide that India's population data is still frozen on the 2011 census, since the real 2021 census was delayed by COVID and isn't complete until 2027 — we growth-adjusted our estimates and cross-checked our highest-priority villages against the most current data available, instead of quietly using 15-year-old numbers.

**Live Demo.** Click a flagged habitation — show the hazard breakdown and the generated summary. Drag the rainfall slider past 100mm/hr (India's official cloudburst threshold) — red zones expand live. Submit/confirm a citizen report of a road blockage — evacuation route reroutes automatically.

**Scope, Honestly.** Piloted across 3 districts — East Khasi Hills and West Kameng (landslide), Dhemaji/Barpeta (flood) — chosen because published hazard research already exists for these areas, so every number is checkable against real prior work. Coastal erosion isn't in this pilot because none of these districts have a coastline; the same architecture extends directly to a coastal state.

**Close.** We're not trying to replace what ISRO already built. NDEM already tells India where the hazards are. What we built is the part that was still missing.
