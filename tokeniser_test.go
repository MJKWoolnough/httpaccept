package httpaccept

import (
	"testing"

	"vimagination.zapto.org/parser"
)

func TestTokeniser(t *testing.T) {
	for n, test := range [...]struct {
		Accept string
		Tokens []parser.Token
	}{
		{
			"",
			[]parser.Token{
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/b",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/b"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/b,*/*",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/b"},
				{Type: tokenMedia, Data: "*/*"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"*/a\t , \tb/*",
			[]parser.Token{
				{Type: tokenMedia, Data: "*/a"},
				{Type: tokenMedia, Data: "b/*"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a;q=1",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenWeight, Data: "1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=1",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenWeight, Data: "1"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=1.",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenWeight, Data: "1."},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=1.0",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenWeight, Data: "1.0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=1.000",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenWeight, Data: "1.000"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=1.0000",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=1.100",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=2",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenInvalidWeight, Data: "2"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=2",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenInvalidWeight, Data: "2"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=0",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenWeight, Data: "0"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=0.123",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenWeight, Data: "0.123"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=0.1234",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"a/a ; q=0.a",
			[]parser.Token{
				{Type: tokenMedia, Data: "a/a"},
				{Type: tokenInvalidWeight, Data: ""},
				{Type: parser.TokenDone, Data: ""},
			},
		},
		{
			"*/*;q=0,abc/def;q=0.512",
			[]parser.Token{
				{Type: tokenMedia, Data: "*/*"},
				{Type: tokenWeight, Data: "0"},
				{Type: tokenMedia, Data: "abc/def"},
				{Type: tokenWeight, Data: "0.512"},
				{Type: parser.TokenDone, Data: ""},
			},
		},
	} {
		p := parseAccept(test.Accept)

		for m, tkn := range test.Tokens {
			if tk, _ := p.GetToken(); tk.Type != tkn.Type {
				if tk.Type == parser.TokenError {
					t.Errorf("test %d.%d: unexpected error: %s", n+1, m+1, tk.Data)
				} else {
					t.Errorf("test %d.%d: Incorrect type, expecting %d, got %d", n+1, m+1, tkn.Type, tk.Type)
				}

				break
			} else if tk.Data != tkn.Data {
				t.Errorf("test %d.%d: Incorrect data, expecting %q, got %q", n+1, m+1, tkn.Data, tk.Data)

				break
			}
		}
	}
}
