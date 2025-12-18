package langserver

import (
	"ocala/language-server/file"
	tt "ocala/testutil"
	"path/filepath"
	"testing"
)

func TestCompletion(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		h := newTestHandler()
		setFileText(h, "file:///test.oc", tt.UnindentBytes(`
			arch z80
			module n001 {
			  data d001 = byte
			  data d002 = struct { abc byte }
			  c0
			  d0
			  d002.
			  d002.ab
			  d002.de
			};
			const c001 = 1
			c0
			d0
			n0
			n001:
			n001:d0
			n001:c0
			modu
			f0
			l0
			aswor
			const f001(x) = 1
			l001: NOP
			//
		`))

		es := []struct {
			s string
			l int
			c int
			n int
		}{
			{s: "-", l: 120, c: 120, n: 0},
			{s: "n001>c0*", l: 4, c: 120, n: 1},
			{s: "n001>d0*", l: 5, c: 120, n: 2},
			{s: "n001>d001.*", l: 6, c: 120, n: 1},
			{s: "n001>d001.ab*", l: 7, c: 120, n: 1},
			{s: "n001>d001.de*", l: 8, c: 120, n: 0},
			{s: ";", l: 9, c: 120, n: 0},
			{s: "c0*", l: 11, c: 120, n: 1},
			{s: "d0*", l: 12, c: 120, n: 0},
			{s: "n0*", l: 13, c: 120, n: 1},
			{s: "n001:*", l: 14, c: 120, n: 3},
			{s: "n001:d0*", l: 15, c: 120, n: 2},
			{s: "n001:c0*", l: 16, c: 120, n: 0},
			{s: "modu*", l: 17, c: 120, n: 1},
			{s: "f0*", l: 18, c: 120, n: 1},
			{s: "l0*", l: 19, c: 120, n: 1},
			{s: "aswor*", l: 20, c: 120, n: 1},
		}
		for _, i := range es {
			result, err := h.completion("file:///test.oc", &CompletionParams{
				TextDocumentPositionParams: TextDocumentPositionParams{
					TextDocument: TextDocumentIdentifier{URI: "file:///test"},
					Position:     Position{Line: i.l, Character: i.c},
				},
			})
			tt.Eq(t, nil, err, i.s)
			tt.Eq(t, i.n, len(result), i.s, result)
		}
	})

	t.Run("ok: chain", func(t *testing.T) {
		h := newTestHandler()
		setFileText(h, "file:///test.oc", tt.UnindentBytes(`
			arch z80
			module n001 {
			  const c001 = 1
			  struct s001 {
			    a001 byte
			    a002 struct {
			      b003 struct {
			        c004 byte
			        c005 byte
			      }
			      b006 struct {
			        d007 byte
			        d008 byte
			      }
			    }
			  }
			}
			n001:
			n001:c00
			n001:s00
			n001:s001.a00
			n001:s001.a002.b00
			n001:s001.a002.b003.
			n001:s001.a002.b003.c00
			n001:s001.a002.b006.
			n001:s001.a002.b00b.d00
			n001:c001.
			n002:c001.
			//
		`))

		es := []struct {
			s string
			l int
			c int
			n int
		}{
			{s: "-", l: 120, c: 120, n: 0},
			{s: "n001:", l: 17, c: 120, n: 2},
			{s: "n001:s00", l: 18, c: 120, n: 1},
			{s: "n001:s001.a00", l: 19, c: 120, n: 1},
			{s: "n001:s001.a002.b00", l: 20, c: 120, n: 2},
			{s: "n001:s001.a002.b003.*", l: 21, c: 120, n: 2},
			{s: "n001:s001.a002.b003.c00*", l: 22, c: 120, n: 2},
			{s: "n001:s001.a002.b006.*", l: 23, c: 120, n: 2},
			{s: "n001:s001.a002.b006.d0*", l: 24, c: 120, n: 2},
			{s: "n001:c001.", l: 25, c: 120, n: 0},
			{s: "n002:c001.", l: 25, c: 120, n: 0},
		}
		for _, i := range es {
			result, err := h.completion("file:///test.oc", &CompletionParams{
				TextDocumentPositionParams: TextDocumentPositionParams{
					Position: Position{Line: i.l, Character: i.c},
				},
			})
			tt.Eq(t, nil, err, i.s)
			tt.Eq(t, i.n, len(result), i.s, result)
		}
	})

	t.Run("ok: include", func(t *testing.T) {
		h := newTestHandler()
		path := tt.Must(filepath.Abs("../file/testdata"))
		h.compileOptions.IncPaths = file.AdjustIncPaths([]string{path})
		setFileText(h, "file:///test.oc", tt.UnindentBytes(`
			arch z80
			include "test3.oc"
			const c001 = 1
			c00
		`))
		result, err := h.completion("file:///test.oc", &CompletionParams{
			TextDocumentPositionParams: TextDocumentPositionParams{
				Position: Position{Line: 3, Character: 120},
			},
		})
		tt.Eq(t, nil, err)
		tt.Eq(t, 3, len(result))
	})

	t.Run("error", func(t *testing.T) {
		h := newTestHandler()
		setFileText(h, "file:///test.oc", []byte("arch z80"))
	})
}

func TestSymbol(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		h := newTestHandler()
		setFileText(h, "file:///test.oc", tt.UnindentBytes(`
			arch z80
			macro mac() { NOP }
			proc fn() { data d001 = byte; RET }
			const a = 1
			const f(x) = x + 1
			data d = byte [a]
			module mod { data d = byte }
			L001:
			struct s { a byte; b byte}
			if 1 { NOP }
		`))
		result, err := h.symbol("file:///test.oc")
		tt.Eq(t, nil, err)
		tt.Eq(t, 8, len(result))
	})

	t.Run("error", func(t *testing.T) {
		h := newTestHandler()
		setFileText(h, "file:///test.oc", []byte("arch z80"))
	})
}

func TestDefinition(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		h := newTestHandler()
		setFileText(h, "file:///test.oc", tt.UnindentBytes(`
			arch z80
			module n001 {
			  data d001 = byte
			  data d002 = struct { abc byte }
			  c001
			  d001
			  d002
			  d002.abc
			  d002.def
			};
			const c001 = 1
			const f001(x) = 1
			c001
			d001
			n001
			n001:d001
			n001:c001
			asword
		`))

		es := []struct {
			s string
			l int
			c int
			n int
		}{
			{s: "-", l: 120, c: 120, n: 0},
			{s: "n001>c001", l: 4, c: 120, n: 1},
			{s: "n001>d001", l: 5, c: 120, n: 1},
			{s: "n001>d002", l: 6, c: 120, n: 1},
			{s: "n001>d002.abc", l: 7, c: 120, n: 1},
			{s: "n001>d002.def", l: 8, c: 120, n: 0},
			{s: ";", l: 9, c: 120, n: 0},
			{s: "c001", l: 12, c: 120, n: 1},
			{s: "d001*", l: 13, c: 120, n: 0},
			{s: "n001", l: 14, c: 120, n: 1},
			{s: "n001:d001", l: 15, c: 120, n: 1},
			{s: "n001:c001", l: 16, c: 120, n: 0},
			{s: "asword", l: 17, c: 120, n: 0},
		}
		for _, i := range es {
			result, err := h.definition("file:///test.oc", &DocumentDefinitionParams{
				TextDocumentPositionParams: TextDocumentPositionParams{
					Position: Position{Line: i.l, Character: i.c},
				},
			})
			tt.Eq(t, nil, err, i.s)
			tt.Eq(t, i.n, len(result), i.s, result)
		}
	})

	t.Run("error", func(t *testing.T) {
		h := newTestHandler()
		setFileText(h, "file:///test.oc", []byte("arch z80"))
	})
}
