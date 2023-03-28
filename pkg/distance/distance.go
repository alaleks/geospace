// Package distance intended for calculating distance
// between two cities (points) by coordinates.
package distance

import "math"

const (
	earthRaidus = 6371 // radius of the earth in kilometers.
)

// CalcGreatCirlcle perfoms calculatin the distance
// between two points by coordinates using the formula:
// https://en.wikipedia.org/wiki/Great-circle_distance
func CalcGreatCirlcle(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitudes and longitude degrees to radians.
	lat1 = degreesToRadians(lat1)
	lon1 = degreesToRadians(lon1)
	lat2 = degreesToRadians(lat2)
	lon2 = degreesToRadians(lon2)

	// Calculate cosines and sines of latitudes and longitude differences.
	diffLon := lon2 - lon1
	diffCos := math.Cos(diffLon)
	diffSin := math.Sin(diffLon)
	cosLat1 := math.Cos(lat1)
	cosLat2 := math.Cos(lat2)
	sinLat1 := math.Sin(lat1)
	sinLat2 := math.Sin(lat2)

	// Calculate the length of the great circle.
	y := math.Sqrt(math.Pow(cosLat2*diffSin, 2) +
		math.Pow(cosLat1*sinLat2-sinLat1*cosLat2*diffCos, 2))
	x := sinLat1*sinLat2 + cosLat1*cosLat2*diffCos

	// Calculate distance.
	a := math.Atan2(y, x)
	dist := a * earthRaidus

	return dist

}

// CalcHaversine perfoms calculatin the distance
// between two points by coordinates using the formula:
// https://en.wikipedia.org/wiki/Haversine_formula
func CalcHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitudes and longitude degrees to radians.
	lat1 = degreesToRadians(lat1)
	lon1 = degreesToRadians(lon1)
	lat2 = degreesToRadians(lat2)
	lon2 = degreesToRadians(lon2)

	// Calculate latitudes and longitude differences.
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	// Calculate distance.
	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dist := c * earthRaidus

	return dist
}

// degreesToRadians perfoms convert degrees to radians.
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}
