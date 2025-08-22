;; -*- mode: Lisp; lisp-indent-offset: 2 -*-
(arch mos6502
  (operand A   RegA  "A"        "A")
  (operand X   RegX  "X"        "X")
  (operand Y   RegY  "Y"        "Y")
  (operand S   RegS  "S"        "S")
  (operand P   RegP  "P"        "P")

  (operand N   ImmN  "%B"       "#%B" NN temp)
  (operand NN  ImmNN "%W"       "#%W" N)

  (operand ZN  MemZN "[%B]"     "%B" AN temp)
  (operand ZX  MemZX "[%B X]"   "%B, X" AX temp)
  (operand ZY  MemZY "[%B Y]"   "%B, Y" AY temp)
  (operand AN  MemAN "[%W]"     "%W" ZN)
  (operand AX  MemAX "[%W X]"   "%W, X" ZX)
  (operand AY  MemAY "[%W Y]"   "%W, Y" ZY)
  (operand IN  MemIN "[[%W]]"   "(%W)")
  (operand IX  MemIX "[[%B X]]" "(%B, X)")
  (operand IY  MemIY "[[%B] Y]" "(%B), Y")

  (operand NE? CondNE "NE?" "NE")
  (operand EQ? CondEQ "EQ?" "EQ")
  (operand CC? CondCC "CC?" "CC")
  (operand CS? CondCS "CS?" "CS")
  (operand VC? CondVC "VC?" "VC")
  (operand VS? CondVS "VS?" "VS")
  (operand PL? CondPL "PL?" "PL")
  (operand MI? CondMI "MI?" "MI")

  (registers A X Y S P)
  (conditions
    (NE? !=? not-zero?)
    (EQ? ==? zero?)
    (CC? <? not-carry? borrow?)
    (CS? >=? carry? not-borrow?)
    (VC? not-over?)
    (VS? over?)
    (PL? plus?)
    (MI? minus?))
  (map CC NE? 0 EQ? 1 CC? 2 CS? 3 VC? 4 VS? 5 PL? 6 MI? 7)

  (example "$prologue" (*)
    "arch mos6502; link { org 512 0 1 ; merge text _ }; flat!; optimize near-jump 0"
    ".org 512")

  (opcode  LDA (a)
    (N)  [0xA9 (=l a)]
    (ZN) [0xA5 (=l a)]
    (ZX) [0xB5 (=l a)]
    (AN) [0xAD (=l a) (=h a)]
    (AX) [0xBD (=l a) (=h a)]
    (AY) [0xB9 (=l a) (=h a)]
    (IX) [0xA1 (=l a)]
    (IY) [0xB1 (=l a)])

  (opcode  LDX (a)
    (N)  [0xA2 (=l a)]
    (ZN) [0xA6 (=l a)]
    (ZY) [0xB6 (=l a)]
    (AN) [0xAE (=l a) (=h a)]
    (AY) [0xBE (=l a) (=h a)])

  (opcode  LDY (a)
    (N)  [0xA0 (=l a)]
    (ZN) [0xA4 (=l a)]
    (ZX) [0xB4 (=l a)]
    (AN) [0xAC (=l a) (=h a)]
    (AX) [0xBC (=l a) (=h a)])

  (opcode  STA (a)
    (ZN) [0x85 (=l a)]
    (ZX) [0x95 (=l a)]
    (AN) [0x8D (=l a) (=h a)]
    (AX) [0x9D (=l a) (=h a)]
    (AY) [0x99 (=l a) (=h a)]
    (IX) [0x81 (=l a)]
    (IY) [0x91 (=l a)])

  (opcode  STX (a)
    (ZN) [0x86 (=l a)]
    (ZY) [0x96 (=l a)]
    (AN) [0x8E (=l a) (=h a)])

  (opcode  STY (a)
    (ZN) [0x84 (=l a)]
    (ZX) [0x94 (=l a)]
    (AN) [0x8C (=l a) (=h a)])

  (opcode  TAX () () [0xAA])
  (opcode  TAY () () [0xA8])
  (opcode  TSX () () [0xBA])
  (opcode  TXA () () [0x8A])
  (opcode  TXS () () [0x9A])
  (opcode  TYA () () [0x98])

  (opcode  PHA () () [0x48])
  (opcode  PHP () () [0x08])

  (opcode  PLP () () [0x28])
  (opcode  PLA () () [0x68])

  (opcode  CLC () () [0x18])
  (opcode  CLI () () [0x58])
  (opcode  CLD () () [0xD8])
  (opcode  CLV () () [0xB8])

  (opcode  SEC () () [0x38])
  (opcode  SEI () () [0x78])
  (opcode  SED () () [0xF8])

  (opcode  BRK () () [0x00])
  (opcode  NOP () () [0xEA])

  (opcode  RTS () () [0x60])
  (opcode  RTI () () [0x40])

  (opcode  JMP (a)
    (AN) [0x4C (=l a) (=h a)]
    (IN) [0x6C (=l a) (=h a)])
  (opcode  JSR (a) (AN) [0x20 (=l a) (=h a)])

  (opcode  BPL (a) (NN) [0x10 (=rl a -2)])
  (opcode  BMI (a) (NN) [0x30 (=rl a -2)])
  (opcode  BVC (a) (NN) [0x50 (=rl a -2)])
  (opcode  BVS (a) (NN) [0x70 (=rl a -2)])
  (opcode  BCC (a) (NN) [0x90 (=rl a -2)])
  (opcode  BCS (a) (NN) [0xB0 (=rl a -2)])
  (opcode  BNE (a) (NN) [0xD0 (=rl a -2)])
  (opcode  BEQ (a) (NN) [0xF0 (=rl a -2)])

  (opcode  #.jump (a) (NN) [0x4C (=l a) (=h a)])
  (opcode  #.jump (a b)
    (NN PL?) [0x30 0x03 0x4C (=l a) (=h a)]
    (NN MI?) [0x10 0x03 0x4C (=l a) (=h a)]
    (NN VC?) [0x70 0x03 0x4C (=l a) (=h a)]
    (NN VS?) [0x50 0x03 0x4C (=l a) (=h a)]
    (NN CC?) [0xB0 0x03 0x4C (=l a) (=h a)]
    (NN CS?) [0x90 0x03 0x4C (=l a) (=h a)]
    (NN NE?) [0xF0 0x03 0x4C (=l a) (=h a)]
    (NN EQ?) [0xD0 0x03 0x4C (=l a) (=h a)])
  (example #.jump (*) "" "")

  (opcode  #.call (a)   (NN)    [0x20 (=l a) (=h a)])
  (opcode  #.call (a b) (NN CC) [(=U)])
  (example #.call (*) "" "")

  ;;
  (opcode  ORA (a)
    (N)  [0x09 (=l a)]
    (ZN) [0x05 (=l a)]
    (ZX) [0x15 (=l a)]
    (AN) [0x0D (=l a) (=h a)]
    (AX) [0x1D (=l a) (=h a)]
    (AY) [0x19 (=l a) (=h a)]
    (IX) [0x01 (=l a)]
    (IY) [0x11 (=l a)])

  (opcode  AND (a)
    (N)  [0x29 (=l a)]
    (ZN) [0x25 (=l a)]
    (ZX) [0x35 (=l a)]
    (AN) [0x2D (=l a) (=h a)]
    (AX) [0x3D (=l a) (=h a)]
    (AY) [0x39 (=l a) (=h a)]
    (IX) [0x21 (=l a)]
    (IY) [0x31 (=l a)])

  (opcode  EOR (a)
    (N)  [0x49 (=l a)]
    (ZN) [0x45 (=l a)]
    (ZX) [0x55 (=l a)]
    (AN) [0x4D (=l a) (=h a)]
    (AX) [0x5D (=l a) (=h a)]
    (AY) [0x59 (=l a) (=h a)]
    (IX) [0x41 (=l a)]
    (IY) [0x51 (=l a)])

  (opcode  ADC (a)
    (N)  [0x69 (=l a)]
    (ZN) [0x65 (=l a)]
    (ZX) [0x75 (=l a)]
    (AN) [0x6D (=l a) (=h a)]
    (AX) [0x7D (=l a) (=h a)]
    (AY) [0x79 (=l a) (=h a)]
    (IX) [0x61 (=l a)]
    (IY) [0x71 (=l a)])

  (opcode  CMP (a)
    (N)  [0xC9 (=l a)]
    (ZN) [0xC5 (=l a)]
    (ZX) [0xD5 (=l a)]
    (AN) [0xCD (=l a) (=h a)]
    (AX) [0xDD (=l a) (=h a)]
    (AY) [0xD9 (=l a) (=h a)]
    (IX) [0xC1 (=l a)]
    (IY) [0xD1 (=l a)])

  (opcode  SBC (a)
    (N)  [0xE9 (=l a)]
    (ZN) [0xE5 (=l a)]
    (ZX) [0xF5 (=l a)]
    (AN) [0xED (=l a) (=h a)]
    (AX) [0xFD (=l a) (=h a)]
    (AY) [0xF9 (=l a) (=h a)]
    (IX) [0xE1 (=l a)]
    (IY) [0xF1 (=l a)])

  (opcode  BIT (a)
    (ZN) [0x24 (=l a)]
    (AN) [0x2C (=l a) (=h a)])

  (opcode  CPX (a)
    (N)  [0xE0 (=l a)]
    (ZN) [0xE4 (=l a)]
    (AN) [0xEC (=l a) (=h a)])

  (opcode  CPY (a)
    (N)  [0xC0 (=l a)]
    (ZN) [0xC4 (=l a)]
    (AN) [0xCC (=l a) (=h a)])

  (opcode  INC (a)
    (ZN) [0xE6 (=l a)]
    (ZX) [0xF6 (=l a)]
    (AN) [0xEE (=l a) (=h a)]
    (AX) [0xFE (=l a) (=h a)])
  (opcode  INX () ()  [0xE8])
  (opcode  INY () ()  [0xC8])

  (opcode  DEC (a)
    (ZN) [0xC6 (=l a)]
    (ZX) [0xD6 (=l a)]
    (AN) [0xCE (=l a) (=h a)]
    (AX) [0xDE (=l a) (=h a)])
  (opcode  DEX () ()  [0xCA])
  (opcode  DEY () ()  [0x88])

  ;;
  (opcode  ASL (a)
    (A)  [0x0A]
    (ZN) [0x06 (=l a)]
    (ZX) [0x16 (=l a)]
    (AN) [0x0E (=l a) (=h a)]
    (AX) [0x1E (=l a) (=h a)])

  (opcode  LSR (a)
    (A)  [0x4A]
    (ZN) [0x46 (=l a)]
    (ZX) [0x56 (=l a)]
    (AN) [0x4E (=l a) (=h a)]
    (AX) [0x5E (=l a) (=h a)])

  (opcode  ROL (a)
    (A)  [0x2A]
    (ZN) [0x26 (=l a)]
    (ZX) [0x36 (=l a)]
    (AN) [0x2E (=l a) (=h a)]
    (AX) [0x3E (=l a) (=h a)])

  (opcode  ROR (a)
    (A)  [0x6A]
    (ZN) [0x66 (=l a)]
    (ZX) [0x76 (=l a)]
    (AN) [0x6E (=l a) (=h a)]
    (AX) [0x7E (=l a) (=h a)])

  ;;
  (operator <- (a b)
    (A X) [(TXA)]
    (A Y) [(TYA)]
    (A _) [(LDA (= b))]
    (X A) [(TAX)]
    (X S) [(TSX)]
    (X _) [(LDX (= b))]
    (Y A) [(TAY)]
    (Y _) [(LDY (= b))]
    (S X) [(TXS)])
  (operator -> (a b)
    (A X) [(TAX)]
    (A Y) [(TAY)]
    (A _) [(STA (= b))]
    (X A) [(TXA)]
    (X S) [(TXS)]
    (X _) [(STX (= b))]
    (Y A) [(TYA)]
    (Y _) [(STY (= b))]
    (S X) [(TSX)])

  (operator -push (a)
    (A) [(PHA)]
    (P) [(PHP)])
  (operator -pop  (a)
    (A) [(PLA)]
    (P) [(PLP)])

  (operator ++ (a)
    (A) [(CLC) (ADC 1)]
    (X) [(INX)]
    (Y) [(INY)]
    (_) [(INC (= a))])
  (operator -- (a)
    (A) [(SEC) (SBC 1)]
    (X) [(DEX)]
    (Y) [(DEY)]
    (_) [(DEC (= a))])
  (operator -not (a) (A) [(EOR 0xFF)])
  (operator -neg (a) (A) [(EOR 0xFF) (CLC) (ADC 1)])

  (operator +  (a b) (A _) [(CLC) (ADC (= b))])
  (operator +$ (a b) (A _) [(ADC (= b))])

  (operator -  (a b) (A _) [(SEC) (SBC (= b))])
  (operator -$ (a b) (A _) [(SBC (= b))])
  (operator -? (a b)
    (A _) [(CMP (= b))]
    (X _) [(CPX (= b))]
    (Y _) [(CPY (= b))])

  (operator &     (a b) (A _) [(AND (= b))])
  (operator "|"   (a b) (A _) [(ORA (= b))])
  (operator ^     (a b) (A _) [(EOR (= b))])
  (operator -bit? (a b) (A _) [(BIT (= b))])

  (operator <*   (a b) (A NN) [(#.REP (= b) `[(CMP 0x80) (ROL A)])])
  (operator <*$  (a b) (_ NN) [(#.REP (= b) `[(ROL (= a))])])
  (operator >*   (a b) (A NN) [(#.REP (= b) `[(LSR A) (#.BYTE 0x90) (#.BYTE 2) (ORA 0x80)])]) ; BCC +2
  (operator >*$  (a b) (_ NN) [(#.REP (= b) `[(ROR (= a))])])

  (operator <<   (a b) (_ NN) [(#.REP (= b) `[(ASL (= a))])])
  (operator >>   (a b) (A NN) [(#.REP (= b) `[(CMP 0x80) (ROR A)])])
  (operator >>>  (a b) (_ NN) [(#.REP (= b) `[(LSR (= a))])])

  (operator -jump (a) (NN) [(#.jump (= a))])
  (example -jump (*) "" "")

  (operator -jump-if (a b) (NN CC) [(#.jump (= a) (= b))])
  (example -jump-if (*) "" "")

  (operator -jump-unless (a b)
    (NN NE?) [(#.jump (= a) EQ?)]
    (NN EQ?) [(#.jump (= a) NE?)]
    (NN CC?) [(#.jump (= a) CS?)]
    (NN CS?) [(#.jump (= a) CC?)]
    (NN VC?) [(#.jump (= a) VS?)]
    (NN VS?) [(#.jump (= a) VC?)]
    (NN PL?) [(#.jump (= a) MI?)]
    (NN MI?) [(#.jump (= a) PL?)])
  (example -jump-unless (*) "" "")

  (example $operators (*)
    "A <* 2"    "CMP #128; ROL A; CMP #128; ROL A"
    "A <*$ 2"   "ROL A; ROL A"
    "[5] <*$ 2" "ROL 5; ROL 5"
    "A >* 2"    "LSR A; BCC :+; ORA #128; :; LSR A; BCC :+; ORA #128; :"
    "A >*$ 2"   "ROR A; ROR A"
    "[5] >*$ 2" "ROR 5; ROR 5"
    "A << 2"    "ASL A; ASL A"
    "[5] << 2"  "ASL 5; ASL 5"
    "A >> 2"    "CMP #128; ROR A; CMP #128; ROR A"
    "A >>> 2"   "LSR A; LSR A"
    "[5] >>> 2" "LSR 5; LSR 5"

    "LBI:; $(LBI) -jump-if NE?" "LBI:; BEQ :+; JMP LBI; :"
    "$(LBI) -jump-if EQ?" "BNE :+; JMP LBI; :"
    "$(LBI) -jump-if CC?" "BCS :+; JMP LBI; :"
    "$(LBI) -jump-if CS?" "BCC :+; JMP LBI; :"
    "$(LBI) -jump-if VC?" "BVS :+; JMP LBI; :"
    "$(LBI) -jump-if VS?" "BVC :+; JMP LBI; :"
    "$(LBI) -jump-if PL?" "BMI :+; JMP LBI; :"
    "$(LBI) -jump-if MI?" "BPL :+; JMP LBI; :"

    "LBU:; $(LBU) -jump-unless NE?" "LBU:; BNE :+; JMP LBU; :"
    "$(LBU) -jump-unless EQ?" "BEQ :+; JMP LBU; :"
    "$(LBU) -jump-unless CC?" "BCC :+; JMP LBU; :"
    "$(LBU) -jump-unless CS?" "BCS :+; JMP LBU; :"
    "$(LBU) -jump-unless VC?" "BVC :+; JMP LBU; :"
    "$(LBU) -jump-unless VS?" "BVS :+; JMP LBU; :"
    "$(LBU) -jump-unless PL?" "BPL :+; JMP LBU; :"
    "$(LBU) -jump-unless MI?" "BMI :+; JMP LBU; :"))
