package core

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unsafe"
)

type Pt struct {
	Pos  int32
	Line int32
}

type Token struct {
	Pt
	From  *Scanner
	Kind  int32
	Value Value
}

func (v *Token) String() string {
	label := tokenLabels[v.Kind]
	return fmt.Sprintf("<Token: %s pos: %d>", label, v.Pos)
}

func (v *Token) Inspect() string {
	return v.String()
}

func (v *Token) Dup() Value {
	return CopyPtr(v)
}

func (v *Token) Column() int32 {
	return v.Pos - v.From.Lines[v.Line]
}

func (v *Token) PtPrefix() string {
	for v.From == InternalParser {
		if id, ok := v.Value.(*Identifier); ok && id.ExpandedBy != nil {
			v = id.ExpandedBy.Token
		} else {
			break
		}
	}
	if v.From == InternalParser {
		return ""
	}
	return fmt.Sprintf("%s:%d:%d: ", v.From.Path, v.Line+1, v.Column())
}

func (v *Token) FormatAsErrorLine(s string) string {
	if v.From == InternalParser {
		s = fmt.Sprintf("  %s internal %s %s\n", s, TypeLabelOf(v.Value), v.Value)
	} else {
		col := v.Column()
		s = fmt.Sprintf("  %s %s:%d:%d\n", s, v.From.Path, v.Line+1, col) +
			fmt.Sprintf("   |%s\n", string(v.From.LineBytes(v.Line, true))) +
			fmt.Sprintf("   |%s^-- ??\n", strings.Repeat(" ", int(col)))
	}
	return s
}

type Scanner struct {
	Pt
	Path    string
	Text    []byte
	Lines   []int32
	Matched []string
	Tokens  []*Token
	OnError func(*InternalError)
	cc      *Compiler
}

func (s *Scanner) Init() {
	s.Lines = []int32{0}
	for x, c := range s.Text {
		if c == '\n' {
			s.Lines = append(s.Lines, int32(x+1))
		}
	}
	s.Lines = append(s.Lines, int32(len(s.Text)))
}

func (s *Scanner) LineBytes(n int32, chomp bool) []byte {
	a := s.Lines[n]
	b := s.Lines[n+1]
	if chomp && b > a && s.Text[b-1] == '\n' { // a == b if len(line) == 0
		b--
	}
	return s.Text[a:b]
}

var InternalParser = &Scanner{Path: "<internal>"}
var ScanErrorToken = &Token{From: InternalParser, Kind: tkSCANERROR, Value: IdUNDER}

func _bstos(a []byte) string {
	if len(a) == 0 {
		return ""
	}
	return unsafe.String(&a[0], len(a))
}

func (s *Scanner) scanError(message string, args ...any) {
	s.scanErrorAt(s.Pt, message, args...)
}

func (s *Scanner) scanErrorAt(pt Pt, message string, args ...any) {
	err := &InternalError{
		tag:     "scan error: ",
		message: fmt.Sprintf(message, args...),
		at:      append([]Value{s.NewToken(0, NIL, pt)}, s.cc.nested...),
	}
	if s.OnError != nil {
		s.OnError(err)
	} else {
		raiseError(err)
	}
}

func (s *Scanner) IsEOF() bool {
	return !s.IsValidPos(s.Pos)
}

func (s *Scanner) IsValidPos(pos int32) bool {
	return int(pos) < len(s.Text)
}

func (s *Scanner) PeekChar() byte {
	return s.Text[s.Pos]
}

func (s *Scanner) ForwardChar() {
	if s.PeekChar() == '\n' {
		s.Line++
	}
	s.Pos++
}

func (s *Scanner) GetChar() byte {
	c := s.PeekChar()
	s.ForwardChar()
	return c
}

func (s *Scanner) SkipChars(n int32) {
	for range n {
		s.ForwardChar()
	}
}

func (s *Scanner) Scan(re *regexp.Regexp) bool {
	s.Matched = re.FindStringSubmatch(_bstos(s.Text[s.Pos:]))
	if s.Matched != nil {
		s.SkipChars(int32(len(s.Matched[0])))
		return true
	}
	return false
}

func (s *Scanner) ScanChar(chars string) bool {
	ok := s.MatchChar(chars)
	if ok {
		s.ForwardChar()
	}
	return ok
}

func (s *Scanner) MatchChar(chars string) bool {
	return s.MatchCharAt(chars, s.Pos)
}

func (s *Scanner) MatchCharAt(chars string, pos int32) bool {
	return s.IsValidPos(pos) && strings.IndexByte(chars, s.Text[pos]) > -1
}

func (s *Scanner) ScanSeq(seq string) bool {
	if s.MatchSeq(seq) {
		s.SkipChars(int32(len(seq)))
		return true
	}
	return false
}

func (s *Scanner) MatchSeq(seq string) bool {
	return bytes.HasPrefix(s.Text[s.Pos:], []byte(seq))
}

func (s *Scanner) SubStringFrom(pos int32) string {
	return _bstos(s.Text[pos:s.Pos])
}

func (s *Scanner) NewToken(kind int32, value Value, pt Pt) *Token {
	return &Token{
		From:  s,
		Kind:  kind,
		Value: value,
		Pt:    pt,
	}
}

func (s *Scanner) NewIdToken(kind int32, id *Identifier, pt Pt) *Token {
	id.Token = s.NewToken(kind, id, pt)
	return id.Token
}

func (s *Scanner) PushToken(token *Token) {
	s.Tokens = append(s.Tokens, token)
}

func (p *Scanner) ConsumeToken() *Token {
	token := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	return token
}

func (p *Scanner) CancelLastTokenIf(b bool) {
	if b {
		p.Tokens = p.Tokens[:len(p.Tokens)-1]
	}
}
