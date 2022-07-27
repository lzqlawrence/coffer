package filesystem

import (
	"context"
	"github.com/pkg/errors"
	"os"
	"path/filepath"

	"golang.org/x/net/webdav"
)

func init() {
	Register("xor", func(folder string, param map[string]string) (webdav.FileSystem, error) {
		key, exist := param["key"]
		if !exist {
			return nil, errors.Errorf("key not exist:%v", param)
		}
		return &xorFS{folder: folder, key: []byte(key)}, nil
	})
}

type xorFS struct {
	folder string
	key    []byte
}

func (fs *xorFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	if name = ResolvePath(fs.folder, name); name == "" {
		return os.ErrNotExist
	}
	return os.Mkdir(name, perm)
}

func (fs *xorFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if name = ResolvePath(fs.folder, name); name == "" {
		return nil, os.ErrNotExist
	}
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (fs *xorFS) RemoveAll(ctx context.Context, name string) error {
	if name = ResolvePath(fs.folder, name); name == "" {
		return os.ErrNotExist
	}
	if name == filepath.Clean(fs.folder) {
		// Prohibit removing the virtual root directory.
		return os.ErrInvalid
	}
	return os.RemoveAll(name)
}
func (fs *xorFS) Rename(ctx context.Context, oldName, newName string) error {
	if oldName = ResolvePath(fs.folder, oldName); oldName == "" {
		return os.ErrNotExist
	}
	if newName = ResolvePath(fs.folder, newName); newName == "" {
		return os.ErrNotExist
	}
	if root := filepath.Clean(fs.folder); root == oldName || root == newName {
		// Prohibit renaming from or to the virtual root directory.
		return os.ErrInvalid
	}
	return os.Rename(oldName, newName)
}
func (fs *xorFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name = ResolvePath(fs.folder, name); name == "" {
		return nil, os.ErrNotExist
	}
	return os.Stat(name)
}

type xorFile struct {
	pos int64
}

func (f *xorFile) Close() error {
	return nil
}

func (f *xorFile) Read(p []byte) (int, error) {
	return 0, nil
}

func (f *xorFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *xorFile) Seek(offset int64, whence int) (int64, error) {

	return 0, nil
}

func (f *xorFile) Stat() (os.FileInfo, error) {
	return nil, nil
}

func (f *xorFile) Write(p []byte) (int, error) {
	return 0, nil
}
