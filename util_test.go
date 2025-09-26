package env

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringFuncs(t *testing.T) {
	os.Setenv("MY_ENV", "hello")

	// env exists
	assert.Equal(t, "hello", Get("MY_ENV"))
	assert.Equal(t, "hello", MustGet("MY_ENV"))
	assert.Equal(t, "hello", GetOr("MY_ENV", "what up?"))

	// env not exists
	assert.Equal(t, "", Get("ENV_NO_EXISTS"))
	assert.Equal(t, "hello world", GetOr("ENV_NO_EXISTS", "hello world"))
	assert.Panics(t, func() { MustGet("ENV_NO_EXISTS") }, "The code did not panic")
}

func TestBoolFuncs(t *testing.T) {
	os.Setenv("BOOL_ENV", "1")

	// env exists
	val, err := GetBool("BOOL_ENV")
	assert.Equal(t, true, val)
	assert.Nil(t, err)
	assert.Equal(t, true, MustGetBool("BOOL_ENV"))
	assert.Equal(t, true, GetOrBool("BOOL_ENV", false))

	// env not exists
	val, err = GetBool("ENV_NO_EXISTS")
	assert.Equal(t, false, val)
	assert.Error(t, err)
	assert.Equal(t, true, GetOrBool("ENV_NO_EXISTS", true))
	assert.Panics(t, func() { MustGetBool("ENV_NO_EXISTS") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_BOOL", "bad_bool")

	val, err = GetBool("BAD_BOOL")
	assert.Equal(t, false, val)
	assert.Error(t, err)
	assert.Equal(t, true, GetOrBool("BAD_BOOL", true))
	assert.Panics(t, func() { MustGetBool("BAD_BOOL") }, "The code did not panic")
}

func TestIntFuncs(t *testing.T) {
	os.Setenv("INT_ENV", "-34")

	// env exists
	val, err := GetInt("INT_ENV")
	assert.Equal(t, -34, val)
	assert.Nil(t, err)
	assert.Equal(t, -34, MustGetInt("INT_ENV"))
	assert.Equal(t, -34, GetOrInt("INT_ENV", 99))

	// env not exists
	val, err = GetInt("ENV_NO_EXISTS")
	assert.Equal(t, 0, val)
	assert.Error(t, err)
	assert.Equal(t, 99, GetOrInt("ENV_NO_EXISTS", 99))
	assert.Panics(t, func() { MustGetInt("ENV_NO_EXISTS") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_INT", "bad_int")

	val, err = GetInt("BAD_INT")
	assert.Equal(t, 0, val)
	assert.Error(t, err)
	assert.Equal(t, 99, GetOrInt("BAD_INT", 99))
	assert.Panics(t, func() { MustGetInt("BAD_INT") }, "The code did not panic")
}

func TestUintFuncs(t *testing.T) {
	os.Setenv("UINT_ENV", "106")

	// env exists
	val, err := GetUint("UINT_ENV")
	assert.Equal(t, uint(106), val)
	assert.Nil(t, err)
	assert.Equal(t, uint(106), MustGetUint("UINT_ENV"))
	assert.Equal(t, uint(106), GetOrUint("UINT_ENV", 99))

	// env not exists
	val, err = GetUint("ENV_NO_EXISTS")
	assert.Equal(t, uint(0), val)
	assert.Error(t, err)
	assert.Equal(t, uint(99), GetOrUint("ENV_NO_EXISTS", 99))
	assert.Panics(t, func() { MustGetUint("ENV_NO_EXISTS") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_UINT", "-10")

	val, err = GetUint("BAD_UINT")
	assert.Equal(t, uint(0), val)
	assert.Error(t, err)
	assert.Equal(t, uint(99), GetOrUint("BAD_UINT", 99))
	assert.Panics(t, func() { MustGetUint("BAD_UINT") }, "The code did not panic")
}

func TestFloat32Funcs(t *testing.T) {
	os.Setenv("FLOAT_ENV", "12.34")

	// env exists
	val, err := GetFloat32("FLOAT_ENV")
	assert.Equal(t, float32(12.34), val)
	assert.Nil(t, err)
	assert.Equal(t, float32(12.34), MustGetFloat32("FLOAT_ENV"))
	assert.Equal(t, float32(12.34), GetOrFloat32("FLOAT_ENV", 66.6))

	// env not exists
	val, err = GetFloat32("ENV_NO_EXISTS")
	assert.Equal(t, float32(0), val)
	assert.Error(t, err)
	assert.Equal(t, float32(66.6), GetOrFloat32("ENV_NO_EXISTS", 66.6))
	assert.Panics(t, func() { MustGetFloat32("ENV_NO_EXISTS") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_FLOAT", "bad_float")

	val, err = GetFloat32("BAD_FLOAT")
	assert.Equal(t, float32(0), val)
	assert.Error(t, err)
	assert.Equal(t, float32(1.23), GetOrFloat32("BAD_FLOAT", 1.23))
	assert.Panics(t, func() { MustGetFloat32("BAD_FLOAT") }, "The code did not panic")
}

func TestFloat64Funcs(t *testing.T) {
	os.Setenv("FLOAT_ENV", "12.34")

	// env exists
	val, err := GetFloat64("FLOAT_ENV")
	assert.Equal(t, float64(12.34), val)
	assert.Nil(t, err)
	assert.Equal(t, float64(12.34), MustGetFloat64("FLOAT_ENV"))
	assert.Equal(t, float64(12.34), GetOrFloat64("FLOAT_ENV", 66.6))

	// env not exists
	val, err = GetFloat64("ENV_NO_EXISTS")
	assert.Equal(t, float64(0), val)
	assert.Error(t, err)
	assert.Equal(t, float64(66.6), GetOrFloat64("ENV_NO_EXISTS", 66.6))
	assert.Panics(t, func() { MustGetFloat64("ENV_NO_EXISTS") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_FLOAT", "bad_float")

	val, err = GetFloat64("BAD_FLOAT")
	assert.Equal(t, float64(0), val)
	assert.Error(t, err)
	assert.Equal(t, float64(1.23), GetOrFloat64("BAD_FLOAT", 1.23))
	assert.Panics(t, func() { MustGetFloat64("BAD_FLOAT") }, "The code did not panic")
}

func TestInt64Funcs(t *testing.T) {
	os.Setenv("INT_ENV", "-34")

	// env exists
	val, err := GetInt64("INT_ENV")
	assert.Equal(t, int64(-34), val)
	assert.Nil(t, err)
	assert.Equal(t, int64(-34), MustGetInt64("INT_ENV"))
	assert.Equal(t, int64(-34), GetOrInt64("INT_ENV", 99))

	// env not exists
	val, err = GetInt64("ENV_NO_EXISTS")
	assert.Equal(t, int64(0), val)
	assert.Error(t, err)
	assert.Equal(t, int64(99), GetOrInt64("ENV_NO_EXISTS", 99))
	assert.Panics(t, func() { MustGetInt64("ENV_NO_EXISTS") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_INT", "bad_int")

	val, err = GetInt64("BAD_INT")
	assert.Equal(t, int64(0), val)
	assert.Error(t, err)
	assert.Equal(t, int64(99), GetOrInt64("BAD_INT", 99))
	assert.Panics(t, func() { MustGetInt64("BAD_INT") }, "The code did not panic")
}

func TestUint64Funcs(t *testing.T) {
	os.Setenv("UINT_ENV", "106")

	// env exists
	val, err := GetUint64("UINT_ENV")
	assert.Equal(t, uint64(106), val)
	assert.Nil(t, err)
	assert.Equal(t, uint64(106), MustGetUint64("UINT_ENV"))
	assert.Equal(t, uint64(106), GetOrUint64("UINT_ENV", 99))

	// env not exists
	val, err = GetUint64("ENV_NO_EXISTS")
	assert.Equal(t, uint64(0), val)
	assert.Error(t, err)
	assert.Equal(t, uint64(99), GetOrUint64("ENV_NO_EXISTS", 99))
	assert.Panics(t, func() { MustGetUint64("ENV_NO_EXISTS") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_UINT", "-10")

	val, err = GetUint64("BAD_UINT")
	assert.Equal(t, uint64(0), val)
	assert.Error(t, err)
	assert.Equal(t, uint64(99), GetOrUint64("BAD_UINT", 99))
	assert.Panics(t, func() { MustGetUint64("BAD_UINT") }, "The code did not panic")
}

func TestDurationFuncs(t *testing.T) {
	os.Setenv("DURATION_ENV", "5s")

	// env exists
	val, err := GetDuration("DURATION_ENV")
	assert.Equal(t, time.Duration(5*time.Second), val)
	assert.Nil(t, err)
	assert.Equal(t, time.Duration(5*time.Second), MustGetDuration("DURATION_ENV"))
	assert.Equal(t, time.Duration(5*time.Second), GetOrDuration("DURATION_ENV", "10m"))

	// env not exists
	val, err = GetDuration("ENV_NO_EXISTS")
	assert.Equal(t, time.Duration(0), val)
	assert.Error(t, err)
	assert.Equal(t, time.Duration(10*time.Minute), GetOrDuration("ENV_NO_EXISTS", "10m"))
	assert.Panics(t, func() { MustGetDuration("ENV_NO_EXISTS") }, "The code did not panic")

	assert.Panics(t, func() { GetOrDuration("ENV_NO_EXISTS", "bad_duration") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_DURATION", "cat_nip")

	val, err = GetDuration("BAD_DURATION")
	assert.Equal(t, time.Duration(0), val)
	assert.Error(t, err)
	assert.Equal(t, time.Duration(3*time.Hour), GetOrDuration("BAD_DURATION", "3h"))
	assert.Panics(t, func() { MustGetDuration("BAD_DURATION") }, "The code did not panic")
}

func TestUrlFuncs(t *testing.T) {
	os.Setenv("URL_ENV", "https://lindenlab.com/foo")

	// env exists
	val, err := GetUrl("URL_ENV")
	assert.Equal(t, "lindenlab.com", val.Hostname())
	assert.Nil(t, err)
	assert.Equal(t, "lindenlab.com", MustGetUrl("URL_ENV").Hostname())
	assert.Equal(t, "lindenlab.com", GetOrUrl("URL_ENV", "http://google.com").Hostname())

	// env not exists
	val, err = GetUrl("ENV_NO_EXISTS")
	assert.Equal(t, (*url.URL)(nil), val)
	assert.Error(t, err)
	assert.Equal(t, "google.com", GetOrUrl("ENV_NO_EXISTS", "http://google.com").Hostname())
	assert.Panics(t, func() { MustGetUrl("ENV_NO_EXISTS") }, "The code did not panic")

	assert.Panics(t, func() { GetOrUrl("ENV_NO_EXISTS", "bad url") }, "The code did not panic")

	// env bad format
	os.Setenv("BAD_URL", "@@@\\foo oo:slk")

	val, err = GetUrl("BAD_URL")
	assert.Equal(t, (*url.URL)(nil), val)
	assert.Error(t, err)
	assert.Equal(t, "google.com", GetOrUrl("BAD_URL", "http://google.com").Hostname())
	assert.Panics(t, func() { MustGetUrl("BAD_URL") }, "The code did not panic")
}

// Test environment variable key validation
func TestEnvVarKeyValidation(t *testing.T) {
	// Valid keys
	assert.True(t, isValidEnvVarKey("TEST_VAR"))
	assert.True(t, isValidEnvVarKey("_PRIVATE_VAR"))
	assert.True(t, isValidEnvVarKey("VAR123"))
	assert.True(t, isValidEnvVarKey("MY_APP_CONFIG"))
	assert.True(t, isValidEnvVarKey("A"))
	assert.True(t, isValidEnvVarKey("_"))

	// Invalid keys
	assert.False(t, isValidEnvVarKey(""))                    // empty
	assert.False(t, isValidEnvVarKey("123VAR"))              // starts with number
	assert.False(t, isValidEnvVarKey("my-var"))              // contains hyphen
	assert.False(t, isValidEnvVarKey("my.var"))              // contains dot
	assert.False(t, isValidEnvVarKey("my var"))              // contains space
	assert.False(t, isValidEnvVarKey("myvar"))               // lowercase (by convention)
	assert.False(t, isValidEnvVarKey("My_Var"))              // mixed case
	assert.False(t, isValidEnvVarKey("VAR!"))                // special character
}

// Test constants are properly defined
func TestParsingConstants(t *testing.T) {
	assert.Equal(t, 10, DecimalBase)
	assert.Equal(t, 32, Int32Bits)
	assert.Equal(t, 64, Int64Bits)
	assert.Equal(t, 32, Float32Bits)
	assert.Equal(t, 64, Float64Bits)
}
