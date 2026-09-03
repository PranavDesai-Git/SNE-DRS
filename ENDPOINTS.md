**GET /habitations**
Purpose: A ranked list of villages that are most aat risk

example file:
``` csv
id,name,district,lat,lng,population,risk_score,tier,slope_score,twi_score,landcover_score,rainfall_score
HAB001,Sorbhog,Barpeta,26.4521,90.9871,3200,8.7,Immediate,0.82,0.75,0.60,0.90
HAB002,Bahari,Barpeta,26.3890,91.0120,1850,6.2,Short-term,0.55,0.60,0.45,0.70
```



``` sql
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


**GET /risk-zones**
Purpose: Visualization on the map. mapped to the habitations. its a polegon zone

example file:
``` json
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
        "coordinates": [
          [
            [90.9850, 26.4500],
            [90.9900, 26.4500],
            [90.9900, 26.4550],
            [90.9850, 26.4550],
            [90.9850, 26.4500]
          ]
        ]
      }
    }
  ]
}
```


``` sql
CREATE TABLE risk_zones (
    zone_id TEXT PRIMARY KEY,
    risk_score FLOAT,
    tier TEXT,
    population INTEGER,
    hazard_type TEXT,
    geometry TEXT
);
```

**GET /sites**
Purpose: Safe site areas.


example file:
```csv
id,name,lat,lng,capacity,suitability_score,land_cover_type,distance_to_road_km
SITE001,Bahjani Chapori,26.5210,91.0340,1200,7.8,grassland,1.2
SITE002,Kalgachia Uplands,26.4780,90.9560,850,6.5,cropland,2.4
```


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

## mode 2
**GET /routes?habitation_id=X**
Purpose: Shortest path that avoids hazards

////TBD MIGHT BE REALTIME, CANT PRECOMPUTE ALL THE PATHS NEED TO REACT TO ENVIRONMENT.
//// MIGHT NOT BE REALTIME IN DEMO WILL BE REALTIME IN THE ACTUAL THING (ROUND 2)

## mode 3
**POST /reports**

**GET /reports**

Purpose: during an actual live disaster, 
static precomputed data goes stale fast.
a road that was fine 10 minutes ago might 
now be blocked by a landslide. This is ground 
truth from actual humans on the ground, feeding back 
into the system live, so the routing (/routes) can reroute 
around a road that JUST got blocked. It's the only endpoint 
where real-time human input changes the app's behavior mid-disaster.

/// JUST SENDING AN SQL QUERY IF SOMETHING IS MARKED AS A DANGER ZONE ONCE CONFIRMED BY A PERSON
/// WILL BE AN ADMIN PANEL WHERE A PERSON CAN APPROVE THIS SHIT

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

