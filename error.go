package jsobj

import "strconv"

type (
	//UnexpectedWordError An unexpected character appearing in a string
	UnexpectedWordError struct {
		text     string
		location int
	}

	//IOEOFError The cursor read from the data has arrived at the end of the stream.
	IOEOFError string
)

func unexpectedWordError(w string, position int) error {
	var err = &UnexpectedWordError{w, position - len([]rune(w))}
	return err
}

func isUnexpectedWordError(err error, w ...string) bool {
	if err == nil {
		return false
	}
	uwe, isuwe := err.(*UnexpectedWordError)
	return isuwe && (len(w) == 0 || (uwe.text == w[0]))
}

func (err *UnexpectedWordError) Error() string {
	return "js: unexpected character '" + err.text + "' at " + strconv.Itoa(err.location)
}

func (err *IOEOFError) Error() string {
	return "io: EOF"
}
