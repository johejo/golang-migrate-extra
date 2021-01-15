// +build go1.16
//go:build go1.16

package iofs_test

import (
	"embed"
	"testing"

	st "github.com/golang-migrate/migrate/v4/source/testing"

	"github.com/johejo/golang-migrate-extra/source/iofs"
)

//go:embed testdata/migrations/*.sql
var fsys embed.FS

func Test(t *testing.T) {
	d, err := iofs.New(fsys, "testdata/migrations")
	if err != nil {
		t.Fatal(err)
	}

	st.Test(t, d)
}
