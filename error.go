package jsobj

type (
	//UnexpectedWordError An unexpected character appearing in a string
	UnexpectedWordError string

	//IOEOFError The cursor read from the data has arrived at the end of the stream.
	IOEOFError string
)

func unexpectedWordError(w string) error {
	var err = UnexpectedWordError(w)
	return &err
}

func (err *UnexpectedWordError) Error() string {
	return "js: unexpected character '" + string(*err) + "'"
}

func (err *IOEOFError) Error() string {
	return "io: EOF"
}
