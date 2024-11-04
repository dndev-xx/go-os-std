package internal

type IRunner interface {
	Run(tasks []Task, n, m int) error
	ShouldExecutionTasks(n, m int) (*bool, error)
	SendTasks(tasks []Task)
	StartWorkers()
}
