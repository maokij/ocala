package ttarch

import . "ocala/core" //lint:ignore ST1001 core

var bmaps = [][]byte{
	{3, 1, 3, 0, 97, 103, 111}, // 0: BMM
}

var kwADD = Intern("ADD")
var kwBIT = Intern("BIT")
var kwBMM = Intern("BMM")
var kwBRA = Intern("BRA")
var kwBRL = Intern("BRL")
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
var kwJSR = Intern("JSR")
var kwLD = Intern("LD")
var kwMemNN = Intern("[%W]")
var kwMemX = Intern("[X]")
var kwMemY = Intern("[Y]")
var kwNOP = Intern("NOP")
var kwRET = Intern("RET")
var kwRRA = Intern("RRA")
var kwRegA = Intern("A")
var kwRegAB = Intern("AB")
var kwRegB = Intern("B")
var kwRegP = Intern("P")
var kwRegPC = Intern("PC")
var kwRegPQ = Intern("PQ")
var kwRegSP = Intern("SP")
var kwRegX = Intern("X")
var kwRegY = Intern("Y")

var asmOperands = map[*Keyword]AsmOperand{
	kwRegA:   {Base: "A", Expand: false},
	kwRegB:   {Base: "B", Expand: false},
	kwRegP:   {Base: "P", Expand: false},
	kwRegAB:  {Base: "AB", Expand: false},
	kwRegSP:  {Base: "SP", Expand: false},
	kwRegPC:  {Base: "PC", Expand: false},
	kwRegPQ:  {Base: "PQ", Expand: false},
	kwRegX:   {Base: "X", Expand: false},
	kwMemX:   {Base: "(X)", Expand: false},
	kwRegY:   {Base: "Y", Expand: false},
	kwMemY:   {Base: "(Y)", Expand: false},
	kwImmN:   {Base: "0+ %", Expand: true},
	kwImmNN:  {Base: "0+ %", Expand: true},
	kwMemNN:  {Base: "(%)", Expand: true},
	kwCondNE: {Base: "NE", Expand: false},
	kwCondEQ: {Base: "EQ", Expand: false},
	kwCondCC: {Base: "CC", Expand: false},
	kwCondCS: {Base: "CS", Expand: false},
}

var tokenWords = [][]string{
	{"A", "B", "P", "X", "Y", "AB", "SP", "PC"},
	{"NE?", "EQ?", "CC?", "CS?"},
	{"-jump", "-return", "-dnnm", "-byte"},
	{"<-", "-jump-if", "-jump-unless", "-return-if", "-return-unless", "-rep"},
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
		nil: InstDat{ // NOP
			{Kind: BcByte, A0: 0x00},
		},
	},
	kwJMP: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // JMP NN
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // JMP NN NE?
					{Kind: BcByte, A0: 0x10},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // JMP NN EQ?
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // JMP NN CC?
					{Kind: BcByte, A0: 0x12},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // JMP NN CS?
					{Kind: BcByte, A0: 0x13},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // JMP N
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // JMP N NE?
					{Kind: BcByte, A0: 0x10},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // JMP N EQ?
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // JMP N CC?
					{Kind: BcByte, A0: 0x12},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // JMP N CS?
					{Kind: BcByte, A0: 0x13},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	kwBRL: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BRL NN
				{Kind: BcByte, A0: 0x02},
				{Kind: BcRlow, A0: 0x00, A1: 0xfd, A2: 0x02},
				{Kind: BcRhigh, A0: 0x00, A1: 0xfd, A2: 0x02},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BRL N
				{Kind: BcByte, A0: 0x02},
				{Kind: BcRlow, A0: 0x00, A1: 0xfd, A2: 0x02},
				{Kind: BcRhigh, A0: 0x00, A1: 0xfd, A2: 0x02},
			},
		},
	},
	kwBRA: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BRA NN
				{Kind: BcByte, A0: 0x03},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BRA N
				{Kind: BcByte, A0: 0x03},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwRET: InstPat{
		nil: InstDat{ // RET
			{Kind: BcByte, A0: 0x04},
		},
		kwCondNE: InstPat{
			nil: InstDat{ // RET NE?
				{Kind: BcByte, A0: 0x20},
			},
		},
		kwCondEQ: InstPat{
			nil: InstDat{ // RET EQ?
				{Kind: BcByte, A0: 0x21},
			},
		},
		kwCondCC: InstPat{
			nil: InstDat{ // RET CC?
				{Kind: BcByte, A0: 0x22},
			},
		},
		kwCondCS: InstPat{
			nil: InstDat{ // RET CS?
				{Kind: BcByte, A0: 0x23},
			},
		},
	},
	kwJSR: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // JSR NN
				{Kind: BcByte, A0: 0x05},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // JSR N
				{Kind: BcByte, A0: 0x05},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	KwJump: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // #.jump NN
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // #.jump NN NE?
					{Kind: BcByte, A0: 0x10},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // #.jump NN EQ?
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // #.jump NN CC?
					{Kind: BcByte, A0: 0x12},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // #.jump NN CS?
					{Kind: BcByte, A0: 0x13},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // #.jump N
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // #.jump N NE?
					{Kind: BcByte, A0: 0x10},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // #.jump N EQ?
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // #.jump N CC?
					{Kind: BcByte, A0: 0x12},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // #.jump N CS?
					{Kind: BcByte, A0: 0x13},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	KwCall: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // #.call NN
				{Kind: BcByte, A0: 0x05},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // #.call NN NE?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // #.call NN EQ?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // #.call NN CC?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // #.call NN CS?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // #.call N
				{Kind: BcByte, A0: 0x05},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // #.call N NE?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // #.call N EQ?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // #.call N CC?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // #.call N CS?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	KwReturn: InstPat{
		nil: InstDat{ // #.return
			{Kind: BcByte, A0: 0x04},
		},
		kwCondNE: InstPat{
			nil: InstDat{ // #.return NE?
				{Kind: BcByte, A0: 0x20},
			},
		},
		kwCondEQ: InstPat{
			nil: InstDat{ // #.return EQ?
				{Kind: BcByte, A0: 0x21},
			},
		},
		kwCondCC: InstPat{
			nil: InstDat{ // #.return CC?
				{Kind: BcByte, A0: 0x22},
			},
		},
		kwCondCS: InstPat{
			nil: InstDat{ // #.return CS?
				{Kind: BcByte, A0: 0x23},
			},
		},
	},
	kwLD: InstPat{
		kwRegA: InstPat{
			kwRegB: InstPat{
				nil: InstDat{ // LD A B
					{Kind: BcByte, A0: 0x20},
				},
			},
			kwMemX: InstPat{
				nil: InstDat{ // LD A X$
					{Kind: BcByte, A0: 0x21},
				},
			},
			kwMemY: InstPat{
				nil: InstDat{ // LD A Y$
					{Kind: BcByte, A0: 0x22},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD A N
					{Kind: BcByte, A0: 0x23},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD A NN$
					{Kind: BcByte, A0: 0x30},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD A NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD B A
					{Kind: BcByte, A0: 0x24},
				},
			},
			kwMemX: InstPat{
				nil: InstDat{ // LD B X$
					{Kind: BcByte, A0: 0x25},
				},
			},
			kwMemY: InstPat{
				nil: InstDat{ // LD B Y$
					{Kind: BcByte, A0: 0x26},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD B N
					{Kind: BcByte, A0: 0x27},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD B NN$
					{Kind: BcByte, A0: 0x31},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD B NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegX: InstPat{
			kwRegAB: InstPat{
				nil: InstDat{ // LD X AB
					{Kind: BcByte, A0: 0x28},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{ // LD X Y
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // LD X SP
					{Kind: BcByte, A0: 0x2a},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD X NN
					{Kind: BcByte, A0: 0x2b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD X NN$
					{Kind: BcByte, A0: 0x32},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD X N
					{Kind: BcByte, A0: 0x2b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegY: InstPat{
			kwRegAB: InstPat{
				nil: InstDat{ // LD Y AB
					{Kind: BcByte, A0: 0x2c},
				},
			},
			kwRegX: InstPat{
				nil: InstDat{ // LD Y X
					{Kind: BcByte, A0: 0x2d},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // LD Y SP
					{Kind: BcByte, A0: 0x2e},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD Y NN
					{Kind: BcByte, A0: 0x2f},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD Y NN$
					{Kind: BcByte, A0: 0x33},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD Y N
					{Kind: BcByte, A0: 0x2f},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegAB: InstPat{
			kwRegX: InstPat{
				nil: InstDat{ // LD AB X
					{Kind: BcByte, A0: 0x38},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{ // LD AB Y
					{Kind: BcByte, A0: 0x39},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // LD AB SP
					{Kind: BcByte, A0: 0x3a},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD AB NN
					{Kind: BcByte, A0: 0x3b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD AB N
					{Kind: BcByte, A0: 0x3b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD NN$ A
					{Kind: BcByte, A0: 0x34},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD NN$ B
					{Kind: BcByte, A0: 0x35},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegX: InstPat{
				nil: InstDat{ // LD NN$ X
					{Kind: BcByte, A0: 0x36},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{ // LD NN$ Y
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
				nil: InstDat{ // ADD A A
					{Kind: BcByte, A0: 0x40},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // ADD A B
					{Kind: BcByte, A0: 0x41},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // ADD A N
					{Kind: BcByte, A0: 0x42},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemX: InstPat{
				nil: InstDat{ // ADD A X$
					{Kind: BcByte, A0: 0x43},
				},
			},
			kwMemY: InstPat{
				nil: InstDat{ // ADD A Y$
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // ADD A NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegAB: InstPat{
			kwRegX: InstPat{
				nil: InstDat{ // ADD AB X
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwRegY: InstPat{
				nil: InstDat{ // ADD AB Y
					{Kind: BcByte, A0: 0x46},
				},
			},
		},
	},
	kwBIT: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // BIT N A
					{Kind: BcImp, A0: 0x00, A1: 0x50, A2: 0x07, A3: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // BIT N B
					{Kind: BcImp, A0: 0x00, A1: 0x58, A2: 0x07, A3: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // BIT NN A
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // BIT NN B
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwBMM: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // BMM N
				{Kind: BcMap, A0: 0x00, A1: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // BMM NN
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwDNN: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // DNN N
				{Kind: BcByte, A0: 0x64},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // DNN NN
				{Kind: BcByte, A0: 0x65},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwRRA: InstPat{
		nil: InstDat{ // RRA
			{Kind: BcByte, A0: 0x63},
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
	Intern("-return"): {
		kwRegPC: {
			nil: {
				{KwReturn},
			},
		},
	},
	Intern("-return-if"): {
		kwRegPC: {
			kwCondNE: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondEQ: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondCC: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondCS: {
				{KwReturn, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-return-unless"): {
		kwRegPC: {
			kwCondNE: {
				{KwReturn, &Operand{Kind: kwCondEQ}},
			},
			kwCondEQ: {
				{KwReturn, &Operand{Kind: kwCondNE}},
			},
			kwCondCC: {
				{KwReturn, &Operand{Kind: kwCondCS}},
			},
			kwCondCS: {
				{KwReturn, &Operand{Kind: kwCondCC}},
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
				nil: InstDat{ // EXT A N
					{Kind: BcByte, A0: 0x70},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // EXT A NN
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
