// Package env provides utilities for working with environment variables and loading
// configuration from .env files.
//
// This package offers three main areas of functionality:
//
// 1. Simple environment variable access with type conversion
// 2. Struct-based configuration parsing using struct tags
// 3. Loading and parsing of .env files
//
// # Basic Usage
//
// Get environment variables with automatic type conversion:
//
//	port, err := env.GetInt("PORT")
//	if err != nil {
//		port = env.GetOrInt("PORT", 8080) // with default
//	}
//
//	// Or panic if missing/invalid
//	dbHost := env.MustGet("DATABASE_HOST")
//
// # Struct-based Configuration
//
// Parse environment variables into structs using tags:
//
//	type Config struct {
//		Host string `env:"HOST" envDefault:"localhost"`
//		Port int    `env:"PORT" envDefault:"8080"`
//		Debug bool  `env:"DEBUG"`
//	}
//
//	var cfg Config
//	if err := env.Parse(&cfg); err != nil {
//		log.Fatal(err)
//	}
//
// # File Loading
//
// Load environment variables from .env files:
//
//	// Load from .env file
//	err := env.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Load from specific files
//	err = env.Load("config.env", "local.env")
//
// # Supported Types
//
// The package supports automatic conversion for:
//   - bool, int, uint, int64, uint64, float32, float64
//   - time.Duration (using time.ParseDuration format)
//   - *url.URL (using url.ParseRequestURI)
//   - []string and other slice types (comma-separated by default)
//   - Any type implementing encoding.TextUnmarshaler
//
// # Struct Tags
//
// When using Parse functions, the following struct tags are supported:
//   - env:"VAR_NAME" - specifies the environment variable name
//   - envDefault:"value" - provides a default value if the variable is not set
//   - required:"true" - makes the field required (parsing fails if missing)
//   - envSeparator:";" - custom separator for slice types (default is comma)
//   - envExpand:"true" - enables variable expansion using os.ExpandEnv
//
// # Error Handling
//
// Functions come in three variants for different error handling approaches:
//   - Get*() functions return (value, error)
//   - GetOr*() functions return value with a fallback default
//   - MustGet*() functions panic if the variable is missing or invalid
package env