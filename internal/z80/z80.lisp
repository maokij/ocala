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
  (operand PQ  RegPQ  "PQ"   "PQ")

  (operand IX  RegIX  "IX"      "IX")
  (operand IX$ MemIX  "[IX %B]" "(IX+%B)")
  (operand IY  RegIY  "IY"      "IY")
  (operand IY$ MemIY  "[IY %B]" "(IY+%B)")

  (operand N   ImmN   "%B"   "0+ %B" NN temp)
  (operand N$  MemN   "[%B]" "(%B)"  NN$ temp)
  (operand NN  ImmNN  "%W"   "0+ %W" N)
  (operand NN$ MemNN  "[%W]" "(%W)"  N$)
  (operand C$  MemC   "[C]"  "(C)")
  (operand I   RegI   "I"    "I")
  (operand R   RegR   "R"    "R")

  (operand NZ? CondNZ "NZ?"  "NZ")
  (operand Z?  CondZ  "Z?"   "Z")
  (operand NC? CondNC "NC?"  "NC")
  (operand C?  CondC  "C?"   "C")
  (operand PO? CondPO "PO?"  "PO")
  (operand PE? CondPE "PE?"  "PE")
  (operand P?  CondP  "P?"   "P")
  (operand M?  CondM  "M?"   "M")

  (registers A B C D E H L I R AF AF- BC DE HL IX IY SP)
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

  (example "$prologue" (*) "arch z80; flat!" "")

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

  (bytemap JP 0xFF 0 0 0xE9)
  (opcode  JP (a)
    (NN)  [0xC3 (=l a) (=h a)]
    (HL$) [0xE9]
    (IX$) [0xDD (=m a JP)]
    (IY$) [0xFD (=m a JP)])
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
    (DD PQ) [(#.LDP (= a) (= b))])
  (operator -> (a b)
    (R8 _)  [(LD (= b) (= a))]
    (I _)   [(LD (= b) (= a))]
    (R _)   [(LD (= b) (= a))]
    (DD _)  [(LD (= b) (= a))]
    (IX _)  [(LD (= b) (= a))]
    (IY _)  [(LD (= b) (= a))]
    (NN _)  [(LD (= b) (= a))]
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

  (operator -jump-if (a b)
    (NN CC) [(JP (= b) (= a))])
  (example -jump-if (*) "" "")

  (operator -jump-unless (a b)
    (NN NZ?) [(JP Z?  (= a))]
    (NN Z?)  [(JP NZ? (= a))]
    (NN NC?) [(JP C?  (= a))]
    (NN C?)  [(JP NC? (= a))]
    (NN PO?) [(JP PE? (= a))]
    (NN PE?) [(JP PO? (= a))]
    (NN M?)  [(JP P?  (= a))]
    (NN P?)  [(JP M?  (= a))])
  (example -jump-unless (*) "" "")

  (example $operators (*)
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
    "B >>> 2" "SRL B; SRL B"

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

    "proc f(!){ RET }" "f: RET"
    "f(!)" "CALL f"
    "NC?.f(!)" "CALL NC, f"))
