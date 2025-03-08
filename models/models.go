package models

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Workspace{},
		&Team{},
		&Project{},
		&ChatGroup{},
		&ChatMessage{},
		&Task{},
		&Comment{},
	}
}
