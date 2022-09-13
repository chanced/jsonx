package jsonx_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/chanced/jsonx"
)

func TestEncodeString(t *testing.T) {
	str := `
<html>
	<body>
  		"hello world!"
	</body>
</html>`

	expected, err := json.Marshal(str)
	if err != nil {
		t.Fatal(err)
	}
	b := bytes.Buffer{}

	jsonx.EncodeAndWriteString(&b, []byte(str))
	if b.String() != string(expected) {
		t.Errorf("expected %s, got %s", expected, b.String())
	}
}
