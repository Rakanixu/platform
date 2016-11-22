// Code generated by go-bindata.
// sources:
// data/bindata.go
// data/es_desktop_flags.json
// data/es_flags.json
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

var _dataBindataGo = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd4\x9b\x6b\x6f\x23\x57\x72\xf7\x5f\x8b\x9f\xa2\x57\x80\xf7\xa1\x76\xb5\x14\x9b\x7d\x21\x5b\x0f\xfc\x62\x7d\x0b\xbc\x88\xbd\x81\xd7\x41\x5e\x78\x8c\x41\x5f\x4e\x6b\x3a\xa6\x48\x85\xa4\xe6\xe2\x81\xbe\x7b\xea\x57\x55\xcd\x8b\x44\x69\x66\x94\x49\x90\x0c\xc0\x11\xfb\xf4\x39\x75\xea\xfa\xaf\xaa\xd3\xcd\x8b\x8b\xe8\xeb\x65\x13\xa2\xab\xb0\x08\xab\x72\x13\x9a\xa8\x7a\x17\x5d\x2d\xff\x52\x75\x8b\xa6\xdc\x94\xa3\x81\x4c\x58\x2f\x6f\x57\x75\x58\x5f\xf2\x9d\xc1\x8b\xfe\xe6\xd5\x72\x3b\x14\xd6\x2f\xdb\x79\x79\xb5\x1e\xfd\xfb\x7a\xb9\xd8\x1f\xbd\x2e\x6f\x6e\xba\xc5\xd5\xcb\xb6\x9b\x87\x0f\xdc\x7d\xb9\x08\x6f\x1e\xcc\x58\x87\xcd\x46\x66\xec\x96\x7e\xf3\xf7\xe8\xc7\xbf\xff\x1c\x7d\xfb\xcd\xf7\x3f\xff\x61\x30\xb8\x29\xeb\xdf\xca\xab\xa0\xd3\x07\x83\xee\xfa\x66\xb9\xda\x44\xc3\xc1\xc9\x69\xf5\x6e\x13\xd6\xa7\xf2\xa5\x5e\x5e\xdf\xac\xc2\x7a\x7d\x71\xf5\x7b\x77\xc3\x40\x7b\xbd\xe1\x4f\xb7\xb4\xff\x2f\xba\xe5\xed\xa6\x9b\x73\xb1\xd4\x05\x37\xe5\xe6\xd5\x05\x0c\xf1\x85\x81\xf5\x66\x05\x07\x7c\xdd\x74\xd7\xe1\x74\x70\x36\x18\xb4\xb7\x8b\x3a\x72\x3d\xfc\x14\xca\x66\xc8\x97\xe8\x97\x5f\xd9\xf6\x3c\x5a\x94\xd7\x21\xb2\x65\x67\xd1\xb0\x1f\x0d\xab\xd5\x72\x75\x16\xbd\x1f\x9c\x5c\xfd\xae\x57\xd1\xe5\x97\x11\x5c\x8d\x7e\x0c\x6f\x20\x12\x56\x43\x65\x9b\xeb\xaf\x6e\xdb\x56\xae\x21\x7b\x76\x36\x38\xe9\x5a\x5d\xf0\x87\x2f\xa3\x45\x37\x87\xc4\xc9\x2a\x6c\x6e\x57\x0b\x2e\xcf\x23\x11\x69\xf4\x2d\xd4\xdb\xe1\x29\x84\xa2\x2f\xfe\xe3\x32\xfa\xe2\xf5\xa9\x71\xa2\x7b\x09\x8d\xbb\xc1\xe0\xe4\x75\xb9\x8a\xaa\xdb\x36\xb2\x7d\x6c\x93\xc1\xc9\x4b\x63\xe7\xcb\xa8\x5b\x8e\xbe\x5e\xde\xbc\x1b\xfe\x51\xe6\x9c\x0b\x6f\xb2\xaa\x9e\x7f\xdb\x73\x3a\xfa\x7a\xbe\x5c\x87\xa1\x88\xff\x99\xf8\x81\x8c\xd1\x7f\x84\x90\x4c\x34\xbe\x7d\x50\xd8\x1a\x7d\x05\xeb\xc3\xb3\x73\x66\x0c\xe4\xde\xe6\xdd\x4d\x88\xca\xb5\x38\x0a\x2a\xbf\xad\x37\x50\x51\xf9\xdc\x1e\xb2\xcd\xa2\x5d\x46\xd1\x72\x3d\xfa\x4e\xcc\xfa\xbd\x5c\x6c\xd7\xb9\x09\xfb\xf1\x3d\x0a\x6a\x43\xf9\x67\x66\x1c\x9c\xac\xbb\xdf\xf5\xba\x5b\x6c\xf2\x74\x70\x72\x4d\xd4\x44\x5b\xa2\x3f\xc8\xa5\x0e\xfe\x2c\x1e\x12\xe1\x26\x23\xbe\xb1\x8f\xba\xca\xb0\xed\xee\xef\x75\x16\xfd\x28\x5b\x0c\xcf\x7c\x07\xf6\x74\x29\xdb\x6e\xc4\xee\xb2\xf8\xf1\xb5\xff\x10\x76\x64\xad\x72\x73\xb8\x14\x46\x9f\x5c\x0a\xaf\xb2\x74\x8f\xf3\x43\x02\x88\xf6\x21\x02\x08\x27\x34\xb6\x82\x3e\xa0\xe0\xd2\x3f\x4e\xe4\xfb\xf5\x37\xdd\x4a\x48\x54\xcb\xe5\x7c\x7f\x75\x39\x5f\x7f\x40\xf2\x77\x6b\x13\x3c\xac\xda\xb2\x0e\xef\xef\xf6\x56\xbb\x4b\xe0\xe5\x2f\x59\xf5\x95\x2d\xfe\xa7\xa5\xf8\xb6\x39\xc3\xf0\xf4\xc5\xdb\xb8\x7d\xf1\x76\x56\xbd\x78\x3b\x9e\xc9\x67\xec\x9f\xe2\xc5\xdb\x3c\xc8\xb8\x8f\xb5\x32\x67\x1c\xef\xee\x73\xad\x63\xe3\xa7\x3e\xa7\x3d\x34\x1c\x6c\xee\x1e\x7b\x0c\x09\x7a\xbf\xde\x43\x12\x09\x81\x43\xde\xcf\x65\xe4\xf4\x1e\xf0\x9e\xca\xe0\xd9\xd6\xbb\x0e\xa6\xb3\xd1\x9f\x34\x1e\xf6\x37\xd2\x80\xd8\xa2\xce\x31\xf6\x3e\x14\xd4\xdb\x58\xd4\x68\x12\x22\xf7\x2c\xf3\x1e\x9f\xbd\x8c\x1e\x72\x1a\xe1\x91\x97\xd1\xf8\x3c\xc2\xb3\x2e\xf7\x1d\x6f\x98\x26\xf9\x99\x8e\xe3\x2f\x97\xe6\x4f\xff\xba\xe8\xde\x0e\xe3\x74\x3a\xce\xb3\x34\x9d\xe6\xe7\xd1\xf8\x4c\xa0\xa2\x64\xcb\x3f\xaa\x5c\xef\x55\x98\xcb\xc8\x65\x82\x9f\x4b\xfd\xff\x6e\xab\xd0\xf2\xfc\x81\x2f\x7c\xeb\x29\xea\x6f\x92\x48\x9e\xe5\x0e\x55\xfd\xe2\x6d\x31\x91\xbf\xe2\x16\xa9\xdc\xab\x13\xf9\xcb\x5c\x19\x9f\xe5\x32\x67\xfa\xe2\x6d\xc2\xf7\x4c\xe6\xc9\xa7\x12\x1a\xa1\x94\xfb\x32\xaf\x92\x79\xb1\x8c\x4d\xa1\x27\x34\x32\x99\x5f\xa7\xb2\x87\x8c\x4f\x84\xd6\x58\x3e\x15\xf7\x65\xfe\x54\xe8\x97\xec\x23\xeb\x82\xec\x1b\x64\xce\x2c\xb5\xb5\x53\xd9\xa3\x90\x3d\xa6\xac\x2d\x8c\xef\xc9\xc4\xd6\xc2\x47\xdc\xc8\x77\xb9\x6e\x65\x6e\x5b\x1b\xdf\x05\xae\x0b\x2d\x99\x3b\x95\xef\x79\x63\x32\x55\xb2\x67\xdb\x18\x7f\x39\x7c\xc9\xbd\x4c\xe8\x36\xe8\x62\x6a\x63\xac\x2d\x84\x9f\x3c\x71\xb9\x64\xcd\x84\xef\x84\x86\xc8\x10\xc3\x87\x7c\x6f\x0b\x9b\x9f\xc8\xba\x5a\x3e\x33\xa1\x3b\x93\xb9\x89\xdc\x9f\xa4\xb6\x6e\x26\xfc\xa4\x72\x2f\x93\xb9\x49\x6e\x7f\x5b\x59\x9b\xc8\xdc\x56\xe4\x4e\x2a\x93\x21\x13\x3a\x75\x6c\x76\x29\x2b\x97\x5d\xe6\xc4\xe5\x4e\xb7\x33\xe1\xbf\x94\xbf\xa5\xf0\x52\xca\xbc\x38\x31\x7d\xb0\xae\x60\xbd\x7c\x4f\x99\x2b\xfb\x24\x13\xe3\x2f\xc0\xeb\xd4\x74\x51\xa0\x47\xa1\x5d\xc2\xef\xcc\x42\x3d\xc8\x1e\xa9\xcc\x9f\xd4\xa6\xd7\x0c\x7d\x25\x26\x13\xba\x4b\xc6\xb6\x3f\xfc\xc6\x32\x56\x09\x3f\x71\x30\xda\xc8\x39\x2b\x4d\x57\x0d\xba\x97\x79\xb5\xcc\x4f\xd9\x13\xbe\xb1\x95\xd0\x9a\xc8\xdf\x09\xfc\x35\xb6\x26\xcc\xcc\xce\x2a\x63\x61\xe3\x69\x6b\xfe\x82\x4e\xb0\x49\x2b\xdf\x83\xc8\x30\x9d\x38\x2f\xc2\x5f\x08\xe6\x6b\x25\xfe\x57\x19\x1f\x59\x6d\x3a\x62\xef\x18\x5f\xc9\x8c\x1e\x7b\x62\xd3\x64\x66\xf3\xa7\x0e\x6b\xcc\xc5\xa7\x27\xb2\xc7\x54\xd6\x87\xd8\xed\x95\x9b\xec\x49\x6c\xe3\x65\x6c\x7e\x8e\x3d\x91\x1f\xfe\xf0\x4f\xe4\xc9\x5c\x37\xf8\x0d\xf3\xe0\x2f\x16\x9a\x4d\x6b\xfe\x5e\xa4\xb6\x16\xbf\x2c\x3d\x8e\xa6\xb9\xf1\x56\xb1\x26\x98\x7f\x33\x86\x0c\xa9\xcc\x2b\xf0\x0b\xa1\x59\x8c\x4d\x57\xc8\x0e\x2f\x93\xcc\x74\x01\x4f\x93\xb1\xf9\x5e\x93\x9b\x3e\xd0\x6f\x28\xcc\xd6\x8d\x5c\xd7\x13\xf3\x53\xb5\x7b\x6d\x36\xc5\x27\xb0\x31\xf6\x20\x06\xd1\x0d\x72\xe2\x0f\x71\xf5\x10\xea\xdb\xdc\xf8\xc7\x9e\x35\x7b\x7a\x2c\x1f\x81\xfa\x7d\x6c\x79\x1e\xda\xef\x53\xd8\x01\xfe\x41\x59\xfd\x10\xf3\xf7\x17\x7d\x2c\xec\x1f\x61\xf5\xb3\x22\xff\x3d\x96\x1d\xfc\x67\xf1\xa7\xc3\xff\x38\x29\x66\xe9\xe7\x84\xff\x83\x6e\xe3\xd9\x69\x00\x77\xca\x3c\x54\xb9\xd7\x38\x34\x00\x5d\x40\x6a\x98\x5a\xc8\xe1\x62\x4a\x13\x28\x9d\x98\x1b\x01\x19\xa5\xaf\x21\xbc\x71\x29\xd2\xc8\xc4\x5d\x3d\xc6\xe5\x80\x97\xb1\xc1\xdc\xc4\xc3\x8d\xb4\x03\x64\x95\xa9\xc1\x77\x13\x0c\x12\x5b\x52\x09\x90\x2e\xf4\xea\xca\xa0\x3f\x24\x16\x5a\xb8\x7e\x43\xd8\xc4\xc6\x43\x99\x5b\x98\xa4\xa4\xb1\xc4\x78\x4b\xe4\xd3\x4c\x2c\xec\x53\x87\xb6\x1a\x19\x73\x0b\x0b\xe4\xd6\x74\x91\x59\x68\x6a\x18\xe6\x06\xfd\xac\x03\xd2\xb3\xc4\xc2\x22\x29\x2d\x7d\xb5\xbe\xa6\x99\x9a\x2e\x34\xf5\x4c\x8d\x57\xa0\x8e\xb0\x06\xde\x80\x76\x52\x43\x1a\x1b\x04\xd4\xd0\x0a\x26\x0b\xe9\x91\x50\x85\x5e\x5a\xda\x3a\xf4\x5b\x67\x26\x63\xd5\xd8\x9c\x74\x66\x72\x17\x9e\xee\x48\x07\x89\x43\x1a\xf0\x10\x17\xa6\x47\x4d\xd5\x8d\xc1\x0a\xb6\x20\x65\x66\x99\xd1\xd6\x71\x99\x9b\x8d\x0d\x42\x08\x79\x20\x08\xbe\x7a\xb8\x9e\x4c\x4d\xe6\xcc\x75\x82\xac\x33\xa1\x91\xe7\x96\x46\x74\x6e\x63\xf6\x83\x37\x85\xb5\xd6\xfd\x03\x1b\xe7\xb6\x0f\x70\x0b\xff\xf1\xc4\x52\x2b\xfa\x47\xef\xc1\xe1\xa5\x2c\x2c\xb5\x92\x46\xaa\xdc\xed\x1e\x8c\x37\x52\x13\x69\x75\xec\x50\x88\x1f\x02\xa3\x6a\x9b\xca\xfc\x0b\x3b\xa2\x17\xe0\x54\x75\x9d\x5b\x9a\x04\x96\xc7\xa5\xf9\x99\xc2\x74\x61\x63\xa4\x27\x74\x98\x92\x46\x67\xc6\x0b\x3c\xe2\x5b\xf8\x2f\x7b\x02\xb1\xa4\x76\xe8\xe2\x7f\xd8\x05\x7e\xd5\x56\xce\x1b\xfe\x8f\xdd\xd0\x1f\xfe\x06\xc4\x96\x7e\xad\xa9\x1d\xbf\x2a\x6d\x4d\x70\x1f\xc3\x2f\xd4\x6f\x67\x9e\xd2\xf0\xc5\xc6\x7c\x35\xb8\xef\xea\xfe\x85\xf9\x38\x36\x25\xd5\x61\x13\xf8\x67\x5f\xca\x07\xa0\x3e\x4f\x3d\xf6\xc6\x96\x22\x18\x43\xee\xac\x31\xd9\x28\x51\x48\x01\x9a\x3e\x4b\x5f\x3b\xf1\x7d\xa7\x66\x0b\xd2\x0c\xfa\xc2\xde\xd8\x99\x54\xa2\xa5\x91\xa7\x1d\x74\x4f\xe9\xc5\x7d\xe6\xa9\x3e\xa6\x2e\x7b\x6b\xe5\x1c\x25\x1a\xa9\x5c\xe3\x3a\xb6\x14\xc8\xf8\xac\x35\x9f\x21\x96\x49\x85\x63\x2f\x77\x48\x89\x55\x65\x29\x8b\x71\xd2\x0c\xbe\x53\x78\x79\x81\xdf\x4c\x1d\x0b\xf8\xa8\xfe\x2a\xa3\x4b\x29\xa5\x3a\x8b\x4d\x4e\xca\x07\x2d\x1d\x52\xf3\x9f\xda\x7d\x14\x3d\xab\x6f\xb5\x26\x0f\x69\x7e\xe6\x98\x41\xbc\xa0\xb3\x71\x9f\xe6\x4a\x2b\xdf\xb0\x33\xf6\x6b\x1c\xc3\x8a\xd2\xca\x05\x7c\x05\xec\xa1\x14\xd0\x12\xa0\xb2\xbf\xa9\x97\x8a\x94\xb1\x94\xc2\x94\x18\x94\x01\x5c\xab\xae\x88\xd9\x99\xc7\x64\x6d\xbe\xad\x31\x5f\x58\x29\x87\x2d\xf8\x40\x1b\x3b\x11\xcf\x8c\x63\x57\xc5\xa2\xca\xee\xf1\x9d\x32\xab\x8f\x79\x8d\xb7\xcc\x74\x02\xef\x60\x4e\xe2\xbe\xc0\x9e\xf8\x0a\x25\xf6\xb4\xb2\x98\x41\x66\xe4\x00\xbf\xf5\x7e\x63\x31\x47\x89\x82\xcf\x6b\x8c\x8e\x1d\xcf\xbc\x34\x2b\x6b\xb3\x61\xf0\x92\x8c\x12\xba\x74\x5a\x6a\xeb\xb0\x5b\x8b\xed\xa0\x97\xd7\x86\x3d\x8a\xf9\x33\xc3\x93\x24\x73\xcc\xc4\xa7\x53\xf3\x9b\xcc\xcb\x3e\xee\xe3\x13\xf8\x3b\xb6\xc5\x66\x55\x6b\x72\x8e\xc1\x5d\x97\x07\x3c\xc0\xa6\xc4\x30\xe5\x89\x96\x8b\xd5\xae\x24\x63\x0d\x3a\xa0\xf4\xa1\x1c\x4c\x1d\x0f\x63\x2f\x21\x19\x43\x0e\xe2\x89\x12\x16\xbf\xc3\xbf\x98\xa7\x71\xde\x58\x8b\x31\xa9\x8c\xdf\x38\x36\x59\x88\x17\x4a\x24\x95\x31\x35\x1c\x01\x8f\x29\x35\x69\x0d\xf0\x15\xfc\x9d\x6b\x95\xdf\xf3\x66\xea\xb1\x96\x7b\x09\x8a\x3f\x80\x2b\x60\x68\x1b\x4c\x36\xf2\x1d\xfa\xc3\x97\x28\x1b\x6b\x8f\x7b\x95\x7b\xe2\x78\x3f\x33\xd9\xb0\x33\x7b\xc1\x33\xf6\x21\x66\x88\xdf\xa9\x97\x7d\xf8\x71\xe9\x2d\x08\xf9\xba\xf1\x98\xa3\x35\x53\x9d\x05\xf3\x53\xe2\x8a\x5c\xa6\x2d\x4b\x63\xf6\x20\x17\x11\x93\xf8\x8a\xc6\x5e\x61\xed\x51\x35\x35\x9e\x19\x27\xb7\xe2\xdb\xe0\x2e\xbc\xa9\xfd\x72\x5b\x47\xdc\x11\xe3\xf8\x2a\xf6\xc2\x47\xf0\x73\xec\x43\xbb\xc6\x9e\xe0\x15\x7e\xad\xba\x6f\x2d\xce\xb1\xb1\xb6\x8d\xde\x56\xe2\x93\xe4\x50\xfc\x74\xec\xfe\x8e\x7f\x6a\x9b\xe2\xf8\x02\x16\xf5\x75\x87\xfa\x77\x6d\xd8\x88\xdf\xb2\x0f\x78\x40\x0b\x90\x78\x9b\x07\xa6\x6a\x8b\x93\x18\x6f\xda\x2a\xce\x2c\xff\x80\x95\x33\xc7\xd2\x24\xd9\xdd\x43\xdf\x33\xbf\x17\xef\x61\x31\xb6\xc0\x36\xa5\xe7\x17\x68\x67\xb1\xe5\x0d\xc5\xdd\x99\xed\x0f\xe6\x16\x5e\xe2\xc3\x07\x18\x44\xbb\x00\x6d\xda\x08\x7c\x8c\xfa\x01\x5d\xa3\x4f\xec\x83\xad\x13\xf7\x7f\x7c\x02\x7c\x43\x77\x1a\x2f\x13\xa3\x35\xf6\x9c\x09\x56\x6a\xa9\xef\x35\x14\x72\xc3\x6b\x33\x36\xbb\x61\x27\x74\xac\x39\xa9\xb1\x1c\x8e\x0e\xc9\xd5\x8d\xe7\xc1\xbe\xf5\xc2\x07\xd0\x8d\xe6\x9e\xdc\xf2\x3a\xf1\x86\x4d\x58\x3f\xeb\x5b\x3a\xcf\x41\xd8\x8c\xbc\x00\x1f\xc4\x5b\xdf\x0e\x53\x1f\x36\x99\xc5\x33\xfc\xa1\x33\x6a\x18\x30\x94\xf8\xa9\x3d\x0f\x13\x6f\xb5\xb7\x8e\x5a\x77\xe5\x86\xd9\xf8\x8a\xb6\x4f\x33\x8f\x2f\xe7\x47\x7d\x34\xb1\x3a\x45\xb1\xc7\x8f\x02\xf0\x13\x7c\xa8\xcf\x3b\xf0\x8c\xdf\x22\x1f\xeb\x34\xbe\x63\xf3\xe1\xd6\x79\x61\x3f\x30\x06\x2c\x98\x38\x7e\xe1\x6b\x8a\xdf\xb9\xc5\x15\xfe\x44\xab\xac\xbe\x95\x9a\x8d\x88\x79\xec\xa0\xed\xf9\xd8\xe3\xb5\x36\x7f\xc8\xfc\x48\x80\x5c\xc4\xf1\x48\xe9\x39\x91\x5a\x0b\xbb\xe1\xc3\xd0\x53\x9f\x76\x5d\xd7\x1e\xbf\xbd\x4c\xc4\x34\xb8\x3c\xf1\xba\x0d\x1c\x4e\xbd\xa5\x84\xf7\xd2\x6b\x4c\xd6\x91\x07\xa9\x01\xc8\xcb\x65\x69\xb5\x2d\x38\x83\x8f\xa3\x0b\xf0\xad\x74\x5e\xa8\x9f\xa8\x77\xa7\x5e\xdf\x50\x43\x4e\xbd\x66\x20\x8f\x4d\xbd\x2e\x21\x87\x29\xa6\x7b\x3d\xc6\x5a\x62\x22\xf5\x56\x9e\x7d\xe1\x8d\xfc\x49\x0d\x80\x6d\xa9\x65\x95\xff\xdc\xf4\x50\xfa\x91\x83\xd6\x80\xad\xd1\x6a\x3d\x9e\xc6\xbe\x2f\x38\x0b\x0e\xe1\x9b\x7a\x24\x32\xb1\x9a\x28\xf6\x63\x09\x70\x1c\x9b\x90\xef\x5a\xaf\xa9\xa9\x59\xc8\x6b\xc4\x93\xe6\xc7\xcc\xea\x65\xe2\xb5\xf5\x7a\x49\x71\xd4\xeb\x34\x74\x01\xed\xb1\x1f\x11\x69\x1d\x96\x5a\x3c\x13\x87\xda\x86\x7b\x8d\x0b\xe6\x93\x2f\xf0\x63\xec\x06\x06\x61\xab\x99\xd7\xa8\x5c\x2b\x2e\xc4\x7e\xe4\xe2\x38\x07\x1e\x82\xf5\xc8\x4e\x5d\x44\x4c\x10\x23\xf4\x18\x60\x81\xd6\xd9\x63\xf3\x11\x72\x04\x31\x4d\xec\xb2\x06\x7b\xd2\x47\x94\x8e\xab\x5c\xeb\xf1\x83\x1f\x7d\x68\x6e\xcf\x2c\xd7\x6a\x7d\x54\x3c\x3c\x7d\xed\xdb\x74\xee\x17\x5e\xcb\x36\xcd\xee\x08\xe3\x78\x9b\xfe\xa0\x07\x7c\x76\xbb\xfe\x80\xd2\x41\xdb\xfe\xf0\xb9\xd7\xd1\xf6\xfd\x01\x91\x4f\x68\xe3\x1f\x13\xe5\x73\xb7\xf3\xc7\x44\xf1\xb6\xbe\x98\xce\xfe\xb7\xf5\xf5\x3c\x45\x7c\x76\x6f\x8f\x6f\x93\xcb\x32\xef\x61\xfa\x63\x3f\xea\xa7\xd6\xf1\x98\x9a\x92\xd8\xe3\x18\x8a\x5c\x4a\xce\x23\xff\x93\xa3\x4b\xaf\xd3\xc0\x3b\x30\x1c\xac\xd4\xe3\xb1\x99\xc5\x42\xea\xfd\x27\x79\x27\xf7\x5a\x9e\x1a\x85\x1c\x4d\x1c\x34\x7e\x74\xaa\xb1\x53\x59\xde\xa2\xae\x00\xdb\xb5\x5e\xf2\x1e\x8f\xfd\xa8\xc9\xc9\x33\xb1\xd7\x9e\x99\xf7\x3b\x1c\xcb\xc5\x7e\xac\xab\xb8\x19\xac\xd6\x01\x67\xfa\x33\x03\xf0\x0b\x5c\x9f\x79\x4d\x00\x3e\xea\x11\xa3\x9f\x01\x30\x4f\x8f\x7a\x2b\xcb\x55\xe8\xa3\x71\xfc\x6d\xfd\xf8\x1a\x8c\x00\x9b\xa9\x81\xc0\x49\x3d\x82\x4c\x0c\xe3\x88\xe5\xac\xde\xc9\x01\xbf\xe4\x38\x72\x12\x38\xd8\xf6\x75\x76\x69\x76\x18\x7b\x9e\x87\x2f\x3d\x2f\x69\xec\xb8\x16\x9a\xcc\x07\xa7\xa9\x1d\xa8\x77\xf5\x78\x37\x35\x1c\xd7\x7c\xef\x36\x44\x97\xd0\x0a\x7e\x4d\xde\x01\x07\x6a\xc7\x4e\x3d\x96\x25\x2f\x66\x86\x8f\xe8\x4f\xfb\x66\x3f\x96\x45\x0f\xad\xf7\xe4\xe4\x55\x3d\xce\x2d\x2d\x1f\x63\x6f\xc5\xde\x89\xe1\x36\xb9\x62\xec\xb9\x9b\xbd\x66\x9e\xcf\xfb\xbc\xdf\x7a\x5d\xc8\xf1\x31\x39\x8d\x1c\xa5\x7d\x6a\xe1\x79\x81\x5c\xe7\xb4\x2b\xef\x77\xf0\x31\x64\x02\x83\x27\x7e\x1c\xaa\xf5\x70\x61\x7c\x20\x3f\x35\x4c\xd1\x1f\xe9\x4e\x5d\x8f\x99\xd5\x16\xd4\x94\xf4\x2a\xd8\x48\x8f\x5b\xbd\xfe\xe2\x31\x04\x74\x5a\xaf\x3f\xf0\x09\x6c\x88\xef\x36\xee\x23\xb9\x3f\x7a\x20\x6f\xb6\x7e\x94\x5c\x7b\x9d\x85\xef\xe6\xde\xaf\xb3\x2f\x79\x4f\xfb\xd8\x99\xd5\x49\xec\x8d\xcd\xf0\x21\x6a\x6f\xb0\x9a\xde\xb4\xf6\x1a\x43\x7b\x8f\xc6\xfa\x1e\xe2\x70\xe2\x47\xdc\x7a\xce\x34\xb5\x9c\x31\xf6\xbc\xaf\x47\xcc\x85\xed\xd9\xf8\x79\x10\x3d\x83\xd6\x3b\x13\x8b\x15\x8e\x6c\xd1\x09\x7e\x45\xad\x45\x1d\xa0\x8f\x52\x72\xab\xff\x2b\x8f\x09\x62\xaa\xf6\xda\x53\x65\x48\xad\x76\x0b\x5e\xb3\x69\x4e\x1f\xdb\x1e\x5a\x1f\xb7\x8e\x09\xfe\xd8\x86\xb5\xf4\x3e\xcc\x29\xfc\xc8\x9a\xdc\xc5\x77\xee\x4f\x92\x23\x4f\x07\x73\xab\xdb\xc8\x87\xfa\xe8\xa7\x7a\xea\xc8\xf8\x28\x6e\x7d\x9e\x7c\xd4\x53\x7b\x3c\x27\x6d\xdf\xb6\xf8\x70\x5e\xea\x89\x3d\x37\x37\xdd\x13\xed\xbf\x35\x3f\xed\x89\xe5\x39\x6a\x3a\x2d\x9e\xf1\xe4\x31\x2e\x8a\xf1\x67\x4c\x51\xfd\x6b\x2c\xcf\xce\x4c\x78\x94\x56\xf7\xd5\xee\xe1\x63\x9f\x99\x40\x10\xa2\x5b\x33\x86\x47\x71\xf0\x4e\x93\x07\x4c\x64\x94\xd8\x1f\x24\xd6\x7e\xa2\x03\x22\xd2\x69\x27\xde\x31\x54\xfe\x40\x8e\x2e\x04\x14\xd1\x8a\xb0\x72\x74\xf3\x93\x48\xaa\x31\xed\x60\x82\x65\x4a\x32\x60\xe1\xa7\x18\xfa\xa0\xa7\xb2\x71\xd0\xb2\x7f\x48\xa5\x0f\x67\x3c\x73\x12\xa5\x89\x23\x06\xf2\xd5\x8e\x82\xac\x99\xf9\xc9\xd4\xd8\x3b\x25\x3d\x69\x6f\x4d\xbe\xb4\xf0\x53\xd3\xc6\xa2\x8b\x28\xcb\xbd\xa3\x06\xc5\xc8\xbc\x20\x0e\x28\xaa\x0f\x4c\xfd\xc1\x0e\x99\x93\xac\x57\x78\x57\x50\xc4\x3b\xa4\xa7\xcb\x21\x2b\xb5\x7e\xb2\x93\x38\x12\x50\xb9\xd3\x7d\x70\xca\xa0\x55\x6e\xb0\x13\x5e\xba\x13\x22\x9a\xaa\x74\xec\xdd\xb1\x9e\x94\x65\xbb\x6e\x2e\xf3\xd3\x10\xd0\x2a\xf8\xa9\x53\xee\x9d\xbd\xa2\x8f\x9f\x24\xe8\x29\x6f\xee\x7c\x64\xc6\x63\xe9\x68\x9c\xfa\x49\x70\xe5\x3e\x31\xf5\x87\xc6\xa5\x9f\x80\x51\xb1\x93\x09\x38\x49\x42\x6e\x7c\x22\xf8\x09\x14\x59\x62\xdc\x3f\x40\x6e\x0d\x2d\x1b\xef\xf8\xd4\x7e\x89\x77\xd3\x63\xd7\x67\xe2\x0f\x39\xe3\x9d\x6d\xb1\x7b\xe9\x27\x45\x74\x98\x9a\xb5\x5b\xd3\x0b\xb6\x24\x93\xe9\xe9\x64\x6e\xa7\x02\xf8\x2e\x88\x9b\x79\x97\xd5\xf4\x0f\x7d\x33\xaf\x78\x62\x43\x74\x32\x0c\x19\x75\xe2\x6b\x82\x3f\x74\x65\xbf\xd4\x7d\x17\xb4\xd4\xce\xdc\x4f\x94\x52\xef\xba\x93\x6a\xd7\x79\xb5\x7e\x5a\x02\x3f\xfa\x70\x73\x66\x9d\x4c\xe2\x4f\x22\x90\x85\x27\x1c\xd8\x87\xae\x57\x9f\x78\x34\x5e\xbd\xf9\x83\x75\x8d\xa5\xcc\xbb\x1c\xaf\xec\x90\x8d\xee\x81\xb1\xc2\xf9\xed\x9f\x58\x04\xcf\x88\xe8\x4f\x4f\x70\x82\x3f\xb4\x6f\xfc\x61\x64\x30\xd9\x5a\xcf\xfc\xc8\xcf\xc9\x64\xea\x15\x99\x76\x72\x7e\x1a\x37\xf1\x8e\x07\x1a\x85\x57\x22\x64\xcd\xd2\x4f\x7b\x52\x3f\x05\x60\x5d\xd2\x3f\xb9\x98\x98\x4c\xf8\x52\xeb\xa7\xaf\x7c\xd7\xd3\xd1\xda\x75\x77\x2f\x13\x91\x9d\xa9\x14\x88\x7d\xed\x90\x67\x4f\x65\xa2\x7d\x78\x7a\x76\x02\xda\x27\x72\x90\x77\x0e\xde\xe1\x3b\x9a\x6e\xf6\x97\x7e\x42\x96\x39\xc2\xf6\xe7\x4e\x2e\xf7\x78\xef\xfb\x9e\x24\x7f\x4e\x4e\xc9\xa7\xd3\xff\x62\x4e\xb9\xb8\x88\xfe\xaa\xaf\xba\xcd\x97\x65\xb3\x8e\xca\x45\x13\xd9\x9c\x75\xb4\x79\xd5\xbf\x06\xd7\x2e\x57\x7a\x75\xd5\xbd\x0e\x0b\x7d\xe5\x4e\xdf\xe3\xfc\x7e\xb3\x9d\x5b\x2e\x4c\xb5\x91\x28\x6b\xb7\xae\x5e\xde\xce\x9b\x68\xb1\xdc\x44\x55\x10\x22\xb7\x42\x7c\xb9\x62\xe5\xc1\x0d\x76\x0e\xcd\xc8\xec\xa7\xcc\x0c\x3f\xf4\xbe\x63\x5d\x2e\x16\xcb\x45\x57\x97\x73\x5e\x79\x43\x74\x7f\xa7\x72\xf4\x53\xb8\x99\x97\x75\x18\xda\x7b\x81\xa7\x2f\x5e\x88\x8a\x4f\x2f\xe4\xbf\xbf\xc4\x66\xca\xf6\x3c\x5a\xfe\xc6\x8a\x97\x6e\xa8\x5f\x0e\x89\xfd\xfa\xff\xb9\x8f\x95\xcb\xad\x77\xb4\xb8\xc1\x11\x3f\x78\xf4\x1d\x45\x53\xe9\x17\xeb\x48\x68\xff\x3f\xb4\x54\xea\x7b\xb0\x2a\xc0\x91\x17\x17\x79\x73\x71\x6b\x9a\x91\xdb\x0d\x03\x9d\xdc\x0d\x3e\xbc\x07\x6a\x54\xe5\x3a\xd5\x33\xb7\xeb\x0f\xb7\xeb\x8d\x4d\xea\xd6\xd1\xbc\xfb\x2d\xb8\xa5\xab\xdb\x4d\x74\x53\x8a\xc0\xeb\xe8\xcd\x2b\xb1\xa7\x8d\xbe\x51\x93\xf4\x4c\xb8\x39\x7b\x33\xaf\xbb\xeb\x9b\x79\xd7\x76\x61\x1d\xad\xcb\x36\x88\x37\x75\x9b\xae\x9c\x77\xbf\x97\x9b\x4e\xea\x8f\x65\x1b\x5d\xcd\x97\x55\x39\x8f\xa4\x44\xe9\xca\x8a\xb6\xde\xcc\xb9\xe5\xe1\xd0\xa4\x66\x51\x74\xb8\x53\xf2\x6e\xda\xd1\x98\x53\x8e\x87\xa7\xea\x58\x97\x3e\xf9\x34\xfa\xb3\xbd\x1a\xfb\xe7\xe8\xf4\x4c\x62\x4c\xfe\xca\x2a\xd3\xd0\xf0\xec\xec\xe0\x4d\xcf\x72\xdf\xdb\xf5\xd5\xcc\xa7\x3c\x5e\x83\xf9\x7f\xda\xed\xe1\xea\x9e\xeb\xef\xbd\x62\xfa\x7f\xce\xff\x55\xc9\xcf\x8c\x81\x4e\x05\xfe\xa8\x10\xe8\xb7\x79\x2c\x0c\x74\x12\x92\xad\x0f\xec\xbc\xd0\x91\xe5\x9e\xdd\xd6\xfb\xb6\xd0\x05\x43\x3c\x75\xf7\x36\xad\x2d\x11\x7d\x5c\x97\xbf\x85\x61\x7f\xe7\x9c\x17\x11\xe7\x61\x31\xec\x15\x8a\xe3\xe1\x3a\x0b\x37\xcd\xaa\x5c\x5c\x85\xad\xba\x55\x71\x46\xe9\xcb\x48\x1a\x8d\xb0\x68\xd4\x58\xeb\x9e\xeb\x7d\x81\x19\x77\x31\xb6\xeb\x25\x98\xcb\x68\x43\x90\x9d\x47\xaf\x96\xf3\x06\xe6\x42\x59\xbf\x72\xe7\xf3\x17\xef\x97\x2b\xc9\x1b\x90\x6f\xa2\xcd\x32\xea\x36\x6b\xf7\x5f\xed\x21\x7a\x52\x48\x72\xf3\x8b\x89\xf1\x2b\xb2\x3f\x4c\x95\xc2\xed\x83\x17\x30\x2f\xa3\x7b\xff\xee\xbf\x5d\x7a\xec\xc5\x9d\xcb\x07\x2b\xee\xbd\xa1\xf4\xd4\xf1\xe0\xe5\xfe\xa2\x63\xe7\xa4\x1f\xea\xdd\x2e\x9f\xec\x23\xf7\x09\x1c\xe6\xe7\xcb\x07\x2c\x1f\x16\x25\x7b\x2e\xf6\x4d\xb7\x3a\x70\x30\x76\x70\x2f\xab\xc2\x7c\xf9\x46\xac\x56\x87\xd5\xa6\xec\xec\xe7\x08\xdd\x2a\xd4\x62\x25\x09\x88\xeb\x2a\x34\x82\x04\x82\x39\xbb\x75\x0f\x7f\x34\xf1\x9d\x78\x54\x78\x5b\x0a\x0e\x07\xd0\xe6\xdd\xf2\x36\x5a\xdd\x2e\xf6\x66\x45\x02\xc4\x2a\xc4\x68\x34\x52\x54\xd3\xc1\x7a\xb9\x60\x4b\x65\x09\x32\xed\x72\x2e\xbc\xe0\x33\xaf\x3a\xf1\x93\x55\xfd\xea\x9d\xfe\x0e\xa3\x97\xf0\xa2\xbf\x88\x64\xea\x72\xb4\x79\xbb\xd9\x0d\x74\xd7\x57\x7b\xb7\x23\x09\xd3\x9b\xc5\xd5\xfe\x40\xd5\x0f\x6c\xb6\x79\x85\x77\xb2\x55\xb7\xa7\x67\x87\x29\xa6\x0f\x9f\xf7\xa7\xbe\x0f\x40\x25\x3b\x9c\xde\xed\x6b\xd4\xd6\x5e\x30\xfe\xe8\x7a\x65\x83\xd5\xba\xfd\xfd\xf5\x3d\xf5\x33\xd5\xc9\x6e\x5c\xe0\x22\xbc\xed\xd6\x9b\xfb\x74\x7b\x48\x3f\xa4\xc2\xac\x6e\x3e\x7f\xb8\xb9\x8a\x76\xb7\x8f\x1c\xcc\xbf\x5f\xbf\xf4\x48\xb1\x43\xf0\x05\x6f\xc9\x3b\x0a\x6f\x56\x21\x28\x36\x83\x22\x0a\x02\x80\xec\x58\x91\xe2\xb9\x48\x7f\xc2\x8f\x4d\xfe\x59\xe4\xdb\x5f\xf2\x0f\x49\xe2\x9b\xe1\x21\x49\x5d\xc4\x02\x20\xeb\xe5\x79\x74\xb3\x83\xac\x2d\x09\xc5\x7a\xe5\x58\xb0\x5f\xfe\x8c\xbe\x7e\xd5\xcd\x9b\x55\x58\xfc\x72\xf3\x2b\xb7\x84\x75\xbb\xbb\x97\x1a\x9e\x53\xb7\x9c\x58\x16\xb8\xeb\x7f\xca\xa1\x7b\x7d\x87\x66\x3f\xe2\x77\x21\x4f\xd1\x05\x52\x5f\x3f\x8e\xdc\x07\x32\xf5\xf0\x5d\x73\xfd\xe3\x01\x86\x1f\xcc\x33\x66\x5e\xef\x50\x7c\xf5\xfa\x7c\xb7\xe8\x00\xc7\xb9\xb3\xff\xe3\x12\xb7\xf9\xde\x8f\x43\x54\x48\x0d\xba\xa3\x28\x3c\x38\xd9\xee\xba\x87\xd8\x7f\xea\x7d\xa7\x3f\x1b\xea\xe9\x4a\x27\xe0\x5f\xdf\xab\x9a\x8e\xac\xe9\x71\x5d\x10\xee\x63\xe6\x9e\x9c\x3e\x85\xff\x5b\x0a\x87\x89\xe0\x28\xad\xbb\x3b\x6d\xe5\x9e\x4a\x0e\x07\xe4\x0e\xb2\xc4\x07\x28\x3e\x9a\x39\xee\x53\x7c\x98\x42\x3e\x85\xf2\x7e\x5a\x79\x92\xf2\x36\xbf\x7c\x80\xfa\xa3\x39\xe7\x3e\xf5\x83\xe4\xf3\x38\x51\xfe\xbb\xb3\xd4\xf4\x53\x58\x4b\x96\x09\x16\x1c\x2b\xbb\xd0\xba\xd5\xaa\x05\x89\x92\xb0\x5f\xe4\x6e\xd3\x92\x01\xda\xfe\xea\xa1\xdc\xbb\xf7\x4b\x34\xab\x7d\xc5\x3d\x60\xef\x63\x0b\x79\x8f\x09\xef\x9b\x4f\xba\xbe\xb4\xdd\x2e\xdd\x96\xc1\x1f\xb3\xdc\x7e\x66\x26\x65\xf2\x0f\xbf\x09\x7f\x7f\x9d\xcf\x87\xaa\xfa\x7f\x11\xe8\x32\x86\xfb\x1f\xde\x8d\x7a\x5c\x3e\x93\x8e\x7a\xbf\xc7\x1e\x4f\xb3\xec\xc9\x5f\xc4\x1d\x6e\x65\xbf\xee\x1b\xfd\xdb\xaa\xdb\x04\x68\xdc\xdf\x4f\xb7\x38\x8f\x4c\x23\x08\x37\xb2\x9f\x47\x7d\x92\x30\x5f\xbf\xa2\xd5\x5f\x3f\x42\xbb\xa7\x6a\xbf\x99\xba\x7f\xfd\x31\x1b\x1d\xfe\xba\xe9\x9e\x9b\xac\x3f\xc5\x4f\x64\x6e\x7d\xbb\x5a\xcb\xe0\xfc\x88\xcf\xac\x9f\x72\x9a\xda\x21\xed\xd0\xfa\x5b\x33\x0d\x4e\xa8\x7a\x44\xfe\x27\xe4\x39\xee\x9f\x86\xbd\xfc\xb4\xb3\x5b\x0d\xfa\xbc\xa6\xbb\xed\xa0\xbc\xde\x47\x71\xd3\xfb\x11\xc6\xb7\xce\xf3\xb7\x65\xb7\xf0\x4c\xab\x2b\xcf\x9e\x6e\x88\x54\xcf\x9e\xc7\x0e\x75\xad\x2a\x3a\x62\xd6\xad\x7a\x76\xed\xc6\x73\x53\xff\xf6\x67\x73\xfb\xbc\x7b\x8e\xda\x16\x2e\xb2\xf1\xdd\xf9\x47\x14\x06\x52\x4c\xea\x47\x78\xff\xcf\x00\x00\x00\xff\xff\xa3\x58\x0b\xfc\xdb\x3b\x00\x00")

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

	info := bindataFileInfo{name: "data/bindata.go", size: 15323, mode: os.FileMode(436), modTime: time.Unix(1470908584, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataEs_desktop_flagsJson = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8a\xe6\x52\x50\xa8\x06\x62\x05\x05\xa5\xec\xd4\x4a\x25\x2b\x05\xa5\xe0\x8c\xfc\x72\x85\x8c\xcc\x94\x94\xd4\x3c\x85\xb4\xcc\x9c\xd4\x62\x25\x1d\x88\x7c\x4a\x6a\x71\x72\x51\x66\x41\x49\x66\x7e\x1e\x56\x75\x0a\x99\x79\x0a\xc5\xa9\x89\x45\xc9\x19\x0a\x45\xa9\xc5\xa5\x39\x25\x70\x9d\x65\x89\x39\xa5\xa9\x40\x3d\x69\x89\x39\xc5\xa9\x40\xa1\x5a\xae\x58\x2e\x40\x00\x00\x00\xff\xff\xf2\x41\xa3\x2c\x79\x00\x00\x00")

func dataEs_desktop_flagsJsonBytes() ([]byte, error) {
	return bindataRead(
		_dataEs_desktop_flagsJson,
		"data/es_desktop_flags.json",
	)
}

func dataEs_desktop_flagsJson() (*asset, error) {
	bytes, err := dataEs_desktop_flagsJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/es_desktop_flags.json", size: 121, mode: os.FileMode(436), modTime: time.Unix(1472209221, 0)}
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

	info := bindataFileInfo{name: "data/es_flags.json", size: 810, mode: os.FileMode(436), modTime: time.Unix(1470908584, 0)}
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
	"data/bindata.go":            dataBindataGo,
	"data/es_desktop_flags.json": dataEs_desktop_flagsJson,
	"data/es_flags.json":         dataEs_flagsJson,
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
		"bindata.go":            &bintree{dataBindataGo, map[string]*bintree{}},
		"es_desktop_flags.json": &bintree{dataEs_desktop_flagsJson, map[string]*bintree{}},
		"es_flags.json":         &bintree{dataEs_flagsJson, map[string]*bintree{}},
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