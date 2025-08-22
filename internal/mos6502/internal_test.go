package mos6502

import (
	"ocala/internal/core"
	"ocala/internal/tt"
	"testing"
)

func TestAdjustOperand(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		id := &core.Identifier{}
		cc := BuildCompiler()
		data := []struct {
			n        int
			from, to *core.Keyword
		}{
			{255, kwImmN, kwImmN},
			{256, kwImmN, kwImmNN},
			{255, kwMemZN, kwMemZN},
			{256, kwMemZN, kwMemAN},
			{255, kwMemZX, kwMemZX},
			{256, kwMemZX, kwMemAX},
			{255, kwMemZY, kwMemZY},
			{256, kwMemZY, kwMemAY},
		}
		for x, i := range data {
			a := &core.Operand{Kind: i.from}
			cc.AdjustOperand(cc, a, i.n, id)
			tt.Eq(t, i.to, a.Kind, x)
		}
	})
}
