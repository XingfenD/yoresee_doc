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

func (s *AuthService) QuerySideBarDisplay(userExternalID, scene string) ([]string, error) {
	if userExternalID == "" || scene == "" {
		return nil, status.StatusParamError
	}

	user, err := s.userRepo.GetByExternalID(userExternalID).Exec()
	if err != nil || user == nil {
		return nil, status.StatusUserNotFound
	}

	isAdmin := strings.EqualFold(user.Username, "admin")
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
