package taxonomies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ecksbee.com/telefacts/pkg/serializables"
	"github.com/google/uuid"
)

var (
	indexingMode bool
	VolumePath   string
)

type Meta struct {
	Name    string
	Zip     string
	Entries []string
	Remap   map[string]string
}

func NewTaxonomy(tm Meta, bytes []byte) (string, error) {
	if VolumePath == "" {
		return "", fmt.Errorf("empty VolumePath")
	}
	err := Remap(bytes, tm.Remap)
	if err != nil {
		return "", err
	}
	err = Discover(tm.Entries)
	if err != nil {
		return "", err
	}
	serializables.GlobalTaxonomySetPath = VolumePath
	workingDir := filepath.Join(VolumePath, "taxonomies")
	_, err = os.Stat(workingDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(workingDir, 0755)
		if err != nil {
			return "", err
		}
	}
	id := uuid.New()
	pathStr := filepath.Join(workingDir, id.String())
	_, err = os.Stat(pathStr)
	for !os.IsNotExist(err) {
		id = uuid.New()
		pathStr = filepath.Join(workingDir, id.String())
		_, err = os.Stat(pathStr)
	}
	err = os.Mkdir(pathStr, 0755)
	if err != nil {
		return "", err
	}
	meta := filepath.Join(pathStr, "_")
	file, _ := os.OpenFile(meta, os.O_CREATE|os.O_WRONLY, 0755)
	defer file.Close()
	encoder := json.NewEncoder(file)
	return id.String(), encoder.Encode(tm)
}
