    LD A, IXH
    LD A, IXL
    LD A, IYH
    LD A, IYL
    LD B, IXH
    LD B, IXL
    LD B, IYH
    LD B, IYL
    LD C, IXH
    LD C, IXL
    LD C, IYH
    LD C, IYL
    LD D, IXH
    LD D, IXL
    LD D, IYH
    LD D, IYL
    LD E, IXH
    LD E, IXL
    LD E, IYH
    LD E, IYL
    LD IXH, A
    LD IXH, B
    LD IXH, C
    LD IXH, D
    LD IXH, E
    LD IXH, 0+ 5
    LD IXH, IXH
    LD IXH, IXL
    LD IXL, A
    LD IXL, B
    LD IXL, C
    LD IXL, D
    LD IXL, E
    LD IXL, 0+ 5
    LD IXL, IXH
    LD IXL, IXL
    LD IYH, A
    LD IYH, B
    LD IYH, C
    LD IYH, D
    LD IYH, E
    LD IYH, 0+ 5
    LD IYH, IYH
    LD IYH, IYL
    LD IYL, A
    LD IYL, B
    LD IYL, C
    LD IYL, D
    LD IYL, E
    LD IYL, 0+ 5
    LD IYL, IYH
    LD IYL, IYL
    LD IXH, A
    LD IXL, A
    LD IYH, A
    LD IYL, A
    LD IXH, B
    LD IXL, B
    LD IYH, B
    LD IYL, B
    LD IXH, C
    LD IXL, C
    LD IYH, C
    LD IYL, C
    LD IXH, D
    LD IXL, D
    LD IYH, D
    LD IYL, D
    LD IXH, E
    LD IXL, E
    LD IYH, E
    LD IYL, E
    LD IXH, 0+ 5
    LD IXL, 0+ 5
    LD IYH, 0+ 5
    LD IYL, 0+ 5
    LD A, IXH
    LD B, IXH
    LD C, IXH
    LD D, IXH
    LD E, IXH
    LD IXH, IXH
    LD IXL, IXH
    LD A, IXL
    LD B, IXL
    LD C, IXL
    LD D, IXL
    LD E, IXL
    LD IXH, IXL
    LD IXL, IXL
    LD A, IYH
    LD B, IYH
    LD C, IYH
    LD D, IYH
    LD E, IYH
    LD IYH, IYH
    LD IYL, IYH
    LD A, IYL
    LD B, IYL
    LD C, IYL
    LD D, IYL
    LD E, IYL
    LD IYH, IYL
    LD IYL, IYL
    EX AF, AF'
    EX DE, HL
    EX DE, HL
    EX (SP), HL
    EX (SP), HL
    PUSH HL
    PUSH BC
    PUSH DE
    PUSH AF
    PUSH IX
    PUSH IY
    POP HL
    POP BC
    POP DE
    POP AF
    POP IX
    POP IY
    INC IXH
    INC IXL
    INC IYH
    INC IYL
    DEC IXH
    DEC IXL
    DEC IYH
    DEC IYL
    CPL 
    NEG 
    ADD A, IXH
    ADD A, IXL
    ADD A, IYH
    ADD A, IYL
    ADC A, IXH
    ADC A, IXL
    ADC A, IYH
    ADC A, IYL
    SUB IXH
    SUB IXL
    SUB IYH
    SUB IYL
    SBC A, IXH
    SBC A, IXL
    SBC A, IYH
    SBC A, IYL
    CP IXH
    CP IXL
    CP IYH
    CP IYL
    AND IXH
    AND IXL
    AND IYH
    AND IYL
    OR IXH
    OR IXL
    OR IYH
    OR IYL
    XOR IXH
    XOR IXL
    XOR IYH
    XOR IYL
    DB 0xED, 0x70
    DB 0xED, 0x71
    INC IXH
    DEC IXH
    INC IXL
    DEC IXL
    INC IYH
    DEC IYH
    INC IYL
    DEC IYL
    RLCA
    RLCA
    RLC B
    RLC B
    RLA
    RLA
    RL B
    RL B
    RRCA
    RRCA
    RRC B
    RRC B
    RRA
    RRA
    RR B
    RR B
    SLA A
    SLA A
    SLA B
    SLA B
    SRA A
    SRA A
    SRA B
    SRA B
    SRL A
    SRL A
    SRL B
    SRL B
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
    f: RET
    CALL f
    CALL NC, f
