package auth_service

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/status"
)

const (
	SideBarSceneHome     = "home"
	SideBarSceneUserInfo = "user_info"
	SideBarSceneManage   = "manage"
)

var (
	homeTabs = []string{
		"home",
		"documents",
		"knowledge-base",
		"templates",
	}
	userInfoTabs = []string{
		"home",
		"user-center",
		"user-notifications",
		"user-invite",
		"user-security",
	}
	manageTabsForAdmin = []string{
		"home",
		"manage-user",
		"manage-user-group",
		"manage-organization",
		"manage-invite",
		"manage-security",
	}
	manageTabsForNormalUser = []string{
		"home",
	}
)

func (s *AuthService) QuerySideBarDisplay(scene string, isAdmin bool) ([]string, error) {
	if scene == "" {
		return nil, status.StatusParamError
	}
	switch scene {
	case SideBarSceneHome:
		return homeTabs, nil
	case SideBarSceneUserInfo:
		return userInfoTabs, nil
	case SideBarSceneManage:
		if isAdmin {
			return manageTabsForAdmin, nil
		}
		return manageTabsForNormalUser, nil
	default:
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "invalid scene")
	}
}

func (s *AuthService) IsAdmin(userExternalID string) (bool, error) {
	if userExternalID == "" {
		return false, status.StatusParamError
	}
	user, err := s.userRepo.GetByExternalID(userExternalID).Exec()
	if err != nil || user == nil {
		return false, status.StatusUserNotFound
	}
	return strings.EqualFold(user.Username, "admin"), nil
}
