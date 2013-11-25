package gaesave

import (
	"reflect"
	"testing"

	"appengine/aetest"
	"appengine/datastore"
)

type Person struct {
	ID   int64
	Name string
	Age  int64
}

func TestSimpleSave(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	p := &Person{ID: 1, Name: "John", Age: 42}
	_, err = SaveStruct(c, p)
	if err != nil {
		t.Error(err)
	}

	ek := datastore.NewKey(c, "Person", "", 1, nil)
	var q Person
	if err := datastore.Get(c, ek, &q); err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(p, q) {
		t.Error("saved %v, got %v", p, q)
	}
}

type Hooks struct {
	Person
	bs, as bool
}

func (p *Hooks) BeforeSave() error {
	p.bs = true
	return nil
}

func (p *Hooks) AfterSave() error {
	p.as = true
	return nil
}

func TestBeforeSave(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	p := &Hooks{Person: Person{Name: "John", Age: 42}}
	_, err = SaveStruct(c, p)
	if err != nil {
		t.Error(err)
	}
	if !p.bs {
		t.Error("BeforeSave wasn't called")
	}
	if !p.as {
		t.Error("AfterSave wasn't called")
	}
}
