package errno

import "fmt"

type ErrNo struct {
	ErrCode int
	ErrMsg  string
	Data    string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int, msg string) ErrNo {
	return ErrNo{code, msg, ""}
}

func (e ErrNo) Copy() ErrNo {
	newErrNo := NewErrNo(e.ErrCode, e.ErrMsg)
	return newErrNo
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	newErrNo := e.Copy()
	newErrNo.ErrMsg = msg
	return newErrNo
}

func (e ErrNo) WithData(data string) ErrNo {
	newErrNo := e.Copy()
	newErrNo.Data = data
	return newErrNo
}
