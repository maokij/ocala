package core

/**:rules
  %options package: core
  program:       ^{ v := &Vec{KwBlock.ToId(_here(p))} }!
    stmt?         { if _1 != nil { v.Push(_1) } }!
      ( ';' stmt? { if _2 != nil { v.Push(_2) } }!
      )*          { RET v }!
  ;
  stmt:
    MACRO identifierp                  { v := &Vec{}; w := &Vec{} }!
      '(' ( identifier                 { v.Push(&Vec{_1, NIL}) }!
          | label ':' dataexpr         { v.Push(&Vec{_1, _3}) }!
          )* REST?:NILTK ')'
      ( '[' ( identifier               { l := _1; r := Value(NIL) }!
                ('=' symexpr           { r = _2 }!
                )?                     { w.Push(l, r) }!
            )* ']' )? block            { RET &Vec{KwMacro.ToId(_1), _2, v, w, _5.Value, _8} }!
  | PROC identifierp '(' sig ')' block { RET &Vec{KwProc.ToId(_1), _2, _4, _6} }!
  | CONST                              { v := &Vec{} }!
      ( identifier                     { RET _1 }!
      | identifierp
          '(' ( identifier             { v.Push(&Vec{_1, NIL}) }!
              | label ':' dataexpr     { v.Push(&Vec{_1, _3}) }!
              )* ')'                   { RET _1 }!
      ) '=' dataexpr                   { RET &Vec{KwConst.ToId(_1), _2, v, _4} }!
  | DATA                               { v := Value(NIL); var w Value }!
      ( identifier                     { w = _1 }!
          ( '=' ( identifier           { v = w; w = _1 }!
                | struct               { v = w; w = _1 }!
                ) )?
      | struct                         { w = _1 }!
      ) databody                       { RET &Vec{KwData.ToId(_1), v, w, _3} }!
  | MODULE identifier block            { RET &Vec{KwModule.ToId(_1), _2, _3} }!
  | label ':'                          { RET &Vec{KwLabel.ToId(_2), _1}; p.state = pstNl }!
  | expr                               { RET _1 }!
  ;
  struct:
    STRUCT identifier? '{'             { v := &Vec{} }!
       labeleddata?                    { if _2 != nil { v.Push(_2) } }!
       ( ';' labeleddata?              { if _2 != nil { v.Push(_2) } }!
       )* '}'                          { RET v }!
  ;
  labeleddata:
    label ':' identifier databody      { RET &Vec{_1, _3, _4} }!
  ;
  databody:                           ^{ here := _here(p) }!
    datalist?:NIL                      { alloc := Value(KwMulOp.ToId(here))
                                         size := _constexpr(Int(1), here) }!
      ( BOP dataval                    { alloc = _1.Value; size = _2 }!
      )?                               { section := Value(NIL) }!
      ( ':' identifier                 { section = _2 }!
      )?                               { RET &Vec{_1, alloc, size, section } }!
  ;
  datalist:
    '['                                { v := &Vec{KwVec.ToId(_1)} }!
        ( dataexpr                     { v.Push(_1) }!
        | datalist                     { v.Push(_1) }!
        )* ']'                         { RET v }!
  ;
  dataexpr:
    constexpr         { RET _1 }!
  | explicitval       { RET _1 }!
  ;
  dataval:
    constval          { RET _1 }!
  | explicitval       { RET _1 }!
  ;
  symexpr:
    '%{' symval       { v := &Vec{KwMakeId.ToId(_1), _2} }!
      ( symval        { v.Push(_1) }!
      )* '}'          { RET _constexpr(v, _1) }!
  ;
  symval:
    IDENTIFIER        { RET _1.Value }!
  | STRING            { RET _1.Value }!
  ;
  sig:               ^{ v := &Vec{NIL, NIL, NIL, NIL} }!
    ( '-*'            { v.SetAt(3, Int(1)) }!
    )?
      regs            { v.SetAt(0, _2) }!
      ( '=>' regs     { v.SetAt(1, _2) }!
      )?
      ( '!' regs      { v.SetAt(2, _2) }!
      )?              { RET v }!
  ;
  regs:              ^{ v := &Vec{} }!
    ( REG             { v.Push(_1.Value) }!
    )*                { RET v }!
  ;
  block:
    '{' program '}'   { v := _2.(*Vec); v.SetAt(0, KwProg.ToId(_1)); RET v }!
  | '={' program '}'  { RET _2 }!
  ;
  expr:
    identifier        { v := &Vec{_idfrom(_1)} }!
      ( oper          { v.Push(_1) }!
      )*              { RET v }!
  | callproc          { RET _1 }!
  | contextexpr       { RET _1 }!
  ;
  callproc:
    CONDDOT?:NILTK
      identifierp
      '(' sig ')'       { RET &Vec{KwCallproc.ToId(_3), _2, _1.Value, _4, KwCall} }!
  ;
  contextexpr:         ^{ v := &Vec{KwWith.ToId(_here(p))} }!
    ( regld             { v = _1.(*Vec) }!
    | mem               { v.Push(_1) }!
    | explicitval       { v.Push(_1) }!
    | '@-' prim         { v.Push(p.cc.KwRegA.ToId(_1), &Vec{KwLeftArrow.ToId(_1), _2}) }!
    ) ( CONDDOT?:NILTK  { c := _1.Value }!
          ( UOP         { v.Push(&Vec{_1.Value, NIL, c}) }!
          | BOP oper    { v.Push(&Vec{_1.Value, _2, c}) }!
          | DOP dotarg  { v.Push(&Vec{_1.Value, _2}) }!
          ) )*          { RET v }!
  ;
  oper:
    prim              { v := _1 }!
      ( ':' prim      { v = &Vec{KwTpl.ToId(_1), v, _2} }!
      )?              { RET v }!
  ;
  prim:
    COND              { RET _1.Value }!
  | regld             { RET _maybeunwrapwith(_1) }!
  | mem               { RET _1 }!
  | dataval           { RET _1 }!
  | block             { RET _1 }!
  ;
  regld:
    REG               { v := &Vec{KwWith.ToId(_1), _1.Value} }!
      ( '-@' prim     { v.Push(&Vec{KwLeftArrow.ToId(_1), _2}) }!
      )?              { RET v }!
  ;
  mem:
    '['               { v := &Vec{KwMem.ToId(_1)} }!
      ( contextexpr   { v.Push(_maybeunwrapwith(_1))}!
      | constexpr     { v.Push(_1)}!
      )* ']'          { RET v }!
  ;
  dotarg:
    callproc          { RET _1 }!
  | regld             { RET _1 }!
  | block             { RET _1 }!
  ;
  explicitval:
    '$-' constval     { RET _2 }!
  | '$$-' constval    { RET &Vec{KwValueOf.ToId(_1), _2} }!
  ;
  constexpr:         ^{ here := _here(p) }!
    iexpr             { RET _constexpr(_1, here) }!
  ;
  constval:          ^{ here := _here(p) }!
    ival              { RET _constexpr(_1, here) }!
  ;
  iexpr:
    ival              { v := []Value{_1} }!
      ( BOP ival      { v = p.cc.orderByPrec(v, _1.Value.(*Identifier), _2) }!
      )*              { RET p.cc.orderByPrec(v, idOpLast, nil)[0] }!
  ;
  ival:
    INTEGER              { RET _1.Value }!
  | STRING               { RET _1.Value }!
  | RESERVED             { RET _1.Value }!
  | IDENTIFIER           { v := _1.Value }!
      ( '.-' IDENTIFIER  { v = &Vec{KwField.ToId(_1), v, _2.Value} }!
      )*                 { RET v }!
  | IDENTIFIERP          { v := &Vec{ _1.Value } }!
      '(' ( iexpr        { v.Push(_1) }!
          )* ')'         { RET v }!
  | '(' iexpr ')'        { RET _2 }!
  | AOP ival             { RET &Vec{ _1.Value, _2 } }!
  ;
  identifier:
    IDENTIFIER           { RET _constexpr(_1.Value, _1) }!
  ;
  identifierp:
    IDENTIFIERP          { RET _constexpr(_1.Value, _1) }!
  ;
  label:
    LABEL                { RET _constexpr(_1.Value, _1) }!
  ;
 **/ //:rules

/**:reservedWords:
  MACRO: macro
  PROC:  proc
  CONST: const
  DATA:  data
  MODULE: module
  STRUCT: struct
  RESERVED: _ _BEG _END _COND __PROC__
 **/ //:reservedWords:

var NILTK = &Token{From: InternalParser, Value: NIL}

func _here(p *Parser) *Token {
	if len(p.Tokens) == 0 {
		pos, line := p.Pos, p.Line
		p.seekToNextToken(false)
		token := p.NewToken(0, NIL, p.Pos)
		p.Pos, p.Line = pos, line
		return token
	}
	return p.Tokens[0]
}

func _constexpr(e Value, token *Token) Value {
	if e, ok := e.(*Identifier); ok && e.PlaceHolder != "" {
		return e
	}
	return &Constexpr{Body: e, Token: token}
}

func _maybeunwrapwith(e Value) Value {
	if e := KwWith.MatchExpr(e); e.Size() == 2 {
		return e.At(1)
	}
	return e
}

func _idfrom(e Value) *Identifier {
	if e, ok := e.(*Constexpr); ok {
		return e.Body.(*Identifier)
	}
	return e.(*Identifier)
}
