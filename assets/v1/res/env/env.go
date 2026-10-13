package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/caarlos0/env/v10"
	"go.yaml.in/yaml/v3"
)

func LoadConfig(filename string, configPtr interface{}, defaults interface{}) (interface{}, error) {
	// 1. Pre-populate configPtr with default values
	defaultBytes, err := json.Marshal(defaults)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(defaultBytes, configPtr); err != nil {
		return nil, err
	}

	// 2. Read the configuration file from disk if it exists
	if _, err := os.Stat(filename); err == nil {
		fileData, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		ext := filepath.Ext(filename)
		if ext == ".yaml" || ext == ".yml" {
			if err := yaml.Unmarshal(fileData, configPtr); err != nil {
				return nil, err
			}
		} else {
			if err := json.Unmarshal(fileData, configPtr); err != nil {
				return nil, err
			}
		}
	}

	// 3. Parse environment variables LAST so they override both defaults and file values
	if err := env.Parse(configPtr); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return configPtr, nil
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
