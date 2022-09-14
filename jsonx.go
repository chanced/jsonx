package jsonx

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

var (
	Null     = RawMessage("null")
	True     = RawMessage("true")
	False    = RawMessage("false")
	trueVal  = []byte("true")  // unnecessary but what if you accidentally override the above?
	falseVal = []byte("false") // ^^
	nullVal  = []byte("null")
)

type Writer interface {
	Write(p []byte) (n int, err error)
	WriteByte(c byte) error
	WriteString(s string) (n int, err error)
}

// IsValid reports whether s is a valid JSON number literal.
//
// Taken from encoding/json
func IsNumber(data []byte) bool {
	// This function implements the JSON numbers grammar.
	// See https://tools.ietf.org/html/rfc7159#section-6
	// and https://www.json.org/img/number.png

	if len(data) == 0 {
		return false
	}

	// Optional -
	if data[0] == '-' {
		data = data[1:]
		if len(data) == 0 {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case data[0] == '0':
		data = data[1:]

	case '1' <= data[0] && data[0] <= '9':
		data = data[1:]
		for len(data) > 0 && '0' <= data[0] && data[0] <= '9' {
			data = data[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(data) >= 2 && data[0] == '.' && '0' <= data[1] && data[1] <= '9' {
		data = data[2:]
		for len(data) > 0 && '0' <= data[0] && data[0] <= '9' {
			data = data[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(data) >= 2 && (data[0] == 'e' || data[0] == 'E') {
		data = data[1:]
		if data[0] == '+' || data[0] == '-' {
			data = data[1:]
			if len(data) == 0 {
				return false
			}
		}
		for len(data) > 0 && '0' <= data[0] && data[0] <= '9' {
			data = data[1:]
		}
	}

	// Make sure we are at the end.
	return len(data) == 0
}

var hex = "0123456789abcdef"

// EncodeString encodes a string to a json string literal, escaping any
// characters that are not permitted.
//
// EncodeString calls EncodeAndWriteString with a bytes.Buffer and returns the
// result.
//
// Logic taken from encoding/json
func EncodeString[S ~string](s S) []byte {
	b := bytes.Buffer{}
	encodeAndWriteString(&b, []byte(s), true)
	return b.Bytes()
}

// EncodeAndWriteString encodes a string to a json string literal, escaping any
// characters that are not permitted. It then writes the result to w.
//
// Logic taken from encoding/json
func EncodeAndWriteString[S ~string](w Writer, s S) {
	encodeAndWriteString(w, []byte(s), true)
}

func EncodeAndWriteStringWithoutHTMLEscape(w Writer, s []byte) {
	encodeAndWriteString(w, s, false)
}

func encodeAndWriteString(w Writer, s []byte, escapeHTML bool) {
	w.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				w.Write(s[start:i])
			}
			w.WriteByte('\\')
			switch b {
			case '\\', '"':
				w.WriteByte(b)
			case '\n':
				w.WriteByte('n')
			case '\r':
				w.WriteByte('r')
			case '\t':
				w.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w.WriteString(`u00`)
				w.WriteByte(hex[b>>4])
				w.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w.Write(s[start:i])
			}
			w.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w.Write(s[start:i])
			}
			w.WriteString(`\u202`)
			w.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		w.Write(s[start:])
	}
	w.WriteByte('"')
}

func IsNull(d []byte) bool {
	return bytes.Equal(d, nullVal)
}

func IsBool(d []byte) bool {
	return IsTrue(d) || IsFalse(d)
}

func IsString(d []byte) bool {
	return startsAndEndsWith(bytes.TrimSpace(d), '"', '"')
}

func IsObject(d []byte) bool {
	return startsAndEndsWith(bytes.TrimSpace(d), '{', '}')
}

func IsArray(d []byte) bool {
	return startsAndEndsWith(bytes.TrimSpace(d), '[', ']')
}

func IsEmptyArray(d []byte) bool {
	return IsArray(d) && isEmpty(d)
}

// IsTrue reports true if data appears to be a json boolean value of true. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsTrue does not parse strings
func IsTrue(d []byte) bool {
	return bytes.Equal(d, True)
}

// IsFalse reports true if data appears to be a json boolean value of false. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsFalse does not parse strings
func IsFalse(d []byte) bool {
	return bytes.Equal(d, falseVal)
}

func isEmpty(d []byte) bool {
	count := 0
	for _, v := range d {
		if !unicode.IsSpace(rune(v)) {
			count += 1
			if count > 2 {
				return false
			}
		}
	}
	return count == 2
}

func startsAndEndsWith(d []byte, start, end byte) bool {
	if len(d) < 2 {
		return false
	}

	return d[0] == start && d[len(d)-1] == end
}
