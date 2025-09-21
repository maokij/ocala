;; -*- mode: Lisp; lisp-indent-offset: 2 -*-
(arch z80
  (operand A   RegA   "A"    "A")
  (operand B   RegB   "B"    "B")
  (operand C   RegC   "C"    "C")
  (operand D   RegD   "D"    "D")
  (operand E   RegE   "E"    "E")
  (operand H   RegH   "H"    "H")
  (operand L   RegL   "L"    "L")

  (operand HL  RegHL  "HL"   "HL")
  (operand HL$ MemHL  "[HL]" "(HL)")
  (operand BC  RegBC  "BC"   "BC")
  (operand BC$ MemBC  "[BC]" "(BC)")
  (operand DE  RegDE  "DE"   "DE")
  (operand DE$ MemDE  "[DE]" "(DE)")
  (operand AF  RegAF  "AF"   "AF")
  (operand AF- AltAF  "AF-"  "AF'")
  (operand SP  RegSP  "SP"   "SP")
  (operand SP$ MemSP  "[SP]" "(SP)")
  (operand PC  RegPC  "PC"   "PC")
  (operand PQ  RegPQ  "PQ"   "PQ")

  (operand IX  RegIX  "IX"      "IX")
  (operand IX$ MemIX  "[IX %B]" "(IX+%B)" WX$ temp)
  (operand WX$ MemWX  "[IX %W]" "(IX+%W)")
  (operand IY  RegIY  "IY"      "IY")
  (operand IY$ MemIY  "[IY %B]" "(IY+%B)" WY$ temp)
  (operand WY$ MemWY  "[IY %W]" "(IY+%W)")

  (operand N   ImmN   "%B"   "0+ %B" NN temp)
  (operand N$  MemN   "[%B]" "(%B)"  NN$ temp)
  (operand NN  ImmNN  "%W"   "0+ %W" N)
  (operand NN$ MemNN  "[%W]" "(%W)"  N$)
  (operand C$  MemC   "[C]"  "(C)")
  (operand I   RegI   "I"    "I")
  (operand R   RegR   "R"    "R")
  (operand F   RegF   "F"    "F")

  (operand NZ? CondNZ "NZ?"  "NZ")
  (operand Z?  CondZ  "Z?"   "Z")
  (operand NC? CondNC "NC?"  "NC")
  (operand C?  CondC  "C?"   "C")
  (operand PO? CondPO "PO?"  "PO")
  (operand PE? CondPE "PE?"  "PE")
  (operand P?  CondP  "P?"   "P")
  (operand M?  CondM  "M?"   "M")

  (registers A B C D E H L I R F AF AF- BC DE HL IX IY SP PC)
  (conditions
    (NZ? !=? not-zero?)
    (Z?  ==? zero?)
    (NC? >=? not-carry?)
    (C?  <? carry?)
    (PO? odd? not-over?)
    (PE? even? over?)
    (P?  plus?)
    (M?  minus?))

  (map R8 A 7 B 0 C 1 D 2 E 3 H 4 L 5)
  (map QQ BC 0 DE 1 HL 2 AF 3)
  (map DD BC 0 DE 1 HL 2 SP 3)
  (map PP BC 0 DE 1 IX 2 SP 3)
  (map RR BC 0 DE 1 IY 2 SP 3)
  (map CC NZ? 0 Z? 1 NC? 2 C? 3 PO? 4 PE? 5 P? 6 M? 7)

  (example "$prologue" (*) "arch z80; flat!; optimize near-jump 0" "")

  (opcode  LD (a b)
    (R8 R8)  [(+ 0b0100_0000 (R8 a 3) (R8 b))]
    (R8 N)   [(+ 0b0000_0110 (R8 a 3)) (=l b)]

    (R8 HL$) [(+ 0b0100_0110 (R8 a 3))]
    (R8 IX$) [0xDD (+ 0b0100_0110 (R8 a 3)) (=l b)]
    (R8 IY$) [0xFD (+ 0b0100_0110 (R8 a 3)) (=l b)]

    (HL$ R8) [(+ 0b0111_0000 (R8 b))]
    (IX$ R8) [0xDD (+ 0b0111_0000 (R8 b)) (=l a)]
    (IY$ R8) [0xFD (+ 0b0111_0000 (R8 b)) (=l a)]

    (HL$ N)  [0x36 (=l b)]
    (IX$ N)  [0xDD 0x36 (=l a) (=l b)]
    (IY$ N)  [0xFD 0x36 (=l a) (=l b)]

    (A BC$)  [0x0A]
    (A DE$)  [0x1A]
    (A NN$)  [0x3A (=l b) (=h b)]
    (BC$ A)  [0x02]
    (DE$ A)  [0x12]
    (NN$ A)  [0x32 (=l a) (=h a)]

    (A I)    [0xED 0x57]
    (A R)    [0xED 0x5F]
    (I A)    [0xED 0x47]
    (R A)    [0xED 0x4F]

    (DD NN)  [(+ 0b0000_0001 (DD a 4)) (=l b) (=h b)]
    (IX NN)  [0xDD 0x21 (=l b) (=h b)]
    (IY NN)  [0xFD 0x21 (=l b) (=h b)]

    (HL NN$) [0x2A (=l b) (=h b)]
    (DD NN$) [0xED (+ 0b0100_1011 (DD a 4)) (=l b) (=h b)]
    (IX NN$) [0xDD 0x2A (=l b) (=h b)]
    (IY NN$) [0xFD 0x2A (=l b) (=h b)]

    (NN$ HL) [0x22 (=l a) (=h a)]
    (NN$ DD) [0xED (+ 0b0100_0011 (DD b 4)) (=l a) (=h a)]

    (NN$ IX) [0xDD 0x22 (=l a) (=h a)]
    (NN$ IY) [0xFD 0x22 (=l a) (=h a)]

    (SP HL)  [0xF9]
    (SP IX)  [0xDD 0xF9]
    (SP IY)  [0xFD 0xF9])

  (opcode  PUSH (a)
    (QQ) [(+ 0b1100_0101 (QQ a 4))]
    (IX) [0xDD 0xE5]
    (IY) [0xFD 0xE5])

  (opcode  POP (a)
    (QQ) [(+ 0b1100_0001 (QQ a 4))]
    (IX) [0xDD 0xE1]
    (IY) [0xFD 0xE1])

  (opcode  EX (a b)
    (DE  HL)  [0xEB]
    (AF  AF-) [0x08]
    (SP$ HL)  [0xE3]
    (SP$ IX)  [0xDD 0xE3]
    (SP$ IY)  [0xFD 0xE3])
  (opcode  EXX  () () [0xD9])

  (opcode  LDI  () () [0xED 0xA0])
  (opcode  LDIR () () [0xED 0xB0])
  (opcode  LDD  () () [0xED 0xA8])
  (opcode  LDDR () () [0xED 0xB8])

  (opcode  CPI  () () [0xED 0xA1])
  (opcode  CPIR () () [0xED 0xB1])
  (opcode  CPD  () () [0xED 0xA9])
  (opcode  CPDR () () [0xED 0xB9])

  (opcode  ADD (a b)
    (A R8)  [(+ 0b1000_0000 (R8 b))]
    (A N)   [0xC6 (=l b)]
    (A HL$) [0x86]
    (A IX$) [0xDD 0x86 (=l b)]
    (A IY$) [0xFD 0x86 (=l b)]

    (HL DD) [(+ 0b0000_1001 (DD b 4))]
    (IX PP) [0xDD (+ 0b0000_1001 (PP b 4))]
    (IY RR) [0xFD (+ 0b0000_1001 (RR b 4))])

  (opcode  ADC (a b)
    (A R8)  [(+ 0b1000_1000 (R8 b))]
    (A N)   [0xCE (=l b)]
    (A HL$) [0x8E]
    (A IX$) [0xDD 0x8E (=l b)]
    (A IY$) [0xFD 0x8E (=l b)]

    (HL DD) [0xED (+ 0b0100_1010 (DD b 4))])

  (opcode  SUB (a)
    (R8)  [(+ 0b1001_0000 (R8 a))]
    (N)   [0xD6 (=l a)]
    (HL$) [0x96]
    (IX$) [0xDD 0x96 (=l a)]
    (IY$) [0xFD 0x96 (=l a)])

  (opcode  SBC (a b)
    (A R8)  [(+ 0b1001_1000 (R8 b))]
    (A N)   [0xDE (=l b)]
    (A HL$) [0x9E]
    (A IX$) [0xDD 0x9E (=l b)]
    (A IY$) [0xFD 0x9E (=l b)]

    (HL DD) [0xED (+ 0b0100_0010 (DD b 4))])

  (opcode  AND (a)
    (R8)  [(+ 0b1010_0000 (R8 a))]
    (N)   [0xE6 (=l a)]
    (HL$) [0xA6]
    (IX$) [0xDD 0xA6 (=l a)]
    (IY$) [0xFD 0xA6 (=l a)])

  (opcode  OR (a)
    (R8)  [(+ 0b1011_0000 (R8 a))]
    (N)   [0xF6 (=l a)]
    (HL$) [0xB6]
    (IX$) [0xDD 0xB6 (=l a)]
    (IY$) [0xFD 0xB6 (=l a)])

  (opcode  XOR (a)
    (R8)  [(+ 0b1010_1000 (R8 a))]
    (N)   [0xEE (=l a)]
    (HL$) [0xAE]
    (IX$) [0xDD 0xAE (=l a)]
    (IY$) [0xFD 0xAE (=l a)])

  (opcode  CP (a)
    (R8)  [(+ 0b1011_1000 (R8 a))]
    (N)   [0xFE (=l a)]
    (HL$) [0xBE]
    (IX$) [0xDD 0xBE (=l a)]
    (IY$) [0xFD 0xBE (=l a)])

  (opcode  INC (a)
    (R8)  [(+ 0b0000_0100 (R8 a 3))]
    (HL$) [0x34]
    (IX$) [0xDD 0x34 (=l a)]
    (IY$) [0xFD 0x34 (=l a)]

    (DD)  [(+ 0b0000_0011 (DD a 4))]
    (IX) [0xDD 0x23]
    (IY) [0xFD 0x23])

  (opcode  DEC (a)
    (R8)  [(+ 0b0000_0101 (R8 a 3))]
    (HL$) [0x35]
    (IX$) [0xDD 0x35 (=l a)]
    (IY$) [0xFD 0x35 (=l a)]

    (DD)  [(+ 0b0000_1011 (DD a 4))]
    (IX)  [0xDD 0x2B]
    (IY)  [0xFD 0x2B])

  (opcode  RLCA () () [0x07])
  (opcode  RLA  () () [0x17])
  (opcode  RRCA () () [0x0F])
  (opcode  RRA  () () [0x1F])

  (opcode  RLC (a)
    (R8)  [0xCB (+ 0b0000_0000 (R8 a))]
    (HL$) [0xCB 0x06]
    (IX$) [0xDD 0xCB (=l a) 0x06]
    (IY$) [0xFD 0xCB (=l a) 0x06])

  (opcode  RL (a)
    (R8)  [0xCB (+ 0b0001_0000 (R8 a))]
    (HL$) [0xCB 0x16]
    (IX$) [0xDD 0xCB (=l a) 0x16]
    (IY$) [0xFD 0xCB (=l a) 0x16])

  (opcode  RRC (a)
    (R8)  [0xCB (+ 0b0000_1000 (R8 a))]
    (HL$) [0xCB 0x0E]
    (IX$) [0xDD 0xCB (=l a) 0x0E]
    (IY$) [0xFD 0xCB (=l a) 0x0E])

  (opcode  RR (a)
    (R8)  [0xCB (+ 0b0001_1000 (R8 a))]
    (HL$) [0xCB 0x1E]
    (IX$) [0xDD 0xCB (=l a) 0x1E]
    (IY$) [0xFD 0xCB (=l a) 0x1E])

  (opcode  SLA (a)
    (R8)  [0xCB (+ 0b0010_0000 (R8 a))]
    (HL$) [0xCB 0x26]
    (IX$) [0xDD 0xCB (=l a) 0x26]
    (IY$) [0xFD 0xCB (=l a) 0x26])

  (opcode  SRA (a)
    (R8)  [0xCB (+ 0b0010_1000 (R8 a))]
    (HL$) [0xCB 0x2E]
    (IX$) [0xDD 0xCB (=l a) 0x2E]
    (IY$) [0xFD 0xCB (=l a) 0x2E])

  (opcode  SRL (a)
    (R8)  [0xCB (+ 0b0011_1000 (R8 a))]
    (HL$) [0xCB 0x3E]
    (IX$) [0xDD 0xCB (=l a) 0x3E]
    (IY$) [0xFD 0xCB (=l a) 0x3E])

  (opcode  RLD  () () [0xED 0x6F])
  (opcode  RRD  () () [0xED 0x67])

  (opcode  BIT (a b)
    (N R8)  [0xCB (=i a (+ 0b0100_0000 (R8 b)) 0x07 3)]
    (N HL$) [0xCB (=i a 0b0100_0110 0x07 3)]
    (N IX$) [0xDD 0xCB (=l b) (=i a 0b0100_0110 0x07 3)]
    (N IY$) [0xFD 0xCB (=l b) (=i a 0b0100_0110 0x07 3)])

  (opcode  SET (a b)
    (N R8)  [0xCB (=i a (+ 0b1100_0000 (R8 b)) 0x07 3)]
    (N HL$) [0xCB (=i a 0b1100_0110 0x07 3)]
    (N IX$) [0xDD 0xCB (=l b) (=i a 0b1100_0110 0x07 3)]
    (N IY$) [0xFD 0xCB (=l b) (=i a 0b1100_0110 0x07 3)])

  (opcode  RES (a b)
    (N R8)  [0xCB (=i a (+ 0b1000_0000 (R8 b)) 0x07 3)]
    (N HL$) [0xCB (=i a 0b1000_0110 0x07 3)]
    (N IX$) [0xDD 0xCB (=l b) (=i a 0b1000_0110 0x07 3)]
    (N IY$) [0xFD 0xCB (=l b) (=i a 0b1000_0110 0x07 3)])

  (opcode  JP (a)
    (NN)  [0xC3 (=l a) (=h a)]
    (HL$) [0xE9]
    (IX$) [0xDD (=i a 0xE9 0 0)]
    (IY$) [0xFD (=i a 0xE9 0 0)])
  (opcode  JP (a b) (CC NN) [(+ 0b1100_0010 (CC a 3)) (=l b) (=h b)])
  (example JP
    (IX$) "JP [IX]" "JP (IX)"
    (IY$) "JP [IY]" "JP (IY)")

  (opcode  JR (a) (NN) [0x18 (=rl a -2)])
  (opcode  JR (a b)
    (C?  NN) [0x38 (=rl b -2)]
    (NC? NN) [0x30 (=rl b -2)]
    (Z?  NN) [0x28 (=rl b -2)]
    (NZ? NN) [0x20 (=rl b -2)])

  (opcode  DJNZ (a) (NN) [0x10 (=rl a -2)])

  (opcode  CALL (a)   (NN)  [0xCD (=l a) (=h a)])
  (opcode  CALL (a b) (CC NN) [(+ 0b1100_0100 (CC a 3)) (=l b) (=h b)])

  (opcode  RET  ()  ()   [0xC9])
  (opcode  RET  (a) (CC) [(+ 0b1100_0000 (CC a 3))])
  (opcode  RETI ()  ()   [0xED 0x4D])
  (opcode  RETN ()  ()   [0xED 0x45])

  (opcode  RST  (a) (N)  [(=i a 0b1100_0111 0b0011_1000 0)])
  (example RST  (N) "RST 16" "RST 16")

  (opcode  #.jump (a)   (NN)    [0xC3 (=l a) (=h a)])
  (opcode  #.jump (a b) (NN CC) [(+ 0b1100_0010 (CC b 3)) (=l a) (=h a)])
  (example #.jump (*) "" "")

  (opcode  #.call (a)   (NN)    [0xCD (=l a) (=h a)])
  (opcode  #.call (a b) (NN CC) [(+ 0b1100_0100 (CC b 3)) (=l a) (=h a)])
  (example #.call (*) "" "")

  (opcode  #.return ()  ()   [0xC9])
  (opcode  #.return (a) (CC) [(+ 0b1100_0000 (CC a 3))])
  (example #.return (*) "" "")

  (opcode  IN (a b)
    (A N$)   [0xDB (=l b)]
    (R8 C$)  [0xED (+ 0b0100_0000 (R8 a 3))])
  (opcode  INI  () () [0xED 0xA2])
  (opcode  INIR () () [0xED 0xB2])
  (opcode  IND  () () [0xED 0xAA])
  (opcode  INDR () () [0xED 0xBA])

  (opcode  OUT (a b)
    (N$ A)   [0xD3 (=l a)]
    (C$ R8)  [0xED (+ 0b0100_0001 (R8 b 3))])
  (opcode  OUTI () () [0xED 0xA3])
  (opcode  OTIR () () [0xED 0xB3])
  (opcode  OUTD () () [0xED 0xAB])
  (opcode  OTDR () () [0xED 0xBB])

  (opcode  DAA  () () [0x27])
  (opcode  CPL  () () [0x2F])
  (opcode  NEG  () () [0xED 0x44])
  (opcode  CCF  () () [0x3F])
  (opcode  SCF  () () [0x37])
  (opcode  NOP  () () [0x00])
  (opcode  HALT () () [0x76])
  (opcode  DI   () () [0xF3])
  (opcode  EI   () () [0xFB])

  (bytemap IM 0b0011 0 2 0x46 0x56 0x5E 0)
  (opcode  IM (a) (N) [0xED (=m a IM)])
  (example IM (N) "IM 1" "IM 1")

  (operator <- (a b)
    (R8 _)  [(LD (= a) (= b))]
    (I _)   [(LD (= a) (= b))]
    (R _)   [(LD (= a) (= b))]
    (DD _)  [(LD (= a) (= b))]
    (IX _)  [(LD (= a) (= b))]
    (IY _)  [(LD (= a) (= b))]
    (BC BC) [(LD B B) (LD C C)]
    (BC DE) [(LD B D) (LD C E)]
    (BC HL) [(LD B H) (LD C L)]
    (DE BC) [(LD D B) (LD E C)]
    (DE DE) [(LD D D) (LD E E)]
    (DE HL) [(LD D H) (LD E L)]
    (HL BC) [(LD H B) (LD L C)]
    (HL DE) [(LD H D) (LD L E)]
    (HL HL) [(LD H H) (LD L L)]
    (DD PQ) [(#.LDP (= a) (= b))])
  (operator -> (a b)
    (R8 _)  [(LD (= b) (= a))]
    (I _)   [(LD (= b) (= a))]
    (R _)   [(LD (= b) (= a))]
    (DD _)  [(LD (= b) (= a))]
    (IX _)  [(LD (= b) (= a))]
    (IY _)  [(LD (= b) (= a))]
    (NN _)  [(LD (= b) (= a))]
    (BC BC) [(LD B B) (LD C C)]
    (BC DE) [(LD D B) (LD E C)]
    (BC HL) [(LD H B) (LD L C)]
    (DE BC) [(LD B D) (LD C E)]
    (DE DE) [(LD D D) (LD E E)]
    (DE HL) [(LD H D) (LD L E)]
    (HL BC) [(LD B H) (LD C L)]
    (HL DE) [(LD D H) (LD E L)]
    (HL HL) [(LD H H) (LD L L)]
    (DD PQ) [(#.LDP (= b) (= a))])
  (operator <-> (a b)
    (AF AF-) [(EX AF AF-)]
    (DE HL)  [(EX DE HL)]
    (HL DE)  [(EX DE HL)]
    (HL SP$) [(EX SP$ HL)]
    (SP$ HL) [(EX SP$ HL)])

  (operator -push (a) (_) [(PUSH (= a))])
  (operator -pop  (a) (_) [(POP (= a))])

  (operator ++ (a) (_) [(INC (= a))])
  (operator -- (a) (_) [(DEC (= a))])
  (operator -not (a) (A) [(CPL)])
  (operator -neg (a) (A) [(NEG)])

  (operator +  (a b) (_ _) [(ADD (= a) (= b))])
  (operator +$ (a b) (_ _) [(ADC (= a) (= b))])

  (operator - (a b)
    (HL _) [(OR A) (SBC HL (= b))]
    (A _)  [(SUB (= b))])
  (operator -$ (a b) (_ _) [(SBC (= a) (= b))])
  (operator -? (a b) (A _) [(CP (= b))])

  (operator &   (a b) (A _) [(AND (= b))])
  (operator "|" (a b) (A _) [(OR (= b))])
  (operator ^   (a b) (A _) [(XOR (= b))])

  (operator <* (a b)
    (A NN) [(#.REP (= b) `[(RLCA)])]
    (_ NN) [(#.REP (= b) `[(RLC (= a))])])
  (operator <*$ (a b)
    (A NN) [(#.REP (= b) `[(RLA)])]
    (_ NN) [(#.REP (= b) `[(RL (= a))])])
  (operator >* (a b)
    (A NN) [(#.REP (= b) `[(RRCA)])]
    (_ NN) [(#.REP (= b) `[(RRC (= a))])])
  (operator >*$ (a b)
    (A NN) [(#.REP (= b) `[(RRA)])]
    (_ NN) [(#.REP (= b) `[(RR (= a))])])

  (operator <<  (a b) (_ NN) [(#.REP (= b) `[(SLA (= a))])])
  (operator >>  (a b) (_ NN) [(#.REP (= b) `[(SRA (= a))])])
  (operator >>> (a b) (_ NN) [(#.REP (= b) `[(SRL (= a))])])

  (operator -set   (a b) (_ _) [(SET (= b) (= a))])
  (operator -reset (a b) (_ _) [(RES (= b) (= a))])
  (operator -bit?  (a b) (_ _) [(BIT (= b) (= a))])

  (operator -in (a b)
    (A C)  [(IN A C$)]
    (_ C)  [(IN (= a) C$)]
    (A NN) [(IN A (= b NN$))])
  (operator -out (a b)
    (A C)  [(OUT C$ A)]
    (_ C)  [(OUT C$ (= a))]
    (A NN) [(OUT (= b NN$) A)])

  (operator -zero? (a)
    (A)  [(AND (= a))]
    (DD) [(#.INVALID (= a))]
    (IX) [(#.INVALID (= a))]
    (IY) [(#.INVALID (= a))]
    (_)  [(INC (= a)) (DEC (= a))])

  (operator -jump (a) (NN) [(#.jump (= a))])
  (example -jump (*) "" "")

  (operator -jump-if (a b) (NN CC) [(#.jump (= a) (= b))])
  (example -jump-if (*) "" "")

  (operator -jump-unless (a b)
    (NN NZ?) [(#.jump (= a) Z? )]
    (NN Z?)  [(#.jump (= a) NZ?)]
    (NN NC?) [(#.jump (= a) C? )]
    (NN C?)  [(#.jump (= a) NC?)]
    (NN PO?) [(#.jump (= a) PE?)]
    (NN PE?) [(#.jump (= a) PO?)]
    (NN M?)  [(#.jump (= a) P? )]
    (NN P?)  [(#.jump (= a) M? )])
  (example -jump-unless (*) "" "")

  (operator -return (a) (PC) [(#.return)])
  (example -return (*) "" "")

  (operator -return-if (a b) (PC CC) [(#.return (= b))])
  (example -return-if (*) "" "")

  (operator -return-unless (a b)
    (PC NZ?) [(#.return Z? )]
    (PC Z?)  [(#.return NZ?)]
    (PC NC?) [(#.return C? )]
    (PC C?)  [(#.return NC?)]
    (PC PO?) [(#.return PE?)]
    (PC PE?) [(#.return PO?)]
    (PC M?)  [(#.return P? )]
    (PC P?)  [(#.return M? )])
  (example -return-unless (*) "" "")

  (example $operators.load/store (*)
    "BC <- D : E" "LD B, D; LD C, E"
    "BC -> D : E" "LD D, B; LD E, C"

    "BC <- BC" "LD B, B; LD C, C"
    "BC <- DE" "LD B, D; LD C, E"
    "BC <- HL" "LD B, H; LD C, L"
    "DE <- BC" "LD D, B; LD E, C"
    "DE <- DE" "LD D, D; LD E, E"
    "DE <- HL" "LD D, H; LD E, L"
    "HL <- BC" "LD H, B; LD L, C"
    "HL <- DE" "LD H, D; LD L, E"
    "HL <- HL" "LD H, H; LD L, L"

    "BC -> BC" "LD B, B; LD C, C"
    "BC -> DE" "LD D, B; LD E, C"
    "BC -> HL" "LD H, B; LD L, C"
    "DE -> BC" "LD B, D; LD C, E"
    "DE -> DE" "LD D, D; LD E, E"
    "DE -> HL" "LD H, D; LD L, E"
    "HL -> BC" "LD B, H; LD C, L"
    "HL -> DE" "LD D, H; LD E, L"
    "HL -> HL" "LD H, H; LD L, L")

  (example $operators.rotate/shift (*)
    "A <* 2"  "RLCA; RLCA"
    "B <* 2"  "RLC B; RLC B"
    "A <*$ 2" "RLA; RLA"
    "B <*$ 2" "RL B; RL B"
    "A >* 2"  "RRCA; RRCA"
    "B >* 2"  "RRC B; RRC B"
    "A >*$ 2" "RRA; RRA"
    "B >*$ 2" "RR B; RR B"
    "A << 2"  "SLA A; SLA A"
    "B << 2"  "SLA B; SLA B"
    "A >> 2"  "SRA A; SRA A"
    "B >> 2"  "SRA B; SRA B"
    "A >>> 2" "SRL A; SRL A"
    "B >>> 2" "SRL B; SRL B")

  (example $operators.misc (*)
    "LBI:; $(LBI) -jump-if NZ?" "LBI:; JP NZ, LBI"
    "$(LBI) -jump-if Z?"  "JP Z, LBI"
    "$(LBI) -jump-if NC?" "JP NC, LBI"
    "$(LBI) -jump-if C?"  "JP C, LBI"
    "$(LBI) -jump-if PO?" "JP PO, LBI"
    "$(LBI) -jump-if PE?" "JP PE, LBI"
    "$(LBI) -jump-if M?"  "JP M, LBI"
    "$(LBI) -jump-if P?"  "JP P, LBI"

    "LBU:; $(LBU) -jump-unless NZ?" "LBU:; JP Z, LBU"
    "$(LBU) -jump-unless Z?"  "JP NZ, LBU"
    "$(LBU) -jump-unless NC?" "JP C, LBU"
    "$(LBU) -jump-unless C?"  "JP NC, LBU"
    "$(LBU) -jump-unless PO?" "JP PE, LBU"
    "$(LBU) -jump-unless PE?" "JP PO, LBU"
    "$(LBU) -jump-unless M?"  "JP P, LBU"
    "$(LBU) -jump-unless P?"  "JP M, LBU"

    "PC -return" "RET"
    "PC -return-if NZ?" "RET NZ"
    "PC -return-if Z?"  "RET Z"
    "PC -return-if NC?" "RET NC"
    "PC -return-if C?"  "RET C"
    "PC -return-if PO?" "RET PO"
    "PC -return-if PE?" "RET PE"
    "PC -return-if M?"  "RET M"
    "PC -return-if P?"  "RET P"
    "PC -return-unless NZ?" "RET Z"
    "PC -return-unless Z?"  "RET NZ"
    "PC -return-unless NC?" "RET C"
    "PC -return-unless C?"  "RET NC"
    "PC -return-unless PO?" "RET PE"
    "PC -return-unless PE?" "RET PO"
    "PC -return-unless M?"  "RET P"
    "PC -return-unless P?"  "RET M"

    "proc f(!){ RET }" "f: RET"
    "f(!)" "CALL f"
    "NC?.f(!)" "CALL NC, f"))

(arch (z80 +undocumented)
  (operand IXH RegIXH "IXH" "IXH")
  (operand IXL RegIXL "IXL" "IXL")
  (operand IYH RegIYH "IYH" "IYH")
  (operand IYL RegIYL "IYL" "IYL")

  (registers IXH IXL IYH IYL)

  (map X8 A 7 B 0 C 1 D 2 E 3 IXH 4 IXL 5)
  (map Y8 A 7 B 0 C 1 D 2 E 3 IYH 4 IYL 5)
  (map A-E A 7 B 0 C 1 D 2 E 3)

  (example "$prologue" (*) "arch z80 +undocumented; flat!; optimize near-jump 0" "")

  (opcode  LD (a b)
    (X8  IXH) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (IXH A-E) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (X8  IXL) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (IXL A-E) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (Y8  IYH) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]
    (IYH A-E) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]
    (Y8  IYL) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]
    (IYL A-E) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]

    (IXH N)   [0xDD (+ 0b0000_0110 (X8 a 3)) (=l b)]
    (IXL N)   [0xDD (+ 0b0000_0110 (X8 a 3)) (=l b)]
    (IYH N)   [0xFD (+ 0b0000_0110 (Y8 a 3)) (=l b)]
    (IYL N)   [0xFD (+ 0b0000_0110 (Y8 a 3)) (=l b)])

  (opcode  ADD (a b)
    (A IXH)  [0xDD (+ 0b1000_0000 (X8 b))]
    (A IXL)  [0xDD (+ 0b1000_0000 (X8 b))]
    (A IYH)  [0xFD (+ 0b1000_0000 (Y8 b))]
    (A IYL)  [0xFD (+ 0b1000_0000 (Y8 b))])

  (opcode  ADC (a b)
    (A IXH)  [0xDD (+ 0b1000_1000 (X8 b))]
    (A IXL)  [0xDD (+ 0b1000_1000 (X8 b))]
    (A IYH)  [0xFD (+ 0b1000_1000 (Y8 b))]
    (A IYL)  [0xFD (+ 0b1000_1000 (Y8 b))])

  (opcode  SUB (a)
    (IXH)  [0xDD (+ 0b1001_0000 (X8 a))]
    (IXL)  [0xDD (+ 0b1001_0000 (X8 a))]
    (IYH)  [0xFD (+ 0b1001_0000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1001_0000 (Y8 a))])

  (opcode  SBC (a b)
    (A IXH)  [0xDD (+ 0b1001_1000 (X8 b))]
    (A IXL)  [0xDD (+ 0b1001_1000 (X8 b))]
    (A IYH)  [0xFD (+ 0b1001_1000 (Y8 b))]
    (A IYL)  [0xFD (+ 0b1001_1000 (Y8 b))])

  (opcode  AND (a)
    (IXH)  [0xDD (+ 0b1010_0000 (X8 a))]
    (IXL)  [0xDD (+ 0b1010_0000 (X8 a))]
    (IYH)  [0xFD (+ 0b1010_0000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1010_0000 (Y8 a))])

  (opcode  OR (a)
    (IXH)  [0xDD (+ 0b1011_0000 (X8 a))]
    (IXL)  [0xDD (+ 0b1011_0000 (X8 a))]
    (IYH)  [0xFD (+ 0b1011_0000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1011_0000 (Y8 a))])

  (opcode  XOR (a)
    (IXH)  [0xDD (+ 0b1010_1000 (X8 a))]
    (IXL)  [0xDD (+ 0b1010_1000 (X8 a))]
    (IYH)  [0xFD (+ 0b1010_1000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1010_1000 (Y8 a))])

  (opcode  CP (a)
    (IXH)  [0xDD (+ 0b1011_1000 (X8 a))]
    (IXL)  [0xDD (+ 0b1011_1000 (X8 a))]
    (IYH)  [0xFD (+ 0b1011_1000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1011_1000 (Y8 a))])

  (opcode  INC (a)
    (IXH)  [0xDD (+ 0b0000_0100 (X8 a 3))]
    (IXL)  [0xDD (+ 0b0000_0100 (X8 a 3))]
    (IYH)  [0xFD (+ 0b0000_0100 (Y8 a 3))]
    (IYL)  [0xFD (+ 0b0000_0100 (Y8 a 3))])
  (example INC
    (IXH) _ "DB 0xDD, 0x24"
    (IXL) _ "DB 0xDD, 0x2C"
    (IYH) _ "DB 0xFD, 0x24"
    (IYL) _ "DB 0xFD, 0x2C")

  (opcode  DEC (a)
    (IXH)  [0xDD (+ 0b0000_0101 (X8 a 3))]
    (IXL)  [0xDD (+ 0b0000_0101 (X8 a 3))]
    (IYH)  [0xFD (+ 0b0000_0101 (Y8 a 3))]
    (IYL)  [0xFD (+ 0b0000_0101 (Y8 a 3))])
  (example DEC
    (IXH) _ "DB 0xDD, 0x25"
    (IXL) _ "DB 0xDD, 0x2D"
    (IYH) _ "DB 0xFD, 0x25"
    (IYL) _ "DB 0xFD, 0x2D")

  (opcode  RLC (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0000_0000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0000_0000 (R8 a))])
  (example RLC
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x07"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x00"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x01"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x02"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x03"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x04"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x05"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x07"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x00"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x01"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x02"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x03"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x04"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x05")

  (opcode  RL (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0001_0000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0001_0000 (R8 a))])
  (example RL
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x17"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x10"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x11"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x12"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x13"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x14"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x15"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x17"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x10"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x11"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x12"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x13"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x14"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x15")

  (opcode  RRC (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0000_1000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0000_1000 (R8 a))])
  (example RRC
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x0F"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x08"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x09"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x0A"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x0B"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x0C"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x0D"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x0F"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x08"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x09"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x0A"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x0B"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x0C"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x0D")

  (opcode  RR (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0001_1000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0001_1000 (R8 a))])
  (example RR
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x1F"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x18"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x19"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x1A"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x1B"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x1C"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x1D"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x1F"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x18"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x19"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x1A"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x1B"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x1C"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x1D")

  (opcode  SLA (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0010_0000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0010_0000 (R8 a))])
  (example SLA
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x27"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x20"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x21"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x22"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x23"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x24"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x25"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x27"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x20"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x21"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x22"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x23"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x24"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x25")

  (opcode SLL (a)
    (R8)  [0xCB (+ 0b0011_0000 (R8 a))]
    (HL$) [0xCB 0x36]
    (IX$) [0xDD 0xCB (=l a) 0x36]
    (IY$) [0xFD 0xCB (=l a) 0x36])
  (example SLL
    (A)   _ "DB 0xCB, 0x37"
    (B)   _ "DB 0xCB, 0x30"
    (C)   _ "DB 0xCB, 0x31"
    (D)   _ "DB 0xCB, 0x32"
    (E)   _ "DB 0xCB, 0x33"
    (H)   _ "DB 0xCB, 0x34"
    (L)   _ "DB 0xCB, 0x35"
    (HL$) _ "DB 0xCB, 0x36"
    (IX$) _ "DB 0xDD, 0xCB, 0x05, 0x36"
    (IY$) _ "DB 0xFD, 0xCB, 0x05, 0x36")

  (opcode SLL (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0011_0000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0011_0000 (R8 a))])
  (example SLL
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x37"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x30"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x31"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x32"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x33"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x34"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x35"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x37"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x30"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x31"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x32"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x33"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x34"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x35")

  (opcode  SRA (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0010_1000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0010_1000 (R8 a))])
  (example SRA
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x2F"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x28"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x29"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x2A"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x2B"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x2C"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x2D"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x2F"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x28"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x29"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x2A"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x2B"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x2C"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x2D")

  (opcode  SRL (a b)
    (R8 IX$) [0xDD 0xCB (=l b) (+ 0b0011_1000 (R8 a))]
    (R8 IY$) [0xFD 0xCB (=l b) (+ 0b0011_1000 (R8 a))])
  (example SRL
    (A IX$) _ "DB 0xDD, 0xCB, 0x05, 0x3F"
    (B IX$) _ "DB 0xDD, 0xCB, 0x05, 0x38"
    (C IX$) _ "DB 0xDD, 0xCB, 0x05, 0x39"
    (D IX$) _ "DB 0xDD, 0xCB, 0x05, 0x3A"
    (E IX$) _ "DB 0xDD, 0xCB, 0x05, 0x3B"
    (H IX$) _ "DB 0xDD, 0xCB, 0x05, 0x3C"
    (L IX$) _ "DB 0xDD, 0xCB, 0x05, 0x3D"
    (A IY$) _ "DB 0xFD, 0xCB, 0x05, 0x3F"
    (B IY$) _ "DB 0xFD, 0xCB, 0x05, 0x38"
    (C IY$) _ "DB 0xFD, 0xCB, 0x05, 0x39"
    (D IY$) _ "DB 0xFD, 0xCB, 0x05, 0x3A"
    (E IY$) _ "DB 0xFD, 0xCB, 0x05, 0x3B"
    (H IY$) _ "DB 0xFD, 0xCB, 0x05, 0x3C"
    (L IY$) _ "DB 0xFD, 0xCB, 0x05, 0x3D")

  (opcode  BIT (a b c)
    (R8 N IX$) [0xDD 0xCB (=l c) (=i b (+ 0b0100_0000 (R8 a)) 0x07 3)]
    (R8 N IY$) [0xFD 0xCB (=l c) (=i b (+ 0b0100_0000 (R8 a)) 0x07 3)]
    (N IX$ R8) [0xDD 0xCB (=l b) (=i a (+ 0b0100_0000 (R8 c)) 0x07 3)]
    (N IY$ R8) [0xFD 0xCB (=l b) (=i a (+ 0b0100_0000 (R8 c)) 0x07 3)])
  (example BIT
    (A N IX$) _ "DB 0xDD, 0xCB, 0x05, 0x6F"
    (B N IX$) _ "DB 0xDD, 0xCB, 0x05, 0x68"
    (C N IX$) _ "DB 0xDD, 0xCB, 0x05, 0x69"
    (D N IX$) _ "DB 0xDD, 0xCB, 0x05, 0x6A"
    (E N IX$) _ "DB 0xDD, 0xCB, 0x05, 0x6B"
    (H N IX$) _ "DB 0xDD, 0xCB, 0x05, 0x6C"
    (L N IX$) _ "DB 0xDD, 0xCB, 0x05, 0x6D"

    (A N IY$) _ "DB 0xFD, 0xCB, 0x05, 0x6F"
    (B N IY$) _ "DB 0xFD, 0xCB, 0x05, 0x68"
    (C N IY$) _ "DB 0xFD, 0xCB, 0x05, 0x69"
    (D N IY$) _ "DB 0xFD, 0xCB, 0x05, 0x6A"
    (E N IY$) _ "DB 0xFD, 0xCB, 0x05, 0x6B"
    (H N IY$) _ "DB 0xFD, 0xCB, 0x05, 0x6C"
    (L N IY$) _ "DB 0xFD, 0xCB, 0x05, 0x6D"

    (N IX$ A) _ "DB 0xDD, 0xCB, 0x05, 0x6F"
    (N IX$ B) _ "DB 0xDD, 0xCB, 0x05, 0x68"
    (N IX$ C) _ "DB 0xDD, 0xCB, 0x05, 0x69"
    (N IX$ D) _ "DB 0xDD, 0xCB, 0x05, 0x6A"
    (N IX$ E) _ "DB 0xDD, 0xCB, 0x05, 0x6B"
    (N IX$ H) _ "DB 0xDD, 0xCB, 0x05, 0x6C"
    (N IX$ L) _ "DB 0xDD, 0xCB, 0x05, 0x6D"

    (N IY$ A) _ "DB 0xFD, 0xCB, 0x05, 0x6F"
    (N IY$ B) _ "DB 0xFD, 0xCB, 0x05, 0x68"
    (N IY$ C) _ "DB 0xFD, 0xCB, 0x05, 0x69"
    (N IY$ D) _ "DB 0xFD, 0xCB, 0x05, 0x6A"
    (N IY$ E) _ "DB 0xFD, 0xCB, 0x05, 0x6B"
    (N IY$ H) _ "DB 0xFD, 0xCB, 0x05, 0x6C"
    (N IY$ L) _ "DB 0xFD, 0xCB, 0x05, 0x6D")

  (opcode  SET (a b c)
    (R8 N IX$) [0xDD 0xCB (=l c) (=i b (+ 0b1100_0000 (R8 a)) 0x07 3)]
    (R8 N IY$) [0xFD 0xCB (=l c) (=i b (+ 0b1100_0000 (R8 a)) 0x07 3)]
    (N IX$ R8) [0xDD 0xCB (=l b) (=i a (+ 0b1100_0000 (R8 c)) 0x07 3)]
    (N IY$ R8) [0xFD 0xCB (=l b) (=i a (+ 0b1100_0000 (R8 c)) 0x07 3)])
  (example SET
    (A N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xEF"
    (B N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xE8"
    (C N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xE9"
    (D N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xEA"
    (E N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xEB"
    (H N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xEC"
    (L N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xED"

    (A N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xEF"
    (B N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xE8"
    (C N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xE9"
    (D N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xEA"
    (E N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xEB"
    (H N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xEC"
    (L N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xED"

    (N IX$ A) _ "DB 0xDD, 0xCB, 0x05, 0xEF"
    (N IX$ B) _ "DB 0xDD, 0xCB, 0x05, 0xE8"
    (N IX$ C) _ "DB 0xDD, 0xCB, 0x05, 0xE9"
    (N IX$ D) _ "DB 0xDD, 0xCB, 0x05, 0xEA"
    (N IX$ E) _ "DB 0xDD, 0xCB, 0x05, 0xEB"
    (N IX$ H) _ "DB 0xDD, 0xCB, 0x05, 0xEC"
    (N IX$ L) _ "DB 0xDD, 0xCB, 0x05, 0xED"

    (N IY$ A) _ "DB 0xFD, 0xCB, 0x05, 0xEF"
    (N IY$ B) _ "DB 0xFD, 0xCB, 0x05, 0xE8"
    (N IY$ C) _ "DB 0xFD, 0xCB, 0x05, 0xE9"
    (N IY$ D) _ "DB 0xFD, 0xCB, 0x05, 0xEA"
    (N IY$ E) _ "DB 0xFD, 0xCB, 0x05, 0xEB"
    (N IY$ H) _ "DB 0xFD, 0xCB, 0x05, 0xEC"
    (N IY$ L) _ "DB 0xFD, 0xCB, 0x05, 0xED")

  (opcode  RES (a b c)
    (R8 N IX$) [0xDD 0xCB (=l c) (=i b (+ 0b1000_0000 (R8 a)) 0x07 3)]
    (R8 N IY$) [0xFD 0xCB (=l c) (=i b (+ 0b1000_0000 (R8 a)) 0x07 3)]
    (N IX$ R8) [0xDD 0xCB (=l b) (=i a (+ 0b1000_0000 (R8 c)) 0x07 3)]
    (N IY$ R8) [0xFD 0xCB (=l b) (=i a (+ 0b1000_0000 (R8 c)) 0x07 3)])
  (example RES
    (A N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xAF"
    (B N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xA8"
    (C N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xA9"
    (D N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xAA"
    (E N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xAB"
    (H N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xAC"
    (L N IX$) _ "DB 0xDD, 0xCB, 0x05, 0xAD"

    (A N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xAF"
    (B N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xA8"
    (C N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xA9"
    (D N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xAA"
    (E N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xAB"
    (H N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xAC"
    (L N IY$) _ "DB 0xFD, 0xCB, 0x05, 0xAD"

    (N IX$ A) _ "DB 0xDD, 0xCB, 0x05, 0xAF"
    (N IX$ B) _ "DB 0xDD, 0xCB, 0x05, 0xA8"
    (N IX$ C) _ "DB 0xDD, 0xCB, 0x05, 0xA9"
    (N IX$ D) _ "DB 0xDD, 0xCB, 0x05, 0xAA"
    (N IX$ E) _ "DB 0xDD, 0xCB, 0x05, 0xAB"
    (N IX$ H) _ "DB 0xDD, 0xCB, 0x05, 0xAC"
    (N IX$ L) _ "DB 0xDD, 0xCB, 0x05, 0xAD"

    (N IY$ A) _ "DB 0xFD, 0xCB, 0x05, 0xAF"
    (N IY$ B) _ "DB 0xFD, 0xCB, 0x05, 0xA8"
    (N IY$ C) _ "DB 0xFD, 0xCB, 0x05, 0xA9"
    (N IY$ D) _ "DB 0xFD, 0xCB, 0x05, 0xAA"
    (N IY$ E) _ "DB 0xFD, 0xCB, 0x05, 0xAB"
    (N IY$ H) _ "DB 0xFD, 0xCB, 0x05, 0xAC"
    (N IY$ L) _ "DB 0xFD, 0xCB, 0x05, 0xAD")

  (opcode  IN (a)   (C$)    [0xED 0x70])
  (example IN (C$) _ "DB 0xED, 0x70")

  (opcode  IN (a b) (F  C$) [0xED 0x70])
  (example IN (F C$) _ "DB 0xED, 0x70")

  (opcode  OUT (a b) (C$ N)  [0xED (=i b 0x71 0 0)])
  (example OUT (C$ N) "OUT [C] 0" "DB 0xED, 0x71")

  (operator <- (a b)
    (IXH _)  [(LD (= a) (= b))]
    (IXL _)  [(LD (= a) (= b))]
    (IYH _)  [(LD (= a) (= b))]
    (IYL _)  [(LD (= a) (= b))])
  (operator -> (a b)
    (IXH _)  [(LD (= b) (= a))]
    (IXL _)  [(LD (= b) (= a))]
    (IYH _)  [(LD (= b) (= a))]
    (IYL _)  [(LD (= b) (= a))])

  (operator <<< (a b) (_ NN) [(#.REP (= b) `[(SLL (= a))])])

  (example -in  (F C) "F -in C" "DB 0xED, 0x70")
  (example -out (N C) "$(0) -out C" "DB 0xED, 0x71"))

(arch (z80 +compat8080)
  (example "$prologue" (*) "arch z80 +compat8080; flat!" "")

  (opcode  LD (a b)
    (R8 IX$) [(=U)]
    (R8 IY$) [(=U)]
    (IX$ R8) [(=U)]
    (IY$ R8) [(=U)]
    (IX$ N)  [(=U)]
    (IY$ N)  [(=U)]
    (A I)    [(=U)]
    (A R)    [(=U)]
    (I A)    [(=U)]
    (R A)    [(=U)]
    (IX NN)  [(=U)]
    (IY NN)  [(=U)]
    (DD NN$) [(=U)]
    (IX NN$) [(=U)]
    (IY NN$) [(=U)]
    (NN$ DD) [(=U)]
    (NN$ IX) [(=U)]
    (NN$ IY) [(=U)]
    (SP IX)  [(=U)]
    (SP IY)  [(=U)])

  (opcode  PUSH (a)
    (IX) [(=U)]
    (IY) [(=U)])

  (opcode  POP (a)
    (IX) [(=U)]
    (IY) [(=U)])

  (opcode  EX (a b)
    (AF  AF-) [(=U)]
    (SP$ IX)  [(=U)]
    (SP$ IY)  [(=U)])
  (opcode  EXX  () () [(=U)])

  (opcode  LDI  () () [(=U)])
  (opcode  LDIR () () [(=U)])
  (opcode  LDD  () () [(=U)])
  (opcode  LDDR () () [(=U)])

  (opcode  CPI  () () [(=U)])
  (opcode  CPIR () () [(=U)])
  (opcode  CPD  () () [(=U)])
  (opcode  CPDR () () [(=U)])

  (opcode  ADD (a b)
    (A IX$) [(=U)]
    (A IY$) [(=U)]
    (IX PP) [(=U)]
    (IY RR) [(=U)])

  (opcode  ADC (a b)
    (A IX$) [(=U)]
    (A IY$) [(=U)]
    (HL DD) [(=U)])

  (opcode  SUB (a)
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  SBC (a b)
    (A IX$) [(=U)]
    (A IY$) [(=U)]
    (HL DD) [(=U)])

  (opcode  AND (a)
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  OR (a)
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  XOR (a)
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  CP (a)
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  INC (a)
    (IX$) [(=U)]
    (IY$) [(=U)]
    (IX)  [(=U)]
    (IY)  [(=U)])

  (opcode  DEC (a)
    (IX$) [(=U)]
    (IY$) [(=U)]
    (IX)  [(=U)]
    (IY)  [(=U)])

  (opcode  RLC (a)
    (R8)  [(=U)]
    (HL$) [(=U)]
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  RL (a)
    (R8)  [(=U)]
    (HL$) [(=U)]
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  RRC (a)
    (R8)  [(=U)]
    (HL$) [(=U)]
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  RR (a)
    (R8)  [(=U)]
    (HL$) [(=U)]
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  SLA (a)
    (R8)  [(=U)]
    (HL$) [(=U)]
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  SRA (a)
    (R8)  [(=U)]
    (HL$) [(=U)]
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  SRL (a)
    (R8)  [(=U)]
    (HL$) [(=U)]
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  RLD  () () [(=U)])
  (opcode  RRD  () () [(=U)])

  (opcode  BIT (a b)
    (N R8)  [(=U)]
    (N HL$) [(=U)]
    (N IX$) [(=U)]
    (N IY$) [(=U)])

  (opcode  SET (a b)
    (N R8)  [(=U)]
    (N HL$) [(=U)]
    (N IX$) [(=U)]
    (N IY$) [(=U)])

  (opcode  RES (a b)
    (N R8)  [(=U)]
    (N HL$) [(=U)]
    (N IX$) [(=U)]
    (N IY$) [(=U)])

  (opcode  JP (a)
    (IX$) [(=U)]
    (IY$) [(=U)])

  (opcode  JR (a) (NN) [(=U)])
  (opcode  JR (a b)
    (C?  NN) [(=U)]
    (NC? NN) [(=U)]
    (Z?  NN) [(=U)]
    (NZ? NN) [(=U)])

  (opcode  DJNZ (a) (NN) [(=U)])

  (opcode  RETI ()  ()   [(=U)])
  (opcode  RETN ()  ()   [(=U)])

  (opcode  IN (a b) (R8 C$)  [(=U)])
  (opcode  INI  () () [(=U)])
  (opcode  INIR () () [(=U)])
  (opcode  IND  () () [(=U)])
  (opcode  INDR () () [(=U)])

  (opcode  OUT (a b) (C$ R8)  [(=U)])
  (opcode  OUTI () () [(=U)])
  (opcode  OTIR () () [(=U)])
  (opcode  OUTD () () [(=U)])
  (opcode  OTDR () () [(=U)])

  (opcode  NEG  () () [(=U)])

  (opcode  IM (a) (N) [(=U)])

  (example $operators.rotate/shift (*)
    "A <* 2"  "RLCA; RLCA"
    "A <*$ 2" "RLA; RLA"
    "A >* 2"  "RRCA; RRCA"
    "A >*$ 2" "RRA; RRA"))

(arch (z80 +r800)
  (operand IXH RegIXH "IXH" "IXH")
  (operand IXL RegIXL "IXL" "IXL")
  (operand IYH RegIYH "IYH" "IYH")
  (operand IYL RegIYL "IYL" "IYL")

  (registers IXH IXL IYH IYL)

  (map X8 A 7 B 0 C 1 D 2 E 3 IXH 4 IXL 5)
  (map Y8 A 7 B 0 C 1 D 2 E 3 IYH 4 IYL 5)
  (map A-E A 7 B 0 C 1 D 2 E 3)

  (example "$prologue" (*) "arch z80 +r800; flat!; optimize near-jump 0" "")

  (opcode  LD (a b)
    (X8  IXH) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (IXH A-E) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (X8  IXL) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (IXL A-E) [0xDD (+ 0b0100_0000 (X8 a 3) (X8 b))]
    (Y8  IYH) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]
    (IYH A-E) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]
    (Y8  IYL) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]
    (IYL A-E) [0xFD (+ 0b0100_0000 (Y8 a 3) (Y8 b))]

    (IXH N)   [0xDD (+ 0b0000_0110 (X8 a 3)) (=l b)]
    (IXL N)   [0xDD (+ 0b0000_0110 (X8 a 3)) (=l b)]
    (IYH N)   [0xFD (+ 0b0000_0110 (Y8 a 3)) (=l b)]
    (IYL N)   [0xFD (+ 0b0000_0110 (Y8 a 3)) (=l b)])

  (opcode  ADD (a b)
    (A IXH)  [0xDD (+ 0b1000_0000 (X8 b))]
    (A IXL)  [0xDD (+ 0b1000_0000 (X8 b))]
    (A IYH)  [0xFD (+ 0b1000_0000 (Y8 b))]
    (A IYL)  [0xFD (+ 0b1000_0000 (Y8 b))])

  (opcode  ADC (a b)
    (A IXH)  [0xDD (+ 0b1000_1000 (X8 b))]
    (A IXL)  [0xDD (+ 0b1000_1000 (X8 b))]
    (A IYH)  [0xFD (+ 0b1000_1000 (Y8 b))]
    (A IYL)  [0xFD (+ 0b1000_1000 (Y8 b))])

  (opcode  SUB (a)
    (IXH)  [0xDD (+ 0b1001_0000 (X8 a))]
    (IXL)  [0xDD (+ 0b1001_0000 (X8 a))]
    (IYH)  [0xFD (+ 0b1001_0000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1001_0000 (Y8 a))])

  (opcode  SBC (a b)
    (A IXH)  [0xDD (+ 0b1001_1000 (X8 b))]
    (A IXL)  [0xDD (+ 0b1001_1000 (X8 b))]
    (A IYH)  [0xFD (+ 0b1001_1000 (Y8 b))]
    (A IYL)  [0xFD (+ 0b1001_1000 (Y8 b))])

  (opcode  AND (a)
    (IXH)  [0xDD (+ 0b1010_0000 (X8 a))]
    (IXL)  [0xDD (+ 0b1010_0000 (X8 a))]
    (IYH)  [0xFD (+ 0b1010_0000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1010_0000 (Y8 a))])

  (opcode  OR (a)
    (IXH)  [0xDD (+ 0b1011_0000 (X8 a))]
    (IXL)  [0xDD (+ 0b1011_0000 (X8 a))]
    (IYH)  [0xFD (+ 0b1011_0000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1011_0000 (Y8 a))])

  (opcode  XOR (a)
    (IXH)  [0xDD (+ 0b1010_1000 (X8 a))]
    (IXL)  [0xDD (+ 0b1010_1000 (X8 a))]
    (IYH)  [0xFD (+ 0b1010_1000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1010_1000 (Y8 a))])

  (opcode  CP (a)
    (IXH)  [0xDD (+ 0b1011_1000 (X8 a))]
    (IXL)  [0xDD (+ 0b1011_1000 (X8 a))]
    (IYH)  [0xFD (+ 0b1011_1000 (Y8 a))]
    (IYL)  [0xFD (+ 0b1011_1000 (Y8 a))])

  (opcode  INC (a)
    (IXH)  [0xDD (+ 0b0000_0100 (X8 a 3))]
    (IXL)  [0xDD (+ 0b0000_0100 (X8 a 3))]
    (IYH)  [0xFD (+ 0b0000_0100 (Y8 a 3))]
    (IYL)  [0xFD (+ 0b0000_0100 (Y8 a 3))])
  (example INC
    (IXH) _ "DB 0xDD, 0x24"
    (IXL) _ "DB 0xDD, 0x2C"
    (IYH) _ "DB 0xFD, 0x24"
    (IYL) _ "DB 0xFD, 0x2C")

  (opcode  DEC (a)
    (IXH)  [0xDD (+ 0b0000_0101 (X8 a 3))]
    (IXL)  [0xDD (+ 0b0000_0101 (X8 a 3))]
    (IYH)  [0xFD (+ 0b0000_0101 (Y8 a 3))]
    (IYL)  [0xFD (+ 0b0000_0101 (Y8 a 3))])
  (example DEC
    (IXH) _ "DB 0xDD, 0x25"
    (IXL) _ "DB 0xDD, 0x2D"
    (IYH) _ "DB 0xFD, 0x25"
    (IYL) _ "DB 0xFD, 0x2D")

  (opcode  IN (a b) (F  C$) [0xED 0x70])
  (example IN (F C$) _ "DB 0xED, 0x70")

  (opcode  MULUB (a b)
    (A B) [0xED 0xC1]
    (A C) [0xED 0xC9]
    (A D) [0xED 0xD1]
    (A E) [0xED 0xD9])
  (example MULUB
    (A B) _ "DB 0xED, 0xC1"
    (A C) _ "DB 0xED, 0xC9"
    (A D) _ "DB 0xED, 0xD1"
    (A E) _ "DB 0xED, 0xD9")

  (opcode  MULUW (a b)
    (HL BC) [0xED 0xC3]
    (HL SP) [0xED 0xF3])
  (example MULUW
    (HL BC) _ "DB 0xED, 0xC3"
    (HL SP) _ "DB 0xED, 0xF3")

  (operator <- (a b)
    (IXH _)  [(LD (= a) (= b))]
    (IXL _)  [(LD (= a) (= b))]
    (IYH _)  [(LD (= a) (= b))]
    (IYL _)  [(LD (= a) (= b))])
  (operator -> (a b)
    (IXH _)  [(LD (= b) (= a))]
    (IXL _)  [(LD (= b) (= a))]
    (IYH _)  [(LD (= b) (= a))]
    (IYL _)  [(LD (= b) (= a))])

  (example -in  (F C) "F -in C" "DB 0xED, 0x70"))
