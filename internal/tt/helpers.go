package tt

import (
	"bytes"
	"io"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func True(t *testing.T, a bool, rest ...any) {
	if !a {
		t.Helper()
		t.Errorf("expected true, but false; %v", rest)
	}
}

func Eq[T comparable](t *testing.T, a, b T, rest ...any) {
	if a != b {
		t.Helper()
		t.Errorf("expected %v, but %v; %v", a, b, rest)
	}
}

func EqSlice[S ~[]E, E comparable](t *testing.T, a, b S, rest ...any) {
	if !slices.Equal(a, b) {
		t.Helper()
		t.Errorf("slices not equal; %v", rest)
	}
}

func EqText(t *testing.T, a, b string, rest ...any) {
	as := strings.Split(a, "\n")
	bs := strings.Split(b, "\n")
	ok := len(as) == len(bs)
	diff := []string{}
	for x, n := 0, min(len(as), len(bs)); x < n; x++ {
		matched := as[x] == bs[x]
		if len(as[x]) > 0 && as[x][0] == '~' {
			matched = regexp.MustCompile(string(as[x][1:])).MatchString(bs[x])
		}

		if matched {
			diff = append(diff, " "+as[x])
		} else {
			diff = append(diff, "-"+as[x], "+"+bs[x])
			ok = false
		}
	}

	if !ok {
		t.Helper()
		t.Errorf("text not equal\n%s; %v", strings.Join(diff, "\n"), rest)
	}
}

func Prefix(t *testing.T, a, b string, rest ...any) {
	if !strings.HasPrefix(b, a) {
		t.Helper()
		t.Errorf("expected %v..., but %v; %v", a, b, rest)
	}
}

func Unindent(s string) string {
	n := 0
	for n < len(s) && s[n] != '\n' {
		n++
	}
	if n == len(s) {
		return s
	}

	m := n + 1
	for m < len(s) && (s[m] == ' ' || s[m] == '\t') {
		m++
	}
	s = strings.ReplaceAll(s, s[n:m], "\n")
	return strings.TrimRight(s[n+1:], " \t\n")
}

func Flush(w io.Writer) []byte {
	b := w.(*bytes.Buffer)
	s := b.Bytes()
	b.Reset()
	return s
}

func FlushString(w io.Writer) string {
	b := w.(*bytes.Buffer)
	s := b.String()
	b.Reset()
	return s
}
