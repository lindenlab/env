package env

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
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
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to boolean", key))
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
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to int", key))
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
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to uint", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetFloat32 - get an environment variable as float32
func GetFloat32(key string) (float32, error) {
	value, err := strconv.ParseFloat(os.Getenv(key), 32)
	return float32(value), err
}

// GetOrFloat32 - get an environment variable or return default value if does not exist
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

// MustGetUFloat32 - get an environment variable or panic if does not exist
func MustGetFloat32(key string) float32 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseFloat(strValue, 32)
		if err == nil {
			return float32(value)
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to float32", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetFloat64 - get an environment variable as float64
func GetFloat64(key string) (float64, error) {
	value, err := strconv.ParseFloat(os.Getenv(key), 64)
	return float64(value), err
}

// GetOrFloat64 - get an environment variable or return default value if does not exist
func GetOrFloat64(key string, defaultValue float64) float64 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseFloat(strValue, 64)
		if err == nil {
			return float64(value)
		}
	}
	return defaultValue
}

// MustGetUFloat64 - get an environment variable or panic if does not exist
func MustGetFloat64(key string) float64 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseFloat(strValue, 64)
		if err == nil {
			return float64(value)
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to float64", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetInt64 - get an environment variable as int64
func GetInt64(key string) (int64, error) {
	value, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	return int64(value), err
}

// GetOrInt64 - get an environment variable or return default value if does not exist
func GetOrInt64(key string, defaultValue int64) int64 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseInt(strValue, 10, 64)
		if err == nil {
			return int64(value)
		}
	}
	return defaultValue
}

// MustGetInt64 - get an environment variable or panic if does not exist
func MustGetInt64(key string) int64 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseInt(strValue, 10, 64)
		if err == nil {
			return int64(value)
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to int64", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetUint64 - get an environment variable as uint
func GetUint64(key string) (uint64, error) {
	value, err := strconv.ParseUint(os.Getenv(key), 10, 64)
	return uint64(value), err
}

// GetOrUint64 - get an environment variable or return default value if does not exist
func GetOrUint64(key string, defaultValue uint64) uint64 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseUint(strValue, 10, 64)
		if err == nil {
			return uint64(value)
		}
	}
	return defaultValue
}

// MustGetUint64 - get an environment variable or panic if does not exist
func MustGetUint64(key string) uint64 {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseUint(strValue, 10, 64)
		if err == nil {
			return uint64(value)
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to uint64", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetDuration - get an environment variable as time.Duration
func GetDuration(key string) (time.Duration, error) {
	value, err := time.ParseDuration(os.Getenv(key))
	return value, err
}

// GetOrDuration - get an environment variable or return default value if does not exist
func GetOrDuration(key string, defaultValue string) time.Duration {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := time.ParseDuration(strValue)
		if err == nil {
			return value
		}
	}
	defaultDuration, err := time.ParseDuration(defaultValue)
	if err != nil {
		panic(fmt.Sprintf("default duration \"%s\" could not be converted to time.Duration", key))
	}
	return defaultDuration
}

// MustGetDuration - get an environment variable or panic if does not exist
func MustGetDuration(key string) time.Duration {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := time.ParseDuration(strValue)
		if err == nil {
			return value
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to time.Duration", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// GetUrl - get an environment variable as url.URL
func GetUrl(key string) (*url.URL, error) {
	value, err := url.ParseRequestURI(os.Getenv(key))
	return value, err
}

// GetOrUrl - get an environment variable or return default value if does not exist
func GetOrUrl(key string, defaultValue string) *url.URL {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := url.ParseRequestURI(strValue)
		if err == nil {
			return value
		}
	}
	defaultUrl, err := url.ParseRequestURI(defaultValue)
	if err != nil {
		panic(fmt.Sprintf("default duration \"%s\" could not be converted to url.URL", key))
	}
	return defaultUrl
}

// MustGetUrl - get an environment variable or panic if does not exist
func MustGetUrl(key string) *url.URL {
	strValue, ok := os.LookupEnv(key)
	if ok {
		value, err := url.ParseRequestURI(strValue)
		if err == nil {
			return value
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to url.URL", key))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}
