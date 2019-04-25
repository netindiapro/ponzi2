// Code generated by "esc -o bindata.go -pkg rect -include .*(ply|png) -modtime 1337 -private data"; DO NOT EDIT.

package rect

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/data/roundedcorner_edges.ply": {
		local:   "data/roundedcorner_edges.ply",
		size:    653,
		modtime: 1337,
		compressed: `
H4sIAAAAAAAC/5yQXW7DIBCE3znFXsCIPzftaSJqJoklG6wNruKevkrUUgu3Dw4S0uhjZ9jdaVjEKfHo
M/lr1/ekpRJdGkfETJzmGBC6xBF8RDjjKu8GDHi8f4AzbtSKidMEzgudhuQz3Wqw1OCzBnHjiRtTXLnm
7uKZGKFGZwZiDd+HGaXr+xzkfkv6Mon+CxqBGI4X+AAWWqrHoSLU/6LUmLb9uUJJY92LtaTkqzsc2rfd
CY2STjujHa3E7ojvz5tVPzszCm2eXwZpockIQ1ZYcl8BAAD//zmz7guNAgAA
`,
	},

	"/data/roundedcorner_faces.ply": {
		local:   "data/roundedcorner_faces.ply",
		size:    1070,
		modtime: 1337,
		compressed: `
H4sIAAAAAAAC/6SRQY+bMBCF7/4Vc9tWKpbHNrbpsf0hEcGTXSRiR8Y0S399RXZFsg5dpSpchsebx/fM
aZjZIaZjm6Edu74H5IJ18XikkOFnojaTh/0MPwYKnhJIbh18Gac9iK9Qwfl85vu3Vzym528wxil1BId+
oO/wlOIUPPkupkBpd2g7Gt/sT4wGunzjF6VMr+DYKcUTpTzDYYhthtdSmEvhdymEu51wtxTutsZSyFdh
6l7aBIl8KT0nolCK+2GitddSFvTVMvRjfvdN/dp71wffdzQyCn73Qq2nxCrBNWqJGm4Gcbk2BrwqaOqm
RhBcSo1OgVhuthoq/HvKVpwVxjmzDKZGdO9xgkuljVIguNPW1s2DaUIIlAuctlIZU8J9wlZtw1mUdjHV
TtXqH+A245qmcba5wXz86D7P+3B21YpV3ZA+FKhRo9Jl3+r/f+6HwksugmQKFGioL88GkCmwgGDYnwAA
AP//PqGt/C4EAAA=
`,
	},

	"/data/squareplane.ply": {
		local:   "data/squareplane.ply",
		size:    625,
		modtime: 1337,
		compressed: `
H4sIAAAAAAAC/6SPQW7rIBiE95xidnlPaiycVGrVZXuQCMM4QSLg/kAT9/RVk8qp7F0DC9DHfAwMYVR9
kqMpMNl6j7bRyqbjkbHgTWgKHboRr4HRUbBpnp7xL9cO+j/WOJ1OTXc9apLsH5BTFUv0PvAFq/xejXAI
JvIaWykGXu7+oBSe8agGSQOljOhDMgXnORjn4HMO4sKJCyn+sqo9GIHQzdFeyDiHXaicXt0bS2xukeBz
+clVP/1q56Pzllkxut2BxlFU2+jLwHra6QnpBbqFvqe6SffZfyq/q3sLjRaby7pFq74CAAD//6eCpu5x
AgAA
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/data": {
		isDir: true,
		local: "data",
	},
}
