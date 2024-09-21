package token

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Slice []T

func (t Slice) Equal(o Slice) bool {
	if len(t) != len(o) {
		return false
	}
	for i, tok := range t {
		if tok != o[i] {
			return false
		}
	}
	return true
}

type T struct {
	Kind  Kind
	Start int
	End   int
}

func (t T) Get(in []byte) []byte {
	return in[t.Start:t.End]
}

func (t T) GetWrap(in []byte, left, right byte) []byte {
	out := make([]byte, 0, len(in)+2)
	out = append(out, left)
	out = append(out, in[t.Start:t.End]...)
	out = append(out, right)
	return out
}

type Kind int

const (
	KindUnset = iota
	KindTitle
	KindText
	KindVar
	KindRawBlock
)

func (k Kind) String() string {
	switch k {
	case KindTitle:
		return "title"
	case KindText:
		return "text"
	case KindVar:
		return "var"
	case KindUnset:
		return "unset"
	}
	return "unknown"
}

var (
	cmdTitle = []byte(".title")
)

func Tokenize(input []byte) (Slice, error) {
	var tokens []T
	var ct T
	startLine := true
	cmdStart := -1
	cmdEnd := cmdStart
	// lookingForEnd := false
	for i, b := range input {
		if startLine {
			startLine = false
			switch b {
			case '~':
			case '.':
				cmdStart = i
				if ct.Kind != KindUnset {
					// trim newlines before the directive
					sub := 1
					for input[i-sub] == '\n' {
						sub++
					}
					ct.End = i - sub + 1
					tokens = append(tokens, ct)
					log.Trace().Msgf("command start: %+v", ct)
				} else {
					log.Trace().Msgf("kind was previously set: %+v", ct)
				}
			case '<':
			case '\n':
			default:
			}
		}

		if cmdStart != -1 {
			switch b {
			case ' ':
				cmdEnd = i
				switch {
				case bytes.Equal(input[cmdStart:cmdEnd], cmdTitle):
					ct.Kind = KindTitle
					ct.Start = i + 1
				}
			}
		}

		if ct.Kind == KindVar && b == '>' {
			ct.End = i
			tokens = append(tokens, ct)
			ct = T{}
			continue
		}
		if b == '<' && len(input) > i && input[i+1] == '!' {
			if ct.Kind != KindUnset {
				ct.End = i
				tokens = append(tokens, ct)
			}
			ct.Kind = KindVar
			ct.Start = i + 2

		}

		if b == '\n' {
			startLine = true
			switch ct.Kind {
			case KindText:
				// lookingForEnd = true
				// log.Trace().Msgf("end of input, adding token: %+v", ct)

			case KindTitle:
				ct.End = i
				tokens = append(tokens, ct)
				ct = T{}
				continue
			case KindUnset:
				if len(input) > i {
					switch input[i+1] {
					case '.':
						continue
					}
				}
			}
		}

		if ct.Kind == KindUnset {
			ct.Kind = KindText
			ct.Start = i
			// log.Trace().Msgf("we are adding a text node: %+v", ct)
		} else if i == len(input)-1 && ct.Kind != KindUnset {
			ct.End = i + 1
			log.Trace().Msgf("end of input, adding token: %+v", ct)
			tokens = append(tokens, ct)
		}
	}
	return tokens, nil
}

func (t T) Stringify(in []byte) string {
	var buf strings.Builder
	buf.Grow(len(in) + 100)
	buf.WriteString(t.Kind.String())
	buf.WriteString(":\t[")
	buf.WriteString(strconv.Itoa(t.Start))
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(t.End))
	buf.WriteString("]\t")
	if t.Start < 0 || t.End > len(in) || t.Start > t.End {
		buf.WriteString("INVALID BOUNDS")
	} else {
		buf.WriteString("\"")
		buf.Write(in[t.Start:t.End])
		buf.WriteString("\"")
	}
	return buf.String()
}

func (s Slice) Stringify(in []byte) string {
	if len(in) == 0 {
		return "no content"
	}
	if len(s) == 0 {
		return "no tokens"
	}
	var buf strings.Builder
	buf.Grow(len(in) + 100)
	for i, tok := range s {
		buf.WriteString(tok.Stringify(in))
		if i != len(s)-1 {
			buf.WriteRune('\n')
		}
	}
	return buf.String()
}
