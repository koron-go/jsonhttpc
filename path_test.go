package jsonhttpc

import "testing"

func TestParse(t *testing.T) {
	for ti, tc := range []struct {
		s   string
		p   Params
		exp string
	}{
		// basic cases
		{"", nil, ""},
		{"/foo/bar/baz", nil, "/foo/bar/baz"},
		// with parameters
		{"/usr/{userID}/items", Params{"userID": 123}, "/usr/123/items"},
		{"/usr/{userID}/items/{itemID}/price", Params{
			"userID": 123,
			"itemID": "xyz999",
		}, "/usr/123/items/xyz999/price"},
		// escape
		{`\\`, nil, `\`},
		{`\a`, nil, `a`},
		{`\あ`, nil, `あ`},
		{`prefix\\suffix`, nil, `prefix\suffix`},
		{`prefix\asuffix`, nil, `prefixasuffix`},
		{`prefix\あsuffix`, nil, `prefixあsuffix`},
		{`\{`, nil, `{`},
		{`prefix\{suffix`, nil, `prefix{suffix`},
		// failure
		{"/usr/{userID}/items", nil,
			"jsonhttpc.Parse error: not found key: key=\"userID\""},
		{"/usr/{userID/items", Params{"userID": 123},
			"jsonhttpc.Parse error: not found '}'"},
		{`\`, nil, `jsonhttpc.Parse error: no chars to escape`},
		{`abc\`, nil, `jsonhttpc.Parse error: no chars to escape`},
		{"\\\xe8", nil, "jsonhttpc.Parse error: rune error to escape"},
		{"abc\\\xe8", nil, "jsonhttpc.Parse error: rune error to escape"},
	} {
		act := Path(tc.s, tc.p)
		if act != tc.exp {
			t.Fatalf("parse failed #%d\nexpect=%s\nactual=%s", ti, tc.exp, act)
		}
	}
}
