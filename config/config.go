package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/BurntSushi/toml"
	"github.com/gopherd/core/encoding"
	"gopkg.in/yaml.v3"
)

const HeaderChecksum = "X-Checksum"

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeYAML ContentType = "application/yaml"
	ContentTypeTOML ContentType = "application/toml"
)

var (
	// ErrDuplicatedKey is the error that the key is duplicated.
	ErrDuplicatedKey = errors.New("duplicated key")
	// ErrNotFound is the error that the key is not found.
	ErrNotFound = errors.New("not found")
	// ErrOperationNotAllowed is the error that the operation is not allowed.
	ErrOperationNotAllowed = errors.New("operation not allowed")
)

// Row is the type of a row in the table.
type Row interface {
	GetID() string
}

// Table is the interface that wraps the basic operations of a table.
type Table interface {
	// Parse parses the data into the table.
	Parse(data []byte, decoder encoding.Decoder) error
	// Lookup looks up the row with the given id.
	Lookup(id string) (row Row, err error)
	// Scan scans the table rows with the given offset and limit.
	Scan(offset, limit int, desc bool) (rows []Row, total int, err error)
	// Insert inserts a new row into the table.
	Insert(rowContent string) (row Row, err error)
	// Update updates the row with the given id.
	Update(id string, content string) error
	// Delete deletes the row with the given id.
	Delete(id string) (deleted bool, err error)
}

// ContentType is the content type of the data or empty.
//
// For example:
//
// - application/json
// - application/yaml; charset=utf-8
// - application/toml; charset=utf-8
type ContentType string

// Parse parses the content type and returns the extension, encoder, decoder.
func (c ContentType) Parse() (ext string, enc encoding.Encoder, dec encoding.Decoder, err error) {
	if c == "" {
		return "json", json.Marshal, json.Unmarshal, nil
	}
	if i := strings.IndexByte(string(c), ';'); i >= 0 {
		c = ContentType(c[:i])
	}
	switch c {
	case ContentTypeJSON:
		return "json", json.Marshal, json.Unmarshal, nil
	case ContentTypeYAML:
		return "yaml", yaml.Marshal, yaml.Unmarshal, nil
	case ContentTypeTOML:
		return "toml", toml.Marshal, toml.Unmarshal, nil
	default:
		return "", nil, nil, fmt.Errorf("unsupported content type")
	}
}

// Scopes represents the scope of the data.
type Scopes []string

// Compact returns a compacted scopes. If the scopes contains "*", the result is "*".
// Otherwise, the result is a sorted and compacted scopes.
func (s Scopes) Compact() Scopes {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == "*" {
			return s[i : i+1]
		}
		if s[i] == "" {
			s = append(s[:i], s[i+1:]...)
		}
	}
	slices.Sort(s)
	return slices.Compact(s)
}

// String returns the string representation of the scopes.
func (s Scopes) String() string {
	return strings.Join(s, ",")
}

// Has returns true if the scopes contains the given scope.
func (s Scopes) Has(scope string) bool {
	if s.Any() {
		return true
	}
	_, ok := slices.BinarySearch(s, scope)
	return ok
}

// Any returns true if the scopes contains any scope.
func (s Scopes) Any() bool {
	return len(s) == 1 && s[0] == "*"
}

// Options represents the options for loading data.
type Options struct {
	// Source is the source of the data or empty. If the Fetch function is not nil, the source is ignored.
	//
	// Supported sources are:
	//
	// - http:// or https://: the data is fetched from the given URL.
	// - file://: the data is fetched from the given directory.
	//
	// For example:
	//
	// - http://example.com/cfg
	// - file:///etc/cfg
	Source string

	// ContentType is the content type of the data or empty (default is "application/json").
	ContentType ContentType

	// Scopes is the scopes to load.
	Scopes Scopes

	// Fetch is the function to fetch data or nil.
	Fetch func(contentType ContentType, scopes Scopes) ([]byte, error)

	// Namer is the function to name the scope. If the Namer is nil, the scope + "." + ext is used.
	Namer func(scope, ext string) string
}

// Hub is the interface that wraps the basic operations of a configuration hub.
type Hub interface {
	// Parse parses the data into the hub.
	Parse(data []byte, decoder encoding.Decoder) error
}

// Config is the configuration.
type Config[H Hub] struct {
	new      func() H
	hub      atomic.Pointer[H]
	checksum string
}

// New creates a new configuration.
func New[H Hub](new func() H) *Config[H] {
	return &Config[H]{new: new}
}

// Latest returns the latest configuration. If the configuration is not loaded, it will panic.
func (c *Config[H]) Latest() H {
	return *c.hub.Load()
}

func (c *Config[H]) parse(data []byte, dec encoding.Decoder) error {
	hub := c.new()
	if err := hub.Parse(data, dec); err != nil {
		return err
	}
	c.hub.Store(&hub)
	return nil
}

// Load loads the data by the given options.
func (c *Config[H]) Load(options Options) (bool, error) {
	options.Scopes = options.Scopes.Compact()
	if len(options.Scopes) == 0 {
		return false, nil
	}
	for _, scope := range options.Scopes {
		if scope == "*" {
			return false, fmt.Errorf("scope * should be resolved before loading")
		}
	}
	if options.Fetch != nil {
		_, _, dec, err := options.ContentType.Parse()
		if err != nil {
			return false, err
		}
		data, err := options.Fetch(options.ContentType, options.Scopes)
		if err != nil {
			return false, err
		}
		return true, c.parse(data, dec)
	}
	if strings.HasPrefix(options.Source, "http://") || strings.HasPrefix(options.Source, "https://") {
		return c.loadHTTP(options)
	}
	if !strings.HasPrefix(options.Source, "file://") {
		dir, err := filepath.Abs(options.Source)
		if err != nil {
			return false, err
		}
		options.Source = "file://" + dir + "?ext=json"
	}
	return true, c.loadDir(options)
}

// loadDir loads the data from the directory.
func (c *Config[H]) loadDir(options Options) error {
	u, err := url.Parse(options.Source)
	if err != nil {
		return err
	}
	dir := u.Path
	ext, _, dec, err := options.ContentType.Parse()
	if err != nil {
		return err
	}
	data := make(map[string]json.RawMessage)
	scopes := []string(options.Scopes)
	for _, scope := range scopes {
		var name string
		if options.Namer != nil {
			name = options.Namer(scope, ext)
		} else {
			name = scope + "." + ext
		}
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		data[scope] = content
	}
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.parse(content, dec)
}

// loadHTTP loads the data from the HTTP.
func (c *Config[H]) loadHTTP(options Options) (bool, error) {
	_, _, dec, err := options.ContentType.Parse()
	if err != nil {
		return false, err
	}
	checksum := c.checksum
	newChecksum, data, err := fetch(checksum, options.Source, string(options.ContentType), options.Scopes)
	if err != nil {
		return false, err
	}
	if newChecksum == checksum {
		return false, nil
	}
	if err := c.parse(data, dec); err != nil {
		return false, err
	}
	c.checksum = newChecksum
	return true, nil
}

func fetch(checksum, url, contentType string, scopes Scopes) (newChecksum string, body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(scopes.String()))
	if err != nil {
		return
	}
	req.Header.Set(HeaderChecksum, checksum)
	if contentType == "" {
		contentType = string(ContentTypeJSON)
	}
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	newChecksum = res.Header.Get(HeaderChecksum)
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	return
}
