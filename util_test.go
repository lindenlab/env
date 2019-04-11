package env

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
)

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

func TestFloatFuncs(t *testing.T) {
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