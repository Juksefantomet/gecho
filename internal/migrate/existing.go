package migrate

import (
	"path/filepath"
	"strings"
)

// FindExisting returns the first .up.sql in dir whose name (after the timestamp) matches one of names.
func FindExisting(dir string, names ...string) (string, bool) {
	files, _ := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	for _, file := range files {
		base := strings.TrimSuffix(filepath.Base(file), ".up.sql")
		i := strings.Index(base, "_")
		if i < 0 {
			continue
		}
		for _, name := range names {
			if base[i+1:] == name {
				return file, true
			}
		}
	}
	return "", false
}
