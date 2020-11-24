// +build go1.16
//go:build go1.16

package iofs_test

import (
	"embed"
	"testing"

	st "github.com/golang-migrate/migrate/v4/source/testing"

	"github.com/johejo/golang-migrate-extra/source/iofs"
)

func Test(t *testing.T) {
	//go:embed testdata/migrations/*.sql
	var fs embed.FS
	d, err := iofs.New(fs, "testdata/migrations")
	if err != nil {
		t.Fatal(err)
	}

	st.Test(t, d)
}
