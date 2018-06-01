# File Source

The file source reads json config from a file

## File Format

The expected file format is json

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

## New Source

Specify file source with path to file. Path is optional, it will default to `config.json`

```go
fileSource := file.NewSource(
	file.WithPath("/tmp/config.json"),
)
```

## Load Source

Load the source into config

```go
// Create new config
conf := config.NewConfig()

// Load file source
conf.Load(fileSource)
```
