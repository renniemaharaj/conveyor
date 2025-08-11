package conveyor

// The worker's consumtion function
func (w *Worker) Consume(j *Job) error {
	return j.Consume(j.Param)
}
