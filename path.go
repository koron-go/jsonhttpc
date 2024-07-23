package jsonhttpc

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Params is a parameters for path templating.
type Params map[string]interface{}

// Path provides path templating.
//
// "{foo}" in `s` are replaced with `fmt.Sprint(p["foo"])`.
// '\' escapes a rune succeeding.
//
// This returns a string with prefix "jsonhttpc.Parse error: " if detects some
// errors.
func Path(s string, p Params) string {
	s2, err := path(s, p)
	if err != nil {
		return fmt.Sprintf("jsonhttpc.Parse error: %s", err)
	}
	return s2
}

func path(s string, p Params) (string, error) {
	b := &strings.Builder{}
	for s != "" {
		n := strings.IndexAny(s, "\\{")
		if n < 0 {
			b.WriteString(s)
			break
		}
		if n > 0 {
			b.WriteString(s[:n])
			s = s[n:]
		}
		switch s[0] {
		case '\\':
			if len(s) <= 1 {
				return "", errors.New("no chars to escape")
			}
			r, m := utf8.DecodeRuneInString(s[1:])
			if r == utf8.RuneError {
				return "", errors.New("rune error to escape")
			}
			b.WriteRune(r)
			s = s[m+1:]
		case '{':
			m := strings.IndexRune(s, '}')
			if m < 0 {
				return "", errors.New("not found '}'")
			}
			k := s[1:m]
			v, ok := p[k]
			if !ok {
				return "", fmt.Errorf("not found key: key=%q", k)
			}
			b.WriteString(fmt.Sprint(v))
			s = s[m+1:]
		}
	}
	return b.String(), nil
}
