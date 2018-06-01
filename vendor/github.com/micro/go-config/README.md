# Config [![GoDoc](https://godoc.org/github.com/micro/go-config?status.svg)](https://godoc.org/github.com/micro/go-config)

Go Config is a pluggable dynamic config library

Most config in applications are statically configured or include complex logic to load from multiple sources. Go-config makes this easy, 
pluggable and mergeable. You'll never have to deal with config in the same way again.

## Features

- Dynamic - load config on the fly as you need it
- Pluggable - choose which source to load from; file, envvar, consul
- Mergeable - merge and override multiple config sources
- Fallback - specify fallback values where keys don't exist
- Watch - Watch the config for changes

## Sources

The following sources for config are supported

- [consul](https://github.com/micro/go-config/tree/master/source/consul) - read from consul
- [envvar](https://github.com/micro/go-config/tree/master/source/envvar) - read from environment variables
- [file](https://github.com/micro/go-config/tree/master/source/file) - read from file
- [flag](https://github.com/micro/go-config/tree/master/source/flag) - read from flags
- [memory](https://github.com/micro/go-config/tree/master/source/memory) - read from memory
- [microcli](https://github.com/micro/go-config/tree/master/source/microcli) - read from micro cli flags

TODO:

- etcd
- vault
- kubernetes config map
- git url

## Config 

Top level config is an interface. It supports multiple sources, watching and fallback values.

### Interface

```go
type Config interface {
        Close() error
        Bytes() []byte
        Get(path ...string) reader.Value
        Load(source ...source.Source) error
        Watch(path ...string) (Watcher, error)
}
```

### Value

The config.Get method returns a reader.Value which can cast to any type with a fallback value

```go
type Value interface {
	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration
	StringSlice(def []string) []string
	StringMap(def map[string]string) map[string]string
	Scan(val interface{}) error
	Bytes() []byte
}
```

## Source

A [Source](https://godoc.org/github.com/micro/go-config/source#Source) is the source of config. 

It can be env vars, a file, a key value store. Anything which conforms to the Source interface.

### Interface

```go
// Source is the source from which config is loaded
type Source interface {
	Read() (*ChangeSet, error)
	Watch() (Watcher, error)
	String() string
}

// ChangeSet represents a set of changes from a source
type ChangeSet struct {
	Data      []byte
	Checksum  string
	Timestamp time.Time
	Source    string
}
```

### Format

Sources are currently expected return config as JSON to operate with the default config reader

The [Reader](https://godoc.org/github.com/micro/go-config/reader#Reader) defaults to json but can be swapped out to any other format.

```
{
	"path": {
		"to": {
			"key": ["foo", "bar"]
		}
	}
}
```

## Usage

Assuming the following config file

```json
{
    "hosts": {
        "database": {
            "address": "10.0.0.1",
            "port": 3306
        },
        "cache": {
            "address": "10.0.0.2",
            "port": 6379
        }
    }
}
```

### Load File

```
import "github.com/micro/go-config/source/file"

// Create new config
conf := config.NewConfig()

// Load file source
conf.Load(file.NewSource(
	file.WithPath("/tmp/config.json"),
))
```

### Scan

```go
type Host struct {
	Address string `json:"address"`
	Port int `json:"port"`
}

var host Host

conf.Get("hosts", "database").Scan(&host)

// 10.0.0.1 3306
fmt.Println(host.Address, host.Port)
```

### Go Vals

```go
// Get address. Set default to localhost as fallback
address := conf.Get("hosts", "database", "address").String("localhost")

// Get port. Set default to 3000 as fallback
port := conf.Get("hosts", "database", "port").Int(3000)
```

### Watch

Watch a path for changes. When the file changes the new value will be made available.

```go
w, err := conf.Watch("hosts", "database")
if err != nil {
	// do something
}

// wait for next value
v, err := w.Next()
if err != nil {
	// do something
}

var host Host

v.Scan(&host)
```

### Merge Sources

Multiple sources can be loaded and merged. Merging priority is in reverse order. 

```go
conf := config.NewConfig()


conf.Load(
	// base config from env
	envvar.NewSource(),
	// override env with flags
	flag.NewSource(),
	// override flags with file
	file.NewSource(
		file.WithPath("/tmp/config.json"),
	),
)
```

## FAQ

### How is this different from Viper?

[Viper](https://github.com/spf13/viper) and go-config are solving the same problem. Go-config provides a different interface and is part of the larger micro 
ecosystem of tooling.

