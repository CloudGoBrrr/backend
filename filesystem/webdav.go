package filesystem

import (
	"context"
	"os"

	"github.com/spf13/afero"
	dav "golang.org/x/net/webdav"
)

// implementation of goland.org/x/net/webdav.FileSystem
type WebdavFilesystem struct {
	Fs afero.Fs
}

func (w WebdavFilesystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return (w.Fs).Mkdir(name, perm)
}

func (w WebdavFilesystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (dav.File, error) {
	return (w.Fs).OpenFile(name, flag, perm)
}

func (w WebdavFilesystem) RemoveAll(ctx context.Context, name string) error {
	return w.Fs.RemoveAll(name)
}

func (w WebdavFilesystem) Rename(ctx context.Context, oldName, newName string) error {
	return (w.Fs).Rename(oldName, newName)
}

func (w WebdavFilesystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return (w.Fs).Stat(name)
}
