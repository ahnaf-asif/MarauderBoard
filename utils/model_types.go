package utils

type Status string

const (
	Pending    Status = "pending"
	InProgress Status = "in_progress"
	Completed  Status = "completed"
	Cancelled  Status = "cancelled"
)
