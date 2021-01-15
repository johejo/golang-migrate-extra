package file_test

import (
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/stub"
	_ "github.com/johejo/golang-migrate-extra/source/file"
)

func Example() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	sourceURL := "file://" + path.Join(wd, "testdata")
	println(sourceURL)
	m, err := migrate.New(sourceURL, "stub://")
	if err != nil {
		panic(err)
	}
	_ = m
	// m.Up()
}
