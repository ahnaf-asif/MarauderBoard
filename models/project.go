package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	Name        string
	Description string
	Status      string
	Workspace   Workspace `gorm:"foreignKey:WorkspaceId"`
	WorkspaceId int
	Teams       []*Team `gorm:"many2many:team_projects;"`
}

func AddNewProject(db *gorm.DB, project *Project) (*Project, error) {
	if err := db.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func GetProjectById(db *gorm.DB, id int) (*Project, error) {
	project := &Project{}
	if err := db.Preload("Workspace").
		Preload("Teams").
		First(project, id).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func GetProjectsByWorkspaceId(db *gorm.DB, workspaceId int) ([]*Project, error) {
	projects := []*Project{}
	if err := db.Preload("Workspace").
		Preload("Teams").
		Where("workspace_id = ?", workspaceId).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func GetProjectsByTeamId(db *gorm.DB, teamId int) ([]*Project, error) {
	projects := []*Project{}
	if err := db.Preload("Workspace").
		Preload("Teams").
		Joins("JOIN team_projects ON projects.id = team_projects.project_id").
		Where("team_projects.team_id = ?", teamId).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func UpdateProject(db *gorm.DB, project *Project) (*Project, error) {
	if err := db.Save(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func DeleteProject(db *gorm.DB, project *Project) error {
	if err := db.Delete(project).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProjectById(db *gorm.DB, id int) error {
	project := &Project{}
	if err := db.First(project, id).Error; err != nil {
		return err
	}
	if err := db.Delete(project).Error; err != nil {
		return err
	}
	return nil
}

func AddTeamToProject(db *gorm.DB, projectId int, teamId int) error {
	project := &Project{}
	team := &Team{}

	if err := db.First(project, projectId).Error; err != nil {
		return err
	}

	if err := db.First(team, teamId).Error; err != nil {
		return err
	}

	if err := db.Model(project).Association("Teams").Append(team); err != nil {
		return err
	}

	return nil
}

func RemoveTeamFromProject(db *gorm.DB, projectId int, teamId int) error {
	project := &Project{}
	team := &Team{}

	if err := db.First(project, projectId).Error; err != nil {
		return err
	}

	if err := db.First(team, teamId).Error; err != nil {
		return err
	}

	if err := db.Model(project).Association("Teams").Delete(team); err != nil {
		return err
	}

	return nil
}

func GetProjectByUserId(db *gorm.DB, userId int) ([]*Project, error) {
	projects := []*Project{}
	if err := db.Preload("Workspace").
		Preload("Teams").
		Joins("JOIN team_projects ON projects.id = team_projects.project_id").
		Joins("JOIN team_users ON team_projects.team_id = team_users.team_id").
		Where("team_users.user_id = ?", userId).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}
