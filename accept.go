// Package httpaccept provides a function to deal with the Accept header.
package httpaccept // import "vimagination.zapto.org/httpaccept"

import (
	"net/http"
	"slices"
	"sort"
	"strings"

	"vimagination.zapto.org/parser"
)

const (
	wcAny    = "*"
	matchAny = "*/*"
	accept   = "Accept"
)

type mimes []mime

func (m mimes) Len() int {
	return len(m)
}

func (m mimes) Less(i, j int) bool {
	return m[j].weight < m[i].weight
}

func (m mimes) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type mime struct {
	mime   Mime
	weight int16
}

// Mime represents a accepted Mime Type.
type Mime string

// Match checks to see whether a given Mime Type matches the value.
//
// The method allows for wildcards in the subtype sections.
func (m Mime) Match(n Mime) bool {
	mMime, mExcl, _ := strings.Cut(string(m), ";")
	nMime, nExcl, _ := strings.Cut(string(n), ";")
	mPrefix, mSuffix, _ := strings.Cut(mMime, "/")
	nPrefix, nSuffix, _ := strings.Cut(nMime, "/")

	if mPrefix != wcAny && nPrefix != wcAny && !strings.EqualFold(mPrefix, nPrefix) {
		return false
	}

	if mSuffix != wcAny && nSuffix != wcAny && !strings.EqualFold(mSuffix, nSuffix) {
		return false
	}

	if mExcl != "" {
		for mEx := range strings.SplitSeq(mExcl, ";") {
			if Mime(nMime).Match(Mime(mEx)) {
				return false
			}
		}
	}

	if nExcl != "" {
		for nEx := range strings.SplitSeq(nExcl, ";") {
			if Mime(mMime).Match(Mime(nEx)) {
				return false
			}
		}
	}

	return true
}

// Handler provides an interface to handle a mime type.
//
// The mime string (e.g. text/html, application/json, text/plain) is passed to
// the handler, which is expected to return true if no more encodings are
// required and false otherwise.
//
// The empty string "" is used to signify when no preference is specified.
type Handler interface {
	Handle(mime Mime) bool
}

// HandlerFunc wraps a func to make it satisfy the Handler interface.
type HandlerFunc func(Mime) bool

// Handle calls the underlying func.
func (h HandlerFunc) Handle(m Mime) bool {
	return h(m)
}

// InvalidAccept writes the 406 header.
func InvalidAccept(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
}

// HandleAccept will process the Accept header and calls the given handler for
// each mime type until the handler returns true.
//
// This function returns true when the Handler returns true, false otherwise.
//
// When no Accept header is given the mime string will be the empty string.
func HandleAccept(r *http.Request, h Handler) bool {
	accepts := parseAccepts(r.Header.Get(accept))

	if len(accepts) == 0 {
		return h.Handle("")
	}

	sort.Stable(accepts)

	for _, accept := range accepts {
		if accept.weight > 0 && h.Handle(accept.mime) {
			return true
		}
	}

	return false
}

func parseAccepts(acceptHeader string) mimes {
	accepts := make(mimes, 0, strings.Count(acceptHeader, delim)+1)

	p := parseAccept(acceptHeader)

	for {
		coding := p.Next()
		if coding.Type == parser.TokenDone {
			break
		}

		name := coding.Data

		if p.Accept(tokenInvalidWeight) {
			continue
		}

		weight := int16(1000)

		if p.Peek().Type == tokenWeight {
			weight = parseQ(p.Next().Data)
		}

		if slices.ContainsFunc(accepts, func(e mime) bool { return e.mime == Mime(name) }) {
			continue
		}

		accepts = append(accepts, mime{mime: Mime(name), weight: weight})
	}

	return processAnys(accepts)
}

func processAnys(accepts mimes) mimes {
	for n := range accepts {
		if accepts[n].weight == 0 {
			continue
		}

		if prefix, suffix, _ := strings.Cut(string(accepts[n].mime), "/"); prefix != wcAny && suffix != wcAny {
			continue
		}

		var nots strings.Builder

		nots.WriteString(string(accepts[n].mime))

		for _, m := range accepts {
			if m.weight > 0 {
				continue
			}

			if accepts[n].mime.Match(m.mime) {
				nots.WriteByte(';')
				nots.WriteString(string(m.mime))
			}
		}

		if nots.Len() > len(accepts[n].mime) {
			accepts[n].mime = Mime(nots.String())
		}
	}

	return accepts
}

var multiplies = [...]int16{100, 10, 1}

func parseQ(q string) int16 {
	if q[0] == '1' {
		return 1000
	}

	if len(q) < 2 {
		return 0
	}

	var qv int16

	for n, v := range q[2:] {
		qv += int16(v-'0') * multiplies[n]
	}

	return qv
}
