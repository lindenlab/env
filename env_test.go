package env

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type unmarshaler struct {
	time.Duration
}

// TextUnmarshaler implements encoding.TextUnmarshaler
func (d *unmarshaler) UnmarshalText(data []byte) (err error) {
	if len(data) != 0 {
		d.Duration, err = time.ParseDuration(string(data))
	} else {
		d.Duration = 0
	}
	return err
}

type Config struct {
	Some            string `env:"somevar"`
	Other           bool   `env:"othervar"`
	Port            int    `env:"PORT"`
	Int64Val        int64  `env:"INT64VAL"`
	UintVal         uint   `env:"UINTVAL"`
	Uint64Val       uint64 `env:"UINT64VAL"`
	NotAnEnv        string
	DatabaseURL     string          `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`
	Strings         []string        `env:"STRINGS"`
	SepStrings      []string        `env:"SEPSTRINGS" envSeparator:":"`
	Numbers         []int           `env:"NUMBERS"`
	Numbers64       []int64         `env:"NUMBERS64"`
	UNumbers64      []uint64        `env:"UNUMBERS64"`
	Bools           []bool          `env:"BOOLS"`
	Duration        time.Duration   `env:"DURATION"`
	Float32         float32         `env:"FLOAT32"`
	Float64         float64         `env:"FLOAT64"`
	Float32s        []float32       `env:"FLOAT32S"`
	Float64s        []float64       `env:"FLOAT64S"`
	Durations       []time.Duration `env:"DURATIONS"`
	Unmarshaler     unmarshaler     `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler    `env:"UNMARSHALER_PTR"`
	Unmarshalers    []unmarshaler   `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler  `env:"UNMARSHALER_PTRS"`
	URL             url.URL         `env:"URL"`
	URLs            []url.URL       `env:"URLS"`
}

type ParentStruct struct {
	InnerStruct *InnerStruct
	unexported  *InnerStruct
	Ignored     *http.Client
}

type InnerStruct struct {
	Inner  string `env:"innervar"`
	Number uint   `env:"innernum"`
}

type DerivedStruct struct {
	BaseStruct
}

type BaseStruct struct {
	Inner  string `env:"innervar"`
	Number uint   `env:"innernum"`
}

func TestParsesEnv(t *testing.T) {
	os.Setenv("somevar", "somevalue")
	os.Setenv("othervar", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("STRINGS", "string1,string2,string3")
	os.Setenv("SEPSTRINGS", "string1:string2:string3")
	os.Setenv("NUMBERS", "1,2,3,4")
	os.Setenv("NUMBERS64", "1,2,2147483640,-2147483640")
	os.Setenv("UNUMBERS64", "1,2,214748364011,9147483641")
	os.Setenv("BOOLS", "t,TRUE,0,1")
	os.Setenv("DURATION", "1s")
	os.Setenv("FLOAT32", "3.40282346638528859811704183484516925440e+38")
	os.Setenv("FLOAT64", "1.797693134862315708145274237317043567981e+308")
	os.Setenv("FLOAT32S", "1.0,2.0,3.0")
	os.Setenv("FLOAT64S", "1.0,2.0,3.0")
	os.Setenv("UINTVAL", "44")
	os.Setenv("UINT64VAL", "6464")
	os.Setenv("INT64VAL", "-7575")
	os.Setenv("DURATIONS", "1s,2s,3s")
	os.Setenv("UNMARSHALER", "1s")
	os.Setenv("UNMARSHALER_PTR", "1m")
	os.Setenv("UNMARSHALERS", "2m,3m")
	os.Setenv("UNMARSHALER_PTRS", "2m,3m")
	os.Setenv("URL", "http://google.com")
	os.Setenv("URLS", "ftp://foo.com:23,https://cat.com")

	defer os.Clearenv()

	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "somevalue", cfg.Some)
	assert.Equal(t, true, cfg.Other)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, uint(44), cfg.UintVal)
	assert.Equal(t, int64(-7575), cfg.Int64Val)
	assert.Equal(t, uint64(6464), cfg.Uint64Val)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.Strings)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.SepStrings)
	assert.Equal(t, []int{1, 2, 3, 4}, cfg.Numbers)
	assert.Equal(t, []int64{1, 2, 2147483640, -2147483640}, cfg.Numbers64)
	assert.Equal(t, []uint64{1, 2, 214748364011, 9147483641}, cfg.UNumbers64)
	assert.Equal(t, []bool{true, true, false, true}, cfg.Bools)
	d1, _ := time.ParseDuration("1s")
	assert.Equal(t, d1, cfg.Duration)
	f32 := float32(3.40282346638528859811704183484516925440e+38)
	assert.Equal(t, f32, cfg.Float32)
	f64 := float64(1.797693134862315708145274237317043567981e+308)
	assert.Equal(t, f64, cfg.Float64)
	assert.Equal(t, []float32{float32(1.0), float32(2.0), float32(3.0)}, cfg.Float32s)
	assert.Equal(t, []float64{float64(1.0), float64(2.0), float64(3.0)}, cfg.Float64s)
	d2, _ := time.ParseDuration("2s")
	d3, _ := time.ParseDuration("3s")
	assert.Equal(t, []time.Duration{d1, d2, d3}, cfg.Durations)
	assert.Equal(t, time.Second, cfg.Unmarshaler.Duration)
	assert.Equal(t, time.Minute, cfg.UnmarshalerPtr.Duration)
	assert.Equal(t, []unmarshaler{{time.Minute * 2}, {time.Minute * 3}}, cfg.Unmarshalers)
	assert.Equal(t, []*unmarshaler{{time.Minute * 2}, {time.Minute * 3}}, cfg.UnmarshalerPtrs)
	assert.Equal(t, "google.com", cfg.URL.Host)
	assert.Equal(t, "ftp", cfg.URLs[0].Scheme)
	assert.Equal(t, "23", cfg.URLs[0].Port())
	assert.Equal(t, "cat.com", cfg.URLs[1].Host)
}

func TestParseWithPrefix(t *testing.T) {
	os.Setenv("FOO_PORT", "1234")
	os.Setenv("FOO_STRINGS", "string1,string2,string3")

	cfg := Config{}
	assert.NoError(t, ParseWithPrefix(&cfg, "FOO_"))
	assert.Equal(t, 1234, cfg.Port)
	assert.Equal(t, []string{"string1", "string2", "string3"}, cfg.Strings)
}

func TestParsesEnvInner(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "someinnervalue", cfg.InnerStruct.Inner)
}

func TestParsesEnvInnerWithPrefix(t *testing.T) {
	os.Setenv("test_innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
		unexported:  &InnerStruct{},
	}
	assert.NoError(t, ParseWithPrefix(&cfg, "test_"))
	assert.Equal(t, "someinnervalue", cfg.InnerStruct.Inner)
}

func TestParsesEnvInnerNil(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := ParentStruct{}
	assert.NoError(t, Parse(&cfg))
}

func TestParsesEnvInnerInvalid(t *testing.T) {
	os.Setenv("innernum", "-547")
	defer os.Clearenv()
	cfg := ParentStruct{
		InnerStruct: &InnerStruct{},
	}
	assert.Error(t, Parse(&cfg))
}

func TestParsesEnvDerived(t *testing.T) {
	os.Setenv("innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := DerivedStruct{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "someinnervalue", cfg.Inner)
}

func TestParsesEnvDerivedWithPrefix(t *testing.T) {
	os.Setenv("test_innervar", "someinnervalue")
	defer os.Clearenv()
	cfg := DerivedStruct{}
	assert.NoError(t, ParseWithPrefix(&cfg, "test_"))
	assert.Equal(t, "someinnervalue", cfg.Inner)
}

func TestParsesEnvDerivedInvalid(t *testing.T) {
	os.Setenv("innernum", "-547")
	defer os.Clearenv()
	cfg := DerivedStruct{}
	assert.Error(t, Parse(&cfg))
}

func TestEmptyVars(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "", cfg.Some)
	assert.Equal(t, false, cfg.Other)
	assert.Equal(t, 0, cfg.Port)
	assert.Equal(t, uint(0), cfg.UintVal)
	assert.Equal(t, uint64(0), cfg.Uint64Val)
	assert.Equal(t, int64(0), cfg.Int64Val)
	assert.Equal(t, 0, len(cfg.Strings))
	assert.Equal(t, 0, len(cfg.SepStrings))
	assert.Equal(t, 0, len(cfg.Numbers))
	assert.Equal(t, 0, len(cfg.Bools))
}

func TestPassAnInvalidPtr(t *testing.T) {
	var thisShouldBreak int
	assert.Error(t, Parse(&thisShouldBreak))
}

func TestPassReference(t *testing.T) {
	cfg := Config{}
	assert.Error(t, Parse(cfg))
}

func TestInvalidBool(t *testing.T) {
	os.Setenv("othervar", "should-be-a-bool")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("PORT", "should-be-an-int")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidUint(t *testing.T) {
	os.Setenv("UINTVAL", "-44")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidFloat32(t *testing.T) {
	os.Setenv("FLOAT32", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidFloat64(t *testing.T) {
	os.Setenv("FLOAT64", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidUint64(t *testing.T) {
	os.Setenv("UINT64VAL", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidInt64(t *testing.T) {
	os.Setenv("INT64VAL", "AAA")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []int64 `env:"BADINTS"`
	}

	os.Setenv("BADINTS", "A,2,3")
	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func TestInvalidUInt64Slice(t *testing.T) {
	type config struct {
		BadFloats []uint64 `env:"BADINTS"`
	}

	os.Setenv("BADFLOATS", "A,2,3")
	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func TestInvalidFloat32Slice(t *testing.T) {
	type config struct {
		BadFloats []float32 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func TestInvalidFloat64Slice(t *testing.T) {
	type config struct {
		BadFloats []float64 `env:"BADFLOATS"`
	}

	os.Setenv("BADFLOATS", "A,2.0,3.0")
	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func TestInvalidBoolsSlice(t *testing.T) {
	type config struct {
		BadBools []bool `env:"BADBOOLS"`
	}

	os.Setenv("BADBOOLS", "t,f,TRUE,faaaalse")
	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func TestInvalidDuration(t *testing.T) {
	os.Setenv("DURATION", "should-be-a-valid-duration")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestInvalidDurations(t *testing.T) {
	os.Setenv("DURATIONS", "1s,contains-an-invalid-duration,3s")
	defer os.Clearenv()

	cfg := Config{}
	assert.Error(t, Parse(&cfg))
}

func TestParsesDefaultConfig(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Equal(t, "postgres://localhost:5432/db", cfg.DatabaseURL)
}

func TestParseStructWithoutEnvTag(t *testing.T) {
	cfg := Config{}
	assert.NoError(t, Parse(&cfg))
	assert.Empty(t, cfg.NotAnEnv)
}

func TestParseStructWithInvalidFieldKind(t *testing.T) {
	type config struct {
		WontWorkByte byte `env:"BLAH"`
	}
	os.Setenv("BLAH", "a")
	cfg := config{}
	assert.Error(t, Parse(&cfg))
}

func TestUnsupportedSliceType(t *testing.T) {
	type config struct {
		WontWork []map[int]int `env:"WONTWORK"`
	}

	os.Setenv("WONTWORK", "1,2,3")
	defer os.Clearenv()

	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func TestBadSeparator(t *testing.T) {
	type config struct {
		WontWork []int `env:"WONTWORK" envSeparator:":"`
	}

	cfg := &config{}
	os.Setenv("WONTWORK", "1,2,3,4")
	defer os.Clearenv()

	assert.Error(t, Parse(cfg))
}

func TestNoErrorRequiredSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED" required:"true"`
	}

	cfg := &config{}

	os.Setenv("IS_REQUIRED", "val")
	defer os.Clearenv()
	assert.NoError(t, Parse(cfg))
	assert.Equal(t, "val", cfg.IsRequired)
}

func TestNoErrorRequiredSetWithPrefix(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED" required:"true"`
	}

	cfg := &config{}

	os.Setenv("MYCLIENT_IS_REQUIRED", "val")
	defer os.Clearenv()
	assert.NoError(t, ParseWithPrefix(cfg, "MYCLIENT_"))
	assert.Equal(t, "val", cfg.IsRequired)
}

func TestErrorRequiredNotSet(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED" required:"True"`
	}

	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func TestErrorRequiredNotSetWithPrefix(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED" required:"True"`
	}

	cfg := &config{}
	assert.Error(t, ParseWithPrefix(cfg, "CLIENT_"))
}

func TestErrorRequiredNotValid(t *testing.T) {
	type config struct {
		IsRequired string `env:"IS_REQUIRED" required:"cat"`
	}

	cfg := &config{}
	err := Parse(cfg)
	assert.Error(t, err)
	// Error now includes field context
	assert.Contains(t, err.Error(), "invalid required tag \"cat\"")
	assert.Contains(t, err.Error(), "field 'IsRequired'")
}

func TestParseExpandOption(t *testing.T) {
	type config struct {
		Host        string `env:"HOST" envDefault:"localhost"`
		Port        int    `env:"PORT" envDefault:"3000" envExpand:"True"`
		SecretKey   string `env:"SECRET_KEY" envExpand:"True"`
		ExpandKey   string `env:"EXPAND_KEY"`
		CompoundKey string `env:"HOST_PORT" envDefault:"${HOST}:${PORT}" envExpand:"True"`
		Default     string `env:"DEFAULT" envDefault:"def1"  envExpand:"True"`
	}
	defer os.Clearenv()

	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "3000")
	os.Setenv("EXPAND_KEY", "qwerty12345")
	os.Setenv("SECRET_KEY", "${EXPAND_KEY}")

	cfg := config{}
	err := Parse(&cfg)

	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 3000, cfg.Port)
	assert.Equal(t, "qwerty12345", cfg.SecretKey)
	assert.Equal(t, "qwerty12345", cfg.ExpandKey)
	assert.Equal(t, "localhost:3000", cfg.CompoundKey)
	assert.Equal(t, "def1", cfg.Default)
}

func TestCustomParser(t *testing.T) {
	type foo struct {
		name string
	}

	type config struct {
		Var foo `env:"VAR"`
	}

	os.Setenv("VAR", "test")

	customParserFunc := func(v string) (interface{}, error) {
		return foo{name: v}, nil
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(foo{}): customParserFunc,
	})

	assert.NoError(t, err)
	assert.Equal(t, cfg.Var.name, "test")
}

func TestParseWithFuncsNoPtr(t *testing.T) {
	type foo struct{}
	err := ParseWithFuncs(foo{}, nil)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNotAStructPtr)
}

func TestParseWithFuncsInvalidType(t *testing.T) {
	var c int
	err := ParseWithFuncs(&c, nil)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNotAStructPtr)
}

func TestCustomParserError(t *testing.T) {
	type foo struct {
		name string
	}

	type config struct {
		Var foo `env:"VAR"`
	}

	os.Setenv("VAR", "test")

	customParserFunc := func(v string) (interface{}, error) {
		return nil, errors.New("something broke")
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(foo{}): customParserFunc,
	})

	assert.Empty(t, cfg.Var.name, "Var.name should not be filled out when parse errors")
	assert.Error(t, err)
	// Error message now includes field context
	assert.Contains(t, err.Error(), "custom parser error: something broke")
	assert.Contains(t, err.Error(), "field 'Var'")
}

func TestCustomParserBasicType(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_VAL"`
	}

	exp := ConstT(123)
	os.Setenv("CONST_VAL", fmt.Sprintf("%d", exp))

	customParserFunc := func(v string) (interface{}, error) {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		r := ConstT(i) //nolint:gosec
		return r, nil
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	assert.NoError(t, err)
	assert.Equal(t, exp, cfg.Const)
}

func TestCustomParserUint64Alias(t *testing.T) {
	type T uint64

	var one T = 1

	type config struct {
		Val T `env:"VAL" envDefault:"1x"`
	}

	parserCalled := false

	tParser := func(value string) (interface{}, error) {
		parserCalled = true
		trimmed := strings.TrimSuffix(value, "x")
		i, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, err
		}
		return T(i), nil //nolint:gosec
	}

	cfg := config{}

	err := ParseWithFuncs(&cfg, CustomParsers{
		reflect.TypeOf(one): tParser,
	})

	assert.True(t, parserCalled, "tParser should have been called")
	assert.NoError(t, err)
	assert.Equal(t, T(1), cfg.Val)
}

func TypeCustomParserBasicInvalid(t *testing.T) { //nolint: unused
	type ConstT int32 //nolint: unused

	type config struct { //nolint: unused
		Const ConstT `env:"CONST_VAL"`
	}

	os.Setenv("CONST_VAL", "foobar")

	expErr := errors.New("Random error")
	customParserFunc := func(_ string) (interface{}, error) {
		return nil, expErr
	}

	cfg := &config{}
	err := ParseWithFuncs(cfg, map[reflect.Type]ParserFunc{
		reflect.TypeOf(ConstT(0)): customParserFunc,
	})

	assert.Empty(t, cfg.Const)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func TestCustomParserNotCalledForNonAlias(t *testing.T) {
	type T uint64
	type U uint64

	type config struct {
		Val   uint64 `env:"VAL" envDefault:"33"`
		Other U      `env:"OTHER" envDefault:"44"`
	}

	tParserCalled := false

	tParser := func(value string) (interface{}, error) {
		tParserCalled = true
		return T(99), nil
	}

	cfg := config{}

	err := ParseWithFuncs(&cfg, CustomParsers{
		reflect.TypeOf(T(0)): tParser,
	})

	assert.False(t, tParserCalled, "tParser should not have been called")
	assert.NoError(t, err)
	assert.Equal(t, uint64(33), cfg.Val)
	assert.Equal(t, U(44), cfg.Other)
}

func TestCustomParserBasicUnsupported(t *testing.T) {
	type ConstT int32

	type config struct {
		Const ConstT `env:"CONST_VAL"`
	}

	exp := ConstT(123)
	os.Setenv("CONST_VAL", fmt.Sprintf("%d", exp))

	cfg := &config{}
	err := Parse(cfg)

	assert.Zero(t, cfg.Const)
	assert.Error(t, err)
	// Error now includes field context and wraps the original error
	if parseErrors, ok := err.(ParseErrors); ok && len(parseErrors) == 1 {
		assert.ErrorIs(t, parseErrors[0], ErrUnsupportedType)
		assert.Contains(t, parseErrors[0].Error(), "field 'Const'")
	} else {
		t.Fatal("Expected ParseErrors")
	}
}

func TestUnsupportedStructType(t *testing.T) {
	type config struct {
		Foo http.Client `env:"FOO"`
	}

	os.Setenv("FOO", "foo")

	cfg := &config{}
	err := Parse(cfg)

	assert.Error(t, err)
	// Error now includes field context and wraps the original error
	if parseErrors, ok := err.(ParseErrors); ok && len(parseErrors) == 1 {
		assert.ErrorIs(t, parseErrors[0], ErrUnsupportedType)
		assert.Contains(t, parseErrors[0].Error(), "field 'Foo'")
	} else {
		t.Fatal("Expected ParseErrors")
	}
}

func TestTextUnmarshalerError(t *testing.T) {
	type config struct {
		Unmarshaler unmarshaler `env:"UNMARSHALER"`
	}
	os.Setenv("UNMARSHALER", "invalid")
	cfg := &config{}
	assert.Error(t, Parse(cfg))
}

func ExampleParse() {
	type config struct {
		Home         string `env:"HOME"`
		Port         int    `env:"PORT" envDefault:"3000"`
		IsProduction bool   `env:"PRODUCTION"`
	}
	os.Setenv("HOME", "/tmp/fakehome")
	cfg := config{}
	_ = Parse(&cfg)
	fmt.Println(cfg)
	// Output: {/tmp/fakehome 3000 false}
}

// Test ParseErrors functionality
func TestParseErrors(t *testing.T) {
	// Test empty ParseErrors
	var emptyErrors ParseErrors
	assert.Equal(t, "", emptyErrors.Error())

	// Test single error
	singleErrors := ParseErrors{errors.New("single error")}
	assert.Equal(t, "single error", singleErrors.Error())

	// Test multiple errors
	multipleErrors := ParseErrors{
		errors.New("first error"),
		errors.New("second error"),
		errors.New("third error"),
	}
	expected := `multiple parsing errors (3):
  1. first error
  2. second error
  3. third error`
	assert.Equal(t, expected, multipleErrors.Error())
}

// Test ParseErrors in actual parsing scenario
func TestParseMultipleErrors(t *testing.T) {
	type config struct {
		RequiredVar1 string `env:"REQUIRED_VAR1" required:"true"`
		RequiredVar2 int    `env:"REQUIRED_VAR2" required:"true"`
		BadInt       int    `env:"BAD_INT_VAR"`
	}

	// Set up environment to cause multiple errors
	os.Unsetenv("REQUIRED_VAR1")
	os.Unsetenv("REQUIRED_VAR2")
	os.Setenv("BAD_INT_VAR", "not_a_number")

	cfg := &config{}
	err := Parse(cfg)

	// Should get a ParseErrors with multiple issues
	assert.Error(t, err)
	parseErrors, ok := err.(ParseErrors)
	assert.True(t, ok, "Error should be of type ParseErrors")
	assert.Greater(t, len(parseErrors), 1, "Should have multiple errors")

	// Check that error message contains information about multiple errors
	errMsg := err.Error()
	assert.Contains(t, errMsg, "multiple parsing errors")
}

// Test prefix validation
func TestPrefixValidation(t *testing.T) {
	type config struct {
		Value string `env:"VALUE"`
	}

	os.Setenv("TEST_VALUE", "test")
	defer os.Unsetenv("TEST_VALUE")

	// Valid prefix (ends with underscore)
	cfg := &config{}
	err := ParseWithPrefix(cfg, "TEST_")
	assert.NoError(t, err)
	assert.Equal(t, "test", cfg.Value)

	// Empty prefix is valid
	cfg = &config{}
	err = ParseWithPrefix(cfg, "")
	assert.NoError(t, err)

	// Invalid prefix (doesn't end with underscore)
	cfg = &config{}
	err = ParseWithPrefix(cfg, "TEST")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "prefix must end with underscore")
}

// Test GetAllVars
func TestGetAllVars(t *testing.T) {
	type config struct {
		Required string `env:"REQUIRED" required:"true"`
		Optional string `env:"OPTIONAL"`
		WithDefault string `env:"WITH_DEFAULT" envDefault:"default_value"`
		Number int `env:"NUMBER" required:"true"`
	}

	vars, err := GetAllVars(&config{}, "")
	assert.NoError(t, err)
	assert.Len(t, vars, 4)

	// Check required field
	requiredVar := findVar(vars, "REQUIRED")
	assert.NotNil(t, requiredVar)
	assert.True(t, requiredVar.Required)
	assert.Equal(t, "Required", requiredVar.FieldName)
	assert.Equal(t, "string", requiredVar.Type)

	// Check optional field
	optionalVar := findVar(vars, "OPTIONAL")
	assert.NotNil(t, optionalVar)
	assert.False(t, optionalVar.Required)

	// Check field with default
	defaultVar := findVar(vars, "WITH_DEFAULT")
	assert.NotNil(t, defaultVar)
	assert.True(t, defaultVar.HasDefault)
	assert.Equal(t, "default_value", defaultVar.Default)

	// Check number field
	numberVar := findVar(vars, "NUMBER")
	assert.NotNil(t, numberVar)
	assert.True(t, numberVar.Required)
	assert.Equal(t, "int", numberVar.Type)
}

// Test GetAllVars with prefix
func TestGetAllVarsWithPrefix(t *testing.T) {
	type config struct {
		Value string `env:"VALUE"`
	}

	vars, err := GetAllVars(&config{}, "PREFIX_")
	assert.NoError(t, err)
	assert.Len(t, vars, 1)
	assert.Equal(t, "PREFIX_VALUE", vars[0].Name)
}

// Test GetRequiredVars
func TestGetRequiredVars(t *testing.T) {
	type config struct {
		Required1 string `env:"REQUIRED1" required:"true"`
		Optional string `env:"OPTIONAL"`
		Required2 int `env:"REQUIRED2" required:"true"`
	}

	required, err := GetRequiredVars(&config{}, "")
	assert.NoError(t, err)
	assert.Len(t, required, 2)
	assert.Contains(t, required, "REQUIRED1")
	assert.Contains(t, required, "REQUIRED2")
	assert.NotContains(t, required, "OPTIONAL")
}

// Test ValidateRequired
func TestValidateRequired(t *testing.T) {
	type config struct {
		Required1 string `env:"REQUIRED1" required:"true"`
		Optional string `env:"OPTIONAL"`
		Required2 int `env:"REQUIRED2" required:"true"`
	}

	// All required vars set
	os.Setenv("REQUIRED1", "value1")
	os.Setenv("REQUIRED2", "42")
	defer os.Unsetenv("REQUIRED1")
	defer os.Unsetenv("REQUIRED2")

	err := ValidateRequired(&config{}, "")
	assert.NoError(t, err)

	// Missing required var
	os.Unsetenv("REQUIRED1")
	err = ValidateRequired(&config{}, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required environment variables")
	assert.Contains(t, err.Error(), "REQUIRED1")
}

// Test EnableDebugLogging
func TestEnableDebugLogging(t *testing.T) {
	type config struct {
		Value string `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "test")
	defer os.Unsetenv("TEST_VALUE")

	var logMessages []string
	logger := func(format string, args ...interface{}) {
		logMessages = append(logMessages, fmt.Sprintf(format, args...))
	}

	EnableDebugLogging(logger)
	defer EnableDebugLogging(nil) // Clean up

	cfg := &config{}
	err := Parse(cfg)
	assert.NoError(t, err)

	// Verify debug logging was called
	assert.NotEmpty(t, logMessages)
	assert.Contains(t, logMessages[0], "TEST_VALUE")
	assert.Contains(t, logMessages[0], "test")
}

// Test nested structs with GetAllVars
func TestGetAllVarsNested(t *testing.T) {
	type NestedConfig struct {
		NestedValue string `env:"NESTED_VALUE"`
	}

	type Config struct {
		TopValue string `env:"TOP_VALUE"`
		Nested NestedConfig
	}

	vars, err := GetAllVars(&Config{}, "")
	assert.NoError(t, err)
	assert.Len(t, vars, 2)

	topVar := findVar(vars, "TOP_VALUE")
	assert.NotNil(t, topVar)
	assert.Equal(t, "TopValue", topVar.FieldName)

	nestedVar := findVar(vars, "NESTED_VALUE")
	assert.NotNil(t, nestedVar)
	assert.Equal(t, "NestedValue", nestedVar.FieldName)
	assert.Contains(t, nestedVar.FieldPath, "Nested")
}

// Helper function to find a var by name
func findVar(vars []VarInfo, name string) *VarInfo {
	for _, v := range vars {
		if v.Name == name {
			return &v
		}
	}
	return nil
}
