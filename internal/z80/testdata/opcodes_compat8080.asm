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
    LD (HL), A
    LD (HL), B
    LD (HL), C
    LD (HL), D
    LD (HL), E
    LD (HL), H
    LD (HL), L
    LD (HL), 0+ 5
    LD A, (BC)
    LD A, (DE)
    LD A, (1234)
    LD (BC), A
    LD (DE), A
    LD (1234), A
    LD BC, 0+ 1234
    LD DE, 0+ 1234
    LD HL, 0+ 1234
    LD SP, 0+ 1234
    LD SP, HL
    LD A, (5)
    LD (5), A
    LD BC, 0+ 5
    LD DE, 0+ 5
    LD HL, 0+ 5
    LD SP, 0+ 5
L2:
    PUSH BC
    PUSH DE
    PUSH HL
    PUSH AF
L3:
    POP BC
    POP DE
    POP HL
    POP AF
L4:
    EX DE, HL
    EX (SP), HL
L5:
    ADD A, A
    ADD A, B
    ADD A, C
    ADD A, D
    ADD A, E
    ADD A, H
    ADD A, L
    ADD A, 0+ 5
    ADD A, (HL)
    ADD HL, BC
    ADD HL, DE
    ADD HL, HL
    ADD HL, SP
L6:
    ADC A, A
    ADC A, B
    ADC A, C
    ADC A, D
    ADC A, E
    ADC A, H
    ADC A, L
    ADC A, 0+ 5
    ADC A, (HL)
L7:
    SUB A
    SUB B
    SUB C
    SUB D
    SUB E
    SUB H
    SUB L
    SUB 0+ 5
    SUB (HL)
L8:
    SBC A, A
    SBC A, B
    SBC A, C
    SBC A, D
    SBC A, E
    SBC A, H
    SBC A, L
    SBC A, 0+ 5
    SBC A, (HL)
L9:
    AND A
    AND B
    AND C
    AND D
    AND E
    AND H
    AND L
    AND 0+ 5
    AND (HL)
L10:
    OR A
    OR B
    OR C
    OR D
    OR E
    OR H
    OR L
    OR 0+ 5
    OR (HL)
L11:
    XOR A
    XOR B
    XOR C
    XOR D
    XOR E
    XOR H
    XOR L
    XOR 0+ 5
    XOR (HL)
L12:
    CP A
    CP B
    CP C
    CP D
    CP E
    CP H
    CP L
    CP 0+ 5
    CP (HL)
L13:
    INC A
    INC B
    INC C
    INC D
    INC E
    INC H
    INC L
    INC (HL)
    INC BC
    INC DE
    INC HL
    INC SP
L14:
    DEC A
    DEC B
    DEC C
    DEC D
    DEC E
    DEC H
    DEC L
    DEC (HL)
    DEC BC
    DEC DE
    DEC HL
    DEC SP
L15:
    RLCA
L16:
    RLA
L17:
    RRCA
L18:
    RRA
L19:
    JP 0+ 1234
    JP (HL)
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
L20:
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
L21:
    RET
    RET NZ
    RET Z
    RET NC
    RET C
    RET PO
    RET PE
    RET P
    RET M
L22:
    RST 16
L23:
L24:
L25:
    IN A, (5)
L26:
    OUT (5), A
L27:
    DAA
L28:
    CPL
L29:
    CCF
L30:
    SCF
L31:
    NOP
L32:
    HALT
L33:
    DI
L34:
    EI
