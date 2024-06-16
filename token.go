package xmltokenizer

import "sync"

var pool = sync.Pool{New: func() any { return new(Token) }}

// GetToken gets token from the pool, don't forget to put it back.
func GetToken() *Token { return pool.Get().(*Token) }

// PutToken puts token back to the pool.
func PutToken(t *Token) { pool.Put(t) }

// Token represent a single token, one of these following:
//   - <?xml version="1.0" encoding="UTF-8"?>
//   - <name attr="value" attr="value">
//   - <name attr="value" attr="value">CharData
//   - <name attr="value" attr="value"><![CDATA[ CharData ]]>
//   - <name attr="value" attr="value"/>
//   - </name>
//   - <!-- a comment -->
//   - <!DOCTYPE library [
//     <!ELEMENT library (book+)>
//     <!ELEMENT book (title, author, year)>
//     ]>
type Token struct {
	Name        Name   // Name can be a StartElement: "name", a EndElement: "/name" or empty when a tag starts with "<?" or "<!" (except "<![CDATA").
	Attrs       []Attr // Attrs exist when len(Attrs) > 0.
	Data        []byte // Data could be a CharData or a CDATA, or maybe a RawToken if a tag starts with "<?" or "<!" (except "<![CDATA").
	SelfClosing bool   // True when a tag ends with "/>" e.g. <c r="E3" s="1" />. Also true when a tag starts with "<?" or "<!" (except "<![CDATA").
}

// IsEndElement checks whether the given token represent an end element (closing tag),
// name start with '/'. e.g. </gpx>
func (t *Token) IsEndElement() bool {
	if len(t.Name.Full) > 0 && t.Name.Full[0] == '/' {
		return true
	}
	return false
}

// IsEndElementOf checks whether the given token represent a
// n end element (closing tag) of given startElement.
func (t *Token) IsEndElementOf(t2 *Token) bool {
	if !t.IsEndElement() {
		return false
	}
	if string(t.Name.Full[1:]) == string(t2.Name.Full) {
		return true
	}
	return false
}

// Copy copies src Token into t, returning t. Attrs should be
// consumed immediately since it's only being shallow copied.
func (t *Token) Copy(src Token) *Token {
	t.Name.Space = append(t.Name.Space[:0], src.Name.Space...)
	t.Name.Local = append(t.Name.Local[:0], src.Name.Local...)
	t.Name.Full = append(t.Name.Full[:0], src.Name.Full...)
	t.Attrs = append(t.Attrs[:0], src.Attrs...) // shallow copy
	t.Data = append(t.Data[:0], src.Data...)
	t.SelfClosing = src.SelfClosing
	return t
}

// Attr represents an XML attribute.
type Attr struct {
	Name  Name
	Value []byte
}

// Name represents an XML name (Local) annotated
// with a name space identifier (Space).
type Name struct {
	Space []byte
	Local []byte
	Full  []byte // Full is combination of "space:local"
}
