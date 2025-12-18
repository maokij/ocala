;; -*- mode: Lisp; lisp-indent-offset: 2 -*-
(arch ttarch
  (operand A   RegA   "A"    "A")
  (operand B   RegB   "B"    "B")
  (operand P   RegP   "P"    "P")
  (operand AB  RegAB  "AB"   "AB")
  (operand SP  RegSP  "SP"   "SP")
  (operand PC  RegPC  "PC"   "PC")
  (operand PQ  RegPQ  "PQ"   "PQ")
  (operand X   RegX   "X"    "X")
  (operand X$  MemX   "[X]"  "(X)")
  (operand Y   RegY   "Y"    "Y")
  (operand Y$  MemY   "[Y]"  "(Y)")
  (operand N   ImmN   "%B"   "0+ %B" NN temp)
  (operand NN  ImmNN  "%W"   "0+ %W" N)
  (operand NN$ MemNN  "[%W]" "(%W)")
  (operand NE? CondNE "NE?"  "NE")
  (operand EQ? CondEQ "EQ?"  "EQ")
  (operand CC? CondCC "CC?"  "CC")
  (operand CS? CondCS "CS?"  "CS")

  (registers A B P X Y AB SP PC)
  (conditions (NE? !=?) (EQ? ==?) (CC? >=?) (CS? <?))
  (map R8 A 0 B 1)
  (map RR AB 0 X 1 Y 2 SP 3)
  (map CO NE? 0 EQ? 1 CC? 2 CS? 3)

  (opcode  NOP ()    ()       [0x00])
  (aliases NOP NOOP)
  (opcode  JMP (a)   (NN)     [0x01 (=l a) (=h a)])
  (opcode  JMP (a b) (NN CO)  [(+ 0x10 (CO b)) (=l a) (=h a)])
  (opcode  BRL (a)   (NN)     [0x02 (=rl a -3 2) (=rh a -3 2)])
  (opcode  BRA (a)   (NN)     [0x03 (=rl a -2)])
  (opcode  RET ()    ()       [0x04])
  (opcode  RET (a)   (CO)     [(+ 0x20 (CO a))])
  (opcode  JSR (a)   (NN)     [0x05 (=l a) (=h a)])

  (opcode  #.jump (a)   (NN)    [0x01 (=l a) (=h a)])
  (opcode  #.jump (a b) (NN CO) [(+ 0x10 (CO b)) (=l a) (=h a)])
  (example #.jump (*) "" "")

  (opcode  #.call (a)   (NN)    [0x05 (=l a) (=h a)])
  (opcode  #.call (a b) (NN CO) [(=U)])
  (example #.call (*) "" "")

  (opcode  #.return ()  ()   [0x04])
  (opcode  #.return (a) (CO) [(+ 0x20 (CO a))])
  (example #.return (*) "" "")

  (opcode  LD (a b)
    (A B)   [0x20]
    (A X$)  [0x21]
    (A Y$)  [0x22]
    (A N)   [0x23 (=l b)]

    (B A)   [0x24]
    (B X$)  [0x25]
    (B Y$)  [0x26]
    (B N)   [0x27 (=l b)]

    (X AB)  [0x28]
    (X Y)   [0x29]
    (X SP)  [0x2A]
    (X NN)  [0x2B (=l b) (=h b)]

    (Y AB)  [0x2C]
    (Y X)   [0x2D]
    (Y SP)  [0x2E]
    (Y NN)  [0x2F (=l b) (=h b)]

    (AB X)  [0x38]
    (AB Y)  [0x39]
    (AB SP) [0x3A]
    (AB NN) [0x3B (=l b) (=h b)]

    (A NN$) [0x30 (=l b) (=h b)]
    (B NN$) [0x31 (=l b) (=h b)]
    (X NN$) [0x32 (=l b) (=h b)]
    (Y NN$) [0x33 (=l b) (=h b)]

    (NN$ A) [0x34 (=l a) (=h a)]
    (NN$ B) [0x35 (=l a) (=h a)]
    (NN$ X) [0x36 (=l a) (=h a)]
    (NN$ Y) [0x37 (=l a) (=h a)])

  (opcode  ADD (a b)
    (A  A)  [0x40]
    (A  B)  [0x41]
    (A  N)  [0x42 (=l b)]
    (A  X$) [0x43]
    (A  Y$) [0x44]
    (AB X)  [0x45]
    (AB Y)  [0x46])

  (opcode  BIT (a b)
    (N A)   [(=i a 0b0101_0000 0x07 0)]
    (N B)   [(=i a 0b0101_1000 0x07 0)])

  (bytemap BMM 0b0011 1 3 0 0x61 0x67 0x6F)
  (opcode  BMM (a) (N) [(=m a BMM)])

  (opcode  DNN (a)
    (N)  [0x64 (=l a)]
    (NN) [0x65 (=l a) (=h a)])

  (opcode  RRA () () [0x63])

  (operator <- (a b)
    (A  _)  [(LD (= a) (= b))]
    (B  _)  [(LD (= a) (= b))]
    (X  _)  [(LD X (= b))]
    (Y  _)  [(LD Y (= b))]
    (AB _)  [(LD (= a) (= b))]
    (AB PQ) [(LD (= a) 0x1234)])

  (operator -jump    (a)   (NN)    [(#.jump (= a))])
  (operator -jump-if (a b) (NN CO) [(#.jump (= a) (= b))])
  (operator -jump-unless (a b)
    (NN NE?) [(#.jump (= a) EQ?)]
    (NN EQ?) [(#.jump (= a) NE?)]
    (NN CC?) [(#.jump (= a) CS?)]
    (NN CS?) [(#.jump (= a) CC?)])

  (operator -return    (a)   (PC)    [(#.return)])
  (operator -return-if (a b) (PC CO) [(#.return (= b))])
  (operator -return-unless (a b)
    (PC NE?) [(#.return EQ?)]
    (PC EQ?) [(#.return NE?)]
    (PC CC?) [(#.return CS?)]
    (PC CS?) [(#.return CC?)])

  (operator -dnnm (a) (NN$) [(DNN (= a NN))])
  (operator -byte (a) (NN) [(#.BYTE (= a))])
  (operator -rep (a b) (NN NN) [(#.REP (= b) `[(#.BYTE (= a))])]))

(arch (ttarch +ext)
  (opcode  EXT (a b) (A N) [0x70 (=l b)])
  (operator <- (a b) (A NN) [(EXT A (= b))])
  (operator -jump (a b) (A NN) [(EXT A (= b))])
  (operator -ext (a b) (A NN) [(EXT A (= b))]))
