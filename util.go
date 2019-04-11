package env

import (
	"fmt"
	"os"
	"strconv"
)

// Set - sets an environment variable
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset - unsets an environment variable
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Get - get an environment variable, empty string if does not exist
func Get(key string) string {
	return os.Getenv(key)
}

// GetOr - get an environment variable or return default value if does not exist
func GetOr(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return defaultValue
}

// MustGet - get an environment variable or panic if does not exist
func MustGet(key string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetBool - get an environment variable as boolean
func GetBool(key string) (bool, error) {
	return strconv.ParseBool(os.Getenv(key))
}

// GetOrBool - get an environment variable or return default value if does not exist
func GetOrBool(key string, defaultValue bool) bool {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseBool(strValue)
		if err == nil {
			return value
		}
	}
	return defaultValue
}

// MustGetBool - get an environment variable or panic if does not exist
func MustGetBool(key string) bool {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseBool(strValue)
		if err == nil {
			return value
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be convert to boolean", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetInt - get an environment variable as int
func GetInt(key string) (int, error) {
	value, err := strconv.ParseInt(os.Getenv(key), 10, 32)
	return int(value), err
}

// GetOrInt - get an environment variable or return default value if does not exist
func GetOrInt(key string, defaultValue int) int {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseInt(strValue, 10, 32)
		if err == nil {
			return int(value)
		}
	}
	return defaultValue
}

// MustGetInt - get an environment variable or panic if does not exist
func MustGetInt(key string) int {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseInt(strValue, 10, 32)
		if err == nil {
			return int(value)
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be convert to int", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetUint - get an environment variable as uint
func GetUint(key string) (uint, error) {
	value, err := strconv.ParseUint(os.Getenv(key), 10, 32)
	return uint(value), err
}

// GetOrUint - get an environment variable or return default value if does not exist
func GetOrUint(key string, defaultValue uint) uint {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseUint(strValue, 10, 32)
		if err == nil {
			return uint(value)
		}
	}
	return defaultValue
}

// MustGetUint - get an environment variable or panic if does not exist
func MustGetUint(key string) uint {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseUint(strValue, 10, 32)
		if err == nil {
			return uint(value)
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be convert to uint", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetFloat - get an environment variable as float32
func GetFloat32(key string) (float32, error) {
	value, err := strconv.ParseFloat(os.Getenv(key), 32)
	return float32(value), err
}

// GetOrFloat - get an environment variable or return default value if does not exist
func GetOrFloat32(key string, defaultValue float32) float32 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseFloat(strValue, 32)
		if err == nil {
			return float32(value)
		}
	}
	return defaultValue
}

// MustGetUFloat - get an environment variable or panic if does not exist
func MustGetFloat32(key string) float32 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseFloat(strValue, 32)
		if err == nil {
			return float32(value)
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be convert to float32", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}
