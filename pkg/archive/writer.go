package archive

import (
	"archive/tar"
	"io"
	"os"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
)

// StreamTo writes the given entries as a tar stream to w.
// Bytes flowing through w can be wrapped by the caller to drive a progress bar.
func StreamTo(w io.Writer, entries []domain.FileEntry) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	for _, e := range entries {
		hdr := &tar.Header{
			Name: e.RelPath,
			Mode: int64(e.Mode),
			Size: e.Size,
		}
		if e.IsDir {
			hdr.Typeflag = tar.TypeDir
			hdr.Size = 0
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			continue
		}

		hdr.Typeflag = tar.TypeReg
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if err := copyFile(tw, e.AbsPath); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(dst io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(dst, f)
	return err
}
