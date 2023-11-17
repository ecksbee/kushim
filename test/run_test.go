package kushim_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"ecksbee.com/kushim/pkg/install"
	"ecksbee.com/kushim/pkg/taxonomies"
	"ecksbee.com/kushim/pkg/throttle"
)

func Test_Run_USGAAP2020(t *testing.T) {
	throttle.StartSECThrottle()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tp := filepath.Join(wd, "us-gaap-2020-01-31.zip")
	volume := filepath.Join(wd, "gts")
	id, err := install.Run(tp, volume, throttle.Throttle)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal(fmt.Errorf("empty id generated"))
	}
}

func Test_Run_USGAAP2019(t *testing.T) {
	throttle.StartSECThrottle()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tp := filepath.Join(wd, "us-gaap-2019-01-31.zip")
	volume := filepath.Join(wd, "gts")
	id, err := install.Run(tp, volume, throttle.Throttle)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal(fmt.Errorf("empty id generated"))
	}
}

func Test_Run_ESEF2017(t *testing.T) {
	throttle.StartSECThrottle()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tp := filepath.Join(wd, "esef_taxonomy_2017.zip")
	volume := filepath.Join(wd, "gts")
	id, err := install.Run(tp, volume, throttle.Throttle)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal(fmt.Errorf("empty id generated"))
	}
}

func Test_Run_CMF_CL_CI2020(t *testing.T) {
	tm := taxonomies.Meta{
		Name: "CMF CL-CI 2020",
		Zip:  "articles-28084_recurso_1.zip",
		Entries: []string{
			"http://www.cmfchile.cl/cl/fr/ci/2020-01-02/cl-ci_shell_2020-01-02.xsd",
		},
		Remap: map[string]string{
			"": "http://www.cmfchile.cl/cl/fr/ci/2020-01-02/",
		},
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tp := filepath.Join(wd, tm.Zip)
	bytes, err := ioutil.ReadFile(tp)
	if err != nil {
		t.Fatal(err)
	}
	taxonomies.VolumePath = filepath.Join(wd, "gts")
	id, err := taxonomies.NewTaxonomy(tm, bytes)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal(fmt.Errorf("empty id generated"))
	}
}
