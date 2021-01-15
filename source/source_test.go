package source_test

import (
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/stub"
	"github.com/golang-migrate/migrate/v4/source"

	"github.com/johejo/golang-migrate-extra/source/file"
)

func Test_source_Open(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	u := "file://" + path.Join(wd, "testdata")
	d, err := source.Open(u)
	if err != nil {
		t.Fatal(err)
	}
	dt := reflect.ValueOf(d).Elem().Type()
	ft := reflect.ValueOf(new(file.File)).Elem().Type()
	if dt != ft {
		t.Errorf("want=%v, but opend driver type=%v", ft, dt)
	}
	_, err = migrate.NewWithSourceInstance("file", d, "stub://")
	if err != nil {
		t.Fatal(err)
	}
}
