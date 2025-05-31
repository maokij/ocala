package core

import (
	"regexp"
	"strconv"
)

type Parser struct {
	Scanner
	state    int
	contexts []byte
}

const (
	pstDefault = iota
	pstWs
	pstNoWs
	pstNl
	pstNoNl
)

func (p *Parser) RaiseParseError(token *Token, expected string) {
	label := tokenLabels[token.Kind]
	p.cc.g.raiseError(token, "parse error: ", "unexpected %s, expected %s\n", label, expected)
}

func (p *Parser) PeekToken() *Token {
	if len(p.Tokens) == 0 {
		p.scanToken()
	}
	return p.Tokens[0]
}

func (p *Parser) seekToNextToken(nl bool) (int32, bool, bool) {
	pos, line, bline, comma, comment := p.Pos, p.Line, int32(-1), false, 0

loop:
	for ; !p.IsEOF(); p.ForwardChar() {
		c := p.PeekChar()
		switch {
		case c == ' ', c == '\t', c == '\r':
			// NOP
		case c == ',' && comment <= 0:
			if comma || p.Line != line {
				p.RaiseScanError("invalid comma")
			}
			comma = true
		case c == '\n':
			p.Line++
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
		p.RaiseScanError("the block comment is not terminated")
	} else if nl {
		nl = !comma && line != p.Line
		if comma && p.Line-line > 1 {
			p.RaiseScanError("the comma followed by blank lines is not allowed")
		} else if bline > -1 {
			p.RaiseScanError("the block comment must be followed by new line")
		}
	}

	return p.Pos, pos != p.Pos, nl
}

var rePlaceholder = regexp.MustCompile(`^%+(=|#|&|>*[<>*])`)
var reInteger = regexp.MustCompile(`^([+-])?((0x[0-9a-fA-F][0-9a-fA-F_]*)|(0b[01][01_]*)|([0-9][0-9_]*))`)
var reIdentifier = regexp.MustCompile(
	`^([_a-zA-Z][-^!$%&*+/<=>?|~_a-zA-Z0-9]*:|::)?([-^!$%&*+/<=>?|~_a-zA-Z][-^!$%&*+/<=>?|~_a-zA-Z0-9]*)`)
var nameChars = `-^!$%&*+/<=>?|~_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
var wsChars = " \t\r\n"
var tailChars = " \t\r\n,;)]}"

func (p *Parser) newIdLikeToken(s, t string, ws, nl bool, context byte, pos int32, ph string) *Token {
	if u, ok := p.cc.TokenAliases[t]; ok {
		t = u
	}
	k := Intern(t)
	v := &Identifier{Name: k, PlaceHolder: ph}
	n := p.cc.Operators[t]

	ns := s != ""
	if ns {
		if ph != "" {
			p.RaiseScanError("placeholders cannot contain namespaces")
		}
		v.Namespace = Intern(s[:len(s)-1])
	}

	if tk, ok := tokenKinds[t]; ok { // =
		return p.NewIdToken(tk, v, pos)
	} else if n&1 != 0 && !ns && ws && p.MatchChar(tailChars) {
		p.CancelLastTokenIf(nl)
		return p.NewIdToken(tkUOP, v, pos)
	} else if n&2 != 0 && !ns && ws && p.MatchChar(wsChars) {
		p.CancelLastTokenIf(nl)
		// p.state = pstNoNl
		return p.NewIdToken(tkBOP, v, pos)
	} else if tk, ok := p.cc.ReservedWords[t]; ok {
		if ph != "" {
			p.RaiseScanError("reserved words cannot be use as placeholders")
		}

		token := p.NewIdToken(tk, v, pos)
		if tk == tkREG && p.ScanChar("@") {
			p.PushToken(token)
			token = p.NewIdToken(tkMIAT, v, p.Pos-1)
		} else if tk == tkCOND && p.ScanChar(".") {
			token.Kind = tkCONDDOT
			p.state = pstWs
		}
		return token
	} else if p.MatchChar("(") {
		return p.NewIdToken(tkIDENTIFIERP, v, pos)
	} else if p.MatchChar(":") && !ns {
		return p.NewIdToken(tkLABEL, v, pos)
	}
	return p.NewIdToken(tkIDENTIFIER, v, pos)
}

func (p *Parser) findTokenKind(s string) int32 {
	tk, ok := tokenKinds[s]
	if !ok {
		p.RaiseScanError("invalid token %s", s)
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

func (p *Parser) scanEscapedChar() byte {
	c := p.GetCharOrError()
	if c, ok := escapeChars[c]; ok {
		return c
	}
	if c == 'x' {
		s := []byte{p.GetCharOrError(), p.GetCharOrError()}
		if n, err := strconv.ParseUint(string(s), 16, 8); err == nil {
			return byte(n)
		}
	}
	p.RaiseScanError("invalid character escape")
	return 0
}

func (p *Parser) scanToken() {
	context := p.contexts[len(p.contexts)-1]
	pos, ws, nl := p.seekToNextToken(context == '{')
	state := p.state
	if state > 0 {
		p.state = 0
	}
	ws = (state != pstNoWs) && (ws || state == pstWs)
	nl = (state != pstNoNl) && (nl || state == pstNl)

	if p.IsEOF() {
		p.PushToken(p.NewToken(tkEOF, NIL, p.Pos))
		return
	}

	if nl {
		p.PushToken(p.NewToken(tkSC, NIL, pos-1))
	}

	// !"#$%&'()*+,-./:;<=>?@[\]^_`{|}~ // -^!$%&*+/<=>?@|~
	var token *Token
	switch {
	case p.Scan(rePlaceholder):
		ph := p.Matched[0]
		if !p.Scan(reIdentifier) {
			p.RaiseScanError("invalid placeholder")
		}
		token = p.newIdLikeToken(p.Matched[1], p.Matched[2], false, false, context, pos, ph)
	case p.ScanSeq("..."):
		token = p.NewToken(tkREST, Int(1), pos)
	case p.ScanSeq(".."):
		token = p.NewIdToken(tkBOP, &Identifier{Name: KwDotDot}, pos)
	case p.ScanChar("."):
		p.CancelLastTokenIf(nl)
		tk := int32(tkDTMI)
		if p.MatchChar(wsChars) {
			if !ws {
				p.RaiseScanError("leading whitespaces are required before the operator `.`")
			}
			tk = tkDOP
		}
		token = p.NewIdToken(tk, &Identifier{Name: KwDot}, pos)
	case p.ScanChar("@"):
		s := p.SubStringFrom(pos) + stringIf(!p.MatchChar(wsChars), "-")
		tk := p.findTokenKind(s)
		token = p.NewToken(tk, NIL, pos)
	case p.ScanSeq("$"):
		p.ScanSeq("$")
		if p.MatchChar("(") {
			s := p.SubStringFrom(pos) + "-"
			tk := p.findTokenKind(s)
			token = p.NewToken(tk, NIL, pos)
			break
		}
		p.Pos = pos
		p.Scan(reIdentifier)
		token = p.newIdLikeToken(p.Matched[1], p.Matched[2], ws, nl, context, pos, "")
	case p.ScanChar("!"):
		here := p.Pos
		if p.ScanChar("=") {
			p.ScanChar("?")
			if !p.MatchChar(nameChars) {
				s := p.SubStringFrom(pos)
				token = p.newIdLikeToken("", s, ws, nl, context, pos, "")
				break
			}
			p.Pos = here
		}
		if p.MatchChar(tailChars) {
			token = p.NewToken(tkEX, NIL, pos)
			break
		}
		token = p.NewIdToken(tkAOP, &Identifier{Name: KwLogicalNotOp}, pos)
	case p.ScanChar("~"):
		if p.MatchChar(wsChars) {
			p.RaiseScanError("no whitespace is allowed after the prefix operator")
		}
		token = p.NewIdToken(tkAOP, &Identifier{Name: KwNotOp}, pos)
	case p.ScanChar(";`"):
		tk := p.findTokenKind(p.SubStringFrom(pos))
		token = p.NewToken(tk, NIL, pos)
	case p.ScanSeq("={"):
		token = p.NewToken(tkEQLC, NIL, pos)
		p.contexts = append(p.contexts, '{')
	case p.ScanSeq("%{"):
		token = p.NewToken(tkPELC, NIL, pos)
		p.contexts = append(p.contexts, '{')
	case p.ScanChar("([{"):
		s := p.SubStringFrom(pos)
		tk := p.findTokenKind(s)
		token = p.NewToken(tk, NIL, pos)
		p.contexts = append(p.contexts, s[0])
	case p.ScanChar(")]}"):
		tk := p.findTokenKind(p.SubStringFrom(pos))
		token = p.NewToken(tk, NIL, pos)
		p.contexts = p.contexts[:len(p.contexts)-1]
	case p.ScanChar(`'`):
		c := p.GetCharOrError()
		if c == '\'' {
			p.RaiseScanError("blank character literal is invalid")
		} else if c == '\\' {
			c = p.scanEscapedChar()
		}
		if p.GetCharOrError() != '\'' {
			p.RaiseScanError("invalid character literal")
		}
		token = p.NewToken(tkINTEGER, Int(c), pos)
	case p.ScanChar(`"`):
		v := []byte{}
		for {
			c := p.GetCharOrError()
			if c == '"' {
				break
			} else if c == '\\' {
				c = p.scanEscapedChar()
			}
			v = append(v, c)
		}
		token = p.NewToken(tkSTRING, NewStr(string(v)), pos)
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
		token = p.NewToken(tkINTEGER, Int(v), pos)
	case p.Scan(reIdentifier):
		token = p.newIdLikeToken(p.Matched[1], p.Matched[2], ws, nl, context, pos, "")
	case p.ScanChar(":"):
		token = p.NewToken(tkCL, NIL, pos)
	default:
		p.RaiseScanError("invalid character '%c'", p.PeekChar())
	}

	// for _, i := range p.Tokens {
	// 	fmt.Printf("%s, ", tokenLabels[i.Kind])
	// }
	// fmt.Println(tokenLabels[token.Kind])
	p.PushToken(token)
}
