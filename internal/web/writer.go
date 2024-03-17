package web

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Creates a writer that writes to a file while
// printing progress to stdout
func progressWriter(file *os.File, size int64) io.Writer {
	name := filepath.Base(file.Name())
	prog := progress{name: name, written: 0, total: size}
	return io.MultiWriter(file, &prog)
}

type progress struct {
	name    string
	written int64
	total   int64
}

func (prog *progress) Write(p []byte) (int, error) {
	n := len(p)
	prog.written += int64(n)
	percent := float64(prog.written) / float64(prog.total) * 100
	fmt.Printf("\rDownloading %s (%.2f%%)", prog.name, percent)
	return n, nil
}
