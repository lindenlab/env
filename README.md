# env

## Credit 
This repo is a heavily modified version of env by https://github.com/caarlos0 .  Also major portions we adapted from github.com/joho/godotenv.

## env utility functions

This repo provides some simple utilities for dealing with environment variables.  There are some basic functions that are simple wrappers to the "os" package's
env functions.  These include:

```
err   = env.Set("KEY", "secret_key")
err   = env.Unset("KEY")
value = env.Get("OTHER_KEY")
value = env.GetOr("OTHER_KEY", "value returned if OTHER_KEY does not exist")
value = env.MustGet("KEY")  // panics if KEY does not exist
```

## Parse environment vars into a struct

A much more powerful approach is to populate an annotated struct with
values from the environment.

Take the following example:
```
type ClientConfig struct {
	Endpoint                  []string      `env:"ENDPOINT" required:"true"`
	HealthCheck               bool          `env:"HEALTH_CHECK" envDefault:"true"`
	HealthCheckTimeout        time.Duration `env:"HEALTH_CHECK_TIMEOUT" envDefault:"1s"`
}
```

This anotated struct can then be populated with the appropiate environment
variables but using the following command:
```
	cfg := ClientConfig{}
	err := env.Parse(&cfg)
```

## Using a prefix

If the struct you are populating is part of a library you may want to have
parse the same struct but with different values.  You can do this by using
a prefix when parsing.
```
	cfg := ClientConfig{}
	err := env.ParseWithPrefix(&cfg, "CLIENT2_")
```

In this case the parser would look for the environment variables CLIENT2_ENDPOINT, CLIENT2_HEALTH_CHECK, etc.

**Note:** Prefixes must end with an underscore (`_`). If you provide a prefix without a trailing underscore, Parse will return an error.

### Supported types and defaults

The environment variables are Parsed to go into the appropiate types (or
an err is thrown if it can not do so).  This sames you from writing a some
basic validations on envirnment varaibles.

The library has built-in support for the following types:

* `string`
* `int`
* `uint`
* `int64`
* `bool`
* `float32`
* `float64`
* `time.Duration`
* `url.URL`
* `[]string`
* `[]int`
* `[]bool`
* `[]float32`
* `[]float64`
* `[]time.Duration`
* `[]url.URL`

### Optional tags

The required tag will cause an error if the environment variable does not exist:
``` `env:"ENDPOINT" required:"true"` ```

The envDefault tag will allow you to provide a default value to use if the environment variable does not exist:
``` `env:"HEALTH_CHECK" envDefault:"true"` ```

When assigning to a slice type the "," is used to seperate fields.  You can override this with the envSeparator:":" tag to use some other character.

## Advanced Features

### Context-aware panics

Use `MustGetWithContext` for better error messages when required variables are missing:

```go
apiBase := env.MustGetWithContext("WALLET_TRANSACTION_API", "generating wallet URIs")
// Panics with: "required env var WALLET_TRANSACTION_API missing (needed for: generating wallet URIs)"
```

### Test helpers

The package provides convenient functions for managing environment variables in tests:

```go
func TestConfig(t *testing.T) {
    // Automatically restores or removes the variable after the test
    env.SetForTest(t, "WALLET_TRANSACTION_API", "https://test.com")

    // Test code here
    // Cleanup happens automatically via t.Cleanup()
}

func TestWithoutEnv(t *testing.T) {
    // Temporarily removes a variable and restores it after test
    env.UnsetForTest(t, "OPTIONAL_CONFIG")

    // Test code that expects variable to not exist
}
```

### Configuration validation

Validate that all required environment variables are set before parsing:

```go
// Check required vars without parsing values
if err := env.ValidateRequired(&Config{}, ""); err != nil {
    log.Fatalf("Missing required configuration: %v", err)
}

// Now safe to parse
cfg := &Config{}
env.Parse(&cfg)
```

### Environment variable discovery

Get information about all environment variables that a struct uses:

```go
// Get all variables (required and optional)
vars, _ := env.GetAllVars(&Config{}, "")
for _, v := range vars {
    fmt.Printf("%s: %s (required: %v, default: %s)\n",
        v.Name, v.Type, v.Required, v.Default)
}

// Get only required variables
required, _ := env.GetRequiredVars(&Config{}, "")
fmt.Println("Required vars:", required)
```

This is useful for:
- Generating documentation
- Creating example .env files
- Debugging configuration issues
- Validating deployment environments

### Debug logging

Enable debug logging to see what environment variables are being read:

```go
// Enable debug logging
env.EnableDebugLogging(log.Printf)

cfg := &Config{}
env.Parse(&cfg)
// Logs: "env: DB_CONNECTION = postgres://... (field: Config.Database)"

// Disable debug logging
env.EnableDebugLogging(nil)
```

This is particularly useful for:
- Debugging configuration issues in production
- Understanding what values are being loaded
- Troubleshooting environment variable naming

### Better error messages

Parse errors now include field context for easier debugging:

```go
type Config struct {
    Connection string `env:"DB_CONNECTION" required:"true"`
}

err := env.Parse(&Config{})
// Error: "field 'Connection' in Config: env var DB_CONNECTION was missing and is required"
```

