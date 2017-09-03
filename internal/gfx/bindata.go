// Code generated by go-bindata.
// sources:
// data/shader.frag
// data/shader.vert
// data/testTexture.png
// data/textPlane.ply
// DO NOT EDIT!

package gfx

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _shaderFrag = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\xc1\x6a\xc3\x30\x0c\x86\xcf\xf5\x53\x08\x76\x89\x4b\x08\x49\xda\x9d\x4c\x0f\xa3\xbd\xee\x21\xbc\xd4\xcd\x04\xb6\x35\x54\x27\x64\x8c\xbd\xfb\x88\x63\xbb\x30\x7a\x31\xf2\x6f\xfd\x9f\x7e\xeb\x65\x36\x7c\x47\xf2\x70\x7c\x6d\x61\x20\x36\x42\x58\xfd\x4d\x53\xa8\x2c\x0d\x3a\xac\x2f\x27\xe8\x5a\x09\x93\xc7\x1b\xb1\x83\xbb\x76\x5f\xd6\x70\x7f\x81\x60\x96\x30\xb1\x51\xcf\x0c\xdd\xc3\x70\xb3\xa4\x03\x0c\x64\x89\xdf\x71\x79\x73\x34\xf9\xf0\xd4\xd3\x3f\x3c\xb3\x19\x0e\x91\x7f\x5e\x6d\x4a\x08\xf4\xab\xd6\xaf\xda\x99\x88\xaf\x2a\x29\xc7\x0d\x9c\xaf\x07\xb0\x38\x7e\x06\xf4\xa3\x12\x82\xa6\xb0\xb5\xdc\x58\x8f\x19\x34\x13\x5e\xc1\x69\xf4\xd5\x5a\x49\xf8\x11\xbb\xd8\x13\xc1\x96\x18\x4e\xd1\x53\x95\xe1\x35\x74\x4d\x2b\x61\xbf\xc9\x5d\xd3\x46\x21\x1d\x69\x05\xfd\xa5\x4a\x55\x5d\x12\xca\x86\xa5\x82\x44\x77\x58\xe8\x0e\x97\x2a\x0f\xab\xb7\xf4\xf5\xbf\xed\x48\x25\x76\x25\x73\x0e\x94\x11\x0d\x8f\x1f\xb0\x2f\xff\xac\x0b\xbb\xd1\x52\x89\xdf\xbf\x00\x00\x00\xff\xff\xd2\xd8\xf3\x9b\xd0\x01\x00\x00")

func shaderFragBytes() ([]byte, error) {
	return bindataRead(
		_shaderFrag,
		"shader.frag",
	)
}

func shaderFrag() (*asset, error) {
	bytes, err := shaderFragBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "shader.frag", size: 464, mode: os.FileMode(438), modTime: time.Unix(1498281760, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _shaderVert = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x92\xcd\x6e\xea\x30\x10\x85\xd7\xf1\x53\x1c\xe9\x6e\x92\x5c\x84\x02\x84\xfb\xa3\x28\x2b\xb6\x6d\xd5\x45\xc5\xb6\x72\x13\x43\x5d\x39\x1e\x64\x0c\x0d\xad\xfa\xee\x95\xf3\x03\x01\x4c\x77\xc9\xf8\x7c\xc7\x3e\x33\xf3\x6b\x2f\xcc\x56\x92\x46\x3a\x4f\x50\x90\x11\x8c\x29\x7e\xa0\x9d\x0d\x15\x15\xdc\xba\x93\x1c\x49\x84\x9d\x96\x2b\x32\x15\x2a\x6e\x53\x6c\x0c\xbd\x89\xc2\x1d\x2e\xa5\x78\xbf\xe7\xd6\xc8\x3a\xf3\x70\x93\x0b\xae\xa2\x52\xa8\xdb\xf2\xe9\x85\x5c\x93\xa9\xf8\x51\xef\x01\x66\x27\x60\x2f\x8a\x19\x78\xf5\x22\x85\xb6\x77\x72\xfd\x6a\x17\xa4\xc8\xf8\x6e\x49\x2f\xa0\x52\x9a\x36\x0b\x57\x3f\x83\xf3\xdb\xe0\x52\x14\xd6\x41\x1e\xea\x4f\x04\xa9\x1d\x90\x62\x43\x5b\xe9\xaa\x3e\xf3\xbf\x27\x59\x1b\xdb\x27\xfa\xd7\x8b\xa6\x90\xfa\x49\xd4\x0b\x22\x53\xfa\x84\xff\x4f\x6e\x52\x77\x81\x18\xed\x6c\xcb\xda\x23\xd9\x95\x52\x14\xad\xa6\xfb\x9f\x41\xb9\x4e\x48\xbd\xce\x18\xdb\x93\x2c\x51\x71\xa9\x43\xf7\x15\xe1\x93\x05\x6b\xf5\xfc\xd8\x65\x41\xee\xdd\x06\xc4\xc3\x61\x23\x1e\x64\x67\x41\x7f\x3f\xf2\xb3\x18\x41\xf3\x88\xa6\xd8\x3f\x39\x68\x1e\x67\x0d\xd7\x5b\xd7\x77\x51\x3e\x34\xbd\x41\x7e\xb6\x1b\x88\x9b\x10\x61\x5b\x1b\xd7\x87\x8f\x11\x26\xe3\x24\xca\x58\xb0\x52\xc4\xed\x70\x52\xc8\x51\xf1\x3a\x2c\xc9\x86\x57\xb6\x2d\x79\x35\xd6\x68\x84\xa4\x75\xeb\xbb\x82\xfc\x7a\xd5\xf0\x1b\xa1\x77\x95\x10\x0f\x2d\xa3\x8c\x7d\x7d\x07\x00\x00\xff\xff\x33\xf6\x94\x00\x74\x03\x00\x00")

func shaderVertBytes() ([]byte, error) {
	return bindataRead(
		_shaderVert,
		"shader.vert",
	)
}

func shaderVert() (*asset, error) {
	bytes, err := shaderVertBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "shader.vert", size: 884, mode: os.FileMode(438), modTime: time.Unix(1498281754, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _testtexturePng = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xea\x0c\xf0\x73\xe7\xe5\x92\xe2\x62\x60\x60\xe0\xf5\xf4\x70\x09\x62\x60\x60\x68\x00\x61\x0e\x36\x06\x06\x86\xc3\x76\x89\xa7\x19\x18\x18\x05\x3c\x5d\x1c\x43\x2a\x6e\xbd\xbd\x73\x90\x97\x41\x91\x83\xc1\xd1\xe4\xdb\x7d\xa3\xba\xb7\x0e\x17\x66\x72\xd9\x6c\x7a\x58\x6a\x65\xa4\xd4\x9c\xef\xab\x2c\xc0\xc0\xe4\xc0\xc0\xc1\xc0\xa8\xc0\xc0\xd2\xc0\x80\x97\xb3\x60\xf6\xd5\x87\x13\x16\xbf\x2c\xb8\xcc\x86\x2c\xdc\xba\x66\xba\xdd\xb1\xec\x70\x33\xdc\x3a\x1b\x18\x7a\x1c\x18\x6a\xde\xac\xac\x67\xca\x8f\x53\xf2\x95\x46\x56\x23\xf2\xbc\x68\x05\x2f\x01\x6b\x71\x70\x34\x56\x1f\x15\x36\x9c\x57\xb7\x91\x91\x68\xcd\xcb\x79\x65\x36\x73\xb0\x36\x5f\x7e\xf8\x81\x81\x81\x81\xc1\xd3\xd5\xcf\x65\x9d\x53\x42\x13\x20\x00\x00\xff\xff\x18\xf7\xe3\x78\x49\x01\x00\x00")

func testtexturePngBytes() ([]byte, error) {
	return bindataRead(
		_testtexturePng,
		"testTexture.png",
	)
}

func testtexturePng() (*asset, error) {
	bytes, err := testtexturePngBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "testTexture.png", size: 329, mode: os.FileMode(438), modTime: time.Unix(1504323831, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _textplanePly = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x8f\x41\x4e\xc3\x30\x14\x44\xf7\x3e\xc5\xec\x0a\x12\x8d\xe2\x14\x09\xc2\x12\x2e\xc0\x0d\x2a\xc7\xfe\xa1\x96\x1c\x3b\xb2\x7f\x48\xc2\xe9\x91\xd2\xaa\x94\xb8\x80\xf0\xc6\xa3\xf7\xff\x8c\xc7\xbd\x9b\x45\x1b\x62\xa7\x18\x2a\x69\x6b\x21\x8b\x52\xe8\xd0\x75\xe4\x19\x2f\x91\x14\x93\x41\x33\xe3\xd9\x91\x37\x14\x51\x15\x0f\x8f\xb8\x49\x43\x83\xf2\x16\x5b\x8c\xe3\x58\x34\xc7\x51\x11\xe2\xdb\x1d\x52\x18\xa2\x26\xb4\xd6\xd1\x13\x36\x4c\x13\xbf\x3a\xe5\xe9\xb8\xb4\x11\xe4\x68\x49\x7e\xa7\xc8\x34\xe1\x5e\xf4\x31\xf4\x14\x79\x46\xeb\x82\x62\x4c\x6b\x30\xaf\xc1\xc7\x1a\xf8\xcc\xe3\x33\x93\xcf\x5c\x69\x0d\xf8\xdc\xad\x55\x9a\x50\x7d\xcd\x9d\x4d\x8c\x41\x1f\x54\xc4\x60\xcf\xdd\xf7\xd6\x1b\xab\x29\x09\xf2\x66\x7f\x20\x65\x28\x8a\xb2\x28\x97\x83\x5c\x6c\x73\x24\xbf\x11\x79\x21\x84\xbc\xbe\xf3\x77\x50\x5d\xd7\xf5\x85\x10\xdb\x1f\x5e\xfb\x47\xa5\x53\x52\xee\xff\x45\x5c\x6b\x74\xfa\xda\x0e\x25\x24\xaa\xe5\xde\x41\x8a\xcf\x00\x00\x00\xff\xff\xcc\x75\x23\x0c\x81\x02\x00\x00")

func textplanePlyBytes() ([]byte, error) {
	return bindataRead(
		_textplanePly,
		"textPlane.ply",
	)
}

func textplanePly() (*asset, error) {
	bytes, err := textplanePlyBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "textPlane.ply", size: 641, mode: os.FileMode(438), modTime: time.Unix(1504402834, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"shader.frag": shaderFrag,
	"shader.vert": shaderVert,
	"testTexture.png": testtexturePng,
	"textPlane.ply": textplanePly,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"shader.frag": &bintree{shaderFrag, map[string]*bintree{}},
	"shader.vert": &bintree{shaderVert, map[string]*bintree{}},
	"testTexture.png": &bintree{testtexturePng, map[string]*bintree{}},
	"textPlane.ply": &bintree{textplanePly, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

