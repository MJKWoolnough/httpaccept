// Package httpaccept provides a function to deal with the Accept header.
package httpaccept // import "vimagination.zapto.org/httpaccept"

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	accept       = "Accept"
	acceptSplit  = ","
	partSplit    = ";"
	weightPrefix = "q="
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
	weight uint16
}

type Mime string

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

// HandlerFunc wraps a func to make it satisfy the Handler interface
type HandlerFunc func(string) bool

// Handle calls the underlying func
func (h HandlerFunc) Handle(e string) bool {
	return h(e)
}

// InvalidAccept writes the 406 header
func InvalidAccept(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotAcceptable)
}

// HandleAccept will process the Accept header and calls the given handler for
// each mime type until the handler returns true.
//
// This function returns true when the Handler returns true, false otherwise
//
// When no Accept header is given the mime string will be the empty string.
func HandleAccept(r *http.Request, h Handler) bool {
	acceptHeader := r.Header.Get(accept)
	accepts := make(mimes, 0, strings.Count(acceptHeader, acceptSplit)+1)
Loop:
	for _, accept := range strings.Split(acceptHeader, acceptSplit) {
		parts := strings.Split(strings.TrimSpace(accept), partSplit)
		name := strings.ToLower(strings.TrimSpace(parts[0]))
		// check mime string format?
		if name == "" {
			continue
		}
		var (
			qVal float64 = 1
			err  error
		)
		for _, part := range parts[1:] {
			if strings.HasPrefix(strings.TrimSpace(part), weightPrefix) {
				qVal, err = strconv.ParseFloat(part[len(weightPrefix):], 32)
				if err != nil || qVal < 0 || qVal >= 2 {
					continue Loop
				}
				break
			}
		}
		accepts = append(accepts, mime{
			mime:   Mime(name),
			weight: uint16(qVal * 1000),
		})
	}
	if len(accepts) == 0 {
		return h.Handle("")
	}
	sort.Stable(accepts)
	for _, accept := range accepts {
		if h.Handle(accept.mime) {
			return true
		}
	}
	return false
}
