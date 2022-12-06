package store

type Task struct {
	ID           string
	Name         string
	Filename     string
	TargetFormat string
	TargetFile   string
	Command      string
	Output       []byte
	Retry        int
	DryRun       bool
}

type Store interface {
	Save(task *Task) string
	Get(id string) *Task
}
