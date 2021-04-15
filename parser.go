package jsobj

import (
	"errors"
	"strings"
)

type parser struct {
	*strings.Reader
	location int
	tokenExt map[rune]bool
}

func (parser *parser) EOF() bool {
	return parser.Len() == 0
}

func (parser *parser) Location() int {
	return parser.location
}

func (parser *parser) ReadObject() (obj Object, err error) {
	var w Word
	w, err = parser.ReadWord()
	if err == nil {
		//expected { / [ or a pure text
		if w.Token() {
			var result interface{}
			switch w.String() {
			case "{":
				result, err = parser.readObj()
				obj = &object{result}
				break
			case "[":
				result, err = parser.readArray()
				obj = &object{result}
				break
			default:
				//other characters
				err = unexpectedWordError(w.String(), parser.location)
			}
		} else {
			//text word can't have a follow-up word.
			obj = &object{w.Typed()}
		}
	}
	return
}

func (parser *parser) ReadObjects() ([]interface{}, error) {
	defer func() {
		parser.tokenExt = map[rune]bool{}
	}()
	w, err := parser.ReadWord()
	if err != nil {
		return nil, err
	}
	if w.String() != "(" {
		return nil, unexpectedWordError(w.String(), parser.Location())
	}
	var results = []interface{}{}
	for {
		obj, err := parser.ReadObject()
		if err != nil {
			if isUnexpectedWordError(err, ",") {
				continue
			} else if isUnexpectedWordError(err, ")") {
				return results, nil
			} else {
				return nil, err
			}
		}
		if obj == nil {
			return nil, err
		}
		results = append(results, obj.Interface())
	}
}

func (parser *parser) readObj() (interface{}, error) {
	// readed first '{'
	//expect token '}', string
	var result, hasComma, key = map[string]interface{}{}, true, ""
	for {
		w, err := parser.ReadWord()
		if err != nil {
			return nil, err
		}
		if w.Token() || !hasComma {
			if w.String() == "}" {
				return result, nil
			}
			//if there is no comma, the object must be finished
			return nil, unexpectedWordError(w.String(), parser.location)
		}
		key = w.String()
		w, err = parser.ReadWord()
		if err != nil {
			return nil, err
		}
		//key : value
		if w.String() != ":" {
			return nil, unexpectedWordError(w.String(), parser.location)
		}
		w, err = parser.ReadWord()
		if err != nil {
			return nil, err
		}
		//an object, an array or just a text
		if w.Token() {
			var obj Object
			parser.UnreadRune()
			obj, err = parser.ReadObject()
			if err != nil {
				return nil, err
			}
			result[key] = obj.Interface()
		} else {
			//text
			result[key] = w.Typed()
		}
		//must be a token in comma or bracket
		w, err = parser.ReadWord()
		if err == nil && w.String() == "(" {
			if result[key], err = parser.callMethod(result[key].(string)); err == nil {
				w, err = parser.ReadWord()
			}
		}
		if err != nil {
			return nil, err
		}
		if !w.Token() {
			return nil, unexpectedWordError(w.String(), parser.location)
		}
		if hasComma = w.String() == ","; !hasComma {
			//meet a none-comma token, unread
			parser.UnreadRune()
		}
	}
}

func (parser *parser) callMethod(method string) (interface{}, error) {
	parser.UnreadRune()
	function, existed := functions[method]
	if !existed {
		function, existed = functions[strings.ToLower(method)]
	}
	if existed {
		if args, err := parser.ReadObjects(); err == nil {
			return function(args...)
		}
	}
	return nil, errors.New("no method named '" + method + "'")
}

func (parser *parser) readArray() (interface{}, error) {
	// readed first '['
	//expect token ']', string, '{', '['
	var result, hasComma = []interface{}{}, true
	for {
		w, err := parser.ReadWord()
		if err != nil {
			return nil, err
		}
		if w.Token() {
			if !hasComma && w.String() != "]" {
				return nil, unexpectedWordError(w.String(), parser.location)
			}
			if w.String() == "]" {
				return result, nil
			}
			parser.UnreadRune()
			var obj Object
			obj, err = parser.ReadObject()
			if err != nil {
				return nil, err
			}
			result = append(result, obj.Interface())
		} else {
			result = append(result, w.Typed())
		}
		//must be a token in comma or bracket
		w, err = parser.ReadWord()
		if err != nil {
			return nil, err
		}
		if !w.Token() {
			return nil, unexpectedWordError(w.String(), parser.location)
		}
		if hasComma = w.String() == ","; !hasComma {
			//meet a none-comma token, unread
			parser.UnreadRune()
		}
	}
}

func (parser *parser) ReadWord() (w Word, err error) {
	var char, quote, last rune
	char, err = parser.NextRune(true)
	if err != nil {
		return
	}
	if tokens[char] || parser.tokenExt[char] {
		w = &word{token: true, text: string(char)}
		return
	}
	var builder = strings.Builder{}
	if char == '\'' || char == '"' {
		quote = char
	} else {
		builder.WriteRune(char)
	}
	for char, err = parser.NextRune(); err == nil; char, err = parser.NextRune() {
		if last == '\\' && char != quote {
			builder.WriteRune(last)
		}
		if char == quote && last != '\\' {
			break
		}
		if quote == 0 {
			if blanks[char] {
				break
			}
			if tokens[char] {
				parser.UnreadRune()
				break
			}
		}
		if last = char; last != '\\' {
			builder.WriteRune(char)
		}
	}
	return &word{token: false, text: builder.String(), must: quote != 0}, nil
}

func (parser *parser) Read(token ...rune) (string, error) {
	builder := strings.Builder{}
	var tokens = map[rune]bool{}
	for _, t := range token {
		tokens[t] = true
	}
	for {
		char, err := parser.NextRune(true)
		if err != nil {
			if _, iseof := err.(*IOEOFError); !iseof {
				return "", err
			}
		}
		if tokens[char] || err != nil {
			return builder.String(), nil
		}
		builder.WriteRune(char)
	}
}

func (parser *parser) NextRune(ignoreBlank ...bool) (char rune, err error) {
	err = new(IOEOFError)
	for parser.Len() > 0 {
		parser.location++
		if char, _, err = parser.ReadRune(); len(ignoreBlank) == 0 || !ignoreBlank[0] || !blanks[char] {
			break
		}
	}
	return
}
