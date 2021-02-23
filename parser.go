package jsobj

import (
	"strings"
)

type parser strings.Reader

func (parser *parser) Len() int {
	return (*strings.Reader)(parser).Len()
}

func (parser *parser) UnreadRune() error {
	return (*strings.Reader)(parser).UnreadRune()
}

func (parser *parser) Parse() (result interface{}, err error) {
	var w *word
	w, err = parser.ReadWord()
	if err == nil {
		//expected { / [ or a pure text
		if w.token {
			switch w.text {
			case "{":
				result, err = parser.ReadObj()
				break
			case "[":
				result, err = parser.ReadArray()
				break
			default:
				//other characters
				err = unexpectedWordError(w.text)
			}
		} else if _, test := parser.ReadWord(); test != nil {
			//text word can't have a follow-up word.
			result = w.getValue()
		} else {
			err = unexpectedWordError(w.text)
		}
	}
	if parser.Len() > 0 {
		if w, err = parser.ReadWord(); err == nil {
			result = nil
			err = unexpectedWordError(w.text)
		}
	}
	return
}

func (parser *parser) ReadObj() (interface{}, error) {
	// readed first '{'
	//expect token '}', string
	var result, hasComma = map[string]interface{}{}, true
	for {
		w, err := parser.ReadWord()
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
		w, err = parser.ReadWord()
		if err != nil {
			return nil, err
		}
		//key : value
		if !w.token || w.text != ":" {
			return nil, unexpectedWordError(w.text)
		}
		w, err = parser.ReadWord()
		if err != nil {
			return nil, err
		}
		//an object, an array or just a text
		if w.token {
			if w.text == "{" {
				//an object
				result[key], err = parser.ReadObj()
			} else if w.text == "[" {
				//an array
				result[key], err = parser.ReadArray()
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
		w, err = parser.ReadWord()
		if err != nil {
			return nil, err
		}
		if !w.token {
			return nil, unexpectedWordError(w.text)
		}
		if hasComma = w.text == ","; !hasComma {
			//meet a none-comma token, unread
			parser.UnreadRune()
		}
	}
}

func (parser *parser) ReadArray() (interface{}, error) {
	// readed first '['
	//expect token ']', string, '{', '['
	var result, hasComma = []interface{}{}, true
	for {
		var element interface{}
		w, err := parser.ReadWord()
		if err != nil {
			return nil, err
		}
		if w.token {
			//if no commas here, the token must be ]
			if w.text == "]" {
				return result, nil
			} else if w.text == "{" && hasComma {
				element, err = parser.ReadObj()
			} else if w.text == "[" && hasComma {
				element, err = parser.ReadArray()
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
		w, err = parser.ReadWord()
		if err != nil {
			return nil, err
		}
		if !w.token {
			return nil, unexpectedWordError(w.text)
		}
		if hasComma = w.text == ","; !hasComma {
			//meet a none-comma token, unread
			parser.UnreadRune()
		}
	}
}

func (parser *parser) ReadWord() (w *word, err error) {
	var char, quote, last rune
	char, err = parser.NextRune(true)
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
	for char, err = parser.NextRune(); err == nil; char, err = parser.NextRune() {
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

func (parser *parser) NextRune(ignoreBlank ...bool) (char rune, err error) {
	err = new(IOEOFError)
	for parser.Len() > 0 {
		if char, _, err = (*strings.Reader)(parser).ReadRune(); len(ignoreBlank) == 0 || !ignoreBlank[0] || !blanks[char] {
			break
		}
	}
	return
}
