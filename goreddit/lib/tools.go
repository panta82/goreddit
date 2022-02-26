package lib

import (
	"os"
	"path"
)

// Given a list of directories, return the first one that exists. Nil if none exist.
// Resolution starts from CWD.
func LookupRelativePath(lookupPaths []string, startDir *string) (*string, error) {
	if startDir == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		startDir = &cwd
	}

	for _, lookupDir := range lookupPaths {
		targetPath := path.Clean(path.Join(*startDir, lookupDir))

		if _, err := os.Stat(lookupDir); err == nil {
			return &targetPath, nil
		}
	}

	return nil, nil
}
