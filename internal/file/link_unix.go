//go:build !windows

package file

import (
	"os"
	"path/filepath"
)

// Link creates or replaces a symlink in the most possible atomic way. Copied from
// https://stackoverflow.com/questions/37345844/how-to-overwrite-a-symlink-in-go
func Link(oldname, newname string) error {
	if err := os.MkdirAll(CacheDir(), os.ModePerm); err != nil {
		return err
	}
	symlinkPathTmp := filepath.Join(CacheDir(), "jdklink.tmp")
	if err := os.RemoveAll(symlinkPathTmp); err != nil {
		return err
	}
	if err := os.Symlink(oldname, symlinkPathTmp); err != nil {
		return err
	}
	if err := os.Rename(symlinkPathTmp, newname); err != nil {
		return err
	}
	return nil
}
