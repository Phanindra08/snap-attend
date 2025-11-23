package utils

import (
	"math"
	"strings"
	"time"
)

// TrimAndConvertToLowerCase - Helps in trimming whitespaces and converting any string to lowercase
func TrimAndConvertToLowerCase(content string) string {
	return strings.TrimSpace(strings.ToLower(content))
}

var getStartOfToday = func() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

func IsEmpty(content string) bool {
	return len(strings.TrimSpace(content)) == 0
}

func CalculateDistanceInMeters(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const earthRadius = 6371000.0 // meters

	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	rLat1 := toRad(lat1)
	rLat2 := toRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
