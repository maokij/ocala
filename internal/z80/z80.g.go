package z80

import . "ocala/internal/core" //lint:ignore ST1001 core

var kwRegA = Intern("A")
var kwRegB = Intern("B")
var kwRegC = Intern("C")
var kwRegD = Intern("D")
var kwRegE = Intern("E")
var kwRegH = Intern("H")
var kwRegL = Intern("L")
var kwRegHL = Intern("HL")
var kwMemHL = Intern("[HL]")
var kwRegBC = Intern("BC")
var kwMemBC = Intern("[BC]")
var kwRegDE = Intern("DE")
var kwMemDE = Intern("[DE]")
var kwRegAF = Intern("AF")
var kwAltAF = Intern("AF-")
var kwRegSP = Intern("SP")
var kwMemSP = Intern("[SP]")
var kwRegPQ = Intern("PQ")
var kwRegIX = Intern("IX")
var kwMemIX = Intern("[IX %B]")
var kwRegIY = Intern("IY")
var kwMemIY = Intern("[IY %B]")
var kwImmN = Intern("%B")
var kwMemN = Intern("[%B]")
var kwImmNN = Intern("%W")
var kwMemNN = Intern("[%W]")
var kwMemC = Intern("[C]")
var kwRegI = Intern("I")
var kwRegR = Intern("R")
var kwCondNZ = Intern("NZ?")
var kwCondZ = Intern("Z?")
var kwCondNC = Intern("NC?")
var kwCondC = Intern("C?")
var kwCondPO = Intern("PO?")
var kwCondPE = Intern("PE?")
var kwCondP = Intern("P?")
var kwCondM = Intern("M?")

var operandToAsmMap = map[*Keyword](struct {
	s string
	t bool
}){
	kwRegA:   {s: "A", t: false},
	kwRegB:   {s: "B", t: false},
	kwRegC:   {s: "C", t: false},
	kwRegD:   {s: "D", t: false},
	kwRegE:   {s: "E", t: false},
	kwRegH:   {s: "H", t: false},
	kwRegL:   {s: "L", t: false},
	kwRegHL:  {s: "HL", t: false},
	kwMemHL:  {s: "(HL)", t: false},
	kwRegBC:  {s: "BC", t: false},
	kwMemBC:  {s: "(BC)", t: false},
	kwRegDE:  {s: "DE", t: false},
	kwMemDE:  {s: "(DE)", t: false},
	kwRegAF:  {s: "AF", t: false},
	kwAltAF:  {s: "AF'", t: false},
	kwRegSP:  {s: "SP", t: false},
	kwMemSP:  {s: "(SP)", t: false},
	kwRegPQ:  {s: "PQ", t: false},
	kwRegIX:  {s: "IX", t: false},
	kwMemIX:  {s: "(IX+%)", t: true},
	kwRegIY:  {s: "IY", t: false},
	kwMemIY:  {s: "(IY+%)", t: true},
	kwImmN:   {s: "0+ %", t: true},
	kwMemN:   {s: "(%)", t: true},
	kwImmNN:  {s: "0+ %", t: true},
	kwMemNN:  {s: "(%)", t: true},
	kwMemC:   {s: "(C)", t: false},
	kwRegI:   {s: "I", t: false},
	kwRegR:   {s: "R", t: false},
	kwCondNZ: {s: "NZ", t: false},
	kwCondZ:  {s: "Z", t: false},
	kwCondNC: {s: "NC", t: false},
	kwCondC:  {s: "C", t: false},
	kwCondPO: {s: "PO", t: false},
	kwCondPE: {s: "PE", t: false},
	kwCondP:  {s: "P", t: false},
	kwCondM:  {s: "M", t: false},
}

var tokenWords = [][]string{
	{"A", "B", "C", "D", "E", "H", "L", "I", "R", "AF", "AF-", "BC", "DE", "HL", "IX", "IY", "SP"},
	{"NZ?", "Z?", "NC?", "C?", "PO?", "PE?", "P?", "M?"},
	{"-push", "-pop", "++", "--", "-not", "-neg", "-zero?"},
	{"<-", "->", "<->", "+$", "-$", "-?", "<*", "<*$", ">*", ">*$", "-set", "-reset", "-bit?", "-in", "-out", "-jump-if", "-jump-unless"},
}

var bmaps = [][]byte{
	{255, 0, 0, 233},         // 0: JP
	{3, 0, 2, 70, 86, 94, 0}, // 1: IM
}

var tokenAliases = map[string]string{
	"!=?":        "NZ?",
	"not-zero?":  "NZ?",
	"==?":        "Z?",
	"zero?":      "Z?",
	">=?":        "NC?",
	"not-carry?": "NC?",
	"<?":         "C?",
	"carry?":     "C?",
	"odd?":       "PO?",
	"not-over?":  "PO?",
	"even?":      "PE?",
	"over?":      "PE?",
	"plus?":      "P?",
	"minus?":     "M?",
}

var instAliases = map[string][]string{}

var instMap = InstPat{
	Intern("LD"): InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x7f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x78},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x79},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x7a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x7b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x7d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x3e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x7e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x0a},
				},
			},
			kwMemDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x1a},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x3a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwRegI: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x57},
				},
			},
			kwRegR: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5f},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x3a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegB: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x47},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x40},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x41},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x42},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x43},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x06},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x46},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x46},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x46},
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
		kwRegC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x4f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x48},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x49},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x4a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x4b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x4d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x0e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x4e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4e},
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
		kwRegD: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x57},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x50},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x51},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x52},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x53},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x55},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x16},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x56},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x56},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x56},
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
		kwRegE: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x5f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x58},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x59},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x5a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x5b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x5d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x1e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x5e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5e},
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
		kwRegH: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x66},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x66},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x66},
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
		kwRegL: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x6e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6e},
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
		kwMemHL: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x77},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x70},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x71},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x72},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x73},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x74},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x75},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x36},
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
		kwMemIX: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x77},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x70},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x71},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x72},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x74},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x75},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x36},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemIY: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x77},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x70},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x71},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x72},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x74},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x75},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x36},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemBC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x02},
				},
			},
		},
		kwMemDE: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x12},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x32},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x43},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x53},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwRegI: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x47},
				},
			},
		},
		kwRegR: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4f},
				},
			},
		},
		kwRegBC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x01},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x01},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegDE: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegHL: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegSP: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x31},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x7b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf9},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xf9},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xf9},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x31},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x7b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegIX: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegIY: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwMemN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x32},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x43},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x53},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	Intern("PUSH"): InstPat{
		kwRegBC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc5},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd5},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe5},
			},
		},
		kwRegAF: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf5},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xe5},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xe5},
			},
		},
	},
	Intern("POP"): InstPat{
		kwRegBC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc1},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd1},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe1},
			},
		},
		kwRegAF: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf1},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xe1},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xe1},
			},
		},
	},
	Intern("EX"): InstPat{
		kwRegDE: InstPat{
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xeb},
				},
			},
		},
		kwRegAF: InstPat{
			kwAltAF: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x08},
				},
			},
		},
		kwMemSP: InstPat{
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe3},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xe3},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xe3},
				},
			},
		},
	},
	Intern("EXX"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xd9},
		},
	},
	Intern("LDI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa0},
		},
	},
	Intern("LDIR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb0},
		},
	},
	Intern("LDD"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa8},
		},
	},
	Intern("LDDR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb8},
		},
	},
	Intern("CPI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa1},
		},
	},
	Intern("CPIR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb1},
		},
	},
	Intern("CPD"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa9},
		},
	},
	Intern("CPDR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb9},
		},
	},
	Intern("ADD"): InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x87},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x80},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x81},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x82},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x83},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x85},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc6},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x86},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x86},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x86},
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
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x39},
				},
			},
		},
		kwRegIX: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x39},
				},
			},
		},
		kwRegIY: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x39},
				},
			},
		},
	},
	Intern("ADC"): InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x8f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x88},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x89},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x8a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x8b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x8d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xce},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x8e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8e},
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
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4a},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5a},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x7a},
				},
			},
		},
	},
	Intern("SUB"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x97},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x90},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x91},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x92},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x93},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x95},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x96},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x96},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x96},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("SBC"): InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x9f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x98},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x99},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x9a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x9b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x9d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xde},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x9e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9e},
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
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x42},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x52},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x72},
				},
			},
		},
	},
	Intern("AND"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa7},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa0},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa1},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa2},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa3},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa5},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa6},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("OR"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb7},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb0},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb1},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb2},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb3},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb5},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb6},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("XOR"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xaf},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa8},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xa9},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xaa},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xab},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xad},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xee},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xae},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xae},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xae},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("CP"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbf},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb8},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xb9},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xba},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbb},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbd},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfe},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xbe},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbe},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbe},
				{Kind: BcLow, A0: 0x00},
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
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x3c},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x04},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x0c},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x14},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x1c},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x34},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x34},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x34},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwRegBC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x03},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x13},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x23},
			},
		},
		kwRegSP: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x33},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x23},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x23},
			},
		},
	},
	Intern("DEC"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x3d},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x05},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x0d},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x15},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x1d},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x35},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x35},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x35},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwRegBC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x0b},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x1b},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x2b},
			},
		},
		kwRegSP: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x3b},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2b},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2b},
			},
		},
	},
	Intern("RLCA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x07},
		},
	},
	Intern("RLA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x17},
		},
	},
	Intern("RRCA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x0f},
		},
	},
	Intern("RRA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x1f},
		},
	},
	Intern("RLC"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x07},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x01},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x02},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x03},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x04},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x05},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x06},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x06},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x06},
			},
		},
	},
	Intern("RL"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x17},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x10},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x11},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x12},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x13},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x14},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x15},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x16},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x16},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x16},
			},
		},
	},
	Intern("RRC"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x08},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x09},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x0e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x0e},
			},
		},
	},
	Intern("RR"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x18},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x19},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x1e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x1e},
			},
		},
	},
	Intern("SLA"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x27},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x20},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x21},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x22},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x23},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x26},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x26},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x26},
			},
		},
	},
	Intern("SRA"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x28},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x29},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x2e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x2e},
			},
		},
	},
	Intern("SRL"): InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x38},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x39},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x3e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x3e},
			},
		},
	},
	Intern("RLD"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x6f},
		},
	},
	Intern("RRD"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x67},
		},
	},
	Intern("BIT"): InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x47, A2: 0x07, A3: 0x03},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x40, A2: 0x07, A3: 0x03},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x41, A2: 0x07, A3: 0x03},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x42, A2: 0x07, A3: 0x03},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x43, A2: 0x07, A3: 0x03},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x44, A2: 0x07, A3: 0x03},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x45, A2: 0x07, A3: 0x03},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x46, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x46, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x46, A2: 0x07, A3: 0x03},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	Intern("SET"): InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x07, A3: 0x03},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc0, A2: 0x07, A3: 0x03},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc1, A2: 0x07, A3: 0x03},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc2, A2: 0x07, A3: 0x03},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc3, A2: 0x07, A3: 0x03},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc4, A2: 0x07, A3: 0x03},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc5, A2: 0x07, A3: 0x03},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc6, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0xc6, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0xc6, A2: 0x07, A3: 0x03},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	Intern("RES"): InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x87, A2: 0x07, A3: 0x03},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x80, A2: 0x07, A3: 0x03},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x81, A2: 0x07, A3: 0x03},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x82, A2: 0x07, A3: 0x03},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x83, A2: 0x07, A3: 0x03},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x84, A2: 0x07, A3: 0x03},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x85, A2: 0x07, A3: 0x03},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x86, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x86, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x86, A2: 0x07, A3: 0x03},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	Intern("JP"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe9},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcMap, A0: 0x00, A1: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcMap, A0: 0x00, A1: 0x00},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPO: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPE: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondP: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondM: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("JR"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x18},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x38},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x38},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x30},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x30},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x28},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x28},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x20},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x20},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x18},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	Intern("DJNZ"): InstPat{
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
	Intern("CALL"): InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPO: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPE: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondP: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondM: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	Intern("RET"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xc9},
		},
		kwCondNZ: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc0},
			},
		},
		kwCondZ: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc8},
			},
		},
		kwCondNC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd0},
			},
		},
		kwCondC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xd8},
			},
		},
		kwCondPO: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe0},
			},
		},
		kwCondPE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xe8},
			},
		},
		kwCondP: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf0},
			},
		},
		kwCondM: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xf8},
			},
		},
	},
	Intern("RETI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x4d},
		},
	},
	Intern("RETN"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x45},
		},
	},
	Intern("RST"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x38, A3: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	Intern("IN"): InstPat{
		kwRegA: InstPat{
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdb},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x78},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x40},
				},
			},
		},
		kwRegC: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x48},
				},
			},
		},
		kwRegD: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x50},
				},
			},
		},
		kwRegE: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x58},
				},
			},
		},
		kwRegH: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x60},
				},
			},
		},
		kwRegL: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x68},
				},
			},
		},
	},
	Intern("INI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa2},
		},
	},
	Intern("INIR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb2},
		},
	},
	Intern("IND"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xaa},
		},
	},
	Intern("INDR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xba},
		},
	},
	Intern("OUT"): InstPat{
		kwMemN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd3},
					{Kind: BcLow, A0: 0x00},
				},
			},
		},
		kwMemC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x79},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x41},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x49},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x51},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x59},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x69},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	Intern("OUTI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa3},
		},
	},
	Intern("OTIR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb3},
		},
	},
	Intern("OUTD"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xab},
		},
	},
	Intern("OTDR"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xbb},
		},
	},
	Intern("DAA"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x27},
		},
	},
	Intern("CPL"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x2f},
		},
	},
	Intern("NEG"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x44},
		},
	},
	Intern("CCF"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x3f},
		},
	},
	Intern("SCF"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x37},
		},
	},
	Intern("NOP"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x00},
		},
	},
	Intern("HALT"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x76},
		},
	},
	Intern("DI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xf3},
		},
	},
	Intern("EI"): InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xfb},
		},
	},
	Intern("IM"): InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xed},
				{Kind: BcMap, A0: 0x00, A1: 0x01},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
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
		kwRegC: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegD: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegE: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegH: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegL: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegI: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegR: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegBC: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegDE: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegHL: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegSP: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIX: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIY: {
			KwAny: {
				{Intern("LD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("->"): {
		kwRegA: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegB: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegC: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegD: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegE: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegH: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegL: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegI: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegR: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegBC: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegDE: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegHL: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegSP: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{Intern("#.LDP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIX: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIY: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwImmNN: {
			KwAny: {
				{Intern("LD"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("<->"): {
		kwRegAF: {
			kwAltAF: {
				{Intern("EX"), &Operand{Kind: kwRegAF}, &Operand{Kind: kwAltAF}},
			},
		},
		kwRegDE: {
			kwRegHL: {
				{Intern("EX"), &Operand{Kind: kwRegDE}, &Operand{Kind: kwRegHL}},
			},
		},
		kwRegHL: {
			kwRegDE: {
				{Intern("EX"), &Operand{Kind: kwRegDE}, &Operand{Kind: kwRegHL}},
			},
			kwMemSP: {
				{Intern("EX"), &Operand{Kind: kwMemSP}, &Operand{Kind: kwRegHL}},
			},
		},
		kwMemSP: {
			kwRegHL: {
				{Intern("EX"), &Operand{Kind: kwMemSP}, &Operand{Kind: kwRegHL}},
			},
		},
	},
	Intern("-push"): {
		KwAny: {
			nil: {
				{Intern("PUSH"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-pop"): {
		KwAny: {
			nil: {
				{Intern("POP"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("++"): {
		KwAny: {
			nil: {
				{Intern("INC"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("--"): {
		KwAny: {
			nil: {
				{Intern("DEC"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-not"): {
		kwRegA: {
			nil: {
				{Intern("CPL")},
			},
		},
	},
	Intern("-neg"): {
		kwRegA: {
			nil: {
				{Intern("NEG")},
			},
		},
	},
	Intern("+"): {
		KwAny: {
			KwAny: {
				{Intern("ADD"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("+$"): {
		KwAny: {
			KwAny: {
				{Intern("ADC"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-"): {
		kwRegHL: {
			KwAny: {
				{Intern("OR"), &Operand{Kind: kwRegA}},
				{Intern("SBC"), &Operand{Kind: kwRegHL}, &Vec{Int(1), nil}},
			},
		},
		kwRegA: {
			KwAny: {
				{Intern("SUB"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-$"): {
		KwAny: {
			KwAny: {
				{Intern("SBC"), &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-?"): {
		kwRegA: {
			KwAny: {
				{Intern("CP"), &Vec{Int(1), nil}},
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
				{Intern("OR"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("^"): {
		kwRegA: {
			KwAny: {
				{Intern("XOR"), &Vec{Int(1), nil}},
			},
		},
	},
	Intern("<*"): {
		kwRegA: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RLCA")},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RLC"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("<*$"): {
		kwRegA: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RLA")},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RL"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">*"): {
		kwRegA: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RRCA")},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RRC"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">*$"): {
		kwRegA: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RRA")},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("RR"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("<<"): {
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("SLA"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">>"): {
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("SRA"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">>>"): {
		KwAny: {
			kwImmNN: {
				{Intern("#.REP"), &Vec{Int(1), nil}, &Vec{
					&Vec{Intern("SRL"), &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("-set"): {
		KwAny: {
			KwAny: {
				{Intern("SET"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-reset"): {
		KwAny: {
			KwAny: {
				{Intern("RES"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-bit?"): {
		KwAny: {
			KwAny: {
				{Intern("BIT"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-in"): {
		kwRegA: {
			kwRegC: {
				{Intern("IN"), &Operand{Kind: kwRegA}, &Operand{Kind: kwMemC}},
			},
			kwImmNN: {
				{Intern("IN"), &Operand{Kind: kwRegA}, &Vec{Int(1), kwMemNN}},
			},
		},
		KwAny: {
			kwRegC: {
				{Intern("IN"), &Vec{Int(0), nil}, &Operand{Kind: kwMemC}},
			},
		},
	},
	Intern("-out"): {
		kwRegA: {
			kwRegC: {
				{Intern("OUT"), &Operand{Kind: kwMemC}, &Operand{Kind: kwRegA}},
			},
			kwImmNN: {
				{Intern("OUT"), &Vec{Int(1), kwMemNN}, &Operand{Kind: kwRegA}},
			},
		},
		KwAny: {
			kwRegC: {
				{Intern("OUT"), &Operand{Kind: kwMemC}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-zero?"): {
		kwRegA: {
			nil: {
				{Intern("AND"), &Vec{Int(0), nil}},
			},
		},
		kwRegBC: {
			nil: {
				{Intern("#.INVALID"), &Vec{Int(0), nil}},
			},
		},
		kwRegDE: {
			nil: {
				{Intern("#.INVALID"), &Vec{Int(0), nil}},
			},
		},
		kwRegHL: {
			nil: {
				{Intern("#.INVALID"), &Vec{Int(0), nil}},
			},
		},
		kwRegSP: {
			nil: {
				{Intern("#.INVALID"), &Vec{Int(0), nil}},
			},
		},
		kwRegIX: {
			nil: {
				{Intern("#.INVALID"), &Vec{Int(0), nil}},
			},
		},
		kwRegIY: {
			nil: {
				{Intern("#.INVALID"), &Vec{Int(0), nil}},
			},
		},
		KwAny: {
			nil: {
				{Intern("INC"), &Vec{Int(0), nil}},
				{Intern("DEC"), &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-jump-if"): {
		kwImmNN: {
			kwCondNZ: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwCondZ: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwCondNC: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwCondC: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwCondPO: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwCondPE: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwCondP: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwCondM: {
				{Intern("JP"), &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-jump-unless"): {
		kwImmNN: {
			kwCondNZ: {
				{Intern("JP"), &Operand{Kind: kwCondZ}, &Vec{Int(0), nil}},
			},
			kwCondZ: {
				{Intern("JP"), &Operand{Kind: kwCondNZ}, &Vec{Int(0), nil}},
			},
			kwCondNC: {
				{Intern("JP"), &Operand{Kind: kwCondC}, &Vec{Int(0), nil}},
			},
			kwCondC: {
				{Intern("JP"), &Operand{Kind: kwCondNC}, &Vec{Int(0), nil}},
			},
			kwCondPO: {
				{Intern("JP"), &Operand{Kind: kwCondPE}, &Vec{Int(0), nil}},
			},
			kwCondPE: {
				{Intern("JP"), &Operand{Kind: kwCondPO}, &Vec{Int(0), nil}},
			},
			kwCondM: {
				{Intern("JP"), &Operand{Kind: kwCondP}, &Vec{Int(0), nil}},
			},
			kwCondP: {
				{Intern("JP"), &Operand{Kind: kwCondM}, &Vec{Int(0), nil}},
			},
		},
	},
}
