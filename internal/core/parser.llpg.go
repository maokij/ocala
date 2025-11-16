package core

/**:rules
  %options package: core
  program:           ^{ v := &Vec{KwBlock.ToId(p.PeekToken())} }!
    statement?        { if _1 != nil { v.Push(_1) } }!
    ( ';' statement?  { if _2 != nil { v.Push(_2) } }!
    )*                { RET v }!
  ;

  statement:
    MACRO identifierp                  { v := &Vec{}; w := &Vec{} }!
    '(' ( identifier                   { v.Push(&Vec{_1, NIL}) }!
        | label_name data_expression   { v.Push(&Vec{_1, _2}) }!
        )* REST?:NILTK ')'
    ( '[' ( identifier                 { l := _1; r := Value(NIL) }!
            ('=' symbol_expression     { r = _2 }!
            )?                         { w.Push(l, r) }!
          )* ']' )? block              { RET &Vec{KwMacro.ToId(_1), _2, v, w, _5.Value, _8} }!
  | PROC identifierp '(' signature ')' { v := &Vec{KwProc.ToId(_1), _2, _4} }!
    ( block                            { v.Push(_1) }!
    | '@' data_expression              { v.Push(_2) }!
    )                                  { RET v }!
  | CONST                              { v := &Vec{} }!
    ( identifier                       { RET _1 }!
    | identifierp
      '(' ( identifier                 { v.Push(&Vec{_1, NIL}) }!
          | label_name data_expression { v.Push(&Vec{_1, _2}) }!
          )* ')'                       { RET _1 }!
    ) '=' data_expression              { RET &Vec{KwConst.ToId(_1), _2, v, _4} }!
  | DATA data_type                     { v := Value(NIL); w := _2 }!
    ( '=' data_type                    { v = w; w = _2 }!
    )? data_body                       { RET &Vec{KwData.ToId(_1), v, w, _4} }!
  | MODULE identifier block            { RET &Vec{KwModule.ToId(_1), _2, _3} }!
  | LABEL label_name                   { RET &Vec{KwLabel.ToId(_1), _2} }!
  | struct                             { v := _1.(*Vec); v.SetAt(1, NIL); RET v }!
  | identifier                         { v := &Vec{_idfrom(_1)} }!
    ( operand                          { v.Push(_1) }!
    )*                                 { RET v }!
  | proc_call                          { RET _1 }!
  | context_expression                 { RET _1 }!
  ;

  struct:
    STRUCT               { v := &Vec{KwStruct.ToId(_1), Int(1), NIL} }!
    identifier?          { if _2 != nil { v.SetAt(2, _2) } }!
    '{' struct_field?    { if _4 != nil { v.Push(_4) } }!
    ( ';' struct_field?  { if _2 != nil { v.Push(_2) } }!
    )* '}'               { RET v }!
  ;
  struct_field:
    identifier data_type { RET &Vec{_1, _2} }!
  ;

  data_type:
    '['              { v := Value(NIL) }!
    ( constexpr      { v = _1 }!
    )? ']' data_type { RET &Vec{KwArray.ToId(_1), _4, v} }!
  | identifier       { RET _1 }!
  | struct           { RET _1 }!
  ;
  data_body:                     ^{ values := Value(NIL) }!
    ( data_list                   { values = _1 }!
    | constval                    { values = _1 }!
    )?                            { alloc := Value(idMulOp);
                                    size := Value(InternalConstexpr(Int(1))) }!
    ( BINARY_OPERATOR data_value  { alloc = _1.Value; size = _2 }!
    )?                            { section := Value(NIL); addr := Value(NIL) }!
    ( ':' identifier              { section = _2 }!
    | '@' data_expression         { addr = _2 }!
    )?                            { RET &Vec{values, alloc, size, section, addr } }!
  ;
  data_list:
    '['               { v := &Vec{KwDataList.ToId(_1)} }!
    ( data_expression { v.Push(_1) }!
    | data_list       { v.Push(_1) }!
    )* ']'            { RET v }!
  | '{'               { v := &Vec{KwStructData.ToId(_1)}; p.SetContext('#') }!
    ( data_expression { v.Push(_1) }!
    | data_list       { v.Push(_1) }!
    )* '}'            { RET v }!
  ;
  data_expression:
    constexpr      { RET _1 }!
  | explicit_value { RET _1 }!
  ;
  data_value:
    constval       { RET _1 }!
  | explicit_value { RET _1 }!
  ;

  symbol_expression:
    '%{' symbol_expression_value { v := &Vec{KwMakeId.ToId(_1), _2} }!
    ( symbol_expression_value    { v.Push(_1) }!
    )* '}'                       { RET _constexpr(v, _1) }!
  ;
  symbol_expression_value:
    IDENTIFIER { RET _1.Value }!
  | STRING     { RET _1.Value }!
  ;

  signature:          ^{ v := &Vec{NIL, NIL, NIL, NIL} }!
    ( '-*'             { v.SetAt(3, Int(1)) }!
    )?
    registers?         { if _2 != nil { v.SetAt(0, _2) } }!
    ( '=>' registers?  { if _2 != nil { v.SetAt(1, _2) } }!
    )?
    ( '!' registers?   { if _2 != nil { v.SetAt(2, _2)
                         } else { v.SetAt(2, &Vec{}) } }!
    )?                 { RET v }!
  ;
  registers:
    REGISTER    { v := &Vec{_1.Value} }!
    ( REGISTER  { v.Push(_1.Value) }!
    )*          { RET v }!
  ;
  block:
    '{' program '}'  { v := _2.(*Vec); v.SetAt(0, KwProg.ToId(_1)); RET v }!
  | '={' program '}' { RET _2 }!
  ;
  proc_call:
    CONDDOT?:NILTK
    identifierp
    '(' signature ')' { RET &Vec{KwCallproc.ToId(_3), _2, _1.Value, _4, KwCall} }!
  ;

  context_expression:          ^{ v := &Vec{KwWith.ToId(p.PeekToken())} }!
    ( decorated_register        { v = _1.(*Vec) }!
    | memory_access             { v.Push(_1) }!
    | explicit_value            { v.Push(_1) }!
    | '@-' primitive            { v.Push(p.cc.KwRegA.ToId(_1), &Vec{_atopid(_2, _1), _2}) }!
    )
    ( POSTFIX_OPERATOR          { v.Push(&Vec{_1.Value, NIL}) }!
    | BINARY_OPERATOR operand   { v.Push(&Vec{_1.Value, _2}) }!
    | DOT_OPERATOR dot_operand  { v.Push(&Vec{_1.Value, _2}) }!
    )*                          { RET v }!
  ;
  operand:
    primitive         { v := _1 }!
      ( ':' primitive { v = &Vec{KwTpl.ToId(_1), v, _2} }!
      )?              { RET v }!
  ;
  primitive:
    CONDITION          { RET _1.Value }!
  | decorated_register { RET _maybeunwrapwith(_1) }!
  | memory_access      { RET _1 }!
  | data_value         { RET _1 }!
  | block              { RET _1 }!
  ;
  <register>
  decorated_register:
    REGISTER         { v := &Vec{KwWith.ToId(_1), _1.Value} }!
    ( '-@' primitive { v.Push(&Vec{_atopid(_2, _1), _2}) }!
    )?               { RET v }!
  ;
  memory_access:
    '['                  { v := &Vec{KwMem.ToId(_1)} }!
    ( context_expression { v.Push(_maybeunwrapwith(_1))}!
    | constexpr          { v.Push(_1)}!
    )* ']'               { RET v }!
  ;
  <operand>
  dot_operand:
    proc_call          { RET _1 }!
  | decorated_register { RET _1 }!
  | block              { RET _1 }!
  ;

  explicit_value:
    '$@'  constval { RET _2 }!
  | '$$@' constval { RET &Vec{KwValueOf.ToId(_1), _2} }!
  ;

  constexpr: ^{ here := p.PeekToken() }!
    iexpr     { RET _constexpr(_1, here) }!
  ;
  constval:  ^{ here := p.PeekToken() }!
    ival      { RET _constexpr(_1, here) }!
  ;
  <constexpr>
  iexpr:
    ival                   { v := []Value{_1} }!
    ( BINARY_OPERATOR ival { v = p.cc.orderByPrec(v, _1.Value.(*Identifier), _2) }!
    )*                     { RET p.cc.orderByPrec(v, idOpLast, nil)[0] }!
  ;
  <constval>
  ival:
    INTEGER              { RET _1.Value }!
  | STRING               { RET _1.Value }!
  | RESERVED             { RET _1.Value }!
  | IDENTIFIER           { v := _1.Value }!
    ( '.-' IDENTIFIER    { w := &Vec{KwField.ToId(_1), v, _2.Value} }!
      ( '.-' IDENTIFIER  { w.Push(_2.Value) }!
      )*                 { v = w }!
    )?                   { RET v }!
  | IDENTIFIERP          { v := &Vec{ _1.Value } }!
    '(' ( iexpr          { v.Push(_1) }!
        )* ')'           { RET v }!
  | '(' iexpr ')'        { RET _2 }!
  | PREFIX_OPERATOR ival { RET &Vec{ _1.Value, _2 } }!
  ;
  identifier:
    IDENTIFIER { RET _constexpr(_1.Value, _1) }!
  ;
  <identifier>
  identifierp:
    IDENTIFIERP { RET _constexpr(_1.Value, _1) }!
  ;
  label_name:
    LABEL_NAME { RET _constexpr(_1.Value, _1) }!
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

/**:symbolDescriptions:
  identifier: IDENTIFIER IDENTIFIERP
  register: REGISTER
  '.': DOT_OPERATOR
  'macro': MACRO
  'proc': PROC
  'const': CONST
  'data': DATA
  'module': MODULE
  'struct': STRUCT
 **/ //:namemap:

var NILTK = &Token{From: InternalParser, Value: NIL}
var idMulOp = InternalId(KwMulOp)

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
