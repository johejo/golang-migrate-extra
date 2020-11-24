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

	"github.com/golang-migrate/migrate/v4/source"
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

	st.TestFirst(t, d)
	st.TestPrev(t, d)
	st.TestNext(t, d)
	testReadUp(t, d)
	testReadDown(t, d)
}

// from github.com/golang-migrate/migrate/v4/source/testing#TestReadUp
func testReadUp(t *testing.T, d source.Driver) {
	tt := []struct {
		version   uint
		expectErr error
		expectUp  bool
	}{
		{version: 0, expectErr: os.ErrNotExist},
		{version: 1, expectErr: nil, expectUp: true},
		{version: 2, expectErr: os.ErrNotExist},
		{version: 3, expectErr: nil, expectUp: true},
		{version: 4, expectErr: nil, expectUp: true},
		{version: 5, expectErr: os.ErrNotExist},
		{version: 6, expectErr: os.ErrNotExist},
		{version: 7, expectErr: nil, expectUp: true},
		{version: 8, expectErr: os.ErrNotExist},
	}

	for i, v := range tt {
		up, identifier, err := d.ReadUp(v.version)
		if (v.expectErr == os.ErrNotExist && !errors.Is(err, os.ErrNotExist)) ||
			(v.expectErr != os.ErrNotExist && err != v.expectErr) {
			t.Errorf("expected %v, got %v, in %v", v.expectErr, err, i)

		} else if err == nil {
			if len(identifier) == 0 {
				t.Errorf("expected identifier not to be empty, in %v", i)
			}

			if v.expectUp && up == nil {
				t.Errorf("expected up not to be nil, in %v", i)
			} else if !v.expectUp && up != nil {
				t.Errorf("expected up to be nil, got %v, in %v", up, i)
			}
		}
		if up != nil {
			defer up.Close()
		}
	}
}

// from github.com/golang-migrate/migrate/v4/source/testing#TestReadDown
func testReadDown(t *testing.T, d source.Driver) {
	tt := []struct {
		version    uint
		expectErr  error
		expectDown bool
	}{
		{version: 0, expectErr: os.ErrNotExist},
		{version: 1, expectErr: nil, expectDown: true},
		{version: 2, expectErr: os.ErrNotExist},
		{version: 3, expectErr: os.ErrNotExist},
		{version: 4, expectErr: nil, expectDown: true},
		{version: 5, expectErr: nil, expectDown: true},
		{version: 6, expectErr: os.ErrNotExist},
		{version: 7, expectErr: nil, expectDown: true},
		{version: 8, expectErr: os.ErrNotExist},
	}

	for i, v := range tt {
		down, identifier, err := d.ReadDown(v.version)
		if (v.expectErr == os.ErrNotExist && !errors.Is(err, os.ErrNotExist)) ||
			(v.expectErr != os.ErrNotExist && err != v.expectErr) {
			t.Errorf("expected %v, got %v, in %v", v.expectErr, err, i)
		} else if err == nil {
			if len(identifier) == 0 {
				t.Errorf("expected identifier not to be empty, in %v", i)
			}

			if v.expectDown && down == nil {
				t.Errorf("expected down not to be nil, in %v", i)
			} else if !v.expectDown && down != nil {
				t.Errorf("expected down to be nil, got %v, in %v", down, i)
			}
		}
		if down != nil {
			defer down.Close()
		}
	}
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
