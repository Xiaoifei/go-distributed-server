package separation

import "time"

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusDone
	StatusFailed
)

type SeperationTask struct {
	ID        string     `json:"id"`
	Source    string     `json:"source"`
	Status    TaskStatus `json:"status"`
	Stems     []string   `json:"stems"`
	Error     string     `json:"error"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
