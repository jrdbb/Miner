package crawler

import (
	"errors"
	"github.com/robertkrimen/otto"
)

func ParseArray(array *otto.Value) ([]*otto.Value, error) {
	if !array.IsObject() {
		return nil, errors.New("ParseArray: array is not an object")
	}
	keys := array.Object().Keys()
	ret := make([]*otto.Value, 0, len(keys))

	for _, k := range keys {
		v, err := array.Object().Get(k)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &v)
	}
	return ret, nil
}
