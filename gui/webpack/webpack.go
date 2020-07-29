package webpack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/tlentz/d2modmaker/gui/generated"
)

// Webpack is a webpack integration
type Webpack struct {
	Manifest Manifest
}

// Manifest reflects the structure of asset-manifest.json
type Manifest struct {
	Files       map[string]string `json:"files"`
	Entrypoints Entrypoints
}

type Entrypoints []string

func readProdFile(filename string) ([]byte, error) {
	f, err := generated.ReactAssets.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// It's a good but not certain bet that FileInfo will tell us exactly how much to
	// read, so let's try it but be prepared for the answer to be wrong.
	var n int64 = bytes.MinRead

	if fi, err := f.Stat(); err == nil {
		// As initial capacity for readAll, use Size + a little extra in case Size
		// is zero, and to avoid another allocation after Read has filled the
		// buffer. The readAll call will read into its allocated internal buffer
		// cheaply. If the size was wrong, we'll either waste some space off the end
		// or reallocate as needed, but in the overwhelmingly common case we'll get
		// it just right.
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}
	return readAll(f, n)
}

func readDevFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// It's a good but not certain bet that FileInfo will tell us exactly how much to
	// read, so let's try it but be prepared for the answer to be wrong.
	var n int64 = bytes.MinRead

	if fi, err := f.Stat(); err == nil {
		// As initial capacity for readAll, use Size + a little extra in case Size
		// is zero, and to avoid another allocation after Read has filled the
		// buffer. The readAll call will read into its allocated internal buffer
		// cheaply. If the size was wrong, we'll either waste some space off the end
		// or reallocate as needed, but in the overwhelmingly common case we'll get
		// it just right.
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}
	return readAll(f, n)
}

// read file based on environment
func ReadFile(isProduction bool, filename string) ([]byte, error) {
	if isProduction {
		return readProdFile(filename)
	} else {
		return readDevFile(filename)
	}
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	var buf bytes.Buffer
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	if int64(int(capacity)) == capacity {
		buf.Grow(int(capacity))
	}
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

func New(isProduction bool, buildPath string) (*Webpack, error) {
	webpack := &Webpack{}
	assetsManifestPath := path.Join(buildPath, "asset-manifest.json")

	if _, err := os.Stat(assetsManifestPath); os.IsNotExist(err) {
		return webpack, nil
	}

	content, err := ReadFile(isProduction, assetsManifestPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file %s: %w", assetsManifestPath, err)
	}

	if err = json.Unmarshal(content, &webpack.Manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest file %s: %w", assetsManifestPath, err)
	}

	return webpack, nil
}

func (e Entrypoints) Scripts() Entrypoints {
	var scripts Entrypoints

	for _, f := range e {
		if strings.HasSuffix(f, ".js") {
			scripts = append(scripts, f)
		}
	}

	return scripts
}

func (e Entrypoints) Styles() Entrypoints {
	var styles Entrypoints

	for _, f := range e {
		if strings.HasSuffix(f, ".css") {
			styles = append(styles, f)
		}
	}

	return styles
}
