package fsutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/viniciusmichelutti/ftransfer/internal/domain"
)

func Walk(sources []string) ([]domain.FileEntry, int64, error) {
	var entries []domain.FileEntry
	var total int64

	for _, src := range sources {
		absSrc, err := filepath.Abs(src)
		if err != nil {
			return nil, 0, err
		}
		info, err := os.Stat(absSrc)
		if err != nil {
			return nil, 0, fmt.Errorf("stat %q: %w", src, err)
		}

		// Use the base name as the root entry name (so "/a/b/photos" arrives as "photos/...").
		root := filepath.Base(absSrc)

		if !info.IsDir() {
			entries = append(entries, domain.FileEntry{
				RelPath: root,
				Size:    info.Size(),
				Mode:    uint32(info.Mode().Perm()),
				AbsPath: absSrc,
			})
			total += info.Size()
			continue
		}

		err = filepath.Walk(absSrc, func(path string, fi os.FileInfo, werr error) error {
			if werr != nil {
				return werr
			}
			rel, err := filepath.Rel(absSrc, path)
			if err != nil {
				return err
			}
			relPath := filepath.ToSlash(filepath.Join(root, rel))
			if fi.IsDir() {
				if rel == "." {
					relPath = root
				}
				entries = append(entries, domain.FileEntry{
					RelPath: relPath,
					Mode:    uint32(fi.Mode().Perm()),
					IsDir:   true,
					AbsPath: path,
				})
				return nil
			}
			entries = append(entries, domain.FileEntry{
				RelPath: relPath,
				Size:    fi.Size(),
				Mode:    uint32(fi.Mode().Perm()),
				AbsPath: path,
			})
			total += fi.Size()
			return nil
		})
		if err != nil {
			return nil, 0, err
		}
	}

	return entries, total, nil
}
