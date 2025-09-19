package httpaccept

import (
	"net/http"
	"reflect"
	"testing"
)

type testMimes []string

func (t *testMimes) Handle(mime Mime) bool {
	*t = append(*t, string(mime))

	return false
}

func TestOrder(t *testing.T) {
	for n, test := range []struct {
		AcceptEncoding string
		Encodings      testMimes
	}{
		{"", testMimes{""}},
		{"a/b", testMimes{"a/b"}},
		{"a/b, c/d, e/f", testMimes{"a/b", "c/d", "e/f"}},
		{"a/b;q=0, c/d;q=0.5, e/f;q=1", testMimes{"e/f", "c/d"}},
	} {
		te := make(testMimes, 0, len(test.Encodings))

		HandleAccept(&http.Request{
			Header: http.Header{
				accept: []string{test.AcceptEncoding},
			},
		}, &te)

		if !reflect.DeepEqual(te, test.Encodings) {
			t.Errorf("test %d: expecting %v, got %v", n+1, test.Encodings, te)
		}
	}
}
