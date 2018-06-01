// Package file is a file source. Expected format is json
package file

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/micro/go-config/source"
)

type file struct {
	path string
	opts source.Options
}

var (
	DefaultPath = "config.json"
)

func (f *file) Read() (*source.ChangeSet, error) {
	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	info, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	// hash the file
	h := md5.New()
	h.Write(b)
	checksum := fmt.Sprintf("%x", h.Sum(nil))

	return &source.ChangeSet{
		Source:    f.String(),
		Timestamp: info.ModTime(),
		Data:      b,
		Checksum:  checksum,
	}, nil
}

func (f *file) String() string {
	return "file"
}

func (f *file) Watch() (source.Watcher, error) {
	if _, err := os.Stat(f.path); err != nil {
		return nil, err
	}
	return newWatcher(f)
}

func NewSource(opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}
	path := DefaultPath
	if options.Context != nil {
		f, ok := options.Context.Value(filePathKey{}).(string)
		if ok {
			path = f
		}
	}
	return &file{opts: options, path: path}
}
