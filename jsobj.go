package jsobj

import (
	"encoding/json"
	"errors"
	"strings"
)

var tokens = map[rune]bool{}
var blanks = map[rune]bool{}

//Parser defines a method of reading JavaScript text based on Rune-based on-demand.
type Parser interface {
	ReadObjects() ([]interface{}, error)
	ReadObject() (interface{}, error)
	ReadWord() (Word, error)
	Read(...rune) (string, error)
	Location() int
	EOF() bool
}

func init() {
	for _, char := range "[{}],.:" {
		tokens[char] = true
	}
	for _, char := range " \t\r\n\b\f" {
		blanks[char] = true
	}
}

//Parse Create a new parser to start reading the JavaScript string.
func Parse(str string) Parser {
	var parser = &parser{Reader: strings.NewReader(str)}
	return parser
}

//Unmarshal Use the system JSON serialization tool to complete the reverse sequence.
func Unmarshal(src []byte, dst interface{}) error {
	parser := Parse(string(src))
	obj, err := parser.ReadObject()
	if err != nil {
		return err
	}
	if !parser.EOF() {
		return errors.New("js: object has not ended")
	}
	bin, _ := json.Marshal(obj)
	return json.Unmarshal(bin, dst)
}
