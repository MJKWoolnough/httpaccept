package httpaccept_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"vimagination.zapto.org/httpaccept"
)

func Example() {
	handler := func(w http.ResponseWriter, r *http.Request) {
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
