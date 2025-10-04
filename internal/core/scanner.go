package core

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unsafe"
)

type Token struct {
	From  *Scanner
	Kind  int32
	Pos   int32
	Value Value
}

func (v *Token) String() string {
	label := tokenLabels[v.Kind]
	return fmt.Sprintf("<Token: %s pos: %d>", label, v.Pos)
}

func (v *Token) LineNumber() int {
	return bytes.Count(v.From.Text[:v.Pos], []byte{'\n'})
}

func (v *Token) Inspect() string {
	return v.String()
}

func (v *Token) Dup() Value {
	return CopyPtr(v)
}

func (v *Token) FormatAsErrorLine(s string) string {
	p := v.From
	bol := p.SeekToBOL(v.Pos)
	col := max(v.Pos-bol, 0)
	p.SkipUntil('\n')

	if v.From == InternalParser {
		s = fmt.Sprintf("  %s internal %s %s\n", s, TypeLabelOf(v.Value), v.Value)
	} else {
		s = fmt.Sprintf("  %s %s:%d:%d\n", s, p.Path, v.LineNumber()+1, col) +
			fmt.Sprintf("   |%s\n", p.SliceFrom(bol)) +
			fmt.Sprintf("   |%s^-- ??\n", strings.Repeat(" ", int(col)))
	}
	return s
}

type Scanner struct {
	Path    string
	Text    []byte
	Pos     int32
	Line    int32
	Matched []string
	Tokens  []*Token
	cc      *Compiler
}

var InternalParser = &Scanner{Path: "<internal>"}

func _bstos(a []byte) string {
	if len(a) == 0 {
		return ""
	}
	return unsafe.String(&a[0], len(a))
}

func (s *Scanner) ErrorWith(message string, args ...any) {
	err := &InternalError{
		tag: "scan error: ",
		at:  append([]Value{s.NewToken(0, NIL, s.Pos)}, s.cc.nested...),
		g:   s.cc.g,
	}
	err.With(message, args...)
}

func (s *Scanner) IsEOF() bool {
	return !s.IsValidPos(s.Pos)
}

func (s *Scanner) IsValidPos(pos int32) bool {
	return int(pos) < len(s.Text)
}

func (s *Scanner) SkipUntil(c byte) bool {
	for ; !s.IsEOF(); s.ForwardChar() {
		if s.PeekChar() == c {
			return true
		}
	}
	return false
}

func (s *Scanner) SeekToBOL(pos int32) int32 {
	for s.Pos = pos; s.Pos > 0; s.Pos-- {
		if s.Text[s.Pos-1] == '\n' {
			break
		}
	}
	return s.Pos
}

func (s *Scanner) PeekChar() byte {
	return s.Text[s.Pos]
}

func (s *Scanner) ForwardChar() {
	s.Pos++
}

func (s *Scanner) GetChar() (byte, bool) {
	if s.IsEOF() {
		return 0, false
	}
	c := s.PeekChar()
	s.ForwardChar()
	return c, true
}

func (s *Scanner) GetCharOrError() byte {
	c, ok := s.GetChar()
	if !ok {
		s.ErrorWith("unexpected EOF")
	}
	return c
}

func (s *Scanner) Scan(re *regexp.Regexp) bool {
	s.Matched = re.FindStringSubmatch(_bstos(s.Text[s.Pos:]))
	if s.Matched != nil {
		s.Pos += int32(len(s.Matched[0]))
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
	return !s.IsEOF() && strings.IndexByte(chars, s.PeekChar()) > -1
}

func (s *Scanner) ScanSeq(seq string) bool {
	if s.MatchSeq(seq) {
		s.Pos += int32(len(seq))
		return true
	}
	return false
}

func (s *Scanner) MatchSeq(seq string) bool {
	return bytes.HasPrefix(s.Text[s.Pos:], []byte(seq))
}

func (s *Scanner) SliceFrom(pos int32) []byte {
	return s.Text[pos:s.Pos]
}

func (s *Scanner) SubStringFrom(pos int32) string {
	return _bstos(s.Text[pos:s.Pos])
}

func (s *Scanner) NewToken(kind int32, value Value, pos int32) *Token {
	return &Token{
		From:  s,
		Kind:  kind,
		Value: value,
		Pos:   pos,
	}
}

func (s *Scanner) NewIdToken(kind int32, id *Identifier, pos int32) *Token {
	id.Token = s.NewToken(kind, id, pos)
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
