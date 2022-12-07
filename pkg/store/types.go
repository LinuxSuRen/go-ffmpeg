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
	Metadata     map[string]string
	BeginTime    string
	EndTime      string
	AudioBitrate int
	Command      string
	Output       string
	ErrOutput    string
	Retry        int
	DryRun       bool
}

type Store interface {
	Save(task *Task) string
	Get(id string) *Task
}
