package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	Name        string
	Description string
	Status      string
	Workspace   Workspace `gorm:"foreignKey:WorkspaceId"`
	WorkspaceId int
	Teams       []*Team `gorm:"many2many:team_projects;constraints:OnDelete: SET NULL;"`
}

func AddNewProject(db *gorm.DB, project *Project) (*Project, error) {
	if err := db.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func GetProjectById(db *gorm.DB, id uint) (*Project, error) {
	project := &Project{}
	if err := db.Preload("Workspace").
		Preload("Teams").
		Preload("Teams.Leader").
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

func AddTeamToProject(db *gorm.DB, projectId uint, teamId uint) error {
	if err := db.Exec("INSERT INTO team_projects (team_id, project_id) VALUES (?, ?)", teamId, projectId).Error; err != nil {
		log.Println("Error adding team to project:", err)
		return err
	}

	team, _ := GetTeamById(db, teamId)
	project, _ := GetProjectById(db, projectId)

	for _, user := range team.Users {
		link := fmt.Sprintf("/workspaces/%d/projects/%d/teams", project.WorkspaceId, project.ID)
		notification := &Notification{
			UserId: user.ID,
			Title:  "Team Added to Project",
			Body:   "Your team " + team.Name + " has been added to the project " + project.Name,
			Seen:   false,
			Link:   &link,
		}
		if _, err := AddNotification(db, notification); err != nil {
			log.Println("Error adding notification:", err)
			return err
		}
	}
	return nil
}

func RemoveTeamFromProject(db *gorm.DB, projectId uint, teamId uint) error {
	if err := db.Exec("DELETE FROM team_projects WHERE team_id = ? AND project_id = ?", teamId, projectId).Error; err != nil {
		log.Println("Error removing team from project:", err)
		return err
	}

	team, _ := GetTeamById(db, teamId)
	project, _ := GetProjectById(db, projectId)

	for _, user := range team.Users {
		notification := &Notification{
			UserId: user.ID,
			Title:  "Team Removed from Project",
			Body:   "Your team " + team.Name + " has been removed from the project " + project.Name,
			Seen:   false,
			Link:   nil,
		}
		if _, err := AddNotification(db, notification); err != nil {
			log.Println("Error adding notification:", err)
		}
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
