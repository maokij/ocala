package z80

import . "ocala/internal/core" //lint:ignore ST1001 core

var bmaps = [][]byte{
	{3, 0, 2, 70, 86, 94, 0}, // 0: IM
}

var kwADC = Intern("ADC")
var kwADD = Intern("ADD")
var kwAND = Intern("AND")
var kwAltAF = Intern("AF-")
var kwBIT = Intern("BIT")
var kwCALL = Intern("CALL")
var kwCCF = Intern("CCF")
var kwCP = Intern("CP")
var kwCPD = Intern("CPD")
var kwCPDR = Intern("CPDR")
var kwCPI = Intern("CPI")
var kwCPIR = Intern("CPIR")
var kwCPL = Intern("CPL")
var kwCondC = Intern("C?")
var kwCondM = Intern("M?")
var kwCondNC = Intern("NC?")
var kwCondNZ = Intern("NZ?")
var kwCondP = Intern("P?")
var kwCondPE = Intern("PE?")
var kwCondPO = Intern("PO?")
var kwCondZ = Intern("Z?")
var kwDAA = Intern("DAA")
var kwDEC = Intern("DEC")
var kwDI = Intern("DI")
var kwDJNZ = Intern("DJNZ")
var kwEI = Intern("EI")
var kwEX = Intern("EX")
var kwEXX = Intern("EXX")
var kwHALT = Intern("HALT")
var kwIM = Intern("IM")
var kwIN = Intern("IN")
var kwINC = Intern("INC")
var kwIND = Intern("IND")
var kwINDR = Intern("INDR")
var kwINI = Intern("INI")
var kwINIR = Intern("INIR")
var kwImmN = Intern("%B")
var kwImmNN = Intern("%W")
var kwJP = Intern("JP")
var kwJR = Intern("JR")
var kwLD = Intern("LD")
var kwLDD = Intern("LDD")
var kwLDDR = Intern("LDDR")
var kwLDI = Intern("LDI")
var kwLDIR = Intern("LDIR")
var kwLDP = Intern("#.LDP")
var kwMULUB = Intern("MULUB")
var kwMULUW = Intern("MULUW")
var kwMemBC = Intern("[BC]")
var kwMemC = Intern("[C]")
var kwMemDE = Intern("[DE]")
var kwMemHL = Intern("[HL]")
var kwMemIX = Intern("[IX %B]")
var kwMemIY = Intern("[IY %B]")
var kwMemN = Intern("[%B]")
var kwMemNN = Intern("[%W]")
var kwMemSP = Intern("[SP]")
var kwNEG = Intern("NEG")
var kwNOP = Intern("NOP")
var kwOR = Intern("OR")
var kwOTDR = Intern("OTDR")
var kwOTIR = Intern("OTIR")
var kwOUT = Intern("OUT")
var kwOUTD = Intern("OUTD")
var kwOUTI = Intern("OUTI")
var kwPOP = Intern("POP")
var kwPUSH = Intern("PUSH")
var kwRES = Intern("RES")
var kwRET = Intern("RET")
var kwRETI = Intern("RETI")
var kwRETN = Intern("RETN")
var kwRL = Intern("RL")
var kwRLA = Intern("RLA")
var kwRLC = Intern("RLC")
var kwRLCA = Intern("RLCA")
var kwRLD = Intern("RLD")
var kwRR = Intern("RR")
var kwRRA = Intern("RRA")
var kwRRC = Intern("RRC")
var kwRRCA = Intern("RRCA")
var kwRRD = Intern("RRD")
var kwRST = Intern("RST")
var kwRegA = Intern("A")
var kwRegAF = Intern("AF")
var kwRegB = Intern("B")
var kwRegBC = Intern("BC")
var kwRegC = Intern("C")
var kwRegD = Intern("D")
var kwRegDE = Intern("DE")
var kwRegE = Intern("E")
var kwRegF = Intern("F")
var kwRegH = Intern("H")
var kwRegHL = Intern("HL")
var kwRegI = Intern("I")
var kwRegIX = Intern("IX")
var kwRegIXH = Intern("IXH")
var kwRegIXL = Intern("IXL")
var kwRegIY = Intern("IY")
var kwRegIYH = Intern("IYH")
var kwRegIYL = Intern("IYL")
var kwRegL = Intern("L")
var kwRegPQ = Intern("PQ")
var kwRegR = Intern("R")
var kwRegSP = Intern("SP")
var kwSBC = Intern("SBC")
var kwSCF = Intern("SCF")
var kwSET = Intern("SET")
var kwSLA = Intern("SLA")
var kwSLL = Intern("SLL")
var kwSRA = Intern("SRA")
var kwSRL = Intern("SRL")
var kwSUB = Intern("SUB")
var kwXOR = Intern("XOR")

var asmOperands = map[*Keyword]AsmOperand{
	kwRegA:   {"A", false},
	kwRegB:   {"B", false},
	kwRegC:   {"C", false},
	kwRegD:   {"D", false},
	kwRegE:   {"E", false},
	kwRegH:   {"H", false},
	kwRegL:   {"L", false},
	kwRegHL:  {"HL", false},
	kwMemHL:  {"(HL)", false},
	kwRegBC:  {"BC", false},
	kwMemBC:  {"(BC)", false},
	kwRegDE:  {"DE", false},
	kwMemDE:  {"(DE)", false},
	kwRegAF:  {"AF", false},
	kwAltAF:  {"AF'", false},
	kwRegSP:  {"SP", false},
	kwMemSP:  {"(SP)", false},
	kwRegPQ:  {"PQ", false},
	kwRegIX:  {"IX", false},
	kwMemIX:  {"(IX+%)", true},
	kwRegIY:  {"IY", false},
	kwMemIY:  {"(IY+%)", true},
	kwImmN:   {"0+ %", true},
	kwMemN:   {"(%)", true},
	kwImmNN:  {"0+ %", true},
	kwMemNN:  {"(%)", true},
	kwMemC:   {"(C)", false},
	kwRegI:   {"I", false},
	kwRegR:   {"R", false},
	kwRegF:   {"F", false},
	kwCondNZ: {"NZ", false},
	kwCondZ:  {"Z", false},
	kwCondNC: {"NC", false},
	kwCondC:  {"C", false},
	kwCondPO: {"PO", false},
	kwCondPE: {"PE", false},
	kwCondP:  {"P", false},
	kwCondM:  {"M", false},
}

var tokenWords = [][]string{
	{"A", "B", "C", "D", "E", "H", "L", "I", "R", "F", "AF", "AF-", "BC", "DE", "HL", "IX", "IY", "SP"},
	{"NZ?", "Z?", "NC?", "C?", "PO?", "PE?", "P?", "M?"},
	{"-push", "-pop", "++", "--", "-not", "-neg", "-zero?", "-jump"},
	{"<-", "->", "<->", "+$", "-$", "-?", "<*", "<*$", ">*", ">*$", "-set", "-reset", "-bit?", "-in", "-out", "-jump-if", "-jump-unless"},
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

var instMap = InstPat{
	kwLD: InstPat{
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
	kwPUSH: InstPat{
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
	kwPOP: InstPat{
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
	kwEX: InstPat{
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
	kwEXX: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xd9},
		},
	},
	kwLDI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa0},
		},
	},
	kwLDIR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb0},
		},
	},
	kwLDD: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa8},
		},
	},
	kwLDDR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb8},
		},
	},
	kwCPI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa1},
		},
	},
	kwCPIR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb1},
		},
	},
	kwCPD: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa9},
		},
	},
	kwCPDR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb9},
		},
	},
	kwADD: InstPat{
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
	kwADC: InstPat{
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
	kwSUB: InstPat{
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
	kwSBC: InstPat{
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
	kwAND: InstPat{
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
	kwOR: InstPat{
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
	kwXOR: InstPat{
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
	kwCP: InstPat{
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
	kwINC: InstPat{
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
	kwDEC: InstPat{
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
	kwRLCA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x07},
		},
	},
	kwRLA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x17},
		},
	},
	kwRRCA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x0f},
		},
	},
	kwRRA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x1f},
		},
	},
	kwRLC: InstPat{
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
	kwRL: InstPat{
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
	kwRRC: InstPat{
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
	kwRR: InstPat{
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
	kwSLA: InstPat{
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
	kwSRA: InstPat{
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
	kwSRL: InstPat{
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
	kwRLD: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x6f},
		},
	},
	kwRRD: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x67},
		},
	},
	kwBIT: InstPat{
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
	kwSET: InstPat{
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
	kwRES: InstPat{
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
	kwJP: InstPat{
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
				{Kind: BcImp, A0: 0x00, A1: 0xe9, A2: 0x00, A3: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcImp, A0: 0x00, A1: 0xe9, A2: 0x00, A3: 0x00},
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
	kwJR: InstPat{
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
	kwDJNZ: InstPat{
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
	kwCALL: InstPat{
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
	kwRET: InstPat{
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
	kwRETI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x4d},
		},
	},
	kwRETN: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x45},
		},
	},
	kwRST: InstPat{
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
	KwJump: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	KwCall: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	kwIN: InstPat{
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
	kwINI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa2},
		},
	},
	kwINIR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb2},
		},
	},
	kwIND: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xaa},
		},
	},
	kwINDR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xba},
		},
	},
	kwOUT: InstPat{
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
	kwOUTI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa3},
		},
	},
	kwOTIR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb3},
		},
	},
	kwOUTD: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xab},
		},
	},
	kwOTDR: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xbb},
		},
	},
	kwDAA: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x27},
		},
	},
	kwCPL: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x2f},
		},
	},
	kwNEG: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x44},
		},
	},
	kwCCF: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x3f},
		},
	},
	kwSCF: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x37},
		},
	},
	kwNOP: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x00},
		},
	},
	kwHALT: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0x76},
		},
	},
	kwDI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xf3},
		},
	},
	kwEI: InstPat{
		nil: InstDat{
			{Kind: BcByte, A0: 0xfb},
		},
	},
	kwIM: InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xed},
				{Kind: BcMap, A0: 0x00, A1: 0x00},
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
		kwRegC: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegD: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegE: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegH: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegL: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegI: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegR: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegBC: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegDE: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegHL: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegSP: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIX: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIY: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("->"): {
		kwRegA: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegB: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegC: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegD: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegE: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegH: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegL: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegI: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegR: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegBC: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegDE: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegHL: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegSP: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIX: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIY: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwImmNN: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("<->"): {
		kwRegAF: {
			kwAltAF: {
				{kwEX, &Operand{Kind: kwRegAF}, &Operand{Kind: kwAltAF}},
			},
		},
		kwRegDE: {
			kwRegHL: {
				{kwEX, &Operand{Kind: kwRegDE}, &Operand{Kind: kwRegHL}},
			},
		},
		kwRegHL: {
			kwRegDE: {
				{kwEX, &Operand{Kind: kwRegDE}, &Operand{Kind: kwRegHL}},
			},
			kwMemSP: {
				{kwEX, &Operand{Kind: kwMemSP}, &Operand{Kind: kwRegHL}},
			},
		},
		kwMemSP: {
			kwRegHL: {
				{kwEX, &Operand{Kind: kwMemSP}, &Operand{Kind: kwRegHL}},
			},
		},
	},
	Intern("-push"): {
		KwAny: {
			nil: {
				{kwPUSH, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-pop"): {
		KwAny: {
			nil: {
				{kwPOP, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("++"): {
		KwAny: {
			nil: {
				{kwINC, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("--"): {
		KwAny: {
			nil: {
				{kwDEC, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-not"): {
		kwRegA: {
			nil: {
				{kwCPL},
			},
		},
	},
	Intern("-neg"): {
		kwRegA: {
			nil: {
				{kwNEG},
			},
		},
	},
	Intern("+"): {
		KwAny: {
			KwAny: {
				{kwADD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("+$"): {
		KwAny: {
			KwAny: {
				{kwADC, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-"): {
		kwRegHL: {
			KwAny: {
				{kwOR, &Operand{Kind: kwRegA}},
				{kwSBC, &Operand{Kind: kwRegHL}, &Vec{Int(1), nil}},
			},
		},
		kwRegA: {
			KwAny: {
				{kwSUB, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-$"): {
		KwAny: {
			KwAny: {
				{kwSBC, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-?"): {
		kwRegA: {
			KwAny: {
				{kwCP, &Vec{Int(1), nil}},
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
				{kwOR, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("^"): {
		kwRegA: {
			KwAny: {
				{kwXOR, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("<*"): {
		kwRegA: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRLCA},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRLC, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("<*$"): {
		kwRegA: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRLA},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRL, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">*"): {
		kwRegA: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRRCA},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRRC, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">*$"): {
		kwRegA: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRRA},
				}},
			},
		},
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwRR, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("<<"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwSLA, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">>"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwSRA, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern(">>>"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwSRL, &Vec{Int(0), nil}},
				}},
			},
		},
	},
	Intern("-set"): {
		KwAny: {
			KwAny: {
				{kwSET, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-reset"): {
		KwAny: {
			KwAny: {
				{kwRES, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-bit?"): {
		KwAny: {
			KwAny: {
				{kwBIT, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-in"): {
		kwRegA: {
			kwRegC: {
				{kwIN, &Operand{Kind: kwRegA}, &Operand{Kind: kwMemC}},
			},
			kwImmNN: {
				{kwIN, &Operand{Kind: kwRegA}, &Vec{Int(1), kwMemNN}},
			},
		},
		KwAny: {
			kwRegC: {
				{kwIN, &Vec{Int(0), nil}, &Operand{Kind: kwMemC}},
			},
		},
	},
	Intern("-out"): {
		kwRegA: {
			kwRegC: {
				{kwOUT, &Operand{Kind: kwMemC}, &Operand{Kind: kwRegA}},
			},
			kwImmNN: {
				{kwOUT, &Vec{Int(1), kwMemNN}, &Operand{Kind: kwRegA}},
			},
		},
		KwAny: {
			kwRegC: {
				{kwOUT, &Operand{Kind: kwMemC}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("-zero?"): {
		kwRegA: {
			nil: {
				{kwAND, &Vec{Int(0), nil}},
			},
		},
		kwRegBC: {
			nil: {
				{KwINVALID, &Vec{Int(0), nil}},
			},
		},
		kwRegDE: {
			nil: {
				{KwINVALID, &Vec{Int(0), nil}},
			},
		},
		kwRegHL: {
			nil: {
				{KwINVALID, &Vec{Int(0), nil}},
			},
		},
		kwRegSP: {
			nil: {
				{KwINVALID, &Vec{Int(0), nil}},
			},
		},
		kwRegIX: {
			nil: {
				{KwINVALID, &Vec{Int(0), nil}},
			},
		},
		kwRegIY: {
			nil: {
				{KwINVALID, &Vec{Int(0), nil}},
			},
		},
		KwAny: {
			nil: {
				{kwINC, &Vec{Int(0), nil}},
				{kwDEC, &Vec{Int(0), nil}},
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
			kwCondNZ: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondZ: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondNC: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondC: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondPO: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondPE: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondP: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwCondM: {
				{KwJump, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-jump-unless"): {
		kwImmNN: {
			kwCondNZ: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondZ}},
			},
			kwCondZ: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondNZ}},
			},
			kwCondNC: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondC}},
			},
			kwCondC: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondNC}},
			},
			kwCondPO: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondPE}},
			},
			kwCondPE: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondPO}},
			},
			kwCondM: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondP}},
			},
			kwCondP: {
				{KwJump, &Vec{Int(0), nil}, &Operand{Kind: kwCondM}},
			},
		},
	},
}

var asmOperandsUndocumented = map[*Keyword]AsmOperand{
	kwRegIXH: {"IXH", false},
	kwRegIXL: {"IXL", false},
	kwRegIYH: {"IYH", false},
	kwRegIYL: {"IYL", false},
}

var tokenWordsUndocumented = [][]string{
	{"IXH", "IXL", "IYH", "IYL"},
	{},
	{},
	{"<<<"},
}

var instMapUndocumented = InstPat{
	kwLD: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
		},
		kwRegB: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x45},
				},
			},
		},
		kwRegC: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
		},
		kwRegD: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x55},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x55},
				},
			},
		},
		kwRegE: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
		},
		kwRegIXH: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIXL: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYH: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYL: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwADD: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x85},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x85},
				},
			},
		},
	},
	kwADC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
		},
	},
	kwSUB: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x95},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x95},
			},
		},
	},
	kwSBC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
		},
	},
	kwAND: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
	},
	kwOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
	},
	kwXOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xad},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xad},
			},
		},
	},
	kwCP: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
	},
	kwINC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
	},
	kwDEC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
	},
	kwRLC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x07},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x07},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x01},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x02},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x02},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x03},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x04},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x04},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x05},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x05},
				},
			},
		},
	},
	kwRL: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x17},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x17},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x10},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x10},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x11},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x11},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x12},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x12},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x13},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x13},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x14},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x14},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x15},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x15},
				},
			},
		},
	},
	kwRRC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0f},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x08},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x08},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x09},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0a},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0b},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0c},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0d},
				},
			},
		},
	},
	kwRR: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1f},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x18},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x18},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x19},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1a},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1b},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1c},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1d},
				},
			},
		},
	},
	kwSLA: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x27},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x27},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x20},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x20},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x21},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x21},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x22},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x22},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x23},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x23},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x24},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x24},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x25},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x25},
				},
			},
		},
	},
	kwSLL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x37},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x37},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x37},
				},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x30},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x30},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x30},
				},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x31},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x31},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x31},
				},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x32},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x32},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x32},
				},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x33},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x33},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x33},
				},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x34},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x34},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x34},
				},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x35},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x35},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x35},
				},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x36},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x36},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x36},
			},
		},
	},
	kwSRA: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2f},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x28},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x28},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x29},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2a},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2b},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2c},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2d},
				},
			},
		},
	},
	kwSRL: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3f},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x38},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x38},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x39},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x39},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3a},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3b},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3c},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3d},
				},
			},
		},
	},
	kwBIT: InstPat{
		kwRegA: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegB: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegC: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegD: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegE: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegH: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegL: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwImmN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
			},
		},
		kwImmNN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
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
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
		},
	},
	kwSET: InstPat{
		kwRegA: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegB: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegC: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegD: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegE: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegH: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegL: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwImmN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
			},
		},
		kwImmNN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
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
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
		},
	},
	kwRES: InstPat{
		kwRegA: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegB: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegC: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegD: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegE: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegH: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwRegL: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwImmNN: InstPat{
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
		kwImmN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
			},
		},
		kwImmNN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
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
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
		},
	},
	kwIN: InstPat{
		kwMemC: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xed},
				{Kind: BcByte, A0: 0x70},
			},
		},
		kwRegF: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x70},
				},
			},
		},
	},
	kwOUT: InstPat{
		kwMemC: InstPat{
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcImp, A0: 0x01, A1: 0x71, A2: 0x00, A3: 0x00},
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

var ctxOpMapUndocumented = CtxOpMap{
	Intern("<-"): {
		kwRegIXH: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIXL: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIYH: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIYL: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("->"): {
		kwRegIXH: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIXL: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIYH: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIYL: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
	Intern("<<<"): {
		KwAny: {
			kwImmNN: {
				{KwREP, &Vec{Int(1), nil}, &Vec{
					&Vec{kwSLL, &Vec{Int(0), nil}},
				}},
			},
		},
	},
}

var instMapCompat8080 = InstPat{
	kwLD: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegI: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegR: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemIX: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemIY: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegI: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegR: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIX: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIY: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegBC: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegDE: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegSP: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemN: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwPUSH: InstPat{
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwPOP: InstPat{
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwEX: InstPat{
		kwRegAF: InstPat{
			kwAltAF: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemSP: InstPat{
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwEXX: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDI: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDIR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDD: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDDR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPI: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPIR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPD: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPDR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwADD: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIX: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIY: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwADC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwSUB: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSBC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwAND: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwOR: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwXOR: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwCP: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwINC: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwDEC: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRLC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRRC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRR: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSLA: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSRA: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSRL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRLD: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwRRD: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwBIT: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwSET: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwRES: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwJP: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwJR: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwDJNZ: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRETI: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwRETN: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwIN: InstPat{
		kwRegA: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwINI: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwINIR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwIND: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwINDR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOUT: InstPat{
		kwMemC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwOUTI: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOTIR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOUTD: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOTDR: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwNEG: InstPat{
		nil: InstDat{
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwIM: InstPat{
		kwImmN: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
}

var asmOperandsR800 = map[*Keyword]AsmOperand{
	kwRegIXH: {"IXH", false},
	kwRegIXL: {"IXL", false},
	kwRegIYH: {"IYH", false},
	kwRegIYL: {"IYL", false},
}

var tokenWordsR800 = [][]string{
	{"IXH", "IXL", "IYH", "IYL"},
	{},
	{},
	{},
}

var instMapR800 = InstPat{
	kwLD: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
		},
		kwRegB: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x45},
				},
			},
		},
		kwRegC: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
		},
		kwRegD: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x55},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x55},
				},
			},
		},
		kwRegE: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
		},
		kwRegIXH: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIXL: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYH: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYL: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwADD: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x85},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x85},
				},
			},
		},
	},
	kwADC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
		},
	},
	kwSUB: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x95},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x95},
			},
		},
	},
	kwSBC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
		},
	},
	kwAND: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
	},
	kwOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
	},
	kwXOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xad},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xad},
			},
		},
	},
	kwCP: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
	},
	kwINC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
	},
	kwDEC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
	},
	kwIN: InstPat{
		kwRegF: InstPat{
			kwMemC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x70},
				},
			},
		},
	},
	kwMULUB: InstPat{
		kwRegA: InstPat{
			kwRegB: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xc1},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xc9},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xd1},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xd9},
				},
			},
		},
	},
	kwMULUW: InstPat{
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xc3},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xf3},
				},
			},
		},
	},
}

var ctxOpMapR800 = CtxOpMap{
	Intern("<-"): {
		kwRegIXH: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIXL: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIYH: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegIYL: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("->"): {
		kwRegIXH: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIXL: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIYH: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegIYL: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
	},
}
