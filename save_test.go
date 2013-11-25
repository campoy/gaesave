package gaesave

import (
	"testing"

	"appengine/aetest"
)

type Person struct {
	Id   int64
	Name string
	Age  int64
}

func (p Person) ID() int64 { return p.ID() }

func TestSimpleSave(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	_, err = SaveStruct(c, &Person{Name: "John", Age: 42})
	if err != nil {
		t.Error(err)
	}
}
