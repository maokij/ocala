package core

import (
	"regexp"
	"strconv"
)

type Parser struct {
	Scanner
	Recovering bool
	state      byte
	contexts   []byte
	cc         *Compiler
}

const (
	pstDefault = iota
	pstWs
	pstNoWs
	pstNl
	pstNoNl
)

func (p *Parser) Recovered() {
	p.Recovering = false
}

func (p *Parser) ErrorAt(token *Token) *InternalError {
	return &InternalError{
		tag: "parse error: ",
		at:  append([]Value{token}, p.cc.nested...),
	}
}

func (p *Parser) ErrorUnexpected(expected string) {
	token := p.PeekToken()
	err := p.ErrorAt(token)
	label := tokenLabels[token.Kind]
	err.With("unexpected %s, expected %s", label, expected)
}

func (p *Parser) PeekToken() *Token {
	if len(p.Tokens) == 0 {
		p.ScanToken()
	}
	return p.Tokens[0]
}

func (p *Parser) Parse() Value {
	return p._parse()
}

func (p *Parser) SetCompiler(cc *Compiler) {
	p.state = 0
	p.contexts = []byte{'{'}
	p.cc = cc
	p.Scanner.Reset()
}

func (p *Parser) SetContext(a byte) {
	p.contexts = append(p.contexts, a)
}

func (p *Parser) seekToNextToken(nl bool) (Pt, bool, bool) {
	pt, bline, comma, comment := p.Pt, int32(-1), false, 0

loop:
	for ; !p.IsEOF(); p.ForwardChar() {
		c := p.PeekChar()
		switch {
		case c == ' ', c == '\t', c == '\r':
			// NOP
		case c == ',' && comment == 0:
			if comma || p.Line != pt.Line {
				p.scanError("invalid comma")
			}
			comma = true
		case c == '\n':
			if comment != 2 {
				comment = 0
				bline = -1
			}
		case c == '/' && comment == 0 && p.IsValidPos(p.Pos+1):
			switch p.Text[p.Pos+1] {
			case '/': // "//"
				p.ForwardChar()
				comment = 1
			case '*': // "/*"
				p.ForwardChar()
				comment = 2
				if bline == -1 {
					bline = p.Line
				}
			default:
				break loop
			}
		case c == '*' && comment == 2 && p.IsValidPos(p.Pos+1):
			if p.Text[p.Pos+1] == '/' {
				p.ForwardChar()
				comment = 0
				if bline == p.Line {
					bline = -1
				}
			}
		case comment == 0:
			break loop
		}
	}

	if comment == 2 {
		p.scanError("the block comment is not terminated")
	} else if nl {
		nl = !comma && pt.Line != p.Line
		if comma && p.Line-pt.Line > 1 {
			p.scanError("the comma followed by blank lines is not allowed")
		} else if bline > -1 {
			p.scanError("the block comment must be followed by new line")
		}
	}

	return p.Pt, pt.Pos != p.Pos, nl
}

var rePlaceholder = regexp.MustCompile(`^%+(=|#|&|>*[<>*])`)
var reInteger = regexp.MustCompile(`^([+-])?((0x[0-9a-fA-F][0-9a-fA-F_]*)|(0b[01][01_]*)|([0-9][0-9_]*))`)
var reIdentifier = regexp.MustCompile(
	`^([_a-zA-Z][-^!$%&*+/<=>?|~_a-zA-Z0-9]*:|::)?([-^!$%&*+/<=>?|~_a-zA-Z][-^!$%&*+/<=>?|~_a-zA-Z0-9]*)`)
var nameChars = `-^!$%&*+/<=>?|~_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
var wsChars = " \t\r\n"
var headChars = " \t\r\n,;([{"
var tailChars = " \t\r\n,;)]}"

func (p *Parser) errorUnlessPlainId(id *Identifier, pt Pt) {
	if id.Namespace != nil {
		p.scanErrorAt(pt, "the name cannot belong to any namespace")
	} else if id.PlaceHolder != "" {
		p.scanErrorAt(pt, "the name cannot be use as placeholders")
	}
}

func (p *Parser) newIdLikeToken(s, t string, ws, nl bool, context byte, pt Pt, ph string) *Token {
	if u, ok := p.cc.TokenAliases[t]; ok {
		t = u
	}
	v := &Identifier{Name: Intern(t), PlaceHolder: ph}

	if s != "" {
		if ph != "" {
			p.scanError("placeholders cannot contain namespaces")
		}
		v.Namespace = Intern(s[:len(s)-1])
	}

	if tk, ok := tokenKinds[t]; ok {
		p.errorUnlessPlainId(v, pt)
		return p.NewIdToken(tk, v, pt)
	} else if n := p.cc.Operators[t]; n != 0 {
		p.errorUnlessPlainId(v, pt)
		if n&1 != 0 && ws && p.MatchChar(tailChars) {
			p.CancelLastTokenIf(nl)
			return p.NewIdToken(tkPOSTFIX_OPERATOR, v, pt)
		} else if n&2 != 0 && ws && p.MatchChar(wsChars) {
			p.CancelLastTokenIf(nl)
			return p.NewIdToken(tkBINARY_OPERATOR, v, pt)
		}
	} else if tk, ok := p.cc.ReservedWords[t]; ok {
		p.errorUnlessPlainId(v, pt)
		if tk == tkCONDITION && p.ScanChar(".") {
			tk = tkCONDDOT
			p.state = pstWs
		}
		return p.NewIdToken(tk, v, pt)
	}

	if p.MatchChar("(") {
		return p.NewIdToken(tkIDENTIFIERP, v, pt)
	} else if s == "" && p.MatchChar(":") {
		if context == '{' {
			p.state = pstNl
			p.PushToken(p.NewToken(tkLABEL, NIL, pt))
		}
		token := p.NewIdToken(tkLABEL_NAME, v, pt)
		p.ForwardChar()
		return token
	}
	return p.NewIdToken(tkIDENTIFIER, v, pt)
}

func (p *Parser) findTokenKind(s string) int32 {
	tk, ok := tokenKinds[s]
	if !ok {
		p.scanError("invalid token '%s'", s)
		return tkSCANERROR
	}
	return tk
}

func stringIf(cond bool, s string) string {
	if cond {
		return s
	}
	return ""
}

var escapeChars = map[byte]byte{
	'0':  0,
	'a':  '\a',
	'b':  '\b',
	'e':  0x1b,
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
}

var invertedEscapeChars = invertEscapeChars()

func invertEscapeChars() map[byte]byte {
	m := make(map[byte]byte, len(escapeChars))
	for k, v := range escapeChars {
		m[v] = k
	}
	delete(m, '\'')
	return m
}

func (p *Parser) scanEscapedChar() (byte, bool) {
	if p.IsEOF() {
		return 0, false
	}

	c := p.GetChar()
	if c, ok := escapeChars[c]; ok {
		return c, true
	}
	if c == 'x' && p.IsValidPos(p.Pos+1) {
		s := []byte{p.GetChar(), p.GetChar()}
		if n, err := strconv.ParseUint(string(s), 16, 8); err == nil {
			return byte(n), true
		}
	}
	return 0, false
}

func (p *Parser) scanQuotedLiteral(close byte, kind string) []byte {
	v := []byte{}
	for !p.IsEOF() {
		c := p.GetChar()
		if c == close {
			return v
		}
		if c == '\n' {
			p.scanError("new line in %s literal", kind)
			continue // skip
		}
		if c == '\\' {
			d, ok := p.scanEscapedChar()
			if !ok {
				p.scanError("invalid character escape")
			}
			c = d
		}
		v = append(v, c)
	}
	p.scanError("%s literal not terminated", kind)
	return nil
}

func (p *Parser) ScanToken() {
	context := p.contexts[len(p.contexts)-1]
	pt, ws, nl := p.seekToNextToken(p.Recovering || context == '{')
	state := p.state
	p.state = 0
	ws = (state != pstNoWs) && (ws || state == pstWs)
	nl = (state != pstNoNl) && (nl || state == pstNl)

	if p.IsEOF() {
		p.PushToken(p.NewToken(tkEOF, NIL, pt))
		return
	}

	if nl {
		nlpt := pt
		if state != pstNl { // always pt.Line > 0
			nlpt = Pt{Pos: p.Lines[pt.Line] - 1, Line: pt.Line - 1}
		}
		p.PushToken(p.NewToken(tkSC, NIL, nlpt))
	}

	// !"#$%&'()*+,-./:;<=>?@[\]^_`{|}~ // -^!$%&*+/<=>?@|~
	token := ScanErrorToken
	switch {
	case p.Scan(rePlaceholder):
		ph := p.Matched[0]
		if !p.Scan(reIdentifier) {
			p.scanError("invalid placeholder")
			break
		}
		token = p.newIdLikeToken(p.Matched[1], p.Matched[2], false, false, context, pt, ph)
	case p.ScanSeq("..."):
		token = p.NewIdToken(tkREST, &Identifier{Name: KwRest}, pt)
	case p.ScanSeq(".."):
		token = p.NewIdToken(tkBINARY_OPERATOR, &Identifier{Name: KwDotDot}, pt)
	case p.ScanChar("."):
		p.CancelLastTokenIf(nl)
		tk := int32(tkDTMI)
		if p.MatchChar(wsChars) {
			if !ws {
				p.scanError("leading whitespaces are required before the operator `.`")
			}
			tk = tkDOT_OPERATOR
		}
		token = p.NewIdToken(tk, &Identifier{Name: KwDot}, pt)
	case p.ScanChar("@"):
		s := p.SubStringFrom(pt.Pos) + stringIf(!p.MatchChar(wsChars), "-")
		if !ws && p.Pos >= 2 && !p.MatchCharAt(headChars, p.Pos-2) {
			s = "-@"
		}
		tk := p.findTokenKind(s)
		token = p.NewToken(tk, NIL, pt)
	case p.ScanChar("$"):
		p.ScanChar("$")
		if p.MatchChar("(") {
			s := p.SubStringFrom(pt.Pos) + "@"
			tk := p.findTokenKind(s)
			token = p.NewToken(tk, NIL, pt)
			break
		}
		p.Pos = pt.Pos
		p.Scan(reIdentifier)
		token = p.newIdLikeToken(p.Matched[1], p.Matched[2], ws, nl, context, pt, "")
	case p.ScanChar("!"):
		if p.ScanChar("=") {
			p.ScanChar("?")
			if p.MatchChar(nameChars) {
				p.scanError("whitespaces required after '!='/'!=?'")
			}
			s := p.SubStringFrom(pt.Pos)
			token = p.newIdLikeToken("", s, ws, nl, context, pt, "")
			break
		}
		if p.MatchChar(tailChars) {
			token = p.NewToken(tkEX, NIL, pt)
			break
		}
		token = p.NewIdToken(tkPREFIX_OPERATOR, &Identifier{Name: KwLogicalNotOp}, pt) // `!`...
	case p.ScanChar("~"):
		if p.ScanChar(wsChars) {
			p.scanError("no whitespace is allowed after the prefix operator")
			break
		}
		token = p.NewIdToken(tkPREFIX_OPERATOR, &Identifier{Name: KwNotOp}, pt)
	case p.ScanChar(";`"):
		tk := p.findTokenKind(p.SubStringFrom(pt.Pos))
		token = p.NewToken(tk, NIL, pt)
	case p.ScanSeq("={"):
		token = p.NewToken(tkEQLC, NIL, pt)
		p.contexts = append(p.contexts, '{')
	case p.ScanSeq("%{"):
		token = p.NewToken(tkPELC, NIL, pt)
		p.contexts = append(p.contexts, '{')
	case p.ScanChar("([{"):
		s := p.SubStringFrom(pt.Pos)
		tk := p.findTokenKind(s)
		token = p.NewToken(tk, NIL, pt)
		p.contexts = append(p.contexts, s[0])
	case p.ScanChar(")]}"):
		tk := p.findTokenKind(p.SubStringFrom(pt.Pos))
		token = p.NewToken(tk, NIL, pt)
		if context == '#' {
			p.contexts = p.contexts[:len(p.contexts)-1]
		}
		p.contexts = p.contexts[:len(p.contexts)-1]
		if len(p.contexts) == 0 {
			p.contexts = append(p.contexts, '{')
		}
	case p.ScanChar(`'`):
		v := p.scanQuotedLiteral('\'', "character")
		c := Int(0)
		if len(v) != 1 {
			p.scanError("invalid character literal")
		} else {
			c = Int(v[0])
		}
		token = p.NewToken(tkINTEGER, c, pt)
	case p.ScanChar(`"`):
		v := p.scanQuotedLiteral('"', "string")
		token = p.NewToken(tkSTRING, NewStr(string(v)), pt)
	case p.Scan(reInteger):
		v := int64(0)
		if p.Matched[3] != "" {
			v, _ = strconv.ParseInt(p.Matched[3], 0, 64)
		} else if p.Matched[4] != "" {
			v, _ = strconv.ParseInt(p.Matched[4], 0, 64)
		} else if p.Matched[5] != "" {
			v, _ = strconv.ParseInt(p.Matched[5], 10, 64)
		}
		if p.Matched[1] == "-" {
			v = -v
		}
		token = p.NewToken(tkINTEGER, Int(v), pt)
	case p.Scan(reIdentifier):
		token = p.newIdLikeToken(p.Matched[1], p.Matched[2], ws, nl, context, pt, "")
	case p.ScanChar(":"):
		if !ws {
			p.scanError("leading whitespaces are required before the operator `:`")
		}
		token = p.NewToken(tkCL, NIL, pt)
	default:
		p.scanError("invalid character '%c'", p.PeekChar())
		p.ForwardChar()
	}

	p.PushToken(token)
	// fmt.Println("INPUT:", p.Tokens)
}
