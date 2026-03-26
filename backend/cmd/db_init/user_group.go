package main

import (
	"errors"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type seedUserGroup struct {
	Name        string
	Description string
	MemberUsers []string
}

func initializeUserGroupsInTx(tx *gorm.DB) error {
	logrus.Println("Creating user groups in transaction...")

	adminUser, err := getUserByUsername(tx, "admin")
	if err != nil {
		return err
	}
	user2, err := getUserByUsername(tx, "user2")
	if err != nil {
		return err
	}

	seedGroups := []seedUserGroup{
		{
			Name:        "Admin Group",
			Description: "Group for admin only",
			MemberUsers: []string{"admin"},
		},
		{
			Name:        "User2 Group",
			Description: "Group for user2 only",
			MemberUsers: []string{"user2"},
		},
		{
			Name:        "Admin & User2 Group",
			Description: "Group for admin and user2",
			MemberUsers: []string{"admin", "user2"},
		},
		{
			Name:        "Empty Group",
			Description: "Group with no members",
			MemberUsers: []string{},
		},
	}

	for _, seed := range seedGroups {
		group, err := findOrCreateGroup(tx, seed.Name, seed.Description, adminUser.ID)
		if err != nil {
			return err
		}
		for _, username := range seed.MemberUsers {
			var userID int64
			switch username {
			case "admin":
				userID = adminUser.ID
			case "user2":
				userID = user2.ID
			default:
				continue
			}
			if err := ensureMembership(tx, group.ID, userID); err != nil {
				return err
			}
		}
	}

	logrus.Println("User groups initialized successfully in transaction.")
	return nil
}

func getUserByUsername(tx *gorm.DB, username string) (*model.User, error) {
	var user model.User
	if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func findOrCreateGroup(tx *gorm.DB, name, description string, creatorID int64) (*model.UserGroupMeta, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, gorm.ErrInvalidValue
	}

	var group model.UserGroupMeta
	if err := tx.Where("name = ?", name).First(&group).Error; err == nil {
		return &group, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	group = model.UserGroupMeta{
		ExternalID:  utils.GenerateExternalID(utils.ExternalIDContextUserGroup),
		Name:        name,
		Description: strings.TrimSpace(description),
		CreatorID:   creatorID,
	}
	if err := tx.Create(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func ensureMembership(tx *gorm.DB, groupID, userID int64) error {
	var count int64
	if err := tx.Model(&model.MembershipRelation{}).
		Where("type = ? AND membership_id = ? AND user_id = ?", model.MembershipType_UserGroup, groupID, userID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return tx.Create(&model.MembershipRelation{
		Type:         model.MembershipType_UserGroup,
		MembershipID: groupID,
		UserID:       userID,
	}).Error
}
