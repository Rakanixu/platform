// Code generated by go-bindata.
// sources:
// data/bindata.go
// data/es_flags.json
// data/es_mapping_files.json
// data/es_settings.json
// DO NOT EDIT!

package data

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

var _dataBindataGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func dataBindataGoBytes() ([]byte, error) {
	return bindataRead(
		_dataBindataGo,
		"data/bindata.go",
	)
}

func dataBindataGo() (*asset, error) {
	bytes, err := dataBindataGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/bindata.go", size: 0, mode: os.FileMode(420), modTime: time.Unix(1469456390, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataEs_flagsJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xbc\x92\xb1\x4e\xc3\x40\x0c\x86\xf7\x3c\x85\x95\xb9\xea\x03\xb0\x15\x78\x02\x56\xc4\x60\x2e\x0e\xb5\x7a\x71\xa2\xb3\xef\xe0\x84\x78\x77\x9c\x74\x69\x8b\x22\xb5\x0c\x1d\xb2\xf8\xfc\xff\x9f\xf3\xeb\x7f\x6d\x00\xbe\xfd\x03\x68\x0f\x54\xdb\x07\x68\x9f\x91\x63\x85\x9e\x23\x81\x06\x14\x61\xf9\x68\x37\xc7\x8d\x8e\x34\x24\x9e\x8c\x47\x59\x36\x59\xf1\x3d\xfa\x3b\x0c\x58\xc1\x1f\xab\x84\x7d\x1a\x85\x95\x80\xa5\xa3\xaf\x13\x69\xc1\x98\xc9\x45\x96\x32\xf9\xe4\x67\xf3\x97\x9b\xa7\xc8\x01\xed\x46\x2c\x8b\x5b\x73\x37\xeb\x30\x85\x3d\x17\xba\x1e\xf9\x34\x8a\x91\xd8\xf5\xc0\x44\x6a\x89\x83\x29\x28\xcd\x34\xe8\xb3\x84\x79\xcd\x4f\xb0\xfa\x0f\xf0\xe2\x72\x37\xec\xee\x3c\xa0\x0b\x1e\x5c\x1a\xf5\x18\x75\xcd\x29\xdb\x38\xa0\x71\x00\x1a\xbc\x2f\x7e\xe1\x34\x26\x5b\x31\x7e\xa1\x40\x8e\x85\x4f\xa2\x83\x57\xeb\x54\xa1\x37\x1c\xdf\x15\x94\x40\x1d\xa0\xff\x76\x75\xb4\xae\xe0\x76\x1e\x4d\x99\xfb\x90\x28\xe2\x31\x25\xf0\x7e\x20\x14\xd6\xec\x89\xe9\x32\xdc\xc2\x23\xf9\xec\x2c\xc9\xed\x6a\x02\xcd\x5b\xf3\x1b\x00\x00\xff\xff\xf6\x94\xd8\xce\x2a\x03\x00\x00")

func dataEs_flagsJsonBytes() ([]byte, error) {
	return bindataRead(
		_dataEs_flagsJson,
		"data/es_flags.json",
	)
}

func dataEs_flagsJson() (*asset, error) {
	bytes, err := dataEs_flagsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/es_flags.json", size: 810, mode: os.FileMode(420), modTime: time.Unix(1468346734, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataEs_mapping_filesJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xec\x59\xcd\x6e\xdb\x30\x0c\xbe\xe7\x29\x02\x1f\x87\x62\xd8\xba\xae\xdb\xfa\x2a\xc3\x20\x28\x16\x9d\x10\x95\x25\x4f\x92\x97\xa4\x6d\xde\x7d\xf2\x5f\x6b\xcb\xb2\xe3\xb4\xf9\xd1\x21\x87\xa6\x89\x4c\x93\x1f\x3f\xd2\x34\x45\x3d\xcf\xe6\xf3\x88\x6d\x05\x4d\x31\x26\x06\xd2\x8c\x53\x03\x3a\x7a\xf8\x6d\xd7\xe7\xf3\xe7\xf2\xb3\x90\xb0\xab\x24\x41\xe0\xcc\x5e\x6b\x56\xed\x7a\x4a\xb3\x0c\xc5\xb2\xbd\x56\x48\xcb\x98\xfc\xa3\x3c\x2f\x14\x19\x95\xc3\x4d\xfb\x62\x22\x55\x4a\x4d\xf4\x50\xe9\x94\x99\x41\x29\x28\x27\x06\x53\x88\x3a\x82\x66\x9b\x41\x2d\x16\xbd\xae\xef\x6e\xda\xc6\x4d\xbc\xb2\x12\x9f\x3f\x45\xee\x2a\xa9\x81\x91\xb6\x92\x9e\x50\x46\x8d\x01\x25\xec\x75\x05\x4b\xd8\x34\x56\x76\xb3\x96\xa5\x37\x0a\x20\xa1\x39\x37\x44\x48\x43\xa8\x45\xbc\x7d\x02\x46\xd6\x68\x56\xa4\xed\xee\x5e\x6e\x50\x30\x6b\xe9\x21\x6a\xab\xe9\xfa\x3d\xca\x5e\xed\xcf\xf3\x6b\xc8\xec\xef\xdd\x87\xe9\xd1\x46\xd9\x5f\x2f\x0b\x29\x39\x50\xf1\xc2\x64\xbe\xe0\xf0\xc2\xa5\x5d\x43\x61\x2c\x37\xea\x50\xf2\xec\xe7\x9f\xe2\x96\x28\x53\x32\x03\x65\xf0\x8d\x9c\x28\x96\x56\xa7\x30\x6f\xc4\x78\x84\x2a\xc1\x15\xc4\x8f\x3a\x4f\x1d\x0a\x3b\xa0\xbb\xdc\x1d\xce\xae\x97\xba\xc6\x70\x95\x95\x5e\xeb\x15\x47\xef\xd1\x5d\x79\x4f\xb4\x91\xca\x82\xf3\x2a\xaf\x03\xe1\x8f\x6b\xa3\xe0\x04\xd8\xca\xc7\xa4\xab\xd2\x1f\x9a\x1a\x48\xde\x8e\xa2\x8b\xa2\x48\x9f\x0e\x86\x71\x1c\x0e\x96\x52\xb8\x14\x1c\xb6\xe0\xc9\x80\xf2\xf2\xde\x2c\x98\x80\x64\xe6\xfb\xde\xe6\x8a\xcb\x98\x16\xb5\xeb\xca\xd7\x34\xbe\x52\x29\x60\x7b\x25\x6b\x1a\x59\x52\x2d\xa9\x40\x7d\x4d\xb0\x03\x38\xb3\xac\xc4\xe0\xfa\x7f\xa5\x6b\x84\x2e\x7d\x4d\xae\xa9\x6c\x19\xd8\xb8\x99\x35\xd2\x86\xd4\x30\x54\x79\x99\x0a\x46\x15\xf3\xbf\xcb\x0d\x3e\x52\x92\x82\xa1\xf6\xcd\x4b\x4f\xa2\xbf\xd7\x23\x5c\x43\x3c\x10\xe2\x5c\xe0\x5f\x17\xdc\x70\x43\xd6\x6d\x77\x2b\x35\x11\x6c\x50\x1b\x4d\xa4\x20\x0c\xf5\x63\xab\xc1\xf5\xa9\x69\x6e\x4a\x50\x69\xdb\x0e\x02\x88\xfe\x0d\x9d\x7d\x8b\xdf\xc1\x46\x0d\x32\xb2\xb8\xbf\xeb\xab\x70\x38\x1d\xe7\x72\xdc\x04\xa7\xc7\x00\xda\x4f\xf7\xa1\xfe\x9f\xc6\x31\x68\x3d\xd4\x24\x3b\x7b\xba\xc9\xfd\xb7\x02\x7b\xe3\x71\x75\x32\x54\x76\x3b\xb4\x3a\xfb\x3e\xa5\xb6\x4b\x74\xc6\xf1\x80\x17\xdf\x97\xcb\x3c\x81\x4e\x2d\xf8\x1a\x06\x8a\x40\xc8\x08\x84\x8d\xdb\x30\x60\x7c\x0b\x03\xc6\x5d\x18\x30\xbe\x87\x01\xe3\x3e\x0c\x18\x3f\xc2\x80\xf1\x33\x0c\x18\xbf\x82\x80\x11\x46\xd9\xb8\x0d\xa3\x96\x87\x51\xbc\xc2\xa8\x5d\x61\x94\xae\x30\x2a\x57\x18\x85\x2b\x8c\xba\x75\xa1\xb2\x35\xf3\x7d\xef\x74\xd4\x32\xae\x0f\x45\xce\xdc\xca\xc3\xc6\x80\xd0\xfd\x71\xdf\xe9\x2d\x27\xc8\x41\xd0\xa1\xb3\x84\x3d\xe3\x87\xe6\x66\x52\x81\xe9\x88\xf5\xcf\x0b\xcb\x65\x45\xd7\x81\x06\xff\xd5\x99\xce\x46\x7a\x1f\x1b\xc7\x0a\x43\xce\xf9\x45\x36\x91\x4b\x1c\xd8\x0c\xf7\xe6\x3d\xd3\xf4\xa5\x68\x77\xf8\x97\x78\x84\x52\xc9\xd0\xe6\xdc\x71\xb7\xf6\x7a\x45\x15\x5c\x24\x2e\x1a\x9f\x06\x38\x7c\x67\x60\x74\xba\x20\xd4\x58\xc0\x8b\xdc\xb8\x53\x81\xd1\xb1\xa0\x6f\x62\x32\x8a\x67\x1a\x0d\xe3\xc8\x1d\xf4\xa5\x70\x59\x24\x19\xb0\x41\x37\xf6\xa3\x3a\xc8\x5e\x39\xf0\xaa\xa6\x50\x01\x78\x5f\xa1\xb1\x9e\x93\x78\x45\xc5\x72\x64\x8e\x7a\x5e\x48\x6b\x85\xee\xb1\xf1\x45\xc0\xc8\x24\xe1\x28\x46\x90\xf4\xa6\xb8\xa5\x96\xbd\xaf\x85\xfc\x18\x05\x72\xd6\xfe\xdf\x8c\x42\xed\x13\x97\xa2\x2e\xde\xfa\x7a\xea\x34\x94\x2c\x95\xcc\x33\xf7\xe9\x9d\x78\x56\x50\x8c\x09\x8b\xe9\xf4\x82\x2a\xff\x79\x41\x6d\x23\xd7\xa0\x4e\x67\x82\x9f\x4e\x33\x01\xa5\xa4\x3a\x7b\xa9\xa6\x9c\xcb\xf5\x89\xbc\x62\x20\xdc\x63\xf3\x0f\xa9\xf6\x26\xa2\x7b\xf6\xe1\x3f\xb0\x98\x15\x7f\xbb\xd9\xff\x00\x00\x00\xff\xff\xbb\x97\x44\xdd\x34\x26\x00\x00")

func dataEs_mapping_filesJsonBytes() ([]byte, error) {
	return bindataRead(
		_dataEs_mapping_filesJson,
		"data/es_mapping_files.json",
	)
}

func dataEs_mapping_filesJson() (*asset, error) {
	bytes, err := dataEs_mapping_filesJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/es_mapping_files.json", size: 9780, mode: os.FileMode(420), modTime: time.Unix(1468506274, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataEs_settingsJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x94\x51\xc1\x4e\xc3\x30\x0c\xbd\xe7\x2b\xaa\x5c\xe1\x80\xd8\x8d\x1f\xe0\xc2\x1f\xac\x25\x0a\xad\x37\x59\xa4\x4e\xd5\x04\xb1\x31\xfa\xef\x38\x40\x57\x27\x68\x53\xb9\xb4\xca\x7b\xb6\xdf\x7b\xf6\x49\x55\x95\x0e\x10\x23\xd2\x3e\xe8\x87\x2a\xbd\x19\xb1\x64\xdd\x31\xe0\x82\xcc\xd8\x07\x8c\x02\x4b\xbd\x83\xc3\x68\x3c\x99\x17\x9b\x33\xcc\xc5\xe3\x00\x8c\xe9\xf6\x2d\x44\xdf\xeb\xdb\x8c\xf3\xaf\x40\xf8\x33\x2e\x1f\x72\xae\x9a\x96\x06\xbd\x43\x07\x64\x7b\x30\x48\x1d\x1c\xfe\x23\xc4\x9d\xf1\x5b\x65\x2b\x50\xc6\x9d\x7f\x87\xb1\xb5\x01\xb2\xf2\x94\x33\xb4\x88\x3b\xef\x3a\x5e\x49\xc9\xd1\x7e\xb4\xbd\xd9\x98\xfb\x3b\x2d\x88\xe6\x72\xb2\xd9\xb7\x48\xa5\x8a\x74\x8b\x43\xb1\x56\x21\x74\x21\x2b\x3d\x72\x45\x1e\xb5\x47\x32\xa9\x2f\xd1\x9b\x82\xb2\x87\x33\x25\xcd\xff\x75\x23\xed\xaf\xbc\xf3\x60\x23\x07\xa0\x34\x7b\xfb\xd9\x14\x77\xfe\xb5\x3b\xd7\x5c\xbd\xee\xb5\xc1\xcf\x75\x3d\x9c\x9e\xa6\xba\xee\x9a\x9b\xb5\x12\x4a\xfe\xd3\x77\x52\x93\x52\x5f\x01\x00\x00\xff\xff\x80\x3c\xd6\x0a\xf4\x02\x00\x00")

func dataEs_settingsJsonBytes() ([]byte, error) {
	return bindataRead(
		_dataEs_settingsJson,
		"data/es_settings.json",
	)
}

func dataEs_settingsJson() (*asset, error) {
	bytes, err := dataEs_settingsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/es_settings.json", size: 756, mode: os.FileMode(420), modTime: time.Unix(1468346734, 0)}
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
	"data/bindata.go": dataBindataGo,
	"data/es_flags.json": dataEs_flagsJson,
	"data/es_mapping_files.json": dataEs_mapping_filesJson,
	"data/es_settings.json": dataEs_settingsJson,
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
	"data": &bintree{nil, map[string]*bintree{
		"bindata.go": &bintree{dataBindataGo, map[string]*bintree{}},
		"es_flags.json": &bintree{dataEs_flagsJson, map[string]*bintree{}},
		"es_mapping_files.json": &bintree{dataEs_mapping_filesJson, map[string]*bintree{}},
		"es_settings.json": &bintree{dataEs_settingsJson, map[string]*bintree{}},
	}},
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

