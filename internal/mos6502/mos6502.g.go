package mos6502

import . "ocala/internal/core" //lint:ignore ST1001 core

var bmaps = [][]byte{}

var kwADC = Intern("ADC")
var kwAND = Intern("AND")
var kwASL = Intern("ASL")
var kwBCC = Intern("BCC")
var kwBCS = Intern("BCS")
var kwBEQ = Intern("BEQ")
var kwBIT = Intern("BIT")
var kwBMI = Intern("BMI")
var kwBNE = Intern("BNE")
var kwBPL = Intern("BPL")
var kwBRK = Intern("BRK")
var kwBVC = Intern("BVC")
var kwBVS = Intern("BVS")
var kwBYTE = Intern("#.BYTE")
var kwCLC = Intern("CLC")
var kwCLD = Intern("CLD")
var kwCLI = Intern("CLI")
var kwCLV = Intern("CLV")
var kwCMP = Intern("CMP")
var kwCPX = Intern("CPX")
var kwCPY = Intern("CPY")
var kwCondCC = Intern("CC?")
var kwCondCS = Intern("CS?")
var kwCondEQ = Intern("EQ?")
var kwCondMI = Intern("MI?")
var kwCondNE = Intern("NE?")
var kwCondPL = Intern("PL?")
var kwCondVC = Intern("VC?")
var kwCondVS = Intern("VS?")
var kwDEC = Intern("DEC")
var kwDEX = Intern("DEX")
var kwDEY = Intern("DEY")
var kwEOR = Intern("EOR")
var kwINC = Intern("INC")
var kwINX = Intern("INX")
var kwINY = Intern("INY")
var kwImmN = Intern("%B")
var kwImmNN = Intern("%W")
var kwJMP = Intern("JMP")
var kwJSR = Intern("JSR")
var kwLDA = Intern("LDA")
var kwLDX = Intern("LDX")
var kwLDY = Intern("LDY")
var kwLSR = Intern("LSR")
var kwMemAN = Intern("[%W]")
var kwMemAX = Intern("[%W X]")
var kwMemAY = Intern("[%W Y]")
var kwMemIN = Intern("[[%W]]")
var kwMemIX = Intern("[[%B X]]")
var kwMemIY = Intern("[[%B] Y]")
var kwMemZN = Intern("[%B]")
var kwMemZX = Intern("[%B X]")
var kwMemZY = Intern("[%B Y]")
var kwNOP = Intern("NOP")
var kwORA = Intern("ORA")
var kwPHA = Intern("PHA")
var kwPHP = Intern("PHP")
var kwPLA = Intern("PLA")
var kwPLP = Intern("PLP")
var kwROL = Intern("ROL")
var kwROR = Intern("ROR")
var kwRTI = Intern("RTI")
var kwRTS = Intern("RTS")
var kwRegA = Intern("A")
var kwRegP = Intern("P")
var kwRegS = Intern("S")
var kwRegX = Intern("X")
var kwRegY = Intern("Y")
var kwSBC = Intern("SBC")
var kwSEC = Intern("SEC")
var kwSED = Intern("SED")
var kwSEI = Intern("SEI")
var kwSTA = Intern("STA")
var kwSTX = Intern("STX")
var kwSTY = Intern("STY")
var kwTAX = Intern("TAX")
var kwTAY = Intern("TAY")
var kwTSX = Intern("TSX")
var kwTXA = Intern("TXA")
var kwTXS = Intern("TXS")
var kwTYA = Intern("TYA")

var asmOperands = map[*Keyword]AsmOperand{
	kwRegA:   {"A", false},
	kwRegX:   {"X", false},
	kwRegY:   {"Y", false},
	kwRegS:   {"S", false},
	kwRegP:   {"P", false},
	kwImmN:   {"#%", true},
	kwImmNN:  {"#%", true},
	kwMemZN:  {"%", true},
	kwMemZX:  {"%, X", true},
	kwMemZY:  {"%, Y", true},
	kwMemAN:  {"%", true},
	kwMemAX:  {"%, X", true},
	kwMemAY:  {"%, Y", true},
	kwMemIN:  {"(%)", true},
	kwMemIX:  {"(%, X)", true},
	kwMemIY:  {"(%), Y", true},
	kwCondNE: {"NE", false},
	kwCondEQ: {"EQ", false},
	kwCondCC: {"CC", false},
	kwCondCS: {"CS", false},
	kwCondVC: {"VC", false},
	kwCondVS: {"VS", false},
	kwCondPL: {"PL", false},
	kwCondMI: {"MI", false},
}

var tokenWords = [][]string{
	{"A", "X", "Y", "S", "P"},
	{"NE?", "EQ?", "CC?", "CS?", "VC?", "VS?", "PL?", "MI?"},
	{"-push", "-pop", "++", "--", "-not", "-neg", "-jump"},
	{"<-", "->", "+$", "-$", "-?", "-bit?", "<*", "<*$", ">*", ">*$", "-jump-if", "-jump-unless"},
}

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

var instMap = InstPat{
	kwLDA: InstPat{
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
	kwLDX: InstPat{
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
	kwLDY: InstPat{
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
	kwSTA: InstPat{
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
	kwSTX: InstPat{
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
	kwSTY: InstPat{
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
	kwTAX: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xaa},
		},
	},
	kwTAY: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xa8},
		},
	},
	kwTSX: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xba},
		},
	},
	kwTXA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x8a},
		},
	},
	kwTXS: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x9a},
		},
	},
	kwTYA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x98},
		},
	},
	kwPHA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x48},
		},
	},
	kwPHP: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x08},
		},
	},
	kwPLP: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x28},
		},
	},
	kwPLA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x68},
		},
	},
	kwCLC: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x18},
		},
	},
	kwCLI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x58},
		},
	},
	kwCLD: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xd8},
		},
	},
	kwCLV: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xb8},
		},
	},
	kwSEC: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x38},
		},
	},
	kwSEI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x78},
		},
	},
	kwSED: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xf8},
		},
	},
	kwBRK: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x00},
		},
	},
	kwNOP: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xea},
		},
	},
	kwRTS: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x60},
		},
	},
	kwRTI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x40},
		},
	},
	kwJMP: InstPat{
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
	kwJSR: InstPat{
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
	kwBPL: InstPat{
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
	kwBMI: InstPat{
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
	kwBVC: InstPat{
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
	kwBVS: InstPat{
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
	kwBCC: InstPat{
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
	kwBCS: InstPat{
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
	kwBNE: InstPat{
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
	kwBEQ: InstPat{
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
	KwJump: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondPL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x30},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x10},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x70},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x50},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xb0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x90},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondPL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x30},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x10},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x70},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x50},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xb0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0x90},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	KwCall: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x20},
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
			kwCondVC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondPL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0x20},
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
			kwCondVC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondPL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwORA: InstPat{
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
	kwAND: InstPat{
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
	kwEOR: InstPat{
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
	kwADC: InstPat{
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
	kwCMP: InstPat{
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
	kwSBC: InstPat{
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
	kwBIT: InstPat{
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
	kwCPX: InstPat{
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
	kwCPY: InstPat{
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
	kwINC: InstPat{
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
	kwINX: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xe8},
		},
	},
	kwINY: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xc8},
		},
	},
	kwDEC: InstPat{
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
	kwDEX: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xca},
		},
	},
	kwDEY: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x88},
		},
	},
	kwASL: InstPat{
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
	kwLSR: InstPat{
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
	kwROL: InstPat{
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
	kwROR: InstPat{
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

var ctxOpMap = CtxOpMap{
	Intern("<-"): {
		kwRegA: {
			kwRegX: {
				{kwTXA},
			},
			kwRegY: {
				{kwTYA},
			},
			KwAny: {
				{kwLDA, &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			kwRegA: {
				{kwTAX},
			},
			kwRegS: {
				{kwTSX},
			},
			KwAny: {
				{kwLDX, &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			kwRegA: {
				{kwTAY},
			},
			KwAny: {
				{kwLDY, &Vec{Int(1), nil}},
			},
		},
		kwRegS: {
			kwRegX: {
				{kwTXS},
			},
		},
	},
	Intern("->"): {
		kwRegA: {
			kwRegX: {
				{kwTAX},
			},
			kwRegY: {
				{kwTAY},
			},
			KwAny: {
				{kwSTA, &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			kwRegA: {
				{kwTXA},
			},
			kwRegS: {
				{kwTXS},
			},
			KwAny: {
				{kwSTX, &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			kwRegA: {
				{kwTYA},
			},
			KwAny: {
				{kwSTY, &Vec{Int(1), nil}},
			},
		},
		kwRegS: {
			kwRegX: {
				{kwTSX},
			},
		},
	},
	Intern("-push"): {
		kwRegA: {
			nil: {
				{kwPHA},
			},
		},
		kwRegP: {
			nil: {
				{kwPHP},
			},
		},
	},
	Intern("-pop"): {
		kwRegA: {
			nil: {
				{kwPLA},
			},
		},
		kwRegP: {
			nil: {
				{kwPLP},
			},
		},
	},
	Intern("++"): {
		kwRegA: {
			nil: {
				{kwCLC},
				{kwADC, Int(1)},
			},
		},
		kwRegX: {
			nil: {
				{kwINX},
			},
		},
		kwRegY: {
			nil: {
				{kwINY},
			},
		},
		KwAny: {
			nil: {
				{kwINC, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("--"): {
		kwRegA: {
			nil: {
				{kwSEC},
				{kwSBC, Int(1)},
			},
		},
		kwRegX: {
			nil: {
				{kwDEX},
			},
		},
		kwRegY: {
			nil: {
				{kwDEY},
			},
		},
		KwAny: {
			nil: {
				{kwDEC, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-not"): {
		kwRegA: {
			nil: {
				{kwEOR, Int(255)},
			},
		},
	},
	Intern("-neg"): {
		kwRegA: {
			nil: {
				{kwEOR, Int(255)},
				{kwCLC},
				{kwADC, Int(1)},
			},
		},
	},
	Intern("+"): {
		kwRegA: {
			KwAny: {
				{kwCLC},
				{kwADC, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("+$"): {
		kwRegA: {
			KwAny: {
				{kwADC, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-"): {
		kwRegA: {
			KwAny: {
				{kwSEC},
				{kwSBC, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-$"): {
		kwRegA: {
			KwAny: {
				{kwSBC, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-?"): {
		kwRegA: {
			KwAny: {
				{kwCMP, &Vec{Int(1), nil}},
			},
		},
		kwRegX: {
			KwAny: {
				{kwCPX, &Vec{Int(1), nil}},
			},
		},
		kwRegY: {
			KwAny: {
				{kwCPY, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("&"): {
		kwRegA: {
			KwAny: {
				{kwAND, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("|"): {
		kwRegA: {
			KwAny: {
				{kwORA, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("^"): {
		kwRegA: {
			KwAny: {
				{kwEOR, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-bit?"): {
		kwRegA: {
			KwAny: {
				{kwBIT, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("<*"): {
		kwRegA: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwCMP, Int(128)},
					&Vec{kwROL, &Operand{Kind: kwRegA}},
				}},
			},
		},
	},
	Intern("<*$"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwROL, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">*"): {
		kwRegA: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwLSR, &Operand{Kind: kwRegA}},
					&Vec{kwBYTE, Int(144)},
					&Vec{kwBYTE, Int(2)},
					&Vec{kwORA, Int(128)},
				}},
			},
		},
	},
	Intern(">*$"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwROR, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("<<"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwASL, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">>"): {
		kwRegA: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwCMP, Int(128)},
					&Vec{kwROR, &Operand{Kind: kwRegA}},
				}},
			},
		},
	},
	Intern(">>>"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwLSR, &Vec{Int(0), nil}},
				}},
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
			kwCondVC: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondVS: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondPL: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondMI: {
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
			kwCondVC: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondVS}},
			},
			kwCondVS: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondVC}},
			},
			kwCondPL: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondMI}},
			},
			kwCondMI: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondPL}},
			},
		},
	},
}
