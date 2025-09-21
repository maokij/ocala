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
var kwMemWX = Intern("[IX %W]")
var kwMemWY = Intern("[IY %W]")
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
var kwRegPC = Intern("PC")
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
	kwRegPC:  {"PC", false},
	kwRegPQ:  {"PQ", false},
	kwRegIX:  {"IX", false},
	kwMemIX:  {"(IX+%)", true},
	kwMemWX:  {"(IX+%)", true},
	kwRegIY:  {"IY", false},
	kwMemIY:  {"(IY+%)", true},
	kwMemWY:  {"(IY+%)", true},
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
	{"A", "B", "C", "D", "E", "H", "L", "I", "R", "F", "AF", "AF-", "BC", "DE", "HL", "IX", "IY", "SP", "PC"},
	{"NZ?", "Z?", "NC?", "C?", "PO?", "PE?", "P?", "M?"},
	{"-push", "-pop", "++", "--", "-not", "-neg", "-zero?", "-jump", "-return"},
	{"<-", "->", "<->", "+$", "-$", "-?", "<*", "<*$", ">*", ">*$", "-set", "-reset", "-bit?", "-in", "-out", "-jump-if", "-jump-unless", "-return-if", "-return-unless"},
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
				nil: InstDat{ // LD A A
					{Kind: BcByte, A0: 0x7f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD A B
					{Kind: BcByte, A0: 0x78},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD A C
					{Kind: BcByte, A0: 0x79},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD A D
					{Kind: BcByte, A0: 0x7a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD A E
					{Kind: BcByte, A0: 0x7b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD A H
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD A L
					{Kind: BcByte, A0: 0x7d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD A N
					{Kind: BcByte, A0: 0x3e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // LD A HL$
					{Kind: BcByte, A0: 0x7e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // LD A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemBC: InstPat{
				nil: InstDat{ // LD A BC$
					{Kind: BcByte, A0: 0x0a},
				},
			},
			kwMemDE: InstPat{
				nil: InstDat{ // LD A DE$
					{Kind: BcByte, A0: 0x1a},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD A NN$
					{Kind: BcByte, A0: 0x3a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwRegI: InstPat{
				nil: InstDat{ // LD A I
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x57},
				},
			},
			kwRegR: InstPat{
				nil: InstDat{ // LD A R
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5f},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD A NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD A N$
					{Kind: BcByte, A0: 0x3a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegB: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD B A
					{Kind: BcByte, A0: 0x47},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD B B
					{Kind: BcByte, A0: 0x40},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD B C
					{Kind: BcByte, A0: 0x41},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD B D
					{Kind: BcByte, A0: 0x42},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD B E
					{Kind: BcByte, A0: 0x43},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD B H
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD B L
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD B N
					{Kind: BcByte, A0: 0x06},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // LD B HL$
					{Kind: BcByte, A0: 0x46},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // LD B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x46},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x46},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD B NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD C A
					{Kind: BcByte, A0: 0x4f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD C B
					{Kind: BcByte, A0: 0x48},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD C C
					{Kind: BcByte, A0: 0x49},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD C D
					{Kind: BcByte, A0: 0x4a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD C E
					{Kind: BcByte, A0: 0x4b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD C H
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD C L
					{Kind: BcByte, A0: 0x4d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD C N
					{Kind: BcByte, A0: 0x0e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // LD C HL$
					{Kind: BcByte, A0: 0x4e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // LD C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD C NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD D A
					{Kind: BcByte, A0: 0x57},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD D B
					{Kind: BcByte, A0: 0x50},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD D C
					{Kind: BcByte, A0: 0x51},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD D D
					{Kind: BcByte, A0: 0x52},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD D E
					{Kind: BcByte, A0: 0x53},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD D H
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD D L
					{Kind: BcByte, A0: 0x55},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD D N
					{Kind: BcByte, A0: 0x16},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // LD D HL$
					{Kind: BcByte, A0: 0x56},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // LD D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x56},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x56},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD D NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD E A
					{Kind: BcByte, A0: 0x5f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD E B
					{Kind: BcByte, A0: 0x58},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD E C
					{Kind: BcByte, A0: 0x59},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD E D
					{Kind: BcByte, A0: 0x5a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD E E
					{Kind: BcByte, A0: 0x5b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD E H
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD E L
					{Kind: BcByte, A0: 0x5d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD E N
					{Kind: BcByte, A0: 0x1e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // LD E HL$
					{Kind: BcByte, A0: 0x5e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // LD E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD E NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD H A
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD H B
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD H C
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD H D
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD H E
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD H H
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD H L
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD H N
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // LD H HL$
					{Kind: BcByte, A0: 0x66},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // LD H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x66},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x66},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD H NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD L A
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD L B
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD L C
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD L D
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD L E
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD L H
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD L L
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD L N
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // LD L HL$
					{Kind: BcByte, A0: 0x6e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // LD L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD L NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemHL: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD HL$ A
					{Kind: BcByte, A0: 0x77},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD HL$ B
					{Kind: BcByte, A0: 0x70},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD HL$ C
					{Kind: BcByte, A0: 0x71},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD HL$ D
					{Kind: BcByte, A0: 0x72},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD HL$ E
					{Kind: BcByte, A0: 0x73},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD HL$ H
					{Kind: BcByte, A0: 0x74},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD HL$ L
					{Kind: BcByte, A0: 0x75},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD HL$ N
					{Kind: BcByte, A0: 0x36},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD HL$ NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemIX: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD IX$ A
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x77},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IX$ B
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x70},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IX$ C
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x71},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IX$ D
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x72},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IX$ E
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD IX$ H
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x74},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD IX$ L
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x75},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IX$ N
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x36},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IX$ NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemIY: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD IY$ A
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x77},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IY$ B
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x70},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IY$ C
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x71},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IY$ D
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x72},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IY$ E
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD IY$ H
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x74},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD IY$ L
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x75},
					{Kind: BcLow, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IY$ N
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x36},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IY$ NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemBC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD BC$ A
					{Kind: BcByte, A0: 0x02},
				},
			},
		},
		kwMemDE: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD DE$ A
					{Kind: BcByte, A0: 0x12},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD NN$ A
					{Kind: BcByte, A0: 0x32},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // LD NN$ HL
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegBC: InstPat{
				nil: InstDat{ // LD NN$ BC
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x43},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // LD NN$ DE
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x53},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // LD NN$ SP
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // LD NN$ IX
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // LD NN$ IY
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwRegI: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD I A
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x47},
				},
			},
		},
		kwRegR: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD R A
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4f},
				},
			},
		},
		kwRegBC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD BC NN
					{Kind: BcByte, A0: 0x01},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD BC NN$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD BC N
					{Kind: BcByte, A0: 0x01},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD BC N$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegDE: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD DE NN
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD DE NN$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD DE N
					{Kind: BcByte, A0: 0x11},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD DE N$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegHL: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD HL NN
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD HL NN$
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD HL N
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD HL N$
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegSP: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD SP NN
					{Kind: BcByte, A0: 0x31},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD SP NN$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x7b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // LD SP HL
					{Kind: BcByte, A0: 0xf9},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // LD SP IX
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xf9},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // LD SP IY
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xf9},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD SP N
					{Kind: BcByte, A0: 0x31},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD SP N$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x7b},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegIX: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD IX NN
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD IX NN$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IX N
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD IX N$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwRegIY: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD IY NN
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD IY NN$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IY N
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x21},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD IY N$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2a},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwMemWX: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD WX$ A
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD WX$ B
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD WX$ C
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD WX$ D
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD WX$ E
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD WX$ H
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD WX$ L
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD WX$ N
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD WX$ NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemWY: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD WY$ A
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD WY$ B
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD WY$ C
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD WY$ D
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD WY$ E
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD WY$ H
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD WY$ L
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD WY$ N
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD WY$ NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD N$ A
					{Kind: BcByte, A0: 0x32},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // LD N$ HL
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegBC: InstPat{
				nil: InstDat{ // LD N$ BC
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x43},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // LD N$ DE
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x53},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // LD N$ SP
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x73},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // LD N$ IX
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x22},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // LD N$ IY
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
			nil: InstDat{ // PUSH BC
				{Kind: BcByte, A0: 0xc5},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{ // PUSH DE
				{Kind: BcByte, A0: 0xd5},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{ // PUSH HL
				{Kind: BcByte, A0: 0xe5},
			},
		},
		kwRegAF: InstPat{
			nil: InstDat{ // PUSH AF
				{Kind: BcByte, A0: 0xf5},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{ // PUSH IX
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xe5},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // PUSH IY
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xe5},
			},
		},
	},
	kwPOP: InstPat{
		kwRegBC: InstPat{
			nil: InstDat{ // POP BC
				{Kind: BcByte, A0: 0xc1},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{ // POP DE
				{Kind: BcByte, A0: 0xd1},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{ // POP HL
				{Kind: BcByte, A0: 0xe1},
			},
		},
		kwRegAF: InstPat{
			nil: InstDat{ // POP AF
				{Kind: BcByte, A0: 0xf1},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{ // POP IX
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xe1},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // POP IY
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xe1},
			},
		},
	},
	kwEX: InstPat{
		kwRegDE: InstPat{
			kwRegHL: InstPat{
				nil: InstDat{ // EX DE HL
					{Kind: BcByte, A0: 0xeb},
				},
			},
		},
		kwRegAF: InstPat{
			kwAltAF: InstPat{
				nil: InstDat{ // EX AF AF-
					{Kind: BcByte, A0: 0x08},
				},
			},
		},
		kwMemSP: InstPat{
			kwRegHL: InstPat{
				nil: InstDat{ // EX SP$ HL
					{Kind: BcByte, A0: 0xe3},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // EX SP$ IX
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xe3},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // EX SP$ IY
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xe3},
				},
			},
		},
	},
	kwEXX: InstPat{
		nil: InstDat{ // EXX
			{Kind: BcByte, A0: 0xd9},
		},
	},
	kwLDI: InstPat{
		nil: InstDat{ // LDI
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa0},
		},
	},
	kwLDIR: InstPat{
		nil: InstDat{ // LDIR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb0},
		},
	},
	kwLDD: InstPat{
		nil: InstDat{ // LDD
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa8},
		},
	},
	kwLDDR: InstPat{
		nil: InstDat{ // LDDR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb8},
		},
	},
	kwCPI: InstPat{
		nil: InstDat{ // CPI
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa1},
		},
	},
	kwCPIR: InstPat{
		nil: InstDat{ // CPIR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb1},
		},
	},
	kwCPD: InstPat{
		nil: InstDat{ // CPD
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa9},
		},
	},
	kwCPDR: InstPat{
		nil: InstDat{ // CPDR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb9},
		},
	},
	kwADD: InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // ADD A A
					{Kind: BcByte, A0: 0x87},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // ADD A B
					{Kind: BcByte, A0: 0x80},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // ADD A C
					{Kind: BcByte, A0: 0x81},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // ADD A D
					{Kind: BcByte, A0: 0x82},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // ADD A E
					{Kind: BcByte, A0: 0x83},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // ADD A H
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // ADD A L
					{Kind: BcByte, A0: 0x85},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // ADD A N
					{Kind: BcByte, A0: 0xc6},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // ADD A HL$
					{Kind: BcByte, A0: 0x86},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // ADD A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x86},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // ADD A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x86},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // ADD A NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // ADD A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // ADD A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // ADD HL BC
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // ADD HL DE
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // ADD HL HL
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // ADD HL SP
					{Kind: BcByte, A0: 0x39},
				},
			},
		},
		kwRegIX: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // ADD IX BC
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // ADD IX DE
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // ADD IX IX
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // ADD IX SP
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x39},
				},
			},
		},
		kwRegIY: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // ADD IY BC
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // ADD IY DE
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // ADD IY IY
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // ADD IY SP
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x39},
				},
			},
		},
	},
	kwADC: InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // ADC A A
					{Kind: BcByte, A0: 0x8f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // ADC A B
					{Kind: BcByte, A0: 0x88},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // ADC A C
					{Kind: BcByte, A0: 0x89},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // ADC A D
					{Kind: BcByte, A0: 0x8a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // ADC A E
					{Kind: BcByte, A0: 0x8b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // ADC A H
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // ADC A L
					{Kind: BcByte, A0: 0x8d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // ADC A N
					{Kind: BcByte, A0: 0xce},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // ADC A HL$
					{Kind: BcByte, A0: 0x8e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // ADC A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // ADC A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // ADC A NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // ADC A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // ADC A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // ADC HL BC
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x4a},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // ADC HL DE
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x5a},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // ADC HL HL
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // ADC HL SP
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x7a},
				},
			},
		},
	},
	kwSUB: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SUB A
				{Kind: BcByte, A0: 0x97},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SUB B
				{Kind: BcByte, A0: 0x90},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SUB C
				{Kind: BcByte, A0: 0x91},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SUB D
				{Kind: BcByte, A0: 0x92},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SUB E
				{Kind: BcByte, A0: 0x93},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SUB H
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SUB L
				{Kind: BcByte, A0: 0x95},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // SUB N
				{Kind: BcByte, A0: 0xd6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SUB HL$
				{Kind: BcByte, A0: 0x96},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SUB IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x96},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SUB IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x96},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // SUB NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SUB WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SUB WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwSBC: InstPat{
		kwRegA: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // SBC A A
					{Kind: BcByte, A0: 0x9f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // SBC A B
					{Kind: BcByte, A0: 0x98},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // SBC A C
					{Kind: BcByte, A0: 0x99},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // SBC A D
					{Kind: BcByte, A0: 0x9a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // SBC A E
					{Kind: BcByte, A0: 0x9b},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // SBC A H
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // SBC A L
					{Kind: BcByte, A0: 0x9d},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // SBC A N
					{Kind: BcByte, A0: 0xde},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // SBC A HL$
					{Kind: BcByte, A0: 0x9e},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SBC A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SBC A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // SBC A NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SBC A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SBC A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // SBC HL BC
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x42},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // SBC HL DE
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x52},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // SBC HL HL
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // SBC HL SP
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x72},
				},
			},
		},
	},
	kwAND: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // AND A
				{Kind: BcByte, A0: 0xa7},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // AND B
				{Kind: BcByte, A0: 0xa0},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // AND C
				{Kind: BcByte, A0: 0xa1},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // AND D
				{Kind: BcByte, A0: 0xa2},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // AND E
				{Kind: BcByte, A0: 0xa3},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // AND H
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // AND L
				{Kind: BcByte, A0: 0xa5},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // AND N
				{Kind: BcByte, A0: 0xe6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // AND HL$
				{Kind: BcByte, A0: 0xa6},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // AND IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // AND IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // AND NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // AND WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // AND WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwOR: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // OR A
				{Kind: BcByte, A0: 0xb7},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // OR B
				{Kind: BcByte, A0: 0xb0},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // OR C
				{Kind: BcByte, A0: 0xb1},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // OR D
				{Kind: BcByte, A0: 0xb2},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // OR E
				{Kind: BcByte, A0: 0xb3},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // OR H
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // OR L
				{Kind: BcByte, A0: 0xb5},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // OR N
				{Kind: BcByte, A0: 0xf6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // OR HL$
				{Kind: BcByte, A0: 0xb6},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // OR IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // OR IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb6},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // OR NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // OR WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // OR WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwXOR: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // XOR A
				{Kind: BcByte, A0: 0xaf},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // XOR B
				{Kind: BcByte, A0: 0xa8},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // XOR C
				{Kind: BcByte, A0: 0xa9},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // XOR D
				{Kind: BcByte, A0: 0xaa},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // XOR E
				{Kind: BcByte, A0: 0xab},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // XOR H
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // XOR L
				{Kind: BcByte, A0: 0xad},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // XOR N
				{Kind: BcByte, A0: 0xee},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // XOR HL$
				{Kind: BcByte, A0: 0xae},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // XOR IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xae},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // XOR IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xae},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // XOR NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // XOR WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // XOR WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwCP: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // CP A
				{Kind: BcByte, A0: 0xbf},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // CP B
				{Kind: BcByte, A0: 0xb8},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // CP C
				{Kind: BcByte, A0: 0xb9},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // CP D
				{Kind: BcByte, A0: 0xba},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // CP E
				{Kind: BcByte, A0: 0xbb},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // CP H
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // CP L
				{Kind: BcByte, A0: 0xbd},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // CP N
				{Kind: BcByte, A0: 0xfe},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // CP HL$
				{Kind: BcByte, A0: 0xbe},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // CP IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbe},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // CP IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbe},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // CP NN
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // CP WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // CP WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwINC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // INC A
				{Kind: BcByte, A0: 0x3c},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // INC B
				{Kind: BcByte, A0: 0x04},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // INC C
				{Kind: BcByte, A0: 0x0c},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // INC D
				{Kind: BcByte, A0: 0x14},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // INC E
				{Kind: BcByte, A0: 0x1c},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // INC H
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // INC L
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // INC HL$
				{Kind: BcByte, A0: 0x34},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // INC IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x34},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // INC IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x34},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwRegBC: InstPat{
			nil: InstDat{ // INC BC
				{Kind: BcByte, A0: 0x03},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{ // INC DE
				{Kind: BcByte, A0: 0x13},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{ // INC HL
				{Kind: BcByte, A0: 0x23},
			},
		},
		kwRegSP: InstPat{
			nil: InstDat{ // INC SP
				{Kind: BcByte, A0: 0x33},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{ // INC IX
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x23},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // INC IY
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x23},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // INC WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // INC WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwDEC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // DEC A
				{Kind: BcByte, A0: 0x3d},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // DEC B
				{Kind: BcByte, A0: 0x05},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // DEC C
				{Kind: BcByte, A0: 0x0d},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // DEC D
				{Kind: BcByte, A0: 0x15},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // DEC E
				{Kind: BcByte, A0: 0x1d},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // DEC H
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // DEC L
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // DEC HL$
				{Kind: BcByte, A0: 0x35},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // DEC IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x35},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // DEC IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x35},
				{Kind: BcLow, A0: 0x00},
			},
		},
		kwRegBC: InstPat{
			nil: InstDat{ // DEC BC
				{Kind: BcByte, A0: 0x0b},
			},
		},
		kwRegDE: InstPat{
			nil: InstDat{ // DEC DE
				{Kind: BcByte, A0: 0x1b},
			},
		},
		kwRegHL: InstPat{
			nil: InstDat{ // DEC HL
				{Kind: BcByte, A0: 0x2b},
			},
		},
		kwRegSP: InstPat{
			nil: InstDat{ // DEC SP
				{Kind: BcByte, A0: 0x3b},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{ // DEC IX
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2b},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // DEC IY
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2b},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // DEC WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // DEC WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwRLCA: InstPat{
		nil: InstDat{ // RLCA
			{Kind: BcByte, A0: 0x07},
		},
	},
	kwRLA: InstPat{
		nil: InstDat{ // RLA
			{Kind: BcByte, A0: 0x17},
		},
	},
	kwRRCA: InstPat{
		nil: InstDat{ // RRCA
			{Kind: BcByte, A0: 0x0f},
		},
	},
	kwRRA: InstPat{
		nil: InstDat{ // RRA
			{Kind: BcByte, A0: 0x1f},
		},
	},
	kwRLC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RLC A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x07},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RLC B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RLC C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x01},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RLC D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x02},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RLC E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x03},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RLC H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x04},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RLC L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x05},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RLC HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x06},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RLC IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x06},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RLC IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x06},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RLC WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RLC WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwRL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RL A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x17},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RL B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x10},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RL C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x11},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RL D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x12},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RL E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x13},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RL H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x14},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RL L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x15},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RL HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x16},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RL IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x16},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RL IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x16},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RL WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RL WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwRRC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RRC A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RRC B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x08},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RRC C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x09},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RRC D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RRC E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RRC H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RRC L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RRC HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x0e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RRC IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x0e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RRC IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x0e},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RRC WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RRC WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwRR: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RR A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RR B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x18},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RR C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x19},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RR D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RR E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RR H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RR L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RR HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x1e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RR IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x1e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RR IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x1e},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RR WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RR WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwSLA: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SLA A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x27},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SLA B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x20},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SLA C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x21},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SLA D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x22},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SLA E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x23},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SLA H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SLA L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SLA HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x26},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SLA IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x26},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SLA IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x26},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SLA WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SLA WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwSRA: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SRA A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SRA B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x28},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SRA C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x29},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SRA D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SRA E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SRA H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SRA L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SRA HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x2e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SRA IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x2e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SRA IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x2e},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SRA WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SRA WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwSRL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SRL A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3f},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SRL B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x38},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SRL C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x39},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SRL D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3a},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SRL E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3b},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SRL H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3c},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SRL L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3d},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SRL HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x3e},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SRL IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x3e},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SRL IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x3e},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SRL WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SRL WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwRLD: InstPat{
		nil: InstDat{ // RLD
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x6f},
		},
	},
	kwRRD: InstPat{
		nil: InstDat{ // RRD
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x67},
		},
	},
	kwBIT: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // BIT N A
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x47, A2: 0x07, A3: 0x03},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // BIT N B
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x40, A2: 0x07, A3: 0x03},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // BIT N C
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x41, A2: 0x07, A3: 0x03},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // BIT N D
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x42, A2: 0x07, A3: 0x03},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // BIT N E
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x43, A2: 0x07, A3: 0x03},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // BIT N H
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x44, A2: 0x07, A3: 0x03},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // BIT N L
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x45, A2: 0x07, A3: 0x03},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // BIT N HL$
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x46, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // BIT N IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x46, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // BIT N IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x46, A2: 0x07, A3: 0x03},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // BIT N WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // BIT N WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // BIT NN A
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // BIT NN B
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // BIT NN C
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // BIT NN D
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // BIT NN E
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // BIT NN H
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // BIT NN L
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // BIT NN HL$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // BIT NN IX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // BIT NN WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // BIT NN IY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // BIT NN WY$
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
				nil: InstDat{ // SET N A
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x07, A3: 0x03},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // SET N B
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc0, A2: 0x07, A3: 0x03},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // SET N C
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc1, A2: 0x07, A3: 0x03},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // SET N D
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc2, A2: 0x07, A3: 0x03},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // SET N E
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc3, A2: 0x07, A3: 0x03},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // SET N H
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc4, A2: 0x07, A3: 0x03},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // SET N L
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc5, A2: 0x07, A3: 0x03},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // SET N HL$
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0xc6, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SET N IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0xc6, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SET N IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0xc6, A2: 0x07, A3: 0x03},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SET N WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SET N WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // SET NN A
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // SET NN B
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // SET NN C
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // SET NN D
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // SET NN E
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // SET NN H
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // SET NN L
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // SET NN HL$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SET NN IX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SET NN WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SET NN IY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SET NN WY$
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
				nil: InstDat{ // RES N A
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x87, A2: 0x07, A3: 0x03},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // RES N B
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x80, A2: 0x07, A3: 0x03},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // RES N C
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x81, A2: 0x07, A3: 0x03},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // RES N D
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x82, A2: 0x07, A3: 0x03},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // RES N E
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x83, A2: 0x07, A3: 0x03},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // RES N H
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x84, A2: 0x07, A3: 0x03},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // RES N L
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x85, A2: 0x07, A3: 0x03},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // RES N HL$
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcImp, A0: 0x00, A1: 0x86, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // RES N IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x86, A2: 0x07, A3: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RES N IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcImp, A0: 0x00, A1: 0x86, A2: 0x07, A3: 0x03},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RES N WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RES N WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // RES NN A
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // RES NN B
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // RES NN C
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // RES NN D
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // RES NN E
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // RES NN H
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // RES NN L
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // RES NN HL$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // RES NN IX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RES NN WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RES NN IY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RES NN WY$
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
			nil: InstDat{ // JP NN
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // JP HL$
				{Kind: BcByte, A0: 0xe9},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // JP IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcImp, A0: 0x00, A1: 0xe9, A2: 0x00, A3: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // JP IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcImp, A0: 0x00, A1: 0xe9, A2: 0x00, A3: 0x00},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP NZ? NN
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP NZ? N
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP Z? NN
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP Z? N
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP NC? NN
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP NC? N
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP C? NN
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP C? N
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPO: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP PO? NN
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP PO? N
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPE: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP PE? NN
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP PE? N
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondP: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP P? NN
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP P? N
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondM: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JP M? NN
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JP M? N
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // JP N
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // JP WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // JP WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwJR: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // JR NN
				{Kind: BcByte, A0: 0x18},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR C? NN
					{Kind: BcByte, A0: 0x38},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR C? N
					{Kind: BcByte, A0: 0x38},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR NC? NN
					{Kind: BcByte, A0: 0x30},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR NC? N
					{Kind: BcByte, A0: 0x30},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR Z? NN
					{Kind: BcByte, A0: 0x28},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR Z? N
					{Kind: BcByte, A0: 0x28},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR NZ? NN
					{Kind: BcByte, A0: 0x20},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR NZ? N
					{Kind: BcByte, A0: 0x20},
					{Kind: BcRlow, A0: 0x01, A1: 0xfe, A2: 0x01},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // JR N
				{Kind: BcByte, A0: 0x18},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwDJNZ: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // DJNZ NN
				{Kind: BcByte, A0: 0x10},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // DJNZ N
				{Kind: BcByte, A0: 0x10},
				{Kind: BcRlow, A0: 0x00, A1: 0xfe, A2: 0x01},
			},
		},
	},
	kwCALL: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // CALL NN
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL NZ? NN
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL NZ? N
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL Z? NN
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL Z? N
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL NC? NN
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL NC? N
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL C? NN
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL C? N
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPO: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL PO? NN
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL PO? N
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondPE: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL PE? NN
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL PE? N
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondP: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL P? NN
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL P? N
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwCondM: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // CALL M? NN
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // CALL M? N
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcHigh, A0: 0x01},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // CALL N
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
		},
	},
	kwRET: InstPat{
		nil: InstDat{ // RET
			{Kind: BcByte, A0: 0xc9},
		},
		kwCondNZ: InstPat{
			nil: InstDat{ // RET NZ?
				{Kind: BcByte, A0: 0xc0},
			},
		},
		kwCondZ: InstPat{
			nil: InstDat{ // RET Z?
				{Kind: BcByte, A0: 0xc8},
			},
		},
		kwCondNC: InstPat{
			nil: InstDat{ // RET NC?
				{Kind: BcByte, A0: 0xd0},
			},
		},
		kwCondC: InstPat{
			nil: InstDat{ // RET C?
				{Kind: BcByte, A0: 0xd8},
			},
		},
		kwCondPO: InstPat{
			nil: InstDat{ // RET PO?
				{Kind: BcByte, A0: 0xe0},
			},
		},
		kwCondPE: InstPat{
			nil: InstDat{ // RET PE?
				{Kind: BcByte, A0: 0xe8},
			},
		},
		kwCondP: InstPat{
			nil: InstDat{ // RET P?
				{Kind: BcByte, A0: 0xf0},
			},
		},
		kwCondM: InstPat{
			nil: InstDat{ // RET M?
				{Kind: BcByte, A0: 0xf8},
			},
		},
	},
	kwRETI: InstPat{
		nil: InstDat{ // RETI
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x4d},
		},
	},
	kwRETN: InstPat{
		nil: InstDat{ // RETN
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x45},
		},
	},
	kwRST: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // RST N
				{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x38, A3: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // RST NN
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	KwJump: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // #.jump NN
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{ // #.jump NN NZ?
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{ // #.jump NN Z?
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{ // #.jump NN NC?
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{ // #.jump NN C?
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{ // #.jump NN PO?
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{ // #.jump NN PE?
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{ // #.jump NN P?
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{ // #.jump NN M?
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // #.jump N
				{Kind: BcByte, A0: 0xc3},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{ // #.jump N NZ?
					{Kind: BcByte, A0: 0xc2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{ // #.jump N Z?
					{Kind: BcByte, A0: 0xca},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{ // #.jump N NC?
					{Kind: BcByte, A0: 0xd2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{ // #.jump N C?
					{Kind: BcByte, A0: 0xda},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{ // #.jump N PO?
					{Kind: BcByte, A0: 0xe2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{ // #.jump N PE?
					{Kind: BcByte, A0: 0xea},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{ // #.jump N P?
					{Kind: BcByte, A0: 0xf2},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{ // #.jump N M?
					{Kind: BcByte, A0: 0xfa},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	KwCall: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // #.call NN
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{ // #.call NN NZ?
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{ // #.call NN Z?
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{ // #.call NN NC?
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{ // #.call NN C?
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{ // #.call NN PO?
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{ // #.call NN PE?
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{ // #.call NN P?
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{ // #.call NN M?
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // #.call N
				{Kind: BcByte, A0: 0xcd},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcHigh, A0: 0x00},
			},
			kwCondNZ: InstPat{
				nil: InstDat{ // #.call N NZ?
					{Kind: BcByte, A0: 0xc4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondZ: InstPat{
				nil: InstDat{ // #.call N Z?
					{Kind: BcByte, A0: 0xcc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondNC: InstPat{
				nil: InstDat{ // #.call N NC?
					{Kind: BcByte, A0: 0xd4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondC: InstPat{
				nil: InstDat{ // #.call N C?
					{Kind: BcByte, A0: 0xdc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPO: InstPat{
				nil: InstDat{ // #.call N PO?
					{Kind: BcByte, A0: 0xe4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondPE: InstPat{
				nil: InstDat{ // #.call N PE?
					{Kind: BcByte, A0: 0xec},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondP: InstPat{
				nil: InstDat{ // #.call N P?
					{Kind: BcByte, A0: 0xf4},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
			kwCondM: InstPat{
				nil: InstDat{ // #.call N M?
					{Kind: BcByte, A0: 0xfc},
					{Kind: BcLow, A0: 0x00},
					{Kind: BcHigh, A0: 0x00},
				},
			},
		},
	},
	KwReturn: InstPat{
		nil: InstDat{ // #.return
			{Kind: BcByte, A0: 0xc9},
		},
		kwCondNZ: InstPat{
			nil: InstDat{ // #.return NZ?
				{Kind: BcByte, A0: 0xc0},
			},
		},
		kwCondZ: InstPat{
			nil: InstDat{ // #.return Z?
				{Kind: BcByte, A0: 0xc8},
			},
		},
		kwCondNC: InstPat{
			nil: InstDat{ // #.return NC?
				{Kind: BcByte, A0: 0xd0},
			},
		},
		kwCondC: InstPat{
			nil: InstDat{ // #.return C?
				{Kind: BcByte, A0: 0xd8},
			},
		},
		kwCondPO: InstPat{
			nil: InstDat{ // #.return PO?
				{Kind: BcByte, A0: 0xe0},
			},
		},
		kwCondPE: InstPat{
			nil: InstDat{ // #.return PE?
				{Kind: BcByte, A0: 0xe8},
			},
		},
		kwCondP: InstPat{
			nil: InstDat{ // #.return P?
				{Kind: BcByte, A0: 0xf0},
			},
		},
		kwCondM: InstPat{
			nil: InstDat{ // #.return M?
				{Kind: BcByte, A0: 0xf8},
			},
		},
	},
	kwIN: InstPat{
		kwRegA: InstPat{
			kwMemN: InstPat{
				nil: InstDat{ // IN A N$
					{Kind: BcByte, A0: 0xdb},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwMemC: InstPat{
				nil: InstDat{ // IN A C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x78},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // IN A NN$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN B C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x40},
				},
			},
		},
		kwRegC: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN C C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x48},
				},
			},
		},
		kwRegD: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN D C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x50},
				},
			},
		},
		kwRegE: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN E C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x58},
				},
			},
		},
		kwRegH: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN H C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x60},
				},
			},
		},
		kwRegL: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN L C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x68},
				},
			},
		},
	},
	kwINI: InstPat{
		nil: InstDat{ // INI
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa2},
		},
	},
	kwINIR: InstPat{
		nil: InstDat{ // INIR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb2},
		},
	},
	kwIND: InstPat{
		nil: InstDat{ // IND
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xaa},
		},
	},
	kwINDR: InstPat{
		nil: InstDat{ // INDR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xba},
		},
	},
	kwOUT: InstPat{
		kwMemN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // OUT N$ A
					{Kind: BcByte, A0: 0xd3},
					{Kind: BcLow, A0: 0x00},
				},
			},
		},
		kwMemC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // OUT C$ A
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x79},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // OUT C$ B
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x41},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // OUT C$ C
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x49},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // OUT C$ D
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x51},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // OUT C$ E
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x59},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // OUT C$ H
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // OUT C$ L
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x69},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // OUT NN$ A
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwOUTI: InstPat{
		nil: InstDat{ // OUTI
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xa3},
		},
	},
	kwOTIR: InstPat{
		nil: InstDat{ // OTIR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xb3},
		},
	},
	kwOUTD: InstPat{
		nil: InstDat{ // OUTD
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xab},
		},
	},
	kwOTDR: InstPat{
		nil: InstDat{ // OTDR
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0xbb},
		},
	},
	kwDAA: InstPat{
		nil: InstDat{ // DAA
			{Kind: BcByte, A0: 0x27},
		},
	},
	kwCPL: InstPat{
		nil: InstDat{ // CPL
			{Kind: BcByte, A0: 0x2f},
		},
	},
	kwNEG: InstPat{
		nil: InstDat{ // NEG
			{Kind: BcByte, A0: 0xed},
			{Kind: BcByte, A0: 0x44},
		},
	},
	kwCCF: InstPat{
		nil: InstDat{ // CCF
			{Kind: BcByte, A0: 0x3f},
		},
	},
	kwSCF: InstPat{
		nil: InstDat{ // SCF
			{Kind: BcByte, A0: 0x37},
		},
	},
	kwNOP: InstPat{
		nil: InstDat{ // NOP
			{Kind: BcByte, A0: 0x00},
		},
	},
	kwHALT: InstPat{
		nil: InstDat{ // HALT
			{Kind: BcByte, A0: 0x76},
		},
	},
	kwDI: InstPat{
		nil: InstDat{ // DI
			{Kind: BcByte, A0: 0xf3},
		},
	},
	kwEI: InstPat{
		nil: InstDat{ // EI
			{Kind: BcByte, A0: 0xfb},
		},
	},
	kwIM: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // IM N
				{Kind: BcByte, A0: 0xed},
				{Kind: BcMap, A0: 0x00, A1: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // IM NN
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
			kwRegBC: {
				{kwLD, &Operand{Kind: kwRegB}, &Operand{Kind: kwRegB}},
				{kwLD, &Operand{Kind: kwRegC}, &Operand{Kind: kwRegC}},
			},
			kwRegDE: {
				{kwLD, &Operand{Kind: kwRegB}, &Operand{Kind: kwRegD}},
				{kwLD, &Operand{Kind: kwRegC}, &Operand{Kind: kwRegE}},
			},
			kwRegHL: {
				{kwLD, &Operand{Kind: kwRegB}, &Operand{Kind: kwRegH}},
				{kwLD, &Operand{Kind: kwRegC}, &Operand{Kind: kwRegL}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegDE: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegBC: {
				{kwLD, &Operand{Kind: kwRegD}, &Operand{Kind: kwRegB}},
				{kwLD, &Operand{Kind: kwRegE}, &Operand{Kind: kwRegC}},
			},
			kwRegDE: {
				{kwLD, &Operand{Kind: kwRegD}, &Operand{Kind: kwRegD}},
				{kwLD, &Operand{Kind: kwRegE}, &Operand{Kind: kwRegE}},
			},
			kwRegHL: {
				{kwLD, &Operand{Kind: kwRegD}, &Operand{Kind: kwRegH}},
				{kwLD, &Operand{Kind: kwRegE}, &Operand{Kind: kwRegL}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
		},
		kwRegHL: {
			KwAny: {
				{kwLD, &Vec{Int(0), nil}, &Vec{Int(1), nil}},
			},
			kwRegBC: {
				{kwLD, &Operand{Kind: kwRegH}, &Operand{Kind: kwRegB}},
				{kwLD, &Operand{Kind: kwRegL}, &Operand{Kind: kwRegC}},
			},
			kwRegDE: {
				{kwLD, &Operand{Kind: kwRegH}, &Operand{Kind: kwRegD}},
				{kwLD, &Operand{Kind: kwRegL}, &Operand{Kind: kwRegE}},
			},
			kwRegHL: {
				{kwLD, &Operand{Kind: kwRegH}, &Operand{Kind: kwRegH}},
				{kwLD, &Operand{Kind: kwRegL}, &Operand{Kind: kwRegL}},
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
			kwRegBC: {
				{kwLD, &Operand{Kind: kwRegB}, &Operand{Kind: kwRegB}},
				{kwLD, &Operand{Kind: kwRegC}, &Operand{Kind: kwRegC}},
			},
			kwRegDE: {
				{kwLD, &Operand{Kind: kwRegD}, &Operand{Kind: kwRegB}},
				{kwLD, &Operand{Kind: kwRegE}, &Operand{Kind: kwRegC}},
			},
			kwRegHL: {
				{kwLD, &Operand{Kind: kwRegH}, &Operand{Kind: kwRegB}},
				{kwLD, &Operand{Kind: kwRegL}, &Operand{Kind: kwRegC}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegDE: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegBC: {
				{kwLD, &Operand{Kind: kwRegB}, &Operand{Kind: kwRegD}},
				{kwLD, &Operand{Kind: kwRegC}, &Operand{Kind: kwRegE}},
			},
			kwRegDE: {
				{kwLD, &Operand{Kind: kwRegD}, &Operand{Kind: kwRegD}},
				{kwLD, &Operand{Kind: kwRegE}, &Operand{Kind: kwRegE}},
			},
			kwRegHL: {
				{kwLD, &Operand{Kind: kwRegH}, &Operand{Kind: kwRegD}},
				{kwLD, &Operand{Kind: kwRegL}, &Operand{Kind: kwRegE}},
			},
			kwRegPQ: {
				{kwLDP, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
		},
		kwRegHL: {
			KwAny: {
				{kwLD, &Vec{Int(1), nil}, &Vec{Int(0), nil}},
			},
			kwRegBC: {
				{kwLD, &Operand{Kind: kwRegB}, &Operand{Kind: kwRegH}},
				{kwLD, &Operand{Kind: kwRegC}, &Operand{Kind: kwRegL}},
			},
			kwRegDE: {
				{kwLD, &Operand{Kind: kwRegD}, &Operand{Kind: kwRegH}},
				{kwLD, &Operand{Kind: kwRegE}, &Operand{Kind: kwRegL}},
			},
			kwRegHL: {
				{kwLD, &Operand{Kind: kwRegH}, &Operand{Kind: kwRegH}},
				{kwLD, &Operand{Kind: kwRegL}, &Operand{Kind: kwRegL}},
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
	Intern("-return"): {
		kwRegPC: {
			nil: {
				{KwReturn},
			},
		},
	},
	Intern("-return-if"): {
		kwRegPC: {
			kwCondNZ: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondZ: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondNC: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondC: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondPO: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondPE: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondP: {
				{KwReturn, &Vec{Int(1), nil}},
			},
			kwCondM: {
				{KwReturn, &Vec{Int(1), nil}},
			},
		},
	},
	Intern("-return-unless"): {
		kwRegPC: {
			kwCondNZ: {
				{KwReturn, &Operand{Kind: kwCondZ}},
			},
			kwCondZ: {
				{KwReturn, &Operand{Kind: kwCondNZ}},
			},
			kwCondNC: {
				{KwReturn, &Operand{Kind: kwCondC}},
			},
			kwCondC: {
				{KwReturn, &Operand{Kind: kwCondNC}},
			},
			kwCondPO: {
				{KwReturn, &Operand{Kind: kwCondPE}},
			},
			kwCondPE: {
				{KwReturn, &Operand{Kind: kwCondPO}},
			},
			kwCondM: {
				{KwReturn, &Operand{Kind: kwCondP}},
			},
			kwCondP: {
				{KwReturn, &Operand{Kind: kwCondM}},
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
				nil: InstDat{ // LD A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
		},
		kwRegB: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD B IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD B IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD B IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD B IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x45},
				},
			},
		},
		kwRegC: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD C IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD C IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD C IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD C IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
		},
		kwRegD: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD D IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD D IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x55},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD D IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD D IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x55},
				},
			},
		},
		kwRegE: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD E IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD E IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD E IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD E IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
		},
		kwRegIXH: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD IXH IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IXH A
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IXH B
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IXH C
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IXH D
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IXH E
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD IXH IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IXH N
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IXH NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIXL: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD IXL IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD IXL IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IXL A
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IXL B
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IXL C
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IXL D
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IXL E
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IXL N
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IXL NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYH: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{ // LD IYH IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IYH A
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IYH B
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IYH C
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IYH D
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IYH E
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD IYH IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IYH N
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IYH NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYL: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{ // LD IYL IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD IYL IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IYL A
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IYL B
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IYL C
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IYL D
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IYL E
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IYL N
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IYL NN
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
				nil: InstDat{ // ADD A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // ADD A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x85},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // ADD A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // ADD A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x85},
				},
			},
		},
	},
	kwADC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // ADC A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // ADC A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // ADC A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // ADC A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
		},
	},
	kwSUB: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // SUB IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // SUB IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x95},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // SUB IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // SUB IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x95},
			},
		},
	},
	kwSBC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // SBC A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // SBC A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // SBC A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // SBC A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
		},
	},
	kwAND: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // AND IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // AND IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // AND IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // AND IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
	},
	kwOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // OR IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // OR IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // OR IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // OR IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
	},
	kwXOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // XOR IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // XOR IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xad},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // XOR IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // XOR IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xad},
			},
		},
	},
	kwCP: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // CP IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // CP IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // CP IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // CP IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
	},
	kwINC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // INC IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // INC IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // INC IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // INC IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
	},
	kwDEC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // DEC IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // DEC IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // DEC IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // DEC IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
	},
	kwRLC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RLC A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x07},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RLC A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x07},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RLC A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RLC A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RLC B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RLC B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RLC B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RLC B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RLC C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x01},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RLC C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x01},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RLC C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RLC C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RLC D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x02},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RLC D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x02},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RLC D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RLC D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RLC E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x03},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RLC E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x03},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RLC E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RLC E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RLC H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x04},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RLC H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x04},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RLC H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RLC H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RLC L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x05},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RLC L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x05},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RLC L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RLC L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwRL: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RL A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x17},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RL A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x17},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RL A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RL A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RL B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x10},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RL B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x10},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RL B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RL B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RL C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x11},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RL C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x11},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RL C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RL C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RL D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x12},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RL D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x12},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RL D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RL D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RL E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x13},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RL E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x13},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RL E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RL E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RL H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x14},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RL H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x14},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RL H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RL H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RL L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x15},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RL L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x15},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RL L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RL L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwRRC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RRC A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RRC A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0f},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RRC A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RRC A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RRC B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x08},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RRC B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x08},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RRC B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RRC B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RRC C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RRC C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x09},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RRC C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RRC C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RRC D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RRC D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0a},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RRC D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RRC D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RRC E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RRC E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0b},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RRC E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RRC E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RRC H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RRC H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0c},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RRC H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RRC H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RRC L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RRC L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x0d},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RRC L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RRC L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwRR: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RR A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RR A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1f},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RR A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RR A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RR B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x18},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RR B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x18},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RR B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RR B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RR C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RR C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x19},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RR C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RR C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RR D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RR D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1a},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RR D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RR D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RR E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RR E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1b},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RR E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RR E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RR H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RR H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1c},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RR H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RR H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // RR L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RR L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x1d},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RR L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RR L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwSLA: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SLA A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x27},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLA A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x27},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLA A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLA A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SLA B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x20},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLA B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x20},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLA B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLA B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SLA C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x21},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLA C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x21},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLA C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLA C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SLA D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x22},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLA D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x22},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLA D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLA D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SLA E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x23},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLA E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x23},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLA E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLA E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SLA H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x24},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLA H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x24},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLA H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLA H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SLA L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x25},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLA L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x25},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLA L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLA L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwSLL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SLL A
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x37},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SLL A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x37},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLL A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x37},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLL A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLL A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SLL B
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x30},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SLL B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x30},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLL B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x30},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLL B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLL B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SLL C
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x31},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SLL C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x31},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLL C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x31},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLL C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLL C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SLL D
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x32},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SLL D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x32},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLL D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x32},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLL D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLL D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SLL E
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x33},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SLL E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x33},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLL E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x33},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLL E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLL E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SLL H
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x34},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SLL H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x34},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLL H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x34},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLL H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLL H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SLL L
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x35},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SLL L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x35},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SLL L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x35},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SLL L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SLL L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SLL HL$
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcByte, A0: 0x36},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SLL IX$
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x36},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SLL IY$
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xcb},
				{Kind: BcLow, A0: 0x00},
				{Kind: BcByte, A0: 0x36},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SLL WX$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SLL WY$
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
				{Kind: BcTemp, A0: 0x00},
			},
		},
	},
	kwSRA: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRA A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRA A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2f},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRA A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRA A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRA B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x28},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRA B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x28},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRA B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRA B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRA C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRA C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x29},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRA C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRA C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRA D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRA D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2a},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRA D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRA D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRA E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRA E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2b},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRA E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRA E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRA H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRA H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2c},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRA H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRA H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRA L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRA L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x2d},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRA L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRA L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwSRL: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRL A IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3f},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRL A IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3f},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRL A WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRL A WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRL B IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x38},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRL B IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x38},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRL B WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRL B WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRL C IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x39},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRL C IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x39},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRL C WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRL C WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRL D IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3a},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRL D IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3a},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRL D WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRL D WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRL E IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3b},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRL E IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3b},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRL E WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRL E WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRL H IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3c},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRL H IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3c},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRL H WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRL H WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SRL L IX$
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3d},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SRL L IY$
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0xcb},
					{Kind: BcLow, A0: 0x01},
					{Kind: BcByte, A0: 0x3d},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SRL L WX$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SRL L WY$
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
	},
	kwBIT: InstPat{
		kwRegA: InstPat{
			kwImmN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT A N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT A N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT A N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT A N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT A NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT A NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT A NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT A NN WY$
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
					nil: InstDat{ // BIT B N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT B N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT B N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT B N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT B NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT B NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT B NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT B NN WY$
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
					nil: InstDat{ // BIT C N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT C N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT C N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT C N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT C NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT C NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT C NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT C NN WY$
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
					nil: InstDat{ // BIT D N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT D N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT D N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT D N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT D NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT D NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT D NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT D NN WY$
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
					nil: InstDat{ // BIT E N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT E N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT E N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT E N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT E NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT E NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT E NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT E NN WY$
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
					nil: InstDat{ // BIT H N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT H N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT H N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT H N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT H NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT H NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT H NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT H NN WY$
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
					nil: InstDat{ // BIT L N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT L N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT L N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT L N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // BIT L NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // BIT L NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // BIT L NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // BIT L NN WY$
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
					nil: InstDat{ // BIT N IX$ A
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT N IX$ B
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT N IX$ C
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT N IX$ D
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT N IX$ E
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT N IX$ H
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT N IX$ L
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // BIT N IY$ A
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x47, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT N IY$ B
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x40, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT N IY$ C
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x41, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT N IY$ D
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x42, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT N IY$ E
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x43, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT N IY$ H
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x44, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT N IY$ L
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x45, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemWX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // BIT N WX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT N WX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT N WX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT N WX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT N WX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT N WX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT N WX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // BIT N WY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT N WY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT N WY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT N WY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT N WY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT N WY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT N WY$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
		},
		kwImmNN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // BIT NN IX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT NN IX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT NN IX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT NN IX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT NN IX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT NN IX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT NN IX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // BIT NN WX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT NN WX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT NN WX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT NN WX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT NN WX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT NN WX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT NN WX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // BIT NN IY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT NN IY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT NN IY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT NN IY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT NN IY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT NN IY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT NN IY$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // BIT NN WY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // BIT NN WY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // BIT NN WY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // BIT NN WY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // BIT NN WY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // BIT NN WY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // BIT NN WY$ L
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
					nil: InstDat{ // SET A N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET A N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET A N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET A N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // SET A NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET A NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET A NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET A NN WY$
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
					nil: InstDat{ // SET B N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET B N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET B N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET B N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // SET B NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET B NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET B NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET B NN WY$
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
					nil: InstDat{ // SET C N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET C N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET C N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET C N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // SET C NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET C NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET C NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET C NN WY$
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
					nil: InstDat{ // SET D N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET D N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET D N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET D N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // SET D NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET D NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET D NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET D NN WY$
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
					nil: InstDat{ // SET E N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET E N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET E N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET E N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // SET E NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET E NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET E NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET E NN WY$
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
					nil: InstDat{ // SET H N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET H N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET H N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET H N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // SET H NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET H NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET H NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET H NN WY$
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
					nil: InstDat{ // SET L N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET L N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET L N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET L N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // SET L NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // SET L NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // SET L NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // SET L NN WY$
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
					nil: InstDat{ // SET N IX$ A
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET N IX$ B
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET N IX$ C
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET N IX$ D
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET N IX$ E
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET N IX$ H
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET N IX$ L
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // SET N IY$ A
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc7, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET N IY$ B
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc0, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET N IY$ C
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc1, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET N IY$ D
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc2, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET N IY$ E
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc3, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET N IY$ H
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc4, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET N IY$ L
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0xc5, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemWX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // SET N WX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET N WX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET N WX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET N WX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET N WX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET N WX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET N WX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // SET N WY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET N WY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET N WY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET N WY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET N WY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET N WY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET N WY$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
		},
		kwImmNN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // SET NN IX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET NN IX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET NN IX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET NN IX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET NN IX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET NN IX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET NN IX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // SET NN WX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET NN WX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET NN WX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET NN WX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET NN WX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET NN WX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET NN WX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // SET NN IY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET NN IY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET NN IY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET NN IY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET NN IY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET NN IY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET NN IY$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // SET NN WY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // SET NN WY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // SET NN WY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // SET NN WY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // SET NN WY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // SET NN WY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // SET NN WY$ L
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
					nil: InstDat{ // RES A N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES A N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES A N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES A N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // RES A NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES A NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES A NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES A NN WY$
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
					nil: InstDat{ // RES B N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES B N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES B N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES B N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // RES B NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES B NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES B NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES B NN WY$
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
					nil: InstDat{ // RES C N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES C N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES C N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES C N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // RES C NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES C NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES C NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES C NN WY$
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
					nil: InstDat{ // RES D N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES D N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES D N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES D N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // RES D NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES D NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES D NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES D NN WY$
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
					nil: InstDat{ // RES E N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES E N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES E N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES E N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // RES E NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES E NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES E NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES E NN WY$
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
					nil: InstDat{ // RES H N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES H N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES H N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES H N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // RES H NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES H NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES H NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES H NN WY$
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
					nil: InstDat{ // RES L N IX$
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES L N IY$
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x02},
						{Kind: BcImp, A0: 0x01, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES L N WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES L N WY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwImmNN: InstPat{
				kwMemIX: InstPat{
					nil: InstDat{ // RES L NN IX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWX: InstPat{
					nil: InstDat{ // RES L NN WX$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemIY: InstPat{
					nil: InstDat{ // RES L NN IY$
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwMemWY: InstPat{
					nil: InstDat{ // RES L NN WY$
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
					nil: InstDat{ // RES N IX$ A
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES N IX$ B
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES N IX$ C
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES N IX$ D
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES N IX$ E
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES N IX$ H
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES N IX$ L
						{Kind: BcByte, A0: 0xdd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // RES N IY$ A
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x87, A2: 0x07, A3: 0x03},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES N IY$ B
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x80, A2: 0x07, A3: 0x03},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES N IY$ C
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x81, A2: 0x07, A3: 0x03},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES N IY$ D
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x82, A2: 0x07, A3: 0x03},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES N IY$ E
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x83, A2: 0x07, A3: 0x03},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES N IY$ H
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x84, A2: 0x07, A3: 0x03},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES N IY$ L
						{Kind: BcByte, A0: 0xfd},
						{Kind: BcByte, A0: 0xcb},
						{Kind: BcLow, A0: 0x01},
						{Kind: BcImp, A0: 0x00, A1: 0x85, A2: 0x07, A3: 0x03},
					},
				},
			},
			kwMemWX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // RES N WX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES N WX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES N WX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES N WX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES N WX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES N WX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES N WX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // RES N WY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES N WY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES N WY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES N WY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES N WY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES N WY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES N WY$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
		},
		kwImmNN: InstPat{
			kwMemIX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // RES NN IX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES NN IX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES NN IX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES NN IX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES NN IX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES NN IX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES NN IX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWX: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // RES NN WX$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES NN WX$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES NN WX$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES NN WX$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES NN WX$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES NN WX$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES NN WX$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemIY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // RES NN IY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES NN IY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES NN IY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES NN IY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES NN IY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES NN IY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES NN IY$ L
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
			},
			kwMemWY: InstPat{
				kwRegA: InstPat{
					nil: InstDat{ // RES NN WY$ A
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegB: InstPat{
					nil: InstDat{ // RES NN WY$ B
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegC: InstPat{
					nil: InstDat{ // RES NN WY$ C
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegD: InstPat{
					nil: InstDat{ // RES NN WY$ D
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegE: InstPat{
					nil: InstDat{ // RES NN WY$ E
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegH: InstPat{
					nil: InstDat{ // RES NN WY$ H
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
						{Kind: BcTemp, A0: 0x00},
					},
				},
				kwRegL: InstPat{
					nil: InstDat{ // RES NN WY$ L
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
			nil: InstDat{ // IN C$
				{Kind: BcByte, A0: 0xed},
				{Kind: BcByte, A0: 0x70},
			},
		},
		kwRegF: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN F C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x70},
				},
			},
		},
	},
	kwOUT: InstPat{
		kwMemC: InstPat{
			kwImmN: InstPat{
				nil: InstDat{ // OUT C$ N
					{Kind: BcByte, A0: 0xed},
					{Kind: BcImp, A0: 0x01, A1: 0x71, A2: 0x00, A3: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // OUT C$ NN
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
				nil: InstDat{ // LD A IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD A IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegI: InstPat{
				nil: InstDat{ // LD A I
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegR: InstPat{
				nil: InstDat{ // LD A R
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD A WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD A WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // LD B IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD B IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD B WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD B WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // LD C IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD C IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD C WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD C WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // LD D IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD D IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD D WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD D WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // LD E IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD E IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD E WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD E WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // LD H IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD H IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD H WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD H WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // LD L IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // LD L IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // LD L WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // LD L WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemIX: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD IX$ A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IX$ B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IX$ C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IX$ D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IX$ E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD IX$ H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD IX$ L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IX$ N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IX$ NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemIY: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD IY$ A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IY$ B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IY$ C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IY$ D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IY$ E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD IY$ H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD IY$ L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IY$ N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IY$ NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegI: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD I A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegR: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD R A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIX: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD IX NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD IX NN$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IX N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD IX N$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIY: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // LD IY NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemNN: InstPat{
				nil: InstDat{ // LD IY NN$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IY N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD IY N$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegBC: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{ // LD BC NN$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD BC N$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegDE: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{ // LD DE NN$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD DE N$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{ // LD HL NN$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD HL N$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegSP: InstPat{
			kwMemNN: InstPat{
				nil: InstDat{ // LD SP NN$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // LD SP IX
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // LD SP IY
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemN: InstPat{
				nil: InstDat{ // LD SP N$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemNN: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // LD NN$ BC
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // LD NN$ DE
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // LD NN$ HL
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // LD NN$ SP
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // LD NN$ IX
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // LD NN$ IY
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemWX: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD WX$ A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD WX$ B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD WX$ C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD WX$ D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD WX$ E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD WX$ H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD WX$ L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD WX$ N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD WX$ NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemWY: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // LD WY$ A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD WY$ B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD WY$ C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD WY$ D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD WY$ E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // LD WY$ H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // LD WY$ L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD WY$ N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD WY$ NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemN: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // LD N$ BC
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // LD N$ DE
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // LD N$ HL
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // LD N$ SP
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // LD N$ IX
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // LD N$ IY
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwPUSH: InstPat{
		kwRegIX: InstPat{
			nil: InstDat{ // PUSH IX
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // PUSH IY
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwPOP: InstPat{
		kwRegIX: InstPat{
			nil: InstDat{ // POP IX
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // POP IY
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwEX: InstPat{
		kwRegAF: InstPat{
			kwAltAF: InstPat{
				nil: InstDat{ // EX AF AF-
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwMemSP: InstPat{
			kwRegIX: InstPat{
				nil: InstDat{ // EX SP$ IX
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // EX SP$ IY
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwEXX: InstPat{
		nil: InstDat{ // EXX
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDI: InstPat{
		nil: InstDat{ // LDI
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDIR: InstPat{
		nil: InstDat{ // LDIR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDD: InstPat{
		nil: InstDat{ // LDD
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwLDDR: InstPat{
		nil: InstDat{ // LDDR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPI: InstPat{
		nil: InstDat{ // CPI
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPIR: InstPat{
		nil: InstDat{ // CPIR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPD: InstPat{
		nil: InstDat{ // CPD
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwCPDR: InstPat{
		nil: InstDat{ // CPDR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwADD: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // ADD A IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // ADD A IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // ADD A WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // ADD A WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIX: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // ADD IX BC
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // ADD IX DE
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIX: InstPat{
				nil: InstDat{ // ADD IX IX
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // ADD IX SP
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegIY: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // ADD IY BC
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // ADD IY DE
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegIY: InstPat{
				nil: InstDat{ // ADD IY IY
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // ADD IY SP
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwADC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // ADC A IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // ADC A IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // ADC A WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // ADC A WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // ADC HL BC
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // ADC HL DE
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // ADC HL HL
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // ADC HL SP
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwSUB: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // SUB IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SUB IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SUB WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SUB WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSBC: InstPat{
		kwRegA: InstPat{
			kwMemIX: InstPat{
				nil: InstDat{ // SBC A IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SBC A IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SBC A WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SBC A WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // SBC HL BC
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegDE: InstPat{
				nil: InstDat{ // SBC HL DE
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegHL: InstPat{
				nil: InstDat{ // SBC HL HL
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // SBC HL SP
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwAND: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // AND IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // AND IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // AND WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // AND WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwOR: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // OR IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // OR IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // OR WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // OR WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwXOR: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // XOR IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // XOR IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // XOR WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // XOR WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwCP: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // CP IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // CP IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // CP WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // CP WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwINC: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // INC IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // INC IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{ // INC IX
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // INC IY
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // INC WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // INC WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwDEC: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // DEC IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // DEC IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIX: InstPat{
			nil: InstDat{ // DEC IX
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegIY: InstPat{
			nil: InstDat{ // DEC IY
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // DEC WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // DEC WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRLC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RLC A
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RLC B
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RLC C
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RLC D
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RLC E
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RLC H
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RLC L
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RLC HL$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RLC IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RLC IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RLC WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RLC WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RL A
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RL B
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RL C
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RL D
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RL E
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RL H
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RL L
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RL HL$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RL IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RL IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RL WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RL WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRRC: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RRC A
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RRC B
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RRC C
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RRC D
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RRC E
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RRC H
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RRC L
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RRC HL$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RRC IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RRC IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RRC WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RRC WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRR: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // RR A
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // RR B
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // RR C
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // RR D
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // RR E
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // RR H
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // RR L
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // RR HL$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // RR IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // RR IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // RR WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // RR WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSLA: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SLA A
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SLA B
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SLA C
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SLA D
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SLA E
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SLA H
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SLA L
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SLA HL$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SLA IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SLA IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SLA WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SLA WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSRA: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SRA A
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SRA B
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SRA C
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SRA D
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SRA E
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SRA H
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SRA L
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SRA HL$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SRA IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SRA IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SRA WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SRA WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwSRL: InstPat{
		kwRegA: InstPat{
			nil: InstDat{ // SRL A
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegB: InstPat{
			nil: InstDat{ // SRL B
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegC: InstPat{
			nil: InstDat{ // SRL C
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegD: InstPat{
			nil: InstDat{ // SRL D
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegE: InstPat{
			nil: InstDat{ // SRL E
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegH: InstPat{
			nil: InstDat{ // SRL H
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwRegL: InstPat{
			nil: InstDat{ // SRL L
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemHL: InstPat{
			nil: InstDat{ // SRL HL$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIX: InstPat{
			nil: InstDat{ // SRL IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // SRL IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // SRL WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // SRL WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRLD: InstPat{
		nil: InstDat{ // RLD
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwRRD: InstPat{
		nil: InstDat{ // RRD
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwBIT: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // BIT N A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // BIT N B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // BIT N C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // BIT N D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // BIT N E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // BIT N H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // BIT N L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // BIT N HL$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // BIT N IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // BIT N IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // BIT N WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // BIT N WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // BIT NN A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // BIT NN B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // BIT NN C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // BIT NN D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // BIT NN E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // BIT NN H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // BIT NN L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // BIT NN HL$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // BIT NN IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // BIT NN WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // BIT NN IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // BIT NN WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwSET: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // SET N A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // SET N B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // SET N C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // SET N D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // SET N E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // SET N H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // SET N L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // SET N HL$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SET N IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SET N IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SET N WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SET N WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // SET NN A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // SET NN B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // SET NN C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // SET NN D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // SET NN E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // SET NN H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // SET NN L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // SET NN HL$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // SET NN IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // SET NN WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // SET NN IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // SET NN WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwRES: InstPat{
		kwImmN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // RES N A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // RES N B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // RES N C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // RES N D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // RES N E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // RES N H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // RES N L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // RES N HL$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // RES N IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RES N IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RES N WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RES N WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmNN: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // RES NN A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // RES NN B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // RES NN C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // RES NN D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // RES NN E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // RES NN H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // RES NN L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemHL: InstPat{
				nil: InstDat{ // RES NN HL$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIX: InstPat{
				nil: InstDat{ // RES NN IX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWX: InstPat{
				nil: InstDat{ // RES NN WX$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemIY: InstPat{
				nil: InstDat{ // RES NN IY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwMemWY: InstPat{
				nil: InstDat{ // RES NN WY$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwJP: InstPat{
		kwMemIX: InstPat{
			nil: InstDat{ // JP IX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemIY: InstPat{
			nil: InstDat{ // JP IY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWX: InstPat{
			nil: InstDat{ // JP WX$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwMemWY: InstPat{
			nil: InstDat{ // JP WY$
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwJR: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // JR NN
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwCondC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR C? NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR C? N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwCondNC: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR NC? NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR NC? N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwCondZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR Z? NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR Z? N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwCondNZ: InstPat{
			kwImmNN: InstPat{
				nil: InstDat{ // JR NZ? NN
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // JR NZ? N
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // JR N
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwDJNZ: InstPat{
		kwImmNN: InstPat{
			nil: InstDat{ // DJNZ NN
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwImmN: InstPat{
			nil: InstDat{ // DJNZ N
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
	},
	kwRETI: InstPat{
		nil: InstDat{ // RETI
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwRETN: InstPat{
		nil: InstDat{ // RETN
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwIN: InstPat{
		kwRegA: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN A C$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegB: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN B C$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegC: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN C C$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegD: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN D C$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegE: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN E C$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegH: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN H C$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
		kwRegL: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN L C$
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwINI: InstPat{
		nil: InstDat{ // INI
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwINIR: InstPat{
		nil: InstDat{ // INIR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwIND: InstPat{
		nil: InstDat{ // IND
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwINDR: InstPat{
		nil: InstDat{ // INDR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOUT: InstPat{
		kwMemC: InstPat{
			kwRegA: InstPat{
				nil: InstDat{ // OUT C$ A
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // OUT C$ B
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // OUT C$ C
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // OUT C$ D
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // OUT C$ E
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegH: InstPat{
				nil: InstDat{ // OUT C$ H
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
			kwRegL: InstPat{
				nil: InstDat{ // OUT C$ L
					{Kind: BcUnsupported, A0: 0x00},
				},
			},
		},
	},
	kwOUTI: InstPat{
		nil: InstDat{ // OUTI
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOTIR: InstPat{
		nil: InstDat{ // OTIR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOUTD: InstPat{
		nil: InstDat{ // OUTD
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwOTDR: InstPat{
		nil: InstDat{ // OTDR
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwNEG: InstPat{
		nil: InstDat{ // NEG
			{Kind: BcUnsupported, A0: 0x00},
		},
	},
	kwIM: InstPat{
		kwImmN: InstPat{
			nil: InstDat{ // IM N
				{Kind: BcUnsupported, A0: 0x00},
			},
		},
		kwImmNN: InstPat{
			nil: InstDat{ // IM NN
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
				nil: InstDat{ // LD A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x7d},
				},
			},
		},
		kwRegB: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD B IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD B IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x45},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD B IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x44},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD B IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x45},
				},
			},
		},
		kwRegC: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD C IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD C IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD C IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD C IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x4d},
				},
			},
		},
		kwRegD: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD D IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD D IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x55},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD D IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x54},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD D IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x55},
				},
			},
		},
		kwRegE: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD E IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD E IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // LD E IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD E IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x5d},
				},
			},
		},
		kwRegIXH: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD IXH IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IXH A
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IXH B
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IXH C
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IXH D
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IXH E
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD IXH IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IXH N
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IXH NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIXL: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // LD IXL IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // LD IXL IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IXL A
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IXL B
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IXL C
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IXL D
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IXL E
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IXL N
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IXL NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYH: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{ // LD IYH IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x64},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IYH A
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x67},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IYH B
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x60},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IYH C
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x61},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IYH D
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x62},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IYH E
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x63},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD IYH IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x65},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IYH N
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x26},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IYH NN
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
					{Kind: BcTemp, A0: 0x00},
				},
			},
		},
		kwRegIYL: InstPat{
			kwRegIYH: InstPat{
				nil: InstDat{ // LD IYL IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // LD IYL IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6d},
				},
			},
			kwRegA: InstPat{
				nil: InstDat{ // LD IYL A
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6f},
				},
			},
			kwRegB: InstPat{
				nil: InstDat{ // LD IYL B
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x68},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // LD IYL C
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x69},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // LD IYL D
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6a},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // LD IYL E
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x6b},
				},
			},
			kwImmN: InstPat{
				nil: InstDat{ // LD IYL N
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x2e},
					{Kind: BcLow, A0: 0x01},
				},
			},
			kwImmNN: InstPat{
				nil: InstDat{ // LD IYL NN
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
				nil: InstDat{ // ADD A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // ADD A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x85},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // ADD A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x84},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // ADD A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x85},
				},
			},
		},
	},
	kwADC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // ADC A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // ADC A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // ADC A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // ADC A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x8d},
				},
			},
		},
	},
	kwSUB: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // SUB IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // SUB IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x95},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // SUB IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x94},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // SUB IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x95},
			},
		},
	},
	kwSBC: InstPat{
		kwRegA: InstPat{
			kwRegIXH: InstPat{
				nil: InstDat{ // SBC A IXH
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIXL: InstPat{
				nil: InstDat{ // SBC A IXL
					{Kind: BcByte, A0: 0xdd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
			kwRegIYH: InstPat{
				nil: InstDat{ // SBC A IYH
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9c},
				},
			},
			kwRegIYL: InstPat{
				nil: InstDat{ // SBC A IYL
					{Kind: BcByte, A0: 0xfd},
					{Kind: BcByte, A0: 0x9d},
				},
			},
		},
	},
	kwAND: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // AND IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // AND IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // AND IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // AND IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xa5},
			},
		},
	},
	kwOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // OR IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // OR IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // OR IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb4},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // OR IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xb5},
			},
		},
	},
	kwXOR: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // XOR IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // XOR IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xad},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // XOR IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xac},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // XOR IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xad},
			},
		},
	},
	kwCP: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // CP IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // CP IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // CP IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbc},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // CP IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0xbd},
			},
		},
	},
	kwINC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // INC IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // INC IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // INC IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x24},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // INC IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2c},
			},
		},
	},
	kwDEC: InstPat{
		kwRegIXH: InstPat{
			nil: InstDat{ // DEC IXH
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIXL: InstPat{
			nil: InstDat{ // DEC IXL
				{Kind: BcByte, A0: 0xdd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
		kwRegIYH: InstPat{
			nil: InstDat{ // DEC IYH
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x25},
			},
		},
		kwRegIYL: InstPat{
			nil: InstDat{ // DEC IYL
				{Kind: BcByte, A0: 0xfd},
				{Kind: BcByte, A0: 0x2d},
			},
		},
	},
	kwIN: InstPat{
		kwRegF: InstPat{
			kwMemC: InstPat{
				nil: InstDat{ // IN F C$
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0x70},
				},
			},
		},
	},
	kwMULUB: InstPat{
		kwRegA: InstPat{
			kwRegB: InstPat{
				nil: InstDat{ // MULUB A B
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xc1},
				},
			},
			kwRegC: InstPat{
				nil: InstDat{ // MULUB A C
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xc9},
				},
			},
			kwRegD: InstPat{
				nil: InstDat{ // MULUB A D
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xd1},
				},
			},
			kwRegE: InstPat{
				nil: InstDat{ // MULUB A E
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xd9},
				},
			},
		},
	},
	kwMULUW: InstPat{
		kwRegHL: InstPat{
			kwRegBC: InstPat{
				nil: InstDat{ // MULUW HL BC
					{Kind: BcByte, A0: 0xed},
					{Kind: BcByte, A0: 0xc3},
				},
			},
			kwRegSP: InstPat{
				nil: InstDat{ // MULUW HL SP
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
