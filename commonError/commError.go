package commonError

type CommError struct {
	errType int
	errStr  string
}

func (ce CommError) Error() string {
	return ce.errStr
}

func (ce CommError) GetType() int {
	return ce.errType
}

func NewCommErr(str string, t int) CommError {
	return CommError{errStr: str, errType: t}
}
