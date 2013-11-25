package gaesave

import (
	"testing"

	"appengine/aetest"
)

type Person struct {
	Name string
	Age  int64
}

func TestSimpleSave(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	p := &Person{Name: "John", Age: 42}
	_, err = SaveStruct(c, p)
	if err != nil {
		t.Error(err)
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
