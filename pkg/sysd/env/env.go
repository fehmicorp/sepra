package env

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

func LoadConfig(filename string, configPtr interface{}, defaults interface{}) (interface{}, error) {
	// 1. Pre-populate configPtr with default values by marshaling and unmarshaling defaults
	defaultBytes, err := json.Marshal(defaults)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(defaultBytes, configPtr); err != nil {
		return nil, err
	}
	// 2. Check if the config file exists. If it doesn't, return the default-populated struct.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return configPtr, nil
	}
	// 3. Read the configuration file from disk
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 4. Unmarshal file contents into the config struct, overriding any defaults
	if err := json.Unmarshal(fileData, configPtr); err != nil {
		return nil, err
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
