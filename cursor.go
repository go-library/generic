package generic

import (
	"fmt"
	"reflect"
)

type Cursor struct {
	// interface
	parent reflect.Value
	// for map or slice
	myKey reflect.Value
}

/* NEW FUNCTIONS */

func NewCursor(root *interface{}) *Cursor {
	c := new(Cursor)
	c.parent = reflect.ValueOf(root).Elem()
	return c
}

/* ACCESS INDEX ELEMENT, IF NOT EXISTED CREATE OR PANIC */

func (c *Cursor) Index(keys ...interface{}) (nextCursor *Cursor) {
	nextCursor = new(Cursor)
	nextCursor.parent = c.parent
	nextCursor.myKey = c.myKey

	for i := range keys {
		k := reflect.ValueOf(keys[i])

		nextCursor.prepareToNext(nextCursor.Value(), k, false)
		nextCursor.parent = nextCursor.Value()
		nextCursor.myKey = k
		if !nextCursor.Value().IsValid() {
			nextCursor.SetEmpty()
		}
	}

	return nextCursor
}

func (c *Cursor) SetIndex(keys ...interface{}) (nextCursor *Cursor) {
	nextCursor = new(Cursor)
	nextCursor.parent = c.parent
	nextCursor.myKey = c.myKey

	for i := range keys {
		k := reflect.ValueOf(keys[i])

		nextCursor.prepareToNext(nextCursor.Value(), k, true)
		nextCursor.parent = nextCursor.Value()
		nextCursor.myKey = k
		if !nextCursor.Value().IsValid() {
			nextCursor.SetEmpty()
		}
	}

	return nextCursor
}

/* GETTERS */

func (c *Cursor) Interface() interface{} {
	return c.Value().Interface()
}

func (c *Cursor) Value() reflect.Value {
	var value reflect.Value

	switch c.parent.Kind() {
	case reflect.Interface:
		value = c.parent
	case reflect.Map:
		value = c.parent.MapIndex(c.myKey)
	case reflect.Slice:
		value = c.parent.Index(int(c.myKey.Int()))
	default:
		value = c.parent
	}

	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	return value
}

/* SLICE FUNCTIONS */

func (c *Cursor) Len() int {
	if c.Value().Kind() != reflect.Slice {
		panic(fmt.Errorf("this value is not slice"))
	}

	return c.Value().Len()
}

func (c *Cursor) Push(values ...interface{}) {
	var vars []reflect.Value
	for i := range values {
		if values[i] == nil {
			vars = append(vars, reflect.Zero(c.Value().Type().Elem()))
		} else {
			vars = append(vars, reflect.ValueOf(values[i]))
		}
	}

	c.PushValues(vars...)
}

func (c *Cursor) PushValues(values ...reflect.Value) {
	if c.Value().Kind() != reflect.Slice {
		panic(fmt.Errorf("this value is not slice"))
	}

	c.SetValue(reflect.Append(c.Value(), values...))
}

func (c *Cursor) Slice(i, j int) (nextCursor *Cursor) {
	v := c.Value().Slice(i, j)
	nc := NewCursor(new(interface{}))
	nc.SetValue(v)
	return nc
}

func (c *Cursor) Slice3(i, j, k int) (nextCursor *Cursor) {
	v := c.Value().Slice3(i, j, k)
	nc := NewCursor(new(interface{}))
	nc.SetValue(v)
	return nc
}

/* MAP FUNCTIONS */

func (c *Cursor) Keys() (keys []string) {
	if c.Value().Kind() != reflect.Map {
		panic(fmt.Errorf("this value is not map"))
	}

	keyValues := c.Value().MapKeys()
	for i := range keyValues {
		keys = append(keys, keyValues[i].String())
	}

	return keys
}

/* SETTERS */

func (c *Cursor) Set(value interface{}) {
	c.SetValue(reflect.ValueOf(value))
}

func (c *Cursor) SetMap() {
	c.SetValue(makeMap())
}

func (c *Cursor) SetSlice() {
	c.SetValue(makeSlice(0, 0))
}

func (c *Cursor) SetEmpty() {
	if c.parent.IsValid() {
		if c.parent.Kind() == reflect.Interface {
			c.SetValue(reflect.Zero(c.parent.Type()))
		} else {
			c.SetValue(reflect.Zero(c.parent.Type().Elem()))
		}
	} else {
		panic(fmt.Errorf("parent value is invalid"))
	}
}

func (c *Cursor) SetValue(value reflect.Value) {
	switch c.parent.Kind() {
	case reflect.Interface:
		c.parent.Set(value)
	case reflect.Map:
		c.parent.SetMapIndex(c.myKey, value)
	case reflect.Slice:
		c.parent.Index(int(c.myKey.Int())).Set(value)
	default:
		panic(fmt.Errorf("unsupported value kind: %s", c.parent.Kind()))
	}
}

/* HIDDEN FUNCTIONS */

func (c *Cursor) prepareToNext(value, key reflect.Value, permitCreate bool) {
	var (
		nextValue reflect.Value
		isCreated = false
		vk        = value.Kind()
		kk        = key.Kind()
	)

	//check permission or type

	switch kk {
	case reflect.String, reflect.Int:
	default:
		panic(fmt.Errorf("key should be string or integer"))
	}

	switch vk {
	case reflect.Map, reflect.Slice, reflect.Invalid:
	default:
		panic(fmt.Errorf("value should be map, slice or invalid"))
	}

	switch {
	case vk == reflect.Map && kk == reflect.String:
	case vk == reflect.Slice && kk == reflect.Int:
	case vk == reflect.Invalid:
		// check create permission
		if !permitCreate {
			panic(fmt.Errorf("implicated creation failure"))
		}
	default:
		// check override permission
		panic(fmt.Errorf("key miss matched in value kind"))
	}

	switch {
	case kk == reflect.String && vk == reflect.Map:
		// nothing to do
		nextValue = value
	case kk == reflect.String && vk == reflect.Invalid:
		// make new map
		nextValue = makeMap()
		isCreated = true
	case kk == reflect.Int && vk == reflect.Slice:
		idx := int(key.Int())
		if idx >= value.Cap() {
			// check increase permission
			nextValue = makeSlice(idx+1, idx+1)
			reflect.Copy(nextValue, value)
			isCreated = true
		} else if idx < value.Cap() && idx >= value.Len() {
			// check increase permission
			nextValue = value.Slice(0, idx+1)
			isCreated = true
		}
	case kk == reflect.Int && vk == reflect.Invalid:
		// make new slice
		idx := int(key.Int())
		nextValue = makeSlice(idx+1, idx+1)
		isCreated = true
	}

	if isCreated {
		c.SetValue(nextValue)
	}
}

func makeMap() reflect.Value {
	var objectInterface map[string]interface{}
	return reflect.MakeMap(reflect.TypeOf(objectInterface))
}

func makeSlice(len int, cap int) reflect.Value {
	var (
		n uint
		i uint
	)

	for i = 1; ; i++ {
		n = 1 << i

		if int(n) >= cap {
			cap = int(n)
			break
		}
	}
	var arrayInterface []interface{}
	return reflect.MakeSlice(reflect.TypeOf(arrayInterface), len, cap)
}