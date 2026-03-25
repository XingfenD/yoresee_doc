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
		return StatusServiceInternalError
	}
	return &Status{
		Code:    statusErr.Code,
		Message: msg,
	}
}

var (
	StatusSuccess = NewStatusErr(0, "success")

	StatusParamError       = NewStatusErr(40000, "invalid parameter")
	StatusTokenInvalid     = NewStatusErr(40001, "invalid token")
	StatusTokenExpired     = NewStatusErr(40002, "token expired")
	StatusPermissionDenied = NewStatusErr(40003, "permission denied")

	StatusUserNotFound      = NewStatusErr(40100, "user not found")
	StatusUserAlreadyExists = NewStatusErr(40101, "user already exists")
	StatusInvalidPassword   = NewStatusErr(40102, "invalid password")
	StatusDocumentNotFound  = NewStatusErr(40120, "document not found")

	StatusMembershipMetaNotFound        = NewStatusErr(40140, "membership meta not found")
	StatusInvalidMembershipType         = NewStatusErr(40141, "invalid membership type")
	StatusOrgNodeHasChildren            = NewStatusErr(40142, "org node has children")
	StatusOrgNodeCannotMoveToDescendant = NewStatusErr(40143, "org node cannot move to descendant")

	StatusKnowledgeBaseNotFound = NewStatusErr(40150, "knowledge base not found")

	StatusInvitationInvalid = NewStatusErr(40200, "invitation invalid")

	StatusPermissionDenied_DocumentRead = NewStatusErr(40300, "permission denied for reading document")

	StatusServiceInternalError = NewStatusErr(50000, "service internal error")
	StatusWriteDBError         = NewStatusErr(50001, "write db error")
	StatusReadDBError          = NewStatusErr(50002, "read db error")
	StatusMQNotInitialized     = NewStatusErr(50002, "message queue not initialized")
	StatusRedisNotInitialized  = NewStatusErr(50003, "redis not initialized")

	StatusInternalParamsError = NewStatusErr(50010, "invalid internal arguments")
	StatusTypeAssertFailed    = NewStatusErr(50011, "type assert failed")
)
