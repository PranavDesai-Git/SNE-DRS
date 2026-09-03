# SIH26191 — Hazard Red Zones & Relocation Platform
### Master Build Guide — 3-District NE Pilot (Meghalaya / Arunachal Pradesh / Assam)

---

## 1. The problem, precisely

India's disaster-hazard data infrastructure is genuinely mature — NDEM (ISRO/NRSC, operational since 2013) already does multi-hazard zonation and real-time monitoring at national scale. What doesn't exist, evidenced concretely by Joshimath (a nationally-watched crisis where, a year in, zero families had been relocated, site selection was ad hoc, and one candidate site provoked resident protest for lack of a transparent case) is the **decision layer**: turning hazard data into a ranked, actionable, capacity-aware relocation plan before a disaster forces the issue.

That decision layer is what this platform is.

---

## 2. Product structure — one platform, three escalating modes

Not three separate apps. One shared data backbone (hazard scores, population, candidate relocation sites), three views that activate progressively. Each mode is a **complete, valid submission on its own** — this is also literally the build order.

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

---

## 3. Tech stack

| Component | Tool | Mode |
|---|---|---|
| GIS processing | QGIS (Bhuvan Web Services, Point Sampling Tool plugins) | 1 |
| Elevation source | CartoDEM (Cartosat-1, via Bhuvan Open Data Archive, 30m, ~8-10m vertical accuracy) | 1 |
| Terrain features | QGIS Slope + SAGA toolbox (TWI, curvature) | 1 |
| ML weight derivation | scikit-learn `RandomForestClassifier` on GSI's historical landslide points | 1 |
| Risk scoring | Transparent MCDA weighted sum (weights informed by the RF's feature importances) | 1 |
| Boundaries | DataMeet (`datta07/INDIAN-SHAPEFILES`) — not GADM, given Arunachal Pradesh border sensitivity | 1 |
| Roads/infra | Geofabrik NE Zone `.osm.pbf`, opened directly in QGIS | 1, 2 |
| Population | GHS-POP 100m grid (JRC/Copernicus) — chosen over WorldPop on India-specific accuracy evidence | 1 |
| Flood hazard (static) | Bhuvan Flood Hazard Zonation Atlas / India-WRIS (registration required — register early if used) | 1 |
| Flood hazard (near-real-time) | Sentinel-1 SAR via Google Earth Engine (free, runs server-side, no local GPU) | 2 |
| Routing | `osmnx` + `networkx`, risk-weighted edge costs | 2 |
| Capacity assignment | Greedy priority-first allocation (Python/pandas) | 2 |
| Citizen reports | Svelte form → Go `POST /reports`, rule-based two-report confirmation | 3 |
| Dashboard | Svelte + Leaflet, served by a Go backend | 1, 2, 3 |
| Data storage | SQLite, populated by a local import script — not a live cross-machine connection | 1, 2, 3 |

**Team split (2-person):** one person owns the data/ML/GIS pipeline (QGIS, the RF model, MCDA scoring, routing, capacity matching) and exports results to static GeoJSON/CSV files. The other owns the backend + frontend (this pilot uses Go + Svelte/Leaflet in place of FastAPI + React, the contract below is framework-agnostic), serving that data via REST endpoints. **Data flow, precisely:** the GIS/ML teammate exports fresh GeoJSON/CSV files whenever they recompute something and hands them off (shared repo, Slack, synced folder — anything that isn't a live connection); the backend teammate runs a small **import script** that loads those files into a local **SQLite** database, and the API queries SQLite rather than re-parsing files per request. This gets you real `WHERE`/`ORDER BY`/`JOIN` ergonomics for things like sorting `/habitations` by risk score, without either machine needing to be network-reachable during dev or demo. Explicitly avoid: a live database server (MySQL/Postgres) with the compute machine pushing writes to the backend machine over the network mid-demo — it adds a cross-laptop dependency that can fail in front of judges for no real benefit over the file → import script → SQLite flow. The backend should serve precomputed data, not compute live on request, so demo-day performance doesn't depend on the GIS pipeline running in the request path. Agree on the exact API contract before either side starts building:
- `GET /habitations` → ranked list (id, name, tier, population, risk_score)
- `GET /risk-zones` → GeoJSON, properties: `risk_score`, `tier`, `population`
- `GET /sites` → candidate relocation sites + capacity
- `GET /routes?habitation_id=X` → Mode 2 routing output
- `POST /reports` / `GET /reports` → Mode 3 citizen reports

No GPU or CUDA required anywhere. Google Earth Engine's compute happens on Google's servers for free under education use, not on your laptop.

**Local vs. remote, precisely:** everything runs fully offline on your machine once source data is downloaded (QGIS processing, the RF model, MCDA scoring, routing, capacity assignment, and the Streamlit dashboard itself) — the one genuine exception is the SAR flood analysis, which runs on Google Earth Engine's servers by design. For demo day, **pre-run the GEE export ahead of time and store the result locally** rather than querying it live in front of judges — removes venue-WiFi risk from your only remote-dependent step, and the dashboard itself needs no internet to run either way.

---

## 3.5 API data contract (finalized)

This is the exact file format and SQLite schema for each endpoint, agreed between the GIS/ML teammate (exports these files) and the backend teammate (imports them into SQLite, serves them via Go). The GIS/ML teammate should export to these formats directly — column names and structure are locked, not placeholders.

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

Status: **not finalized.** Route computation itself (graph pathfinding over roads) is cheap enough to run live per-request in the Go server — this doesn't violate the "no live compute" rule, since it's pathfinding over an already-built graph, not model inference. What genuinely can't be live is the hazard layer feeding into which edges are weighted as dangerous, if that layer depends on live SAR reprocessing. Resolve with the ML pipeline owner: is the road graph (with hazard-weighted edges) handed off once as a static file, or does it need updating mid-event? For the hackathon demo, this may not need to be real-time; the production version (Round 2) should be.

### `POST /reports` / `GET /reports` (Mode 3)
**Purpose:** ground-truth citizen reports during an active disaster — feeding live updates into the routing layer faster than any precomputed hazard data can react (e.g., a bridge that just collapsed).

**Confirmation model — deliberately not the two-report-corroboration approach from the original PetaBencana-style design.** In an actual disaster, waiting for a second independent report within 500m/30min before a hazard shows up anywhere is a bigger risk than trusting a single ground report — the failure mode is a vehicle driving into a hazard that a human already reported but the system hadn't surfaced yet. Instead:
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

### Phase 0 — Pre-hackathon prep (do this before the clock starts, if your round allows it)
Downloading and clipping data during the actual 3 days is pure time risk. If your specific SIH round permits pre-work (check the rules — the internal/idea round typically does): pull CartoDEM, GSI points, Bhuvan layers, WorldCover, GHS-POP, DataMeet boundaries, and the Geofabrik NE extract for your 3 districts now, clip and reproject them in QGIS, and export clean GeoJSON/GeoParquet files ahead of time. Everything below then starts from clean local data instead of racing to acquire it.

**CartoDEM specifically needs extra lead time here** — unlike SRTM's one-click QGIS plugin, CartoDEM requires creating a Bhuvan account, browsing/ordering tiles per district (not a single bounding-box pull), and manually mosaicking multiple tiles in QGIS if your district spans more than one tile. Do this download during Phase 0, not Day 1 — Bhuvan's portal is also known to be slow/intermittent, so budget real slack time.

**DataMeet's village/habitation-level coverage should also be checked in Phase 0, before you commit to it as the join layer.** Its state/district boundaries are solid, but the village-level layer is community-maintained, self-described as incomplete, and hasn't been independently verified for coverage — this could matter specifically for Meghalaya/Arunachal Pradesh, since volunteer digitization skews toward better-mapped states. Open the layer for your 3 districts and visually confirm habitation-level polygons actually exist and align sensibly before building the join pipeline on top of it. If coverage is thin, the honest fallback is manual digitization from Census 2011 village shapefiles or Survey of India topo sheets — not a drop-in dataset swap.

### Mode 1 — Long-term (Days 1–2, must complete)

**Day 1 — data acquisition, all in QGIS:**
1. Install QGIS + the Bhuvan Web Services and Point Sampling Tool plugins. (No SRTM-Downloader needed — CartoDEM isn't served through that plugin; it's pulled manually from Bhuvan's Open Data Archive, ideally already done in Phase 0.)
2. If not already done in Phase 0: log into Bhuvan's Open Data Archive, select Cartosat-1 → CartoDEM Versions, download tiles covering each of your 3 districts, and mosaic (Raster → Miscellaneous → Merge) if a district spans multiple tiles. Pull GSI landslide points, Bhuvan flood hazard, ESA WorldCover, GHS-POP, and DataMeet boundaries for your 3 districts (see prior conversation for exact sources — datasets get finalized in the next step now that this architecture is fixed).
3. Open the Geofabrik NE Zone `.osm.pbf` directly in QGIS for roads/hospitals.
4. Reproject everything to EPSG:4326 (CartoDEM tiles may come in a different native projection/tiling scheme than SRTM — check this explicitly); verify layers visually overlap correctly.

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
7. Export the final risk zones, ranked habitations, and sites as clean GeoJSON/CSV matching the agreed API contract (Section 3) — your teammate's backend imports these into SQLite via their import script, dashboard work happening in parallel with your Day 2 work above.

*Checkpoint: a complete, working, defensible submission exists here even if nothing else gets built.*

### Mode 2 — Pre-disaster (attempt after Mode 1 is solid)

1. **Near-real-time flood layer — EOS-04 (RISAT-1A) via Bhoonidhi, preferred.** ISRO's own SAR satellite, explicitly used in published studies for flood mapping in Assam specifically — a strong, genuinely on-point "Indian data" pitch line. **Register at bhoonidhi.nrsc.gov.in on Day 1**, before you need it — registration turnaround is the one real unknown here. Browse & Order the EOS-04 Medium Resolution ScanSAR product for your Assam district, then apply a simple backscatter threshold (same logic as the Sentinel-1 approach below) to extract flood extent — don't chase the more sophisticated published ML methods, a basic threshold is enough for a working demo.

   **If you also plan to use India-WRIS as your Section 7 static-flood-hazard fallback:** its gauge-station water-level data requires a registered account, not open download — register at indiawris.gov.in during Phase 0/Day 1 alongside the Bhoonidhi registration, same unknown-turnaround risk category.

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

1. Svelte form (map-click or GPS location, category dropdown, optional photo, auto-timestamp) posting to `POST /reports`.
2. Human review, not auto-confirmation: a report goes live on the map and into the routing graph immediately as `unconfirmed` (visually distinct marker), so it's actionable the moment it's submitted — no waiting on a second corroborating report, since that delay is the actual danger in a live event (see Section 3.5 for the reasoning). An admin panel lets a person mark it `confirmed`/`rejected` afterward as a credibility step, not a gate.
3. **Any** unconfirmed-or-confirmed road-blockage report removes that edge from the Mode 2 routing graph immediately, triggering automatic reroute — confirmation status affects the marker's visual trust level, not whether it acts. Hazard reports temporarily boost that area's live risk score the same way, on submission.
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
| Elevation (DEM) | CartoDEM (Cartosat-1) 30m via Bhuvan Open Data Archive | GeoTIFF | Slope/TWI/curvature derivation, site suitability (Mode 1) |
| Landslide inventory | GSI NLSM points via Bhuvan Web Services (WMS), or GSI Bhukosh portal export | Vector points | Hazard layer AND ML training labels (Mode 1) — dual purpose |
| Flood hazard (static) | Bhuvan Flood Hazard Zonation Atlas (WMS, no registration); fall back to India-WRIS if Assam isn't covered in the WMS layer — India-WRIS gauge data requires a registered account, register early (Phase 0) if this fallback is needed | Vector/raster | Assam district hazard layer (Mode 1) |
| Land cover | ESA WorldCover 10m — esa-worldcover.org directly, not via Bhuvan | GeoTIFF | Exposure scoring, relocation-site exclusion (forest/water) |
| Population | GHS-POP 100m grid (JRC/Copernicus) | GeoTIFF | Population weighting for priority ranking, site capacity |
| Admin boundaries | DataMeet (`datta07/INDIAN-SHAPEFILES`) — not GADM, given Arunachal Pradesh border sensitivity; verify village-level coverage for your districts before committing | Shapefile | Joining population/risk to habitation-level units |
| Roads/hospitals | Geofabrik North-Eastern Zone `.osm.pbf` — known rural-coverage gap in NE India, see note below | PBF | Site accessibility (Mode 1), routing graph (Mode 2) |

**Derived locally, not downloaded:** slope, TWI, curvature (from DEM via QGIS/SAGA), distance to drainage/roads (from DEM + Geofabrik layer).

**Live remote source (Mode 2 only):** EOS-04 (ISRO) via Bhoonidhi, preferred — register early given unknown turnaround time. Sentinel-1 via Google Earth Engine as fallback — pre-run and cache locally before demo day rather than querying live.

**Optional:** NASA POWER API for historical rainfall context, to sanity-check the slider's cloudburst marker — not required since the slider is a simulated input.

**Resolution — the honest limiting factor:** your finest layer doesn't set overall precision, your coarsest one does. GHS-POP's 100m population grid is the bottleneck for habitation-level ranking specifically — small NE hill hamlets can be smaller than a single grid cell, so nearby villages can blur together. GHS-POP was chosen over WorldPop specifically because an India-focused ground-truth study (IIHS, ~600,000 villages/towns) found it outperforms WorldPop, GPW, and HRSL across nearly all Indian states — WorldPop has documented underestimation for small rural settlements, which is close to the worst case for hill hamlets this size. Still cross-check your top-ranked habitations (the ones that actually appear in the demo/pitch) against Census 2011's village-level population tables, which give exact named-village counts rather than a modeled estimate.

**Hazard types covered, and why one isn't:** the PS scope mentions landslide, flood, coastal erosion, and cloudburst. This pilot covers landslide (Meghalaya, Arunachal Pradesh) and flood (Assam) directly; cloudburst is handled as the dynamic escalation trigger (the rainfall slider crossing IMD's 100mm/hr threshold) rather than a separate static layer. Coastal erosion is deliberately out of scope — all 3 pilot districts are inland/hill terrain with no coastline. Have this answer ready verbatim if asked: pilot districts are inland, coastal erosion doesn't apply there, the same architecture would extend to it for a coastal pilot state.

**Handling outdated Census data — and why this is a strength, not a weakness, if framed correctly:** Census 2011 is 15 years stale as your population baseline. Don't cite it raw:
1. Use GHS-POP's growth-modeled estimates, which derive population from satellite-detected built-up area change, not just the frozen 2011 count.
2. Cite the National Commission on Population's official "Population Projections for India and States, 2011–2036" for district-level scaling — a government-published number, defensible if a judge asks about currency.
3. For your top-ranked habitations specifically, sanity-check against recent built-up-area growth visible in ESA WorldCover's time series.

**Handling OSM's known rural-mapping gap in NE India — state this proactively, don't wait to be asked:** OpenStreetMap's own community documentation acknowledges Meghalaya and Arunachal Pradesh are among the least-mapped regions in the country, with rural roads often missing even where they physically exist — this is exactly why a dedicated OSM North East India mapping community had to be formed in 2023 to address it. This matters beyond convenience: your site-accessibility scoring (Mode 1) and risk-weighted routing (Mode 2) both inherit this gap, so a habitation could be marked poorly-connected or routed around incorrectly simply because the road wasn't digitized, not because it's actually unreachable. This is a silent failure mode — nothing in the output flags it as a data gap rather than ground truth. Name this explicitly in the data-provenance panel and in your pitch: a production deployment would close this gap with BRO (Border Roads Organisation) or state PWD road records, which this pilot doesn't have access to. Stating it upfront is stronger than having a judge notice a suspiciously sparse road network and ask why.

**The framing that matters:** India's actual 2021 census was delayed by COVID and is only happening now — house-listing began April 2026, full enumeration results aren't expected until after March 2027. As of your build, Census 2011 genuinely is the most current official population data that exists for India — this isn't a shortcut you're taking, it's the real state of India's data infrastructure, and growth-adjusting it transparently is more responsible than silently citing a stale number. State this explicitly to judges; it converts a weakness into a credibility point.

---

## 8. Final pitch structure

**1. Hook (30–45s):** Open with Joshimath — a year after India's most-watched land subsidence crisis, not a single family had been relocated, not because the hazard was unknown (GSI/ISRO had the data throughout) but because nothing connected that data to an actual, capacity-checked, ranked relocation plan.

**2. The gap (30s):** NDEM already gives India excellent hazard-monitoring infrastructure. This platform doesn't duplicate that — it builds the decision layer that turns monitoring into action.

**3. What was built (1–2 min):** One platform, three escalating modes — Long-term baseline (ML-informed, fully auditable MCDA scoring, ranked habitations, relocation-site suitability), Pre-disaster (risk-weighted evacuation routing, capacity-matched displacement assignment), During-disaster (verified citizen reports rerouting around confirmed blockages) — piloted across 3 districts chosen specifically because real published hazard studies exist to validate against.

**4. Technical credibility (1 min):**
- ML-informed but transparent: a Random Forest trained on GSI's own historical landslide inventory derives the scoring weights, but the final output is a fully auditable formula — not a black box a government official has to trust blindly
- Genuinely Indian data where it counts: EOS-04, GSI, CartoDEM, Bhuvan — not just global proxies standing in, and the DEM choice specifically (CartoDEM over SRTM) was made on accuracy grounds for this terrain, not sourcing preference alone
- Fully offline-capable, zero paid cloud dependency — matters specifically because these are the exact regions that lose connectivity when disaster strikes
- Handles India's real population-data currency gap honestly (2011 census, 2027 replacement still pending) instead of pretending the problem doesn't exist
- Discloses the one known data limitation upfront: OSM's rural road coverage in NE India has real gaps, which the accessibility/routing layers inherit — named proactively, with BRO/state PWD data identified as the production-grade fix

**5. Live demo (2–3 min):** click a habitation → ranked hazard breakdown and matched relocation site appear; drag the rainfall slider → watch red zones dynamically escalate past the IMD cloudburst threshold; submit/confirm a citizen report → watch the evacuation route reroute around it live.

**6. Scaling roadmap, stated honestly (30s):** 3-district pilot validated against published hazard research; architecture extends directly to full Northeast coverage and to coastal states for the erosion module not included in this pilot.

**7. Close:** "We didn't rebuild India's hazard data — NDEM already does that well. We built the part that was still missing."
