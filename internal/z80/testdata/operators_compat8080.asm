    LD A, A
    LD A, B
    LD A, C
    LD A, D
    LD A, E
    LD A, H
    LD A, L
    LD A, (HL)
    LD A, (BC)
    LD A, (DE)
    LD A, 0+ 5
    LD A, (5)
    LD A, (1234)
    LD B, A
    LD B, B
    LD B, C
    LD B, D
    LD B, E
    LD B, H
    LD B, L
    LD B, (HL)
    LD B, 0+ 5
    LD C, A
    LD C, B
    LD C, C
    LD C, D
    LD C, E
    LD C, H
    LD C, L
    LD C, (HL)
    LD C, 0+ 5
    LD D, A
    LD D, B
    LD D, C
    LD D, D
    LD D, E
    LD D, H
    LD D, L
    LD D, (HL)
    LD D, 0+ 5
    LD E, A
    LD E, B
    LD E, C
    LD E, D
    LD E, E
    LD E, H
    LD E, L
    LD E, (HL)
    LD E, 0+ 5
    LD H, A
    LD H, B
    LD H, C
    LD H, D
    LD H, E
    LD H, H
    LD H, L
    LD H, (HL)
    LD H, 0+ 5
    LD L, A
    LD L, B
    LD L, C
    LD L, D
    LD L, E
    LD L, H
    LD L, L
    LD L, (HL)
    LD L, 0+ 5
    LD BC, 0+ 5
    LD BC, 0+ 1234
    LD DE, 0+ 5
    LD DE, 0+ 1234
    LD HL, 0+ 5
    LD HL, 0+ 1234
    LD SP, HL
    LD SP, 0+ 5
    LD SP, 0+ 1234
    LD A, A
    LD B, A
    LD C, A
    LD D, A
    LD E, A
    LD H, A
    LD L, A
    LD (HL), A
    LD (BC), A
    LD (DE), A
    LD (5), A
    LD (1234), A
    LD A, B
    LD B, B
    LD C, B
    LD D, B
    LD E, B
    LD H, B
    LD L, B
    LD (HL), B
    LD A, C
    LD B, C
    LD C, C
    LD D, C
    LD E, C
    LD H, C
    LD L, C
    LD (HL), C
    LD A, D
    LD B, D
    LD C, D
    LD D, D
    LD E, D
    LD H, D
    LD L, D
    LD (HL), D
    LD A, E
    LD B, E
    LD C, E
    LD D, E
    LD E, E
    LD H, E
    LD L, E
    LD (HL), E
    LD A, H
    LD B, H
    LD C, H
    LD D, H
    LD E, H
    LD H, H
    LD L, H
    LD (HL), H
    LD A, L
    LD B, L
    LD C, L
    LD D, L
    LD E, L
    LD H, L
    LD L, L
    LD (HL), L
    LD SP, HL
    LD HL, 0+ 1234
    LD BC, 0+ 1234
    LD DE, 0+ 1234
    LD SP, 0+ 1234
    LD A, 0+ 5
    LD B, 0+ 5
    LD C, 0+ 5
    LD D, 0+ 5
    LD E, 0+ 5
    LD H, 0+ 5
    LD L, 0+ 5
    LD HL, 0+ 5
    LD (HL), 0+ 5
    LD BC, 0+ 5
    LD DE, 0+ 5
    LD SP, 0+ 5
    EX DE, HL
    EX DE, HL
    EX (SP), HL
    EX (SP), HL
    PUSH HL
    PUSH BC
    PUSH DE
    PUSH AF
    POP HL
    POP BC
    POP DE
    POP AF
    INC A
    INC B
    INC C
    INC D
    INC E
    INC H
    INC L
    INC HL
    INC (HL)
    INC BC
    INC DE
    INC SP
    DEC A
    DEC B
    DEC C
    DEC D
    DEC E
    DEC H
    DEC L
    DEC HL
    DEC (HL)
    DEC BC
    DEC DE
    DEC SP
    CPL
    ADD A, A
    ADD A, B
    ADD A, C
    ADD A, D
    ADD A, E
    ADD A, H
    ADD A, L
    ADD A, (HL)
    ADD A, 0+ 5
    ADD HL, HL
    ADD HL, BC
    ADD HL, DE
    ADD HL, SP
    ADC A, A
    ADC A, B
    ADC A, C
    ADC A, D
    ADC A, E
    ADC A, H
    ADC A, L
    ADC A, (HL)
    ADC A, 0+ 5
    SUB A
    SUB B
    SUB C
    SUB D
    SUB E
    SUB H
    SUB L
    SUB (HL)
    SUB 0+ 5
    SBC A, A
    SBC A, B
    SBC A, C
    SBC A, D
    SBC A, E
    SBC A, H
    SBC A, L
    SBC A, (HL)
    SBC A, 0+ 5
    CP A
    CP B
    CP C
    CP D
    CP E
    CP H
    CP L
    CP (HL)
    CP 0+ 5
    AND A
    AND B
    AND C
    AND D
    AND E
    AND H
    AND L
    AND (HL)
    AND 0+ 5
    OR A
    OR B
    OR C
    OR D
    OR E
    OR H
    OR L
    OR (HL)
    OR 0+ 5
    XOR A
    XOR B
    XOR C
    XOR D
    XOR E
    XOR H
    XOR L
    XOR (HL)
    XOR 0+ 5
    AND A
    INC B
    DEC B
    INC C
    DEC C
    INC D
    DEC D
    INC E
    DEC E
    INC H
    DEC H
    INC L
    DEC L
    INC (HL)
    DEC (HL)
    LD B, D
    LD C, E
    LD D, B
    LD E, C
    LD B, B
    LD C, C
    LD B, D
    LD C, E
    LD B, H
    LD C, L
    LD D, B
    LD E, C
    LD D, D
    LD E, E
    LD D, H
    LD E, L
    LD H, B
    LD L, C
    LD H, D
    LD L, E
    LD H, H
    LD L, L
    LD B, B
    LD C, C
    LD D, B
    LD E, C
    LD H, B
    LD L, C
    LD B, D
    LD C, E
    LD D, D
    LD E, E
    LD H, D
    LD L, E
    LD B, H
    LD C, L
    LD D, H
    LD E, L
    LD H, H
    LD L, L
    RLCA
    RLCA
    RLA
    RLA
    RRCA
    RRCA
    RRA
    RRA
    LBI:
    JP NZ, LBI
    JP Z, LBI
    JP NC, LBI
    JP C, LBI
    JP PO, LBI
    JP PE, LBI
    JP M, LBI
    JP P, LBI
    LBU:
    JP Z, LBU
    JP NZ, LBU
    JP C, LBU
    JP NC, LBU
    JP PE, LBU
    JP PO, LBU
    JP P, LBU
    JP M, LBU
    RET
    RET NZ
    RET Z
    RET NC
    RET C
    RET PO
    RET PE
    RET M
    RET P
    RET Z
    RET NZ
    RET C
    RET NC
    RET PE
    RET PO
    RET P
    RET M
    f: RET
    CALL f
    CALL NC, f
