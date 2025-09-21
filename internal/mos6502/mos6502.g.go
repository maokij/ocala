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
var kwRegPC = Intern("PC")
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
	kwRegPC:  {"PC", false},
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
	{"A", "X", "Y", "S", "P", "PC"},
	{"NE?", "EQ?", "CC?", "CS?", "VC?", "VS?", "PL?", "MI?"},
	{"-push", "-pop", "++", "--", "-not", "-neg", "-jump", "-return"},
	{"<-", "->", "+$", "-$", "-?", "-bit?", "<*", "<*$", ">*", ">*$", "-jump-if", "-jump-unless", "-return-if", "-return-unless"},
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
			nil: InstDat{ // LDA N
				{Kind: BcByte, A0: 0xa9},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // LDA ZN
				{Kind: BcByte, A0: 0xa5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // LDA ZX
				{Kind: BcByte, A0: 0xb5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // LDA AN
				{Kind: BcByte, A0: 0xad},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // LDA AX
				{Kind: BcByte, A0: 0xbd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // LDA AY
				{Kind: BcByte, A0: 0xb9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // LDA IX
				{Kind: BcByte, A0: 0xa1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // LDA IY
				{Kind: BcByte, A0: 0xb1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // LDA NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // LDA ZY
				{Kind: BcByte, A0: 0xb9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwLDX: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // LDX N
				{Kind: BcByte, A0: 0xa2},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // LDX ZN
				{Kind: BcByte, A0: 0xa6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // LDX ZY
				{Kind: BcByte, A0: 0xb6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // LDX AN
				{Kind: BcByte, A0: 0xae},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // LDX AY
				{Kind: BcByte, A0: 0xbe},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // LDX NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwLDY: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // LDY N
				{Kind: BcByte, A0: 0xa0},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // LDY ZN
				{Kind: BcByte, A0: 0xa4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // LDY ZX
				{Kind: BcByte, A0: 0xb4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // LDY AN
				{Kind: BcByte, A0: 0xac},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // LDY AX
				{Kind: BcByte, A0: 0xbc},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // LDY NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwSTA: InstPat{
		kwMemZN: InstPat{
			nil: InstDat{ // STA ZN
				{Kind: BcByte, A0: 0x85},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // STA ZX
				{Kind: BcByte, A0: 0x95},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // STA AN
				{Kind: BcByte, A0: 0x8d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // STA AX
				{Kind: BcByte, A0: 0x9d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // STA AY
				{Kind: BcByte, A0: 0x99},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // STA IX
				{Kind: BcByte, A0: 0x81},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // STA IY
				{Kind: BcByte, A0: 0x91},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // STA ZY
				{Kind: BcByte, A0: 0x99},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwSTX: InstPat{
		kwMemZN: InstPat{
			nil: InstDat{ // STX ZN
				{Kind: BcByte, A0: 0x86},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // STX ZY
				{Kind: BcByte, A0: 0x96},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // STX AN
				{Kind: BcByte, A0: 0x8e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // STX AY
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwSTY: InstPat{
		kwMemZN: InstPat{
			nil: InstDat{ // STY ZN
				{Kind: BcByte, A0: 0x84},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // STY ZX
				{Kind: BcByte, A0: 0x94},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // STY AN
				{Kind: BcByte, A0: 0x8c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // STY AX
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwTAX: InstPat{
		nil: InstDat{ // TAX
			{Kind: BcByte, A0: 0xaa},
		},
	},
	kwTAY: InstPat{
		nil: InstDat{ // TAY
			{Kind: BcByte, A0: 0xa8},
		},
	},
	kwTSX: InstPat{
		nil: InstDat{ // TSX
			{Kind: BcByte, A0: 0xba},
		},
	},
	kwTXA: InstPat{
		nil: InstDat{ // TXA
			{Kind: BcByte, A0: 0x8a},
		},
	},
	kwTXS: InstPat{
		nil: InstDat{ // TXS
			{Kind: BcByte, A0: 0x9a},
		},
	},
	kwTYA: InstPat{
		nil: InstDat{ // TYA
			{Kind: BcByte, A0: 0x98},
		},
	},
	kwPHA: InstPat{
		nil: InstDat{ // PHA
			{Kind: BcByte, A0: 0x48},
		},
	},
	kwPHP: InstPat{
		nil: InstDat{ // PHP
			{Kind: BcByte, A0: 0x08},
		},
	},
	kwPLP: InstPat{
		nil: InstDat{ // PLP
			{Kind: BcByte, A0: 0x28},
		},
	},
	kwPLA: InstPat{
		nil: InstDat{ // PLA
			{Kind: BcByte, A0: 0x68},
		},
	},
	kwCLC: InstPat{
		nil: InstDat{ // CLC
			{Kind: BcByte, A0: 0x18},
		},
	},
	kwCLI: InstPat{
		nil: InstDat{ // CLI
			{Kind: BcByte, A0: 0x58},
		},
	},
	kwCLD: InstPat{
		nil: InstDat{ // CLD
			{Kind: BcByte, A0: 0xd8},
		},
	},
	kwCLV: InstPat{
		nil: InstDat{ // CLV
			{Kind: BcByte, A0: 0xb8},
		},
	},
	kwSEC: InstPat{
		nil: InstDat{ // SEC
			{Kind: BcByte, A0: 0x38},
		},
	},
	kwSEI: InstPat{
		nil: InstDat{ // SEI
			{Kind: BcByte, A0: 0x78},
		},
	},
	kwSED: InstPat{
		nil: InstDat{ // SED
			{Kind: BcByte, A0: 0xf8},
		},
	},
	kwBRK: InstPat{
		nil: InstDat{ // BRK
			{Kind: BcByte, A0: 0x00},
		},
	},
	kwNOP: InstPat{
		nil: InstDat{ // NOP
			{Kind: BcByte, A0: 0xea},
		},
	},
	kwRTS: InstPat{
		nil: InstDat{ // RTS
			{Kind: BcByte, A0: 0x60},
		},
	},
	kwRTI: InstPat{
		nil: InstDat{ // RTI
			{Kind: BcByte, A0: 0x40},
		},
	},
	kwJMP: InstPat{
		kwMemAN: InstPat{
			nil: InstDat{ // JMP AN
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIN: InstPat{
			nil: InstDat{ // JMP IN
				{Kind: BcByte, A0: 0x6c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // JMP ZN
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwJSR: InstPat{
		kwMemAN: InstPat{
			nil: InstDat{ // JSR AN
				{Kind: BcByte, A0: 0x20},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // JSR ZN
				{Kind: BcByte, A0: 0x20},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwBPL: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BPL NN
				{Kind: BcByte, A0: 0x10},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BPL N
				{Kind: BcByte, A0: 0x10},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwBMI: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BMI NN
				{Kind: BcByte, A0: 0x30},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BMI N
				{Kind: BcByte, A0: 0x30},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwBVC: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BVC NN
				{Kind: BcByte, A0: 0x50},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BVC N
				{Kind: BcByte, A0: 0x50},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwBVS: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BVS NN
				{Kind: BcByte, A0: 0x70},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BVS N
				{Kind: BcByte, A0: 0x70},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwBCC: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BCC NN
				{Kind: BcByte, A0: 0x90},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BCC N
				{Kind: BcByte, A0: 0x90},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwBCS: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BCS NN
				{Kind: BcByte, A0: 0xb0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BCS N
				{Kind: BcByte, A0: 0xb0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwBNE: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BNE NN
				{Kind: BcByte, A0: 0xd0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BNE N
				{Kind: BcByte, A0: 0xd0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwBEQ: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // BEQ NN
				{Kind: BcByte, A0: 0xf0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // BEQ N
				{Kind: BcByte, A0: 0xf0},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	KwJump: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // #.jump NN
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondPL: InstPat{
				nil: InstDat{ // #.jump NN PL?
					{Kind: BcByte, A0: 0x30},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{ // #.jump NN MI?
					{Kind: BcByte, A0: 0x10},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVC: InstPat{
				nil: InstDat{ // #.jump NN VC?
					{Kind: BcByte, A0: 0x70},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{ // #.jump NN VS?
					{Kind: BcByte, A0: 0x50},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // #.jump NN CC?
					{Kind: BcByte, A0: 0xb0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // #.jump NN CS?
					{Kind: BcByte, A0: 0x90},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // #.jump NN NE?
					{Kind: BcByte, A0: 0xf0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // #.jump NN EQ?
					{Kind: BcByte, A0: 0xd0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // #.jump N
				{Kind: BcByte, A0: 0x4c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondPL: InstPat{
				nil: InstDat{ // #.jump N PL?
					{Kind: BcByte, A0: 0x30},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{ // #.jump N MI?
					{Kind: BcByte, A0: 0x10},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVC: InstPat{
				nil: InstDat{ // #.jump N VC?
					{Kind: BcByte, A0: 0x70},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{ // #.jump N VS?
					{Kind: BcByte, A0: 0x50},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCC: InstPat{
				nil: InstDat{ // #.jump N CC?
					{Kind: BcByte, A0: 0xb0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondCS: InstPat{
				nil: InstDat{ // #.jump N CS?
					{Kind: BcByte, A0: 0x90},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNE: InstPat{
				nil: InstDat{ // #.jump N NE?
					{Kind: BcByte, A0: 0xf0},
					{Kind: BcByte, A0: 0x03},
					{Kind: BcByte, A0: 0x4c},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondEQ: InstPat{
				nil: InstDat{ // #.jump N EQ?
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
			nil: InstDat{ // #.call NN
				{Kind: BcByte, A0: 0x20},
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
			kwCondVC: InstPat{
				nil: InstDat{ // #.call NN VC?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{ // #.call NN VS?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondPL: InstPat{
				nil: InstDat{ // #.call NN PL?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{ // #.call NN MI?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // #.call N
				{Kind: BcByte, A0: 0x20},
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
			kwCondVC: InstPat{
				nil: InstDat{ // #.call N VC?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondVS: InstPat{
				nil: InstDat{ // #.call N VS?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondPL: InstPat{
				nil: InstDat{ // #.call N PL?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwCondMI: InstPat{
				nil: InstDat{ // #.call N MI?
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	KwReturn: InstPat{
		nil: InstDat{ // #.return
			{Kind: BcByte, A0: 0x60},
		},
		kwCondPL: InstPat{
			nil: InstDat{ // #.return PL?
				{Kind: BcByte, A0: 0x30},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
		kwCondMI: InstPat{
			nil: InstDat{ // #.return MI?
				{Kind: BcByte, A0: 0x10},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
		kwCondVC: InstPat{
			nil: InstDat{ // #.return VC?
				{Kind: BcByte, A0: 0x70},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
		kwCondVS: InstPat{
			nil: InstDat{ // #.return VS?
				{Kind: BcByte, A0: 0x50},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
		kwCondCC: InstPat{
			nil: InstDat{ // #.return CC?
				{Kind: BcByte, A0: 0xb0},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
		kwCondCS: InstPat{
			nil: InstDat{ // #.return CS?
				{Kind: BcByte, A0: 0x90},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
		kwCondNE: InstPat{
			nil: InstDat{ // #.return NE?
				{Kind: BcByte, A0: 0xf0},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
		kwCondEQ: InstPat{
			nil: InstDat{ // #.return EQ?
				{Kind: BcByte, A0: 0xd0},
				{Kind: BcByte, A0: 0x01},
				{Kind: BcByte, A0: 0x60},
			},
		},
	},
	kwORA: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // ORA N
				{Kind: BcByte, A0: 0x09},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // ORA ZN
				{Kind: BcByte, A0: 0x05},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // ORA ZX
				{Kind: BcByte, A0: 0x15},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // ORA AN
				{Kind: BcByte, A0: 0x0d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // ORA AX
				{Kind: BcByte, A0: 0x1d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // ORA AY
				{Kind: BcByte, A0: 0x19},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // ORA IX
				{Kind: BcByte, A0: 0x01},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // ORA IY
				{Kind: BcByte, A0: 0x11},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // ORA NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // ORA ZY
				{Kind: BcByte, A0: 0x19},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwAND: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // AND N
				{Kind: BcByte, A0: 0x29},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // AND ZN
				{Kind: BcByte, A0: 0x25},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // AND ZX
				{Kind: BcByte, A0: 0x35},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // AND AN
				{Kind: BcByte, A0: 0x2d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // AND AX
				{Kind: BcByte, A0: 0x3d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // AND AY
				{Kind: BcByte, A0: 0x39},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // AND IX
				{Kind: BcByte, A0: 0x21},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // AND IY
				{Kind: BcByte, A0: 0x31},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // AND NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // AND ZY
				{Kind: BcByte, A0: 0x39},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwEOR: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // EOR N
				{Kind: BcByte, A0: 0x49},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // EOR ZN
				{Kind: BcByte, A0: 0x45},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // EOR ZX
				{Kind: BcByte, A0: 0x55},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // EOR AN
				{Kind: BcByte, A0: 0x4d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // EOR AX
				{Kind: BcByte, A0: 0x5d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // EOR AY
				{Kind: BcByte, A0: 0x59},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // EOR IX
				{Kind: BcByte, A0: 0x41},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // EOR IY
				{Kind: BcByte, A0: 0x51},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // EOR NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // EOR ZY
				{Kind: BcByte, A0: 0x59},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwADC: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // ADC N
				{Kind: BcByte, A0: 0x69},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // ADC ZN
				{Kind: BcByte, A0: 0x65},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // ADC ZX
				{Kind: BcByte, A0: 0x75},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // ADC AN
				{Kind: BcByte, A0: 0x6d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // ADC AX
				{Kind: BcByte, A0: 0x7d},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // ADC AY
				{Kind: BcByte, A0: 0x79},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // ADC IX
				{Kind: BcByte, A0: 0x61},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // ADC IY
				{Kind: BcByte, A0: 0x71},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // ADC NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // ADC ZY
				{Kind: BcByte, A0: 0x79},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwCMP: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // CMP N
				{Kind: BcByte, A0: 0xc9},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // CMP ZN
				{Kind: BcByte, A0: 0xc5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // CMP ZX
				{Kind: BcByte, A0: 0xd5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // CMP AN
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // CMP AX
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // CMP AY
				{Kind: BcByte, A0: 0xd9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // CMP IX
				{Kind: BcByte, A0: 0xc1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // CMP IY
				{Kind: BcByte, A0: 0xd1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // CMP NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // CMP ZY
				{Kind: BcByte, A0: 0xd9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwSBC: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // SBC N
				{Kind: BcByte, A0: 0xe9},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // SBC ZN
				{Kind: BcByte, A0: 0xe5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // SBC ZX
				{Kind: BcByte, A0: 0xf5},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // SBC AN
				{Kind: BcByte, A0: 0xed},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // SBC AX
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAY: InstPat{
			nil: InstDat{ // SBC AY
				{Kind: BcByte, A0: 0xf9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SBC IX
				{Kind: BcByte, A0: 0xe1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SBC IY
				{Kind: BcByte, A0: 0xf1},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // SBC NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemZY: InstPat{
			nil: InstDat{ // SBC ZY
				{Kind: BcByte, A0: 0xf9},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwBIT: InstPat{
		kwMemZN: InstPat{
			nil: InstDat{ // BIT ZN
				{Kind: BcByte, A0: 0x24},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // BIT AN
				{Kind: BcByte, A0: 0x2c},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwCPX: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // CPX N
				{Kind: BcByte, A0: 0xe0},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // CPX ZN
				{Kind: BcByte, A0: 0xe4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // CPX AN
				{Kind: BcByte, A0: 0xec},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // CPX NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwCPY: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // CPY N
				{Kind: BcByte, A0: 0xc0},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // CPY ZN
				{Kind: BcByte, A0: 0xc4},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // CPY AN
				{Kind: BcByte, A0: 0xcc},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // CPY NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwINC: InstPat{
		kwMemZN: InstPat{
			nil: InstDat{ // INC ZN
				{Kind: BcByte, A0: 0xe6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // INC ZX
				{Kind: BcByte, A0: 0xf6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // INC AN
				{Kind: BcByte, A0: 0xee},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // INC AX
				{Kind: BcByte, A0: 0xfe},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwINX: InstPat{
		nil: InstDat{ // INX
			{Kind: BcByte, A0: 0xe8},
		},
	},
	kwINY: InstPat{
		nil: InstDat{ // INY
			{Kind: BcByte, A0: 0xc8},
		},
	},
	kwDEC: InstPat{
		kwMemZN: InstPat{
			nil: InstDat{ // DEC ZN
				{Kind: BcByte, A0: 0xc6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // DEC ZX
				{Kind: BcByte, A0: 0xd6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // DEC AN
				{Kind: BcByte, A0: 0xce},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // DEC AX
				{Kind: BcByte, A0: 0xde},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwDEX: InstPat{
		nil: InstDat{ // DEX
			{Kind: BcByte, A0: 0xca},
		},
	},
	kwDEY: InstPat{
		nil: InstDat{ // DEY
			{Kind: BcByte, A0: 0x88},
		},
	},
	kwASL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // ASL A
				{Kind: BcByte, A0: 0x0a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // ASL ZN
				{Kind: BcByte, A0: 0x06},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // ASL ZX
				{Kind: BcByte, A0: 0x16},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // ASL AN
				{Kind: BcByte, A0: 0x0e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // ASL AX
				{Kind: BcByte, A0: 0x1e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwLSR: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // LSR A
				{Kind: BcByte, A0: 0x4a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // LSR ZN
				{Kind: BcByte, A0: 0x46},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // LSR ZX
				{Kind: BcByte, A0: 0x56},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // LSR AN
				{Kind: BcByte, A0: 0x4e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // LSR AX
				{Kind: BcByte, A0: 0x5e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwROL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // ROL A
				{Kind: BcByte, A0: 0x2a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // ROL ZN
				{Kind: BcByte, A0: 0x26},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // ROL ZX
				{Kind: BcByte, A0: 0x36},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // ROL AN
				{Kind: BcByte, A0: 0x2e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // ROL AX
				{Kind: BcByte, A0: 0x3e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwROR: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // ROR A
				{Kind: BcByte, A0: 0x6a},
			},
		},
		kwMemZN: InstPat{
			nil: InstDat{ // ROR ZN
				{Kind: BcByte, A0: 0x66},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemZX: InstPat{
			nil: InstDat{ // ROR ZX
				{Kind: BcByte, A0: 0x76},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemAN: InstPat{
			nil: InstDat{ // ROR AN
				{Kind: BcByte, A0: 0x6e},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemAX: InstPat{
			nil: InstDat{ // ROR AX
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
			kwCondVC: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondVS: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondPL: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondMI: {
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
			kwCondVC: {
				{KwReturn, &Operand{Kind: kwCondVS}},
			},
			kwCondVS: {
				{KwReturn, &Operand{Kind: kwCondVC}},
			},
			kwCondPL: {
				{KwReturn, &Operand{Kind: kwCondMI}},
			},
			kwCondMI: {
				{KwReturn, &Operand{Kind: kwCondPL}},
			},
		},
	},
}
