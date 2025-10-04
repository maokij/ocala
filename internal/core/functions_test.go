package core_test

import (
	"ocala/internal/tt"
	"ocala/internal/tt/ttarch"
	"testing"
)

func TestCompileFunctions(t *testing.T) {
	t.Run("ok: operators", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			link-as-tests
			expect 6   (2 * 3)
			expect -15 (3 * -5)
			expect 25  (-5 * -5)
			expect 3   (10 / 3)
			expect -3  (12 / -4)
			expect 3   (-20 / -6)
			expect 1   (10 % 3)
			expect 0   (12 % -4)
			expect -2  (-20 % -6)
			expect 3   (1 + 2)
			expect 5   (10 + -5)
			expect -3  (-1 + -2)
			expect 8   (10 - 2)
			expect 15  (10 - -5)
			expect 2   (-3 - -5)
			expect -32 -(32)
			expect 32  +(32)

			expect 8     (1 << 3)
			expect 32    byte(0x11 << 5)
			expect 0     byte(0xd0 << 65)
			expect -0x1a (-0xcd >> 3)
			expect 0     (0xaa >> 10)
			expect 0     (0xd0 >> 65)
			expect 0x0f  (-0xcd >>> 60)
			expect 0     (0xaa >>> 10)
			expect 0     (0xd0 >>> 65)
			expect 0xff  byte(~0)
			expect 0xfe  byte(~1)
			expect 0     ~-1

			expect 1 (1 > 0)
			expect 0 (1 > 1)
			expect 1 (2 > 1)
			expect 1 (1 >= 0)
			expect 0 (0 >= 1)
			expect 1 (2 >= 2)
			expect 0 (1 < 0)
			expect 0 (1 < 1)
			expect 1 (1 < 2)
			expect 0 (1 <= 0)
			expect 1 (1 <= 1)
			expect 0 (2 <= 1)
			expect 1 (1 == 1)
			expect 0 (1 == 2)
			expect 1 ("a" == "a")
			expect 0 ("a" == "b")
			expect 0 (1 != 1)
			expect 1 (1 != 2)
			expect 0 ("a" != "a")
			expect 1 ("a" != "b")
			expect 1 !0
			expect 0 !1
			expect 0 !-1

			expect 2 (1 && 2)
			expect 0 (0 && 2)
			expect 1 (1 || 2)
			expect 2 (0 || 2)
			expect 2 (3 & 2)
			expect 7 (3 | 4)

			expect 0x66 byte(0xab ^ 0xcd)
			expect 0xef byte(~(0x10))
			expect 0xff byte(~(0))
			expect 0x5f byte(~(-0x60))
			expect 0xcd lobyte(0xabcd)
			expect 0xab hibyte(0xabcd)
			expect 0xabcd asword(0xab 0xcd)
		`)
		tt.Eq(t, "ok", string(dat))
	})

	t.Run("ok: expand-binary", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			$(0xFF) -byte; $(0xFE) -byte; $(0x01) -rep 6
		`)
		tt.EqSlice(t, []byte{0xff, 0xfe, 1, 1, 1, 1, 1, 1}, dat)
	})

	t.Run("ok: typecasts", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			dw byte(-2) byte(-1) byte(0) byte(1) byte(2) byte(255)
			dw byte(-0xffff) byte(-0x7ffe) byte(-0x7fff) byte(0x1ff) byte(0xffff)
			dw word(-2) word(-1) word(0) word(1) word(2) word(255)
			dw word(-0x7ffe) word(-0x7fff) word(0x1ff) word(0xffff)
		`)
		tt.EqSlice(t, []byte{
			254, +0, 255, +0, 0, +0, 1, +0, 2, +0, 255, +0,
			1, +0, 2, +0, 1, +0, 255, +0, 255, +0,
			254, 255, 255, 255, 0, +0, 1, +0, 2, +0, 255, +0,
			2, 128, 1, 128, 255, 1, 255, 255,
		}, dat)
	})

	t.Run("ok: functions", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			link-as-tests
			module ModA { const c001 = 1 }
			module ModB { const c002 = 2 }
			data d001 = byte * 5 : bss
			expect (5 == sizeof(d001))
			expect ("d001" == nameof(d001))
			expect defined?(d001)
			expect defined?(ModA:c001)
			expect !defined?(ModB:c001)
			expect !defined?(ModA:c002)
			expect !defined?(ModC:c002)
		`)
		tt.Eq(t, "ok", string(dat))
	})

	t.Run("ok: others", func(t *testing.T) {
		dat := expectCompileOk(t, `
			link { org 0 0 1; merge text ModA ModB _ }
			module ModA { const c001 = 1 }
			module ModB { import ModA; const c002 = 2 }
			macro m(a ...) {}
			macro m001() ={ macro m002(a) { db %%=a } }
			macro isform(a b) { if (formtypeof(%&a) != %=b) { compile-error %=b } }
			macro isexpr(a b) { if (exprtypeof(%=a) != %=b) { compile-error %=b } }
			section text
			proc f001(!) {
				data byte ["long long string"]
				expand-loop (1 + 2) { db 0xab}
				NOOP
				X <- 1
				AB <- 1 : 2
				A . { db 1 }
				A . f001(!)
				A . B@1
				A <- B@{m}
				[0x1234] -dnnm
				@1
				JMP f001 ==?
				m [+ X] [X +] (1 / 1) -10 0b1010 '\0' "ok\0" {
					A
					  + 1
					  . { db 1 }
					$(1)
					  -byte
				}
				m001; m002 1
				L0: BRL L0
				db 'o' 'k'
				isform A "reg"
				isform EQ? "cond"
				isform 1 "constexpr"
				isform { do } "block-form"
				isform [1] "mem-form"
				isform A@1 "unknown"
				isexpr 1 "int"
				isexpr "str" "str"
				isexpr id "identifier"
				isexpr (1 + 1) "unknown"
				RET
			}
		`)
		tt.True(t, len(dat) > 0)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"unknown namespace nothing", `flat!
				import nothing
			`,
			"qualified name is not allowed in this context", `flat!
				import ns:nothing
			`,
			"loop counter must be constexpr", `flat!
				expand-loop A { NOP }
			`,
			"loop counter must be integer", `flat!
				expand-loop "A" { NOP }
			`,
			"this is compile-error", `flat!
				compile-error "this is compile-error"
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileSpecials(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			db $$(__FILE__) $$(loaded-as-main?)
			db __ORG__ __PC__
			proc fn() {
				dw __PROC__
				RET
			}
		`)
		tt.EqSlice(t, []byte{'-', 1, 0, 3, 4, 0, 0x04}, dat)
	})

	t.Run("ok: __FILE__", func(t *testing.T) {
		g := ttarch.BuildGenerator("ttarch", `flat!
			db $$(__FILE__)
		`)
		dat, _, mes := ttarch.DoCompile(g, "test.oc")
		tt.EqSlice(t, []byte("./test.oc"), dat, mes)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"cannot use __FILE__ in link phase", `flat!
				db __FILE__
			`,
			"cannot use loaded-as-main? in link phase", `flat!
				db loaded-as-main?
			`,
			"cannot use __ORG__ in compile phase", `flat!
				db $$(__ORG__)
			`,
			"cannot use __PC__ in compile phase", `flat!
				db $$(__PC__)
			`,
			"undefined name __PROC__", `flat!
				db $$(__PROC__)
			`,
			"cannot use label address in compile phase", `flat!
				proc f() { db $$(__PROC__) }
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileBlock(t *testing.T) {
	t.Run("ok: do", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			L01: JMP L01
		`)
		dat := expectCompileOk(t, `flat!
			do ={ L01: }
			JMP L01
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("error: prog", func(t *testing.T) {
		mes := expectCompileError(t, `flat!
			do { L01: }
			JMP L01
		`)
		tt.Eq(t, "compile error: undefined name L01", mes)
	})
}

func TestCompileLoop(t *testing.T) {
	t.Run("ok: loop", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			L0: JMP L0 ; JMP L2 ; JMP L1
			L1: JMP L0
			L2:
		`)
		dat := expectCompileOk(t, `flat!
			loop {
				JMP _BEG
				JMP _END
				JMP _COND
			}
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("ok: loop cond", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			L0: JMP L1
			L1:
		`)
		dat := expectCompileOk(t, `flat!
			loop {
			} JMP _END
		`)
		tt.EqSlice(t, expected, dat)
	})
}

func TestCompileApply(t *testing.T) {
	t.Run("ok: apply", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			apply do { NOP }
		`)
		tt.EqSlice(t, []byte{0x00}, dat)
	})
}

func TestCompileIf(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			if (1 == 2) { db 0 }
			if (1 == 1) { db 0 }
			if (1 == 2) { db 0 } else { db 1 }
			if (1 == 2) { db 0 } else if (1 == 1) { db 2 }
			if (1 == 2) { db 0 } else if (1 == 3) { db 1 } else { db 3 }
		`)
		tt.EqSlice(t, []byte{0, 1, 2, 3}, dat)
	})

	t.Run("ok: if-cond", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			if EQ? { db 0 }
			if EQ? { db 0 } else { db 1 }
		`)
		dat := expectCompileOk(t, `flat!
			l00: $(l02) -jump-unless EQ?
			l01: db 0
			l02:

			l10: $(l12) -jump-unless EQ?
			l11: db 0; JMP $(l13)
			l12: db 1
			l13:
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"then-body must be block-form", `flat!
				if (1 == 1) A
			`,
			"invalid if form, expected `else`", `flat!
				if (1 == 1) { db 0 } ell
			`,
			"invalid if form, else body required", `flat!
				if (1 == 1) { db 0 } else
			`,
			"else-body must be block-form", `flat!
				if (1 == 1) { db 0 } else iff
			`,
			"invalid else-if form", `flat!
				if (1 == 1) { db 0 } else if
			`,
			"cond must be constexpr", `flat!
				if (1 == 1) { db 0 } else if EQ? {}
			`,
			"then-body must be block-form", `flat!
				if (1 == 1) { db 0 } else if 1 B
			`,
			"invalid if-else form", `flat!
				if 1 {} else {} 1
			`,
			"invalid condition", `flat!
				if A { db 0 }
			`,
			"invalid if form, expected `else`", `flat!
				if EQ? {} else_ {}
			`,
			"invalid if form, else body required", `flat!
				if EQ? {} else
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileCase(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			case 1 when 2 { db 0 }
			case 1 when 1 { db 0 }
			case 1 when 2 { db 0 } else { db 1 }
			case 1 when 2 { db 0 } when 1 { db 2 }
			case 1 when 2 { db 0 } when 3 { db 1 } else { db 3 }
		`)
		tt.EqSlice(t, []byte{0, 1, 2, 3}, dat)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"case: 3+ argument(s) required, but given 1", `flat!
				case 1
			`,
			"value must be constexpr", `flat!
				case A when 1 {}
			`,
			"invalid case-when form", `flat!
				case 1 when 1 {} when
			`,
			"pattern must be constexpr", `flat!
				case 1 when EQ? {}
			`,
			"when-body must be block-form", `flat!
				case 1 when 1 A
			`,
			"invalid case form. when/else required", `flat!
				case 1 when 1 { db 0 } ell
			`,
			"invalid case-else form", `flat!
				case 1 when 1 { db 0 } else
			`,
			"else-body must be block-form", `flat!
				case 1 when 1 { db 0 } else A
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileWhen(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			when { db 0; &&- EQ? },
			then { db 1 }

			when { db 0; ||- EQ?; &&- CC? },
			then { db 1 }

			when { db 0; &&- EQ? },
			then { db 1 },
			else { db 2 }
		`)
		dat := expectCompileOk(t, `flat!
			l00: db 0; $(l02) -jump-unless EQ?
			l01: db 1
			l02:

			l10: db 0; $(l11) -jump-if EQ?; $(l12) -jump-unless CC?
			l11: db 1
			l12:

			l20: db 0; $(l22) -jump-unless EQ?
			l21: db 1; JMP $(l23)
			l22: db 2
			l23:
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"invalid when form, expected `then`", `flat!
				when { &&- EQ? } then_ {}
			`,
			"invalid when form, expected `else`", `flat!
				when { &&- EQ? } then {} else_ {}
			`,
			"invalid when form, else body required", `flat!
				when { &&- EQ? } then {} else
			`,
			"at least one condition required", `flat!
				when {} then {}
			`,
			"the last condition must be an AND expression", `flat!
				when { ||- EQ? } then {}
			`,
			"unknown form name &&-", `flat!
				when { &&- EQ? } then { &&- EQ? }
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompilePragma(t *testing.T) {
	t.Run("ok: list-constants", func(t *testing.T) {
		list := expectGenListOk(t, `
				module Mod {
					const a = 1
					pragma list-constants 0
					const b = 2
					pragma list-constants 1
					const c = 3
					pragma list-constants 0 "comment"
					const d = 4
					pragma list-constants 1 "comment"
					const e = 5
				}
		`)
		tt.EqText(t, tt.Unindent(`
			|                                            ; generated by ocala
			                                            __ARCH__ = "ttarch"
			                                            Mod:a = 1
			                                            Mod:c = 3
			                                            ; comment
			                                            Mod:e = 5

			     - 0000                                 .org 0
		`)[1:], list)
	})

	t.Run("ok: comment", func(t *testing.T) {
		list := expectGenListOk(t, `
				pragma comment "comment"
				pragma comment "value" 1
		`)
		tt.EqText(t, tt.Unindent(`
			|                                            ; generated by ocala
			                                            __ARCH__ = "ttarch"

			     - 0000                                 .org 0
			                                            ; comment
			                                            ; value 1
		`)[1:], list)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"unknown pragma: invalid", `flat!
				pragma invalid
			`,
			"qualified name is not allowed in this context", `flat!
				pragma invalit:comment
			`,
			"pragma is not allowed in this context", `flat!
				pragma list-constants 1
			`,
			"value must be constexpr", `flat!
				module Mod { pragma list-constants A }
			`,
			"value must be integer", `flat!
				module Mod { pragma list-constants "1" }
			`,
			"message must be constexpr", `flat!
				module Mod { pragma list-constants 1 A }
			`,
			"message must be string", `flat!
				module Mod { pragma list-constants 1 1 }
			`,
			"comment must be constexpr", `flat!
				pragma comment A
			`,
			"comment must be string", `flat!
				pragma comment 1
			`,
			"comment element must be constexpr", `flat!
				pragma comment "message" A
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileTco(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			proc f001(!) { RET }
			JMP f001
		`)
		dat := expectCompileOk(t, `flat!
			proc f001(!) { RET }
			tco { f001(!) }
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"only one statement is allowed within the tco form", `flat!
				proc f001(!) { RET }
				tco { f001(!); f001(!) }
			`,
			"proc call required", `flat!
				proc f001(!) { RET }
				tco { @1 }
			`,
			"the conditional call is not a tail call", `flat!
				proc f001(!) { RET }
				tco { EQ?.f001(!) }
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompilePatch(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			data d001 = byte @ <reserved>
			LD [d001] A@10
			A <- 0; *patch* d001 byte
		`)
		dat := expectCompileOk(t, `flat!
			LD [l001 + 1] A@10
			l001: A <- 0
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("ok: word", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			data d001 = word @ <reserved>
			LD [d001] A@10
			AB <- 0; *patch* d001 word
		`)
		dat := expectCompileOk(t, `flat!
			LD [(l001 + 1)] A@10
			l001: AB <- 0
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("ok: index", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			data d001 = word @ <reserved>
			LD [d001] A@10
			AB <- 0; *patch* d001 -2
		`)
		dat := expectCompileOk(t, `flat!
			LD [(l001 + 1)] A@10
			l001: AB <- 0
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("ok: proc", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			proc f001() @ <reserved>
			proc f002() {
				NOP;
				*patch* f001
				RET;
			}
			f001(); f002()
		`)
		dat := expectCompileOk(t, `flat!
			l002: NOP
			l001: RET
			JSR l001; JSR l002;
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("ok: module", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			module ModA {
				proc f() { A <- [ModB:d001]; RET }
			}
			module ModB {
				data d001 = byte @ <reserved>
				proc f() { NOP; *patch* d001 byte; RET }
			}
		`)
		tt.True(t, len(dat) > 0)
	})

	t.Run("ok: generate list", func(t *testing.T) {
		list := expectGenListOk(t, `flat!
			data d001 = word @ <reserved>
			LD [d001] A@10
			AB <- 0; *patch* d001 word
		`)
		tt.EqText(t, tt.Unindent(`
			|                                            ; generated by ocala
			                                            __ARCH__ = "ttarch"
			                                            d001 = patch\.d001\.G\d+ \+ -2 ~

			     - 0000                                 .org 0
			000000 0000[2] 23 0a                        LD     A, 0+ 10
			000002 0002[3] 34 06 00                     LD     (d001), A
			000005 0005[3] 3b 00 00                     LD     AB, 0+ 0
			     - 0008                             patch\.d001\.G\d+: ~
		`)[1:], list)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"unknown label: invalid", `flat!
				*patch* invalid
			`,
			"d001 already defined", `flat!
				data d001 = byte @ <reserved>
				A <- 0; *patch* d001 byte
				A <- 0; *patch* d001 byte
			`,
			"reserved but undefined", `flat!
				const d001 = <reserved>
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileMakeCounter(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			make-counter c 0
			L1: db c c
		`)
		tt.EqSlice(t, []byte{0, 1}, dat)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"start must be constexpr", `flat!
				make-counter c A
			`,
			"start must be integer", `flat!
				make-counter c "A"
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileAssert(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			beg: NOP; end:
			assert (end - beg == 1)
		`)
		tt.Eq(t, 1, len(dat))
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"assertion failed", `flat!
				beg: NOP; end:
				assert (end - beg != 1)
			`,
			"error message", `flat!
				beg: NOP; end:
				assert (end - beg != 1) "error message"
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileSizeof(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dat := expectCompileOk(t, `flat!
			link-as-tests
			data d001 = byte * 4 : bss
			expect 4 sizeof(d001)

			data d002 = byte [0 1 2 3 4 5] : bss
			expect 6 sizeof(d002)

			data d003 = byte load-file("./testdata/embed01.dat") : bss
			expect 16 sizeof(d003)

			struct s001 { x byte; y word }
			expect 3 sizeof(s001)
			expect 1 sizeof(s001.x)
			expect 2 sizeof(s001.y)

			struct s002 { x [4]byte; y [4]word; z [4]s001; a s001 }
			expect 27 sizeof(s002)
			expect 4  sizeof(s002.x)
			expect 8  sizeof(s002.y)
			expect 12 sizeof(s002.z)
			expect 1  sizeof(s002.z.x)
			expect 2  sizeof(s002.z.y)
			expect 3  sizeof(s002.a)
			expect 1  sizeof(s002.a.x)
			expect 2  sizeof(s002.a.y)

			struct s003 { x [4][8]byte; y [4][8]word }
			expect 96 sizeof(s003)
			expect 32 sizeof(s003.x)
			expect 64 sizeof(s003.y)

			data d004 = s001 {} : bss
			expect 1 sizeof(d004.x)
			expect 2 sizeof(d004.y)

			data d005 = s002 {} : bss
			expect 27 sizeof(d005)
			expect 4  sizeof(d005.x)
			expect 8  sizeof(d005.y)
			expect 12 sizeof(d005.z)
			expect 1  sizeof(d005.z.x)
			expect 2  sizeof(d005.z.y)
			expect 3  sizeof(d005.a)
			expect 1  sizeof(d005.a.x)
			expect 2  sizeof(d005.a.y)

			data d006 = s003 {} : bss
			expect 96 sizeof(d006)
			expect 32 sizeof(d006.x)
			expect 64 sizeof(d006.y)
		`)
		tt.Eq(t, "ok", string(dat))
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"unknown label name L1", `flat!
				db sizeof(L1)
			`,
			"sizeof requires a data label or struct", `flat!
				const a = 0
				db sizeof(a)
			`,
			"L1 is not a data label", `flat!
				L1: db sizeof(L1)
			`,
			"d001 is not a struct type", `flat!
				data d001 = byte * 10
				db sizeof(d001.a)
			`,
			"byte is not a struct type", `flat!
				data d001 = byte * 10
				db sizeof(byte.a)
			`,
			"unknown field b", `flat!
				data d001 = struct { a byte }
				db sizeof(d001.b)
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileNametypeof(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		dat := expectCompileOk(t, `
			link-as-tests
			module ModA { const c002 = 2 }
			macro m001() {}
			const c001 = 1
			proc f001(!) {
				data d001 = byte [0]
				L001:
				expect "module"  nametypeof(ModA)
				expect "syntax"  nametypeof(if)
				expect "macro"   nametypeof(m001)
				expect "func"    nametypeof(+)
				expect "inst"    nametypeof(LD)
				expect "const"   nametypeof(c001)
				expect "label"   nametypeof(d001)
				expect "label"   nametypeof(f001)
				expect "label"   nametypeof(L001)
				expect "special" nametypeof(__PC__)
				expect "const"   nametypeof(ModA:c002)
				RET
			}
		`)
		tt.Eq(t, "ok", string(dat))
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"unknown name unknown", `flat!
				db nametypeof(unknown)
			`,
			"unknown name unknown:unknown", `flat!
				db nametypeof(unknown:unknown)
			`,
			"unknown name ModA:unknown", `flat!
				module ModA {}
				db nametypeof(ModA:unknown)
			`,
			"unknown name unknown:c001", `flat!
				const c001 = 1
				db nametypeof(unknown:c001)
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}

func TestCompileDebugInspect(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := ttarch.BuildGenerator("ttarch", `flat!
			debug-inspect "hello"
			debug-inspect EQ?
			debug-inspect quote(LD)
		`)
		ttarch.DoCompile(g, "-")
		tt.Eq(t, "[DEBUG] \"hello\"\n"+
			"[DEBUG] EQ?\n"+
			"[DEBUG] LD\n", tt.FlushString(g.ErrWriter))
	})
}

func TestCompileOpcode(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expected := expectCompileOk(t, `flat!
			data byte [0]
			data byte [0]
			NOP; NOP; NOP
			LD A B; RET EQ?
			L1: BRA L1
			JMP 0
		`)
		dat := expectCompileOk(t, `flat!
			data byte [opcode("LD X 0x8000" 0)]
			data byte [opcode("BRA 0x8000" 0)]
			data byte [opcode("NOP") opcode("NOP" 1) opcode("NOP" 2)]
			data byte [opcode("LD A B") opcode("RET EQ?")]
			const b = opcode("BRA 0" 2)
			data byte [hibyte(b) lobyte(b)]
			const c = opcode("JMP 0" 3)
			data byte [((c >>> 16) & 0xff) ((c >>> 8) & 0xff) (c & 0xff)]
		`)
		tt.EqSlice(t, expected, dat)
	})

	t.Run("error", func(t *testing.T) {
		es := []string{
			"unknown mnemonic", `flat!
				db opcode("UNKOWN A")
			`,
			"unknown mnemonic", `flat!
				db opcode("x:LD A B")
			`,
			"unknown mnemonic", `flat!
				db opcode("L: BRA L")
			`,
			"unknown mnemonic", `flat!
				db opcode("const a = 1")
			`,
			"#.const is not allowed in this context", `flat!
				db opcode("LD { const a = 1 }")
			`,
			"invalid operands for 'LD'", `flat!
				db opcode("LD A 0x8000")
			`,
			"invalid operands for 'JMP'", `flat!
				db opcode("JMP A")
			`,
			"the relative address is out of range", `flat!
				db opcode("BRA 0x8000" 2)
			`,
		}
		for x := 0; x < len(es); x += 2 {
			mes := expectCompileError(t, es[x+1])
			tt.Eq(t, "compile error: "+es[x], mes, es[x+1])
		}
	})
}
