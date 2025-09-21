L1:
    LD A, A
    LD A, B
    LD A, C
    LD A, D
    LD A, E
    LD A, H
    LD A, L
    LD B, A
    LD B, B
    LD B, C
    LD B, D
    LD B, E
    LD B, H
    LD B, L
    LD C, A
    LD C, B
    LD C, C
    LD C, D
    LD C, E
    LD C, H
    LD C, L
    LD D, A
    LD D, B
    LD D, C
    LD D, D
    LD D, E
    LD D, H
    LD D, L
    LD E, A
    LD E, B
    LD E, C
    LD E, D
    LD E, E
    LD E, H
    LD E, L
    LD H, A
    LD H, B
    LD H, C
    LD H, D
    LD H, E
    LD H, H
    LD H, L
    LD L, A
    LD L, B
    LD L, C
    LD L, D
    LD L, E
    LD L, H
    LD L, L
    LD A, 0+ 5
    LD B, 0+ 5
    LD C, 0+ 5
    LD D, 0+ 5
    LD E, 0+ 5
    LD H, 0+ 5
    LD L, 0+ 5
    LD A, (HL)
    LD B, (HL)
    LD C, (HL)
    LD D, (HL)
    LD E, (HL)
    LD H, (HL)
    LD L, (HL)
    LD A, (IX+5)
    LD B, (IX+5)
    LD C, (IX+5)
    LD D, (IX+5)
    LD E, (IX+5)
    LD H, (IX+5)
    LD L, (IX+5)
    LD A, (IY+5)
    LD B, (IY+5)
    LD C, (IY+5)
    LD D, (IY+5)
    LD E, (IY+5)
    LD H, (IY+5)
    LD L, (IY+5)
    LD (HL), A
    LD (HL), B
    LD (HL), C
    LD (HL), D
    LD (HL), E
    LD (HL), H
    LD (HL), L
    LD (IX+5), A
    LD (IX+5), B
    LD (IX+5), C
    LD (IX+5), D
    LD (IX+5), E
    LD (IX+5), H
    LD (IX+5), L
    LD (IY+5), A
    LD (IY+5), B
    LD (IY+5), C
    LD (IY+5), D
    LD (IY+5), E
    LD (IY+5), H
    LD (IY+5), L
    LD (HL), 0+ 5
    LD (IX+5), 0+ 5
    LD (IY+5), 0+ 5
    LD A, (BC)
    LD A, (DE)
    LD A, (1234)
    LD (BC), A
    LD (DE), A
    LD (1234), A
    LD A, I
    LD A, R
    LD I, A
    LD R, A
    LD BC, 0+ 1234
    LD DE, 0+ 1234
    LD HL, 0+ 1234
    LD SP, 0+ 1234
    LD IX, 0+ 1234
    LD IY, 0+ 1234
    LD HL, (1234)
    LD BC, (1234)
    LD DE, (1234)
    LD SP, (1234)
    LD IX, (1234)
    LD IY, (1234)
    LD (1234), HL
    LD (1234), BC
    LD (1234), DE
    LD (1234), SP
    LD (1234), IX
    LD (1234), IY
    LD SP, HL
    LD SP, IX
    LD SP, IY
    LD A, (5)
    LD (5), A
    LD BC, 0+ 5
    LD DE, 0+ 5
    LD HL, 0+ 5
    LD SP, 0+ 5
    LD IX, 0+ 5
    LD IY, 0+ 5
    LD HL, (5)
    LD BC, (5)
    LD DE, (5)
    LD SP, (5)
    LD IX, (5)
    LD IY, (5)
    LD (5), HL
    LD (5), BC
    LD (5), DE
    LD (5), SP
    LD (5), IX
    LD (5), IY
L2:
    PUSH BC
    PUSH DE
    PUSH HL
    PUSH AF
    PUSH IX
    PUSH IY
L3:
    POP BC
    POP DE
    POP HL
    POP AF
    POP IX
    POP IY
L4:
    EX DE, HL
    EX AF, AF'
    EX (SP), HL
    EX (SP), IX
    EX (SP), IY
L5:
    EXX
L6:
    LDI
L7:
    LDIR
L8:
    LDD
L9:
    LDDR
L10:
    CPI
L11:
    CPIR
L12:
    CPD
L13:
    CPDR
L14:
    ADD A, A
    ADD A, B
    ADD A, C
    ADD A, D
    ADD A, E
    ADD A, H
    ADD A, L
    ADD A, 0+ 5
    ADD A, (HL)
    ADD A, (IX+5)
    ADD A, (IY+5)
    ADD HL, BC
    ADD HL, DE
    ADD HL, HL
    ADD HL, SP
    ADD IX, BC
    ADD IX, DE
    ADD IX, IX
    ADD IX, SP
    ADD IY, BC
    ADD IY, DE
    ADD IY, IY
    ADD IY, SP
L15:
    ADC A, A
    ADC A, B
    ADC A, C
    ADC A, D
    ADC A, E
    ADC A, H
    ADC A, L
    ADC A, 0+ 5
    ADC A, (HL)
    ADC A, (IX+5)
    ADC A, (IY+5)
    ADC HL, BC
    ADC HL, DE
    ADC HL, HL
    ADC HL, SP
L16:
    SUB A
    SUB B
    SUB C
    SUB D
    SUB E
    SUB H
    SUB L
    SUB 0+ 5
    SUB (HL)
    SUB (IX+5)
    SUB (IY+5)
L17:
    SBC A, A
    SBC A, B
    SBC A, C
    SBC A, D
    SBC A, E
    SBC A, H
    SBC A, L
    SBC A, 0+ 5
    SBC A, (HL)
    SBC A, (IX+5)
    SBC A, (IY+5)
    SBC HL, BC
    SBC HL, DE
    SBC HL, HL
    SBC HL, SP
L18:
    AND A
    AND B
    AND C
    AND D
    AND E
    AND H
    AND L
    AND 0+ 5
    AND (HL)
    AND (IX+5)
    AND (IY+5)
L19:
    OR A
    OR B
    OR C
    OR D
    OR E
    OR H
    OR L
    OR 0+ 5
    OR (HL)
    OR (IX+5)
    OR (IY+5)
L20:
    XOR A
    XOR B
    XOR C
    XOR D
    XOR E
    XOR H
    XOR L
    XOR 0+ 5
    XOR (HL)
    XOR (IX+5)
    XOR (IY+5)
L21:
    CP A
    CP B
    CP C
    CP D
    CP E
    CP H
    CP L
    CP 0+ 5
    CP (HL)
    CP (IX+5)
    CP (IY+5)
L22:
    INC A
    INC B
    INC C
    INC D
    INC E
    INC H
    INC L
    INC (HL)
    INC (IX+5)
    INC (IY+5)
    INC BC
    INC DE
    INC HL
    INC SP
    INC IX
    INC IY
L23:
    DEC A
    DEC B
    DEC C
    DEC D
    DEC E
    DEC H
    DEC L
    DEC (HL)
    DEC (IX+5)
    DEC (IY+5)
    DEC BC
    DEC DE
    DEC HL
    DEC SP
    DEC IX
    DEC IY
L24:
    RLCA
L25:
    RLA
L26:
    RRCA
L27:
    RRA
L28:
    RLC A
    RLC B
    RLC C
    RLC D
    RLC E
    RLC H
    RLC L
    RLC (HL)
    RLC (IX+5)
    RLC (IY+5)
L29:
    RL A
    RL B
    RL C
    RL D
    RL E
    RL H
    RL L
    RL (HL)
    RL (IX+5)
    RL (IY+5)
L30:
    RRC A
    RRC B
    RRC C
    RRC D
    RRC E
    RRC H
    RRC L
    RRC (HL)
    RRC (IX+5)
    RRC (IY+5)
L31:
    RR A
    RR B
    RR C
    RR D
    RR E
    RR H
    RR L
    RR (HL)
    RR (IX+5)
    RR (IY+5)
L32:
    SLA A
    SLA B
    SLA C
    SLA D
    SLA E
    SLA H
    SLA L
    SLA (HL)
    SLA (IX+5)
    SLA (IY+5)
L33:
    SRA A
    SRA B
    SRA C
    SRA D
    SRA E
    SRA H
    SRA L
    SRA (HL)
    SRA (IX+5)
    SRA (IY+5)
L34:
    SRL A
    SRL B
    SRL C
    SRL D
    SRL E
    SRL H
    SRL L
    SRL (HL)
    SRL (IX+5)
    SRL (IY+5)
L35:
    RLD
L36:
    RRD
L37:
    BIT 0+ 5, A
    BIT 0+ 5, B
    BIT 0+ 5, C
    BIT 0+ 5, D
    BIT 0+ 5, E
    BIT 0+ 5, H
    BIT 0+ 5, L
    BIT 0+ 5, (HL)
    BIT 0+ 5, (IX+5)
    BIT 0+ 5, (IY+5)
L38:
    SET 0+ 5, A
    SET 0+ 5, B
    SET 0+ 5, C
    SET 0+ 5, D
    SET 0+ 5, E
    SET 0+ 5, H
    SET 0+ 5, L
    SET 0+ 5, (HL)
    SET 0+ 5, (IX+5)
    SET 0+ 5, (IY+5)
L39:
    RES 0+ 5, A
    RES 0+ 5, B
    RES 0+ 5, C
    RES 0+ 5, D
    RES 0+ 5, E
    RES 0+ 5, H
    RES 0+ 5, L
    RES 0+ 5, (HL)
    RES 0+ 5, (IX+5)
    RES 0+ 5, (IY+5)
L40:
    JP 0+ 1234
    JP (HL)
    JP (IX)
    JP (IY)
    JP NZ, 0+ 1234
    JP Z, 0+ 1234
    JP NC, 0+ 1234
    JP C, 0+ 1234
    JP PO, 0+ 1234
    JP PE, 0+ 1234
    JP P, 0+ 1234
    JP M, 0+ 1234
    JP 0+ 5
    JP NZ, 0+ 5
    JP Z, 0+ 5
    JP NC, 0+ 5
    JP C, 0+ 5
    JP PO, 0+ 5
    JP PE, 0+ 5
    JP P, 0+ 5
    JP M, 0+ 5
L41:
    JR L40
    JR L43
    JR C, L40
    JR C, L43
    JR NC, L40
    JR NC, L43
    JR Z, L40
    JR Z, L43
    JR NZ, L40
    JR NZ, L43
    JR L40
    JR L43
    JR C, L40
    JR C, L43
    JR NC, L40
    JR NC, L43
    JR Z, L40
    JR Z, L43
    JR NZ, L40
    JR NZ, L43
L42:
    DJNZ L41
    DJNZ L44
    DJNZ L41
    DJNZ L44
L43:
    CALL 0+ 1234
    CALL NZ, 0+ 1234
    CALL Z, 0+ 1234
    CALL NC, 0+ 1234
    CALL C, 0+ 1234
    CALL PO, 0+ 1234
    CALL PE, 0+ 1234
    CALL P, 0+ 1234
    CALL M, 0+ 1234
    CALL 0+ 5
    CALL NZ, 0+ 5
    CALL Z, 0+ 5
    CALL NC, 0+ 5
    CALL C, 0+ 5
    CALL PO, 0+ 5
    CALL PE, 0+ 5
    CALL P, 0+ 5
    CALL M, 0+ 5
L44:
    RET
    RET NZ
    RET Z
    RET NC
    RET C
    RET PO
    RET PE
    RET P
    RET M
L45:
    RETI
L46:
    RETN
L47:
    RST 16
L48:
    IN A, (5)
    IN A, (C)
    IN B, (C)
    IN C, (C)
    IN D, (C)
    IN E, (C)
    IN H, (C)
    IN L, (C)
L49:
    INI
L50:
    INIR
L51:
    IND
L52:
    INDR
L53:
    OUT (5), A
    OUT (C), A
    OUT (C), B
    OUT (C), C
    OUT (C), D
    OUT (C), E
    OUT (C), H
    OUT (C), L
L54:
    OUTI
L55:
    OTIR
L56:
    OUTD
L57:
    OTDR
L58:
    DAA
L59:
    CPL
L60:
    NEG
L61:
    CCF
L62:
    SCF
L63:
    NOP
L64:
    HALT
L65:
    DI
L66:
    EI
L67:
    IM 1
