//go:build windows

package file

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Link creates or replaces a symlink in the most possible atomic way. If it fails
// to create the link, as last resort, it makes a complete copy instead.
func Link(oldname, newname string) error {
	if err := link(oldname, newname); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create link (reason: %s)\n", err)
		fmt.Println("Copying folder")
		if err := os.RemoveAll(newname); err != nil {
			return err
		}
		if err := copyFolder(oldname, newname); err != nil {
			return err
		}
	}
	return nil
}

func link(oldname, newname string) error {
	if err := os.MkdirAll(CacheDir(), os.ModePerm); err != nil {
		return err
	}
	symlinkPathTmp := filepath.Join(CacheDir(), "jdklink.tmp")
	if err := os.RemoveAll(symlinkPathTmp); err != nil {
		return err
	}
	if err := os.Symlink(oldname, symlinkPathTmp); err != nil {
		// Unable to create symlink, try to create a junction
		cmd := exec.Command("cmd", "/C", "mklink", "/J", symlinkPathTmp, oldname)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	if err := os.RemoveAll(newname); err != nil {
		return err
	}
	if err := os.Rename(symlinkPathTmp, newname); err != nil {
		return err
	}
	return nil
}

func copyFolder(srcFolder string, destFolder string) error {
	return filepath.WalkDir(srcFolder, func(srcPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(destFolder, srcPath[len(srcFolder):])

		if d.IsDir() {
			return os.MkdirAll(destPath, d.Type())
		} else {
			return copyFile(srcPath, destPath)
		}
	})
}

func copyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	return os.Chmod(destPath, srcInfo.Mode())
}
