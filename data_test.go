package jsonx_test

// TODO: rewrite these tests

// type Obj struct {
// 	V string `json:"v"`
// }

// type Nested struct {
// 	Received []byte
// }

// func (n *Nested) UnmarshalJSON(data []byte) error {
// 	n.Received = data
// 	return nil
// }

// type Parent struct {
// 	unmarshalFn func([]byte) error
// 	Nested      Nested `json:"nested"`
// }

// func TestIsObject(t *testing.T) {
// 	m := map[string]map[string]string{
// 		"nested": {
// 			"val": "value",
// 		},
// 		"somethingElse": {},
// 		"other":         {},
// 	}
// 	b, _ := json.Marshal(m)
// 	assert := require.New(t)

// 	assert.True(jay.IsObject(b))
// 	var p Parent
// 	json.Unmarshal(b, &p)

// 	assert.Equal(`{"val":"value"}`, string(p.Nested.Received))
// }

// func TestJSON(t *testing.T) {
// 	data, err := json.Marshal(nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if
// 	assert.Equal("null", string(data))
// 	jd := jay.JSON(data)

// 	assert.True(jd.IsNull())
// 	assert.False(jd.IsEmptyObject())
// 	assert.False(jd.IsObject())
// 	assert.False(jd.IsEmptyArray())
// 	assert.False(jd.IsArray())

// 	data, err = json.Marshal([]string{"1,2,3"})
// 	assert.NoError(err)
// 	jd = jay.JSON(data)
// 	assert.True(jd.IsArray())
// 	assert.False(jd.IsNull())
// 	assert.False(jd.IsEmptyObject())
// 	assert.False(jd.IsObject())
// 	assert.False(jd.IsEmptyArray())

// 	data, err = json.Marshal([]string{})
// 	assert.NoError(err)
// 	jd = jay.JSON(data)
// 	assert.True(jd.IsArray())
// 	assert.False(jd.IsNull())
// 	assert.False(jd.IsEmptyObject())
// 	assert.False(jd.IsObject())
// 	assert.True(jd.IsEmptyArray())

// 	data, err = json.Marshal(map[string]string{})
// 	assert.NoError(err)
// 	jd = jay.JSON(data)
// 	assert.False(jd.IsArray())
// 	assert.False(jd.IsNull())
// 	assert.True(jd.IsEmptyObject())
// 	assert.True(jd.IsObject())
// 	assert.False(jd.IsEmptyArray())

// 	data, err = json.Marshal(map[string]string{"key": "val"})
// 	assert.NoError(err)
// 	jd = jay.JSON(data)
// 	assert.False(jd.IsArray())
// 	assert.False(jd.IsNull())
// 	assert.False(jd.IsEmptyObject())
// 	assert.True(jd.IsObject())
// 	assert.False(jd.IsEmptyArray())
// }

// func TestJSONOBject(t *testing.T) {
// 	assert := require.New(t)
// 	o := Obj{V: "value"}

// 	od, err := json.Marshal(o)
// 	assert.NoError(err)
// 	obj := jay.Object{
// 		"key": od,
// 	}
// 	objData, err := json.Marshal(obj)
// 	assert.NoError(err)
// 	assert.Equal(`{"key":{"v":"value"}}`, string(objData))
// }
