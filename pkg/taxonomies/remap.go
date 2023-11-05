package taxonomies

import (
	"fmt"
	"os"
	"path"
	"strings"

	"ecksbee.com/kushim/internal/actions"
	"ecksbee.com/telefacts/pkg/serializables"
)

func Remap(bytes []byte, remap map[string]string) error {
	if VolumePath == "" {
		return fmt.Errorf("empty VolumePath")
	}
	serializables.GlobalTaxonomySetPath = VolumePath
	unZipFiles, err := actions.Unzip(bytes)
	if err != nil {
		return err
	}
	for _, unZipFile := range unZipFiles {
		ext := path.Ext(unZipFile.Name)
		if ext != ".xsd" && ext != ".xml" {
			continue
		}
		dir := path.Dir(unZipFile.Name)
		for oldPrefix, newPrefix := range remap {
			if strings.HasPrefix(dir, oldPrefix) {
				url := strings.Replace(unZipFile.Name, oldPrefix,
					newPrefix, 1)
				dest, err := serializables.UrlToFilename(url)
				if err != nil {
					return err
				}
				_, err = os.Stat(dest)
				if !os.IsNotExist(err) {
					continue
				}
				targetDir := path.Dir(dest)
				_, err = os.Stat(targetDir)
				if os.IsNotExist(err) {
					err = os.MkdirAll(targetDir, 0755)
					if err != nil {
						return err
					}
				}
				unzipped, err := actions.UnzipFile(unZipFile)
				if err != nil {
					return err
				}
				err = actions.WriteFile(dest, unzipped)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
