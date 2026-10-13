package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ToJSONString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(bytes), nil
}

func ToPrettyJSONString(v interface{}) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal to pretty JSON: %w", err)
	}
	return string(bytes), nil
}

func FromJSONString(jsonStr string, target interface{}) error {
	if err := json.Unmarshal([]byte(jsonStr), target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON string: %w", err)
	}
	return nil
}

func StringToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

func StringToInt64(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func StringToBool(s string) bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return val
}

func StringToFloat64(s string, defaultVal float64) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

func IntToString(n int) string {
	return strconv.Itoa(n)
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func FloatToString(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}
