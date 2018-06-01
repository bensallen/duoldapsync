package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/imdario/mergo"
	"github.com/micro/go-config/reader"
	"github.com/micro/go-config/source"
	hash "github.com/mitchellh/hashstructure"
)

type jsonReader struct{}

func (j *jsonReader) Parse(changes ...*source.ChangeSet) (*source.ChangeSet, error) {
	var merged map[string]interface{}

	for _, m := range changes {
		if m == nil {
			continue
		}

		if len(m.Data) == 0 {
			m.Data = []byte(`{}`)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(m.Data, &data); err != nil {
			return nil, err
		}
		if err := mergo.Map(&merged, data, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	b, err := json.Marshal(merged)
	if err != nil {
		return nil, err
	}

	h, err := hash.Hash(merged, nil)
	if err != nil {
		return nil, err
	}

	return &source.ChangeSet{
		Timestamp: time.Now(),
		Data:      b,
		Checksum:  fmt.Sprintf("%x", h),
		Source:    "json",
	}, nil
}

func (j *jsonReader) Values(ch *source.ChangeSet) (reader.Values, error) {
	if ch == nil {
		return nil, errors.New("changeset is nil")
	}
	return newValues(ch)
}

func (j *jsonReader) String() string {
	return "json"
}

// NewReader creates a json reader
func NewReader() reader.Reader {
	return &jsonReader{}
}
