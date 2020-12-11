// +build go1.16
//go:build go1.16

package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	st "github.com/golang-migrate/migrate/v4/source/testing"
)

func Test(t *testing.T) {
	tmpDir := tmpDir(t)
	// write files that meet driver test requirements
	mustWriteFile(t, tmpDir, "1_foobar.up.sql", "1 up")
	mustWriteFile(t, tmpDir, "1_foobar.down.sql", "1 down")

	mustWriteFile(t, tmpDir, "3_foobar.up.sql", "3 up")

	mustWriteFile(t, tmpDir, "4_foobar.up.sql", "4 up")
	mustWriteFile(t, tmpDir, "4_foobar.down.sql", "4 down")

	mustWriteFile(t, tmpDir, "5_foobar.down.sql", "5 down")

	mustWriteFile(t, tmpDir, "7_foobar.up.sql", "7 up")
	mustWriteFile(t, tmpDir, "7_foobar.down.sql", "7 down")

	f := &File{}
	d, err := f.Open("file://" + tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Log(err)
		}
	})

	st.Test(t, d)
}

func TestOpen(t *testing.T) {
	tmpDir := tmpDir(t)

	mustWriteFile(t, tmpDir, "1_foobar.up.sql", "")
	mustWriteFile(t, tmpDir, "1_foobar.down.sql", "")

	if !filepath.IsAbs(tmpDir) {
		t.Fatal("expected tmpDir to be absolute path")
	}

	f := &File{}
	d, err := f.Open("file://" + tmpDir) // absolute path
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Log(err)
		}
	})
}

func TestOpenWithRelativePath(t *testing.T) {
	tmpDir := tmpDir(t)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		// rescue working dir after we are done
		if err := os.Chdir(wd); err != nil {
			t.Log(err)
		}
	})

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(filepath.Join(tmpDir, "foo"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	mustWriteFile(t, filepath.Join(tmpDir, "foo"), "1_foobar.up.sql", "")

	f := &File{}

	// dir: foo
	d, err := f.Open("file://foo")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.First()
	if err != nil {
		t.Fatalf("expected first file in working dir %v for foo", tmpDir)
	}

	// dir: ./foo
	d, err = f.Open("file://./foo")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Log(err)
		}
	})
	_, err = d.First()
	if err != nil {
		t.Fatalf("expected first file in working dir %v for ./foo", tmpDir)
	}
}

func TestOpenDefaultsToCurrentDirectory(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	f := &File{}
	d, err := f.Open("file://")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Log(err)
		}
	})

	if d.(*File).path != wd {
		t.Fatal("expected driver to default to current directory")
	}
}

func TestOpenWithDuplicateVersion(t *testing.T) {
	tmpDir := tmpDir(t)

	mustWriteFile(t, tmpDir, "1_foo.up.sql", "") // 1 up
	mustWriteFile(t, tmpDir, "1_bar.up.sql", "") // 1 up

	f := &File{}
	_, err := f.Open("file://" + tmpDir)
	if err == nil {
		t.Fatal("expected err")
	}
}

func TestClose(t *testing.T) {
	tmpDir := tmpDir(t)

	f := &File{}
	d, err := f.Open("file://" + tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if d.Close() != nil {
		t.Fatal("expected nil")
	}
}

func mustWriteFile(tb testing.TB, dir, file string, body string) {
	tb.Helper()
	if err := ioutil.WriteFile(path.Join(dir, file), []byte(body), 06444); err != nil {
		tb.Fatal(err)
	}
}

func mustCreateBenchmarkDir(b *testing.B) (dir string) {
	b.Helper()
	tmpDir := tmpDir(b)
	for i := 0; i < 1000; i++ {
		mustWriteFile(b, tmpDir, fmt.Sprintf("%v_foobar.up.sql", i), "")
		mustWriteFile(b, tmpDir, fmt.Sprintf("%v_foobar.down.sql", i), "")
	}
	return tmpDir
}

func tmpDir(tb testing.TB) string {
	tb.Helper()
	return filepath.ToSlash(tb.TempDir())
}

func BenchmarkOpen(b *testing.B) {
	dir := mustCreateBenchmarkDir(b)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		f := &File{}
		_, err := f.Open("file://" + dir)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkNext(b *testing.B) {
	dir := mustCreateBenchmarkDir(b)
	f := &File{}
	d, err := f.Open("file://" + dir)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if err := d.Close(); err != nil {
			b.Log(err)
		}
	})
	b.ResetTimer()
	v, err := d.First()
	for n := 0; n < b.N; n++ {
		for !errors.Is(err, os.ErrNotExist) {
			v, err = d.Next(v)
		}
	}
	b.StopTimer()
}
