package jsonx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

// Marshal calls json.Marshal
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal calls json.Unmarshal
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// NewEncoder calls json.NewEncoder
func NewEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}

// NewDecoder calls json.NewDecoder
func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

type RawMessage []byte

func (d RawMessage) Len() int {
	return len(d)
}

func (d RawMessage) MarshalJSON() ([]byte, error) {
	if d == nil {
		return Null, nil
	}

	return d, nil
}

func (d *RawMessage) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("jay: UnmarshalJSON on nil pointer")
	}
	*d = append((*d)[0:0], data...)
	return nil
}

func (d RawMessage) IsObject() bool {
	return IsObject(d)
}

func (d RawMessage) IsEmptyObject() bool {
	return IsEmptyObject(d) && isEmpty(d)
}

func (d RawMessage) IsEmptyArray() bool {
	return IsEmptyArray(d)
}

// IsArray reports whether the data is a json array. It does not check whether
// the json is malformed.
func (d RawMessage) IsArray() bool {
	return IsArray(d)
}

func (d RawMessage) IsNull() bool {
	return IsNull(d)
}

// IsBool reports true if data appears to be a json boolean value. It is
// possible that it will report false positives of malformed json.
//
// IsBool does not parse strings
func (d RawMessage) IsBool() bool {
	return d.IsTrue() || d.IsFalse()
}

// IsTrue reports true if data appears to be a json boolean value of true. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsTrue does not parse strings
func (d RawMessage) IsTrue() bool {
	return IsTrue(d)
}

// IsFalse reports true if data appears to be a json boolean value of false. It is
// possible that it will report false positives of malformed json as it only
// checks the first character and length.
//
// IsFalse does not parse strings
func (d RawMessage) IsFalse() bool {
	return IsFalse(d)
}

func (d RawMessage) Equal(data []byte) bool {
	return bytes.Equal(d, data)
}

// ContainsEscapeRune reports whether the string value of d contains "\"
// It returns false if d is not a quoted string.
func (d RawMessage) ContainsEscapeRune() bool {
	for i := 0; i < len(d); i++ {
		if d[i] == '\\' {
			return true
		}
	}
	return false
}

func (d RawMessage) IsNumber() bool {
	return IsNumber(d)
}

func (d RawMessage) IsString() bool {
	return IsString(d)
}

type Object map[string]RawMessage

func (obj Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]RawMessage(obj))
}

func (obj *Object) UnmarshalJSON(data []byte) error {
	var m map[string]RawMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	*obj = m
	return nil
}

func IsEmptyObject(d []byte) bool {
	return IsObject(bytes.TrimSpace(d)) && len(d) == 2
}
