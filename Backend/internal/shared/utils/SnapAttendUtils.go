package utils

import (
	"math"
	"strings"
	"time"

	"github.com/phanindra08/snap-attend/internal/shared/models"
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

func IsValidLatAndLon(lat float64, lon float64) bool {
	if lat < -90 || lat > 90 {
		return false
	} else if lon < -180 || lon > 180 {
		return false
	} else if lat == 0 && lon == 0 {
		return false
	}
	return true
}

func IsClassInSessionNow(now time.Time, sectionSchedules []models.SectionSchedule) bool {
	if len(sectionSchedules) == 0 {
		return false
	}

	weekday := now.Weekday().String()
	timeInStringFormat := now.Format("15:04:05")
	nowTime, err := time.Parse("15:04:05", timeInStringFormat)
	if err != nil {
		return false
	}

	for _, schedule := range sectionSchedules {
		if string(schedule.DaysOfTheClass) != weekday {
			continue
		}

		startStr := schedule.StartTime.String()
		endStr := schedule.EndTime.String()

		start, err := time.Parse("15:04:05", startStr)
		if err != nil {
			continue
		}
		end, err := time.Parse("15:04:05", endStr)
		if err != nil {
			continue
		}

		if !nowTime.Before(start) && !nowTime.After(end) {
			return true
		}
	}
	return false
}
