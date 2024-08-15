# httpaccept
--
    import "vimagination.zapto.org/httpaccept"

Package httpaccept provides a function to deal with the Accept header.

## Usage

#### func  HandleAccept

```go
func HandleAccept(r *http.Request, h Handler) bool
```
HandleAccept will process the Accept header and calls the given handler for each
mime type until the handler returns true.

This function returns true when the Handler returns true, false otherwise.

When no Accept header is given the mime string will be the empty string.

#### func  InvalidAccept

```go
func InvalidAccept(w http.ResponseWriter)
```
InvalidAccept writes the 406 header.

#### type Handler

```go
type Handler interface {
	Handle(mime Mime) bool
}
```

Handler provides an interface to handle a mime type.

The mime string (e.g. text/html, application/json, text/plain) is passed to the
handler, which is expected to return true if no more encodings are required and
false otherwise.

The empty string "" is used to signify when no preference is specified.

#### type HandlerFunc

```go
type HandlerFunc func(Mime) bool
```

HandlerFunc wraps a func to make it satisfy the Handler interface.

#### func (HandlerFunc) Handle

```go
func (h HandlerFunc) Handle(m Mime) bool
```
Handle calls the underlying func.

#### type Mime

```go
type Mime string
```

Mime represents a accepted Mime Type.

#### func (Mime) Match

```go
func (m Mime) Match(n Mime) bool
```
Match checks to see whether a given Mime Type matches the value.

The method allows for wildcards in the subtype sections.
