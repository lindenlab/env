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

// Parser is a function type that converts a string value to type T.
// It is used by the generic parsing functions to provide type-safe
// conversion from environment variable string values.
type Parser[T any] func(string) (T, error)

// GetParsed retrieves an environment variable and parses it using the provided parser function.
// It returns the parsed value and any parsing error. If the environment variable is not set,
// the parser receives an empty string.
//
// This is a generic function that can be used with any type that has a corresponding parser.
// For common types, use the specific Get* functions which are more convenient.
func GetParsed[T any](key string, parser Parser[T]) (T, error) {
	value := os.Getenv(key)
	return parser(value)
}

// GetOrParsed retrieves an environment variable and parses it using the provided parser function.
// If the environment variable is not set or parsing fails, it returns the default value.
//
// This function never returns an error; it falls back to the default value on any failure.
func GetOrParsed[T any](key string, defaultValue T, parser Parser[T]) T {
	strValue, ok := os.LookupEnv(key)
	if ok {
		if value, err := parser(strValue); err == nil {
			return value
		}
	}
	return defaultValue
}

// MustGetParsed retrieves an environment variable and parses it using the provided parser function.
// If the environment variable is not set or parsing fails, it panics with a descriptive message.
//
// The typeName parameter is used in panic messages to identify the expected type.
// Use this function when the environment variable is required for the application to function.
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

// Set sets an environment variable to the specified value.
// It returns an error if the operation fails.
//
// This is a simple wrapper around os.Setenv for consistency with other functions in this package.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset removes an environment variable.
// It returns an error if the operation fails.
//
// This is a simple wrapper around os.Unsetenv for consistency with other functions in this package.
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Get retrieves the value of an environment variable.
// If the variable is not set, it returns an empty string.
//
// This is a simple wrapper around os.Getenv for consistency with other functions in this package.
func Get(key string) string {
	return os.Getenv(key)
}

// GetOr retrieves the value of an environment variable.
// If the variable is not set, it returns the provided default value.
//
// This function distinguishes between unset variables and variables set to empty strings.
func GetOr(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return defaultValue
}

// MustGet retrieves the value of an environment variable.
// If the variable is not set, it panics with a descriptive message.
//
// Use this function when the environment variable is required for the application to function.
func MustGet(key string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	panic(fmt.Sprintf("expected environment variable \"%s\" does not exist", key))
}

// MustGetWithContext retrieves the value of an environment variable with additional context.
// If the variable is not set, it panics with a descriptive message that includes the context.
//
// The context parameter should describe why this variable is needed, making errors more helpful.
// For example: "generating wallet URIs" or "connecting to database".
//
// Use this function when the environment variable is required and you want to provide
// additional context in error messages.
func MustGetWithContext(key, context string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	panic(fmt.Sprintf("required env var %s missing (needed for: %s)", key, context))
}

// GetBool retrieves an environment variable and parses it as a boolean.
// It accepts values like "true", "false", "1", "0", "t", "f", "T", "F", "TRUE", "FALSE" (case-insensitive).
// Returns the parsed boolean value and any parsing error.
func GetBool(key string) (bool, error) {
	return GetParsed(key, ParseBool)
}

// GetOrBool retrieves an environment variable and parses it as a boolean.
// If the variable is not set or parsing fails, it returns the default value.
func GetOrBool(key string, defaultValue bool) bool {
	return GetOrParsed(key, defaultValue, ParseBool)
}

// MustGetBool retrieves an environment variable and parses it as a boolean.
// If the variable is not set or parsing fails, it panics.
func MustGetBool(key string) bool {
	return MustGetParsed(key, ParseBool, "bool")
}

// GetInt retrieves an environment variable and parses it as a signed integer.
// It accepts decimal integers that fit in the int type (platform-dependent size).
// Returns the parsed integer value and any parsing error.
func GetInt(key string) (int, error) {
	return GetParsed(key, ParseInt)
}

// GetOrInt retrieves an environment variable and parses it as a signed integer.
// If the variable is not set or parsing fails, it returns the default value.
func GetOrInt(key string, defaultValue int) int {
	return GetOrParsed(key, defaultValue, ParseInt)
}

// MustGetInt retrieves an environment variable and parses it as a signed integer.
// If the variable is not set or parsing fails, it panics.
func MustGetInt(key string) int {
	return MustGetParsed(key, ParseInt, "int")
}

// GetUint retrieves an environment variable and parses it as an unsigned integer.
// It accepts non-negative decimal integers that fit in the uint type (platform-dependent size).
// Returns the parsed unsigned integer value and any parsing error.
func GetUint(key string) (uint, error) {
	return GetParsed(key, ParseUint)
}

// GetOrUint retrieves an environment variable and parses it as an unsigned integer.
// If the variable is not set or parsing fails, it returns the default value.
func GetOrUint(key string, defaultValue uint) uint {
	return GetOrParsed(key, defaultValue, ParseUint)
}

// MustGetUint retrieves an environment variable and parses it as an unsigned integer.
// If the variable is not set or parsing fails, it panics.
func MustGetUint(key string) uint {
	return MustGetParsed(key, ParseUint, "uint")
}

// GetFloat32 retrieves an environment variable and parses it as a 32-bit floating point number.
// It accepts decimal numbers in standard or scientific notation.
// Returns the parsed float32 value and any parsing error.
func GetFloat32(key string) (float32, error) {
	return GetParsed(key, ParseFloat32)
}

// GetOrFloat32 retrieves an environment variable and parses it as a 32-bit floating point number.
// If the variable is not set or parsing fails, it returns the default value.
func GetOrFloat32(key string, defaultValue float32) float32 {
	return GetOrParsed(key, defaultValue, ParseFloat32)
}

// MustGetFloat32 retrieves an environment variable and parses it as a 32-bit floating point number.
// If the variable is not set or parsing fails, it panics.
func MustGetFloat32(key string) float32 {
	return MustGetParsed(key, ParseFloat32, "float32")
}

// GetFloat64 retrieves an environment variable and parses it as a 64-bit floating point number.
// It accepts decimal numbers in standard or scientific notation.
// Returns the parsed float64 value and any parsing error.
func GetFloat64(key string) (float64, error) {
	return GetParsed(key, ParseFloat64)
}

// GetOrFloat64 retrieves an environment variable and parses it as a 64-bit floating point number.
// If the variable is not set or parsing fails, it returns the default value.
func GetOrFloat64(key string, defaultValue float64) float64 {
	return GetOrParsed(key, defaultValue, ParseFloat64)
}

// MustGetFloat64 retrieves an environment variable and parses it as a 64-bit floating point number.
// If the variable is not set or parsing fails, it panics.
func MustGetFloat64(key string) float64 {
	return MustGetParsed(key, ParseFloat64, "float64")
}

// GetInt64 retrieves an environment variable and parses it as a 64-bit signed integer.
// It accepts decimal integers in the range -9223372036854775808 to 9223372036854775807.
// Returns the parsed int64 value and any parsing error.
func GetInt64(key string) (int64, error) {
	return GetParsed(key, ParseInt64)
}

// GetOrInt64 retrieves an environment variable and parses it as a 64-bit signed integer.
// If the variable is not set or parsing fails, it returns the default value.
func GetOrInt64(key string, defaultValue int64) int64 {
	return GetOrParsed(key, defaultValue, ParseInt64)
}

// MustGetInt64 retrieves an environment variable and parses it as a 64-bit signed integer.
// If the variable is not set or parsing fails, it panics.
func MustGetInt64(key string) int64 {
	return MustGetParsed(key, ParseInt64, "int64")
}

// GetUint64 retrieves an environment variable and parses it as a 64-bit unsigned integer.
// It accepts non-negative decimal integers in the range 0 to 18446744073709551615.
// Returns the parsed uint64 value and any parsing error.
func GetUint64(key string) (uint64, error) {
	return GetParsed(key, ParseUint64)
}

// GetOrUint64 retrieves an environment variable and parses it as a 64-bit unsigned integer.
// If the variable is not set or parsing fails, it returns the default value.
func GetOrUint64(key string, defaultValue uint64) uint64 {
	return GetOrParsed(key, defaultValue, ParseUint64)
}

// MustGetUint64 retrieves an environment variable and parses it as a 64-bit unsigned integer.
// If the variable is not set or parsing fails, it panics.
func MustGetUint64(key string) uint64 {
	return MustGetParsed(key, ParseUint64, "uint64")
}

// GetDuration retrieves an environment variable and parses it as a time.Duration.
// It accepts duration strings like "5s", "2m30s", "1h", "300ms", etc.
// See time.ParseDuration for the complete format specification.
// Returns the parsed duration value and any parsing error.
func GetDuration(key string) (time.Duration, error) {
	return GetParsed(key, ParseDuration)
}

// GetOrDuration retrieves an environment variable and parses it as a time.Duration.
// If the variable is not set or parsing fails, it parses and returns the default value.
// The default value must be a valid duration string, or the function will panic.
//
// Note: This function takes a string default value (unlike other GetOr* functions)
// to maintain backwards compatibility with existing APIs.
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

// MustGetDuration retrieves an environment variable and parses it as a time.Duration.
// If the variable is not set or parsing fails, it panics.
func MustGetDuration(key string) time.Duration {
	return MustGetParsed(key, ParseDuration, "time.Duration")
}

// GetUrl retrieves an environment variable and parses it as a URL.
// It accepts absolute URLs and parses them using url.ParseRequestURI.
// Returns the parsed *url.URL value and any parsing error.
func GetUrl(key string) (*url.URL, error) {
	return GetParsed(key, ParseURL)
}

// GetOrUrl retrieves an environment variable and parses it as a URL.
// If the variable is not set or parsing fails, it parses and returns the default value.
// The default value must be a valid URL string, or the function will panic.
//
// Note: This function takes a string default value (unlike other GetOr* functions)
// to maintain backwards compatibility with existing APIs.
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

// MustGetUrl retrieves an environment variable and parses it as a URL.
// If the variable is not set or parsing fails, it panics.
func MustGetUrl(key string) *url.URL {
	return MustGetParsed(key, ParseURL, "url.URL")
}

// TestHelper is an interface that represents a testing object (typically *testing.T or *testing.B).
// It provides the minimal interface needed for test helper functions.
type TestHelper interface {
	Helper()
	Cleanup(func())
}

// SetForTest sets an environment variable for testing and automatically cleans it up.
// It returns a cleanup function that can be called manually, though cleanup is also
// registered with t.Cleanup() if t is non-nil.
//
// If the environment variable already exists, it will be restored to its original value.
// If it doesn't exist, it will be removed during cleanup.
//
// Example usage:
//
//	func TestConfig(t *testing.T) {
//	    cleanup := env.SetForTest(t, "WALLET_TRANSACTION_API", "https://test.com")
//	    defer cleanup() // Optional: cleanup is also registered with t.Cleanup()
//	    // test code
//	}
func SetForTest(t TestHelper, key, value string) func() {
	if t != nil {
		t.Helper()
	}

	old, existed := os.LookupEnv(key)
	os.Setenv(key, value)

	cleanup := func() {
		if existed {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	}

	if t != nil {
		t.Cleanup(cleanup)
	}

	return cleanup
}

// UnsetForTest removes an environment variable for testing and automatically restores it.
// It returns a cleanup function that can be called manually, though cleanup is also
// registered with t.Cleanup() if t is non-nil.
//
// If the environment variable exists, it will be restored to its original value during cleanup.
//
// Example usage:
//
//	func TestWithoutEnv(t *testing.T) {
//	    cleanup := env.UnsetForTest(t, "OPTIONAL_CONFIG")
//	    defer cleanup()
//	    // test code that expects the variable to not exist
//	}
func UnsetForTest(t TestHelper, key string) func() {
	if t != nil {
		t.Helper()
	}

	old, existed := os.LookupEnv(key)
	os.Unsetenv(key)

	cleanup := func() {
		if existed {
			os.Setenv(key, old)
		}
	}

	if t != nil {
		t.Cleanup(cleanup)
	}

	return cleanup
}
