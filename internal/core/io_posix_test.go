//go:build !wasm

package core

import (
	"ocala/internal/tt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReplacePathExt(t *testing.T) {
	data := []struct {
		expected string
		path     string
		ext      string
	}{
		{"/.ext", "/", ".ext"},
		{"/a.ext", "/a", ".ext"},
		{"/a.ext", "/a.oc", ".ext"},
		{"/tmp/.ext", "/tmp/", ".ext"},
		{"/tmp/a.ext", "/tmp/a", ".ext"},
		{"/tmp/a.ext", "/tmp/a.oc", ".ext"},
	}
	for x, i := range data {
		actual := replacePathExt(i.path, i.ext)
		tt.Eq(t, i.expected, actual, x, i)
	}
}

func lastPathComponents(s string, n int) string {
	s = filepath.ToSlash(s)
	x := len(s)
	for range n {
		x = strings.LastIndex(s[:x], "/")
		if x == -1 {
			return ""
		}
	}
	return s[x:]
}

func TestRegularizePath(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		wd, _ := os.Getwd()
		paths := []string{
			filepath.Join(wd, "testdata/include"),
			filepath.Join(wd, "testdata/include2"),
		}
		data := []struct {
			expected string
			path     string
			n        int
		}{
			{"/io_posix_test.go", "./io_posix_test.go", 1},
			{"/testdata/test.oc", "./testdata/test.oc", 2},
			{"/testdata/include/inc01.oc", "./testdata/include/inc01.oc", 3},
			{"/testdata/include2/inc01.oc", "./testdata/include2/inc01.oc", 3},
			{"/testdata/include/inc01.oc", "inc01.oc", 3},
			{"/testdata/include2/inc21.oc", "inc21.oc", 3},
		}
		for x, i := range data {
			s, _ := regularizePath(i.path, wd, paths)
			actual := lastPathComponents(s, i.n)
			tt.Eq(t, i.expected, actual, x, i)
		}
	})

	t.Run("error", func(t *testing.T) {
		wd, _ := os.Getwd()
		paths := []string{}
		data := []struct {
			expected string
			path     string
		}{
			{"invalid path `/`", "/"},
			{"invalid path `/a`", "/a"},

			{"invalid path `.`", "."},
			{"invalid path `..`", ".."},
			{"invalid path `...`", "..."},

			{"invalid path `./`", "./"},
			{"invalid path `../`", "../"},
			{"invalid path `.../`", ".../"},

			{"invalid path `../a`", "../a"},
			{"invalid path `.../a`", ".../a"},

			{"invalid path `a/`", "a/"},
			{"invalid path `a/.`", "a/."},
			{"invalid path `a/..`", "a/.."},
			{"invalid path `a/...`", "a/..."},

			{"invalid path `a/./`", "a/./"},
			{"invalid path `a/../`", "a/../"},
			{"invalid path `a/.../`", "a/.../"},

			{"invalid path `a/./a`", "a/./a"},
			{"invalid path `a/../a`", "a/../a"},
			{"invalid path `a/.../a`", "a/.../a"},

			{"the file `./nothing` not found", "./nothing"},
			{"the file `./nothing/nothing` not found", "./nothing/nothing"},
			{"the file `./testdata/nothing` not found", "./testdata/nothing"},
		}
		for x, i := range data {
			_, err := regularizePath(i.path, wd, paths)
			tt.Eq(t, i.expected, err.Error(), x, i)
		}
	})

	t.Run("error: windows", func(t *testing.T) {
		wd, _ := os.Getwd()
		paths := []string{}
		data := []struct {
			expected string
			path     string
		}{
			{"invalid path `C:/`", `C:/`},
			{"invalid path `//./C:/`", `//./C:/`},
			{"invalid path `C:/a`", `C:/a`},
			{"invalid path `//./C:/a`", `//./C:/a`},

			{"invalid path `C:\\`", `C:\`},
			{"invalid path `\\\\.\\C:\\`", `\\.\C:\`},
			{"invalid path `C:\\a`", `C:\a`},
			{"invalid path `\\\\.\\C:\\a`", `\\.\C:\a`},
		}
		for x, i := range data {
			if runtime.GOOS == "windows" {
				_, err := regularizePath(i.path, wd, paths)
				tt.Eq(t, i.expected, err.Error(), x, i)
			} else {
				tt.Eq(t, 1, 1, x, i)
			}
		}
	})
}
