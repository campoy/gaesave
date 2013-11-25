// This package provides the Save and SaveStruct methods that provide a
// simplified interface to save object into the Google App Engine datastore.
//
// Save is a more general method, while SaveStruct offers a more constrained
// use case.
package gaesave

import (
	"fmt"
	"reflect"
	"strings"

	"appengine"
	"appengine/datastore"
)

// Savables can be saved by this package.
type Savable interface {
	ID() int64
	SetID(int64)
	Kind() string
}

// BeforeSave is called before putting the Savable in the datastore.
type BeforeSaver interface {
	BeforeSave() error
}

// BeforeSave is called after putting the Savable in the datastore and
// setting its ID field.
type AfterSaver interface {
	AfterSave() error
}

// If a type implements Elem the interface returned by the method will be
// put into the datastore instead of the instance itself.
type Elem interface {
	Elem() interface{}
}

// Save saves a Savable to the datastore taking into account before and
// after save methods.
func Save(c appengine.Context, obj Savable) (key *datastore.Key, err error) {
	if id := obj.ID(); id == 0 {
		key = datastore.NewKey(c, obj.Kind(), "", id, nil)
	} else {
		key = datastore.NewIncompleteKey(c, obj.Kind(), nil)
	}

	// BeforeSave hook.
	if bs, ok := obj.(BeforeSaver); ok {
		if err := bs.BeforeSave(); err != nil {
			return nil, err
		}
	}

	// Actual datastore Put.
	if e, ok := obj.(Elem); ok {
		fmt.Printf("put %T\n", e.Elem())
		key, err = datastore.Put(c, key, e.Elem())
	} else {
		key, err = datastore.Put(c, key, obj)
	}
	if err != nil {
		return nil, err
	}
	obj.SetID(key.IntID())

	// AfterSave hook.
	if as, ok := obj.(AfterSaver); ok {
		if err := as.AfterSave(); err != nil {
			return nil, err
		}
	}

	return key, nil
}

// structSaver implements Savable given a reflect.Value of Kind struct.
type structSaver struct {
	e interface{}
	v reflect.Value
}

// ID returns the value of the ID field of the struct if it exists or zero.
func (s structSaver) ID() int64 {
	f := s.v.FieldByName("ID")
	if f.IsValid() {
		return f.Int()
	}
	return 0
}

// SetID sets the ID field of the struct to id if it exists.
func (s structSaver) SetID(id int64) {
	f := s.v.FieldByName("ID")
	if f.IsValid() {
		f.SetInt(id)
	}
}

// Kind returns the type name without package name.
func (s structSaver) Kind() string {
	typ := reflect.TypeOf(s.v).String()
	ps := strings.Split(typ, ".")
	return ps[len(ps)-1]
}

func (s structSaver) Elem() interface{} { return s.e }

func NewSavableFromStruct(obj interface{}) (Savable, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can't save %T", obj)
	}
	return structSaver{obj, val}, nil
}

// SaveStruct saves the given struct to the datastore.
func SaveStruct(c appengine.Context, obj interface{}) (*datastore.Key, error) {
	s, err := NewSavableFromStruct(obj)
	if err != nil {
		return nil, err
	}
	return Save(c, s)
}
