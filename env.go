package env

import (
	"encoding"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Constants for parsing operations (shared with util.go)
const (
	DecimalBase = 10
	Int32Bits   = 32
	Int64Bits   = 64
	Float32Bits = 32
	Float64Bits = 64
)

var (
	// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
	// Struct to Parse
	ErrNotAStructPtr = errors.New("expected a pointer to a Struct")
	// ErrUnsupportedType if the struct field type is not supported by env
	ErrUnsupportedType = errors.New("type is not supported")
	// ErrUnsupportedSliceType if the slice element type is not supported by env
	ErrUnsupportedSliceType = errors.New("unsupported slice type")
)

// ParseErrors represents multiple errors that occurred during parsing
type ParseErrors []error

// Error implements the error interface for ParseErrors
func (pe ParseErrors) Error() string {
	if len(pe) == 0 {
		return ""
	}
	if len(pe) == 1 {
		return pe[0].Error()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("multiple parsing errors (%d):", len(pe)))
	for i, err := range pe {
		sb.WriteString(fmt.Sprintf("\n  %d. %s", i+1, err.Error()))
	}
	return sb.String()
}

var (
	// OnEnvVarSet is an optional convenience callback, such as for logging purposes.
	// If not nil, it's called after successfully setting the given field from the given value.
	OnEnvVarSet func(reflect.StructField, string)
	// Friendly names for reflect types
	sliceOfInts      = reflect.TypeOf([]int(nil))
	sliceOfInt64s    = reflect.TypeOf([]int64(nil))
	sliceOfUint64s   = reflect.TypeOf([]uint64(nil))
	sliceOfStrings   = reflect.TypeOf([]string(nil))
	sliceOfBools     = reflect.TypeOf([]bool(nil))
	sliceOfFloat32s  = reflect.TypeOf([]float32(nil))
	sliceOfFloat64s  = reflect.TypeOf([]float64(nil))
	sliceOfDurations = reflect.TypeOf([]time.Duration(nil))
	sliceOfURLs      = reflect.TypeOf([]url.URL(nil))
)

// CustomParsers maps Go types to custom parsing functions.
// It allows you to provide custom logic for parsing environment variables
// into specific types that aren't supported by default.
//
// The key is the reflect.Type of the target type, and the value is a ParserFunc
// that knows how to convert a string to that type.
type CustomParsers map[reflect.Type]ParserFunc

// ParserFunc defines the signature of a custom parsing function.
// It takes a string value from an environment variable and returns
// the parsed value as an interface{} and any parsing error.
//
// The returned value should be of the type that the parser is designed to handle.
type ParserFunc func(v string) (interface{}, error)

// Parse populates a struct's fields from environment variables.
// The struct fields must be tagged with `env:"VAR_NAME"` to specify
// which environment variable to read.
//
// Supported struct tags:
//   - env:"VAR_NAME" - specifies the environment variable name (required)
//   - envDefault:"value" - default value if the environment variable is not set
//   - required:"true" - makes the field required (causes error if missing)
//   - envSeparator:"," - separator for slice types (default is comma)
//   - envExpand:"true" - enables variable expansion using os.ExpandEnv
//
// The function supports nested structs and pointers to structs.
// It returns an error if required fields are missing or if type conversion fails.
func Parse(v interface{}) error {
	return ParseWithPrefixFuncs(v, "", make(map[reflect.Type]ParserFunc))
}

// ParseWithPrefix populates a struct's fields from environment variables with a prefix.
// This is useful for loading different configurations for the same struct type.
//
// For example, with prefix "CLIENT2_", a field tagged `env:"ENDPOINT"` will
// read from the environment variable "CLIENT2_ENDPOINT".
//
// See Parse for details on supported struct tags and behavior.
func ParseWithPrefix(v interface{}, prefix string) error {
	return ParseWithPrefixFuncs(v, prefix, make(map[reflect.Type]ParserFunc))
}

// ParseWithFuncs populates a struct's fields from environment variables,
// using custom parsing functions for specific types.
//
// This allows you to handle types that aren't supported by default.
// The funcMap parameter maps reflect.Type values to ParserFunc implementations.
//
// See Parse for details on supported struct tags and behavior.
func ParseWithFuncs(v interface{}, funcMap CustomParsers) error {
	return ParseWithPrefixFuncs(v, "", funcMap)
}

// ParseWithPrefixFuncs populates a struct's fields from environment variables
// with both a prefix and custom parsing functions.
//
// This combines the functionality of ParseWithPrefix and ParseWithFuncs,
// allowing both prefixed variable names and custom type parsing.
//
// See Parse for details on supported struct tags and behavior.
func ParseWithPrefixFuncs(v interface{}, prefix string, funcMap CustomParsers) error {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}
	return doParse(ref, prefix, funcMap)
}

func doParse(ref reflect.Value, prefix string, funcMap CustomParsers) error {
	refType := ref.Type()
	var parseErrors ParseErrors

	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		if reflect.Ptr == refField.Kind() && !refField.IsNil() && refField.CanSet() {
			err := ParseWithPrefixFuncs(refField.Interface(), prefix, funcMap)
			if nil != err {
				return err
			}
			continue
		}
		refTypeField := refType.Field(i)
		value, err := get(refTypeField, prefix)
		if err != nil {
			parseErrors = append(parseErrors, err)
			continue
		}
		if value == "" {
			if reflect.Struct == refField.Kind() {
				if err := doParse(refField, prefix, funcMap); err != nil {
					parseErrors = append(parseErrors, err)
				}
			}
			continue
		}
		if err := set(refField, refTypeField, value, funcMap); err != nil {
			parseErrors = append(parseErrors, err)
			continue
		}
		if OnEnvVarSet != nil {
			OnEnvVarSet(refTypeField, value)
		}
	}
	if len(parseErrors) == 0 {
		return nil
	}
	return parseErrors
}

func get(field reflect.StructField, prefix string) (string, error) {
	key := prefix + field.Tag.Get("env")

	var envRequired = false
	reqTag, hasRequiredTag := field.Tag.Lookup("required")
	if hasRequiredTag {
		var b bool
		var err error
		if b, err = strconv.ParseBool(reqTag); err != nil {
			// The value provided for the required tag is not a valid
			// Boolean, so inform the user.
			return "", fmt.Errorf("invalid required tag %q: %v", reqTag, err)
		}
		if b {
			envRequired = true
		}
	}

	value, envFound := os.LookupEnv(key)
	if !envFound && envRequired {
		return "", fmt.Errorf("env var %s was missing and is required", key)
	}

	if !envFound {
		// apply default if one exists
		value = field.Tag.Get("envDefault")
	}

	expandVar := field.Tag.Get("envExpand")
	if strings.ToLower(expandVar) == "true" {
		value = os.ExpandEnv(value)
	}

	return value, nil
}

func set(field reflect.Value, refType reflect.StructField, value string, funcMap CustomParsers) error {
	// use custom parser if configured for this type
	parserFunc, ok := funcMap[refType.Type]
	if ok {
		val, err := parserFunc(value)
		if err != nil {
			return fmt.Errorf("custom parser error: %v", err)
		}
		field.Set(reflect.ValueOf(val))
		return nil
	}

	if refType.Type == reflect.TypeOf(url.URL{}) {
		u, err := url.Parse(value)
		if err != nil {
			return fmt.Errorf("unable to complete URL parse: %v", err)
		}
		field.Set(reflect.ValueOf(*u))
		return nil
	}

	// fall back to built-in parsers
	switch field.Kind() {
	case reflect.Slice:
		separator := refType.Tag.Get("envSeparator")
		return handleSlice(field, value, separator)
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, DecimalBase, Int32Bits)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Uint:
		uintValue, err := strconv.ParseUint(value, DecimalBase, Int32Bits)
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, Float32Bits)
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, Float64Bits)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(v))
	case reflect.Int64:
		if refType.Type.String() == "time.Duration" {
			dValue, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(dValue))
		} else {
			intValue, err := strconv.ParseInt(value, DecimalBase, Int64Bits)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		}
	case reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, DecimalBase, Int64Bits)
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
	default:
		return handleTextUnmarshaler(field, value)
	}
	return nil
}

func handleSlice(field reflect.Value, value, separator string) error {
	if separator == "" {
		separator = ","
	}

	splitData := strings.Split(value, separator)

	switch field.Type() {
	case sliceOfStrings:
		field.Set(reflect.ValueOf(splitData))
	case sliceOfInts:
		intData, err := parseInts(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(intData))
	case sliceOfInt64s:
		int64Data, err := parseInt64s(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(int64Data))
	case sliceOfUint64s:
		uint64Data, err := parseUint64s(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(uint64Data))
	case sliceOfFloat32s:
		data, err := parseFloat32s(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(data))
	case sliceOfFloat64s:
		data, err := parseFloat64s(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(data))
	case sliceOfBools:
		boolData, err := parseBools(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(boolData))
	case sliceOfDurations:
		durationData, err := parseDurations(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(durationData))
	case sliceOfURLs:
		urlData, err := parseUrls(splitData)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(urlData))
	default:
		elemType := field.Type().Elem()
		// Ensure we test *type as we can always address elements in a slice.
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if _, ok := reflect.New(elemType).Interface().(encoding.TextUnmarshaler); !ok {
			return ErrUnsupportedSliceType
		}
		return parseTextUnmarshalers(field, splitData)

	}
	return nil
}

func handleTextUnmarshaler(field reflect.Value, value string) error {
	if reflect.Ptr == field.Kind() {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	} else if field.CanAddr() {
		field = field.Addr()
	}

	tm, ok := field.Interface().(encoding.TextUnmarshaler)
	if !ok {
		return ErrUnsupportedType
	}

	return tm.UnmarshalText([]byte(value))
}

func parseInts(data []string) ([]int, error) {
	intSlice := make([]int, 0, len(data))

	for _, v := range data {
		intValue, err := strconv.ParseInt(v, DecimalBase, Int32Bits)
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, int(intValue))
	}
	return intSlice, nil
}

func parseInt64s(data []string) ([]int64, error) {
	intSlice := make([]int64, 0, len(data))

	for _, v := range data {
		intValue, err := strconv.ParseInt(v, DecimalBase, Int64Bits)
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, int64(intValue))
	}
	return intSlice, nil
}

func parseUint64s(data []string) ([]uint64, error) {
	uintSlice := make([]uint64, 0, len(data))

	for _, v := range data {
		uintValue, err := strconv.ParseUint(v, DecimalBase, Int64Bits)
		if err != nil {
			return nil, err
		}
		uintSlice = append(uintSlice, uint64(uintValue))
	}
	return uintSlice, nil
}

func parseFloat32s(data []string) ([]float32, error) {
	float32Slice := make([]float32, 0, len(data))

	for _, v := range data {
		data, err := strconv.ParseFloat(v, Float32Bits)
		if err != nil {
			return nil, err
		}
		float32Slice = append(float32Slice, float32(data))
	}
	return float32Slice, nil
}

func parseFloat64s(data []string) ([]float64, error) {
	float64Slice := make([]float64, 0, len(data))

	for _, v := range data {
		data, err := strconv.ParseFloat(v, Float64Bits)
		if err != nil {
			return nil, err
		}
		float64Slice = append(float64Slice, float64(data))
	}
	return float64Slice, nil
}

func parseBools(data []string) ([]bool, error) {
	boolSlice := make([]bool, 0, len(data))

	for _, v := range data {
		bvalue, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}

		boolSlice = append(boolSlice, bvalue)
	}
	return boolSlice, nil
}

func parseDurations(data []string) ([]time.Duration, error) {
	durationSlice := make([]time.Duration, 0, len(data))

	for _, v := range data {
		dvalue, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}

		durationSlice = append(durationSlice, dvalue)
	}
	return durationSlice, nil
}

func parseUrls(data []string) ([]url.URL, error) {
	urlSlice := make([]url.URL, 0, len(data))

	for _, v := range data {
		uvalue, err := url.Parse(v)
		if err != nil {
			return nil, err
		}

		urlSlice = append(urlSlice, *uvalue)
	}
	return urlSlice, nil
}

func parseTextUnmarshalers(field reflect.Value, data []string) error {
	s := len(data)
	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(elemType), s, s)
	for i, v := range data {
		sv := slice.Index(i)
		kind := sv.Kind()
		if kind == reflect.Ptr {
			sv = reflect.New(elemType.Elem())
		} else {
			sv = sv.Addr()
		}
		tm := sv.Interface().(encoding.TextUnmarshaler)
		if err := tm.UnmarshalText([]byte(v)); err != nil {
			return err
		}
		if kind == reflect.Ptr {
			slice.Index(i).Set(sv)
		}
	}

	field.Set(slice)

	return nil
}
