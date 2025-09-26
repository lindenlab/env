package env

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Regular expression for validating environment variable names
// Typically follows pattern: [A-Z_][A-Z0-9_]*
var envVarNameRegex = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// isValidEnvVarKey validates that an environment variable key follows standard naming conventions
// This is a helper function for internal validation - not exported to avoid breaking changes
func isValidEnvVarKey(key string) bool {
	if key == "" {
		return false
	}
	return envVarNameRegex.MatchString(key)
}

// Generic types and functions for reducing code duplication

// Parser function type for converting strings to any type
type Parser[T any] func(string) (T, error)

// GetParsed - generic function for getting parsed environment variables
func GetParsed[T any](key string, parser Parser[T]) (T, error) {
	value := os.Getenv(key)
	return parser(value)
}

// GetOrParsed - generic function with default value
func GetOrParsed[T any](key string, defaultValue T, parser Parser[T]) T {
	strValue, ok := os.LookupEnv(key)
	if ok {
		if value, err := parser(strValue); err == nil {
			return value
		}
	}
	return defaultValue
}

// MustGetParsed - generic function that panics on missing/invalid values
func MustGetParsed[T any](key string, parser Parser[T], typeName string) T {
	strValue, ok := os.LookupEnv(key)
	if ok {
		if value, err := parser(strValue); err == nil {
			return value
		} else {
			panic(fmt.Sprintf("environment variable \"%s\" could not be converted to %s", key, typeName))
		}
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// Type-specific parser functions
var (
	ParseBool = func(s string) (bool, error) {
		return strconv.ParseBool(s)
	}
	ParseInt = func(s string) (int, error) {
		v, err := strconv.ParseInt(s, DecimalBase, Int32Bits)
		return int(v), err
	}
	ParseUint = func(s string) (uint, error) {
		v, err := strconv.ParseUint(s, DecimalBase, Int32Bits)
		return uint(v), err
	}
	ParseFloat32 = func(s string) (float32, error) {
		v, err := strconv.ParseFloat(s, Float32Bits)
		return float32(v), err
	}
	ParseFloat64 = func(s string) (float64, error) {
		return strconv.ParseFloat(s, Float64Bits)
	}
	ParseInt64 = func(s string) (int64, error) {
		return strconv.ParseInt(s, DecimalBase, Int64Bits)
	}
	ParseUint64 = func(s string) (uint64, error) {
		return strconv.ParseUint(s, DecimalBase, Int64Bits)
	}
	ParseDuration = func(s string) (time.Duration, error) {
		return time.ParseDuration(s)
	}
	ParseURL = func(s string) (*url.URL, error) {
		return url.ParseRequestURI(s)
	}
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

// Bool functions
func GetBool(key string) (bool, error) {
	return GetParsed(key, ParseBool)
}

func GetOrBool(key string, defaultValue bool) bool {
	return GetOrParsed(key, defaultValue, ParseBool)
}

func MustGetBool(key string) bool {
	return MustGetParsed(key, ParseBool, "bool")
}

// Int functions
func GetInt(key string) (int, error) {
	return GetParsed(key, ParseInt)
}

func GetOrInt(key string, defaultValue int) int {
	return GetOrParsed(key, defaultValue, ParseInt)
}

func MustGetInt(key string) int {
	return MustGetParsed(key, ParseInt, "int")
}

// Uint functions
func GetUint(key string) (uint, error) {
	return GetParsed(key, ParseUint)
}

func GetOrUint(key string, defaultValue uint) uint {
	return GetOrParsed(key, defaultValue, ParseUint)
}

func MustGetUint(key string) uint {
	return MustGetParsed(key, ParseUint, "uint")
}

// Float32 functions
func GetFloat32(key string) (float32, error) {
	return GetParsed(key, ParseFloat32)
}

func GetOrFloat32(key string, defaultValue float32) float32 {
	return GetOrParsed(key, defaultValue, ParseFloat32)
}

func MustGetFloat32(key string) float32 {
	return MustGetParsed(key, ParseFloat32, "float32")
}

// Float64 functions
func GetFloat64(key string) (float64, error) {
	return GetParsed(key, ParseFloat64)
}

func GetOrFloat64(key string, defaultValue float64) float64 {
	return GetOrParsed(key, defaultValue, ParseFloat64)
}

func MustGetFloat64(key string) float64 {
	return MustGetParsed(key, ParseFloat64, "float64")
}

// Int64 functions
func GetInt64(key string) (int64, error) {
	return GetParsed(key, ParseInt64)
}

func GetOrInt64(key string, defaultValue int64) int64 {
	return GetOrParsed(key, defaultValue, ParseInt64)
}

func MustGetInt64(key string) int64 {
	return MustGetParsed(key, ParseInt64, "int64")
}

// Uint64 functions
func GetUint64(key string) (uint64, error) {
	return GetParsed(key, ParseUint64)
}

func GetOrUint64(key string, defaultValue uint64) uint64 {
	return GetOrParsed(key, defaultValue, ParseUint64)
}

func MustGetUint64(key string) uint64 {
	return MustGetParsed(key, ParseUint64, "uint64")
}

// Duration functions
func GetDuration(key string) (time.Duration, error) {
	return GetParsed(key, ParseDuration)
}

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
		panic(fmt.Sprintf("default duration \"%s\" could not be converted to time.Duration", defaultValue))
	}
	return defaultDuration
}

func MustGetDuration(key string) time.Duration {
	return MustGetParsed(key, ParseDuration, "time.Duration")
}

// URL functions
func GetUrl(key string) (*url.URL, error) {
	return GetParsed(key, ParseURL)
}

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
		panic(fmt.Sprintf("default url \"%s\" could not be converted to url.URL", defaultValue))
	}
	return defaultUrl
}

func MustGetUrl(key string) *url.URL {
	return MustGetParsed(key, ParseURL, "url.URL")
}