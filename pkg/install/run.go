package install

import (
	"io/ioutil"

	"ecksbee.com/kushim/internal/actions"
	"ecksbee.com/kushim/pkg/taxonomies"
	"ecksbee.com/telefacts/pkg/serializables"
)

func Run(taxonomyPackage string, volumePath string, throttle func(string)) (string, error) {
	bytes, err := ioutil.ReadFile(taxonomyPackage)
	if err != nil {
		return "", err
	}
	unZipFiles, err := actions.Unzip(bytes)
	if err != nil {
		return "", err
	}
	name, err := name(unZipFiles)
	if err != nil {
		return "", err
	}
	entries, err := entries(unZipFiles)
	if err != nil {
		return "", err
	}
	remap, err := remap(unZipFiles)
	if err != nil {
		return "", err
	}
	cleanedBytes, err := clean(unZipFiles)
	if err != nil {
		return "", err
	}
	taxonomies.VolumePath = volumePath
	serializables.GlobalTaxonomySetPath = volumePath
	err = DownloadUTR(throttle)
	if err != nil {
		return "", err
	}
	err = DownloadLRR(throttle)
	if err != nil {
		return "", err
	}
	err = DownloadDTRs(throttle)
	if err != nil {
		return "", err
	}
	return taxonomies.NewTaxonomy(taxonomies.Meta{
		Name:    name,
		Zip:     taxonomyPackage,
		Entries: entries,
		Remap:   remap,
	}, cleanedBytes)
}
