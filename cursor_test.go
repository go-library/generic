package generic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func CheckWithJson(t *testing.T, stream string, value interface{}) error {
	var (
		err error
		j   interface{}
	)

	err = json.Unmarshal([]byte(stream), &j)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(value, j) {
		t.Logf("JSON:  %+v\n", j)
		t.Logf("VALUE: %+v\n", value)
		return fmt.Errorf("not equal")
	}

	return nil
}

func TestCursor_Set_Interface(t *testing.T) {
	var (
		err error
		v   interface{}
		c   *Cursor
	)

	// init
	c = NewCursor(&v)
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// empty
	c.SetEmpty()
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// string
	c.Set("value")
	if err = CheckWithJson(t, `"value"`, v); err != nil {
		t.Error(err)
	}

	// float
	c.Set(float64(100))
	if err = CheckWithJson(t, `100`, v); err != nil {
		t.Error(err)
	}

	// boolean
	c.Set(true)
	if err = CheckWithJson(t, `true`, v); err != nil {
		t.Error(err)
	}

	// empty
	c.SetEmpty()
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// map
	c.SetMap()
	if err = CheckWithJson(t, `{}`, v); err != nil {
		t.Error(err)
	}

	// slice
	c.SetSlice()
	if err = CheckWithJson(t, `[]`, v); err != nil {
		t.Error(err)
	}
}

func TestCursor_Set_Map(t *testing.T) {
	var (
		err error
		v   interface{}
		c   *Cursor
	)

	// init
	c = NewCursor(&v)
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// map
	c.SetMap()
	if err = CheckWithJson(t, `{}`, v); err != nil {
		t.Error(err)
	}

	// map - empty
	c.SetIndex("value").SetEmpty()
	if err = CheckWithJson(t, `{"value":null}`, v); err != nil {
		t.Error(err)
	}

	// map - string
	c.SetIndex("value").Set("string")
	if err = CheckWithJson(t, `{"value":"string"}`, v); err != nil {
		t.Error(err)
	}

	// map - float
	c.SetIndex("value").Set(float64(100))
	if err = CheckWithJson(t, `{"value":100}`, v); err != nil {
		t.Error(err)
	}

	// map - boolean
	c.SetIndex("value").Set(true)
	if err = CheckWithJson(t, `{"value":true}`, v); err != nil {
		t.Error(err)
	}

	// map - empty
	c.SetIndex("value").SetEmpty()
	if err = CheckWithJson(t, `{"value":null}`, v); err != nil {
		t.Error(err)
	}

	// map - map
	c.SetIndex("value").SetMap()
	if err = CheckWithJson(t, `{"value":{}}`, v); err != nil {
		t.Error(err)
	}

	// map - slice
	c.SetIndex("value").SetSlice()
	if err = CheckWithJson(t, `{"value":[]}`, v); err != nil {
		t.Error(err)
	}

	// map - delete
	c.SetIndex("value").Set(nil)
	if err = CheckWithJson(t, `{}`, v); err != nil {
		t.Error(err)
	}
}

func TestCursor_Set_Slice(t *testing.T) {
	var (
		err error
		v   interface{}
		c   *Cursor
	)

	// init
	c = NewCursor(&v)
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// slice
	c.SetSlice()
	if err = CheckWithJson(t, `[]`, v); err != nil {
		t.Error(err)
	}

	// slice - empty
	c.SetIndex(0).SetEmpty()
	if err = CheckWithJson(t, `[null]`, v); err != nil {
		t.Error(err)
	}

	// slice - string
	c.SetIndex(0).Set("string")
	if err = CheckWithJson(t, `["string"]`, v); err != nil {
		t.Error(err)
	}

	// slice - float
	c.SetIndex(0).Set(float64(100))
	if err = CheckWithJson(t, `[100]`, v); err != nil {
		t.Error(err)
	}

	// slice - boolean
	c.SetIndex(0).Set(true)
	if err = CheckWithJson(t, `[true]`, v); err != nil {
		t.Error(err)
	}

	// slice - empty
	c.SetIndex(0).SetEmpty()
	if err = CheckWithJson(t, `[null]`, v); err != nil {
		t.Error(err)
	}

	// slice -map
	c.SetIndex(0).SetMap()
	if err = CheckWithJson(t, `[{}]`, v); err != nil {
		t.Error(err)
	}

	// slice - slice
	c.SetIndex(0).SetSlice()
	if err = CheckWithJson(t, `[[]]`, v); err != nil {
		t.Error(err)
	}
}

func TestCursor_Slice_Increase(t *testing.T) {
	var (
		err error
		v   interface{}
		c   *Cursor
	)
	// init
	c = NewCursor(&v)
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// slice
	c.SetSlice()
	if err = CheckWithJson(t, `[]`, v); err != nil {
		t.Error(err)
	}

	var items []string
	for i := 0; i < 256; i++ {
		c.SetIndex(i).Set(float64(i))

		items = append(items, fmt.Sprintf("%d", i))

		if err = CheckWithJson(t, fmt.Sprintf(`[%s]`, strings.Join(items, ",")), v); err != nil {
			t.Error(err)
		}
	}

	if 256 != c.Len() {
		t.Log("Len() is incorrect")
	}
}

func TestCursor_Map_Keys(t *testing.T) {
	var (
		err error
		v   interface{}
		c   *Cursor
	)
	// init
	c = NewCursor(&v)
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// map
	c.SetMap()
	if err = CheckWithJson(t, `{}`, v); err != nil {
		t.Error(err)
	}

	c.SetIndex("a")
	c.SetIndex("b")
	c.SetIndex("c")
	keys := c.Keys()
	match := 0
	for i := range keys {
		key := keys[i]
		if key == "a" {
			match++
		} else if key == "b" {
			match++
		} else if key == "c" {
			match++
		} else {
			match = 10
		}
	}

	if match != 3 {
		t.Error("Keys() is incorrect:", c.Keys())
	}

}

func TestCursor_Push(t *testing.T) {
	var (
		err error
		v   interface{}
		c   *Cursor
	)
	// init
	c = NewCursor(&v)
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// slice
	c.SetSlice()
	if err = CheckWithJson(t, `[]`, v); err != nil {
		t.Error(err)
	}

	c.Push(true, "string", float64(100), nil)
	if err = CheckWithJson(t, `[true,"string",100,null]`, v); err != nil {
		t.Error(err)
	}

}
func TestCursor_Index(t *testing.T) {
	var (
		err error
		v   interface{}
		c   *Cursor
	)
	// init
	c = NewCursor(&v)
	if err = CheckWithJson(t, `null`, v); err != nil {
		t.Error(err)
	}

	// index function test
	c.SetIndex("a", "b", "c", "d", "e").Set("var")
	if err = CheckWithJson(t, `{"a":{"b":{"c":{"d":{"e":"var"}}}}}`, v); err != nil {
		t.Error(err)
	}
}
