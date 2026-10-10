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

func CalcDuration(tp string, valueStr string, unit string, roundoff bool, baseTime ...time.Time) (time.Time, error) {
	// Handle optional baseTime (default to time.Now() if not provided or zero)
	var t time.Time
	if len(baseTime) == 0 || baseTime[0].IsZero() {
		t = time.Now()
	} else {
		t = baseTime[0]
	}

	val, err := strconv.Atoi(valueStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid value format %q: %w", valueStr, err)
	}

	var result time.Time
	tFunc := func(d time.Duration) time.Time {
		switch tp {
		case "add":
			return t.Add(d)
		case "sub":
			return t.Add(-d)
		default:
			return t.Add(d)
		}
	}

	switch unit {
	case "s", "sec", "seconds":
		result = tFunc(time.Duration(val) * time.Second)
		if roundoff {
			interval := 5 * time.Second
			result = result.Round(interval)
		}
	case "m", "min", "minutes":
		result = tFunc(time.Duration(val) * time.Minute)
		if roundoff {
			interval := 30 * time.Second
			result = result.Round(interval)
		}
	case "h", "hour", "hours":
		result = tFunc(time.Duration(val) * time.Hour)
		if roundoff {
			result = result.Round(time.Hour)
		}
	case "d", "day", "days":
		result = tFunc(time.Duration(val) * 24 * time.Hour)
		if roundoff {
			// Round to the start of the local day (Midnight 00:00:00)
			y, m, d := result.Date()
			result = time.Date(y, m, d, 0, 0, 0, 0, result.Location())
		}
	case "w", "week", "weeks":
		result = tFunc(time.Duration(val) * 7 * 24 * time.Hour)
		if roundoff {
			// Round to the start of the local day (Midnight 00:00:00)
			y, m, d := result.Date()
			result = time.Date(y, m, d, 0, 0, 0, 0, result.Location())
		}
	case "y", "year", "years":
		result = tFunc(time.Duration(val) * 365 * 24 * time.Hour)
		if roundoff {
			y := result.Year()
			result = time.Date(y, time.January, 1, 0, 0, 0, 0, result.Location())
		}
	default:
		return time.Time{}, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return result, nil
}

func CalcDifference(startTime, endTime time.Time, unit string) (float64, error) {
	if endTime.Before(startTime) {
		return 0, fmt.Errorf("end time %v is before start time %v", endTime, startTime)
	}
	diff := endTime.Sub(startTime)

	switch unit {
	case "s", "sec", "seconds":
		return diff.Seconds(), nil
	case "m", "min", "minutes":
		return diff.Minutes(), nil
	case "h", "hour", "hours":
		return diff.Hours(), nil
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}
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
