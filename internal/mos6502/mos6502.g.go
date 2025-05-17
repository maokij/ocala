package mos6502

import . "ocala/internal/core" //lint:ignore ST1001 core

var kwRegA = Intern("A")
var kwRegX = Intern("X")
var kwRegY = Intern("Y")
var kwRegS = Intern("S")
var kwRegP = Intern("P")
var kwImmN = Intern("%B")
var kwImmNN = Intern("%W")
var kwMemZN = Intern("[%B]")
var kwMemZX = Intern("[%B X]")
var kwMemZY = Intern("[%B Y]")
var kwMemAN = Intern("[%W]")
var kwMemAX = Intern("[%W X]")
var kwMemAY = Intern("[%W Y]")
var kwMemIN = Intern("[[%W]]")
var kwMemIX = Intern("[[%B X]]")
var kwMemIY = Intern("[[%B] Y]")
var kwCondNE = Intern("NE?")
var kwCondEQ = Intern("EQ?")
var kwCondCC = Intern("CC?")
var kwCondCS = Intern("CS?")
var kwCondVC = Intern("VC?")
var kwCondVS = Intern("VS?")
var kwCondPL = Intern("PL?")
var kwCondMI = Intern("MI?")

var operandToAsmMap = map[*Keyword](struct {
	s string
	t bool
}){
	kwRegA:   {s: "A", t: false},
	kwRegX:   {s: "X", t: false},
	kwRegY:   {s: "Y", t: false},
	kwRegS:   {s: "S", t: false},
	kwRegP:   {s: "P", t: false},
	kwImmN:   {s: "#%", t: true},
	kwImmNN:  {s: "#%", t: true},
	kwMemZN:  {s: "%", t: true},
	kwMemZX:  {s: "%, X", t: true},
	kwMemZY:  {s: "%, Y", t: true},
	kwMemAN:  {s: "%", t: true},
	kwMemAX:  {s: "%, X", t: true},
	kwMemAY:  {s: "%, Y", t: true},
	kwMemIN:  {s: "(%)", t: true},
	kwMemIX:  {s: "(%, X)", t: true},
	kwMemIY:  {s: "(%), Y", t: true},
	kwCondNE: {s: "NE", t: false},
	kwCondEQ: {s: "EQ", t: false},
	kwCondCC: {s: "CC", t: false},
	kwCondCS: {s: "CS", t: false},
	kwCondVC: {s: "VC", t: false},
	kwCondVS: {s: "VS", t: false},
	kwCondPL: {s: "PL", t: false},
	kwCondMI: {s: "MI", t: false},
}

var tokenWords = [][]string{
	{"A", "X", "Y", "S", "P"},
	{"NE?", "EQ?", "CC?", "CS?", "VC?", "VS?", "PL?", "MI?"},
	{"-push", "-pop", "++", "--", "-not", "-neg"},
	{"<-", "->", "+$", "-$", "-?", "-bit?", "<*", "<*$", ">*", ">*$", "-jump-if", "-jump-unless"},
}

var oppositeConds = map[*Keyword]*Keyword{
	kwCondNE: kwCondEQ,
	kwCondEQ: kwCondNE,
	kwCondCC: kwCondCS,
	kwCondCS: kwCondCC,
	kwCondVC: kwCondVS,
	kwCondVS: kwCondVC,
	kwCondPL: kwCondMI,
	kwCondMI: kwCondPL,
}

var bmaps = [][]byte{}

var tokenAliases = map[string]string{
	"!=?":         "NE?",
	"not-zero?":   "NE?",
	"==?":         "EQ?",
	"zero?":       "EQ?",
	"<?":          "CC?",
	"not-carry?":  "CC?",
	"borrow?":     "CC?",
	">=?":         "CS?",
	"carry?":      "CS?",
	"not-borrow?": "CS?",
	"not-over?":   "VC?",
	"over?":       "VS?",
	"plus?":       "PL?",
	"minus?":      "MI?",
}

var instAliases = map[string][]string{}

var instMap = InstPat{
	Intern("LDA"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa9},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xad},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("LDX"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa2},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xae},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbe},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("LDY"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa0},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xac},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbc},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("STA"): InstPat{
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x85},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x95},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x8d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x9d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x99},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x81},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x91},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x99},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("STX"): InstPat{
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x86},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x96},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x8e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("STY"): InstPat{
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x84},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x94},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x8c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("TAX"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xaa},
		},
	},
	Intern("TAY"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xa8},
		},
	},
	Intern("TSX"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xba},
		},
	},
	Intern("TXA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x8a},
		},
	},
	Intern("TXS"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x9a},
		},
	},
	Intern("TYA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x98},
		},
	},
	Intern("PHA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x48},
		},
	},
	Intern("PHP"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x08},
		},
	},
	Intern("PLP"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x28},
		},
	},
	Intern("PLA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x68},
		},
	},
	Intern("CLC"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x18},
		},
	},
	Intern("CLI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x58},
		},
	},
	Intern("CLD"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xd8},
		},
	},
	Intern("CLV"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xb8},
		},
	},
	Intern("SEC"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x38},
		},
	},
	Intern("SEI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x78},
		},
	},
	Intern("SED"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xf8},
		},
	},
	Intern("BRK"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x00},
		},
	},
	Intern("NOP"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xea},
		},
	},
	Intern("RTS"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x60},
		},
	},
	Intern("RTI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x40},
		},
	},
	Intern("JMP"): InstPat{
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x6c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("JSR"): InstPat{
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x20},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x20},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("BPL"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x10},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x10},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BMI"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x30},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x30},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BVC"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x50},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x50},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BVS"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x70},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x70},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BCC"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x90},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x90},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BCS"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BNE"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BEQ"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("ORA"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x09},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x05},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x15},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x0d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x1d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x19},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x11},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x19},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("AND"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x29},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x25},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x35},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x2d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x3d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x39},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x21},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x31},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x39},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("EOR"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x49},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x45},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x55},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x4d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x5d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x59},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x41},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x51},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x59},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("ADC"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x69},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x65},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x75},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x6d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x7d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x79},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x61},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x71},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x79},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("CMP"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc9},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("SBC"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe9},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xed},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("BIT"): InstPat{
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x24},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x2c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("CPX"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe0},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xec},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("CPY"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc0},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcc},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("INC"): InstPat{
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xee},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfe},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("INX"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xe8},
		},
	},
	Intern("INY"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xc8},
		},
	},
	Intern("DEC"): InstPat{
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xce},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xde},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("DEX"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xca},
		},
	},
	Intern("DEY"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x88},
		},
	},
	Intern("ASL"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x0a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x06},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x16},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x0e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x1e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("LSR"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x4a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x46},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x56},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x4e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x5e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("ROL"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x2a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x26},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x36},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x2e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x3e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("ROR"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x6a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x66},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x76},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x6e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x7e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
}

var ctxOpMap = map[*Keyword]map[*Keyword]map[*Keyword][][]Value{
	Intern("<-"): {
		kwRegA: {
			kwRegX: {
				{Intern("TXA")},
			},
			kwRegY: {
				{Intern("TYA")},
			},
			KwAny: {
				{Intern("LDA"), &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			kwRegA: {
				{Intern("TAX")},
			},
			kwRegS: {
				{Intern("TSX")},
			},
			KwAny: {
				{Intern("LDX"), &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			kwRegA: {
				{Intern("TAY")},
			},
			KwAny: {
				{Intern("LDY"), &Vec{Int(1), nil}},
			},
		},
		kwRegS: {
			kwRegX: {
				{Intern("TXS")},
			},
		},
	},
	Intern("->"): {
		kwRegA: {
			kwRegX: {
				{Intern("TAX")},
			},
			kwRegY: {
				{Intern("TAY")},
			},
			KwAny: {
				{Intern("STA"), &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			kwRegA: {
				{Intern("TXA")},
			},
			kwRegS: {
				{Intern("TXS")},
			},
			KwAny: {
				{Intern("STX"), &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			kwRegA: {
				{Intern("TYA")},
			},
			KwAny: {
				{Intern("STY"), &Vec{Int(1), nil}},
			},
		},
		kwRegS: {
			kwRegX: {
				{Intern("TSX")},
			},
		},
	},
	Intern("-push"): {
		kwRegA: {
			nil: {
				{Intern("PHA")},
			},
		},
		kwRegP: {
			nil: {
				{Intern("PHP")},
			},
		},
	},
	Intern("-pop"): {
		kwRegA: {
			nil: {
				{Intern("PLA")},
			},
		},
		kwRegP: {
			nil: {
				{Intern("PLP")},
			},
		},
	},
	Intern("++"): {
		kwRegA: {
			nil: {
				{Intern("CLC")},
				{Intern("ADC"), Int(1)},
			},
		},
		kwRegX: {
			nil: {
				{Intern("INX")},
			},
		},
		kwRegY: {
			nil: {
				{Intern("INY")},
			},
		},
		KwAny: {
			nil: {
				{Intern("INC"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("--"): {
		kwRegA: {
			nil: {
				{Intern("SEC")},
				{Intern("SBC"), Int(1)},
			},
		},
		kwRegX: {
			nil: {
				{Intern("DEX")},
			},
		},
		kwRegY: {
			nil: {
				{Intern("DEY")},
			},
		},
		KwAny: {
			nil: {
				{Intern("DEC"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-not"): {
		kwRegA: {
			nil: {
				{Intern("EOR"), Int(255)},
			},
		},
	},
	Intern("-neg"): {
		kwRegA: {
			nil: {
				{Intern("EOR"), Int(255)},
				{Intern("CLC")},
				{Intern("ADC"), Int(1)},
			},
		},
	},
	Intern("+"): {
		kwRegA: {
			KwAny: {
				{Intern("CLC")},
				{Intern("ADC"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("+$"): {
		kwRegA: {
			KwAny: {
				{Intern("ADC"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-"): {
		kwRegA: {
			KwAny: {
				{Intern("SEC")},
				{Intern("SBC"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-$"): {
		kwRegA: {
			KwAny: {
				{Intern("SBC"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-?"): {
		kwRegA: {
			KwAny: {
				{Intern("CMP"), &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			KwAny: {
				{Intern("CPX"), &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			KwAny: {
				{Intern("CPY"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("&"): {
		kwRegA: {
			KwAny: {
				{Intern("AND"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("|"): {
		kwRegA: {
			KwAny: {
				{Intern("ORA"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("^"): {
		kwRegA: {
			KwAny: {
				{Intern("EOR"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-bit?"): {
		kwRegA: {
			KwAny: {
				{Intern("BIT"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("<*"): {
		kwRegA: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("CMP"), Int(128)},
					&Vec{Intern("ROL"), &Operand{Kind: kwRegA}},
				}},
			},
		},
	},
	Intern("<*$"): {
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("ROL"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">*"): {
		kwRegA: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("LSR"), &Operand{Kind: kwRegA}},
					&Vec{Intern("#.BYTE"), Int(144)},
					&Vec{Intern("#.BYTE"), Int(2)},
					&Vec{Intern("ORA"), Int(128)},
				}},
			},
		},
	},
	Intern(">*$"): {
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("ROR"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("<<"): {
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("ASL"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">>"): {
		kwRegA: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("CMP"), Int(128)},
					&Vec{Intern("ROR"), &Operand{Kind: kwRegA}},
				}},
			},
		},
	},
	Intern(">>>"): {
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("LSR"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("-jump-if"): {
		kwImmNN: {
			kwCondNE: {
				{Intern("BNE"), &Vec{Int(0), nil}},
			},
			kwCondEQ: {
				{Intern("BEQ"), &Vec{Int(0), nil}},
			},
			kwCondCC: {
				{Intern("BCC"), &Vec{Int(0), nil}},
			},
			kwCondCS: {
				{Intern("BCS"), &Vec{Int(0), nil}},
			},
			kwCondVC: {
				{Intern("BVC"), &Vec{Int(0), nil}},
			},
			kwCondVS: {
				{Intern("BVS"), &Vec{Int(0), nil}},
			},
			kwCondPL: {
				{Intern("BPL"), &Vec{Int(0), nil}},
			},
			kwCondMI: {
				{Intern("BMI"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-jump-unless"): {
		kwImmNN: {
			kwCondNE: {
				{Intern("BEQ"), &Vec{Int(0), nil}},
			},
			kwCondEQ: {
				{Intern("BNE"), &Vec{Int(0), nil}},
			},
			kwCondCC: {
				{Intern("BCS"), &Vec{Int(0), nil}},
			},
			kwCondCS: {
				{Intern("BCC"), &Vec{Int(0), nil}},
			},
			kwCondVC: {
				{Intern("BVS"), &Vec{Int(0), nil}},
			},
			kwCondVS: {
				{Intern("BVC"), &Vec{Int(0), nil}},
			},
			kwCondPL: {
				{Intern("BMI"), &Vec{Int(0), nil}},
			},
			kwCondMI: {
				{Intern("BPL"), &Vec{Int(0), nil}},
			},
		},
	},
}
