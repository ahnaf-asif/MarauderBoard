package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model

	Name        string
	Description string
	Status      string
	Project     Project `gorm:"foreignKey:ProjectId"`
	ProjectId   int
	Assignee    User `gorm:"foreignKey:AssigneeId"`
	AssigneeId  int
	Reporter    User `gorm:"foreignKey:ReporterId"`
	ReporterId  int
	Team        Team `gorm:"foreignKey:TeamId"`
	TeamId      int
	Comments    []*Comment `gorm:"foreignKey:TaskId"`
}

func AddNewTask(db *gorm.DB, task *Task) (*Task, error) {
	if err := db.Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func GetTaskById(db *gorm.DB, id int) (*Task, error) {
	task := &Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		First(task, id).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func GetTasksByProjectId(db *gorm.DB, projectId int) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("project_id = ?", projectId).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByTeamId(db *gorm.DB, teamId int) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("team_id = ?", teamId).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func UpdateTask(db *gorm.DB, task *Task) (*Task, error) {
	if err := db.Save(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func DeleteTask(db *gorm.DB, task *Task) error {
	if err := db.Delete(task).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTaskById(db *gorm.DB, id int) error {
	task := &Task{}
	if err := db.First(task, id).Error; err != nil {
		return err
	}
	if err := db.Delete(task).Error; err != nil {
		return err
	}
	return nil
}

func AddCommentToTask(db *gorm.DB, taskId int, comment *Comment) error {
	task := &Task{}
	if err := db.First(task, taskId).Error; err != nil {
		return err
	}

	if err := db.Model(task).Association("Comments").Append(comment); err != nil {
		return err
	}

	return nil
}

func RemoveCommentFromTask(db *gorm.DB, taskId int, commentId int) error {
	task := &Task{}
	comment := &Comment{}

	if err := db.First(task, taskId).Error; err != nil {
		return err
	}

	if err := db.First(comment, commentId).Error; err != nil {
		return err
	}

	if err := db.Model(task).Association("Comments").Delete(comment); err != nil {
		return err
	}

	return nil
}

func GetTasksByAssigneeId(db *gorm.DB, assigneeId int) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("assignee_id = ?", assigneeId).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByReporterId(db *gorm.DB, reporterId int) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("reporter_id = ?", reporterId).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByStatus(db *gorm.DB, status string) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("status = ?", status).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByAssigneeAndStatus(db *gorm.DB, assigneeId int, status string) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("assignee_id = ? AND status = ?", assigneeId, status).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByReporterAndStatus(db *gorm.DB, reporterId int, status string) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("reporter_id = ? AND status = ?", reporterId, status).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByTeamAndStatus(db *gorm.DB, teamId int, status string) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("team_id = ? AND status = ?", teamId, status).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByAssigneeAndTeam(db *gorm.DB, assigneeId int, teamId int) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("assignee_id = ? AND team_id = ?", assigneeId, teamId).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByReporterAndTeam(db *gorm.DB, reporterId int, teamId int) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("reporter_id = ? AND team_id = ?", reporterId, teamId).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByAssigneeAndReporter(db *gorm.DB, assigneeId int, reporterId int) ([]*Task, error) {
	tasks := []*Task{}
	if err := db.Preload("Project").
		Preload("Assignee").
		Preload("Reporter").
		Preload("Team").
		Preload("Comments").
		Where("assignee_id = ? AND reporter_id = ?", assigneeId, reporterId).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
