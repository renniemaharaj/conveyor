package conveyor

// Struct type for a worker
type Worker struct {
	canWork bool
}

// Start function of a worker
func (w *Worker) Start() {
	w.canWork = true
	for {
		if w.canWork {
			j := <-CONVEYOR_BELT
			w.Consume(&j)
			continue
		}

		break
	}
}

// Safe stop function of a worker
func (w *Worker) Stop() {
	w.canWork = false
}
