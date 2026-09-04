import json

filepath = '/home/host0/Projects/GO/SNE-DRS/data/risk-zones.geojson'

with open(filepath, 'r') as f:
    data = json.load(f)

SCALE_FACTOR = 4.0

def get_centroid(points):
    x_sum = sum([p[0] for p in points])
    y_sum = sum([p[1] for p in points])
    return x_sum / len(points), y_sum / len(points)

for feature in data.get('features', []):
    geom_type = feature.get('geometry', {}).get('type')
    if geom_type == 'Polygon':
        coordinates = feature['geometry']['coordinates']
        for ring in coordinates:
            # Skip the last point in centroid calculation as it duplicates the first
            centroid_x, centroid_y = get_centroid(ring[:-1])
            for i, point in enumerate(ring):
                dx = point[0] - centroid_x
                dy = point[1] - centroid_y
                ring[i] = [centroid_x + dx * SCALE_FACTOR, centroid_y + dy * SCALE_FACTOR]
    elif geom_type == 'MultiPolygon':
        coordinates = feature['geometry']['coordinates']
        for polygon in coordinates:
            for ring in polygon:
                centroid_x, centroid_y = get_centroid(ring[:-1])
                for i, point in enumerate(ring):
                    dx = point[0] - centroid_x
                    dy = point[1] - centroid_y
                    ring[i] = [centroid_x + dx * SCALE_FACTOR, centroid_y + dy * SCALE_FACTOR]

with open(filepath, 'w') as f:
    json.dump(data, f, indent=2)

print("Scaled all shapes.")
