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

var _shaderFrag = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8f\xc1\x6a\xc3\x30\x10\x44\xcf\xd9\xaf\x58\xe8\x45\x02\x51\x8a\x9b\xb4\x07\xe1\x43\x49\xae\xfd\x08\xa1\xc8\x45\x20\x69\xcb\x46\x32\x2e\xa5\xff\x5e\x24\xc7\x36\x04\xdf\x96\xd1\xbc\x19\xcd\xd3\xe8\xf8\xe6\x29\xe1\xf1\xf4\x82\x96\xd8\x01\x04\xf3\x43\x25\x8b\x40\xd6\xe4\xfa\xd2\xe3\x49\x62\x49\x7e\x20\x8e\x78\x33\xf1\x3b\x38\xee\x2e\x98\xdd\x94\x0b\x3b\xbd\xe3\x7f\xdb\xfc\x43\x20\x93\xd1\x52\x20\xfe\xf4\xd3\x47\xa4\x92\xf2\x1e\xf2\xbe\x21\xa3\xb3\xaf\x2d\xfd\x5c\x29\x0d\xe0\x53\xd5\xba\xaa\x9d\x89\xf8\xaa\xef\xca\x71\xce\xd5\x00\x54\xf2\x2c\x0c\x6c\xbe\x16\x6c\x24\x7f\xc5\x68\x7c\x12\xf5\x92\xf8\x0b\x87\xe6\x69\x31\x81\x18\xfb\xc6\x88\xb5\x4a\x2d\x9b\xba\x8b\xb8\x5f\x6a\x2d\x95\xcf\x2c\x35\x1c\xd6\x02\xec\x31\xfa\x49\x2c\x61\x6a\xfe\x8b\x7a\x98\x2a\x35\xfc\xfd\x07\x00\x00\xff\xff\x5d\x75\xd7\xca\x63\x01\x00\x00")

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

	info := bindataFileInfo{name: "shader.frag", size: 355, mode: os.FileMode(438), modTime: time.Unix(1505281516, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _shaderVert = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\x41\x4b\xc5\x30\x10\x84\xcf\xcd\xaf\x18\xf0\xf2\xfa\x4e\xcf\x1a\x4f\xa5\xa7\x9e\x05\x0f\xe2\x55\x42\x1a\x65\x25\xc9\x96\x34\xad\x15\xf1\xbf\x4b\x6a\xaa\x55\x72\x0b\xb3\xdf\xec\xce\xe4\x6a\x31\x61\x22\xf6\x90\xb7\x17\x68\x0e\x46\x08\xab\xde\x79\x8e\x27\xcb\x5a\xc5\x34\xe9\x70\xa9\x31\x7b\x7a\xe6\xe0\xe0\x54\x94\x18\x03\xbf\x1a\x9d\x86\x8f\x64\xde\xee\x54\x0c\xb4\xb6\x05\xdf\xf5\x3f\x9f\xe3\xc1\xd8\x1d\x2f\xf0\x4d\x0d\xf2\x58\x8c\x96\x18\x79\xa2\xa4\x96\xd6\xde\xec\x58\x03\xf2\x0f\x66\xed\x99\xc3\x50\x02\xe5\xef\x3e\xf2\x3d\x5b\x0e\xad\x10\x3c\xc7\x6f\x6f\xfc\x71\x66\x49\x42\x67\x66\x61\x1a\xe0\x14\xf9\x53\x7a\xd5\xf8\x10\xd5\x8b\x7d\xba\xcf\x99\xd0\x15\x7f\x00\xe7\x63\x41\x9c\x0f\x1d\xaa\xfd\x16\xba\x3f\x91\xab\xed\xe0\x26\xe6\x78\x9f\x5f\x01\x00\x00\xff\xff\x92\x05\x12\x42\x91\x01\x00\x00")

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

	info := bindataFileInfo{name: "shader.vert", size: 401, mode: os.FileMode(438), modTime: time.Unix(1505279225, 0)}
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

