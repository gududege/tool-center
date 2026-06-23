package task

// Repository 任务仓储接口
type Repository interface {
	Get(id string) (*Task, error)
	List() ([]*Task, error)
	Save(task *Task) error
	Delete(id string) error
}
