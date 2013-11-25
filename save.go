package gaesave

import (
	"errors"
	"reflect"
	"strings"

	"appengine"
	"appengine/datastore"
)

type IDer interface {
	ID() int64
}

type IDSetter interface {
	SetID(int64)
}

type BeforeSaver interface {
	BeforeSave() error
}

type AfterSaver interface {
	AfterSave() error
}

func Save(c appengine.Context, obj IDer) (*datastore.Key, error) {
	kind, val := reflect.TypeOf(obj), reflect.ValueOf(obj)
	str := val
	if val.Kind() == reflect.Ptr {
		kind, str = kind.Elem(), val.Elem()
	}
	if str.Kind() != reflect.Struct {
		return nil, errors.New("Must pass a struct or pointer to struct")
	}
	dsKind := kind.String()
	if li := strings.LastIndex(dsKind, "."); li >= 0 {
		//Format kind to be in a standard format used for datastore
		dsKind = dsKind[li+1:]
	}

	var key *datastore.Key
	if id := obj.ID(); id == 0 {
		key = datastore.NewKey(c, dsKind, "", id, nil)
	} else {
		key = datastore.NewIncompleteKey(c, dsKind, nil)
	}

	if bs, ok := obj.(BeforeSaver); ok {
		if err := bs.BeforeSave(); err != nil {
			return nil, err
		}
	}

	key, err := datastore.Put(c, key, obj)
	if err != nil {
		return nil, err
	}

	if is, ok := obj.(IDSetter); ok {
		is.SetID(key.IntID())
	}

	if as, ok := obj.(AfterSaver); ok {
		if err := as.AfterSave(); err != nil {
			return nil, err
		}
	}

	return key, nil
}
