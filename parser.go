package jsobj

import (
	"strings"
)

type parser struct {
	reader *strings.Reader
}

func (obj *parser) Parse() (result interface{}, err error) {
	var w *word
	w, err = obj.ReadWord()
	if err == nil {
		//expected { / [ or a pure text
		if w.token {
			switch w.text {
			case "{":
				result, err = obj.ReadObj()
				break
			case "[":
				result, err = obj.ReadArray()
				break
			default:
				//other characters
				err = unexpectedWordError(w.text)
			}
		} else if _, test := obj.ReadWord(); test != nil {
			//text word can't have a follow-up word.
			result = w.getValue()
		} else {
			err = unexpectedWordError(w.text)
		}
	}
	if obj.reader.Len() > 0 {
		if w, err = obj.ReadWord(); err == nil {
			result = nil
			err = unexpectedWordError(w.text)
		}
	}
	return
}

func (obj *parser) ReadObj() (interface{}, error) {
	// readed first '{'
	//expect token '}', string
	var result, hasComma = map[string]interface{}{}, true
	for {
		w, err := obj.ReadWord()
		if err != nil {
			return nil, err
		}
		if w.token || !hasComma {
			if w.text == "}" {
				return result, nil
			}
			//if there is no comma, the object must be finished
			return nil, unexpectedWordError(w.text)
		}
		var key = w.text
		w, err = obj.ReadWord()
		if err != nil {
			return nil, err
		}
		//key : value
		if !w.token || w.text != ":" {
			return nil, unexpectedWordError(w.text)
		}
		w, err = obj.ReadWord()
		if err != nil {
			return nil, err
		}
		//an object, an array or just a text
		if w.token {
			if w.text == "{" {
				//an object
				result[key], err = obj.ReadObj()
			} else if w.text == "[" {
				//an array
				result[key], err = obj.ReadArray()
			} else {
				err = unexpectedWordError(w.text)
			}
			if err != nil {
				return nil, err
			}
		} else {
			//text
			result[key] = w.getValue()
		}
		//must be a token in comma or bracket
		w, err = obj.ReadWord()
		if err != nil {
			return nil, err
		}
		if !w.token {
			return nil, unexpectedWordError(w.text)
		}
		if hasComma = w.text == ","; !hasComma {
			//meet a none-comma token, unread
			obj.reader.UnreadRune()
		}
	}
}

func (obj *parser) ReadArray() (interface{}, error) {
	// readed first '['
	//expect token ']', string, '{', '['
	var result, hasComma = []interface{}{}, true
	for {
		var element interface{}
		w, err := obj.ReadWord()
		if err != nil {
			return nil, err
		}
		if w.token {
			//if no commas here, the token must be ]
			if w.text == "]" {
				return result, nil
			} else if w.text == "{" && hasComma {
				element, err = obj.ReadObj()
			} else if w.text == "[" && hasComma {
				element, err = obj.ReadArray()
			} else {
				err = unexpectedWordError(w.text)
			}
			if err != nil {
				return nil, err
			}
			result = append(result, element)
		} else {
			element = w.getValue()
			result = append(result, element)
		}
		//must be a token in comma or bracket
		w, err = obj.ReadWord()
		if err != nil {
			return nil, err
		}
		if !w.token {
			return nil, unexpectedWordError(w.text)
		}
		if hasComma = w.text == ","; !hasComma {
			//meet a none-comma token, unread
			obj.reader.UnreadRune()
		}
	}
}

func (obj *parser) ReadWord() (w *word, err error) {
	var char, quote, last rune
	char, err = obj.NextRune(true)
	if err != nil {
		return
	}
	if tokens[char] {
		w = &word{token: true, text: string(char)}
		return
	}
	var builder = strings.Builder{}
	if char == '\'' || char == '"' {
		quote = char
	} else {
		builder.WriteRune(char)
	}
	for char, err = obj.NextRune(); err == nil; char, err = obj.NextRune() {
		if last == '\\' && quote != 0 {
			builder.WriteRune(char)
			last = char
			continue
		}
		if char == quote {
			break
		}
		if quote == 0 {
			if blanks[char] {
				break
			}
			if tokens[char] {
				obj.reader.UnreadRune()
				break
			}
		}
		if last = char; last != '\\' {
			builder.WriteRune(char)
		}
	}
	return &word{token: false, text: builder.String(), must: quote != 0}, nil
}

func (obj *parser) NextRune(ignoreBlank ...bool) (char rune, err error) {
	err = new(IOEOFError)
	for obj.reader.Len() > 0 {
		if char, _, err = obj.reader.ReadRune(); len(ignoreBlank) == 0 || !ignoreBlank[0] || !blanks[char] {
			break
		}
	}
	return
}
