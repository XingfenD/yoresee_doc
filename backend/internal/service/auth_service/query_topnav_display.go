package auth_service

import "github.com/XingfenD/yoresee_doc/internal/status"

const (
	TopNavMenuUserCenter   = "user-center"
	TopNavMenuSystemManage = "system-manage"
)

var (
	topNavMenusForAdmin = []string{
		TopNavMenuUserCenter,
		TopNavMenuSystemManage,
	}
	topNavMenusForNormalUser = []string{
		TopNavMenuUserCenter,
	}
)

func (s *AuthService) QueryTopNavDisplay(isAdmin bool) ([]string, error) {
	if isAdmin {
		return topNavMenusForAdmin, nil
	}
	if topNavMenusForNormalUser == nil {
		return nil, status.StatusParamError
	}
	return topNavMenusForNormalUser, nil
}
