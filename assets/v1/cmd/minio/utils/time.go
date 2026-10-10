package utils

import (
	"fmt"
	"strconv"
	"time"
)

func GetCurrentTime() time.Time {
	return time.Now()
}

func GetCurrentTimeFormatted(layout string) string {
	return time.Now().Format(layout)
}

func GetCurrentTimeInLocation(locationName string) (time.Time, error) {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location %q: %w", locationName, err)
	}
	return time.Now().In(loc), nil
}

func GetCurrentTimeInLocationFormatted(locationName, layout string) (string, error) {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return "", fmt.Errorf("failed to load location %q: %w", locationName, err)
	}
	return time.Now().In(loc).Format(layout), nil
}

func CalcDuration(baseTime time.Time, valueStr string, unit string, roundoff bool) (time.Time, error) {
	// If baseTime is not provided (zero value), default to current time
	if baseTime.IsZero() {
		baseTime = time.Now()
	}

	val, err := strconv.Atoi(valueStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid value format %q: %w", valueStr, err)
	}

	var result time.Time

	switch unit {
	case "s", "sec", "seconds":
		result = baseTime.Add(time.Duration(val) * time.Second)
		if roundoff {
			result = result.Round(time.Second)
		}
	case "m", "min", "minutes":
		result = baseTime.Add(time.Duration(val) * time.Minute)
		if roundoff {
			result = result.Round(time.Minute)
		}
	case "h", "hour", "hours":
		result = baseTime.Add(time.Duration(val) * time.Hour)
		if roundoff {
			result = result.Round(time.Hour)
		}
	case "d", "day", "days":
		result = baseTime.AddDate(0, 0, val)
		if roundoff {
			result = result.Round(24 * time.Hour)
		}
	case "w", "week", "weeks":
		result = baseTime.AddDate(0, 0, val*7)
		if roundoff {
			result = result.Round(24 * time.Hour)
		}
	case "mo", "month", "months":
		result = baseTime.AddDate(0, val, 0)
		if roundoff {
			// Round to the start of the month
			y, m, _ := result.Date()
			result = time.Date(y, m, 1, 0, 0, 0, 0, result.Location())
		}
	case "y", "year", "years":
		result = baseTime.AddDate(val, 0, 0)
		if roundoff {
			// Round to the start of the year
			y := result.Year()
			result = time.Date(y, time.January, 1, 0, 0, 0, 0, result.Location())
		}
	default:
		return time.Time{}, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return result, nil
}

func ConvertDuration(valueStr string, targetUnit string) (float64, error) {
	// Parse the string into a standard time.Duration
	d, err := time.ParseDuration(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format %q: %w", valueStr, err)
	}

	// Convert based on the target unit
	switch targetUnit {
	case "ns", "nanoseconds":
		return float64(d.Nanoseconds()), nil
	case "us", "µs", "microseconds":
		return float64(d.Microseconds()), nil
	case "ms", "milliseconds":
		return float64(d.Milliseconds()), nil
	case "s", "seconds":
		return d.Seconds(), nil
	case "m", "minutes":
		return d.Minutes(), nil
	case "h", "hours":
		return d.Hours(), nil
	case "d", "days":
		return d.Hours() / 24, nil
	case "w", "weeks":
		return d.Hours() / (24 * 7), nil
	case "mo", "months":
		return d.Hours() / (24 * 30), nil
	case "y", "years":
		return d.Hours() / (24 * 365), nil
	case "decade", "decades":
		return d.Hours() / (24 * 365 * 10), nil
	case "century", "centuries":
		return d.Hours() / (24 * 365 * 100), nil
	default:
		return 0, fmt.Errorf("unsupported target unit: %s", targetUnit)
	}
}
