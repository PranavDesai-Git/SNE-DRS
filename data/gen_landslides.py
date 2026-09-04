import json
import random
import math

filepath = '/home/host0/Projects/GO/SNE-DRS/data/risk-zones.geojson'

with open(filepath, 'r') as f:
    data = json.load(f)

# Bounding box roughly around Guwahati / Assam area mapped
LAT_MIN, LAT_MAX = 26.1, 26.6
LNG_MIN, LNG_MAX = 90.5, 91.8

num_landslides = 15
tiers = [("Immediate", 8.5, 9.9), ("Short-term", 6.0, 8.4), ("Medium-term", 4.0, 5.9)]

for i in range(num_landslides):
    center_lat = random.uniform(LAT_MIN, LAT_MAX)
    center_lng = random.uniform(LNG_MIN, LNG_MAX)
    
    tier_info = random.choice(tiers)
    tier_name = tier_info[0]
    risk_score = round(random.uniform(tier_info[1], tier_info[2]), 1)
    
    num_points = random.randint(7, 14)
    coordinates = []
    
    for j in range(num_points):
        angle = (j / num_points) * 2 * math.pi
        # Irregular radius to make jagged shapes
        radius = random.uniform(0.005, 0.02)
        
        pt_lat = center_lat + radius * math.sin(angle)
        pt_lng = center_lng + radius * math.cos(angle)
        coordinates.append([pt_lng, pt_lat])
    
    # Close the polygon
    coordinates.append(coordinates[0])
    
    feature = {
        "type": "Feature",
        "properties": {
            "zone_id": f"LS{100 + i}",
            "risk_score": risk_score,
            "tier": tier_name,
            "population": random.randint(100, 2500),
            "hazard_type": "landslide"
        },
        "geometry": {
            "type": "Polygon",
            "coordinates": [coordinates]
        }
    }
    
    data["features"].append(feature)

with open(filepath, 'w') as f:
    json.dump(data, f, indent=2)

print(f"Added {num_landslides} irregular landslide polygons to {filepath}.")
