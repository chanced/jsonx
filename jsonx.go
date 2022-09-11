package jsonx

import (
	"bytes"
	"unicode"
)

var (
	Null     = RawMessage("null")
	True     = RawMessage("true")
	False    = RawMessage("false")
	trueVal  = []byte("true")  // unnecessary but what if you accidentally override the above?
	falseVal = []byte("false") // ^^
	nullVal  = []byte("null")
)

// IsValid reports whether s is a valid JSON number literal.
//
// Taken from encoding/json/scanner.go
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
