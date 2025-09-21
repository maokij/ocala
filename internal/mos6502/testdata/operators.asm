    .org 512
    TXA
    TYA
    LDA #5
    LDA 5
    LDA 5, X
    LDA 5, Y
    LDA 1234
    LDA 1234, X
    LDA 1234, Y
    LDA (5, X)
    LDA (5), Y
    TAX
    TSX
    LDX #5
    LDX 5
    LDX 5, Y
    LDX 1234
    LDX 1234, Y
    TAY
    LDY #5
    LDY 5
    LDY 5, X
    LDY 1234
    LDY 1234, X
    TXS
    TAX
    TAY
    STA 5
    STA 5, X
    STA 5, Y
    STA 1234
    STA 1234, X
    STA 1234, Y
    STA (5, X)
    STA (5), Y
    TXA
    TXS
    STX 5
    STX 5, Y
    STX 1234
    TYA
    STY 5
    STY 5, X
    STY 1234
    TSX
    PHA
    PHP
    PLA
    PLP
    CLC
    ADC #1
    INX
    INY
    INC 5
    INC 5, X
    INC 1234
    INC 1234, X
    SEC
    SBC #1
    DEX
    DEY
    DEC 5
    DEC 5, X
    DEC 1234
    DEC 1234, X
    EOR #255
    EOR #255
    CLC
    ADC #1
    CLC
    ADC #5
    CLC
    ADC 5
    CLC
    ADC 5, X
    CLC
    ADC 5, Y
    CLC
    ADC 1234
    CLC
    ADC 1234, X
    CLC
    ADC 1234, Y
    CLC
    ADC (5, X)
    CLC
    ADC (5), Y
    ADC #5
    ADC 5
    ADC 5, X
    ADC 5, Y
    ADC 1234
    ADC 1234, X
    ADC 1234, Y
    ADC (5, X)
    ADC (5), Y
    SEC
    SBC #5
    SEC
    SBC 5
    SEC
    SBC 5, X
    SEC
    SBC 5, Y
    SEC
    SBC 1234
    SEC
    SBC 1234, X
    SEC
    SBC 1234, Y
    SEC
    SBC (5, X)
    SEC
    SBC (5), Y
    SBC #5
    SBC 5
    SBC 5, X
    SBC 5, Y
    SBC 1234
    SBC 1234, X
    SBC 1234, Y
    SBC (5, X)
    SBC (5), Y
    CMP #5
    CMP 5
    CMP 5, X
    CMP 5, Y
    CMP 1234
    CMP 1234, X
    CMP 1234, Y
    CMP (5, X)
    CMP (5), Y
    CPX #5
    CPX 5
    CPX 1234
    CPY #5
    CPY 5
    CPY 1234
    AND #5
    AND 5
    AND 5, X
    AND 5, Y
    AND 1234
    AND 1234, X
    AND 1234, Y
    AND (5, X)
    AND (5), Y
    ORA #5
    ORA 5
    ORA 5, X
    ORA 5, Y
    ORA 1234
    ORA 1234, X
    ORA 1234, Y
    ORA (5, X)
    ORA (5), Y
    EOR #5
    EOR 5
    EOR 5, X
    EOR 5, Y
    EOR 1234
    EOR 1234, X
    EOR 1234, Y
    EOR (5, X)
    EOR (5), Y
    BIT 5
    BIT 1234
    CMP #128
    ROL A
    CMP #128
    ROL A
    ROL A
    ROL A
    ROL 5
    ROL 5
    LSR A
    BCC :+
    ORA #128
    :
    LSR A
    BCC :+
    ORA #128
    :
    ROR A
    ROR A
    ROR 5
    ROR 5
    ASL A
    ASL A
    ASL 5
    ASL 5
    CMP #128
    ROR A
    CMP #128
    ROR A
    LSR A
    LSR A
    LSR 5
    LSR 5
    LBI:
    BEQ :+
    JMP LBI
    :
    BNE :+
    JMP LBI
    :
    BCS :+
    JMP LBI
    :
    BCC :+
    JMP LBI
    :
    BVS :+
    JMP LBI
    :
    BVC :+
    JMP LBI
    :
    BMI :+
    JMP LBI
    :
    BPL :+
    JMP LBI
    :
    LBU:
    BNE :+
    JMP LBU
    :
    BEQ :+
    JMP LBU
    :
    BCC :+
    JMP LBU
    :
    BCS :+
    JMP LBU
    :
    BVC :+
    JMP LBU
    :
    BVS :+
    JMP LBU
    :
    BPL :+
    JMP LBU
    :
    BMI :+
    JMP LBU
    :
    RTS
    BEQ :+
    RTS
    :
    BNE :+
    RTS
    :
    BCS :+
    RTS
    :
    BCC :+
    RTS
    :
    BVS :+
    RTS
    :
    BVC :+
    RTS
    :
    BMI :+
    RTS
    :
    BPL :+
    RTS
    :
    BNE :+
    RTS
    :
    BEQ :+
    RTS
    :
    BCC :+
    RTS
    :
    BCS :+
    RTS
    :
    BVC :+
    RTS
    :
    BVS :+
    RTS
    :
    BPL :+
    RTS
    :
    BMI :+
    RTS
    :
