// Code generated by go-bindata. DO NOT EDIT.
// sources:
// 1_setup.up.sql (856B)
// 2_timestamp.up.sql (53B)
// 3_found.up.sql (86B)
// 4_fleet.up.sql (158B)
// doc.go (74B)

package migrations

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
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

var __1_setupUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x91\x3d\x6f\xc2\x30\x10\x86\x67\xfc\x2b\x6e\x4c\x24\x96\x76\x65\x0a\xe0\x52\xab\x60\xaa\xd8\x48\x30\x1a\x6c\xa5\x96\xf2\x41\x73\x76\xd5\xfe\xfb\x0a\xb7\x8a\x42\x12\x88\xda\xf9\x7c\xcf\xfb\xfa\xb9\x45\x4a\x13\x49\x41\x26\xf3\x35\x05\xf6\x04\x7c\x2b\x81\xee\x99\x90\x02\xf0\xab\x3c\xc9\xea\x6c\x4f\xc2\x29\xe7\x11\x22\x32\x39\xe5\x1e\x9d\xa9\x99\x06\xc6\x25\x5d\xd1\x34\xbc\xe7\xbb\xf5\x7a\x4a\x26\x67\x7f\x44\x7f\x0c\x1b\x20\xe9\x5e\xb6\x67\xb9\x42\x27\x2e\x3c\x5b\x18\x74\xaa\x38\xc3\x9c\xad\x18\xbf\x7a\xf3\x9a\xb2\x4d\x92\x1e\xe0\x85\x1e\x20\x6a\x92\xa6\xd0\xe2\xc6\x24\x9e\x11\x72\xa7\x73\x61\x11\x6d\x99\x6d\x0c\xa2\xca\x4c\xe8\x5c\xfb\x92\xe9\x5e\xa1\x7f\xff\xa4\xf8\x41\x3f\x2b\x7c\xeb\xcf\x30\xbb\xf7\x41\x74\x55\x6d\xca\x4a\x9b\xa1\xc5\x5f\xc9\xdd\x49\xd8\xd1\x89\xeb\xe1\xa0\x23\xac\xd5\x6b\x0a\x4d\xd2\x95\x30\xc6\x97\x74\xdf\x11\x66\xf5\xe7\x06\xb3\x07\xd8\xf2\xae\xbb\xa8\x89\x5e\x52\xb1\x88\x67\xe3\x94\xc7\x21\x4a\xd0\x3f\x72\xb5\x90\xc4\x2b\x6d\x76\xa5\xfa\x50\x36\x57\xc7\xdc\xdc\x3e\xdd\x6d\x8d\xb5\x79\xf7\x06\xdd\xe5\x06\xa3\xbe\x02\xfb\x8f\xa6\x84\xab\x83\xa9\xa1\xbe\x51\x2b\x3c\x9e\x91\xef\x00\x00\x00\xff\xff\x8d\xc8\x7a\x90\x58\x03\x00\x00")

func _1_setupUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_setupUpSql,
		"1_setup.up.sql",
	)
}

func _1_setupUpSql() (*asset, error) {
	bytes, err := _1_setupUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_setup.up.sql", size: 856, mode: os.FileMode(0664), modTime: time.Unix(1716213156, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe, 0x35, 0xf8, 0x8e, 0x74, 0xf4, 0xdc, 0x81, 0xc1, 0x30, 0x44, 0x54, 0x70, 0x36, 0xdb, 0x7c, 0x41, 0x8b, 0xba, 0xd0, 0x52, 0x26, 0x59, 0xbf, 0x56, 0xa, 0x84, 0x8d, 0x89, 0xfe, 0xf1, 0x69}}
	return a, nil
}

var __2_timestampUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\xc8\xcd\x2c\x2e\xce\xcc\x4b\xf7\x4d\x2d\x2e\x4e\x4c\x4f\x2d\x56\x70\x09\xf2\x0f\x50\x70\xf6\xf7\x09\xf5\xf5\x53\xc8\x2d\x4e\x0f\xc9\xcc\x4d\x2d\x2e\x49\xcc\x2d\xb0\x06\x04\x00\x00\xff\xff\x15\xd6\x7e\xf9\x35\x00\x00\x00")

func _2_timestampUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__2_timestampUpSql,
		"2_timestamp.up.sql",
	)
}

func _2_timestampUpSql() (*asset, error) {
	bytes, err := _2_timestampUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "2_timestamp.up.sql", size: 53, mode: os.FileMode(0664), modTime: time.Unix(1720472316, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x81, 0x97, 0xcf, 0xf4, 0x2e, 0x3c, 0xca, 0x2, 0xd6, 0x9e, 0x7e, 0xf8, 0x7c, 0x18, 0x9c, 0x82, 0x25, 0x54, 0xd6, 0xd6, 0x21, 0x25, 0xdf, 0xf9, 0xdf, 0x24, 0x7a, 0xca, 0xbf, 0xeb, 0xe5, 0xd2}}
	return a, nil
}

var __3_foundUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x04\xc0\x41\x0e\x82\x40\x0c\x05\xd0\x3d\xa7\xf8\xf7\x70\x55\x9c\xb2\xfa\xb4\x09\x76\x0e\x60\xb0\x22\x31\x8e\x8b\xc6\xfb\xfb\x84\xa1\x1b\x42\x66\x2a\x3e\x67\xd5\x39\x8e\x35\xab\xee\x47\xd6\x24\xad\xe1\xea\xec\xab\xe1\xf9\xfd\x8d\x87\x8f\x2d\xf7\x57\xee\x6f\xcc\xee\x54\x31\x34\x5d\xa4\x33\xb0\x08\x6f\x0a\xf3\x80\x75\xf2\x32\xfd\x03\x00\x00\xff\xff\xc0\x0e\x74\x77\x56\x00\x00\x00")

func _3_foundUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__3_foundUpSql,
		"3_found.up.sql",
	)
}

func _3_foundUpSql() (*asset, error) {
	bytes, err := _3_foundUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "3_found.up.sql", size: 86, mode: os.FileMode(0664), modTime: time.Unix(1723057360, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x65, 0x30, 0x49, 0xb0, 0x22, 0x8c, 0xab, 0xc4, 0x3a, 0x4e, 0x39, 0xe, 0x31, 0x1, 0x53, 0xbe, 0x5c, 0x7a, 0x9a, 0xcd, 0x84, 0xe9, 0x28, 0x7d, 0xeb, 0xb4, 0x2, 0xc2, 0x3d, 0xee, 0x5d, 0x7a}}
	return a, nil
}

var __4_fleetUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\xc8\xcd\x2c\x2e\xce\xcc\x4b\xf7\x4d\x2d\x2e\x4e\x4c\x4f\x2d\x56\x70\x74\x71\x51\x70\xf6\xf7\x09\xf5\xf5\x53\x48\xcb\x49\x4d\x2d\x51\x08\x71\x8d\x08\xb1\xe6\x42\xd6\x52\x5c\x92\x5f\x94\xea\x97\x9f\x92\x1a\x9a\x97\x58\x96\x98\x99\x93\x98\x94\x93\x4a\x94\xbe\xca\xbc\xe4\x90\xfc\x82\xcc\xe4\xe0\x92\xc4\x92\x52\x9c\x56\x01\x02\x00\x00\xff\xff\x89\x58\xdc\x2d\x9e\x00\x00\x00")

func _4_fleetUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__4_fleetUpSql,
		"4_fleet.up.sql",
	)
}

func _4_fleetUpSql() (*asset, error) {
	bytes, err := _4_fleetUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4_fleet.up.sql", size: 158, mode: os.FileMode(0664), modTime: time.Unix(1724337954, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xae, 0x83, 0x6f, 0x67, 0x41, 0xfe, 0xb8, 0xd6, 0x2b, 0x27, 0xb6, 0xee, 0xa5, 0xe9, 0x52, 0x1c, 0xd4, 0xdc, 0xb5, 0xa4, 0x79, 0x15, 0x33, 0xd0, 0x8a, 0x56, 0x0, 0xbc, 0x48, 0x2f, 0x98, 0x9c}}
	return a, nil
}

var _docGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\xc9\xb1\x0d\xc4\x20\x0c\x05\xd0\x9e\x29\xfe\x02\xd8\xfd\x6d\xe3\x4b\xac\x2f\x44\x82\x09\x78\x7f\xa5\x49\xfd\xa6\x1d\xdd\xe8\xd8\xcf\x55\x8a\x2a\xe3\x47\x1f\xbe\x2c\x1d\x8c\xfa\x6f\xe3\xb4\x34\xd4\xd9\x89\xbb\x71\x59\xb6\x18\x1b\x35\x20\xa2\x9f\x0a\x03\xa2\xe5\x0d\x00\x00\xff\xff\x60\xcd\x06\xbe\x4a\x00\x00\x00")

func docGoBytes() ([]byte, error) {
	return bindataRead(
		_docGo,
		"doc.go",
	)
}

func docGo() (*asset, error) {
	bytes, err := docGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "doc.go", size: 74, mode: os.FileMode(0664), modTime: time.Unix(1715888633, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xde, 0x7c, 0x28, 0xcd, 0x47, 0xf2, 0xfa, 0x7c, 0x51, 0x2d, 0xd8, 0x38, 0xb, 0xb0, 0x34, 0x9d, 0x4c, 0x62, 0xa, 0x9e, 0x28, 0xc3, 0x31, 0x23, 0xd9, 0xbb, 0x89, 0x9f, 0xa0, 0x89, 0x1f, 0xe8}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"1_setup.up.sql": _1_setupUpSql,

	"2_timestamp.up.sql": _2_timestampUpSql,

	"3_found.up.sql": _3_foundUpSql,

	"4_fleet.up.sql": _4_fleetUpSql,

	"doc.go": docGo,
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
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"1_setup.up.sql":     &bintree{_1_setupUpSql, map[string]*bintree{}},
	"2_timestamp.up.sql": &bintree{_2_timestampUpSql, map[string]*bintree{}},
	"3_found.up.sql":     &bintree{_3_foundUpSql, map[string]*bintree{}},
	"4_fleet.up.sql":     &bintree{_4_fleetUpSql, map[string]*bintree{}},
	"doc.go":             &bintree{docGo, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
