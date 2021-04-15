package jsobj

import "strings"

var functions = map[string]Function{}

type Function = func(args ...interface{}) (obj interface{}, err error)

func RegisterMethod(name string, function Function, caseNotSensitive ...bool) (err error) {
	if len(caseNotSensitive) > 0 && caseNotSensitive[0] {
		name = strings.ToLower(name)
	}
	functions[name] = function
	return
}
