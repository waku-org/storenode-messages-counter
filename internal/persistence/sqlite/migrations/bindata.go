// Code generated by go-bindata. DO NOT EDIT.
// sources:
// 1_setup.up.sql (927B)
// 2_timestamp.up.sql (581B)
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

var __1_setupUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x91\x4d\x4f\xc2\x30\x18\xc7\xcf\xf4\x53\x3c\xc7\x2d\xe1\xa2\xd7\x9d\x26\x54\x69\x84\xcd\xb4\x45\xe0\x58\x68\x33\x9b\xec\x05\xf7\xb4\x46\xbf\xbd\xa1\x12\x82\x58\xc7\x81\x73\xfb\xfc\x5f\x7e\xff\x09\xa7\xb9\xa4\x20\xf3\x87\x39\x05\xf6\x08\x45\x29\x81\xae\x99\x90\x02\xf0\xab\xdd\xc9\x6e\x6f\x77\xc2\x29\xe7\x11\x12\x32\xda\xd5\x1e\x9d\xe9\x99\x06\x56\x48\xfa\x44\x79\xf8\x5f\x2c\xe7\xf3\x31\x19\xed\xfd\x16\xfd\x36\x5c\xc0\x6b\xce\x27\xb3\xfc\xd7\x73\xad\xd0\x89\x83\xa4\x6d\x0c\x3a\xd5\xec\x63\x1a\x2f\x9c\x2d\x72\xbe\x81\x67\xba\x81\xe4\xe4\x36\x86\x33\xed\x94\xa4\xb0\x62\x72\x56\x2e\x25\xf0\x72\xc5\xa6\x19\x21\x64\xa0\x46\x63\x11\x6d\x5b\x2d\x0c\xa2\xaa\x4c\xa8\xd1\xfb\x96\xe9\x58\xc6\x5b\xfa\x35\x3f\x06\x33\x85\x6f\xd1\x67\xac\x06\x9b\xa3\xeb\x7a\xd3\x76\xda\xfc\x73\x7c\x1c\x21\xf2\x18\x2e\x75\xee\xfe\xaa\xc2\x05\xd0\xb3\x88\x63\x38\x19\xc6\x80\x1e\x79\xb2\x62\x4a\xd7\x17\x3c\xad\xfe\x5c\x60\x75\x07\x65\x71\x89\x36\x39\x25\x99\x52\x31\x49\xb3\xeb\x2a\xf7\x31\x95\xb0\x4e\x9a\x0d\x8e\x1a\x9c\x8a\x4e\x9b\x65\xab\x3e\x94\xad\xd5\xb6\x36\x83\xcb\x0e\xe2\xed\xcd\xbb\x37\xe8\x0e\xfb\x5c\x87\x18\x2c\x6e\xc2\x27\x5c\x1f\xf0\xc5\x4a\x24\x67\x59\xd2\x8c\x7c\x07\x00\x00\xff\xff\x6e\xd3\x44\xe6\x9f\x03\x00\x00")

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

	info := bindataFileInfo{name: "1_setup.up.sql", size: 927, mode: os.FileMode(0664), modTime: time.Unix(1716212183, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xcc, 0x43, 0xbc, 0x2c, 0x21, 0xab, 0xc4, 0xe7, 0xc5, 0x35, 0xea, 0xbb, 0x4b, 0x1c, 0x98, 0x4, 0x62, 0x56, 0x6d, 0xaf, 0x18, 0x28, 0xe5, 0xf8, 0x8, 0xfd, 0xa1, 0x3, 0xa8, 0xe0, 0x65, 0xe1}}
	return a, nil
}

var __2_timestampUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x91\x41\x6f\x83\x30\x0c\x85\xcf\xf5\xaf\xf0\xb1\x48\xf9\x07\x9c\xd2\xce\xed\xa2\x41\xa8\x82\x27\xb5\xa7\x89\x96\x88\x21\xad\x50\xe1\xa0\xfd\xfd\xa9\x5d\x85\xd0\xc2\x71\xd7\x7c\x79\xcf\x7e\xcf\x5b\x47\x9a\x09\x59\x6f\x32\x42\xb3\x43\x5b\x30\xd2\xd1\x94\x5c\xe2\xb5\x15\x69\xbb\x26\xf7\x22\x55\xe3\xe5\xa3\xf3\xdf\xb8\x86\xd5\x30\x76\xa6\x46\xa6\x23\x3f\x3e\xdb\xf7\x2c\x53\xb0\xba\x7c\x8d\x12\xfc\x60\x6a\x34\x96\x69\x4f\x6e\x0e\x6f\xe3\x59\xc6\x33\xf7\xb7\xf6\x12\x09\xaf\xbf\xf6\xaf\x95\x7c\x46\x4c\x42\x3f\xf8\xae\xaf\x7d\xac\x92\xa6\x0c\x55\x18\x65\x59\x53\xeb\x80\x1b\xb3\x37\x76\x86\x10\x56\x07\x67\x72\xed\x4e\xf8\x46\x27\x5c\xcf\x06\x2b\x9c\x26\x25\x90\xa4\x00\xc6\x96\xe4\xf8\x1e\xa5\x58\xae\xe1\x51\x82\xc2\x29\xb5\xc2\x59\x46\x85\xcb\xde\x0a\xa7\xb5\x9f\xaf\xb5\x0e\x09\x94\x94\xd1\x96\xf1\xff\x2c\x61\xe7\x8a\xfc\xef\xde\x29\xc0\x8b\x2b\x0e\xcf\x4b\xc7\x50\x67\x4c\x6e\x99\xde\x23\x83\x23\xab\x73\xc2\xb8\x90\x14\xe0\x27\x00\x00\xff\xff\x7a\xcc\xb2\x8c\x45\x02\x00\x00")

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

	info := bindataFileInfo{name: "2_timestamp.up.sql", size: 581, mode: os.FileMode(0664), modTime: time.Unix(1720472437, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x54, 0x76, 0x6e, 0xf4, 0x2a, 0x56, 0x88, 0xd0, 0xe3, 0x5e, 0x7d, 0xbd, 0xec, 0x5c, 0x59, 0xfa, 0x44, 0x18, 0x82, 0xae, 0x55, 0x4c, 0xcf, 0x41, 0xa6, 0x7, 0x63, 0xba, 0x41, 0xa4, 0xfc, 0x3}}
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

	info := bindataFileInfo{name: "doc.go", size: 74, mode: os.FileMode(0664), modTime: time.Unix(1715177003, 0)}
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
