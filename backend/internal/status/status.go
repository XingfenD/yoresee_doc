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
	StatusCommon = NewStatusErr(200, "success")

	StatusParamError = NewStatusErr(400, "invalid parameter")
)
