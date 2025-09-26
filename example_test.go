package env_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lindenlab/env"
)

func ExampleGet() {
	os.Setenv("HOME", "/home/user")
	fmt.Println(env.Get("HOME"))
	fmt.Println(env.Get("NONEXISTENT"))
	// Output:
	// /home/user
	//
}

func ExampleGetOr() {
	os.Setenv("PORT", "8080")
	fmt.Println(env.GetOr("PORT", "3000"))
	fmt.Println(env.GetOr("MISSING_PORT", "3000"))
	// Output:
	// 8080
	// 3000
}

func ExampleGetInt() {
	os.Setenv("PORT", "8080")
	os.Setenv("INVALID_PORT", "not-a-number")

	port, err := env.GetInt("PORT")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port:", port)

	_, err = env.GetInt("INVALID_PORT")
	fmt.Println("Error:", err != nil)
	// Output:
	// Port: 8080
	// Error: true
}

func ExampleGetOrInt() {
	os.Setenv("PORT", "8080")
	fmt.Println("Existing:", env.GetOrInt("PORT", 3000))
	fmt.Println("Missing:", env.GetOrInt("MISSING_PORT", 3000))
	fmt.Println("Invalid:", env.GetOrInt("INVALID_PORT", 3000))
	// Output:
	// Existing: 8080
	// Missing: 3000
	// Invalid: 3000
}

func ExampleGetBool() {
	os.Setenv("DEBUG", "true")
	os.Setenv("VERBOSE", "1")
	os.Setenv("QUIET", "false")

	debug, _ := env.GetBool("DEBUG")
	verbose, _ := env.GetBool("VERBOSE")
	quiet, _ := env.GetBool("QUIET")

	fmt.Println("Debug:", debug)
	fmt.Println("Verbose:", verbose)
	fmt.Println("Quiet:", quiet)
	// Output:
	// Debug: true
	// Verbose: true
	// Quiet: false
}

func ExampleGetDuration() {
	os.Setenv("TIMEOUT", "30s")
	os.Setenv("INTERVAL", "5m30s")

	timeout, _ := env.GetDuration("TIMEOUT")
	interval, _ := env.GetDuration("INTERVAL")

	fmt.Println("Timeout:", timeout)
	fmt.Println("Interval:", interval)
	// Output:
	// Timeout: 30s
	// Interval: 5m30s
}

func ExampleLoad() {
	// Create a temporary .env file
	envContent := `# Database configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=myapp

# Application settings
DEBUG=true
LOG_LEVEL=info`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(".env")

	// Load the .env file
	err = env.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Access the loaded variables
	fmt.Println("Database Host:", env.Get("DATABASE_HOST"))
	fmt.Println("Database Port:", env.GetOrInt("DATABASE_PORT", 0))
	fmt.Println("Debug Mode:", env.GetOrBool("DEBUG", false))
	// Output:
	// Database Host: localhost
	// Database Port: 5432
	// Debug Mode: true
}

func ExampleParse() {
	// Set up environment variables
	os.Setenv("HOME", "/tmp/fakehome")
	os.Setenv("PORT", "8080")
	os.Setenv("DEBUG", "true")
	os.Setenv("TAGS", "web,api,database")
	os.Setenv("TIMEOUT", "30s")

	type Config struct {
		Home     string        `env:"HOME"`
		Port     int           `env:"PORT" envDefault:"3000"`
		Debug    bool          `env:"DEBUG"`
		Tags     []string      `env:"TAGS" envSeparator:","`
		Timeout  time.Duration `env:"TIMEOUT" envDefault:"10s"`
		Version  string        `env:"VERSION" envDefault:"1.0.0"`
		Required string        `env:"REQUIRED_VAR" required:"true"`
	}

	// This will fail because REQUIRED_VAR is not set
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Error:", err != nil)
	}

	// Set the required variable and try again
	os.Setenv("REQUIRED_VAR", "important-value")
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Home: %s\n", cfg.Home)
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("Debug: %t\n", cfg.Debug)
	fmt.Printf("Tags: %v\n", cfg.Tags)
	fmt.Printf("Timeout: %v\n", cfg.Timeout)
	fmt.Printf("Version: %s\n", cfg.Version)
	fmt.Printf("Required: %s\n", cfg.Required)
	// Output:
	// Error: true
	// Home: /tmp/fakehome
	// Port: 8080
	// Debug: true
	// Tags: [web api database]
	// Timeout: 30s
	// Version: 1.0.0
	// Required: important-value
}

func ExampleParseWithPrefix() {
	// Set up environment variables with prefixes
	os.Setenv("CLIENT1_HOST", "api.example.com")
	os.Setenv("CLIENT1_PORT", "443")
	os.Setenv("CLIENT2_HOST", "internal.example.com")
	os.Setenv("CLIENT2_PORT", "8080")

	type ClientConfig struct {
		Host string `env:"HOST" envDefault:"localhost"`
		Port int    `env:"PORT" envDefault:"80"`
	}

	var client1, client2 ClientConfig

	env.ParseWithPrefix(&client1, "CLIENT1_")
	env.ParseWithPrefix(&client2, "CLIENT2_")

	fmt.Printf("Client 1: %s:%d\n", client1.Host, client1.Port)
	fmt.Printf("Client 2: %s:%d\n", client2.Host, client2.Port)
	// Output:
	// Client 1: api.example.com:443
	// Client 2: internal.example.com:8080
}

func ExampleRead() {
	// Create a temporary .env file
	envContent := `DATABASE_URL=postgres://localhost/myapp
REDIS_URL=redis://localhost:6379
API_KEY=secret123
DEBUG=true`

	err := os.WriteFile("config.env", []byte(envContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("config.env")

	// Read the file without setting environment variables
	envMap, err := env.Read("config.env")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database URL:", envMap["DATABASE_URL"])
	fmt.Println("Redis URL:", envMap["REDIS_URL"])
	fmt.Println("Debug:", envMap["DEBUG"])
	// Output:
	// Database URL: postgres://localhost/myapp
	// Redis URL: redis://localhost:6379
	// Debug: true
}

func ExampleMarshal() {
	envMap := map[string]string{
		"DATABASE_URL": "postgres://localhost/myapp",
		"DEBUG":        "true",
		"PORT":         "8080",
		"API_KEY":      "secret with spaces",
	}

	content, err := env.Marshal(envMap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(content)
	// Output:
	// API_KEY="secret with spaces"
	// DATABASE_URL="postgres://localhost/myapp"
	// DEBUG="true"
	// PORT="8080"
}