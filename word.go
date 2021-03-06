package jsobj

import (
	"reflect"
	"strconv"
)

//Word represents a JavaScript token or a JavaScript value.
type Word interface {
	Token() bool
	String() string
	Type() reflect.Kind
	Typed() interface{}
}

type word struct {
	token bool
	must  bool
	text  string
	value interface{}
	t     reflect.Kind
}

func (w *word) Typed() interface{} {
	if w.Type(); w.t != reflect.Invalid {
		return w.value
	}
	return w.Typed()
}

func (w *word) Token() bool {
	return w.token
}

func (w *word) String() string {
	return w.text
}

func (w *word) Type() (kind reflect.Kind) {
	if w.t != reflect.Invalid {
		return w.t
	}
	var obj interface{}
	var err error
	defer func() {
		w.t = kind
		if kind == reflect.String {
			w.value = w.text
		} else {
			w.value = obj
		}
	}()
	if w.must {
		return reflect.String
	}
	if w.text == "null" {
		obj = nil
		return reflect.Ptr
	}
	if obj, err = strconv.ParseInt(w.text, 10, 64); err == nil {
		return reflect.Int
	}
	if obj, err = strconv.ParseUint(w.text, 10, 64); err == nil {
		return reflect.Uint
	}
	if obj, err = strconv.ParseFloat(w.text, 64); err == nil {
		return reflect.Float64
	}
	if obj, err = strconv.ParseBool(w.text); err == nil {
		return reflect.Bool
	}
	return reflect.String
}
