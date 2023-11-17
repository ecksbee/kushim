package kushim_test

import (
	"os"
	"path/filepath"
	"testing"

	"ecksbee.com/kushim/pkg/install"
)

func Test_InstallSECTaxonomies(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	gts := filepath.Join(wd, "gts")
	err = install.InstallSECTaxonomies(gts)
	if err != nil {
		t.Fatal(err)
	}
}
