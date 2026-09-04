import json
import math

def create_circle(cx, cy, radius_deg, num_points=32):
    points = []
    for i in range(num_points):
        angle = (2 * math.pi * i) / num_points
        dx = radius_deg * math.cos(angle)
        dy = radius_deg * math.sin(angle)
        points.append([cx + dx, cy + dy])
    points.append(points[0]) # Close the loop
    return [points]

with open('data/risk-zones.geojson', 'r') as f:
    data = json.load(f)

for feature in data['features']:
    if feature['properties'].get('hazard_type') == 'tornado':
        # Get center from first point of existing geometry
        cx = feature['geometry']['coordinates'][0][0][0]
        cy = feature['geometry']['coordinates'][0][0][1]
        
        # Huge circle (0.3 degrees ~ 33km radius)
        feature['geometry']['coordinates'] = create_circle(cx, cy, 0.4)

with open('data/risk-zones.geojson', 'w') as f:
    json.dump(data, f, indent=2)
