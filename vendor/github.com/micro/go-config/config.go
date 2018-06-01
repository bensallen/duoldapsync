// Package config is an interface for dynamic configuration.
package config

import (
	"context"

	"github.com/micro/go-config/reader"
	"github.com/micro/go-config/source"
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	Close() error
	Bytes() []byte
	Get(path ...string) reader.Value
	Load(source ...source.Source) error
	Watch(path ...string) (Watcher, error)
}

// Watcher is the config watcher
type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

type Options struct {
	Reader reader.Reader
	Source []source.Source

	// for alternative data
	Context context.Context
}

type Option func(o *Options)

// NewConfig returns new config
func NewConfig(opts ...Option) Config {
	return newConfig(opts...)
}
