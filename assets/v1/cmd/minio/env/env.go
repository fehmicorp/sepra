package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

func LoadConfig(filename string, configPtr interface{}, defaults interface{}) (interface{}, error) {
	val := reflect.ValueOf(configPtr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil, fmt.Errorf("configPtr must be a non-nil pointer to a struct")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("configPtr must point to a struct")
	}

	// 1. Populate with default values object if provided
	if defaults != nil {
		defaultBytes, err := json.Marshal(defaults)
		if err == nil {
			_ = json.Unmarshal(defaultBytes, configPtr)
		}
	}

	// Apply field-level `default` tags
	applyDefaultsFromTags(elem)

	// 2. Read from configuration file if it exists
	if filename != "" {
		if _, err := os.Stat(filename); err == nil {
			fileData, err := os.ReadFile(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}

			ext := strings.ToLower(filepath.Ext(filename))
			if ext == ".yaml" || ext == ".yml" {
				if err := yaml.Unmarshal(fileData, configPtr); err != nil {
					return nil, fmt.Errorf("failed to parse YAML config: %w", err)
				}
			} else if ext == ".json" {
				if err := json.Unmarshal(fileData, configPtr); err != nil {
					return nil, fmt.Errorf("failed to parse JSON config: %w", err)
				}
			}
		}
	}

	// 3. Override with Environment Variables (`env` tags)
	applyEnvOverrides(elem)

	return configPtr, nil
}

func applyDefaultsFromTags(elem reflect.Value) {
	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		fieldVal := elem.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		structField := typ.Field(i)
		defaultTag := structField.Tag.Get("default")
		if defaultTag != "" {
			setFieldValue(fieldVal, defaultTag)
		}
	}
}

func applyEnvOverrides(elem reflect.Value) {
	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		fieldVal := elem.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		structField := typ.Field(i)
		envTag := structField.Tag.Get("env")
		if envTag == "" {
			continue
		}

		if envVal, exists := os.LookupEnv(envTag); exists && envVal != "" {
			setFieldValue(fieldVal, envVal)
		}
	}
}

func setFieldValue(fieldVal reflect.Value, valStr string) {
	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(valStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(valStr, 10, 64); err == nil {
			fieldVal.SetInt(intVal)
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(valStr); err == nil {
			fieldVal.SetBool(boolVal)
		}
	case reflect.Float32, reflect.Float64:
		if floatVal, err := strconv.ParseFloat(valStr, 64); err == nil {
			fieldVal.SetFloat(floatVal)
		}
	}
}

func GetString(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func GetBool(key string, defaultValue bool) bool {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		if durationVal, err := time.ParseDuration(val); err == nil {
			return durationVal
		}
	}
	return defaultValue
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func UnsetEnv(key string) error {
	return os.Unsetenv(key)
}
