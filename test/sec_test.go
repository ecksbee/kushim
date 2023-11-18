package kushim_test

import (
	"os"
	"path/filepath"
	"testing"

	"ecksbee.com/kushim/pkg/install"
	"ecksbee.com/kushim/pkg/librarian"
)

func Test_InstallSECTaxonomies(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	librarian.IndexingMode = true
	gts := filepath.Join(wd, "gts")
	err = install.InstallSECTaxonomies(gts)
	if err != nil {
		t.Fatal(err)
	}
	librarian.ProcessIndex()
}
