package jsobj

import "strconv"

type word struct {
	token bool
	must  bool
	text  string
}

func (w *word) getValue() interface{} {
	if w.must {
		return w.text
	}
	if w.text == "null" {
		return nil
	}
	if obj, err := strconv.ParseInt(w.text, 10, 64); err == nil {
		return obj
	}
	if obj, err := strconv.ParseUint(w.text, 10, 64); err == nil {
		return obj
	}
	if obj, err := strconv.ParseFloat(w.text, 64); err == nil {
		return obj
	}
	if obj, err := strconv.ParseBool(w.text); err == nil {
		return obj
	}
	return w.text
}
