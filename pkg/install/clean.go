package install

import (
	"archive/zip"
	"path"

	"ecksbee.com/kushim/internal/actions"
)

var (
	METAINF string = "META-INF"
)

func clean(unZipFiles []*zip.File) ([]byte, error) {
	cleanedFiles := make([]*zip.File, 0)
	for _, file := range unZipFiles {
		dir := path.Base(path.Dir(file.Name))
		if METAINF != dir {
			cleanedFiles = append(cleanedFiles, file)
		}
	}
	return actions.Zip(cleanedFiles)
}
