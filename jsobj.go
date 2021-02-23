package jsobj

import (
	"encoding/json"
	"strings"
)

var tokens = map[rune]bool{}
var blanks = map[rune]bool{}

func init() {
	for _, char := range "[{}],.:" {
		tokens[char] = true
	}
	for _, char := range " \t\r\n\b\f" {
		blanks[char] = true
	}
}

//Parse Parse JavaScript object string to an object (map/array)
func Parse(str string) (interface{}, error) {
	var obj = (*parser)(strings.NewReader(str))
	return obj.Parse()
}

//Unmarshal Use the system JSON serialization tool to complete the reverse sequence.
func Unmarshal(src []byte, dst interface{}) error {
	obj, err := Parse(string(src))
	if err != nil {
		return err
	}
	bin, _ := json.Marshal(obj)
	return json.Unmarshal(bin, dst)
}
