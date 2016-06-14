package parser

import (
	"reflect"
	"strings"
)

// A structLexer lexes over the tags of struct fields while tracking the current field.
type structLexer struct {
	s     reflect.Type
	field int
	lexer Lexer
	r     *strings.Reader
}

func newStructLexer(s reflect.Type) *structLexer {
	return &structLexer{
		s:     s,
		lexer: LexString(string(s.Field(0).Tag)),
	}
}

func (s *structLexer) NumField() int {
	return s.s.NumField()
}

func (s *structLexer) Field() reflect.StructField {
	return s.s.Field(s.field)
}

func (s *structLexer) Peek() Token {
	field := s.field
	lexer := s.lexer
	for {
		token := lexer.Peek()
		if !token.EOF() {
			return token
		}
		field++
		if field >= s.s.NumField() {
			return EOFToken
		}
		lexer = LexString(string(s.s.Field(field).Tag))
	}
}

func (s *structLexer) Next() Token {
	token := s.lexer.Next()
	if !token.EOF() {
		return token
	}
	if s.field+1 >= s.s.NumField() {
		return EOFToken
	}
	s.field++
	s.lexer = LexString(string(s.s.Field(s.field).Tag))
	return s.Next()
}

func (s *structLexer) Position() Position {
	pos := s.lexer.Position()
	pos.Line = s.field + 1
	return pos
}

func (s *structLexer) Symbols() map[rune]string {
	return s.lexer.Symbols()
}
