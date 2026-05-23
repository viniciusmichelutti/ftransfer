package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/viniciusmichelutti/ftransfer/pkg/fsutil"
)

// ExtractTo reads a tar stream from r and extracts entries into dir.
// Every entry path is validated by fsutil.SafeJoin to block path traversal.
func ExtractTo(r io.Reader, dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		target, err := fsutil.SafeJoin(dir, hdr.Name)
		if err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)&os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if err := writeFile(target, tr, os.FileMode(hdr.Mode)&os.ModePerm); err != nil {
				return err
			}
		}
	}
}

func writeFile(path string, r io.Reader, mode os.FileMode) error {
	if mode == 0 {
		mode = 0o644
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
