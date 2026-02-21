package status

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Status) Error() string {
	return s.Message
}

func NewStatusErr(code int, msg string) error {
	return &Status{
		Code:    code,
		Message: msg,
	}
}

func GenErrWithCustomMsg(err error, msg string) error {
	statusErr, ok := err.(*Status)
	if !ok {
		return &Status{
			Code:    500,
			Message: "internal server error",
		}
	}
	return &Status{
		Code:    statusErr.Code,
		Message: msg,
	}
}

var (
	StatusSuccess = NewStatusErr(20000, "success")

	StatusParamError   = NewStatusErr(40000, "invalid parameter")
	StatusTokenInvalid = NewStatusErr(40001, "invalid token")
	StatusTokenExpired = NewStatusErr(40002, "token expired")

	StatusUserNotFound      = NewStatusErr(40100, "user not found")
	StatusUserAlreadyExists = NewStatusErr(40101, "user already exists")
	StatusInvalidPassword   = NewStatusErr(40102, "invalid password")

	StatusInvitationInvalid = NewStatusErr(40200, "invitation invalid")

	StatusWriteDBError = NewStatusErr(50000, "write db error")
	StatusReadDBError  = NewStatusErr(50001, "read db error")
)
