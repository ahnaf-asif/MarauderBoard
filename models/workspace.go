package models

import (
	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model

	Name            string
	Description     string
	Administrator   *User `gorm:"constraint:OnDelete:SET NULL;foreignKey:AdministratorId;references:ID"`
	AdministratorId *uint
	Users           []*User    `gorm:"many2many:user_workspaces;constraint:OnDelete:SET NULL;"`
	Teams           []*Team    `gorm:"foreignKey:WorkspaceId;constraint:OnDelete:SET NULL;"`
	Projects        []*Project `gorm:"foreignKey:WorkspaceId;constraint:OnDelete:SET NULL;"`
	ChatGroup       *ChatGroup `gorm:"constraing:foreignKey:ChatGroupId;OnDelete:Set NULL;"`
	ChatGroupId     *uint
}

func getAllWorkspacesByAdminId(db *gorm.DB, adminId uint) ([]*Workspace, error) {
	var workspaces []*Workspace
	if err := db.Where("administrator_id = ?", adminId).Find(&workspaces).Error; err != nil {
		return nil, err
	}
	return workspaces, nil
}

func GetAllWorkspacesByUserId(db *gorm.DB, userId uint) ([]*Workspace, error) {
	var workspaces []*Workspace
	if err := db.Preload("Administrator").
		Preload("Teams").
		Preload("Teams.Users").
		Preload("Projects").
		Preload("ChatGroup").
		Find(&workspaces).Error; err != nil {
		return nil, err
	}

	var filteredWorkspaces []*Workspace
	for _, workspace := range workspaces {
		if *workspace.AdministratorId == userId {
			filteredWorkspaces = append(filteredWorkspaces, workspace)
			continue
		}
		for _, team := range workspace.Teams {
			for _, user := range team.Users {
				if user.ID == userId {
					filteredWorkspaces = append(filteredWorkspaces, workspace)
					break
				}
			}
		}
	}
	return filteredWorkspaces, nil
}

func AddWorkspace(db *gorm.DB, workspace *Workspace) (*Workspace, error) {
	if err := db.Create(workspace).Error; err != nil {
		return nil, err
	}
	return workspace, nil
}

func GetWorkspaceById(db *gorm.DB, id uint) (*Workspace, error) {
	var workspace Workspace
	if err := db.Preload("Administrator").
		Preload("Users").
		Preload("Teams").
		Preload("Projects").
		Preload("ChatGroup").
		Preload("ChatGroup.Messages").
		Preload("ChatGroup.Users").
		Preload("ChatGroup.Messages.User").
		First(&workspace, id).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}
