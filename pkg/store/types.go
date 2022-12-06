package store

type Task struct {
	ID           string
	Name         string
	Filename     string
	Info         string
	TargetFormat string
	TargetWidth  int
	TargetHeight int
	TargetFile   string
	BeginTime    string
	EndTime      string
	Command      string
	Output       []byte
	Retry        int
	DryRun       bool
}

type Store interface {
	Save(task *Task) string
	Get(id string) *Task
}
