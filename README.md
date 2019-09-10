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
a prfic when parsing.
```
	cfg := ClientConfig{}
	err := env.ParseWithPrefix(&cfg, "CLIENT2_")
```

In this case the parser would look for the envionment variables CLIENT2_ENDPOINT, CLIENT2_HEALTH_CHECK, etc.

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


