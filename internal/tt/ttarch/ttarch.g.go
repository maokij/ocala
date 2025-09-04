package ttarch

import . "ocala/internal/core" //lint:ignore ST1001 core

var bmaps = [][]byte{
	{3, 1, 3, 0, 97, 103, 111}, // 0: BMM
}

var kwADD = Intern("ADD")
var kwBCO = Intern("BCO")
var kwBIT = Intern("BIT")
var kwBMM = Intern("BMM")
var kwBYTE = Intern("#.BYTE")
var kwCondCC = Intern("CC?")
var kwCondCS = Intern("CS?")
var kwCondEQ = Intern("EQ?")
var kwCondNE = Intern("NE?")
var kwDNN = Intern("DNN")
var kwEXT = Intern("EXT")
var kwImmN = Intern("%B")
var kwImmNN = Intern("%W")
var kwJMP = Intern("JMP")
var kwJPR = Intern("JPR")
var kwJR = Intern("JR")
var kwJSR = Intern("JSR")
var kwLD = Intern("LD")
var kwMemNN = Intern("[%W]")
var kwMemX = Intern("[X]")
var kwMemY = Intern("[Y]")
var kwNOP = Intern("NOP")
var kwRET = Intern("RET")
var kwRegA = Intern("A")
var kwRegAB = Intern("AB")
var kwRegB = Intern("B")
var kwRegP = Intern("P")
var kwRegPQ = Intern("PQ")
var kwRegSP = Intern("SP")
var kwRegX = Intern("X")
var kwRegY = Intern("Y")

var asmOperands = map[*Keyword]AsmOperand{
	kwRegA:   {"A", false},
	kwRegB:   {"B", false},
	kwRegP:   {"P", false},
	kwRegAB:  {"AB", false},
	kwRegSP:  {"SP", false},
	kwRegPQ:  {"PQ", false},
	kwRegX:   {"X", false},
	kwMemX:   {"(X)", false},
	kwRegY:   {"Y", false},
	kwMemY:   {"(Y)", false},
	kwImmN:   {"0+ %", true},
	kwImmNN:  {"0+ %", true},
	kwMemNN:  {"(%)", true},
	kwCondNE: {"NE", false},
	kwCondEQ: {"EQ", false},
	kwCondCC: {"CC", false},
	kwCondCS: {"CS", false},
}

var tokenWords = [][]string{
	{"A", "B", "P", "X", "Y", "AB", "SP"},
	{"NE?", "EQ?", "CC?", "CS?"},
	{"-jump", "-dnnm", "-byte"},
	{"<-", "-jump-if", "-jump-unless", "-rep"},
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
	kwNOP: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x00},
		},
	},
	kwJMP: InstPat{
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
	kwJPR: InstPat{
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
	kwJSR: InstPat{
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
	kwRET: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x04},
		},
	},
	kwJR: InstPat{
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
	kwBCO: InstPat{
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
	KwJump: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
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
			nil: InstDat{
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
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
	KwCall: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x03},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x03},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwLD: InstPat{
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
	kwADD: InstPat{
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
	kwBIT: InstPat{
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
	kwBMM: InstPat{
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
	kwDNN: InstPat{
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

var ctxOpMap = CtxOpMap{
	Intern("<-"): {
		kwRegA: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegB: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			KwAny: {
				{kwLD, &Operand{Kind: kwRegX}, &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			KwAny: {
				{kwLD, &Operand{Kind: kwRegY}, &Vec{Int(1), nil}},
			},
		},
		kwRegAB: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{kwLD, &Vec{Int(0), nil}, Int(4660)},
			},
		},
	},
	Intern("-jump"): {
		kwImmNN: {
			nil: {
				{KwJump, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-jump-if"): {
		kwImmNN: {
			kwCondNE: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondEQ: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondCC: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondCS: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-jump-unless"): {
		kwImmNN: {
			kwCondNE: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondEQ}},
			},
			kwCondEQ: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondNE}},
			},
			kwCondCC: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondCS}},
			},
			kwCondCS: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondCC}},
			},
		},
	},
	Intern("-dnnm"): {
		kwMemNN: {
			nil: {
				{kwDNN, &Vec{Int(0), kwImmNN}},
			},
		},
	},
	Intern("-byte"): {
		kwImmNN: {
			nil: {
				{kwBYTE, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-rep"): {
		kwImmNN: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwBYTE, &Vec{Int(0), nil}},
				}},
			},
		},
	},
}

var tokenWordsExt = [][]string{
	{},
	{},
	{},
	{"-jump", "-ext"},
}

var instMapExt = InstPat{
	kwEXT: InstPat{
		kwRegA: InstPat{
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x70},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
}

var ctxOpMapExt = CtxOpMap{
	Intern("<-"): {
		kwRegA: {
			kwImmNN: {
				{kwEXT, &Operand{Kind: kwRegA}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-jump"): {
		kwRegA: {
			kwImmNN: {
				{kwEXT, &Operand{Kind: kwRegA}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-ext"): {
		kwRegA: {
			kwImmNN: {
				{kwEXT, &Operand{Kind: kwRegA}, &Vec{Int(1), nil}},
			},
		},
	},
}
