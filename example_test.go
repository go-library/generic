package generic

import (
	"encoding/json"
	"fmt"
	"log"
)

func ExampleCursor_Index() {
	var (
		err error
		v   interface{}
		c   *Cursor
	)

	stream := `{"results": [{"name": "foo" }]}`

	err = json.Unmarshal([]byte(stream), &v)
	if err != nil {
		log.Fatal(err)
	}

	c = NewCursor(&v)
	c.Index("results", 0, "name").Set("bar")
	fmt.Printf("%+v\n", v)

	// Output: map[results:[map[name:bar]]]
}
