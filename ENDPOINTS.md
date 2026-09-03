**GET /habitations**
Purpose: A ranked list of villages that are most aat risk

**GET /risk-zones**
Purpose: Visualization on the map. mapped to the habitations. its a polegon zone

**GET /sites**
Purpose: Safe site areas.

## mode 2
**GET /routes?habitation_id=X**
Purpose: Shortest path that avoids hazards

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

