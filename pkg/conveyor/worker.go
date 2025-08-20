package conveyor

// Struct type for a worker
type Worker struct {
	canWork bool
	B       ConveyorBelt // The worker's assigned channel
}

// Creates and returns a worker with an assigned channel
func NewWorker(b *ConveyorBelt, id int) *Worker {
	return &Worker{canWork: false, B: *b}
}

// Start function of a worker
func (w *Worker) Start() {
	w.canWork = true
	for {
		if w.canWork {
			j := <-w.B.C // worker will take from it's assigned channel
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
