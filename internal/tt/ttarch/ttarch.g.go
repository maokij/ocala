package ttarch

import . "ocala/internal/core" //lint:ignore ST1001 core

var kwRegA = Intern("A")
var kwRegB = Intern("B")
var kwRegP = Intern("P")
var kwRegAB = Intern("AB")
var kwRegSP = Intern("SP")
var kwRegPQ = Intern("PQ")
var kwRegX = Intern("X")
var kwMemX = Intern("[X]")
var kwRegY = Intern("Y")
var kwMemY = Intern("[Y]")
var kwImmN = Intern("%B")
var kwImmNN = Intern("%W")
var kwMemNN = Intern("[%W]")
var kwCondNE = Intern("NE?")
var kwCondEQ = Intern("EQ?")
var kwCondCC = Intern("CC?")
var kwCondCS = Intern("CS?")

var operandToAsmMap = map[*Keyword](struct {
	s string
	t bool
}){
	kwRegA:   {s: "A", t: false},
	kwRegB:   {s: "B", t: false},
	kwRegP:   {s: "P", t: false},
	kwRegAB:  {s: "AB", t: false},
	kwRegSP:  {s: "SP", t: false},
	kwRegPQ:  {s: "PQ", t: false},
	kwRegX:   {s: "X", t: false},
	kwMemX:   {s: "(X)", t: false},
	kwRegY:   {s: "Y", t: false},
	kwMemY:   {s: "(Y)", t: false},
	kwImmN:   {s: "0+ %", t: true},
	kwImmNN:  {s: "0+ %", t: true},
	kwMemNN:  {s: "(%)", t: true},
	kwCondNE: {s: "NE", t: false},
	kwCondEQ: {s: "EQ", t: false},
	kwCondCC: {s: "CC", t: false},
	kwCondCS: {s: "CS", t: false},
}

var tokenWords = [][]string{
	{"A", "B", "P", "X", "Y", "AB", "SP"},
	{"NE?", "EQ?", "CC?", "CS?"},
	{"-dnnm", "-byte"},
	{"<-", "-jump-if", "-jump-unless", "-rep"},
}

var bmaps = [][]byte{
	{3, 1, 3, 0, 97, 103, 111}, // 0: BMM
}

var tokenAliases = map[string]string{
	"!=?": "NE?",
	"==?": "EQ?",
	">=?": "CC?",
	"<?":  "CS?",
}

var instAliases = map[string][]string{
	"NOP": {"NOOP"},
}

var instMap = InstPat{
	Intern("NOP"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x00},
		},
	},
	Intern("JMP"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("JPR"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x02},
				{Kind: BcRlow, A0: 0x00, A1: 0xfd, A2: 0x02},
				{Kind: BcRhigh, A0: 0x00, A1: 0xfd, A2: 0x02},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x02},
				{Kind: BcRlow, A0: 0x00, A1: 0xfd, A2: 0x02},
				{Kind: BcRhigh, A0: 0x00, A1: 0xfd, A2: 0x02},
			},
		},
	},
	Intern("JSR"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x03},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x03},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("RET"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x04},
		},
	},
	Intern("JR"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x05},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x05},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("BCO"): InstPat{
		kwImmNN: InstPat{
			kwCondNE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x10},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x12},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x13},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			kwCondNE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x10},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x12},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x13},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	Intern("LD"): InstPat{
		kwRegA: InstPat{
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x20},
				},
			},
			kwMemX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x21},
				},
			},
			kwMemY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x22},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x23},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x30},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x24},
				},
			},
			kwMemX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x25},
				},
			},
			kwMemY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x26},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x27},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x31},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegX: InstPat{
			kwRegAB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x28},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2a},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x32},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegY: InstPat{
			kwRegAB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2c},
				},
			},
			kwRegX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2d},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2e},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2f},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x33},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2f},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegAB: InstPat{
			kwRegX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x38},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x39},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x3a},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x3b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x3b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x34},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x35},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x36},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x37},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	Intern("ADD"): InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x40},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x41},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x42},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x43},
				},
			},
			kwMemY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegAB: InstPat{
			kwRegX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x46},
				},
			},
		},
	},
	Intern("BIT"): InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcImp, A0: 0x00, A1: 0x50, A2: 0x07, A3: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcImp, A0: 0x00, A1: 0x58, A2: 0x07, A3: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	Intern("BMM"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcMap, A0: 0x00, A1: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("DNN"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x64},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x65},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
}

var ctxOpMap = map[*Keyword]map[*Keyword]map[*Keyword][][]Value{
	Intern("<-"): {
		kwRegA: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegB: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			KwAny: {
				{Intern("LD"), &Operand{Kind: kwRegX}, &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			KwAny: {
				{Intern("LD"), &Operand{Kind: kwRegY}, &Vec{Int(1), nil}},
			},
		},
		kwRegAB: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{Intern("LD"), &Vec{Int(0), nil}, Int(4660)},
			},
		},
	},
	Intern("-jump-if"): {
		kwImmNN: {
			kwCondNE: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondEQ: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondCC: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondCS: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-jump-unless"): {
		kwImmNN: {
			kwCondNE: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Operand{Kind: kwCondEQ}},
			},
			kwCondEQ: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Operand{Kind: kwCondNE}},
			},
			kwCondCC: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Operand{Kind: kwCondCS}},
			},
			kwCondCS: {
				{Intern("BCO"), &Vec{Int(0), nil}, &Operand{Kind: kwCondCC}},
			},
		},
	},
	Intern("-dnnm"): {
		kwMemNN: {
			nil: {
				{Intern("DNN"), &Vec{Int(0), kwImmNN}},
			},
		},
	},
	Intern("-byte"): {
		kwImmNN: {
			nil: {
				{Intern("#.BYTE"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-rep"): {
		kwImmNN: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("#.BYTE"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
}
