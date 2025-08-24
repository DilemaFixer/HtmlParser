package htmlparser

import (
	"fmt"
)

type Scanner struct {
	input    string
	position int
	line     int
	column   int
	ch       rune
}

func NewScanner(input string) *Scanner {
	s := &Scanner{
		input:  input,
		line:   1,
		column: 1,
	}
	s.updateCh()
	return s
}

func (s *Scanner) updateCh() {
	if s.position >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = rune(s.input[s.position])
	}
}

func (s *Scanner) Current() rune {
	return s.ch
}

func (s *Scanner) PeekNext() rune {
	if s.position+1 >= len(s.input) {
		return 0
	}
	return rune(s.input[s.position+1])
}

func (s *Scanner) Previous() rune {
	if s.position <= 0 {
		return 0
	}
	return rune(s.input[s.position-1])
}

func (s *Scanner) Take() rune {
	if s.position >= len(s.input) {
		return 0
	}

	current := s.ch
	s.position++

	if current == '\n' {
		s.line++
		s.column = 1
	} else {
		s.column++
	}

	s.updateCh()
	return current
}

func (s *Scanner) Retreat() rune {
	if s.position <= 0 {
		return 0
	}

	s.position--
	s.updateCh()

	if s.ch == '\n' {
		s.line--
		s.column = s.ColumnAt(s.position)
	} else {
		s.column--
	}

	return s.ch
}

func (s *Scanner) ColumnAt(pos int) int {
	col := 1
	for i := pos - 1; i >= 0 && s.input[i] != '\n'; i-- {
		col++
	}
	return col
}

func (s *Scanner) Peek(offset int) rune {
	pos := s.position + offset
	if pos < 0 || pos >= len(s.input) {
		return 0
	}
	return rune(s.input[pos])
}

func (s *Scanner) Skip() {
	s.Take()
}

func (s *Scanner) SkipN(n int) {
	for i := 0; i < n && !s.EOF(); i++ {
		s.Take()
	}
}

func (s *Scanner) EOF() bool {
	return s.position >= len(s.input)
}

func (s *Scanner) Position() int {
	return s.position
}

func (s *Scanner) Line() int {
	return s.line
}

func (s *Scanner) Column() int {
	return s.column
}

func (s *Scanner) Mark() int {
	return s.position
}

func (s *Scanner) Reset(mark int) {
	if mark < 0 || mark > len(s.input) {
		return
	}

	s.position = mark
	s.line = 1
	s.column = 1

	for i := 0; i < mark; i++ {
		if s.input[i] == '\n' {
			s.line++
			s.column = 1
		} else {
			s.column++
		}
	}

	s.updateCh()
}

func (s *Scanner) Slice(start, end int) string {
	if start < 0 || end > len(s.input) || start > end {
		return ""
	}
	return s.input[start:end]
}

func (s *Scanner) SliceFrom(start int) string {
	return s.Slice(start, s.position)
}

func (s *Scanner) Match(expected rune) bool {
	if s.ch == expected {
		s.Take()
		return true
	}
	return false
}

func (s *Scanner) MatchAny(chars ...rune) bool {
	for _, char := range chars {
		if s.ch == char {
			s.Take()
			return true
		}
	}
	return false
}

func (s *Scanner) MatchString(expected string) bool {
	if s.position+len(expected) > len(s.input) {
		return false
	}

	if s.input[s.position:s.position+len(expected)] == expected {
		for range expected {
			s.Take()
		}
		return true
	}
	return false
}

func (s *Scanner) ConsumeWhile(predicate func(rune) bool) string {
	start := s.position
	for !s.EOF() && predicate(s.ch) {
		s.Take()
	}
	return s.input[start:s.position]
}

func (s *Scanner) ConsumeUntil(predicate func(rune) bool) string {
	start := s.position
	for !s.EOF() && !predicate(s.ch) {
		s.Take()
	}
	return s.input[start:s.position]
}

func (s *Scanner) ConsumeN(n int) string {
	start := s.position
	for i := 0; i < n && !s.EOF(); i++ {
		s.Take()
	}
	return s.input[start:s.position]
}

func (s *Scanner) Find(target rune) bool {
	for !s.EOF() {
		if s.ch == target {
			return true
		}
		s.Take()
	}
	return false
}

func (s *Scanner) FindString(target string) bool {
	targetLen := len(target)
	for s.position+targetLen <= len(s.input) {
		if s.input[s.position:s.position+targetLen] == target {
			return true
		}
		s.Take()
	}
	return false
}

func (s *Scanner) Remaining() string {
	return s.input[s.position:]
}

func (s *Scanner) Len() int {
	return len(s.input)
}

func (s *Scanner) Location() string {
	return fmt.Sprintf("%d:%d", s.line, s.column)
}

func (s *Scanner) SetLocation(line, column int) {
	s.line = line
	s.column = column
	s.updateCh()
}

func (s *Scanner) SkipWhitespace() {
	for !s.EOF() && (s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r') {
		s.Take()
	}
}

func (s *Scanner) ConsumeUntilString(target string) string {
	start := s.position
	targetLen := len(target)
	for s.position+targetLen <= len(s.input) {
		if s.input[s.position:s.position+targetLen] == target {
			break
		}
		s.Take()
	}
	return s.input[start:s.position]
}
