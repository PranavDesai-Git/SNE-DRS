# SIH26191 — Hazard Red Zones & Relocation Platform
### Master Build Guide — 2-District NE Pilot, Round 1 (Meghalaya / Arunachal Pradesh) — Landslide Focus, Flood & Coastal Erosion on Roadmap

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
- Near-real-time satellite monitoring — Phase 2 roadmap item (landslide-equivalent to SAR flood detection would be InSAR ground-deformation or optical change-detection for fresh scarps; not verified/built for Round 1)

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
| Terrain features | QGIS Slope + SAGA toolbox (TWI, curvature) — from CartoDEM for both landslide districts | 1 |
| ML weight derivation | scikit-learn `RandomForestClassifier` on GSI's historical landslide points | 1 |
| Risk scoring | Transparent MCDA weighted sum (weights informed by the RF's feature importances) | 1 |
| Boundaries | `yashveeeeeeer/india-geodata` (merges LGD, Survey of India, Bhuvan, DataMeet) — not GADM, given Arunachal Pradesh border sensitivity | 1 |
| Roads/infra | Geofabrik NE Zone `.osm.pbf`, opened directly in QGIS | 1, 2 |
| Population | WorldPop 100m raster (regional base) + Census 2011 village-level figures via `india-geodata` (top-ranked habitations) | 1 |
| Flood/coastal erosion layers | Not built for Round 1 (landslide-only scope) — Phase 2 roadmap | — |
| Routing | `osmnx` + `networkx`, risk-weighted edge costs | 2 |
| Capacity assignment | Greedy priority-first allocation (Python/pandas) | 2 |
| Citizen reports | React form → FastAPI `POST /reports`, rule-based two-report confirmation | 3 |
| Dashboard | React (Vite) + TailwindCSS + React-Leaflet, served by a FastAPI backend | 1, 2, 3 |

**Team split (2-person):** one person owns the data/ML/GIS pipeline (QGIS, the RF model, MCDA scoring, routing, capacity matching) and exports results to static GeoJSON/CSV files. The other owns FastAPI + React/Leaflet, serving those precomputed files via REST endpoints — FastAPI should serve precomputed data, not compute live on request, so demo-day performance doesn't depend on the GIS pipeline running in the request path. Agree on the exact API contract before either side starts building:
- `GET /habitations` → ranked list (id, name, tier, population, risk_score)
- `GET /risk-zones` → GeoJSON, properties: `risk_score`, `tier`, `population`
- `GET /sites` → candidate relocation sites + capacity
- `GET /routes?habitation_id=X` → Mode 2 routing output
- `POST /reports` / `GET /reports` → Mode 3 citizen reports
- `GET /summary?habitation_id=X` → triggers the template-based Action Center summary

No GPU or CUDA required anywhere. Google Earth Engine's compute happens on Google's servers for free under education use, not on your laptop.

**Local vs. remote, precisely:** everything runs fully offline on your machine once source data is downloaded — QGIS processing, the RF model, MCDA scoring, routing, capacity assignment, and the dashboard itself. With flood dropped from Round 1's scope, there's no remaining remote-compute dependency at all — the whole pipeline is local, which is worth stating plainly in the pitch as a reliability point.

---

## 4. Build sequence

### Before Day 1 — get organized first
Set up a clean local folder structure before opening QGIS: `/raw_data/<district>/`, `/processed/`, `/qgis_project/`, `/exports/`. Set your QGIS project CRS to EPSG:4326 before loading a single layer. Then work **one layer at a time** — load, clip to your district boundary, visually verify it actually overlaps correctly, only then move to the next layer. Debugging one layer at a time is dramatically easier than discovering a CRS mismatch after five layers are already loaded. Do district boundaries first (from `india-geodata`) — they become your clipping mask for everything else — then DEM, then the remaining layers in the order listed in Section 7. Fully complete and verify one district before starting the next, rather than processing all 3 in parallel on the first pass.

### Phase 0 — Pre-hackathon prep (do this before the clock starts, if your round allows it)
Downloading and clipping data during the actual 3 days is pure time risk. If your specific SIH round permits pre-work (check the rules — the internal/idea round typically does): pull CartoDEM, GSI points, WorldCover, WorldPop, `india-geodata` boundaries, and the Geofabrik NE extract for your 2 districts now, clip and reproject them in QGIS, and export clean GeoJSON/GeoParquet files ahead of time. Everything below then starts from clean local data instead of racing to acquire it.

### Mode 1 — Long-term (Days 1–2, must complete)

**Detailed QGIS walkthrough:**

**Part 1 — Install and initial setup**
1. Download QGIS from qgis.org/download — pick the **Long Term Release (LTR)**, not the latest release, for stability.
2. Open QGIS → **Project → New**.
3. Set the project CRS immediately: **Project → Properties → CRS** → search "4326" → select **WGS 84 / EPSG:4326** → OK. Do this before adding any layer.
4. **Project → Save As** → save into `/qgis_project/`, name it `sih26191.qgz`.

**Part 2 — Install plugins**
**Plugins → Manage and Install Plugins**, search and install: **SRTM-Downloader**, **Point Sampling Tool**, **QuickOSM**. For Bhuvan's WMS layers, no plugin needed — **Layer → Add Layer → Add WMS/WMTS Layer → New**, paste Bhuvan's WMS URL, connect.

**Part 3 — Load and split your district boundaries**
1. **Layer → Add Layer → Add Vector Layer** → your `india-geodata` district shapefile (`2011_Dist.shp`).
2. Right-click → **Properties → Information** → confirm CRS is EPSG:4326.
3. Open attribute table, find the district name field, **Select by Expression** → e.g. `"DISTRICT" = 'East Khasi Hills'`.
4. Right-click layer → **Export → Save Selected Features As** → Shapefile → `/processed/east_khasi_hills/boundary.shp`.
5. Repeat for West Kameng.
6. **For terrain-relevant layers specifically (DEM, slope/TWI/curvature), buffer this selection to also include neighboring districts** rather than the tight administrative line — see Part 4 for why. Select by Expression can include an `OR` for adjacent district names (stays in 4326, fine), or use **Vector → Geoprocessing Tools → Buffer** with a real distance value (e.g. 5km) — if using the distance-buffer option specifically, reproject to **EPSG:32646** (UTM Zone 46N, covers both districts) first, buffer there, then reproject the result back to 4326. A degree-based buffer doesn't give a consistent real-world distance across latitude.

**Part 4 — Acquire and clip your DEM**

**Clip to the buffered boundary (district + neighbors), not the tight administrative line.** Slope/TWI/curvature are sensitive to edge effects — clipping exactly to a district's border leaves no terrain context beyond it, so calculations right at that edge become inaccurate, since water flow and terrain shape don't respect administrative lines. Buffering with neighboring districts avoids this. Do one additional tight clip to the exact district boundary only at final export time (Part 11), so your displayed risk map doesn't bleed into neighboring districts you never actually assessed.

- **East Khasi Hills / West Kameng (CartoDEM):** Add Raster Layer → your downloaded CartoDEM GeoTIFF. **Raster → Extraction → Clip Raster by Mask Layer** → Mask = your *buffered* boundary → Run → save to `/processed/<district>/dem.tif`. Repeat for both districts.

**Part 5 — Derive terrain features (from the clipped, buffered DEM) — CRS-critical, read this before running any of these tools**

Slope, TWI, and curvature all measure real ground gradient/distance — computed directly on a 4326 (degree-based) DEM, they silently give wrong values with no error (slope reads artificially near-flat everywhere) or an outright error (GRASS's `r.topidx` refuses lat/long input for TWI). Fix the same way for all three:

1. **Reproject the clipped DEM to EPSG:32646** (UTM Zone 46N — covers both districts): right-click DEM layer → **Export → Save As** → set CRS to EPSG:32646 → save as a new file.
2. **Slope:** Raster → Analysis → Slope → Input = the *reprojected* DEM → Run.
3. **TWI:** Processing Toolbox (Ctrl+Alt+T) → search "Wetness" → SAGA's Topographic Wetness Index (or GRASS's `r.topidx` if SAGA isn't available) → Input = the *reprojected* DEM → Run.
4. **Curvature:** Processing Toolbox → search "Curvature" → SAGA's Curvature tool → Input = the *reprojected* DEM → Run.
5. **Reproject all three outputs (slope, TWI, curvature) back to EPSG:4326** to match the rest of your layers before moving on.

Repeat Parts 4–5 once per district.

**Part 6 — Load hazard layers**
GSI landslide points via the Bhuvan WMS connection, or Add Vector Layer if you exported from Bhukosh directly. Clip to their buffered district boundary (**Vector → Geoprocessing Tools → Clip** — geometric masking, fine in 4326). Bhuvan Flood Hazard Zonation isn't needed for Round 1's landslide-only scope.

**Part 7 — Load land cover and population**
ESA WorldCover GeoTIFF and WorldPop GeoTIFF — Add Raster, clip to each buffered district boundary same as the DEM (masking only, fine in 4326).

**Part 8 — Load roads and habitations (PMGSY primary, Geofabrik supplementary), and compute `dist_to_drainage`**
PMGSY shapefile/Parquet from `india-geodata` — Add Vector Layer (QGIS reads GeoParquet natively in recent versions; use Shapefile if not). Filter/clip same as Part 3 (fine in 4326). Geofabrik's `.pbf` — use **QuickOSM** to query specific tags (`highway=*`, `amenity=hospital`) directly rather than loading the whole file.

**`dist_to_drainage` is one of your actual ML training features (Part 9) and is CRS-critical the same way slope was** — computing it in 4326 silently gives wrong values, same failure mode as slope. Reproject your drainage/stream layer and your points layer to EPSG:32646, run **Vector → Analysis → Distance to Nearest Hub (Line to Hub)** (or equivalent) there, then reproject the distance field's source points back to 4326 to continue.

**Part 9 — Build your ML training points**
1. GSI landslide points (Part 6) = positive samples.
2. **Vector → Research Tools → Random Points in Polygons** → Input = district boundary → generate roughly equal count as negative samples (point generation is containment-based, fine in 4326).
3. Merge both point layers, add field `is_landslide` (1/0) via Field Calculator.
4. **Point Sampling Tool**: Input = merged points, Raster layers to sample = slope, TWI, curvature, land cover, elevation, plus your `dist_to_drainage` field from Part 8 → Run. This step itself doesn't compute distances, just reads existing values, so it's fine as long as the inputs feeding it were correctly computed in Parts 5 and 8.
5. Export result as CSV — this is the file your Random Forest script trains on:
```python
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import roc_auc_score

train = points_df[points_df["district"] != "West Kameng"]
test = points_df[points_df["district"] == "West Kameng"]
features = ["slope", "twi", "curvature", "landcover_class", "elevation", "dist_to_drainage"]

rf = RandomForestClassifier(n_estimators=300, random_state=42).fit(train[features], train["is_landslide"])
print("AUC:", roc_auc_score(test["is_landslide"], rf.predict_proba(test[features])[:,1]))
```
6. Extract feature importances, rescale to sum to 1, use as MCDA weights. If positive-sample count is under ~50, blend with domain-reasoned weights rather than trusting them outright:
```python
importances = rf.feature_importances_
weights = dict(zip(features, importances / importances.sum()))
```

**Part 10 — Back in QGIS: the transparent risk score**
**Raster Calculator** (Raster menu) → build your weighted-sum formula using the ML-derived weights from Part 9, referencing your slope/TWI/curvature/land-cover rasters — this stays fully auditable, which is the whole point of keeping it a formula rather than raw model output.

**Part 11 — Join population and export final files**
1. **Vector → Data Management Tools → Join Attributes by Location** → join risk zones/habitation points to population data.
2. Compute `Priority = Risk Score × Population Density`.
3. **Do the tight clip to the exact (unbuffered) district boundary now** — this is the display/reporting extent, separate from the buffered computation extent used in Parts 4–9.
4. Export final layers: right-click each → **Export → Save Features As → GeoJSON**, matching property names in your API contract (Section 3) — your teammate's React/FastAPI dashboard consumes these directly.

*Checkpoint: a complete, working, defensible submission exists here even if nothing else gets built.*

### Mode 2 — Pre-disaster (attempt after Mode 1 is solid)

1. **Near-real-time satellite monitoring — Phase 2 roadmap, not built for Round 1.** With flood dropped from Round 1's scope, EOS-04/Sentinel-1 SAR isn't needed. If extending to landslide-specific near-real-time monitoring later, InSAR ground-deformation or optical change-detection for fresh scarps is the equivalent technique — not verified for free accessibility yet, treat as something to check into for Phase 2 rather than commit build time to now.
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

1. React form (map-click or GPS location, category dropdown, optional photo, auto-timestamp) posting to `POST /reports`.
2. Rule-based two-step confirmation (simplified from PetaBencana's chatbot flow, achievable without one): a report is "pending" until a second independent report lands within ~500m and ~30 minutes, then flips to "confirmed."
3. Confirmed road-blockage reports remove that edge from the Mode 2 routing graph, triggering automatic reroute. Confirmed hazard reports temporarily boost that area's live risk score.
4. For the demo: seed 5–8 realistic sample reports. State plainly that this is a functional pipeline shown with sample data, not a claim of a live user base — that honesty is a strength, not a weakness, in front of judges.

---

## 5. Methodology notes worth stating explicitly in your pitch

- **On existing infrastructure:** you're not duplicating NDEM's hazard-monitoring layer — you're building the relocation-decision layer on top of it. Naming NDEM explicitly shows awareness of the ecosystem rather than reinventing what already works.
- **On "real-time":** honestly means an automated refresh pipeline on the best available public data (SAR re-pulls, rainfall thresholds), not literal live government IoT — no team gets that access in days, and judges know it.
- **On ML vs. MCDA:** the final score is a transparent, auditable formula. ML's role is deriving the weights from real historical data, not replacing the formula with a black box — this matters specifically because the output drives real relocation decisions.

---

## 6. If you fall behind schedule

- Mode 1 alone is complete and defensible — never sacrifice finishing it for a partial Mode 2.
- Within Mode 2, risk-weighted routing and capacity-matched assignment are the highest-impact additions since Mode 1 finished — build both before attempting Mode 3.
- The dashboard is never optional at any mode — an unfinished model with no live way to show it reads as incomplete regardless of backend sophistication.
- **Cut Mode 2/3 features before cutting a landslide district.** Dropping to one landslide district breaks the spatially-blocked train/test split (Section 4, Part 9) — your core ML rigor point — forcing a weaker within-district sub-area holdout instead of a genuine held-out district. That's a real methodological cost, not just less work, and it's worth protecting over additional polish elsewhere.

---

## 7. Dataset list (final)

| Dataset | Source | Format | Used for |
|---|---|---|---|
| Elevation (DEM) — East Khasi Hills & West Kameng | CartoDEM (Cartosat-1) via Bhoonidhi | GeoTIFF, ~32m | Accuracy evidence favors CartoDEM specifically in hilly terrain — used for both Round 1 districts |
| Landslide inventory | GSI NLSM points via Bhuvan Web Services (WMS), or GSI Bhukosh portal export | Vector points | Hazard layer AND ML training labels (Mode 1) — dual purpose |
| Land cover | ESA WorldCover 10m — esa-worldcover.org directly, not via Bhuvan | GeoTIFF | Exposure scoring, relocation-site exclusion (forest/water) |
| Population (regional base layer) | WorldPop 100m grid | GeoTIFF | Population weighting for the broader regional view only — see resolution note below |
| Population (top-ranked habitations) | Census 2011 village-level population figures, joined via `india-geodata`'s village boundaries | CSV + Shapefile | Real named-village precision for the specific habitations shown in the demo/pitch |
| Admin boundaries | `yashveeeeeeer/india-geodata` (merges LGD, Survey of India, Bhuvan, and DataMeet) — not GADM, given Arunachal Pradesh border sensitivity. Census 2011 district shapefile confirmed working: `.shp/.shx/.dbf/.prj`, CRS is WGS84/EPSG:4326, matching every other layer in this build | Shapefile | Joining population/risk to habitation-level units |
| Roads/hospitals | PMGSY rural roads + habitation points (field-surveyed, via `india-geodata` or `geosadak-pmgsy.nic.in/opendata` directly) — primary; Geofabrik North-Eastern Zone `.osm.pbf` as supplementary for larger-settlement infrastructure | PMGSY: Shapefile/Parquet; OSM: PBF | Site accessibility (Mode 1), routing graph (Mode 2), habitation points |

**Not needed for Round 1 (landslide-only scope):** Bhuvan Flood Hazard Zonation Atlas / India-WRIS, EOS-04/Sentinel-1 SAR, SRTM (CartoDEM covers both remaining districts). Kept as Phase 2 roadmap items — see hazard-types note below.

**If time allows — re-adding flood as a Day 3 stretch, not a default:** roughly 6–10 hours of combined work (Assam district through the same established pipeline, Bhuvan flood WMS layer, one SAR flood-extent step, one more API endpoint/map layer for your teammate). Use **Sentinel-1/GEE, not EOS-04/Bhoonidhi**, if attempting this — Bhoonidhi's unknown registration-approval turnaround is the one real open-ended risk in this estimate, and you've already hit friction with its 5-year date-range cap this session. Only attempt this after Mode 1 and Mode 2 are solid for the 2 landslide districts — never let it compete with finishing those.

**Derived locally, not downloaded:** slope, TWI, curvature (from DEM via QGIS/SAGA), distance to drainage/roads (from DEM + Geofabrik layer).

**Why PMGSY over OSM as primary for roads/habitations:** OSM is crowdsourced, so completeness tracks contributor density — remote NE hill districts are exactly the kind of area where volunteer mapping can be sparse. PMGSY is field-surveyed by the Ministry of Rural Development specifically for rural connectivity — 1M+ habitation points, 2.5M+ km of rural roads, ~770,000 geo-tagged rural facilities — completeness doesn't depend on whether a volunteer happened to visit. Already available through the same `india-geodata` repo used for boundaries.

**Land cover currency, worth naming if asked:** ESA WorldCover is a 2021 snapshot (76.7% independently validated global accuracy, no India-specific breakdown published) — 5 years old as of this build. Land cover changes faster than terrain; this is a real, secondary staleness point, smaller than the census gap but worth acknowledging rather than implying continuous freshness.

**Optional:** NASA POWER API for historical rainfall context, to sanity-check the slider's cloudburst marker — not required since the slider is a simulated input.

**Resolution — the honest limiting factor:** your finest layer doesn't set overall precision, your coarsest one does. WorldPop's 100m population grid is the bottleneck for habitation-level ranking specifically — small NE hill hamlets can be smaller than a single grid cell, so nearby villages can blur together. Primary mitigation: `india-geodata`'s village/habitation-level boundary polygons, joined directly to Census population figures, give you real named-village precision for your top-ranked habitations instead of relying on the WorldPop grid there — use WorldPop only as the base layer for the broader regional view.

**Hazard types — Round 1 scope, and what's on the roadmap:** the PS scope mentions landslide, flood, coastal erosion, and cloudburst. **Round 1 covers landslide only** (East Khasi Hills, West Kameng), built with real rigor — the spatially-blocked train/test split, ML-informed weights, full Mode 1–3 architecture. Cloudburst is still handled within this scope, as the dynamic escalation trigger (the rainfall slider crossing IMD's 100mm/hr threshold) — rainfall is a genuine landslide trigger, not just a flood one, so this isn't a stretch. Flood and coastal erosion are explicitly Phase 2 roadmap items, not silently dropped — state this plainly if asked: Round 1 concentrates depth and validation on one hazard type rather than spreading thin across four; the same architecture (Modes 1–3, the scoring methodology, the routing/capacity logic) extends directly to flood and coastal states in the next phase.

**Handling outdated Census data — and why this is a strength, not a weakness, if framed correctly:** Census 2011 is 15 years stale as your population baseline. Don't cite it raw:
1. Use WorldPop's growth-modeled estimates, which redistribute/project population using satellite-detected built-up area change, not just the frozen 2011 count.
2. Cite the National Commission on Population's official "Population Projections for India and States, 2011–2036" for district-level scaling — a government-published number, defensible if a judge asks about currency.
3. For your top-ranked habitations specifically, sanity-check against recent built-up-area growth visible in ESA WorldCover's time series.

**The framing that matters:** India's actual 2021 census was delayed by COVID and is only happening now — house-listing began April 2026, full enumeration results aren't expected until after March 2027. As of your build, Census 2011 genuinely is the most current official population data that exists for India — this isn't a shortcut you're taking, it's the real state of India's data infrastructure, and growth-adjusting it transparently is more responsible than silently citing a stale number. State this explicitly to judges; it converts a weakness into a credibility point.

---

## 8. Final pitch — full script (~6–7 min spoken)

Working name used below: **Suraksha Setu** ("safety bridge") — swap for whatever name you settle on.

**The Hook.** In January 2023, cracks started tearing through the walls of Joshimath, Uttarakhand. Within weeks it became a national crisis — buildings tilting, hundreds of families evacuated overnight. The state government promised relocation within 20 days. A year later, journalists went back. Not a single family had actually been moved. Officials were still saying they needed new surveys to figure out where people could even go. And when one relocation site was finally proposed, residents protested — because there was no transparent case for why that site, and not somewhere else. This wasn't a data problem. ISRO and GSI already had the hazard information the whole time. It was a decision problem.

**The Gap.** India already has excellent hazard-monitoring infrastructure. NDEM, run by ISRO, has been mapping floods, landslides, and other hazards nationally since 2013. We're not here to reinvent that. What we built is the layer that's missing on top of it — the system that turns hazard data into an actual, ranked, ready-to-act-on relocation plan. Before the next Joshimath, not two years after.

**What We Built.** One platform, three modes that scale with urgency. Baseline mode continuously assesses risk and ranks specific habitations — not districts, actual named villages — by urgency, and identifies exactly which relocation sites nearby can take them, with real computed capacity, not guesses. When conditions escalate — rainfall approaching cloudburst intensity — it switches into pre-disaster mode: evacuation routes that specifically avoid other hazard zones, populations matched to shelters respecting real capacity limits. During an unfolding event, it takes in verified citizen reports — a bridge down, a road flooding — the same crowdsourcing approach running successfully in Jakarta for years — and reroutes around what's actually happening, live.

**Technical Credibility.** None of this is a black box — the risk score is a fully transparent formula a district officer can open and understand. We used ML responsibly: trained on GSI's own historical landslide records to learn which factors mattered, using that to set the formula's weights instead of guessing. AI-informed, not AI-obscured. Built on India's own data wherever it counts — CartoDEM, GSI, Bhuvan — not borrowed global proxies. Runs fully offline on a single laptop, which matters specifically because these districts lose connectivity exactly when disaster strikes. And we didn't hide that India's population data is still frozen on the 2011 census, since the real 2021 census was delayed by COVID and isn't complete until 2027 — we growth-adjusted our estimates and cross-checked our highest-priority villages against the most current data available, instead of quietly using 15-year-old numbers.

**Live Demo.** Click a flagged habitation — show the hazard breakdown and the generated summary. Drag the rainfall slider past 100mm/hr (India's official cloudburst threshold) — red zones expand live. Submit/confirm a citizen report of a road blockage — evacuation route reroutes automatically.

**Scope, Honestly.** Round 1 is piloted across 2 districts — East Khasi Hills and West Kameng, both landslide risk — chosen because published hazard research already exists for these areas, so every number is checkable against real prior work. We deliberately concentrated depth on one hazard type rather than spreading across four: the spatially-blocked validation, the ML-informed scoring, the full three-mode architecture are all built with real rigor here. Flood and coastal erosion are explicit next-phase roadmap items — the same architecture extends directly to them.

**Close.** We're not trying to replace what ISRO already built. NDEM already tells India where the hazards are. What we built is the part that was still missing.
