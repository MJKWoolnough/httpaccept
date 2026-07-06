# httpaccept

[![CI](https://github.com/MJKWoolnough/httpaccept/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/httpaccept/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/httpaccept.svg)](https://pkg.go.dev/vimagination.zapto.org/httpaccept)

--
    import "vimagination.zapto.org/httpaccept"

Package httpaccept provides a function to deal with the Accept header.

## Highlights

 - Simple handling of `Accept` HTTP header.
 - Supports wildcards, and q-values.

## Usage

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"vimagination.zapto.org/httpaccept"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if !httpaccept.HandleAccept(r, httpaccept.HandlerFunc(func(m httpaccept.Mime) bool {
		if m.Match("image/png") {
			io.WriteString(w, "png")
		} else if m.Match("image/*") {
			io.WriteString(w, "image")
		} else if m.Match("*/other") {
			io.WriteString(w, "other")
		} else {
			return false
		}

		return true
	})) {
		io.WriteString(w, "none")
	}
}

func Example() {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept", "image/png")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept", "image/*")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept", "image/gif")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept", "text/other")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept", "image/png, text/other")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept", "image/png;q=0.5, text/other")
	handler(w, r)
	fmt.Println(w.Body)

	w = httptest.NewRecorder()
	r.Header.Set("Accept", "*/*;q=0")
	handler(w, r)
	fmt.Println(w.Body)

	// Output:
	// png
	// png
	// image
	// other
	// png
	// other
	// none
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/httpaccept
