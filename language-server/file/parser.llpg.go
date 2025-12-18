package file

import "ocala/core"

/**:rules
  %options package: file, parser_type: "*File"
  %options value_type: "core.Value", token_type: "*core.Token"
  %options prologue: 'import "ocala/core"'
  %options recover: ["'={'", "'{'", "'}'", "'('", "')'", "'['", "']'", \
                     BINARY_OPERATOR, DOT_OPERATOR, "'='", "';'"]
  program:
    statement?
    ( ';'        { p.Recovered(); p.pending = nil }!
      statement?
    )*
  ;
  statement:
    MACRO identifierp                { r := &MacroNode{}; p.HoldNode(r, _2) }!
    '(' ( identifier                 { r.AddArg(_1) }!
        | label_name data_expression { r.AddArg(_1) }!
        )* REST? ')'                 { if _4 != nil { r.AddArg(_4) } }!
    ( '[' ( identifier               { r.AddVar(_1) }!
            ('=' symbol_expression
            )?
          )* ']' )? block
  | PROC identifierp                 { r := &ProcNode{}; p.HoldNode(r, _2) }!
    '(' signature ')'                { r.SetSignature(_4) }!
    ( block
    | '@' data_expression
    )
  | CONST
    ( identifier                     { p.AddNode(&ConstNode{}, _1) }!
    | identifierp '('                { r := &ConstFnNode{}; p.AddNode(r, _1) }!
      ( identifier                   { r.AddArg(_1) }!
      | label_name data_expression   { r.AddArg(_1) }!
      )* ')'
    ) '=' data_expression
  | DATA data_type                   { r := &DataNode{}; t := _2 }!
    ( '=' data_type                  { p.AddNode(r, t); t = _2 }!
    )?                               { r.Type = _expectid(t) }!
    data_body
  | MODULE identifier                { p.HoldNode(&ModuleNode{}, _2) }!
    block
  | LABEL label_name                 { p.AddNode(&LabelNode{}, _2) }!
  | struct
  | identifier                       { v := _1 }!
    ( operand                        { _checkinclude(p, v, _1) }!
    )*
  | proc_call
  | context_expression
  ;
  struct:
    STRUCT identifier?  { name := _ensureid(_2); p.HoldNode(&StructNode{}, name) }!
    '{'                 { p.EnterNode(_3) }!
    struct_field?
    ( ';' struct_field?
    )* '}'              { p.LeaveNode(); RET name }!
  ;
  struct_field:
    identifier data_type { r := &StructFieldNode{}; p.AddNode(r, _1); r.Type = _expectid(_2) }!
  ;
  data_type:
    '['
    ( constexpr
    )? ']' data_type { RET _4 }!
  | identifier       { RET _1 }!
  | struct           { RET _1 }!
  ;
  data_body:
    ( data_list
    | constval
    )?
    ( BINARY_OPERATOR data_value
    )?
    ( ':' identifier
    | '@' data_expression
    )?
  ;
  data_list:
    '[' ( data_expression
        | data_list
        )* ']'
  | '{' { p.SetContext('#') }!
    ( data_expression
    | data_list
    )* '}'
  ;
  data_expression:
    constexpr
  | explicit_value
  ;
  data_value:
    constval       { RET _1 }!
  | explicit_value
  ;
  symbol_expression:
    '%{' symbol_expression_value
    ( symbol_expression_value
    )* '}'
  ;
  symbol_expression_value:
    IDENTIFIER
  | STRING
  ;
  signature:          ^{ var t, u, v, w core.Value }!
    ( '-*'             { t = core.Int(1) }!
    )?
    registers?         { u = _2 }!
    ( '=>' registers?  { v = _2 }!
    )?
    ( '!' registers?   { w = _siginvalidated(_2) }!
    )?                 { RET _newsig(t, u, v, w) }!
  ;
  registers:
    REGISTER   { v:= &core.Vec{_1.Value} }!
    ( REGISTER { v.Push(_1.Value) }!
    )*         { RET v }!
  ;
  block:
    ( '{'        { p.EnterNode(_1) }!
    | '={'       { p.EnterNode(_1) }!
    )
    statement?
    ( ';'        { p.Recovered() }!
      statement?
    )* '}'       { p.LeaveNode() }!
  ;
  proc_call:
    CONDDOT?
    identifierp
    '(' signature ')'
  ;
  context_expression:
    ( decorated_register
    | memory_access
    | explicit_value
    | '@-' primitive
    )
    ( POSTFIX_OPERATOR
    | BINARY_OPERATOR operand
    | DOT_OPERATOR dot_operand
    )*
  ;
  operand:
    primitive
    ( ':' primitive
    )?                 { RET _1 }!
  ;
  primitive:
    CONDITION          { RET _1.Value }!
  | decorated_register { RET _1 }!
  | memory_access      { RET _1 }!
  | data_value         { RET _1 }!
  | block              { RET _1 }!
  ;
  <register>
  decorated_register:
    REGISTER
    ( '-@' primitive
    )?
  ;
  memory_access:
    '['
    ( context_expression
    | constexpr
    )* ']'
  ;
  <operand>
  dot_operand:
    proc_call
  | decorated_register
  | block
  ;
  explicit_value:
    '$@' constval
  | '$$@' constval
  ;
  constexpr:
    iexpr
  ;
  constval:
    ival { RET _1 }!
  ;
  <constexpr>
  iexpr:
    ival
    ( BINARY_OPERATOR ival
    )*
  ;
  <constval>
  ival:
    INTEGER              { RET _1.Value }!
  | STRING               { RET _1.Value }!
  | RESERVED             { RET _1.Value }!
  | IDENTIFIER           { RET _1.Value }!
    ( '.-' IDENTIFIER
      ( '.-' IDENTIFIER
      )*
    )?
  | IDENTIFIERP          { RET _1.Value }!
    '(' ( iexpr
        )* ')'
  | '(' iexpr ')'        { RET _1.Value }!
  | PREFIX_OPERATOR ival
  ;
  identifier:
    IDENTIFIER { RET _1.Value }!
  ;
  <identifier>
  identifierp:
    IDENTIFIERP { RET _1.Value }!
  ;
  label_name:
    LABEL_NAME { RET _1.Value }!
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

var SYNTAX_ERROR = &core.Token{}

func _siginvalidated(v core.Value) core.Value {
	if v == nil {
		return &core.Vec{core.IdUNDER}
	}
	return v
}

func _newsig(t, u, v, w core.Value) *core.Sig {
	cast := func(v core.Value) []*core.Keyword {
		if v == nil {
			return nil
		}
		r := []*core.Keyword{}
		for _, i := range *v.(*core.Vec) {
			r = append(r, i.(*core.Identifier).Name)
		}
		return r
	}
	return &core.Sig{
		IsInline:    t != nil,
		Required:    cast(u),
		Results:     cast(v),
		Invalidated: cast(w),
	}
}

var kwInclude = core.Intern("include")

func _checkinclude(f *File, id core.Value, arg core.Value) {
	if id := kwInclude.MatchId(id); id != nil {
		if s, ok := arg.(*core.Str); ok {
			f.Includes = append(f.Includes, s.String())
			// fmt.Fprintln(os.Stderr, "include:", id, arg)
		}
	}
}

func _ensureid(id core.Value) *core.Identifier {
	if id, ok := id.(*core.Identifier); ok {
		return id
	}
	return &core.Identifier{Name: core.Gensym("id"), Token: invalidToken}
}

func _expectid(v core.Value) *core.Identifier {
	if id, ok := v.(*core.Identifier); ok {
		return id
	}
	return nil
}
