package z80

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
			{255, kwImmNN, kwImmN},
			{256, kwImmNN, kwImmNN},
			{255, kwMemN, kwMemN},
			{256, kwMemN, kwMemNN},
			{255, kwMemNN, kwMemN},
			{256, kwMemNN, kwMemNN},
			{255, kwMemIX, kwMemIX},
			{256, kwMemIX, kwMemWX},
			{255, kwMemWX, kwMemIX},
			{256, kwMemWX, kwMemWX},
			{255, kwMemIY, kwMemIY},
			{256, kwMemIY, kwMemWY},
			{255, kwMemWY, kwMemIY},
			{256, kwMemWY, kwMemWY},
		}
		for x, i := range data {
			a := &core.Operand{Kind: i.from}
			cc.AdjustOperand(cc, a, i.n, id)
			tt.Eq(t, i.to, a.Kind, x)
		}
	})
}
