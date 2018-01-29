package eimb

//Job ...
type Job interface {
	Execute()
}

//WithJob ...
type WithJob func()

//Execute ...
func (w WithJob) Execute() {
	w()
}

//JobQueue ...
var JobQueue = make(chan Job, 1)

//WorkerPool ..
type WorkerPool interface {
	Start()
}

//NewWorkerPool ..
func NewWorkerPool(size int) WorkerPool {
	return &workerPoolImpl{
		size: size,
	}
}

type workerPoolImpl struct {
	size int
}

func worker(id int) {
	for job := range JobQueue {
		job.Execute()
	}
}

func (w *workerPoolImpl) Start() {
	for i := 1; i <= w.size; i++ {
		go worker(i)
	}

}
