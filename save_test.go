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

type CheckedPerson struct {
	Person
	saved bool
}

func (p *CheckedPerson) BeforeSave() error {
	p.saved = true
	return nil
}

func TestBeforeSave(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	p := &CheckedPerson{Person{Name: "John", Age: 42}, false}
	_, err = SaveStruct(c, p)
	if err != nil {
		t.Error(err)
	}
	if p.saved == false {
		t.Error("BeforeSave wasn't called")
	}
}
