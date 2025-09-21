package core

/**:rules
  %options package: core
  program:       ^{ v := &Vec{_markerid(p, KwBlock)} }!
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
  | PROC identifierp '(' sig ')'       { v := &Vec{KwProc.ToId(_1), _2, _4} }!
      ( block                          { v.Push(_1) }!
      | '@' dataexpr                   { v.Push(_2) }!
      )                                { RET v }!
  | CONST                              { v := &Vec{} }!
      ( identifier                     { RET _1 }!
      | identifierp
          '(' ( identifier             { v.Push(&Vec{_1, NIL}) }!
              | label ':' dataexpr     { v.Push(&Vec{_1, _3}) }!
              )* ')'                   { RET _1 }!
      ) '=' dataexpr                   { RET &Vec{KwConst.ToId(_1), _2, v, _4} }!
  | DATA datatype                      { v := Value(NIL); w := _2 }!
      ( '=' datatype                   { v = w; w = _2 }!
      )? databody                      { RET &Vec{KwData.ToId(_1), v, w, _4} }!
  | MODULE identifier block            { RET &Vec{KwModule.ToId(_1), _2, _3} }!
  | label ':'                          { RET &Vec{KwLabel.ToId(_2), _1}; p.state = pstNl }!
  | struct                             { v := _1.(*Vec); v.SetAt(1, NIL); RET v }!
  | expr                               { RET _1 }!
  ;
  struct:
    STRUCT                             { v := &Vec{KwStruct.ToId(_1), Int(1), NIL} }!
      identifier?                      { if _2 != nil { v.SetAt(2, _2) } }!
      '{' structfield?                 { if _4 != nil { v.Push(_4) } }!
          ( ';' structfield?           { if _2 != nil { v.Push(_2) } }!
          )* '}'                       { RET v }!
  ;
  structfield:
    identifier datatype                { RET &Vec{_1, _2} }!
  ;
  datatype:
    '['                                { v := Value(NIL) }!
      ( constexpr                      { v = _1 }!
      )? ']' datatype                  { RET &Vec{KwArray.ToId(_1), _4, v} }!
  | identifier                         { RET _1 }!
  | struct                             { RET _1 }!
  ;
  databody:                           ^{ values := Value(NIL); }!
    ( datalist                         { values = _1 }!
    | constval                         { values = _1 }!
    )?                                 { alloc := Value(idMulOp);
                                         size := Value(InternalConstexpr(Int(1))) }!
      ( BOP dataval                    { alloc = _1.Value; size = _2 }!
      )?                               { section := Value(NIL); addr := Value(NIL) }!
      ( ':' identifier                 { section = _2 }!
      | '@' dataexpr                   { addr = _2 }!
      )?                               { RET &Vec{values, alloc, size, section, addr } }!
  ;
  datalist:
    '['                                { v := &Vec{KwDataList.ToId(_1)} }!
        ( dataexpr                     { v.Push(_1) }!
        | datalist                     { v.Push(_1) }!
        )* ']'                         { RET v }!
  | '{'                                { v := &Vec{KwStructData.ToId(_1)}
                                         p.contexts = append(p.contexts, '#')}!
        ( dataexpr                     { v.Push(_1) }!
        | datalist                     { v.Push(_1) }!
        )* '}'                         { RET v }!
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
  contextexpr:         ^{ v := &Vec{_markerid(p, KwWith)} }!
    ( regld             { v = _1.(*Vec) }!
    | mem               { v.Push(_1) }!
    | explicitval       { v.Push(_1) }!
    | '@-' prim         { v.Push(p.cc.KwRegA.ToId(_1), &Vec{_atopid(_2, _1), _2}) }!
    ) ( UOP             { v.Push(&Vec{_1.Value, NIL}) }!
      | BOP oper        { v.Push(&Vec{_1.Value, _2}) }!
      | DOP dotarg      { v.Push(&Vec{_1.Value, _2}) }!
      )*                { RET v }!
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
      ( '-@' prim     { v.Push(&Vec{_atopid(_2, _1), _2}) }!
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
    '$@' constval     { RET _2 }!
  | '$$@' constval    { RET &Vec{KwValueOf.ToId(_1), _2} }!
  ;
  constexpr:         ^{ here := _marker(p) }!
    iexpr             { RET _constexpr(_1, here) }!
  ;
  constval:          ^{ here := _marker(p) }!
    ival              { RET _constexpr(_1, here) }!
  ;
  iexpr:
    ival              { v := []Value{_1} }!
      ( BOP ival      { v = p.cc.orderByPrec(v, _1.Value.(*Identifier), _2) }!
      )*              { RET p.cc.orderByPrec(v, idOpLast, nil)[0] }!
  ;
  ival:
    INTEGER                { RET _1.Value }!
  | STRING                 { RET _1.Value }!
  | RESERVED               { RET _1.Value }!
  | IDENTIFIER             { v := _1.Value }!
      ( '.-' IDENTIFIER    { w := &Vec{KwField.ToId(_1), v, _2.Value} }!
        ( '.-' IDENTIFIER  { w.Push(_2.Value) }!
        )*                 { v = w }!
      )?                   { RET v }!
  | IDENTIFIERP            { v := &Vec{ _1.Value } }!
      '(' ( iexpr          { v.Push(_1) }!
          )* ')'           { RET v }!
  | '(' iexpr ')'          { RET _2 }!
  | AOP ival               { RET &Vec{ _1.Value, _2 } }!
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
  RESERVED: _ _BEG _END _COND
 **/ //:reservedWords:

var NILTK = &Token{From: InternalParser, Value: NIL}
var idMulOp = InternalId(KwMulOp)

func _marker(p *Parser) *Token {
	if len(p.Tokens) == 0 {
		token := p.NewToken(0, NIL, 0)
		p.pending = append(p.pending, token)
		return token
	}
	return p.Tokens[0]
}

func _markerid(p *Parser, v *Keyword) Value {
	if len(p.Tokens) == 0 {
		token := p.NewToken(0, NIL, 0)
		p.pending = append(p.pending, token)
		id := &Identifier{Name: v, Token: token}
		token.Value = id
		return id
	}
	return v.ToId(p.Tokens[0])
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

func _atopid(e Value, token *Token) *Identifier {
	if v := AsBlockForm(e); v != nil {
		return KwDot.ToId(token)
	}
	return KwLeftArrow.ToId(token)
}
